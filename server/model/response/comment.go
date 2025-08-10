package response

import (
	"server/global"
	"server/model/database"
	"time"
)

// CommentResponse 单个评论响应结构体
type CommentResponse struct {
	ID              uint              `json:"id"`
	UserID          uint              `json:"user_id"`
	UserName        string            `json:"user_name"`
	UserAvatar      string            `json:"user_avatar"` // 添加用户头像字段
	ArticleID       uint              `json:"article_id"`
	Content         string            `json:"content"`
	ParentID        *uint             `json:"parent_id,omitempty"`
	ParentUserName  string            `json:"parent_user_name,omitempty"` // 父评论用户名
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Children        []CommentResponse `json:"children,omitempty"` // 新增子评论字段
}

// CommentListResponse 评论列表响应结构体
type CommentListResponse struct {
	List  []CommentResponse `json:"list"`  // 评论列表
	Total int64             `json:"total"` // 总条数
	Page  int               `json:"page"`  // 当前页码
	Size  int               `json:"size"`  // 每页条数
}

// ToCommentResponse 将数据库评论模型转换为响应结构体
func ToCommentResponse(comment database.Comment) CommentResponse {
	// 根据UserID查询用户名
	var user database.User
	global.DB.Where("id = ?", comment.UserID).First(&user)

	// 如果有父评论，查询父评论的用户名
	var parentUserName string
	if comment.ParentID != nil {
		var parentComment database.Comment
		if global.DB.Where("id = ?", *comment.ParentID).First(&parentComment).Error == nil {
			var parentUser database.User
			if global.DB.Where("id = ?", parentComment.UserID).First(&parentUser).Error == nil {
				parentUserName = parentUser.Username
			}
		}
	}

	// 处理子评论
	var children []CommentResponse
	for _, child := range comment.Children {
		children = append(children, ToCommentResponse(child))
	}

	return CommentResponse{
		ID:             comment.ID,
		ArticleID:      comment.ArticleID,
		UserID:         comment.UserID,
		UserName:       user.Username, // 赋值用户名
		UserAvatar:     user.Avatar,   // 赋值用户头像
		Content:        comment.Content,
		ParentID:       comment.ParentID,
		ParentUserName: parentUserName, // 赋值父评论用户名
		CreatedAt:      comment.CreatedAt,
		UpdatedAt:      comment.UpdatedAt,
		Children:       children, // 添加子评论
	}
}

// ToCommentListResponse 将数据库评论列表转换为响应结构体
func ToCommentListResponse(comments []database.Comment, total int64, page, size int) CommentListResponse {
	var commentResponses []CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, ToCommentResponse(comment))
	}

	return CommentListResponse{
		List:  commentResponses,
		Total: total,
		Page:  page,
		Size:  size,
	}
}
