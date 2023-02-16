package datasets

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (s Store) DownloadTargets(ctx context.Context, datasetID uuid.UUID, format int) ([]string, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDatasetNotFound
		}

		return nil, fmt.Errorf("failed to find dataset: %v", err)
	}

	if dataset.Type == dbmodel.DatasetTypePhoneNumbers {
		var fullTargets []dbmodel.FullTarget
		if dataset.PhoneCode != nil {
			fullTargets, err = q.FindFullTargetsWithCode(ctx, dbmodel.FindFullTargetsWithCodeParams{
				DatasetID:              datasetID,
				PublicPhoneCountryCode: strconv.Itoa(int(*dataset.PhoneCode)),
			})
		} else {
			fullTargets, err = q.FindFullTargets(ctx, datasetID)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to find full targets: %v", err)
		}

		formattedTargets := make([]string, len(fullTargets))
		for i, target := range fullTargets {
			formattedTargets[i] = domain.FullUser(target).Format(format)
		}

		return formattedTargets, nil
	}

	targets, err := q.FindTargetsForDataset(ctx, datasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find targets: %v", err)
	}

	return domain.Targets(targets).ToProto(format), nil
}
