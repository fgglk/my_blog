package initialize

import (
	"server/global"
	"server/model/database"

	"go.uber.org/zap"
)

// InitSeedData 初始化种子数据
func InitSeedData() {
	// 初始化默认分类
	initDefaultCategories()
}

// initDefaultCategories 初始化默认分类
func initDefaultCategories() {
	var count int64
	global.DB.Model(&database.Category{}).Count(&count)

	// 如果已经有分类数据，跳过初始化
	if count > 0 {
		return
	}

	// 创建固定的分类
	defaultCategories := []database.Category{
		{
			Name: "技术杂谈",
			Slug: "tech-talk",
			Sort: 1,
		},
		{
			Name: "生活随笔",
			Slug: "life-notes",
			Sort: 2,
		},
		{
			Name: "读书笔记",
			Slug: "reading-notes",
			Sort: 3,
		},
	}

	for _, category := range defaultCategories {
		if err := global.DB.Create(&category).Error; err != nil {
			global.ZapLog.Error("创建固定分类失败", zap.Error(err))
		}
	}

	global.ZapLog.Info("固定分类初始化完成")
}
