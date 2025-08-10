package request

// CommentCreateRequest 创建评论请求结构体
type CommentCreateRequest struct {
	ArticleID uint   `json:"article_id" binding:"required,min=1"`      // 文章ID
	UserID    uint   `json:"user_id" binding:""`                       // 用户ID（完全移除验证标签）
	Content   string `json:"content" binding:"required,min=1,max=500"` // 评论内容
	ParentID  *uint  `json:"parent_id" binding:"omitempty,min=0"`      // 父评论ID，用于回复
}

// CommentUpdateRequest 更新评论请求结构体
type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"` // 评论内容
}

// CommentQueryRequest 查询评论列表请求结构体
type CommentQueryRequest struct {
	ArticleID uint `form:"article_id" binding:"required,min=1"`   // 文章ID
	Page      int  `form:"page" binding:"omitempty,min=1"`       // 页码（移除了required标签，添加了omitempty）
	Size      int  `form:"size" binding:"omitempty,min=1,max=100"` // 每页条数（移除了required标签，添加了omitempty）
}

// CommentReplyRequest 回复评论请求结构体
type CommentReplyRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"` // 回复内容
}
