package routers

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// CommentRouter 注册评论相关路由
func CommentRouter(router *gin.RouterGroup) {
	commentRouter := router.Group("comments")
	{
		// 公开路由
		commentRouter.GET("", (&api.CommentApi{}).GetCommentList) // 获取评论列表
		commentRouter.GET("/:id", (&api.CommentApi{}).GetComment) // 获取单个评论

		// 需认证路由
		authRouter := commentRouter.Use(middleware.InitJWT())
		{
			authRouter.POST("", (&api.CommentApi{}).CreateComment)            // 创建评论
			authRouter.PUT("/:id", (&api.CommentApi{}).UpdateComment)         // 更新评论
			authRouter.DELETE("/:id", (&api.CommentApi{}).DeleteComment)      // 删除评论
			authRouter.POST("/:id/reply", (&api.CommentApi{}).ReplyToComment) // 回复评论
		}
	}
}
