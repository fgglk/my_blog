package database

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含所有表的公共字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 带状态的基础模型
type BaseModelWithStatus struct {
	BaseModel
	Status uint8 `gorm:"default:1" json:"status"` // 状态(0-禁用,1-正常)
}