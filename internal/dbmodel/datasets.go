package dbmodel

type datasetStatus int16

const (
	// DraftDatasetStatus задача только создана, нужно загрузить список ботов, прокси и получателей
	DraftDatasetStatus datasetStatus = 1
	// FindingSimilarStarted запустили поиска похожих блогеров
	FindingSimilarStarted datasetStatus = 2
	// ReadyForParsingDatasetStatus закончили поиск похожих блогеров
	ReadyForParsingDatasetStatus datasetStatus = 3
	// ParsingTargetsStartedDatasetStatus начали парсинг юзеров
	ParsingTargetsStartedDatasetStatus datasetStatus = 4
	// ParsingDoneDatasetStatus - закончили парсинг целевых пользователей
	ParsingDoneDatasetStatus datasetStatus = 5
	// ErrorDatasetStatus проблема при парсинге/поиска похожих в датасете
	ErrorDatasetStatus datasetStatus = 6
)
