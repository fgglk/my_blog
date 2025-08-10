package main

import (
	"server/core"
	"server/flag"
	"server/initialize"
	"server/middleware"
)

func main() {
	// 使用Air热重载工具进行开发
	core.InitConfig()
	initialize.InitZap()
	initialize.InitMysql()
	initialize.InitRedis()
	initialize.InitEs()
	initialize.InitCron()
	initialize.InitSeedData() // 初始化种子数据
	middleware.Init()
	initialize.InitRouter()
	flag.Run()

	core.StartServer()

}
