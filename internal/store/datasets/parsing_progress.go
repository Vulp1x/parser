package datasets

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (s Store) ParsingProgress(ctx context.Context, datasetID uuid.UUID) (domain.ParsingProgress, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ParsingProgress{}, ErrDatasetNotFound
		}

		return domain.ParsingProgress{}, err
	}

	progress, err := q.GetParsingProgress(ctx, datasetID)
	if err != nil {
		return domain.ParsingProgress{}, err
	}

	return domain.ParsingProgress{
		BloggersParsed:     int(progress.ParsedBloggersCount),
		TargetsSaved:       int(progress.TargetsSavedCount),
		TotalBloggersCount: int32(progress.TotalBloggers),
		Done:               dataset.Status == dbmodel.ParsingDoneDatasetStatus,
	}, nil
}
