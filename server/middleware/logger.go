package middleware

import (
	"server/global"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

// 初始化日志中间件
func InitLogger() gin.HandlerFunc {
	return ginzap.Ginzap(global.ZapLog, time.RFC3339, true)
}
