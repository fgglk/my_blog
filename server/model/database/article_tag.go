package database

import (
	"time"

	"gorm.io/gorm"
)

// ArticleTag 文章与标签的多对多关联表
// 注意: 多对多关联表不应包含自增ID字段，仅需联合主键
type ArticleTag struct {
	ArticleID uint           `gorm:"primaryKey;index" json:"article_id"` // 文章ID
	TagID     uint           `gorm:"primaryKey;index" json:"tag_id"`     // 标签ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 自定义表名
func (ArticleTag) TableName() string {
	return "article_tags"
}
