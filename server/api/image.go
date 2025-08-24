package api

import (
	"fmt"
	"path/filepath"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	// 获取当前用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	file, header, err := c.Request.FormFile("avatar")
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
	_, _, media, err := service.UploadAvatar(file, header, userID)
	if err != nil {
		response.FailWithMessage("头像上传失败: "+err.Error(), c)
		return
	}

	// 生成直接访问的头像URL
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filepath.Base(media.StoragePath))
	fmt.Printf("生成的头像URL: %s\n", avatarURL)
	fmt.Printf("媒体文件路径: %s\n", media.StoragePath)

	// 更新用户头像
	if err := global.DB.Model(&database.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error; err != nil {
		response.FailWithMessage("更新用户头像失败: "+err.Error(), c)
		return
	}

	fmt.Printf("用户头像更新成功，用户ID: %d, 头像URL: %s\n", userID, avatarURL)
	response.OkWithDetailed(gin.H{
		"url": avatarURL,
	}, "头像上传成功", c)
}

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	// 获取当前用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		response.FailWithMessage("获取图片文件失败: "+err.Error(), c)
		return
	}

	// 检查文件类型
	if !utils.IsImageType(header.Header.Get("Content-Type")) {
		response.FailWithMessage("不支持的图片格式", c)
		return
	}

	// 调用Service层上传图片
	imageURL, imageID, media, err := service.UploadImage(file, header, userID)
	if err != nil {
		response.FailWithMessage("图片上传失败: "+err.Error(), c)
		return
	}

	// 在使用 imageID 前添加转换逻辑
	imageIDUint, err := utils.StringToUint(imageID)
	if err != nil {
		response.FailWithMessage("图片ID格式错误: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.ImageResponse{
		ID:         imageIDUint,
		URL:        imageURL,                                      // 使用正确的API URL
		UploadTime: media.CreatedAt.Format("2006-01-02 15:04:05"), // 添加上传时间
		UpdateTime: media.UpdatedAt.Format("2006-01-02 15:04:05"), // 添加更新时间
	}, "图片上传成功", c)
}

// ShowImage 查看图片
func ShowImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	media, err := service.GetImageByID(id)
	if err != nil {
		response.FailWithMessage("图片不存在或已过期", c)
		return
	}

	// 设置响应头并返回图片
	c.Header("Content-Type", media.FileType)
	c.File(media.StoragePath)
}

// GetImageList 获取图片列表
func GetImageList(c *gin.Context) {
	// 获取当前用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	var req request.ImageListRequest
	_ = c.ShouldBindQuery(&req)

	// 设置分页默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 10
	}

	list, total, err := service.GetImageList(req, userID)
	if err != nil {
		response.FailWithMessage("获取图片列表失败", c)
		return
	}

	// 转换为前端需要的格式
	imageList := response.ToImageInfoList(list)

	response.OkWithDetailed(response.ImageListResponse{
		List:     imageList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.Size,
	}, "获取成功", c)
}

// DeleteImage 删除图片
func DeleteImage(c *gin.Context) {
	// 获取当前用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	idStr := c.Param("id")
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := service.DeleteImage(id, userID); err != nil {
		response.FailWithMessage("图片删除失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("图片删除成功", c)
}

// UpdateImage 更新图片信息
func UpdateImage(c *gin.Context) {
	// 获取当前用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	var imageInfo request.UpdateImageInfo
	_ = c.ShouldBindJSON(&imageInfo)
	idStr := c.Param("id")
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	imageInfo.ID = id

	if err := service.UpdateImage(imageInfo, userID); err != nil {
		response.FailWithMessage("图片信息更新失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("图片信息更新成功", c)
}
