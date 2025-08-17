package response

import (
	"time"

	"server/global"
	"server/model/database"
)

// 文章详情数据结构
type ArticleResponse struct {
	ID            uint           `json:"id"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Summary       string         `json:"summary"`
	CategoryID    uint           `json:"category_id"`
	Category      CategorySimple `json:"category"`
	Tags          []TagResponse  `json:"tags"`
	AuthorID      uint           `json:"author_id"`
	AuthorName    string         `json:"author_name"`
	AuthorAvatar  string         `json:"author_avatar"`
	CoverImage    string         `json:"cover_image"`
	ViewCount     int            `json:"view_count"`
	LikeCount     int            `json:"like_count"`
	CommentCount  int            `json:"comment_count"`
	FavoriteCount int            `json:"favorite_count"`
	IsPublished   bool           `json:"is_published"`
	IsLiked       bool           `json:"is_liked"`
	IsFavorited   bool           `json:"is_favorited"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// 文章列表数据结构
type ArticleListResponse struct {
	List      []ArticleResponse `json:"list"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	Size      int               `json:"size"`
	TotalPage int               `json:"total_page"`
}

// 分类简要信息
type CategorySimple struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 标签信息
type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 分类响应（带文章数量）
type CategoryWithCountResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int64  `json:"article_count"`
}

// 标签响应（带文章数量）
type TagWithCountResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int64  `json:"article_count"`
}

// 分类响应
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// 数据库模型转换为响应模型
func ToArticleResponse(article database.Article, category database.Category, tags []database.Tag, authorName string, currentUserID uint) ArticleResponse {
	var tagResponses []TagResponse
	for _, tag := range tags {
		tagResponses = append(tagResponses, TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	// 检查用户是否已点赞
	isLiked := false
	if currentUserID > 0 {
		var like database.Like
		if global.DB.Where("article_id = ? AND user_id = ? AND deleted_at IS NULL", article.ID, currentUserID).First(&like).Error == nil {
			isLiked = true
		}
	}

	// 检查用户是否已收藏
	isFavorited := false
	if currentUserID > 0 {
		var favorite database.Favorite
		if global.DB.Where("article_id = ? AND user_id = ? AND deleted_at IS NULL", article.ID, currentUserID).First(&favorite).Error == nil {
			isFavorited = true
		}
	}

	return ArticleResponse{
		ID:         article.ID,
		Title:      article.Title,
		Content:    article.Content,
		Summary:    article.Summary,
		CategoryID: article.CategoryID,
		Category: CategorySimple{
			ID:   category.ID,
			Name: category.Name,
		},
		Tags:          tagResponses,
		AuthorID:      article.AuthorID,
		AuthorName:    authorName,
		AuthorAvatar:  article.Author.Avatar,
		CoverImage:    article.CoverImage,
		ViewCount:     article.ViewCount,
		LikeCount:     article.LikeCount,
		CommentCount:  article.CommentCount,
		FavoriteCount: article.FavoriteCount,
		IsPublished:   article.Status == 1,
		IsLiked:       isLiked,
		IsFavorited:   isFavorited,
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
	}
}

// 收藏记录响应
type FavoriteResponse struct {
	ID        uint            `json:"id"`
	ArticleID uint            `json:"article_id"`
	UserID    uint            `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Article   ArticleResponse `json:"article"`
}

// 收藏列表响应
type FavoriteListResponse struct {
	List      []FavoriteResponse `json:"list"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	Size      int                `json:"size"`
	TotalPage int                `json:"total_page"`
}

// ToCategoryWithCountResponse 转换为带计数的分类响应
func ToCategoryWithCountResponse(category database.CategoryWithCount) CategoryWithCountResponse {
	return CategoryWithCountResponse{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		ArticleCount: category.ArticleCount,
	}
}

// ToTagWithCountResponse 转换为带计数的标签响应
func ToTagWithCountResponse(tag database.TagWithCount) TagWithCountResponse {
	return TagWithCountResponse{
		ID:           tag.ID,
		Name:         tag.Name,
		Slug:         tag.Slug,
		ArticleCount: tag.ArticleCount,
	}
}

// ToCategoryResponse 转换为分类响应
func ToCategoryResponse(category database.Category) CategoryResponse {
	return CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

// ToTagResponse 转换为标签响应
func ToTagResponse(tag database.Tag) TagResponse {
	return TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

// ToFavoriteResponse 转换为收藏响应
func ToFavoriteResponse(favorite database.Favorite, currentUserID uint) FavoriteResponse {
	return FavoriteResponse{
		ID:        favorite.ID,
		ArticleID: favorite.ArticleID,
		UserID:    favorite.UserID,
		CreatedAt: favorite.CreatedAt,
		UpdatedAt: favorite.UpdatedAt,
		Article:   ToArticleResponse(favorite.Article, favorite.Article.Category, favorite.Article.Tags, favorite.Article.Author.Username, currentUserID),
	}
}
