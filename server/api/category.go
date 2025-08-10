package api

import (
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type CategoryApi struct{}

var categoryService = service.ServiceGroups.CategoryService

// @Summary 获取分类列表
// @Description 获取所有分类及其文章数量
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.CategoryWithCountResponse}
// @Router /api/categories [get]
func (c *CategoryApi) GetCategoryList(ctx *gin.Context) {
	categories, err := categoryService.GetCategoryListWithCount()
	if err != nil {
		response.FailWithMessage("获取分类列表失败: "+err.Error(), ctx)
		return
	}

	var categoryResponses []response.CategoryWithCountResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, response.ToCategoryWithCountResponse(category))
	}

	response.OkWithData(categoryResponses, ctx)
}

// @Summary 获取分类详情
// @Description 根据ID获取分类详情
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=response.CategoryResponse}
// @Router /api/categories/{id} [get]
func (c *CategoryApi) GetCategory(ctx *gin.Context) {
	id, err := utils.StringToUint(ctx.Param("id"))
	if err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	category, err := categoryService.GetCategoryByID(id)
	if err != nil {
		response.FailWithMessage("分类不存在", ctx)
		return
	}

	response.OkWithData(response.ToCategoryResponse(category), ctx)
} 