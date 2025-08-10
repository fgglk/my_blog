package api

import (
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CommentApi struct{}

var commentService = service.ServiceGroups.CommentService

// @Summary 创建评论
// @Description 创建新评论，需要认证
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.CommentCreateRequest true "评论创建信息"
// @Success 200 {object} response.Response{data=response.CommentResponse}
// @Router /api/comments [post]
func (a *CommentApi) CreateComment(c *gin.Context) {
	var req request.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			response.FailWithMessage("参数错误: "+utils.TranslateValidationError(validateErr), c)
			return
		}
		response.FailWithMessage("请求参数格式错误", c)
		return
	}

	// 使用 utils.GetUserID 获取用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	req.UserID = userID
	comment, err := commentService.CreateComment(req)
	if err != nil {
		response.FailWithMessage("创建评论失败: "+err.Error(), c)
		return
	}

	response.OkWithData(response.ToCommentResponse(comment), c)
}

// @Summary 获取评论列表
// @Description 分页获取评论列表
// @Tags comment
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Param page query int false "页码，默认为1"
// @Param size query int false "每页条数，默认为10，最大100"
// @Success 200 {object} response.Response{data=response.CommentListResponse}
// @Router /api/comments [get]
func (a *CommentApi) GetCommentList(c *gin.Context) {
	var req request.CommentQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: " + err.Error(), c)
		return
	}

	// 设置默认分页值
	if req.Page <= 0 {
		req.Page = 1 // 默认第一页
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 10 // 默认每页10条
	}

	comments, total, err := commentService.GetCommentList(req)
	if err != nil {
		response.FailWithMessage("获取评论列表失败: "+err.Error(), c)
		return
	}

	// 将 database.Comment 切片转换为 response.CommentResponse 切片
	var commentResponses []response.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, response.ToCommentResponse(comment))
	}

	response.OkWithData(response.CommentListResponse{
		List:  commentResponses,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, c)
}

// @Summary 获取单个评论
// @Description 根据ID获取单个评论
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response{data=response.CommentResponse}
// @Router /api/comments/{id} [get]
func (a *CommentApi) GetComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取主评论
	comment, err := commentService.GetCommentByID(id)
	if err != nil {
		response.FailWithMessage("获取评论失败: " + err.Error(), c)
		return
	}

	// 转换为主评论响应
	commentResp := response.ToCommentResponse(comment)

	// 获取子评论
	childComments, err := commentService.GetChildComments(id)
	if err != nil {
		response.FailWithMessage("获取子评论失败: " + err.Error(), c)
		return
	}

	// 转换子评论并添加到主评论响应中
	for _, child := range childComments {
		commentResp.Children = append(commentResp.Children, response.ToCommentResponse(child))
	}

	response.OkWithData(commentResp, c)
}

// @Summary 更新评论
// @Description 更新评论内容，需要认证
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param data body request.CommentUpdateRequest true "评论更新信息"
// @Success 200 {object} response.Response{data=response.CommentResponse}
// @Router /api/comments/{id} [put]
func (a *CommentApi) UpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	// 使用 utils.StringToUint 转换ID
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var req request.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			response.FailWithMessage("参数错误: "+utils.TranslateValidationError(validateErr), c)
			return
		}
		response.FailWithMessage("请求参数格式错误", c)
		return
	}

	// 使用 utils.GetUserID 获取用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	comment, err := commentService.UpdateComment(id, userID, req)
	if err != nil {
		response.FailWithMessage("更新评论失败: "+err.Error(), c)
		return
	}

	response.OkWithData(response.ToCommentResponse(comment), c)
}

// @Summary 删除评论
// @Description 删除评论，需要认证
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/comments/{id} [delete]
func (a *CommentApi) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	// 使用 utils.StringToUint 转换ID
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 使用 utils.GetUserID 获取用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	if err := commentService.DeleteComment(id, userID); err != nil {
		response.FailWithMessage("删除评论失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("评论删除成功", c)
}

// @Summary 回复评论
// @Description 回复评论，需要认证
// @Tags comment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param data body request.CommentReplyRequest true "回复信息"
// @Success 200 {object} response.Response{data=response.CommentResponse}
// @Router /api/comments/{id}/reply [post]
func (a *CommentApi) ReplyToComment(c *gin.Context) {
	idStr := c.Param("id")
	// 使用 utils.StringToUint 转换ID
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var req request.CommentReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			response.FailWithMessage("参数错误: "+utils.TranslateValidationError(validateErr), c)
			return
		}
		response.FailWithMessage("请求参数格式错误", c)
		return
	}

	// 使用 utils.GetUserID 获取用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}
	comment, err := commentService.ReplyToComment(id, userID, req.Content)
	if err != nil {
		response.FailWithMessage("回复评论失败: "+err.Error(), c)
		return
	}

	response.OkWithData(response.ToCommentResponse(comment), c)
}
