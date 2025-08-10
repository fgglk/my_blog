package response

import (
	"server/model/database"
	"time"
)

// ImageResponse 图片响应
type ImageResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Size       int64     `json:"size"`
	FileType   string    `json:"file_type"`
	Width      int       `json:"width,omitempty"`
	Height     int       `json:"height,omitempty"`
	UploadTime time.Time `json:"upload_time"`
	UserID     uint      `json:"user_id,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`
}

// ToImageResponse 转换Media为ImageResponse
func ToImageResponse(media database.Media) ImageResponse {
	return ImageResponse{
		ID:         media.ID,
		Name:       media.Filename,
		URL:        media.StoragePath,
		Size:       media.FileSize,
		FileType:   media.FileType,
		Width:      media.Width,
		Height:     media.Height,
		UploadTime: media.CreatedAt,
		UserID:     media.UserID,
		UpdateTime: media.UpdatedAt,
	}
}
