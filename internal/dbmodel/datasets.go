package dbmodel

type datasetStatus int16

const (
	// DraftDatasetStatus задача только создана, нужно загрузить список ботов, прокси и получателей
	DraftDatasetStatus datasetStatus = 1
	// DataUploadedDatasetStatus в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
	DataUploadedDatasetStatus datasetStatus = 2
	// ReadyDatasetStatus задача готова к запуску
	ReadyDatasetStatus   datasetStatus = 3
	StartedDatasetStatus datasetStatus = 4
	// StoppedDatasetStatus - задача остановлена
	StoppedDatasetStatus datasetStatus = 5
	// DoneDatasetStatus задача выполнена
	DoneDatasetStatus datasetStatus = 6
)
