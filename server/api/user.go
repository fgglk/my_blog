package api

import (
	"strconv"
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserApi 用户API
type UserApi struct{}

// 定义用户服务变量，通过ServiceGroups调用
var userService = service.ServiceGroups.UserService

// Register 用户注册
func (u *UserApi) Register(c *gin.Context) {
	var registerReq request.RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 参数验证
	if errMsg := utils.ValidateStruct(registerReq); errMsg != "" {
		response.FailWithMessage(errMsg, c)
		return
	}

	// 验证图片验证码
	if !utils.VerifyCaptcha(registerReq.CaptchaID, registerReq.CaptchaCode) {
		response.FailWithMessage("图片验证码错误或已过期", c)
		return
	}

	// 验证邮箱验证码
	if !utils.VerifyEmailCode(registerReq.Email, registerReq.EmailCode) {
		response.FailWithMessage("邮箱验证码错误或已过期", c)
		return
	}

	// 调用服务层
	if err, user := userService.Register(registerReq); err != nil {
		global.ZapLog.Error("注册失败", zap.Error(err))
		response.FailWithMessage("注册失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(response.ToUserResponse(user), "注册成功", c)
	}
}

// Login 用户登录
func (u *UserApi) Login(c *gin.Context) {
	var loginReq request.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 参数验证
	if errMsg := utils.ValidateStruct(loginReq); errMsg != "" {
		response.FailWithMessage(errMsg, c)
		return
	}

	// 验证图片验证码
	if !utils.VerifyCaptcha(loginReq.CaptchaID, loginReq.CaptchaCode) {
		response.FailWithMessage("图片验证码错误或已过期", c)
		return
	}

	// 调用服务层(内部实现登录类型自动判断)
	if err, user := userService.Login(loginReq); err != nil {
		response.FailWithMessage("登录失败: "+err.Error(), c)
	} else {
		// 生成JWT令牌
		token, err := utils.GenerateToken(user.ID, user.Username)
		if err != nil {
			response.FailWithMessage("生成令牌失败", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:  response.ToUserResponse(user),
			Token: token,
		}, "登录成功", c)
	}
}

// GetUserInfo 获取用户信息
func (u *UserApi) GetUserInfo(c *gin.Context) {
	userId, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	if err, user := userService.GetUserInfo(userId); err != nil {
		response.FailWithMessage("获取用户信息失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(response.ToUserResponse(user), "获取成功", c)
	}
}

// GetUserById 根据ID获取用户信息
func (u *UserApi) GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的用户ID", c)
		return
	}

	if err, user := userService.GetUserInfo(uint(id)); err != nil {
		response.FailWithMessage("获取用户信息失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(response.ToUserResponse(user), "获取成功", c)
	}
}

// UpdateUserInfo 更新用户信息
func (u *UserApi) UpdateUserInfo(c *gin.Context) {
	var updateReq request.UserUpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 参数验证
	if errMsg := utils.ValidateStruct(updateReq); errMsg != "" {
		response.FailWithMessage(errMsg, c)
		return
	}

	userId, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	// 只有管理员或本人可以更新信息
	if userId != updateReq.ID && !utils.IsAdmin(userId) {
		response.FailWithMessage("没有权限进行此操作", c)
		return
	}

	// 新增：非管理员用户不能修改Role字段
	if !utils.IsAdmin(userId) {
		// 获取原始用户信息
		if err, originalUser := userService.GetUserInfo(updateReq.ID); err != nil {
			response.FailWithMessage("获取用户信息失败", c)
			return
		} else {
			// 保留原始角色
			updateReq.Role = originalUser.Role
		}
	}

	if err, user := userService.UpdateUserInfo(updateReq); err != nil {
		response.FailWithMessage("更新用户信息失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(response.ToUserResponse(user), "更新成功", c)
	}
}

// ChangePassword 修改密码
func (u *UserApi) ChangePassword(c *gin.Context) {
	var passwordReq request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&passwordReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 参数验证
	if errMsg := utils.ValidateStruct(passwordReq); errMsg != "" {
		response.FailWithMessage(errMsg, c)
		return
	}

	// 新增：检查新密码和确认密码是否一致
	if passwordReq.NewPassword != passwordReq.ConfirmPassword {
		response.FailWithMessage("新密码和确认密码不一致", c)
		return
	}

	userId, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	if err := userService.ChangePassword(userId, passwordReq); err != nil {
		response.FailWithMessage("修改密码失败: "+err.Error(), c)
	} else {
		response.OkWithMessage("修改密码成功", c)
	}
}

// DeleteUser 删除用户
func (u *UserApi) DeleteUser(c *gin.Context) {
	var idReq request.IdRequest
	if err := c.ShouldBindJSON(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	// 只有管理员或用户本人可以删除账户
	if !utils.IsAdmin(userId) && userId != idReq.ID {
		response.FailWithMessage("没有权限进行此操作", c)
		return
	}

	if err := userService.DeleteUser(idReq.ID); err != nil {
		response.FailWithMessage("删除用户失败: "+err.Error(), c)
	} else {
		response.OkWithMessage("删除用户成功", c)
	}
}

// GetUserList 获取用户列表
func (u *UserApi) GetUserList(c *gin.Context) {
	var listReq request.UserListRequest
	if err := c.ShouldBindQuery(&listReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if listReq.Page <= 0 {
		listReq.Page = 1
	}
	if listReq.Size <= 0 || listReq.Size > 100 {
		listReq.Size = 10
	}
	// 参数验证
	if errMsg := utils.ValidateStruct(listReq); errMsg != "" {
		response.FailWithMessage(errMsg, c)
		return
	}

	userId, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	// 只有管理员可以获取用户列表
	if !utils.IsAdmin(userId) {
		response.FailWithMessage("没有权限进行此操作", c)
		return
	}

	if err, list, total := userService.GetUserList(listReq); err != nil {
		response.FailWithMessage("获取用户列表失败: "+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     listReq.Page,
			PageSize: listReq.Size,
		}, "获取成功", c)
	}
}

// GetCaptcha 获取图片验证码
func (u *UserApi) GetCaptcha(c *gin.Context) {
	// 生成验证码
	captchaID, b64s, err := utils.GenerateCaptcha()
	if err != nil {
		response.FailWithMessage("验证码生成失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(map[string]string{
		"captcha_id": captchaID,
		"image":      b64s,
	}, "获取验证码成功", c)
}

// SendEmailCode 发送邮箱验证码
func (u *UserApi) SendEmailCode(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.FailWithMessage("邮箱不能为空", c)
		return
	}

	// 生成6位数字验证码
	code := utils.GenerateEmailCode()

	// 存储验证码到Redis
	if err := utils.StoreEmailCodeInRedis(email, code); err != nil {
		response.FailWithMessage("验证码存储失败: "+err.Error(), c)
		return
	}

	// 发送邮件
	if err := utils.SendEmailCode(email, code, "注册验证码", "注册"); err != nil {
		response.FailWithMessage("邮件发送失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("验证码已发送至邮箱，请注意查收", c)
}

// ForgotPassword 处理忘记密码请求
func (u *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ZapLog.Error("请求参数错误", zap.Error(err))
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 验证图片验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		response.FailWithMessage("验证码错误或已过期", c)
		return
	}

	// 检查用户是否存在
	user, err := userService.FindUserByEmail(req.Email)
	if err != nil || user.ID == 0 {
		response.FailWithMessage("该邮箱未注册", c)
		return
	}

	// 生成邮箱验证码
	code := utils.GenerateEmailCode()

	// 存储验证码到Redis
	if err := utils.StoreEmailCodeInRedis(req.Email, code); err != nil {
		global.ZapLog.Error("存储验证码失败", zap.Error(err))
		response.FailWithMessage("发送验证码失败，请重试", c)
		return
	}

	// 发送邮箱验证码
	if err := utils.SendEmailCode(req.Email, code, "重置密码验证码", "重置密码"); err != nil {
		global.ZapLog.Error("发送邮件失败", zap.Error(err))
		response.FailWithMessage("发送验证码失败，请重试", c)
		return
	}

	response.OkWithMessage("验证码已发送至您的邮箱，请注意查收", c)
}

// ResetPassword 处理重置密码请求
func (u *UserApi) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ZapLog.Error("请求参数错误", zap.Error(err))
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 验证邮箱验证码
	if !utils.VerifyEmailCode(req.Email, req.EmailCode) {
		response.FailWithMessage("邮箱验证码错误或已过期", c)
		return
	}

	// 更新用户密码
	if err := userService.ResetUserPassword(req.Email, req.NewPassword); err != nil {
		global.ZapLog.Error("重置密码失败", zap.Error(err))
		response.FailWithMessage("重置密码失败，请重试", c)
		return
	}

	response.OkWithMessage("密码重置成功，请使用新密码登录", c)
}

// UploadAvatar 上传头像
func (u *UserApi) UploadAvatar(c *gin.Context) {
	// 获取当前用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	_, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.FailWithMessage("获取头像文件失败: "+err.Error(), c)
		return
	}

	// 检查文件类型
	if !utils.IsImageType(header.Header.Get("Content-Type")) {
		response.FailWithMessage("不支持的图片格式", c)
		return
	}

	// 调用Service层上传头像
	avatarURL, err := userService.UploadAvatar(header, userID)
	if err != nil {
		response.FailWithMessage("头像上传失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"url": avatarURL,
	}, "头像上传成功", c)
}
