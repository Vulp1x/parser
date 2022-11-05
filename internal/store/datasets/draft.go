package datasets

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
)

// CreateDraftDataset создаем драфт
func (s Store) CreateDraftDataset(ctx context.Context, userID uuid.UUID, title string) (uuid.UUID, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	taskID, err := q.CreateDraftDataset(ctx, dbmodel.CreateDraftDatasetParams{UserID: userID, Title: title})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create dataset draft: %w", err)
	}

	return taskID, nil
}
