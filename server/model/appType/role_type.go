package appType

// RoleType 用户角色类型
type RoleType string

// 角色常量定义
const (
	RoleAdmin   RoleType = "admin"   // 管理员
	RoleUser    RoleType = "user"    // 普通用户
	RoleVisitor RoleType = "visitor" // 游客
)

// 所有角色列表
var AllRoles = []RoleType{
	RoleAdmin,
	RoleUser,
	RoleVisitor,
}

// 检查角色是否有效
func (r RoleType) IsValid() bool {
	for _, role := range AllRoles {
		if r == role {
			return true
		}
	}
	return false
}
