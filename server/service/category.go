package service

import (
	"server/global"
	"server/model/database"
)

type CategoryService struct{}

// GetCategoryListWithCount 获取分类列表及其文章数量
func (s *CategoryService) GetCategoryListWithCount() ([]database.CategoryWithCount, error) {
	var categories []database.CategoryWithCount

	// 使用原生SQL查询分类及其文章数量
	// 只统计未删除且已发布的文章
	query := `
		SELECT 
			c.id,
			c.name,
			c.slug,
			c.parent_id,
			c.sort,
			c.created_at,
			c.updated_at,
			COUNT(a.id) as article_count
		FROM categories c
		LEFT JOIN articles a ON c.id = a.category_id AND a.status = 1 AND a.deleted_at IS NULL
		GROUP BY c.id, c.name, c.slug, c.parent_id, c.sort, c.created_at, c.updated_at
		ORDER BY c.sort DESC, c.id ASC
	`

	err := global.DB.Raw(query).Scan(&categories).Error
	return categories, err
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(id uint) (database.Category, error) {
	var category database.Category
	err := global.DB.Where("id = ?", id).First(&category).Error
	return category, err
}
