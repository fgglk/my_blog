package api

import (
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type PageApi struct{}

var pageService = service.ServiceGroups.PageService

// CreatePage 创建页面
func (p *PageApi) CreatePage(c *gin.Context) {
	var req request.CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 检查Slug是否已存在
	if _, err := pageService.GetPageBySlug(req.Slug); err == nil {
		response.FailWithMessage("URL路径已存在: "+req.Slug, c)
		return
	}

	page := database.Page{
		Title:     req.Title,
		Slug:      req.Slug,
		Content:   req.Content,
		Template:  req.Template,
		ShowInNav: req.ShowInNav,
		Sort:      req.Sort,
	}
	page.Status = req.Status
	if err := pageService.CreatePage(page); err != nil {
		response.FailWithMessage("创建页面失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("创建页面成功", c)
}

// GetPage 获取单个页面
func (p *PageApi) GetPage(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	page, err := pageService.GetPageByID(id)
	if err != nil {
		response.FailWithMessage("页面不存在", c)
		return
	}

	resp := response.ToPageResponse(page)
	response.OkWithData(resp, c)
}

// GetPageBySlug 通过Slug获取页面（前台用）
func (p *PageApi) GetPageBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.FailWithMessage("参数错误", c)
		return
	}

	page, err := pageService.GetPageBySlug(slug)
	if err != nil {
		response.FailWithMessage("页面不存在或已禁用", c)
		return
	}

	resp := response.ToPageResponse(page)
	response.OkWithData(resp, c)
}

// UpdatePage 更新页面
func (p *PageApi) UpdatePage(c *gin.Context) {
	var req request.UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	page, err := pageService.GetPageByID(req.ID)
	if err != nil {
		response.FailWithMessage("页面不存在: "+err.Error(), c)
		return
	}

	// 检查Slug是否已被其他页面使用
	if req.Slug != page.Slug {
		if _, err := pageService.GetPageBySlug(req.Slug); err == nil {
			response.FailWithMessage("URL路径已存在", c)
			return
		}
	}

	// 更新字段
	page.Title = req.Title
	page.Slug = req.Slug
	page.Content = req.Content
	page.Template = req.Template
	page.ShowInNav = req.ShowInNav
	page.Sort = req.Sort
	page.Status = req.Status

	if err := pageService.UpdatePage(page); err != nil {
		response.FailWithMessage("更新页面失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("更新页面成功", c)
}

// DeletePage 删除页面
func (p *PageApi) DeletePage(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 先检查页面是否存在
	_, err = pageService.GetPageByID(id)
	if err != nil {
		response.FailWithMessage("页面不存在或已被删除", c)
		return
	}

	if err := pageService.DeletePage(id); err != nil {
		response.FailWithMessage("删除页面失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除页面成功", c)
}

// ListPages 分页查询页面列表
func (p *PageApi) ListPages(c *gin.Context) {
	var req request.PageQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 10
	}

	pages, total, err := pageService.ListPages(req)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}

	resp := response.ToPageListResponse(pages, total, req.Page, req.Size)
	response.OkWithData(resp, c)
}

// GetNavPages 获取导航栏页面列表
func (p *PageApi) GetNavPages(c *gin.Context) {
	pages, err := pageService.GetNavPages()
	if err != nil {
		response.FailWithMessage("获取导航页面失败: "+err.Error(), c)
		return
	}
	resp := response.ToPageListResponse(pages, int64(len(pages)), 1, len(pages))
	response.OkWithData(resp, c)
}
