package domain

import (
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/dbmodel"
)

func NewDatasetWithBloggers(dataset dbmodel.Dataset, bloggers []dbmodel.Blogger) DatasetWithBloggers {
	return DatasetWithBloggers{Dataset: dataset, Bloggers: bloggers}
}

type DatasetWithBloggers struct {
	Dataset  dbmodel.Dataset
	Bloggers []dbmodel.Blogger
}

func (b DatasetWithBloggers) ToProto() *datasetsservice.Dataset {
	var bloggers = make([]*datasetsservice.Blogger, len(b.Bloggers))

	for i, blogger := range b.Bloggers {
		bloggers[i] = &datasetsservice.Blogger{
			ID:        blogger.ID.String(),
			Username:  blogger.Username,
			UserID:    blogger.UserID,
			DatasetID: blogger.DatasetID.String(),
			IsInitial: blogger.IsInitial,
		}
	}

	return &datasetsservice.Dataset{
		ID:               b.Dataset.ID.String(),
		Bloggers:         bloggers,
		Status:           datasetsservice.DatasetStatus(b.Dataset.Status),
		Title:            b.Dataset.Title,
		PostsPerBlogger:  b.Dataset.PostsPerBlogger,
		LikedPerPost:     b.Dataset.LikedPerPost,
		CommentedPerPost: b.Dataset.CommentedPerPost,
	}
}

func (b DatasetWithBloggers) ToBloggersProto() []*datasetsservice.Blogger {
	var bloggers = make([]*datasetsservice.Blogger, len(b.Bloggers))

	for i, blogger := range b.Bloggers {
		bloggers[i] = &datasetsservice.Blogger{
			ID:        blogger.ID.String(),
			Username:  blogger.Username,
			UserID:    blogger.UserID,
			DatasetID: blogger.DatasetID.String(),
			IsInitial: blogger.IsInitial,
		}
	}

	return bloggers
}

// Usernames возвращает usernames блогеров
func (b DatasetWithBloggers) Usernames() []string {
	var usernames = make([]string, len(b.Bloggers))

	for i, blogger := range b.Bloggers {
		usernames[i] = blogger.Username
	}

	return usernames
}

func (b DatasetWithBloggers) IsReadyForParsing() bool {
	return b.Dataset.Status == dbmodel.ReadyForParsingDatasetStatus
}

type Datasets []dbmodel.Dataset

func (d Datasets) ToProto() []*datasetsservice.Dataset {
	var protoDatasets = make([]*datasetsservice.Dataset, len(d))

	for i, dataset := range d {
		protoDatasets[i] = &datasetsservice.Dataset{
			ID:               dataset.ID.String(),
			Status:           datasetsservice.DatasetStatus(dataset.Status),
			Title:            dataset.Title,
			PostsPerBlogger:  dataset.PostsPerBlogger,
			LikedPerPost:     dataset.LikedPerPost,
			CommentedPerPost: dataset.CommentedPerPost,
			PhoneCode:        dataset.PhoneCode,
		}
	}

	return protoDatasets
}

type Dataset dbmodel.Dataset
