package service

import (
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/dbmodel"
)

var datasetTypesToDBType = map[datasetsservice.DatasetType]dbmodel.DatasetType{
	1: dbmodel.DatasetTypeFollowers,
	2: dbmodel.DatasetTypePhoneNumbers,
	3: dbmodel.DatasetTypeLikesAndComments,
}
