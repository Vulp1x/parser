package datasets

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/dbtx"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/workers"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/inst-api/parser/pkg/pgqueue"
	"github.com/jackc/pgx/v4"
)

func (s *Store) ParseTargetUsers(ctx context.Context, datasetID uuid.UUID) (domain.DatasetWithBloggers, error) {
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

	if dataset.Status != dbmodel.ReadyForParsingDatasetStatus && dataset.Status != dbmodel.DraftDatasetStatus {
		return domain.DatasetWithBloggers{}, fmt.Errorf("%w: ожидали статус готов к парсингу (%d) или драфт (%d), а получили %d",
			ErrDatasetInvalidStatus, dbmodel.ReadyForParsingDatasetStatus, dbmodel.DraftDatasetStatus, dataset.Status,
		)
	}

	bloggers, err := q.FindBloggersForParsing(ctx, datasetID)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to find initial bloggers for dataset: %v", err)
	}

	if len(bloggers) == 0 {
		return domain.DatasetWithBloggers{}, ErrNoBlogers
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.ParsingTargetsStartedDatasetStatus, ID: datasetID})
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to update dataset status: %v", err)
	}

	dataset.Status = dbmodel.ParsingTargetsStartedDatasetStatus

	logger.Infof(ctx, "adding tasks for %d bloggers", len(bloggers))

	tasks := make([]pgqueue.Task, len(bloggers))
	for i, blogger := range bloggers {
		tasks[i] = pgqueue.Task{
			Kind: workers.ParseBloggersMediaTaskKind, ExternalKey: fmt.Sprintf("%s::%s", datasetID, blogger.Username), Payload: workers.EmptyPayload,
		}
	}

	err = s.queue.PushTasksTx(ctx, tx, tasks)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to push tasks to queue: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return domain.NewDatasetWithBloggers(dataset, bloggers), nil
}
