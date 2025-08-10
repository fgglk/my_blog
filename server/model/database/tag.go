package database

import (
	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	BaseModel // 嵌入基础模型（ID, CreatedAt, UpdatedAt, DeletedAt）
	Name      string `gorm:"size:50;uniqueIndex;not null" json:"name"` // 标签名称
	Slug      string `gorm:"size:50;uniqueIndex;not null" json:"slug"` // URL友好名称
	Count     int    `gorm:"default:0" json:"count"`                   // 关联文章数量
}

// 标签状态常量
const (
	TagStatusNormal = 1 // 正常
	TagStatusHidden = 0 // 隐藏
)

// BeforeDelete 删除标签前清理关联关系
func (t *Tag) BeforeDelete(tx *gorm.DB) error {
	// 删除标签前先删除关联关系
	return tx.Where("tag_id = ?", t.ID).Delete(&ArticleTag{}).Error
}

// TagWithCount 带文章数量的标签
type TagWithCount struct {
	Tag
	ArticleCount int64 `json:"article_count"` // 文章数量
}
