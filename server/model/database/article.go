package database

// Article 文章模型
type Article struct {
	BaseModelWithStatus        // 嵌入带状态的基础模型
	Title               string `gorm:"size:200;not null" json:"title"`
	Slug                string `gorm:"size:255;uniqueIndex" json:"slug"`
	Content             string `gorm:"type:longtext;not null" json:"content"`
	Summary             string `gorm:"type:text" json:"summary"`
	CoverImage          string `gorm:"size:255" json:"cover_image"`
	AuthorID            uint   `gorm:"index;not null" json:"author_id"`
	CategoryID          uint   `gorm:"index;not null" json:"category_id"`
	ViewCount           int    `gorm:"default:0" json:"view_count"`
	CommentCount        int    `gorm:"default:0" json:"comment_count"`
	LikeCount           int    `gorm:"default:0" json:"like_count"`
	FavoriteCount       int    `gorm:"default:0" json:"favorite_count"`

	// 关联
	Author   User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag    `gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:ArticleID;References:ID;joinReferences:TagID" json:"tags,omitempty"`
}
