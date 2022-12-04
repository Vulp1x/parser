package domain

import (
	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/pb/instaproxy"
)

type Media dbmodel.Media

type Medias []Media

func (m Medias) ToSaveMediasParams() []dbmodel.SaveMediasParams {
	saveParams := make([]dbmodel.SaveMediasParams, len(m))
	for i, media := range m {
		saveParams[i] = dbmodel.SaveMediasParams{
			Pk:              media.Pk,
			ID:              media.ID,
			DatasetID:       media.DatasetID,
			MediaType:       media.MediaType,
			Code:            media.Code,
			HasMoreComments: media.HasMoreComments,
			Caption:         media.Caption,
			Width:           media.Width,
			Height:          media.Height,
			LikeCount:       media.LikeCount,
			TakenAt:         media.TakenAt,
		}
	}

	return saveParams
}

func MediasFromProto(protoMedias []*instaproxy.Media, datasetID uuid.UUID) Medias {
	domainMedias := make([]Media, len(protoMedias))
	for i, media := range protoMedias {
		domainMedias[i] = Media{
			Pk:              media.Pk,
			ID:              media.Id,
			DatasetID:       datasetID,
			MediaType:       int32(media.MediaType),
			Code:            media.Code,
			HasMoreComments: media.HasMoreComments,
			Caption:         media.Caption.Text,
			Width:           int32(media.OriginalWidth),
			Height:          int32(media.OriginalHeight),
			LikeCount:       int32(media.LikeCount),
			TakenAt:         int32(media.TakenAt),
		}
	}

	return domainMedias
}
