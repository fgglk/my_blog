package database

import (
	"server/model/appType"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	BaseModelWithStatus                   // 嵌入带状态的基础模型
	UUID                string            `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	Username            string            `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email               string            `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password            string            `gorm:"size:100;not null" json:"-"`
	Nickname            string            `gorm:"size:50" json:"nickname"`
	Avatar              string            `gorm:"size:255" json:"avatar"`
	Bio                 string            `gorm:"type:text" json:"bio"`
	Address             string            `gorm:"size:255" json:"address"`
	Role                appType.RoleType  `gorm:"size:20;default:'user'" json:"role"`
	LoginMethod         appType.LoginType `gorm:"size:20;default:'password'" json:"login_method"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UUID = uuid.New().String()
	return nil
}
