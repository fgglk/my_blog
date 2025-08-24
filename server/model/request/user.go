package request

import (
	"server/model/appType"
	"time"
)

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=1,max=50"`
	Password    string `json:"password" validate:"required,min=6,max=20"`
	Nickname    string `json:"nickname" validate:"max=50"`
	Email       string `json:"email" validate:"required,email"`
	CaptchaID   string `json:"captcha_id" validate:"required"`
	CaptchaCode string `json:"captcha_code" validate:"required,len=6"`
	EmailCode   string `json:"email_code" validate:"required,len=6"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	CaptchaID   string `json:"captcha_id" validate:"required"`
	CaptchaCode string `json:"captcha_code" validate:"required,len=6"`
}

// UserUpdateRequest 用户更新信息请求
type UserUpdateRequest struct {
	ID        uint             `json:"id" validate:"required"`
	Nickname  string           `json:"nickname" validate:"max=50"`
	Email     string           `json:"email" validate:"email"`
	Avatar    string           `json:"avatar" validate:"omitempty,url"`
	Bio       string           `json:"bio" validate:"max=200"`
	Address   string           `json:"address" validate:"max=100"`
	Role      appType.RoleType `json:"role" validate:"omitempty"` // 改为可选
	Username  string           `json:"username" validate:"max=50"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required,min=6,max=20"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=20"`
}

// UserListRequest 用户列表查询请求
type UserListRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	Size     int    `form:"size" validate:"min=1,max=100"`
	Username string `form:"username" validate:"omitempty"`
	Email    string `form:"email" validate:"omitempty,email"`
	Status   *uint8 `form:"status" validate:"omitempty,oneof=0 1"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email     string `json:"email" binding:"required,email"`
	CaptchaID string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	EmailCode       string `json:"email_code" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// CreateUserRequest 管理员创建用户请求
type CreateUserRequest struct {
	Username string           `json:"username" validate:"required,min=1,max=50"`
	Password string           `json:"password" validate:"required,min=6,max=20"`
	Nickname string           `json:"nickname" validate:"max=50"`
	Email    string           `json:"email" validate:"required,email"`
	Role     appType.RoleType `json:"role" validate:"omitempty"`
	Bio      string           `json:"bio" validate:"max=200"`
	Address  string           `json:"address" validate:"max=100"`
}
