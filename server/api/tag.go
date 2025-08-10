package api

import (
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type TagApi struct{}

var tagService = service.ServiceGroups.TagService

// @Summary 获取标签列表
// @Description 获取所有标签及其文章数量
// @Tags tag
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.TagWithCountResponse}
// @Router /api/tags [get]
func (t *TagApi) GetTagList(ctx *gin.Context) {
	tags, err := tagService.GetTagListWithCount()
	if err != nil {
		response.FailWithMessage("获取标签列表失败: "+err.Error(), ctx)
		return
	}

	var tagResponses []response.TagWithCountResponse
	for _, tag := range tags {
		tagResponses = append(tagResponses, response.ToTagWithCountResponse(tag))
	}

	response.OkWithData(tagResponses, ctx)
}

// @Summary 获取标签详情
// @Description 根据ID获取标签详情
// @Tags tag
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response{data=response.TagResponse}
// @Router /api/tags/{id} [get]
func (t *TagApi) GetTag(ctx *gin.Context) {
	id, err := utils.StringToUint(ctx.Param("id"))
	if err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	tag, err := tagService.GetTagByID(id)
	if err != nil {
		response.FailWithMessage("标签不存在", ctx)
		return
	}

	response.OkWithData(response.ToTagResponse(tag), ctx)
}

// @Summary 清理孤儿标签
// @Description 清理不再被任何文章使用的标签，需要管理员权限
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=map[string]int64}
// @Router /api/tags/cleanup [delete]
func (t *TagApi) CleanupOrphanTags(ctx *gin.Context) {
	// 获取用户ID并检查管理员权限
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		response.NoAuth(err.Error(), ctx)
		return
	}

	if !utils.IsAdmin(userID) {
		response.Forbidden("需要管理员权限", ctx)
		return
	}

	// 执行清理
	deletedCount, err := tagService.CleanupOrphanTags()
	if err != nil {
		response.FailWithMessage("清理孤儿标签失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(map[string]int64{"deleted_count": deletedCount}, ctx)
} 