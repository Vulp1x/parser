package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/dbtx"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/pb/instaproxy"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
)

type ParseFollowersHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
}

func (h *ParseFollowersHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	datasetID, bloggerUsername, err := h.parseTaskKey(ctx, task)
	if err != nil {
		return err
	}

	db := h.dbTxF(ctx)
	q := dbmodel.New(db)
	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find datatset with id '%s': %v", datasetID, err)
	}

	if dataset.FollowersCount <= 0 {
		return fmt.Errorf("%w: got %d followers count for dataset6 expected at least 0", pgqueue.ErrMustCancelTask, dataset.FollowersCount)
	}

	targetsResp, err := h.cli.ParseFollowers(ctx, &instaproxy.ParseFollowersRequest{
		UserName:       bloggerUsername,
		FollowersCount: dataset.FollowersCount,
	})
	if err != nil {
		return fmt.Errorf("failed to parse targets from media '%s': %v", bloggerUsername, err)
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)
	q = dbmodel.New(tx)

	fullTargetParams := domain.FullUserFromProto(targetsResp.Blogger).ToSaveFullTargetParams(datasetID)
	err = q.SaveFullTarget(ctx, fullTargetParams)
	if err != nil {
		return fmt.Errorf("failed to save full user: %v with params %v", err, fullTargetParams)
	}

	err = q.SaveFakeMedia(ctx, dbmodel.SaveFakeMediaParams{
		Pk: targetsResp.Blogger.Pk, DatasetID: datasetID, Caption: bloggerUsername,
	})

	targets := domain.ShortUsersFromProto(targetsResp.GetTargets())

	count, err := q.SaveTargetUsers(ctx, targets.ToSaveTargetsParams(targetsResp.Blogger.Pk, datasetID))
	if err != nil {
		return fmt.Errorf("failed to save targets: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	logger.Infof(ctx, "saved %d/%d targets", count, len(targets))

	return nil
}

func (h *ParseFollowersHandler) parseTaskKey(ctx context.Context, task pgqueue.Task) (uuid.UUID, string, error) {
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
