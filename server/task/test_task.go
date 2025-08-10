package task

import (
	"server/global"
	"time"

	"github.com/robfig/cron/v3"
)

// TestTask 测试任务的具体逻辑
func TestTask() {
	global.ZapLog.Info("Test cron task executed at " + time.Now().Format("2006-01-02 15:04:05"))
}

// RegisterTestTask 注册测试任务
// 返回 error 便于上层处理注册失败情况
func RegisterTestTask(c *cron.Cron) error {
	// 每10秒执行一次
	_, err := c.AddFunc("*/10 * * * * *", TestTask)
	if err != nil {
		return err
	}
	global.ZapLog.Info("Test task registered successfully")
	return nil
}