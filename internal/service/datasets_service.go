package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/sessions"
	"github.com/inst-api/parser/internal/store/datasets"
	"github.com/inst-api/parser/pkg/logger"
	"goa.design/goa/v3/security"
)

type datasetsStore interface {
	CreateDraftDataset(ctx context.Context, userID uuid.UUID, title string) (uuid.UUID, error)
	GetDataset(ctx context.Context, datasetID uuid.UUID) (dbmodel.Dataset, error)
}

// datasets_service service example implementation.
// The example methods log the requests and return zero values.
type datasetsServicesrvc struct {
	auth  *authService
	store datasetsStore
}

// NewDatasetsService returns the datasets_service service implementation.
func NewDatasetsService(cfg sessions.Configuration, store *datasets.Store) datasetsservice.Service {
	return &datasetsServicesrvc{
		auth: &authService{securityCfg: cfg},
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
		return "", datasetsservice.InternalError(err.Error())
	}

	taskID, err := s.store.CreateDraftDataset(ctx, userID, p.Title)
	if err != nil {
		return "", datasetsservice.InternalError(err.Error())
	}

	return taskID.String(), nil
}

// UpdateDataset обновляет информацию о задаче. Не меняет статус задачи, можно вызывать
// сколько угодно раз.
// Нельзя вызвать для задачи, которая уже выполняется, для этого надо сначала
// остановить выполнение.
func (s *datasetsServicesrvc) UpdateDataset(ctx context.Context, p *datasetsservice.UpdateDatasetPayload) (res *datasetsservice.Dataset, err error) {
	res = &datasetsservice.Dataset{}
	logger.Info(ctx, "datasetsService.update dataset")

	return
}

// начать выполнение задачи
func (s *datasetsServicesrvc) FindSimilar(ctx context.Context, p *datasetsservice.FindSimilarPayload) (res *datasetsservice.FindSimilarResult, err error) {
	res = &datasetsservice.FindSimilarResult{}
	logger.Info(ctx, "datasetsService.find similar")
	return
}

// получить базу доноров для выбранных блогеров
func (s *datasetsServicesrvc) ParseDataset(ctx context.Context, p *datasetsservice.ParseDatasetPayload) (res *datasetsservice.ParseDatasetResult, err error) {
	res = &datasetsservice.ParseDatasetResult{}
	logger.Info(ctx, "datasetsService.parse dataset")
	return
}

// получить задачу по id
func (s *datasetsServicesrvc) GetDataset(ctx context.Context, p *datasetsservice.GetDatasetPayload) (*datasetsservice.Dataset, error) {
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

		return nil, datasetsservice.InternalError(err.Error())
	}

	return &datasetsservice.Dataset{
		ID: dataset.ID.String(),
	}, nil
}

// получить статус выполнения задачи по id
func (s *datasetsServicesrvc) GetProgress(ctx context.Context, p *datasetsservice.GetProgressPayload) (res *datasetsservice.DatasetProgress, err error) {
	res = &datasetsservice.DatasetProgress{}
	logger.Info(ctx, "datasetsService.get progress")
	return
}

// получить все задачи для текущего пользователя
func (s *datasetsServicesrvc) ListDatasets(ctx context.Context, p *datasetsservice.ListDatasetsPayload) (res []*datasetsservice.Dataset, err error) {
	logger.Info(ctx, "datasetsService.list datasets")
	return
}
