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

// List возвращает все датасеты
func (s Store) List(ctx context.Context, managerID uuid.UUID) (domain.Datasets, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	datasets, err := q.FindUserDatasets(ctx, managerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDatasetNotFound
		}

		return nil, fmt.Errorf("failed to find datasets for user '%s': %v", managerID, err)
	}

	return datasets, nil
}
