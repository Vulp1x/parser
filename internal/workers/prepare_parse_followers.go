package workers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/pb/instaproxy"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PrepareParseFollowersHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
	queue *pgqueue.Queue
}

func (h *PrepareParseFollowersHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	taskKey, err := h.parseTaskKey(ctx, task)
	if err != nil {
		return err
	}

	ctx = logger.WithFields(ctx, logger.Fields{"blogger": taskKey.bloggerUsername, "dataset_id": taskKey.datasetID})

	db := h.dbTxF(ctx)
	q := dbmodel.New(db)
	dataset, err := q.GetDatasetByID(ctx, taskKey.datasetID)
	if err != nil {
		return fmt.Errorf("failed to find datatset with id '%s': %v", taskKey.datasetID, err)
	}

	if dataset.FollowersCount <= 0 {
		return fmt.Errorf("%w: got %d followers count for dataset6 expected at least 0", pgqueue.ErrMustCancelTask, dataset.FollowersCount)
	}

	blogger, err := h.cli.GetFullUser(ctx, &instaproxy.GetUserRequest{Username: taskKey.bloggerUsername})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			logger.WarnKV(ctx, "blogger not found, going to set invalid status")

			if err2 := q.SetBloggerStatusToInvalid(ctx, dbmodel.SetBloggerStatusToInvalidParams{
				Username:  taskKey.bloggerUsername,
				DatasetID: taskKey.datasetID,
			}); err2 != nil {
				return fmt.Errorf("failed to set blogger status to invalid: %v", err2)
			}

			return fmt.Errorf("%w: blogger not found: %v", pgqueue.ErrMustCancelTask, err)
		}

		return fmt.Errorf("failed to get full info about blogger: %v", err)
	}

	fullTargetParams := domain.FullUserFromProto(blogger).ToSaveFullTargetParams(taskKey.datasetID)
	err = q.SaveFullTarget(ctx, fullTargetParams)
	if err != nil {
		return fmt.Errorf("failed to save full user: %v with params %+v", err, fullTargetParams)
	}

	err = q.SaveFakeMedia(ctx, dbmodel.SaveFakeMediaParams{
		Pk: blogger.Pk, ID: strconv.FormatInt(blogger.Pk, 10),
		DatasetID: taskKey.datasetID, Caption: taskKey.bloggerUsername,
	})
	if err != nil {
		return fmt.Errorf("failed to save fake media: %v", err)
	}

	newTaskKey := parseFollowersTaskKey{
		datasetID:   dataset.ID,
		bloggerPk:   blogger.Pk,
		orderNumber: 0,
	}

	err = h.queue.PushTask(ctx, pgqueue.Task{
		Kind:        ParseFollowersTaskKind,
		Payload:     EmptyPayload,
		ExternalKey: newTaskKey.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to push new task to queue: %v", err)
	}

	return nil
}

func (h *PrepareParseFollowersHandler) parseTaskKey(ctx context.Context, task pgqueue.Task) (preparesToParseFollowersTaskKey, error) {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) != 2 {
		return preparesToParseFollowersTaskKey{}, fmt.Errorf(
			`%w: expected external key in format "<dataset id>::<blogger username>", got '%s'`,
			pgqueue.ErrMustCancelTask, task.ExternalKey)
	}

	datasetID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return preparesToParseFollowersTaskKey{}, fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	return preparesToParseFollowersTaskKey{datasetID: datasetID, bloggerUsername: externalKeyParts[1]}, nil
}

type preparesToParseFollowersTaskKey struct {
	datasetID       uuid.UUID
	bloggerUsername string
}

func (p preparesToParseFollowersTaskKey) String() string {
	return fmt.Sprintf("%s::%s", p.datasetID.String(), p.bloggerUsername)
}
