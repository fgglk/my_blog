package database

type Favorite struct {
  BaseModel
  UserID    uint `gorm:"primaryKey;index" json:"user_id"`
  ArticleID uint `gorm:"primaryKey;index" json:"article_id"`
}

// TableName 设置表名
func (Favorite) TableName() string {
  return "favorites"
}