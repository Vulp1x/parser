package workers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
)

type TransitToSimilarFoundHandler struct {
	dbTxF dbmodel.DBTXFunc
	queue *pgqueue.Queue
}

func (h *TransitToSimilarFoundHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	datasetID, err := uuid.Parse(task.ExternalKey)
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, task.ExternalKey, err)
	}

	q := dbmodel.New(h.dbTxF(ctx))

	notReadyBloggers, err := q.FindNotReadyBloggers(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find not ready bloggers: %v", err)
	}

	if len(notReadyBloggers) != 0 {
		if err = h.queue.RetryTasks(ctx, FindSimilarBloggersTaskKind, 10, 1000); err != nil {
			return fmt.Errorf("failed to retry tasks for dataset %s with %d not ready bloggers: %v: %v", datasetID,
				len(notReadyBloggers), domain.DatasetWithBloggers{Bloggers: notReadyBloggers}.Usernames(), err)
		}

		return fmt.Errorf("dataset %s still has %d not ready bloggers: %v", datasetID,
			len(notReadyBloggers), domain.DatasetWithBloggers{Bloggers: notReadyBloggers}.Usernames())
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.ReadyForParsingDatasetStatus, ID: datasetID})
	if err != nil {
		return fmt.Errorf("failed to update dataset status: %v", err)
	}

	return nil
}
