package utils

import (
	"image"
	"mime/multipart"
	"os"
	"path/filepath"
	"io"

	// 导入图片格式支持
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// IsImageType 检查是否为图片类型
func IsImageType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp":
		return true
	default:
		return false
	}
}

// GetImageDimensions 获取图片宽高
func GetImageDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 确保目标目录存在
	dir := dst[:len(dst)-len(filepath.Base(dst))]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 复制文件内容
	_, err = io.Copy(out, src)
	return err
}