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
	FindSimilar(context.Context, *FindSimilarPayload) (res *Dataset, err error)
	// получить статус выполнения поиска похожих аккаунтов по айди датасета
	GetProgress(context.Context, *GetProgressPayload) (res *DatasetProgress, err error)
	// получить базу доноров для выбранных блогеров
	ParseDataset(context.Context, *ParseDatasetPayload) (res *ParseDatasetResult, err error)
	// получить задачу по id
	GetDataset(context.Context, *GetDatasetPayload) (res *Dataset, err error)
	// получить статус выполнения парсинга аккаунтов
	GetParsingProgress(context.Context, *GetParsingProgressPayload) (res *ParsingProgress, err error)
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
var MethodNames = [8]string{"create dataset draft", "update dataset", "find similar", "get progress", "parse dataset", "get dataset", "get parsing progress", "list datasets"}

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

// CreateDatasetDraftPayload is the payload type of the datasets_service
// service create dataset draft method.
type CreateDatasetDraftPayload struct {
	// JWT used for authentication
	Token string
}

// Dataset is the result type of the datasets_service service update dataset
// method.
type Dataset struct {
	ID       string
	Bloggers []*Blogger
	Status   DatasetStatus
	// название задачи
	Title string
	// имена аккаунтов, для которых ищем похожих
	PostsPerBlogger int32 `json:"posts_per_blogger"`
	// сколько лайкнувших для каждого поста брать
	LikedPerPost int32 `json:"liked_per_post"`
	// сколько прокоментировааших для каждого поста брать
	CommentedPerPost int32 `json:"commented_per_post"`
}

// DatasetProgress is the result type of the datasets_service service get
// progress method.
type DatasetProgress struct {
	// блогеры, которых уже нашли
	Bloggers []*Blogger
	// количество блогеров, которые были изначально
	InitialBloggers int `json:"initial_bloggers"`
	// количество блогеров, которых нашли
	NewBloggers int `json:"new_bloggers"`
	// количество блогеров, которые проходят проверку по коду региона
	FilteredBloggers int `json:"filtered_bloggers"`
	// закончена ли задача
	Done bool
}

// 1 - датасет только создан
// 2- начали поиск блогеров
// 3- успешно закончили поиска похожих блогеров
// 4- начали парсинг юзеров у блогеров
// 5- успешно закончили парсинг юзеров
// 6- всё сломалось
type DatasetStatus int

// FindSimilarPayload is the payload type of the datasets_service service find
// similar method.
type FindSimilarPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// GetDatasetPayload is the payload type of the datasets_service service get
// dataset method.
type GetDatasetPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// GetParsingProgressPayload is the payload type of the datasets_service
// service get parsing progress method.
type GetParsingProgressPayload struct {
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

// ParsingProgress is the result type of the datasets_service service get
// parsing progress method.
type ParsingProgress struct {
	// количество блогеров, у которых спарсили пользователей
	BloggersParsed int `json:"bloggers_parsed"`
	// количество сохраненных доноров
	TargetsSaved int `json:"filtered_bloggers"`
	// закончен ли парсинг блогеров
	Done bool
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
	// имена аккаунтов, для которых ищем похожих
	PostsPerBlogger *uint `json:"posts_per_blogger"`
	// сколько лайкнувших для каждого поста брать
	LikedPerPost *uint `json:"liked_per_post"`
	// сколько прокоментировааших для каждого поста брать
	CommentedPerPost *uint `json:"commented_per_post"`
	// код региона, по которому будем сортировать
	PhoneCode *int32 `json:"phone_code"`
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
