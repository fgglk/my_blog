package appType

// LoginType 登录方式类型
type LoginType string

// 登录方式常量
const (
	LoginTypePassword LoginType = "password" // 密码登录
	LoginTypeQQ       LoginType = "qq"       // QQ登录
	LoginTypeWeChat   LoginType = "wechat"   // 微信登录
)

// 支持的登录方式列表
var SupportedLoginTypes = []LoginType{
	LoginTypePassword,
	LoginTypeQQ,
	LoginTypeWeChat,
}
