package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/sessions"
	"github.com/inst-api/parser/internal/store/datasets"
	"github.com/inst-api/parser/pkg/logger"
	"goa.design/goa/v3/security"
)

type datasetsStore interface {
	CreateDraftDataset(ctx context.Context, userID uuid.UUID, datasetType dbmodel.DatasetType) (uuid.UUID, error)
	GetDataset(ctx context.Context, datasetID uuid.UUID) (domain.DatasetWithBloggers, error)
	UpdateDataset(ctx context.Context, datasetID uuid.UUID, originalAccounts []string, opts ...datasets.UpdateOption) (domain.DatasetWithBloggers, error)
	List(ctx context.Context, managerID uuid.UUID) (domain.Datasets, error)
	FindSimilarBloggers(ctx context.Context, datasetID uuid.UUID) (domain.DatasetWithBloggers, error)
	ParseTargetUsers(ctx context.Context, datasetID uuid.UUID) (domain.DatasetWithBloggers, error)
	ParsingProgress(ctx context.Context, datasetID uuid.UUID) (domain.ParsingProgress, error)
	DownloadTargets(ctx context.Context, datasetID uuid.UUID) (domain.Targets, error)
}

// datasets_service service example implementation.
// The example methods log the requests and return zero values.
type datasetsServicesrvc struct {
	auth  *authService
	store datasetsStore
}

func (s *datasetsServicesrvc) UploadFiles(ctx context.Context, p *datasetsservice.UploadFilesPayload) (*datasetsservice.UploadFilesResult, error) {
	logger.Infof(ctx, "got upload files for dataset %s", p.DatasetID)
	return &datasetsservice.UploadFilesResult{}, nil
}

// NewDatasetsService returns the datasets_service service implementation.
func NewDatasetsService(cfg sessions.Configuration, store datasetsStore) datasetsservice.Service {
	return &datasetsServicesrvc{
		auth:  &authService{securityCfg: cfg},
		store: store,
	}
}

// JWTAuth implements the authorization logic for service "datasets_service"
// for the "jwt" security scheme.
func (s *datasetsServicesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	if s.auth == nil {
		logger.Error(ctx, "datasets service has nil auther")
		return ctx, datasetsservice.Unauthorized("internal error")
	}

	return s.auth.JWTAuth(ctx, token, scheme)
}

// CreateDatasetDraft создать драфт задачи
func (s *datasetsServicesrvc) CreateDatasetDraft(ctx context.Context, p *datasetsservice.CreateDatasetDraftPayload) (string, error) {
	logger.Info(ctx, "datasetsService.create dataset draft")

	userID, err := UserIDFromContext(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to get user id from context: %v", err)
		return "", internalErr(err)
	}

	dbType, ok := datasetTypesToDBType[p.Type]
	if !ok {
		logger.Errorf(ctx, "invalid dataset type %d: failed to find it in %v", p.Type, datasetTypesToDBType)
		return "", datasetsservice.BadRequest(fmt.Sprintf("unexpected dataset type %d", p.Type))
	}

	taskID, err := s.store.CreateDraftDataset(ctx, userID, dbType)
	if err != nil {
		return "", internalErr(err)
	}

	return taskID.String(), nil
}

// UpdateDataset обновляет информацию о задаче. Не меняет статус задачи, можно вызывать
// сколько угодно раз.
// Нельзя вызвать для задачи, которая уже выполняется, для этого надо сначала
// остановить выполнение.
func (s *datasetsServicesrvc) UpdateDataset(ctx context.Context, p *datasetsservice.UpdateDatasetPayload) (*datasetsservice.Dataset, error) {
	ctx = logger.WithFields(ctx, logger.Fields{"dataset_id": p.DatasetID})
	logger.Info(ctx, "datasetsService.update dataset")

	datasetID, err := uuid.Parse(p.DatasetID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, datasetsservice.BadRequest(err.Error())
	}

	dataset, err := s.store.UpdateDataset(ctx, datasetID, p.OriginalAccounts,
		datasets.WithUpdatePostsPerBloggerOption(p.PostsPerBlogger),
		datasets.WithUpdateCommentedPerPostOption(p.CommentedPerPost),
		datasets.WithUpdateLikedPerPostOption(p.LikedPerPost),
		datasets.WithUpdatePhoneCodeOption(p.PhoneCode),
		datasets.WithUpdateTitleOption(p.Title),
	)
	if err != nil {
		logger.Errorf(ctx, "failed to update dataset: %v", err)

		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	return dataset.ToProto(), nil
}

// FindSimilar начать выполнение задачи
func (s *datasetsServicesrvc) FindSimilar(ctx context.Context, p *datasetsservice.FindSimilarPayload) (*datasetsservice.Dataset, error) {
	ctx = logger.WithFields(ctx, logger.Fields{"dataset_id": p.DatasetID})

	datasetID, err := uuid.Parse(p.DatasetID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, datasetsservice.BadRequest(err.Error())
	}

	datasetWithBloggers, err := s.store.FindSimilarBloggers(ctx, datasetID)
	if err != nil {
		logger.Errorf(ctx, "failed to find similar bots: %v", err)
		if errors.Is(err, datasets.ErrDatasetInvalidStatus) || errors.Is(err, datasets.ErrNoBlogers) || errors.Is(err, datasets.ErrNoReadyBots) {
			return nil, datasetsservice.BadRequest(err.Error())
		}

		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	return datasetWithBloggers.ToProto(), nil
}

// ParseDataset получить базу доноров для выбранных блогеров
func (s *datasetsServicesrvc) ParseDataset(ctx context.Context, p *datasetsservice.ParseDatasetPayload) (*datasetsservice.ParseDatasetResult, error) {
	ctx = logger.WithFields(ctx, logger.Fields{"dataset_id": p.DatasetID})

	datasetID, err := uuid.Parse(p.DatasetID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, datasetsservice.BadRequest(err.Error())
	}

	dataset, err := s.store.ParseTargetUsers(ctx, datasetID)
	if err != nil {
		logger.Errorf(ctx, "failed to start parsing users: %v", err)
		if errors.Is(err, datasets.ErrDatasetInvalidStatus) || errors.Is(err, datasets.ErrNoBlogers) || errors.Is(err, datasets.ErrNoReadyBots) {
			return nil, datasetsservice.BadRequest(err.Error())
		}

		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	return &datasetsservice.ParseDatasetResult{
		Status:    datasetsservice.DatasetStatus(dataset.Dataset.Status),
		DatasetID: dataset.Dataset.ID.String(),
	}, nil
}

func (s *datasetsServicesrvc) GetParsingProgress(ctx context.Context, p *datasetsservice.GetParsingProgressPayload) (*datasetsservice.ParsingProgress, error) {
	ctx = logger.WithFields(ctx, logger.Fields{"dataset_id": p.DatasetID})
	logger.Debug(ctx, "datasetsService.getParsingProgress")

	datasetID, err := uuid.Parse(p.DatasetID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, datasetsservice.BadRequest(err.Error())
	}

	progress, err := s.store.ParsingProgress(ctx, datasetID)
	if err != nil {
		logger.Errorf(ctx, "failed to get dataset's progress: %v", err)
		if errors.Is(err, datasets.ErrDatasetInvalidStatus) || errors.Is(err, datasets.ErrNoBlogers) || errors.Is(err, datasets.ErrNoReadyBots) {
			return nil, datasetsservice.BadRequest(err.Error())
		}

		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	return progress.ToProto(), nil
}

// GetDataset получить задачу по id
func (s *datasetsServicesrvc) GetDataset(ctx context.Context, p *datasetsservice.GetDatasetPayload) (*datasetsservice.Dataset, error) {
	ctx = logger.WithFields(ctx, logger.Fields{"dataset_id": p.DatasetID})
	logger.Info(ctx, "datasetsService.get dataset")

	datasetID, err := uuid.Parse(p.DatasetID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, datasetsservice.BadRequest(err.Error())
	}

	dataset, err := s.store.GetDataset(ctx, datasetID)
	if err != nil {
		logger.Errorf(ctx, "failed to find dataset by id: %v", err)

		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	return dataset.ToProto(), nil
}

// GetProgress получить статус выполнения задачи по id
func (s *datasetsServicesrvc) GetProgress(ctx context.Context, p *datasetsservice.GetProgressPayload) (*datasetsservice.DatasetProgress, error) {
	ctx = logger.WithFields(ctx, logger.Fields{"dataset_id": p.DatasetID})
	logger.Info(ctx, "datasetsService.get progress")

	datasetID, err := uuid.Parse(p.DatasetID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, datasetsservice.BadRequest(err.Error())
	}

	dataset, err := s.store.GetDataset(ctx, datasetID)
	if err != nil {
		logger.Errorf(ctx, "failed to find dataset by id: %v", err)

		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	var initialBloggersCount, newBloggersCount, filteredBloggersCount int

	for _, blogger := range dataset.Bloggers {
		if blogger.IsInitial {
			initialBloggersCount++
		} else {
			newBloggersCount++
		}

		if dataset.Dataset.PhoneCode != nil &&
			blogger.PublicPhoneCountryCode != nil &&
			fmt.Sprintf("%d", dataset.Dataset.PhoneCode) == *blogger.PublicPhoneCountryCode {
			filteredBloggersCount++
		}
	}

	return &datasetsservice.DatasetProgress{
		Bloggers:         dataset.ToBloggersProto(),
		InitialBloggers:  initialBloggersCount,
		NewBloggers:      newBloggersCount,
		FilteredBloggers: filteredBloggersCount,
		Done:             dataset.IsReadyForParsing(),
	}, nil
}

// ListDatasets получить все задачи для текущего пользователя
func (s *datasetsServicesrvc) ListDatasets(ctx context.Context, p *datasetsservice.ListDatasetsPayload) ([]*datasetsservice.Dataset, error) {
	logger.Info(ctx, "datasetsService.list datasets")

	managerID, err := UserIDFromContext(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to get user id from context: %v", err)
		return nil, internalErr(err)
	}

	domainDatasets, err := s.store.List(ctx, managerID)
	if err != nil {
		logger.Errorf(ctx, "failed to list datasets: %v", err)
		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	return domainDatasets.ToProto(), nil
}

func (s *datasetsServicesrvc) DownloadTargets(ctx context.Context, p *datasetsservice.DownloadTargetsPayload) ([]string, error) {
	ctx = logger.WithFields(ctx, logger.Fields{"dataset_id": p.DatasetID})
	logger.Info(ctx, "datasetsService.get progress")

	datasetID, err := uuid.Parse(p.DatasetID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, datasetsservice.BadRequest(err.Error())
	}

	targets, err := s.store.DownloadTargets(ctx, datasetID)
	if err != nil {
		logger.Errorf(ctx, "failed to find download targets: %v", err)
		if errors.Is(err, datasets.ErrDatasetNotFound) {
			return nil, datasetsservice.DatasetNotFound("")
		}

		return nil, internalErr(err)
	}

	return targets.ToProto(p.Format), nil
}

func internalErr(err error) datasetsservice.InternalError {
	return datasetsservice.InternalError(err.Error())
}
