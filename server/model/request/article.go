package request

import (
	"github.com/go-playground/validator/v10"
)

// ArticleCreateRequest 文章创建请求结构体
type ArticleCreateRequest struct {
	Title        string   `json:"title" binding:"required,min=1,max=100"`
	Content      string   `json:"content" binding:"required"`
	Summary      string   `json:"summary" binding:"max=500"`
	CategoryID   uint     `json:"category_id"` // 固定分类ID
	AuthorID     uint     `json:"author_id"`
	Tags         []uint   `json:"tags"`        // 保留ID数组，可选
	TagNames     []string `json:"tag_names"`   // 新增标签名称数组，可选
	CoverImage   string   `json:"cover_image"` // 封面图片URL
	Status       uint8    `json:"status" binding:"oneof=0 1"`
	ViewCount    int      `json:"view_count"`
	CommentCount int      `json:"comment_count"`
	LikeCount    int      `json:"like_count"`
}

// ArticleUpdateRequest 更新文章请求结构体
type ArticleUpdateRequest struct {
	ID         uint     `json:"id" binding:"" comment:"文章ID"` // 移除required标签
	Title      string   `json:"title" binding:"omitempty,min=1,max=100" comment:"文章标题"`
	Content    string   `json:"content" binding:"omitempty" comment:"文章内容"`
	CategoryID uint     `json:"category_id" binding:"omitempty" comment:"分类ID"`
	Tags       []uint   `json:"tags" binding:"omitempty" comment:"标签ID列表"`
	TagNames   []string `json:"tag_names" binding:"omitempty" comment:"标签名称列表"` // 新增标签名称字段
	CoverImage string   `json:"cover_image" binding:"omitempty" comment:"封面图片URL"`
	Summary    string   `json:"summary" binding:"omitempty" comment:"文章摘要"`
	Status     uint8    `json:"status" binding:"omitempty" comment:"文章状态"`
}

// ArticleQueryRequest 文章查询请求结构体
type ArticleQueryRequest struct {
	Title      string `form:"title" binding:"omitempty" comment:"文章标题模糊查询"`
	CategoryID uint   `form:"category_id" binding:"omitempty" comment:"分类ID"`
	TagID      uint   `form:"tag_id" binding:"omitempty" comment:"标签ID"`
	Page       int    `form:"page" binding:"omitempty,min=1" comment:"页码"`           // 修改：移除required
	Size       int    `form:"size" binding:"omitempty,min=1,max=100" comment:"每页条数"` // 修改：移除required
	Status     uint8  `form:"status" binding:"omitempty" comment:"文章状态"`
	AuthorID   uint   `json:"author_id" binding:"omitempty" comment:"作者ID"` // 修改：移除required
}

// ArticleIDRequest 文章ID请求结构体
type ArticleIDRequest struct {
	ID uint `json:"id" binding:"required" comment:"文章ID"`
}

// ToggleLikeRequest 切换点赞状态请求模型
type ToggleLikeRequest struct {
	ArticleID uint `json:"article_id" binding:"required,min=1"`
}

type ToggleFavoriteRequest struct {
	ArticleID uint `json:"article_id" binding:"required,min=1"`
}

// 自定义验证器示例
func (a *ArticleCreateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

// SearchArticleRequest 文章搜索请求结构体
type SearchArticleRequest struct {
	Keyword    string `form:"keyword" binding:"omitempty" comment:"搜索关键词"` // 移除required,min=1
	Page       int    `form:"page" binding:"omitempty,min=1" comment:"页码"`
	Size       int    `form:"size" binding:"omitempty,min=1,max=100" comment:"每页条数"`
	CategoryID uint   `form:"category_id" binding:"omitempty" comment:"分类ID筛选"`
	Tag        string `form:"tag" binding:"omitempty" comment:"标签筛选"`
	Sort       string `form:"sort" binding:"omitempty,oneof=time view comment like"`
	Order      string `form:"order" binding:"omitempty,oneof=asc desc"`
}
