package service

import (
	"server/global"
	"server/model/database"
	"server/model/request"
)

type PageService struct{}

// CreatePage 创建页面
func (s *PageService) CreatePage(page database.Page) error {
	return global.DB.Create(&page).Error
}

// GetPageByID 根据ID获取页面
func (s *PageService) GetPageByID(id uint) (database.Page, error) {
	var page database.Page
	err := global.DB.Where("id = ?", id).First(&page).Error
	return page, err
}

// GetPageBySlug 根据Slug获取页面
func (s *PageService) GetPageBySlug(slug string) (database.Page, error) {
	var page database.Page
	err := global.DB.Where("slug = ? AND status = 1", slug).First(&page).Error
	return page, err
}

// UpdatePage 更新页面
func (s *PageService) UpdatePage(page database.Page) error {
	return global.DB.Save(&page).Error
}

// DeletePage 删除页面
func (s *PageService) DeletePage(id uint) error {
	return global.DB.Delete(&database.Page{}, id).Error
}

// ListPages 分页查询页面
func (s *PageService) ListPages(req request.PageQueryRequest) ([]database.Page, int64, error) {
	db := global.DB.Model(&database.Page{})

	// 条件过滤
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Slug != "" {
		db = db.Where("slug LIKE ?", "%"+req.Slug+"%")
	}
	if req.ShowInNav != nil {
		db = db.Where("show_in_nav = ?", req.ShowInNav)
	}
	if req.Status != nil {
		db = db.Where("status = ?", req.Status)
	}

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var pages []database.Page
	offset := (req.Page - 1) * req.Size
	err := db.Order("sort DESC, id DESC").Offset(offset).Limit(req.Size).Find(&pages).Error

	return pages, total, err
}

// GetNavPages 获取导航栏页面
func (s *PageService) GetNavPages() ([]database.Page, error) {
	var pages []database.Page
	err := global.DB.Where("show_in_nav = true AND status = 1").Order("sort DESC, id DESC").Find(&pages).Error
	return pages, err
}
