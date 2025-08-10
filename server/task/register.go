package task

import (
	"server/global"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func RegisterTask(c *cron.Cron) {
	if err := RegisterSyncArticleStatsTask(c); err != nil {
		global.ZapLog.Error("注册文章统计数据同步任务失败", zap.Error(err))
	}
}
