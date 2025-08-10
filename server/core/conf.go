package core

import (
	"log"
	"server/config"
	"server/global"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败:%v", err)
	}
	global.Config = &config.Config{}
	if err := viper.Unmarshal(global.Config); err != nil {
		log.Fatalf("配置文件解析失败:%v", err)
	}
}
