package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/pb/instaproxy"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
)

type ParseFullUsersHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
}

func (h *ParseFullUsersHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	datasetID, bloggerUsername, err := h.parseTaskKey(ctx, task)
	if err != nil {
		return err
	}

	q := dbmodel.New(h.dbTxF(ctx))
	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find datatset with id '%s': %v", datasetID, err)
	}

	if dataset.FollowersCount <= 0 {
		return fmt.Errorf("%w: got %d followers count for dataset6 expected at least 0", pgqueue.ErrMustCancelTask, dataset.FollowersCount)
	}

	fullUserResp, err := h.cli.GetFullUser(ctx, &instaproxy.GetUserRequest{Username: bloggerUsername})
	if err != nil {
		return fmt.Errorf("failed to get full user info from media '%s': %v", bloggerUsername, err)
	}

	fullTargetParams := domain.FullUserFromProto(fullUserResp).ToSaveFullTargetParams(datasetID)
	err = q.SaveFullTarget(ctx, fullTargetParams)
	if err != nil {
		return fmt.Errorf("failed to save full user: %v with params %v", err, fullTargetParams)
	}

	logger.InfoKV(ctx, "saved full target")

	return nil
}

func (h *ParseFullUsersHandler) parseTaskKey(ctx context.Context, task pgqueue.Task) (uuid.UUID, string, error) {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) != 2 {
		return uuid.UUID{}, "", fmt.Errorf(
			`%w: expected external key in format "<dataset id>::<blogger username>", got '%s'`,
			pgqueue.ErrMustCancelTask, task.ExternalKey)
	}

	datasetID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return uuid.UUID{}, "", fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	bloggerUsername := externalKeyParts[1]

	return datasetID, bloggerUsername, nil
}
