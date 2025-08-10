package utils

import (
	"errors"
	"server/global"
	"server/model/appType"
	"server/model/database"

	"github.com/gin-gonic/gin"
)

// IsAdmin 检查用户是否为管理员
func IsAdmin(userID uint) bool {
	var user database.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return false
	}
	// 使用appType中定义的管理员角色
	return user.Role == appType.RoleAdmin
}

// GetUserID 从gin上下文中获取用户ID
func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("未登录或登录已过期")
	}

	userIdUint, ok := userID.(uint)
	if !ok {
		return 0, errors.New("用户ID格式错误")
	}

	return userIdUint, nil
}
