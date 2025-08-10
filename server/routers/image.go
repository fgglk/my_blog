package routers

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// ImageRouter 图片路由配置
func ImageRouter(Router *gin.RouterGroup) {
	imageRouter := Router.Group("image")
	{
		// 公开路由

		imageRouter.GET("show/:id", api.ShowImage) // 查看图片

		// 需要认证的路由
		authRouter := imageRouter.Use(middleware.InitJWT())
		{
			authRouter.POST("upload", api.UploadImage)       // 上传图片
			authRouter.POST("avatar", api.UploadAvatar)      // 上传头像
			authRouter.DELETE("delete/:id", api.DeleteImage) // 删除图片
			authRouter.PUT("update/:id", api.UpdateImage)    // 更新图片信息
			authRouter.GET("list", api.GetImageList)         // 获取图片列表
		}
	}
}
