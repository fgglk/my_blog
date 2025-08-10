package flag

import (
	"server/global"
	"server/model/database"

	"go.uber.org/zap"
)

func migrateDatabase() error {

	err := global.DB.AutoMigrate(
		&database.User{},
		&database.Comment{},
		&database.Article{},
		&database.Like{},
		&database.Favorite{},
		&database.Category{},
		&database.Tag{},
		&database.ArticleTag{},
		&database.Media{},
		&database.Page{},
	)
	if err != nil {
		global.ZapLog.Error("数据库表结构迁移失败", zap.Error(err))
		return err
	}
	global.ZapLog.Info("数据库表结构迁移成功")
	return nil
}
