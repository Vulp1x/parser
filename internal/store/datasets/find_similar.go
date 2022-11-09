package datasets

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/dbtx"
	"github.com/inst-api/parser/internal/domain"
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
			ErrTaskInvalidStatus, dbmodel.DraftDatasetStatus, dataset.Status,
		)
	}

	bloggers, err := q.FindInitialBloggersForDataset(ctx, datasetID)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to find initial bloggers for dataset: %v", err)
	}

	if len(bloggers) == 0 {
		return domain.DatasetWithBloggers{}, ErrNoBlogers
	}

	count, err := q.CountAvailableBots(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to count available bots")
	}

	if count == 0 {
		return domain.DatasetWithBloggers{}, ErrNoReadyBots
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.FindingSimilarStarted, ID: datasetID})
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to update dataset status: %v", err)
	}

	dataset.Status = dbmodel.FindingSimilarStarted

	err = s.findSimilarQueue.AddFindSimilarTask(ctx, dataset.ID, bloggers, 5)
	if err != nil {

	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return domain.NewDatasetWithBloggers(dataset, bloggers), nil
}

func findAvailableBots(ctx context.Context, q *dbmodel.Queries, i int) (dbmodel.Bot, error) {
	return dbmodel.Bot{}, nil
}
