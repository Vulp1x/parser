package domain

import (
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
)

type ParsingProgress struct {
	BloggersParsed     int
	TargetsSaved       int
	TotalBloggersCount int32
	Done               bool
}

func (p ParsingProgress) ToProto() *datasetsservice.ParsingProgress {
	return &datasetsservice.ParsingProgress{
		BloggersParsed: p.BloggersParsed,
		TargetsSaved:   p.TargetsSaved,
		Done:           p.Done,
	}
}
