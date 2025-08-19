package database

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
