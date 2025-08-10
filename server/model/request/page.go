package request

// CreatePageRequest 创建页面请求
type CreatePageRequest struct {
	Title     string `json:"title" binding:"required,min=1,max=200"`
	Slug      string `json:"slug" binding:"required,min=1,max=255"`
	Content   string `json:"content" binding:"required"`
	Template  string `json:"template" binding:"max=100"`
	ShowInNav bool   `json:"show_in_nav"`
	Sort      int    `json:"sort"`
	Status    uint8  `json:"status" binding:"oneof=0 1"`
}

// UpdatePageRequest 更新页面请求
type UpdatePageRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Title     string `json:"title" binding:"required,min=1,max=200"`
	Slug      string `json:"slug" binding:"required,min=1,max=255"`
	Content   string `json:"content" binding:"required"`
	Template  string `json:"template" binding:"max=100"`
	ShowInNav bool   `json:"show_in_nav"`
	Sort      int    `json:"sort"`
	Status    uint8  `json:"status" binding:"oneof=0 1"`
}

// PageQueryRequest 页面查询请求
type PageQueryRequest struct {
	Title     string `form:"title"`
	Slug      string `form:"slug"`
	ShowInNav *bool  `form:"show_in_nav"`
	Status    *uint8 `form:"status"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Size      int    `form:"size" binding:"omitempty,min=1,max=100"`
}
