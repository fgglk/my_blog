package routers

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// 注册标签相关路由
func TagRouter(Router *gin.RouterGroup) {
	tagRouter := Router.Group("tags")
	{
		// 前台路由（无需认证）
		tagRouter.GET("", (&api.TagApi{}).GetTagList) // 获取标签列表
		tagRouter.GET("/:id", (&api.TagApi{}).GetTag) // 获取标签详情

		// 管理员路由（需要认证和管理员权限）
		tagRouter.DELETE("/cleanup", middleware.InitJWT(), (&api.TagApi{}).CleanupOrphanTags) // 清理孤儿标签
	}
}
