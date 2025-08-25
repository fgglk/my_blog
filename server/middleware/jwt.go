package middleware

import (
	"fmt"
	"net/http"
	"server/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 初始化JWT中间件
func InitJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 对于GET请求，允许未登录用户访问，但不设置用户信息
			if c.Request.Method == "GET" {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未提供token"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			// 对于GET请求，允许未登录用户访问，但不设置用户信息
			if c.Request.Method == "GET" {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token格式错误"})
			c.Abort()
			return
		}

		// 解析token - 使用utils包中的ParseToken函数
		claims, err := utils.ParseToken(parts[1], false)
		if err != nil {
			// 对于GET请求，允许未登录用户访问，但不设置用户信息
			if c.Request.Method == "GET" {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "无效的token"})
			c.Abort()
			return
		}

		// 检查token是否过期
		if claims.ExpiresAt.Before(time.Now()) {
			// 对于GET请求，允许未登录用户访问，但不设置用户信息
			if c.Request.Method == "GET" {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token已过期"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// 添加调试信息
		if c.Request.Method == "POST" && c.Request.URL.Path == "/api/articles" {
			fmt.Printf("JWT中间件: 用户ID=%d, 用户名=%s\n", claims.UserID, claims.Username)
		}

		c.Next()
	}
}
