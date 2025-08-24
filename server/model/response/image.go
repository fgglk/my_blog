package response

import "server/model/database"

// ImageInfo 图片信息响应
type ImageInfo struct {
	ID           uint   `json:"id"`
	UUID         string `json:"uuid"`
	UserID       uint   `json:"user_id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	URL          string `json:"url"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	Width        *int   `json:"width"`
	Height       *int   `json:"height"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ImageListResponse 图片列表响应
type ImageListResponse struct {
	List     []ImageInfo `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// UploadImageResponse 上传图片响应
type UploadImageResponse struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// 将数据库模型转换为响应模型
func ToImageInfo(media database.Media) ImageInfo {
	return ImageInfo{
		ID:           media.ID,
		UUID:         "", // Media模型没有UUID字段
		UserID:       media.UserID,
		Filename:     media.Filename,
		OriginalName: media.Filename,          // 使用Filename作为原始名称
		URL:          "/" + media.StoragePath, // 构建静态文件URL
		Size:         media.FileSize,
		MimeType:     media.FileType,
		Width:        &media.Width,
		Height:       &media.Height,
		CreatedAt:    media.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    media.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// 将数据库模型切片转换为响应模型切片
func ToImageInfoList(medias []database.Media) []ImageInfo {
	var result []ImageInfo
	for _, media := range medias {
		result = append(result, ToImageInfo(media))
	}
	return result
}
