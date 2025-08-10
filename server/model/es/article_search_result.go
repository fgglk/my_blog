package es

// ArticleSearchResult 文章搜索结果
type ArticleSearchResult struct {
	Articles  []ArticleES `json:"articles"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	Size      int         `json:"size"`
	TotalPage int         `json:"total_page"`
}
