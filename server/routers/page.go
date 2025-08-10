package routers

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// 注册页面相关路由
func PageRouter(Router *gin.RouterGroup) {
	pageRouter := Router.Group("pages")
	{
		// 前台路由（无需认证）
		pageRouter.GET("/slug/:slug", (&api.PageApi{}).GetPageBySlug) // 通过slug获取页面
		pageRouter.GET("/nav", (&api.PageApi{}).GetNavPages)          // 获取导航页面

		// 后台路由（需要管理员认证）
		authPageRouter := pageRouter.Use(middleware.InitJWT())
		{
			authPageRouter.POST("", (&api.PageApi{}).CreatePage)       // 创建页面
			authPageRouter.GET("/:id", (&api.PageApi{}).GetPage)       // 获取单个页面
			authPageRouter.PUT("", (&api.PageApi{}).UpdatePage)        // 更新页面
			authPageRouter.DELETE("/:id", (&api.PageApi{}).DeletePage) // 删除页面
			authPageRouter.GET("", (&api.PageApi{}).ListPages)         // 分页查询页面
		}
	}
}
