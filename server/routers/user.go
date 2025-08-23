package routers

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// UserRouter 注册用户相关路由
func UserRouter(router *gin.RouterGroup) {
	userApi := api.UserApi{}
	// 公开路由
	publicRouter := router.Group("users")
	{
		publicRouter.POST("register", userApi.Register)
		publicRouter.POST("login", userApi.Login)
		publicRouter.GET("captcha", userApi.GetCaptcha)       // 图片验证码
		publicRouter.GET("email/code", userApi.SendEmailCode) // 发送邮箱验证码
		publicRouter.POST("forgot", userApi.ForgotPassword)   // 忘记密码
		publicRouter.POST("reset", userApi.ResetPassword)     // 重置密码
		publicRouter.GET(":id", userApi.GetUserById)          // 根据ID获取用户信息
	}

	// 需认证路由
	authRouter := router.Group("users").Use(middleware.InitJWT())
	{
		authRouter.GET("info", userApi.GetUserInfo)
		authRouter.PUT("update", userApi.UpdateUserInfo)
		authRouter.PUT("password", userApi.ChangePassword)
		authRouter.DELETE("delete", userApi.DeleteUser)
		authRouter.GET("list", userApi.GetUserList)
		authRouter.POST("avatar", userApi.UploadAvatar)      // 上传头像
		authRouter.PUT(":uuid/approve", userApi.ApproveUser) // 启用用户
		authRouter.PUT(":uuid/reject", userApi.RejectUser)   // 禁用用户
	}
}
