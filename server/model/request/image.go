package request

// ImageListRequest 图片列表查询请求
type ImageListRequest struct {
	Page      int    `form:"page" validate:"min=1"`
	Size      int    `form:"size" validate:"min=1,max=100"`
	Keyword   string `form:"keyword" validate:"omitempty"`
	SortBy    string `form:"sortBy" validate:"omitempty,oneof=createdAt size filename"`
	SortOrder string `form:"sortOrder" validate:"omitempty,oneof=asc desc"`
}

// DeleteImageRequest 删除图片请求
type DeleteImageRequest struct {
	ImageID uint `json:"image_id" validate:"required"`
}

// BatchDeleteImageRequest 批量删除图片请求
type BatchDeleteImageRequest struct {
	ImageIDs []uint `json:"image_ids" validate:"required,min=1"`
}

// UpdateImageInfo 更新图片信息请求
type UpdateImageInfo struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
