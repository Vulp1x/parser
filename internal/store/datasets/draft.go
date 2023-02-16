package datasets

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
)

// CreateDraftDataset создаем драфт
func (s Store) CreateDraftDataset(ctx context.Context, managerID uuid.UUID, dbType dbmodel.DatasetType) (uuid.UUID, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	// var datasetStatus = dbmodel.DraftDatasetStatus
	// if dbType == dbmodel.DatasetTypePhoneNumbers {
	// 	datasetStatus = dbmodel.ReadyForParsingDatasetStatus
	// }

	taskID, err := q.CreateDraftDataset(ctx, dbmodel.CreateDraftDatasetParams{ManagerID: managerID, Type: dbType})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create dataset draft: %w", err)
	}

	return taskID, nil
}
