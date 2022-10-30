package model

type CreateArticleForm struct {
	Title       string
	Summary     string
	Content     string
	ContentType string
	Tags        []string
	SourceUrl   string
}

type ImageDTO struct {
	Url string `json:"url"`
}

// CreateCommentForm 发表评论
type CreateCommentForm struct {
	EntityType  string     `form:"entityType"`
	EntityId    int64      `form:"entityId"`
	Content     string     `form:"content"`
	ImageList   []ImageDTO `form:"imageList"`
	QuoteId     int64      `form:"quoteId"`
	ContentType string     `form:"contentType"`
	UserAgent   string     `form:"userAgent"`
	Ip          string     `form:"ip"`
}
