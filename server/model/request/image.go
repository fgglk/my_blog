package request

// UpdateImageInfo 更新图片信息请求
type UpdateImageInfo struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name"`
}
