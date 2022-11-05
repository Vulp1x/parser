package datasets

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/jackc/pgx/v4"
)

var ErrDatasetNotFound = errors.New("dataset not found")

func (s Store) GetDataset(ctx context.Context, datasetID uuid.UUID) (dbmodel.Dataset, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbmodel.Dataset{}, ErrDatasetNotFound
		}

		return dbmodel.Dataset{}, err
	}

	return dataset, nil
}
