package request

// IdRequest ID请求
type IdRequest struct {
	ID uint `json:"id" validate:"required"`
}

type PageInfo struct {
	Page  int    `json:"page" form:"page"`   // 页码
	Size  int    `json:"size" form:"size"`   // 每页大小
	Sort  string `json:"sort" form:"sort"`   // 排序字段
	Order string `json:"order" form:"order"` // 排序方式: asc/desc
}
