package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/utils"

	"mime/multipart"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UploadAvatar 上传头像
func UploadAvatar(file io.Reader, header *multipart.FileHeader, userID uint) (string, string, database.Media, error) {
	// 生成唯一文件名，使用用户ID和时间戳
	fileName := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), filepath.Ext(header.Filename))
	
	// 头像存储在 avatars 目录
	uploadPath := "uploads/avatars"

	// 创建目录
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", "", database.Media{}, fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存文件到本地
	filePath := filepath.Join(uploadPath, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return "", "", database.Media{}, fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", "", database.Media{}, fmt.Errorf("保存文件失败: %v", err)
	}

	// 保存到数据库
	media := database.Media{
		Filename:    header.Filename,
		StoragePath: filePath,
		FileSize:    header.Size,
		FileType:    header.Header.Get("Content-Type"),
		UserID:      userID,
	}

	// 如果是图片，获取宽高信息
	if utils.IsImageType(media.FileType) {
		width, height, err := utils.GetImageDimensions(filePath)
		if err == nil {
			media.Width = width
			media.Height = height
		}
	}

	if err := global.DB.Create(&media).Error; err != nil {
		// 保存失败，删除文件
		os.Remove(filePath)
		return "", "", database.Media{}, fmt.Errorf("保存媒体信息失败: %v", err)
	}

	// 生成访问URL
	imageURL := fmt.Sprintf("/api/image/show/%d", media.ID)

	return imageURL, fmt.Sprintf("%d", media.ID), media, nil
}

// UploadImage 上传图片到本地存储
// 修改UploadImage函数返回值，包含完整的media对象
func UploadImage(file io.Reader, header *multipart.FileHeader, userID uint) (string, string, database.Media, error) {
	// 生成唯一文件名
	fileName := fmt.Sprintf("image_%d_%d%s", userID, time.Now().Unix(), filepath.Ext(header.Filename))
	
	// 文章图片存储在 images 目录
	uploadPath := "uploads/images"

	// 创建目录
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		// 返回空的media对象作为第三个返回值
		return "", "", database.Media{}, fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存文件到本地
	filePath := filepath.Join(uploadPath, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return "", "", database.Media{}, fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", "", database.Media{}, fmt.Errorf("保存文件失败: %v", err)
	}

	// 保存到数据库
	media := database.Media{
		Filename:    header.Filename,
		StoragePath: filePath,
		FileSize:    header.Size,
		FileType:    header.Header.Get("Content-Type"),
		UserID:      userID,
	}

	// 如果是图片，获取宽高信息
	if utils.IsImageType(media.FileType) {
		width, height, err := utils.GetImageDimensions(filePath)
		if err == nil {
			media.Width = width
			media.Height = height
		}
	}

	if err := global.DB.Create(&media).Error; err != nil {
		// 保存失败，删除文件
		os.Remove(filePath)
		return "", "", database.Media{}, fmt.Errorf("保存媒体信息失败: %v", err)
	}

	// 生成访问URL
	imageURL := fmt.Sprintf("/api/image/show/%d", media.ID)

	return imageURL, fmt.Sprintf("%d", media.ID), media, nil
}

// GetImageByID 根据ID获取图片信息
func GetImageByID(id uint) (database.Media, error) {
	var media database.Media
	if err := global.DB.Where("id = ?", id).First(&media).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return media, errors.New("图片不存在")
		}
		return media, err
	}
	return media, nil
}

// GetImageList 获取图片列表
func GetImageList(pageInfo request.PageInfo, userID uint) ([]database.Media, int64, error) {
	var list []database.Media
	var total int64
	db := global.DB.Model(&database.Media{}).Where("user_id = ?", userID)

	// 只查询图片类型
	db = db.Where("file_type LIKE ?", "image/%")

	// 分页
	offset := (pageInfo.Page - 1) * pageInfo.Size
	db.Count(&total)
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageInfo.Size).Find(&list).Error; err != nil {
		return list, total, err
	}

	return list, total, nil
}

// DeleteImage 删除图片
func DeleteImage(id uint, userID uint) error {
	// 获取图片信息
	media, err := GetImageByID(id)
	if err != nil {
		return err
	}

	// 检查权限
	if media.UserID != userID {
		return errors.New("没有权限删除此图片")
	}

	// 删除数据库记录
	if err := global.DB.Delete(&database.Media{}, "id = ?", id).Error; err != nil {
		return err
	}

	// 删除文件
	if err := os.Remove(media.StoragePath); err != nil {
		// 记录日志但不阻止删除操作
		global.ZapLog.Error("删除图片文件失败", zap.Error(err))
	}

	return nil
}

// UpdateImage 更新图片信息
func UpdateImage(imageInfo request.UpdateImageInfo, userID uint) error {
	// 检查权限
	media, err := GetImageByID(imageInfo.ID)
	if err != nil {
		return err
	}

	if media.UserID != userID {
		return errors.New("没有权限更新此图片")
	}

	return global.DB.Model(&database.Media{}).Where("id = ?", imageInfo.ID).Updates(map[string]interface{}{
		"filename": imageInfo.Name,
	}).Error
}
