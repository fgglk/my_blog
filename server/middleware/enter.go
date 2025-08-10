package middleware

import (
	"github.com/gin-gonic/gin"
)

// 所有中间件的集合，便于统一导入
var (
	Cors   gin.HandlerFunc
	Logger gin.HandlerFunc
	JWT    gin.HandlerFunc
)

// 初始化所有中间件
func Init() {
	Cors = InitCors()
	Logger = InitLogger()
	JWT = InitJWT()
}
