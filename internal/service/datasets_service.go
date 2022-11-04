package service

import (
	"context"

	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/sessions"
	"github.com/inst-api/parser/internal/store/datasets"
	"github.com/inst-api/parser/pkg/logger"
	"goa.design/goa/v3/security"
)

// datasets_service service example implementation.
// The example methods log the requests and return zero values.
type datasetsServicesrvc struct {
	auth *authService
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

// создать драфт задачи
func (s *datasetsServicesrvc) CreateDatasetDraft(ctx context.Context, p *datasetsservice.CreateDatasetDraftPayload) (res string, err error) {
	logger.Info(ctx, "datasetsService.create dataset draft")
	return
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
func (s *datasetsServicesrvc) GetDataset(ctx context.Context, p *datasetsservice.GetDatasetPayload) (res *datasetsservice.Dataset, err error) {
	res = &datasetsservice.Dataset{}
	logger.Info(ctx, "datasetsService.get dataset")
	return
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
