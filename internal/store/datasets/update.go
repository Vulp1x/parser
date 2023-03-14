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

// UpdateDataset обновляет датасет, записывая новых блогеров (старые будут удалены)
func (s Store) UpdateDataset(ctx context.Context, datasetID uuid.UUID, originalAccounts []string, opts ...UpdateOption) (domain.DatasetWithBloggers, error) {
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

	updateConfig := toUpdateConfig(dataset)
	for _, opt := range opts {
		opt(&updateConfig)
	}

	dataset, err = q.UpdateDataset(ctx, updateConfig.toUpdateParams())
	if err != nil {
		return domain.DatasetWithBloggers{}, err
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
		return fmt.Errorf("failed to delete previous bloggers: %v", err)
	}

	logger.Infof(ctx, "deleted %d bloggers", tag.RowsAffected())

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

type UpdateConfig struct {
	datasetID           uuid.UUID
	phoneCode           *int32
	title               string
	postsPerBlogger     int32
	likedPerPost        int32
	commentedPerPost    int32
	followersPerBlogger int32
}

func toUpdateConfig(dataset dbmodel.Dataset) UpdateConfig {
	return UpdateConfig{
		datasetID:        dataset.ID,
		phoneCode:        dataset.PhoneCode,
		title:            dataset.Title,
		postsPerBlogger:  dataset.PostsPerBlogger,
		likedPerPost:     dataset.LikedPerPost,
		commentedPerPost: dataset.CommentedPerPost,
	}
}

func (c UpdateConfig) toUpdateParams() dbmodel.UpdateDatasetParams {
	return dbmodel.UpdateDatasetParams{
		PhoneCode:        c.phoneCode,
		Title:            c.title,
		PostsPerBlogger:  c.postsPerBlogger,
		LikedPerPost:     c.likedPerPost,
		CommentedPerPost: c.commentedPerPost,
		ID:               c.datasetID,
		FollowersCount:   c.followersPerBlogger,
	}
}

type UpdateOption func(config *UpdateConfig)

func WithUpdatePhoneCodeOption(phoneCode *int32) UpdateOption {
	return func(config *UpdateConfig) {
		config.phoneCode = phoneCode
	}
}

func WithUpdateTitleOption(title *string) UpdateOption {
	return func(config *UpdateConfig) {
		if title != nil {
			config.title = *title
		}
	}
}

func WithUpdatePostsPerBloggerOption(postsPerBlogger *uint) UpdateOption {
	return func(config *UpdateConfig) {
		if postsPerBlogger != nil {
			config.postsPerBlogger = int32(*postsPerBlogger)
		}
	}
}

func WithUpdateLikedPerPostOption(likedPerPost *uint) UpdateOption {
	return func(config *UpdateConfig) {
		if likedPerPost != nil {
			config.likedPerPost = int32(*likedPerPost)
		}
	}
}

func WithUpdateCommentedPerPostOption(commentedPerPost *uint) UpdateOption {
	return func(config *UpdateConfig) {
		if commentedPerPost != nil {
			config.commentedPerPost = int32(*commentedPerPost)
		}
	}
}

func WithUpdateFollowersPerBloggerOption(followersPerBlogger *uint) UpdateOption {
	return func(config *UpdateConfig) {
		if followersPerBlogger != nil {
			config.followersPerBlogger = int32(*followersPerBlogger)
		}
	}
}
