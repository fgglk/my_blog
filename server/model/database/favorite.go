package database

// Favorite 收藏模型
// 设计说明：
// 1. ArticleID: 外键，关联文章表
// 2. Article: 关联字段，用于预加载文章信息，避免N+1查询
// 3. 使用 omitempty 标签，在不需要文章详情时不会序列化该字段
type Favorite struct {
	BaseModel
	UserID    uint    `gorm:"primaryKey;index" json:"user_id"`
	ArticleID uint    `gorm:"primaryKey;index" json:"article_id"`
	Article   Article `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"article,omitempty"`
}

// TableName 设置表名
func (Favorite) TableName() string {
	return "favorites"
}

// GetArticleInfo 获取文章基本信息（轻量级）
func (f *Favorite) GetArticleInfo() map[string]interface{} {
	if f.Article.ID == 0 {
		return nil
	}

	return map[string]interface{}{
		"id":             f.Article.ID,
		"title":          f.Article.Title,
		"summary":        f.Article.Summary,
		"cover_image":    f.Article.CoverImage,
		"author_id":      f.Article.AuthorID,
		"view_count":     f.Article.ViewCount,
		"like_count":     f.Article.LikeCount,
		"comment_count":  f.Article.CommentCount,
		"favorite_count": f.Article.FavoriteCount,
		"created_at":     f.Article.CreatedAt,
	}
}
