package flag

import (
	"server/global"
	"server/service"

	"go.uber.org/zap"
)

func createEsIndex() error {

	esService := service.NewArticleESService()

	if err := esService.CreateIndex(); err != nil {
		global.ZapLog.Error("创建Elasticsearch索引失败", zap.Error(err))
		return err
	}

	global.ZapLog.Info("Elasticsearch索引创建成功")
	return nil
}
