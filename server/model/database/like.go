package database

// Like 点赞模型
type Like struct {
	BaseModel                 // 嵌入基础模型
	UserID    uint            `gorm:"index;not null" json:"user_id"`
	ArticleID uint            `gorm:"index;not null" json:"article_id"`
}

// TableName 设置表名
func (Like) TableName() string {
	return "likes"
}