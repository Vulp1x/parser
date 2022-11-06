package domain

import (
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/dbmodel"
)

func NewDatasetWithBloggers(dataset dbmodel.Dataset, bloggers []dbmodel.Blogger) DatasetWithBloggers {
	return DatasetWithBloggers{dataset: dataset, bloggers: bloggers}
}

type DatasetWithBloggers struct {
	dataset  dbmodel.Dataset
	bloggers []dbmodel.Blogger
}

func (b DatasetWithBloggers) ToProto() *datasetsservice.Dataset {
	var bloggers = make([]*datasetsservice.Blogger, len(b.bloggers))

	for i, blogger := range b.bloggers {
		bloggers[i] = &datasetsservice.Blogger{
			ID:        blogger.ID.String(),
			Username:  blogger.Username,
			UserID:    blogger.UserID,
			DatasetID: blogger.DatasetID.String(),
			IsInitial: blogger.IsInitial,
		}
	}

	return &datasetsservice.Dataset{
		ID:       b.dataset.ID.String(),
		Bloggers: bloggers,
		Status:   datasetsservice.DatasetStatus(b.dataset.Status),
		Title:    b.dataset.Title,
	}
}
