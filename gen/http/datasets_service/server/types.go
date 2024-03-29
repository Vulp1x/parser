// Code generated by goa v3.11.3, DO NOT EDIT.
//
// datasets_service HTTP server types
//
// Command:
// $ goa gen github.com/inst-api/parser/design

package server

import (
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	goa "goa.design/goa/v3/pkg"
)

// CreateDatasetDraftRequestBody is the type of the "datasets_service" service
// "create dataset draft" endpoint HTTP request body.
type CreateDatasetDraftRequestBody struct {
	Type *int `form:"type,omitempty" json:"type,omitempty" xml:"type,omitempty"`
}

// UpdateDatasetRequestBody is the type of the "datasets_service" service
// "update dataset" endpoint HTTP request body.
type UpdateDatasetRequestBody struct {
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
	Title *string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
	// сколько подписчиков для каждого блоггера брать
	SubscribersPerBlogger *uint `json:"subscribers_per_blogger"`
}

// UploadFilesRequestBody is the type of the "datasets_service" service "upload
// files" endpoint HTTP request body.
type UploadFilesRequestBody struct {
	ProxiesFilename *string `json:"proxies_filename"`
	BotsFilename    *string `json:"bots_filename"`
	// список ботов
	Bots []*BotAccountRecordRequestBody `form:"bots,omitempty" json:"bots,omitempty" xml:"bots,omitempty"`
	// список проксей для использования
	Proxies []*ProxyRecordRequestBody `form:"proxies,omitempty" json:"proxies,omitempty" xml:"proxies,omitempty"`
}

// UpdateDatasetOKResponseBody is the type of the "datasets_service" service
// "update dataset" endpoint HTTP response body.
type UpdateDatasetOKResponseBody struct {
	ID     string `form:"id" json:"id" xml:"id"`
	Status int    `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// имена аккаунтов, для которых ищем похожих
	PostsPerBlogger int32 `json:"posts_per_blogger"`
	// сколько лайкнувших для каждого поста брать
	LikedPerPost int32 `json:"liked_per_post"`
	// сколько прокоментировааших для каждого поста брать
	CommentedPerPost int32 `json:"commented_per_post"`
	// является ли блоггер изначально в датасете или появился при парсинге
	PhoneCode *int32 `json:"phone_code"`
	// сколько подписчиков будем парсить у каждого блогера
	SubscribersPerBlogger int32                  `json:"subscribers_per_blogger"`
	Bloggers              []*BloggerResponseBody `form:"bloggers" json:"bloggers" xml:"bloggers"`
	Type                  int                    `form:"type" json:"type" xml:"type"`
}

// FindSimilarOKResponseBody is the type of the "datasets_service" service
// "find similar" endpoint HTTP response body.
type FindSimilarOKResponseBody struct {
	ID     string `form:"id" json:"id" xml:"id"`
	Status int    `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// имена аккаунтов, для которых ищем похожих
	PostsPerBlogger int32 `json:"posts_per_blogger"`
	// сколько лайкнувших для каждого поста брать
	LikedPerPost int32 `json:"liked_per_post"`
	// сколько прокоментировааших для каждого поста брать
	CommentedPerPost int32 `json:"commented_per_post"`
	// является ли блоггер изначально в датасете или появился при парсинге
	PhoneCode *int32 `json:"phone_code"`
	// сколько подписчиков будем парсить у каждого блогера
	SubscribersPerBlogger int32                  `json:"subscribers_per_blogger"`
	Bloggers              []*BloggerResponseBody `form:"bloggers" json:"bloggers" xml:"bloggers"`
	Type                  int                    `form:"type" json:"type" xml:"type"`
}

// GetProgressOKResponseBody is the type of the "datasets_service" service "get
// progress" endpoint HTTP response body.
type GetProgressOKResponseBody struct {
	// блогеры, которых уже нашли
	Bloggers []*BloggerResponseBody `form:"bloggers" json:"bloggers" xml:"bloggers"`
	// количество блогеров, которые были изначально
	InitialBloggers int `json:"initial_bloggers"`
	// количество блогеров, которых нашли
	NewBloggers int `json:"new_bloggers"`
	// количество блогеров, которые проходят проверку по коду региона
	FilteredBloggers int `json:"filtered_bloggers"`
	// закончена ли задача
	Done bool `form:"done" json:"done" xml:"done"`
}

// ParseDatasetOKResponseBody is the type of the "datasets_service" service
// "parse dataset" endpoint HTTP response body.
type ParseDatasetOKResponseBody struct {
	Status int `form:"status" json:"status" xml:"status"`
	// id задачи
	DatasetID string `json:"dataset_id"`
}

// GetDatasetOKResponseBody is the type of the "datasets_service" service "get
// dataset" endpoint HTTP response body.
type GetDatasetOKResponseBody struct {
	ID     string `form:"id" json:"id" xml:"id"`
	Status int    `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// имена аккаунтов, для которых ищем похожих
	PostsPerBlogger int32 `json:"posts_per_blogger"`
	// сколько лайкнувших для каждого поста брать
	LikedPerPost int32 `json:"liked_per_post"`
	// сколько прокоментировааших для каждого поста брать
	CommentedPerPost int32 `json:"commented_per_post"`
	// является ли блоггер изначально в датасете или появился при парсинге
	PhoneCode *int32 `json:"phone_code"`
	// сколько подписчиков будем парсить у каждого блогера
	SubscribersPerBlogger int32                  `json:"subscribers_per_blogger"`
	Bloggers              []*BloggerResponseBody `form:"bloggers" json:"bloggers" xml:"bloggers"`
	Type                  int                    `form:"type" json:"type" xml:"type"`
}

// GetParsingProgressOKResponseBody is the type of the "datasets_service"
// service "get parsing progress" endpoint HTTP response body.
type GetParsingProgressOKResponseBody struct {
	// количество блогеров, у которых спарсили пользователей
	BloggersParsed int `json:"bloggers_parsed"`
	// количество сохраненных доноров
	TargetsSaved int `json:"targets_saved"`
	// закончен ли парсинг блогеров
	Done bool `form:"done" json:"done" xml:"done"`
}

// ListDatasetsResponseBody is the type of the "datasets_service" service "list
// datasets" endpoint HTTP response body.
type ListDatasetsResponseBody []*DatasetResponse

// UploadFilesOKResponseBody is the type of the "datasets_service" service
// "upload files" endpoint HTTP response body.
type UploadFilesOKResponseBody struct {
	// ошибки, которые возникли при загрузке файлов
	UploadErrors []*UploadErrorResponseBody `json:"upload_errors"`
}

// BloggerResponseBody is used to define fields on response body types.
type BloggerResponseBody struct {
	ID string `form:"id" json:"id" xml:"id"`
	// имя аккаунта в инстаграме
	Username string `form:"username" json:"username" xml:"username"`
	// user_id в инстаграме, -1 если неизвестен
	UserID int64 `json:"user_id"`
	// айди датасета, к которому принадлежит блоггер
	DatasetID string `json:"dataset_id"`
	// является ли блоггер изначально в датасете или появился при парсинге
	IsInitial bool `json:"is_initial"`
}

// DatasetResponse is used to define fields on response body types.
type DatasetResponse struct {
	ID     string `form:"id" json:"id" xml:"id"`
	Status int    `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// имена аккаунтов, для которых ищем похожих
	PostsPerBlogger int32 `json:"posts_per_blogger"`
	// сколько лайкнувших для каждого поста брать
	LikedPerPost int32 `json:"liked_per_post"`
	// сколько прокоментировааших для каждого поста брать
	CommentedPerPost int32 `json:"commented_per_post"`
	// является ли блоггер изначально в датасете или появился при парсинге
	PhoneCode *int32 `json:"phone_code"`
	// сколько подписчиков будем парсить у каждого блогера
	SubscribersPerBlogger int32              `json:"subscribers_per_blogger"`
	Bloggers              []*BloggerResponse `form:"bloggers" json:"bloggers" xml:"bloggers"`
	Type                  int                `form:"type" json:"type" xml:"type"`
}

// BloggerResponse is used to define fields on response body types.
type BloggerResponse struct {
	ID string `form:"id" json:"id" xml:"id"`
	// имя аккаунта в инстаграме
	Username string `form:"username" json:"username" xml:"username"`
	// user_id в инстаграме, -1 если неизвестен
	UserID int64 `json:"user_id"`
	// айди датасета, к которому принадлежит блоггер
	DatasetID string `json:"dataset_id"`
	// является ли блоггер изначально в датасете или появился при парсинге
	IsInitial bool `json:"is_initial"`
}

// UploadErrorResponseBody is used to define fields on response body types.
type UploadErrorResponseBody struct {
	// 1 - список ботов
	// 2 - список прокси
	// 3 - список получателей рекламы
	Type int `form:"type" json:"type" xml:"type"`
	Line int `form:"line" json:"line" xml:"line"`
	// номер порта
	Input  string `form:"input" json:"input" xml:"input"`
	Reason string `form:"reason" json:"reason" xml:"reason"`
}

// BotAccountRecordRequestBody is used to define fields on request body types.
type BotAccountRecordRequestBody struct {
	Record []string `form:"record,omitempty" json:"record,omitempty" xml:"record,omitempty"`
	// номер строки в исходном файле
	LineNumber *int `json:"line_number"`
}

// ProxyRecordRequestBody is used to define fields on request body types.
type ProxyRecordRequestBody struct {
	Record []string `form:"record,omitempty" json:"record,omitempty" xml:"record,omitempty"`
	// номер строки в исходном файле
	LineNumber *int `json:"line_number"`
}

// NewUpdateDatasetOKResponseBody builds the HTTP response body from the result
// of the "update dataset" endpoint of the "datasets_service" service.
func NewUpdateDatasetOKResponseBody(res *datasetsservice.Dataset) *UpdateDatasetOKResponseBody {
	body := &UpdateDatasetOKResponseBody{
		ID:                    res.ID,
		Status:                int(res.Status),
		Title:                 res.Title,
		PostsPerBlogger:       res.PostsPerBlogger,
		LikedPerPost:          res.LikedPerPost,
		CommentedPerPost:      res.CommentedPerPost,
		PhoneCode:             res.PhoneCode,
		SubscribersPerBlogger: res.SubscribersPerBlogger,
		Type:                  int(res.Type),
	}
	if res.Bloggers != nil {
		body.Bloggers = make([]*BloggerResponseBody, len(res.Bloggers))
		for i, val := range res.Bloggers {
			body.Bloggers[i] = marshalDatasetsserviceBloggerToBloggerResponseBody(val)
		}
	}
	return body
}

// NewFindSimilarOKResponseBody builds the HTTP response body from the result
// of the "find similar" endpoint of the "datasets_service" service.
func NewFindSimilarOKResponseBody(res *datasetsservice.Dataset) *FindSimilarOKResponseBody {
	body := &FindSimilarOKResponseBody{
		ID:                    res.ID,
		Status:                int(res.Status),
		Title:                 res.Title,
		PostsPerBlogger:       res.PostsPerBlogger,
		LikedPerPost:          res.LikedPerPost,
		CommentedPerPost:      res.CommentedPerPost,
		PhoneCode:             res.PhoneCode,
		SubscribersPerBlogger: res.SubscribersPerBlogger,
		Type:                  int(res.Type),
	}
	if res.Bloggers != nil {
		body.Bloggers = make([]*BloggerResponseBody, len(res.Bloggers))
		for i, val := range res.Bloggers {
			body.Bloggers[i] = marshalDatasetsserviceBloggerToBloggerResponseBody(val)
		}
	}
	return body
}

// NewGetProgressOKResponseBody builds the HTTP response body from the result
// of the "get progress" endpoint of the "datasets_service" service.
func NewGetProgressOKResponseBody(res *datasetsservice.DatasetProgress) *GetProgressOKResponseBody {
	body := &GetProgressOKResponseBody{
		InitialBloggers:  res.InitialBloggers,
		NewBloggers:      res.NewBloggers,
		FilteredBloggers: res.FilteredBloggers,
		Done:             res.Done,
	}
	if res.Bloggers != nil {
		body.Bloggers = make([]*BloggerResponseBody, len(res.Bloggers))
		for i, val := range res.Bloggers {
			body.Bloggers[i] = marshalDatasetsserviceBloggerToBloggerResponseBody(val)
		}
	}
	return body
}

// NewParseDatasetOKResponseBody builds the HTTP response body from the result
// of the "parse dataset" endpoint of the "datasets_service" service.
func NewParseDatasetOKResponseBody(res *datasetsservice.ParseDatasetResult) *ParseDatasetOKResponseBody {
	body := &ParseDatasetOKResponseBody{
		Status:    int(res.Status),
		DatasetID: res.DatasetID,
	}
	return body
}

// NewGetDatasetOKResponseBody builds the HTTP response body from the result of
// the "get dataset" endpoint of the "datasets_service" service.
func NewGetDatasetOKResponseBody(res *datasetsservice.Dataset) *GetDatasetOKResponseBody {
	body := &GetDatasetOKResponseBody{
		ID:                    res.ID,
		Status:                int(res.Status),
		Title:                 res.Title,
		PostsPerBlogger:       res.PostsPerBlogger,
		LikedPerPost:          res.LikedPerPost,
		CommentedPerPost:      res.CommentedPerPost,
		PhoneCode:             res.PhoneCode,
		SubscribersPerBlogger: res.SubscribersPerBlogger,
		Type:                  int(res.Type),
	}
	if res.Bloggers != nil {
		body.Bloggers = make([]*BloggerResponseBody, len(res.Bloggers))
		for i, val := range res.Bloggers {
			body.Bloggers[i] = marshalDatasetsserviceBloggerToBloggerResponseBody(val)
		}
	}
	return body
}

// NewGetParsingProgressOKResponseBody builds the HTTP response body from the
// result of the "get parsing progress" endpoint of the "datasets_service"
// service.
func NewGetParsingProgressOKResponseBody(res *datasetsservice.ParsingProgress) *GetParsingProgressOKResponseBody {
	body := &GetParsingProgressOKResponseBody{
		BloggersParsed: res.BloggersParsed,
		TargetsSaved:   res.TargetsSaved,
		Done:           res.Done,
	}
	return body
}

// NewListDatasetsResponseBody builds the HTTP response body from the result of
// the "list datasets" endpoint of the "datasets_service" service.
func NewListDatasetsResponseBody(res []*datasetsservice.Dataset) ListDatasetsResponseBody {
	body := make([]*DatasetResponse, len(res))
	for i, val := range res {
		body[i] = marshalDatasetsserviceDatasetToDatasetResponse(val)
	}
	return body
}

// NewUploadFilesOKResponseBody builds the HTTP response body from the result
// of the "upload files" endpoint of the "datasets_service" service.
func NewUploadFilesOKResponseBody(res *datasetsservice.UploadFilesResult) *UploadFilesOKResponseBody {
	body := &UploadFilesOKResponseBody{}
	if res.UploadErrors != nil {
		body.UploadErrors = make([]*UploadErrorResponseBody, len(res.UploadErrors))
		for i, val := range res.UploadErrors {
			body.UploadErrors[i] = marshalDatasetsserviceUploadErrorToUploadErrorResponseBody(val)
		}
	}
	return body
}

// NewCreateDatasetDraftPayload builds a datasets_service service create
// dataset draft endpoint payload.
func NewCreateDatasetDraftPayload(body *CreateDatasetDraftRequestBody, token string) *datasetsservice.CreateDatasetDraftPayload {
	v := &datasetsservice.CreateDatasetDraftPayload{
		Type: datasetsservice.DatasetType(*body.Type),
	}
	v.Token = token

	return v
}

// NewUpdateDatasetPayload builds a datasets_service service update dataset
// endpoint payload.
func NewUpdateDatasetPayload(body *UpdateDatasetRequestBody, datasetID string, token string) *datasetsservice.UpdateDatasetPayload {
	v := &datasetsservice.UpdateDatasetPayload{
		PostsPerBlogger:       body.PostsPerBlogger,
		LikedPerPost:          body.LikedPerPost,
		CommentedPerPost:      body.CommentedPerPost,
		PhoneCode:             body.PhoneCode,
		Title:                 body.Title,
		SubscribersPerBlogger: body.SubscribersPerBlogger,
	}
	if body.OriginalAccounts != nil {
		v.OriginalAccounts = make([]string, len(body.OriginalAccounts))
		for i, val := range body.OriginalAccounts {
			v.OriginalAccounts[i] = val
		}
	}
	v.DatasetID = datasetID
	v.Token = token

	return v
}

// NewFindSimilarPayload builds a datasets_service service find similar
// endpoint payload.
func NewFindSimilarPayload(datasetID string, token string) *datasetsservice.FindSimilarPayload {
	v := &datasetsservice.FindSimilarPayload{}
	v.DatasetID = datasetID
	v.Token = token

	return v
}

// NewGetProgressPayload builds a datasets_service service get progress
// endpoint payload.
func NewGetProgressPayload(datasetID string, token string) *datasetsservice.GetProgressPayload {
	v := &datasetsservice.GetProgressPayload{}
	v.DatasetID = datasetID
	v.Token = token

	return v
}

// NewParseDatasetPayload builds a datasets_service service parse dataset
// endpoint payload.
func NewParseDatasetPayload(datasetID string, token string) *datasetsservice.ParseDatasetPayload {
	v := &datasetsservice.ParseDatasetPayload{}
	v.DatasetID = datasetID
	v.Token = token

	return v
}

// NewGetDatasetPayload builds a datasets_service service get dataset endpoint
// payload.
func NewGetDatasetPayload(datasetID string, token string) *datasetsservice.GetDatasetPayload {
	v := &datasetsservice.GetDatasetPayload{}
	v.DatasetID = datasetID
	v.Token = token

	return v
}

// NewGetParsingProgressPayload builds a datasets_service service get parsing
// progress endpoint payload.
func NewGetParsingProgressPayload(datasetID string, token string) *datasetsservice.GetParsingProgressPayload {
	v := &datasetsservice.GetParsingProgressPayload{}
	v.DatasetID = datasetID
	v.Token = token

	return v
}

// NewDownloadTargetsPayload builds a datasets_service service download targets
// endpoint payload.
func NewDownloadTargetsPayload(datasetID string, format int, token string) *datasetsservice.DownloadTargetsPayload {
	v := &datasetsservice.DownloadTargetsPayload{}
	v.DatasetID = datasetID
	v.Format = format
	v.Token = token

	return v
}

// NewListDatasetsPayload builds a datasets_service service list datasets
// endpoint payload.
func NewListDatasetsPayload(token string) *datasetsservice.ListDatasetsPayload {
	v := &datasetsservice.ListDatasetsPayload{}
	v.Token = token

	return v
}

// NewUploadFilesPayload builds a datasets_service service upload files
// endpoint payload.
func NewUploadFilesPayload(body *UploadFilesRequestBody, datasetID string, token string) *datasetsservice.UploadFilesPayload {
	v := &datasetsservice.UploadFilesPayload{
		ProxiesFilename: *body.ProxiesFilename,
		BotsFilename:    *body.BotsFilename,
	}
	v.Bots = make([]*datasetsservice.BotAccountRecord, len(body.Bots))
	for i, val := range body.Bots {
		v.Bots[i] = unmarshalBotAccountRecordRequestBodyToDatasetsserviceBotAccountRecord(val)
	}
	v.Proxies = make([]*datasetsservice.ProxyRecord, len(body.Proxies))
	for i, val := range body.Proxies {
		v.Proxies[i] = unmarshalProxyRecordRequestBodyToDatasetsserviceProxyRecord(val)
	}
	v.DatasetID = datasetID
	v.Token = token

	return v
}

// ValidateCreateDatasetDraftRequestBody runs the validations defined on Create
// Dataset DraftRequestBody
func ValidateCreateDatasetDraftRequestBody(body *CreateDatasetDraftRequestBody) (err error) {
	if body.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "body"))
	}
	if body.Type != nil {
		if !(*body.Type == 1 || *body.Type == 2 || *body.Type == 3) {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.type", *body.Type, []any{1, 2, 3}))
		}
	}
	return
}

// ValidateUpdateDatasetRequestBody runs the validations defined on Update
// DatasetRequestBody
func ValidateUpdateDatasetRequestBody(body *UpdateDatasetRequestBody) (err error) {
	if body.PhoneCode != nil {
		if *body.PhoneCode < 1 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.phone_code", *body.PhoneCode, 1, true))
		}
	}
	if body.PhoneCode != nil {
		if *body.PhoneCode > 1000 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.phone_code", *body.PhoneCode, 1000, false))
		}
	}
	return
}

// ValidateUploadFilesRequestBody runs the validations defined on Upload
// FilesRequestBody
func ValidateUploadFilesRequestBody(body *UploadFilesRequestBody) (err error) {
	if body.Bots == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bots", "body"))
	}
	if body.Proxies == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("proxies", "body"))
	}
	if body.ProxiesFilename == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("proxies_filename", "body"))
	}
	if body.BotsFilename == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bots_filename", "body"))
	}
	for _, e := range body.Bots {
		if e != nil {
			if err2 := ValidateBotAccountRecordRequestBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range body.Proxies {
		if e != nil {
			if err2 := ValidateProxyRecordRequestBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateBotAccountRecordRequestBody runs the validations defined on
// BotAccountRecordRequestBody
func ValidateBotAccountRecordRequestBody(body *BotAccountRecordRequestBody) (err error) {
	if body.Record == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("record", "body"))
	}
	if body.LineNumber == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("line_number", "body"))
	}
	if len(body.Record) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, true))
	}
	if len(body.Record) > 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, false))
	}
	return
}

// ValidateProxyRecordRequestBody runs the validations defined on
// ProxyRecordRequestBody
func ValidateProxyRecordRequestBody(body *ProxyRecordRequestBody) (err error) {
	if body.Record == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("record", "body"))
	}
	if body.LineNumber == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("line_number", "body"))
	}
	if len(body.Record) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, true))
	}
	if len(body.Record) > 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, false))
	}
	return
}
