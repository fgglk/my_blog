package service

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"

	"go.uber.org/zap"
)

// CommentService 评论服务结构体
type CommentService struct{}

// CreateComment 创建评论
func (s *CommentService) CreateComment(req request.CommentCreateRequest) (database.Comment, error) {
	comment := database.Comment{
		ArticleID: req.ArticleID,
		UserID:    req.UserID,
		Content:   req.Content,
		ParentID:  req.ParentID,
	}

	if err := global.DB.Create(&comment).Error; err != nil {
		global.ZapLog.Error("创建评论失败", zap.Error(err))
		return comment, errors.New("创建评论失败")
	}

	// 更新文章评论数
	s.updateArticleCommentCount(req.ArticleID)

	return comment, nil
}

// GetCommentList 获取评论列表（层级结构）
func (s *CommentService) GetCommentList(req request.CommentQueryRequest) ([]database.Comment, int64, error) {
	var comments []database.Comment
	var total int64
	
	// 查询所有评论（包括子评论）的总数
	if err := global.DB.Model(&database.Comment{}).Where("article_id = ?", req.ArticleID).Count(&total).Error; err != nil {
		global.ZapLog.Error("获取评论总数失败", zap.Error(err))
		return comments, 0, errors.New("获取评论总数失败")
	}

	// 只查询顶级评论（parent_id为null的评论）用于分页显示
	db := global.DB.Model(&database.Comment{}).Where("article_id = ? AND parent_id IS NULL", req.ArticleID)

	// 分页查询顶级评论
	pageSize := req.Size
	pageOffset := (req.Page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(pageOffset).Limit(pageSize).Find(&comments).Error; err != nil {
		global.ZapLog.Error("获取评论列表失败", zap.Error(err))
		return comments, 0, errors.New("获取评论列表失败")
	}

	// 为每个顶级评论加载子评论
	for i := range comments {
		childComments, err := s.GetChildComments(comments[i].ID)
		if err != nil {
			global.ZapLog.Error("获取子评论失败", zap.Error(err))
			continue
		}
		// 将子评论按时间排序
		for _, child := range childComments {
			comments[i].Children = append(comments[i].Children, child)
		}
	}

	return comments, total, nil
}

// GetCommentByID 根据ID获取评论
func (s *CommentService) GetCommentByID(id uint) (database.Comment, error) {
	var comment database.Comment
	if err := global.DB.Where("id = ?", id).First(&comment).Error; err != nil {
		return comment, err
	}
	return comment, nil
}

// GetChildComments 获取某个评论的所有子评论（递归获取多层嵌套）
func (s *CommentService) GetChildComments(parentID uint) ([]database.Comment, error) {
	var comments []database.Comment
	if err := global.DB.Where("parent_id = ?", parentID).Find(&comments).Error; err != nil {
		return nil, err
	}
	
	// 递归获取每个子评论的子评论
	for i := range comments {
		childComments, err := s.GetChildComments(comments[i].ID)
		if err != nil {
			global.ZapLog.Error("获取子评论失败", zap.Error(err))
			continue
		}
		comments[i].Children = childComments
	}
	
	return comments, nil
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(id, userID uint, req request.CommentUpdateRequest) (database.Comment, error) {
	var comment database.Comment
	if err := global.DB.Where("id = ?", id).First(&comment).Error; err != nil {
		global.ZapLog.Error("获取评论失败", zap.Error(err))
		return comment, errors.New("评论不存在")
	}

	// 检查权限（只能修改自己的评论）
	if comment.UserID != userID {
		return comment, errors.New("没有权限修改此评论")
	}

	// 更新评论内容
	comment.Content = req.Content
	if err := global.DB.Save(&comment).Error; err != nil {
		global.ZapLog.Error("更新评论失败", zap.Error(err))
		return comment, errors.New("更新评论失败")
	}

	return comment, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(id, userID uint) error {
	var comment database.Comment
	if err := global.DB.Where("id = ?", id).First(&comment).Error; err != nil {
		global.ZapLog.Error("获取评论失败", zap.Error(err))
		return errors.New("评论不存在")
	}

	// 获取当前用户信息
	var currentUser database.User
	if err := global.DB.Where("id = ?", userID).First(&currentUser).Error; err != nil {
		global.ZapLog.Error("获取用户信息失败", zap.Error(err))
		return errors.New("用户信息获取失败")
	}

	// 获取文章信息
	var article database.Article
	if err := global.DB.Where("id = ?", comment.ArticleID).First(&article).Error; err != nil {
		global.ZapLog.Error("获取文章信息失败", zap.Error(err))
		return errors.New("文章信息获取失败")
	}

	// 检查权限：管理员、文章作者或评论作者可以删除评论
	if currentUser.Role != "admin" && article.AuthorID != userID && comment.UserID != userID {
		return errors.New("没有权限删除此评论")
	}

	// 递归获取所有要删除的评论ID（包括所有层级的子评论）
	commentIDsToDelete := s.getAllChildCommentIDs(id)
	commentIDsToDelete = append(commentIDsToDelete, id)

	// 删除所有相关评论
	if err := global.DB.Where("id IN ?", commentIDsToDelete).Delete(&database.Comment{}).Error; err != nil {
		global.ZapLog.Error("删除评论失败", zap.Error(err))
		return errors.New("删除评论失败")
	}

	// 更新文章评论数
	s.updateArticleCommentCount(comment.ArticleID)

	return nil
}

// getAllChildCommentIDs 递归获取所有子评论的ID
func (s *CommentService) getAllChildCommentIDs(parentID uint) []uint {
	var childIDs []uint
	
	// 获取直接子评论
	var children []database.Comment
	if err := global.DB.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		global.ZapLog.Error("获取子评论失败", zap.Error(err))
		return childIDs
	}
	
	// 递归获取每个子评论的子评论
	for _, child := range children {
		childIDs = append(childIDs, child.ID)
		// 递归获取更深层的子评论
		deeperChildren := s.getAllChildCommentIDs(child.ID)
		childIDs = append(childIDs, deeperChildren...)
	}
	
	return childIDs
}

// ReplyToComment 回复评论
func (s *CommentService) ReplyToComment(parentID, userID uint, content string) (database.Comment, error) {
	// 获取父评论
	var parentComment database.Comment
	if err := global.DB.Where("id = ?", parentID).First(&parentComment).Error; err != nil {
		global.ZapLog.Error("获取父评论失败", zap.Error(err))
		return database.Comment{}, errors.New("父评论不存在")
	}

	// 创建回复评论
	comment := database.Comment{
		ArticleID: parentComment.ArticleID,
		UserID:    userID,
		Content:   content,
		ParentID:  &parentID,
	}

	if err := global.DB.Create(&comment).Error; err != nil {
		global.ZapLog.Error("创建回复评论失败", zap.Error(err))
		return comment, errors.New("创建回复评论失败")
	}

	// 更新文章评论数
	s.updateArticleCommentCount(parentComment.ArticleID)

	return comment, nil
}

// updateArticleCommentCount 更新文章评论数
func (s *CommentService) updateArticleCommentCount(articleID uint) {
	// 重新计算文章的实际评论数
	var commentCount int64
	if err := global.DB.Model(&database.Comment{}).Where("article_id = ?", articleID).Count(&commentCount).Error; err != nil {
		global.ZapLog.Error("计算文章评论数失败", zap.Error(err))
		return
	}

	// 更新文章的评论数
	if err := global.DB.Model(&database.Article{}).Where("id = ?", articleID).Update("comment_count", commentCount).Error; err != nil {
		global.ZapLog.Error("更新文章评论数失败", zap.Error(err))
	}
}
