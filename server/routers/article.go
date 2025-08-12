package routers

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// 文章路由
func ArticleRouter(Router *gin.RouterGroup) {
	articleRouter := Router.Group("articles")
	{
		// 公开路由
		articleRouter.GET("", (&api.ArticleApi{}).GetArticleList)
		articleRouter.GET("/search", (&api.ArticleApi{}).SearchArticles)
		articleRouter.GET("/stats", (&api.ArticleApi{}).GetWebsiteStats)
		articleRouter.GET("/:id/related", (&api.ArticleApi{}).GetRelatedArticles)

		// 需要认证的路由（包括可选的认证）
		authArticleRouter := articleRouter.Use(middleware.InitJWT())
		{
			authArticleRouter.GET("/:id", (&api.ArticleApi{}).GetArticle)
			authArticleRouter.GET("/my", (&api.ArticleApi{}).GetUserArticles)
			authArticleRouter.POST("", (&api.ArticleApi{}).CreateArticle)
			authArticleRouter.PUT("/:id", (&api.ArticleApi{}).UpdateArticle)
			authArticleRouter.DELETE("/:id", (&api.ArticleApi{}).DeleteArticle)
			authArticleRouter.POST("/like", (&api.ArticleApi{}).ToggleLike)
			authArticleRouter.POST("/favorite", (&api.ArticleApi{}).ToggleFavorite)
		}
	}
}
