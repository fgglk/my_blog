package routers

import (
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// 统一路由注册入口
// SetupRouter 初始化路由
func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), middleware.InitLogger(), middleware.InitCors())

	// 静态文件服务
	router.Static("/uploads", "./uploads")

	// 公开路由组
	publicGroup := router.Group("/api")
	{
		// 注册文章路由
		ArticleRouter(publicGroup)
		// 注册用户路由
		UserRouter(publicGroup)
		// 注册图片路由
		ImageRouter(publicGroup)
		// 注册评论路由
		CommentRouter(publicGroup)
		// 注册页面路由
		PageRouter(publicGroup)
		// 注册分类路由
		CategoryRouter(publicGroup)
		// 注册标签路由
		TagRouter(publicGroup)
	}

	return router
}
