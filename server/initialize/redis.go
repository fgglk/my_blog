package initialize

import (
	"server/global"

	"context"
	"os"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func InitRedis() {
	redisConfig := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Address,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.PoolSize,
		MinIdleConns: redisConfig.MinIdleConns,
		IdleTimeout:  redisConfig.IdleTimeout,
	})
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping().Result()
	if err != nil {
		global.ZapLog.Error("failed to connect redis", zap.Error(err))
		os.Exit(1)
	}
	global.Redis = client
	global.ZapLog.Info("redis connected successfully")
}
