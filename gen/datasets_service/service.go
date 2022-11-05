// Code generated by goa v3.10.2, DO NOT EDIT.
//
// datasets_service service
//
// Command:
// $ goa gen github.com/inst-api/parser/design

package datasetsservice

import (
	"context"

	"goa.design/goa/v3/security"
)

// сервис для создания, редактирования и работы с задачами (рекламными
// компаниями)
type Service interface {
	// создать драфт задачи
	CreateDatasetDraft(context.Context, *CreateDatasetDraftPayload) (res string, err error)
	// обновить информацию о задаче. Не меняет статус задачи, можно вызывать
	// сколько угодно раз.
	// Нельзя вызвать для задачи, которая уже выполняется, для этого надо сначала
	// остановить выполнение.
	UpdateDataset(context.Context, *UpdateDatasetPayload) (res *Dataset, err error)
	// начать выполнение задачи
	FindSimilar(context.Context, *FindSimilarPayload) (res *FindSimilarResult, err error)
	// получить базу доноров для выбранных блогеров
	ParseDataset(context.Context, *ParseDatasetPayload) (res *ParseDatasetResult, err error)
	// получить задачу по id
	GetDataset(context.Context, *GetDatasetPayload) (res *Dataset, err error)
	// получить статус выполнения задачи по id
	GetProgress(context.Context, *GetProgressPayload) (res *DatasetProgress, err error)
	// получить все задачи для текущего пользователя
	ListDatasets(context.Context, *ListDatasetsPayload) (res []*Dataset, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "datasets_service"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [7]string{"create dataset draft", "update dataset", "find similar", "parse dataset", "get dataset", "get progress", "list datasets"}

type Blogger struct {
	ID string
	// имя аккаунта в инстаграме
	Username string
	// user_id в инстаграме, -1 если неизвестен
	UserID int64 `json:"user_id"`
	// айди датасета, к которому принадлежит блоггер
	DatasetID string `json:"dataset_id"`
	// является ли блоггер изначально в датасете или появился при парсинге
	IsInitial bool `json:"is_initial"`
}

type BloggersProgress struct {
	// имя пользователя бота
	UserName string `json:"user_name"`
	// количество выложенных постов
	PostsCount int `json:"posts_count"`
	// текущий статус бота, будут ли выкладываться посты
	Status int
}

// CreateDatasetDraftPayload is the payload type of the datasets_service
// service create dataset draft method.
type CreateDatasetDraftPayload struct {
	// JWT used for authentication
	Token string
	// название задачи
	Title string
}

// Dataset is the result type of the datasets_service service update dataset
// method.
type Dataset struct {
	ID       string
	Bloggers []*Blogger
	Status   DatasetStatus
	// название задачи
	Title string
}

// DatasetProgress is the result type of the datasets_service service get
// progress method.
type DatasetProgress struct {
	// результат работы по каждому боту, ключ- имя бота
	BotsProgresses map[string]*BloggersProgress `json:"bots_progresses"`
	// количество аккаунтов, которых упомянули в постах
	TargetsNotified int `json:"targets_notified"`
	// количество аккаунтов, которых не получилось упомянуть, при перезапуске
	// задачи будут использованы заново
	TargetsFailed int `json:"targets_failed"`
	// количество аккаунтов, которых не выбрали для постов
	TargetsWaiting int `json:"targets_waiting,targets_waiting"`
	// закончена ли задача
	Done bool
}

// 1 - задача только создана, нужно загрузить список ботов, прокси и получателей
// 2- в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
// 3- задача готова к запуску
// 4- задача запущена
// 5 - задача остановлена
// 6 - задача завершена
type DatasetStatus int

// FindSimilarPayload is the payload type of the datasets_service service find
// similar method.
type FindSimilarPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// FindSimilarResult is the result type of the datasets_service service find
// similar method.
type FindSimilarResult struct {
	Status DatasetStatus
	// id задачи
	DatasetID string `json:"dataset_id"`
	Bloggers  []*Blogger
}

// GetDatasetPayload is the payload type of the datasets_service service get
// dataset method.
type GetDatasetPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// GetProgressPayload is the payload type of the datasets_service service get
// progress method.
type GetProgressPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// ListDatasetsPayload is the payload type of the datasets_service service list
// datasets method.
type ListDatasetsPayload struct {
	// JWT used for authentication
	Token string
}

// ParseDatasetPayload is the payload type of the datasets_service service
// parse dataset method.
type ParseDatasetPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// ParseDatasetResult is the result type of the datasets_service service parse
// dataset method.
type ParseDatasetResult struct {
	Status DatasetStatus
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// UpdateDatasetPayload is the payload type of the datasets_service service
// update dataset method.
type UpdateDatasetPayload struct {
	// JWT used for authentication
	Token string
	// id задачи, которую хотим обновить
	DatasetID string `json:"dataset_id"`
	// имена аккаунтов, для которых ищем похожих
	OriginalAccounts []string `json:"original_accounts"`
	// код региона, по которому будем сортировать
	PhoneCode *int `json:"phone_code"`
	// название задачи
	Title *string
}

// Invalid request
type BadRequest string

// Not found
type DatasetNotFound string

// internal error
type InternalError string

// Credentials are invalid
type Unauthorized string

// Error returns an error description.
func (e BadRequest) Error() string {
	return "Invalid request"
}

// ErrorName returns "bad request".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e BadRequest) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "bad request".
func (e BadRequest) GoaErrorName() string {
	return "bad request"
}

// Error returns an error description.
func (e DatasetNotFound) Error() string {
	return "Not found"
}

// ErrorName returns "dataset not found".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e DatasetNotFound) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "dataset not found".
func (e DatasetNotFound) GoaErrorName() string {
	return "dataset not found"
}

// Error returns an error description.
func (e InternalError) Error() string {
	return "internal error"
}

// ErrorName returns "internal error".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e InternalError) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "internal error".
func (e InternalError) GoaErrorName() string {
	return "internal error"
}

// Error returns an error description.
func (e Unauthorized) Error() string {
	return "Credentials are invalid"
}

// ErrorName returns "unauthorized".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e Unauthorized) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "unauthorized".
func (e Unauthorized) GoaErrorName() string {
	return "unauthorized"
}
