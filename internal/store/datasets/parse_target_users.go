package datasets

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/dbtx"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/pkg/logger"
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

	if dataset.Status != dbmodel.DraftDatasetStatus {
		return domain.DatasetWithBloggers{}, fmt.Errorf("%w: ожидали статус драфт (%d), а получили %d",
			ErrDatasetInvalidStatus, dbmodel.DraftDatasetStatus, dataset.Status,
		)
	}

	bloggers, err := q.FindInitialBloggersForDataset(ctx, datasetID)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to find initial bloggers for dataset: %v", err)
	}

	if len(bloggers) == 0 {
		return domain.DatasetWithBloggers{}, ErrNoBlogers
	}

	countAvailableBots, err := q.CountAvailableBots(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to countAvailableBots available bots")
	}

	if countAvailableBots == 0 {
		return domain.DatasetWithBloggers{}, ErrNoReadyBots
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.FindingSimilarStarted, ID: datasetID})
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to update dataset status: %v", err)
	}

	dataset.Status = dbmodel.ReadyForParsingDatasetStatus

	logger.Infof(ctx, "adding task for %d bloggers, expected maximum %d bots (available %d)", len(bloggers), botsPerDataset, countAvailableBots)

	err = s.findSimilarService.AddParseTargetsTask(ctx, dataset.ID, bloggers, botsPerDataset)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to add task: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return domain.NewDatasetWithBloggers(dataset, bloggers), nil
}
