package datasets

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/jackc/pgx/v4"
)

var ErrDatasetNotFound = errors.New("dataset not found")

func (s Store) GetDataset(ctx context.Context, datasetID uuid.UUID) (domain.DatasetWithBloggers, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.DatasetWithBloggers{}, ErrDatasetNotFound
		}

		return domain.DatasetWithBloggers{}, err
	}

	bloggers, err := q.FindBloggersForDataset(ctx, datasetID)
	if err != nil {
		return domain.DatasetWithBloggers{}, fmt.Errorf("failed to find bloggers: %v", err)
	}

	return domain.NewDatasetWithBloggers(dataset, bloggers), nil
}
