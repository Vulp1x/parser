package datasets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/dbtx"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/workers"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
	"github.com/jackc/pgx/v4"
)

const maxRetriesCount = 5

// ErrNoBlogers не смогли найти таску
var ErrNoBlogers = errors.New("no initial bloggers")

// ErrNoReadyBots не смогли найти таску
var ErrNoReadyBots = errors.New("all bots are blocked")

func (s *Store) FindSimilarBloggers(ctx context.Context, datasetID uuid.UUID) (domain.DatasetWithBloggers, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to start transaction: %v", err)
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.DatasetWithBloggers{}, ErrDatasetNotFound
		}

		return domain.DatasetWithBloggers{}, err
	}

	if dataset.Status != dbmodel.DraftDatasetStatus {
		return domain.DatasetWithBloggers{}, fmt.Errorf("%w: ожидали статус драфт (%d), а получили %d",
			ErrDatasetInvalidStatus, dbmodel.DraftDatasetStatus, dataset.Status,
		)
	}

	initialBloggers, err := q.FindInitialBloggersForDataset(ctx, datasetID)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to find initial bloggers for dataset: %v", err)
	}

	if len(initialBloggers) == 0 {
		return domain.DatasetWithBloggers{}, ErrNoBlogers
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.FindingSimilarStarted, ID: datasetID})
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to update dataset status: %v", err)
	}

	dataset.Status = dbmodel.FindingSimilarStarted

	logger.Infof(ctx, "adding task for %d initial bloggers", len(initialBloggers))

	var tasks = make([]pgqueue.Task, len(initialBloggers))

	for i, initialBlogger := range initialBloggers {
		bloggerBytes, err := json.Marshal(initialBlogger)
		if err != nil {
			return domain.DatasetWithBloggers{}, fmt.Errorf("failed to marshal blogger %s: %v", initialBlogger.Username, err)
		}

		tasks[i] = pgqueue.Task{
			Kind:        workers.FindSimilarBloggersTaskKind,
			Payload:     bloggerBytes,
			ExternalKey: fmt.Sprintf("%s::%s", dataset.ID, initialBlogger.Username),
		}
	}

	err = s.queue.PushTasksTx(ctx, tx, tasks)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to push tasks to queue: %v", err)
	}

	err = s.queue.PushTaskTx(ctx, tx, pgqueue.Task{
		Kind:        workers.TransitToSimilarFoundTaskKind,
		Payload:     workers.EmptyPayload,
		ExternalKey: datasetID.String(),
	}, pgqueue.WithDelay(time.Duration(len(initialBloggers))*time.Second))
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to push TransitToSimilarFound task to queue: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return domain.NewDatasetWithBloggers(dataset, initialBloggers), nil
}
