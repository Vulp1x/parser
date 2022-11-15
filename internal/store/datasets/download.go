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

func (s Store) DownloadTargets(ctx context.Context, datasetID uuid.UUID) (domain.Targets, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	_, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDatasetNotFound
		}

		return nil, fmt.Errorf("failed to find dataset: %v", err)
	}

	targets, err := q.FindTargetsForDataset(ctx, datasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find targets: %v", err)
	}

	return targets, nil
}
