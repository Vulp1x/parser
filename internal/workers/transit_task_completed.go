package workers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
)

type TransitToTaskCompleteddHandler struct {
	dbTxF dbmodel.DBTXFunc
	queue *pgqueue.Queue
}

// HandleTask переводит датасет в конечный статус как только все блогеры распаршены
func (h *TransitToTaskCompleteddHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	datasetID, err := uuid.Parse(task.ExternalKey)
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, task.ExternalKey, err)
	}

	taskKind, err := strconv.ParseInt(string(task.Payload), 10, 16)
	if err != nil {
		return fmt.Errorf("%w: failed to parse task kind from %s: %v", pgqueue.ErrMustCancelTask, string(task.Payload), err)
	}

	q := dbmodel.New(h.dbTxF(ctx))

	notReadyBloggers, err := q.FindNotReadyBloggers(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find not ready bloggers: %v", err)
	}

	if len(notReadyBloggers) != 0 {
		if err = h.queue.RetryTasks(ctx, int16(taskKind), 10, 1000); err != nil {
			return fmt.Errorf("failed to retry tasks for dataset %s with %d not ready bloggers: %v: %v", datasetID,
				len(notReadyBloggers), domain.DatasetWithBloggers{Bloggers: notReadyBloggers}.Usernames(), err)
		}

		return fmt.Errorf("dataset %s still has %d not ready bloggers: %v", datasetID,
			len(notReadyBloggers), domain.DatasetWithBloggers{Bloggers: notReadyBloggers}.Usernames())
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.ParsingDoneDatasetStatus, ID: datasetID})
	if err != nil {
		return fmt.Errorf("failed to update dataset status: %v", err)
	}

	return nil
}
