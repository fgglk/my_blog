package response

import (
	"server/model/appType"
	"server/model/database"
	"time" // 添加 time 包导入
)

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint               `json:"id"`
	UUID      string             `json:"uuid"`
	Username  string             `json:"username"`
	Nickname  string             `json:"nickname"`
	Email     string             `json:"email"`
	Avatar    string             `json:"avatar"`
	Bio       string             `json:"bio"`
	Address   string             `json:"address"`
	Role      appType.RoleType   `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// 转换数据库用户模型为响应模型
func ToUserResponse(user database.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		UUID:      user.UUID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Address:   user.Address,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}


// UserListResponse 用户列表响应
type UserListResponse struct {
	List  []UserResponse `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}
