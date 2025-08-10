package response

import (
	"server/model/database"
)

// PageResponse 页面响应结构体
type PageResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Content   string `json:"content"`
	Template  string `json:"template"`
	ShowInNav bool   `json:"show_in_nav"`
	Sort      int    `json:"sort"`
	Status    uint8  `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ToPageResponse 转换为页面响应
type PageListResponse struct {
	List  []PageResponse `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// ToPageResponse 转换单个页面
func ToPageResponse(page database.Page) PageResponse {
	return PageResponse{
		ID:        page.ID,
		Title:     page.Title,
		Slug:      page.Slug,
		Content:   page.Content,
		Template:  page.Template,
		ShowInNav: page.ShowInNav,
		Sort:      page.Sort,
		Status:    page.Status,
		CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToPageListResponse 转换页面列表
func ToPageListResponse(pages []database.Page, total int64, page, size int) PageListResponse {

	list := make([]PageResponse, 0, len(pages))
	for _, page := range pages {
		list = append(list, ToPageResponse(page))
	}
	return PageListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}
}
