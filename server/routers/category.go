package routers

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

// 注册分类相关路由
func CategoryRouter(Router *gin.RouterGroup) {
	categoryRouter := Router.Group("categories")
	{
		// 前台路由（无需认证）
		categoryRouter.GET("", (&api.CategoryApi{}).GetCategoryList)     // 获取分类列表
		categoryRouter.GET("/:id", (&api.CategoryApi{}).GetCategory)     // 获取分类详情
	}
} 