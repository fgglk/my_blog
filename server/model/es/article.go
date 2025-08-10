package es

import (
	"server/model/appType"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// ArticleES 文章搜索文档结构
type ArticleES struct {
	ID            uint64                `json:"id"`
	Title         string                `json:"title"`
	Content       string                `json:"content"`
	Summary       string                `json:"summary"`
	Tags          []string              `json:"tags"`
	CategoryID    uint                  `json:"category_id"`
	UserID        uint                  `json:"user_id"`
	Status        appType.ArticleStatus `json:"status"`
	ViewCount     int                   `json:"view_count"`
	LikeCount     int                   `json:"like_count"`
	CommentCount  int                   `json:"comment_count"`
	FavoriteCount int                   `json:"favorite_count"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// IndexName 返回索引名称
func (a *ArticleES) IndexName() string {
	return "articles"
}

// GetMapping 获取文章索引映射
func GetMapping() *types.TypeMapping {
	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"id": types.LongNumberProperty{},
			"title": types.TextProperty{
				Fields: map[string]types.Property{
					"keyword": types.KeywordProperty{},
				},
			},
			"content":        types.TextProperty{},
			"summary":        types.TextProperty{},
			"tags":           types.KeywordProperty{},
			"category_id":    types.IntegerNumberProperty{},
			"user_id":        types.IntegerNumberProperty{},
			"status":         types.IntegerNumberProperty{},
			"view_count":     types.IntegerNumberProperty{},
			"like_count":     types.IntegerNumberProperty{},
			"comment_count":  types.IntegerNumberProperty{},
			"favorite_count": types.IntegerNumberProperty{},
			"created_at":     types.DateProperty{},
			"updated_at":     types.DateProperty{},
		},
	}
}
