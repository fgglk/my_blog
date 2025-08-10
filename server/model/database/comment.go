package database

import (
	"server/model/appType"
)

// Comment 评论模型
type Comment struct {
	BaseModel                               // 嵌入基础模型(包含ID、CreatedAt、UpdatedAt、DeletedAt)
	ArticleID     uint                      `gorm:"index;not null" json:"article_id"`  // 文章ID
	UserID        uint                      `gorm:"index;not null" json:"user_id"`     // 用户ID
	ParentID      *uint                     `gorm:"index" json:"parent_id,omitempty"`  // 父评论ID(支持回复)
	Content       string                    `gorm:"type:text;not null" json:"content"` // 评论内容
	CommentStatus appType.CommentStatusType `gorm:"type:tinyint;default:0;comment:'评论状态：0-待审核，1-已发布，2-已拒绝'" json:"comment_status"`
	User          User                      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article       Article                   `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Replies       []Comment                 `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Children      []Comment                 `gorm:"-" json:"children,omitempty"`        // 子评论（不存储在数据库中）
}
