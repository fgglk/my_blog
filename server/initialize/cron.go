package initialize

import (
	"server/global"
	"server/task"

	"context"
	"server/hooks"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger)))

	task.RegisterTask(c)

	c.Start()

	global.Cron = c
	global.ZapLog.Info("cron started successfully")

	hooks.GetHookManager().RegisterHook(hooks.ShutdownHook, func(ctx context.Context) error {
		stopCtx := c.Stop()
		<-stopCtx.Done()
		global.ZapLog.Info("cron stopped successfully")
		return nil
	})
}
