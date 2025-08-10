package database

type Category struct {
	BaseModel
	Name     string `gorm:"size:50;uniqueIndex;not null" json:"name"` // 分类名称
	Slug     string `gorm:"size:100;uniqueIndex" json:"slug"`         // URL友好名称
	ParentID *uint  `gorm:"index" json:"parent_id,omitempty"`         // 父分类ID(支持多级分类)
	Sort     int    `gorm:"default:0" json:"sort"`                    // 排序

	// 关联
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// CategoryWithCount 带文章数量的分类
type CategoryWithCount struct {
	Category
	ArticleCount int64 `json:"article_count"` // 文章数量
}
