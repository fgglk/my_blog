package task

import (
	"server/global"
	"server/model/database"
	"server/service"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// SyncArticleStatsTask 同步文章统计数据到ES
func SyncArticleStatsTask() {
	startTime := time.Now()
	global.ZapLog.Info("开始同步文章统计数据到ES", zap.String("start_time", startTime.Format("2006-01-02 15:04:05")))

	articleService := service.ArticleService{}

	// 获取最近24小时内有更新的文章ID
	var articleIDs []uint
	twelveHoursAgo := time.Now().Add(-24 * time.Hour)

	err := global.DB.Model(&database.Article{}).
		Where("updated_at >= ?", twelveHoursAgo).
		Pluck("id", &articleIDs).Error

	if err != nil {
		global.ZapLog.Error("获取需要同步的文章ID失败", zap.Error(err))
		return
	}

	global.ZapLog.Info("获取到需要同步的文章数量", zap.Int("count", len(articleIDs)))

	// 批量同步文章统计数据
	successCount := 0
	failCount := 0

	for _, id := range articleIDs {
		if err := articleService.SyncArticleStatsToES(id); err != nil {
			global.ZapLog.Error("同步文章统计数据失败", zap.Uint("article_id", id), zap.Error(err))
			failCount++
		} else {
			successCount++
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	global.ZapLog.Info("文章统计数据同步完成",
		zap.String("end_time", endTime.Format("2006-01-02 15:04:05")),
		zap.Duration("duration", duration),
		zap.Int("success_count", successCount),
		zap.Int("fail_count", failCount),
	)
}

// RegisterSyncArticleStatsTask 注册文章统计数据同步任务
func RegisterSyncArticleStatsTask(c *cron.Cron) error {
	// 每5分钟执行一次，可根据需求调整
	_, err := c.AddFunc("0 */5 * * * *", SyncArticleStatsTask)
	if err != nil {
		return err
	}
	global.ZapLog.Info("文章统计数据同步任务注册成功")
	return nil
}
