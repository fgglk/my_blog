package service

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/utils"
	"strings"
	"time"

	"go.uber.org/zap"
)

type UserService struct{}

// Register 用户注册
func (u *UserService) Register(registerReq request.RegisterRequest) (err error, user database.User) {
	// 开始事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户名是否已存在
	var count int64
	if err := tx.Model(&database.User{}).Where("username = ? AND deleted_at IS NULL", registerReq.Username).Count(&count).Error; err != nil {
		tx.Rollback()
		return err, user
	}
	if count > 0 {
		tx.Rollback()
		return errors.New("用户名已存在"), user
	}

	// 检查邮箱是否已存在
	if err := tx.Model(&database.User{}).Where("email = ?", registerReq.Email).Count(&count).Error; err != nil {
		tx.Rollback()
		return err, user
	}
	if count > 0 {
		tx.Rollback()
		return errors.New("邮箱已被注册"), user
	}

	// 创建用户
	user = database.User{
		Username: registerReq.Username,
		Password: utils.BcryptHash(registerReq.Password),
		Nickname: registerReq.Nickname,
		Email:    registerReq.Email,
	}

	if err = tx.Create(&user).Error; err != nil {
		tx.Rollback()
		// 明确处理唯一约束错误
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "idx_users_username") {
				return errors.New("用户名已存在"), user
			} else if strings.Contains(err.Error(), "idx_users_email") {
				return errors.New("邮箱已被注册"), user
			}
		}
		return err, user
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err, user
	}

	return nil, user
}

// Login 用户登录
func (u *UserService) Login(loginReq request.LoginRequest) (err error, user database.User) {
	// 查询用户
	if err = global.DB.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		return errors.New("用户不存在"), user
	}

	// 验证密码
	if !utils.BcryptCheck(loginReq.Password, user.Password) {
		return errors.New("密码错误"), user
	}

	// 检查用户状态
	if user.Status == 0 {
		return errors.New("用户已被禁用，请联系管理员"), user
	}

	// 更新最后登录时间
	now := time.Now()
	if err = global.DB.Model(&user).Update("last_login_at", now).Error; err != nil {
		// 即使更新失败也不影响登录，只记录错误
		global.ZapLog.Error("更新最后登录时间失败", zap.Error(err))
	} else {
		user.LastLoginAt = &now
	}

	return nil, user
}

// GetUserInfo 获取用户信息
func (u *UserService) GetUserInfo(id uint) (err error, user database.User) {
	if err = global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return errors.New("用户不存在"), user
	}
	return nil, user
}

// UpdateUserInfo 更新用户信息
func (u *UserService) UpdateUserInfo(updateReq request.UserUpdateRequest) (err error, user database.User) {
	// 检查用户是否存在
	if err = global.DB.Where("id = ?", updateReq.ID).First(&user).Error; err != nil {
		return errors.New("用户不存在"), user
	}

	// 更新用户信息
	updateMap := make(map[string]interface{})
	if updateReq.Nickname != "" {
		updateMap["nickname"] = updateReq.Nickname
	}
	if updateReq.Email != "" {
		updateMap["email"] = updateReq.Email
	}
	if updateReq.Avatar != "" {
		updateMap["avatar"] = updateReq.Avatar
	}
	if updateReq.Bio != "" {
		updateMap["bio"] = updateReq.Bio
	}
	if updateReq.Address != "" {
		updateMap["address"] = updateReq.Address
	}
	// 更新用户信息方法中，修改角色相关代码
	if updateReq.Role != "" {
		// 验证角色是否有效
		if !updateReq.Role.IsValid() {
			return errors.New("无效的角色类型"), user
		}
		updateMap["role"] = updateReq.Role
	}
	if updateReq.Username != "" {
		updateMap["username"] = updateReq.Username
	}
	updateMap["updated_at"] = time.Now()

	if err = global.DB.Model(&user).Updates(updateMap).Error; err != nil {
		return err, user
	}

	// 重新获取用户信息
	if err = global.DB.Where("id = ?", updateReq.ID).First(&user).Error; err != nil {
		return err, user
	}

	return nil, user
}

// ChangePassword 修改密码
func (u *UserService) ChangePassword(id uint, passwordReq request.ChangePasswordRequest) (err error) {
	var user database.User
	if err = global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !utils.BcryptCheck(passwordReq.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	// 更新新密码
	user.Password = utils.BcryptHash(passwordReq.NewPassword)
	if err = global.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// DeleteUser 删除用户
func (u *UserService) DeleteUser(id uint) (err error) {
	var user database.User
	if err = global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	if err = global.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

// GetUserList 获取用户列表
func (u *UserService) GetUserList(listReq request.UserListRequest) (err error, list []database.User, total int64) {
	// 构建查询
	query := global.DB.Model(&database.User{})

	// 条件过滤
	if listReq.Username != "" {
		query = query.Where("username LIKE ?", "%"+listReq.Username+"%")
	}
	if listReq.Email != "" {
		query = query.Where("email LIKE ?", "%"+listReq.Email+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (listReq.Page - 1) * listReq.Size
	if err = query.Offset(offset).Limit(listReq.Size).Find(&list).Error; err != nil {
		return err, list, total
	}

	return nil, list, total
}

// ResetUserPassword 重置用户密码
func (u *UserService) ResetUserPassword(email, newPassword string) error {
	// 加密密码
	hashedPassword := utils.BcryptHash(newPassword)

	// 更新数据库中的密码
	return global.DB.Model(&database.User{}).Where("email = ?", email).Update("password", string(hashedPassword)).Error
}

// FindUserByEmail 根据邮箱查找用户
func (u *UserService) FindUserByEmail(email string) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// UploadAvatar 上传头像
func (u *UserService) UploadAvatar(file *multipart.FileHeader, userID uint) (string, error) {
	// 创建头像存储目录
	avatarDir := "uploads/avatars"
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return "", err
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	fileName := utils.GenerateUUID() + ext
	filePath := filepath.Join(avatarDir, fileName)

	// 保存文件
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	// 更新用户头像
	avatarURL := "/uploads/avatars/" + fileName
	if err := global.DB.Model(&database.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error; err != nil {
		// 如果数据库更新失败，删除已上传的文件
		os.Remove(filePath)
		return "", err
	}

	return avatarURL, nil
}

// ApproveUser 启用用户
func (u *UserService) ApproveUser(userUUID string) error {
	// 根据UUID查找用户
	var user database.User
	if err := global.DB.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 启用用户（设置状态为1）
	if err := global.DB.Model(&user).Update("status", 1).Error; err != nil {
		return err
	}

	return nil
}

// RejectUser 禁用用户
func (u *UserService) RejectUser(userUUID string) error {
	// 根据UUID查找用户
	var user database.User
	if err := global.DB.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 禁用用户（设置状态为0）
	if err := global.DB.Model(&user).Update("status", 0).Error; err != nil {
		return err
	}

	return nil
}
