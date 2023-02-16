package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/pb/instaproxy"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SimilarBloggersHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
}

func (s *SimilarBloggersHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) == 0 {
		return fmt.Errorf("%w: expected '::' in external key after dataset id in '%s'", pgqueue.ErrMustCancelTask, task.ExternalKey)
	}

	datasetID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	var blogger dbmodel.Blogger

	err = json.Unmarshal(task.Payload, &blogger)
	if err != nil {
		return fmt.Errorf("failed to unmarshal task: %v", err)
	}

	q := dbmodel.New(s.dbTxF(ctx))

	similarBloggersResp, err := s.cli.FindSimilarBloggers(ctx, &instaproxy.SimilarBloggersRequest{Username: blogger.Username})
	if err != nil {
		logger.Errorf(ctx, "failed to find similar bloggers: %v", err)

		if statusCode, ok := status.FromError(err); ok && statusCode.Code() == codes.NotFound {
			logger.Warnf(ctx, "blogger %s wasn't found: %v", blogger.Username, err)
			err = q.SetBloggerIsParsed(ctx, dbmodel.SetBloggerIsParsedParams{IsCorrect: false, ID: blogger.ID})
			if err != nil {
				return fmt.Errorf("failed to set bot incorrect (%s): %v", blogger.ID, err)
			}
		}

		return fmt.Errorf("failed to find similar bloggers: %v", err)
	}

	initialBlogger := similarBloggersResp.GetInitialBlogger()
	err = q.UpdateBlogger(ctx, dbmodel.UpdateBloggerParams{
		UserID:     initialBlogger.Pk,
		ParsedAt:   domain.Ptr(time.Now()),
		IsCorrect:  true,
		IsPrivate:  initialBlogger.IsPrivate,
		IsVerified: initialBlogger.IsVerified,
		ID:         blogger.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to update initial blogger (%s): %v", blogger.ID, err)
	}

	fullTargetParams := domain.FullUserFromProto(initialBlogger).ToSaveFullTargetParams(datasetID)
	err = q.SaveFullTarget(ctx, fullTargetParams)
	if err != nil {
		return fmt.Errorf("failed to save full user: %v with params %v", err, fullTargetParams)
	}

	domainparsedBloggers := domain.ShortUsersFromProto(similarBloggersResp.SimilarBloggers)

	var savedBloggersCount int

	saveBloggersBatch := q.SaveBloggers(ctx, domainparsedBloggers.ToSaveBloggersParmas(datasetID))
	saveBloggersBatch.Exec(func(j int, err error) {
		if err != nil {
			logger.Errorf(ctx, "failed to save %d parsed blogger from blogger (%s): %v", j, blogger.Username, err)
			return
		}

		savedBloggersCount++
	})

	logger.Infof(ctx, "saved %d/%d bloggers from initial blogger '%s'", savedBloggersCount, len(domainparsedBloggers), blogger.Username)

	err = q.MarkBloggerAsSimilarAccountsFound(ctx, blogger.ID)
	if err != nil {
		return fmt.Errorf("failed to mark blogger (%s) as ready for target's parsing : %v", blogger.ID, err)
	}

	return nil
}
