package flag

import (
	"errors"
	"fmt"
	"syscall"

	"golang.org/x/term"

	"server/global"
	"server/model/appType"
	"server/model/database"
	"server/utils"
	"go.uber.org/zap"
)

// CreateAdministrator 初始化系统管理员账户
func CreateAdministrator() error {
	admin := database.User{}

	// 获取管理员邮箱
	fmt.Print("管理员邮箱: ")
	var adminEmail string
	if _, err := fmt.Scan(&adminEmail); err != nil {
		errMsg := "邮箱输入错误: " + err.Error()
		global.ZapLog.Error(errMsg)
		return errors.New(errMsg)
	}
	admin.Email = adminEmail

	// 配置终端为原始模式（无回显）
	stdin := int(syscall.Stdin)
	oldState, err := term.MakeRaw(stdin)
	if err != nil {
		return err
	}
	defer term.Restore(stdin, oldState)

	// 读取密码
	pwd, err := inputSecurely("设置密码: ")
	if err != nil {
		return err
	}

	// 确认密码
	confirmPwd, err := inputSecurely("确认密码: ")
	if err != nil {
		return err
	}

	// 验证密码
	if pwd != confirmPwd {
		errMsg := "两次密码输入不一致"
		global.ZapLog.Error(errMsg)
		return errors.New(errMsg)
	}
	if len(pwd) < 8 || len(pwd) > 20 {
		errMsg := "密码长度必须在8-20个字符之间"
		global.ZapLog.Error(errMsg)
		return errors.New(errMsg)
	}

	// 构建管理员信息
	admin.UUID = utils.GenerateUUID()
	admin.Username = global.Config.Website.Name
	admin.Password = utils.BcryptHash(pwd)
	admin.Role = appType.RoleAdmin
	admin.Avatar = "image/admin.jpg"
	admin.Address = global.Config.Website.Address

	// 保存到数据库
	if err := global.DB.Create(&admin).Error; err != nil {
		errMsg := "创建管理员失败: " + err.Error()
		global.ZapLog.Error(errMsg)
		return errors.New(errMsg)
	}

	global.ZapLog.Info("管理员账户创建成功", zap.String("email", admin.Email), zap.String("username", admin.Username))
	return nil
}

// inputSecurely 安全地读取用户输入（无回显）
func inputSecurely(prompt string) (string, error) {
	fmt.Print(prompt)
	pwdBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", errors.New("读取输入失败: " + err.Error())
	}
	return string(pwdBytes), nil
}
