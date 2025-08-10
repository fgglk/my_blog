package global

import (
	"server/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
	ZapLog *zap.Logger
	Redis  *redis.Client
	ES     *elasticsearch.TypedClient
	Cron   *cron.Cron
)
