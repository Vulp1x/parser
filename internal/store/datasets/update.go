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

func (s Store) UpdateDataset(ctx context.Context, datasetID uuid.UUID, title *string, phoneCode *int16, originalAccounts []string) (domain.DatasetWithBloggers, error) {
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
	}

	if title != nil || phoneCode != nil {
		params := dbmodel.UpdateDatasetParams{PhoneCode: dataset.PhoneCode, Title: dataset.Title, ID: datasetID}

		if title != nil {
			params.Title = *title
		}

		if phoneCode != nil {
			params.PhoneCode = phoneCode
		}

		dataset, err = q.UpdateDataset(ctx, params)
		if err != nil {
			return domain.DatasetWithBloggers{}, err
		}
	}

	err = s.insertInitialBloggers(ctx, q, datasetID, originalAccounts)
	if err != nil {
		return domain.DatasetWithBloggers{}, err
	}

	bloggers, err := q.FindBloggersForDataset(ctx, datasetID)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to find bloggers: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return domain.NewDatasetWithBloggers(dataset, bloggers), nil
}

func (s Store) insertInitialBloggers(ctx context.Context, q *dbmodel.Queries, datasetID uuid.UUID, originalAccounts []string) error {
	if len(originalAccounts) == 0 {
		logger.Info(ctx, "no accounts to insert")
		return nil
	}

	tag, err := q.DeleteBloggersPerDataset(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to delete previous bloggers")
	}

	logger.Infof(ctx, "deleted %d bloggers: %v", tag.RowsAffected())

	var initialBloggersParams = make([]dbmodel.InsertInitialBloggersParams, len(originalAccounts))

	for i, account := range originalAccounts {
		initialBloggersParams[i] = dbmodel.InsertInitialBloggersParams{
			DatasetID: datasetID,
			Username:  account,
			UserID:    -1,
			IsInitial: true,
		}
	}

	insertedBloggersCount, err := q.InsertInitialBloggers(ctx, initialBloggersParams)
	if err != nil {
		return fmt.Errorf("failed to insert new initial bloggers: %v", err)
	}

	logger.Infof(ctx, "inserted %d initial bloggers", insertedBloggersCount)

	return nil
}
