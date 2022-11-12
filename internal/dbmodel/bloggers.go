package dbmodel

type bloggerStatus int16

const (
	// DraftBloggerStatus драфт блогера, возможно нет user_id
	DraftBloggerStatus bloggerStatus = 1
	// FoundBloggerStatus
	FoundBloggerStatus bloggerStatus = 2
	// TargetsParsedBloggerStatus закончили поиск похожих блогеров
	TargetsParsedBloggerStatus bloggerStatus = 3
)
