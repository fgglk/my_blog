package initialize

import (
	"context"
	"os"
	"server/global"
	"server/service"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

func InitEs() {
	esConfig := global.Config.ES
	cfg := elasticsearch.Config{
		Addresses: []string{esConfig.Url},
		Username:  esConfig.Username,
		Password:  esConfig.Password,
	}

	// 1. 将NewClient改为NewTypedClient
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		global.ZapLog.Error("创建ES类型化客户端失败", zap.Error(err))
		os.Exit(1)
	}

	// 2. 调整连接验证逻辑（TypedClient的Ping方法签名不同）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用TypedClient的Ping方法
	pingSuccess, err := client.Ping().Do(ctx)
	if err != nil {
		global.ZapLog.Error("ES连接测试失败", zap.Error(err))
		os.Exit(1)
	}

	// 检查Ping是否成功
	if !pingSuccess {
		global.ZapLog.Error("ES服务Ping失败，服务可能不可用")
		os.Exit(1)
	}

	global.ES = client
	global.ZapLog.Info("es connected successfully")

	// 检查并创建索引（如果不存在）
	esService := service.NewArticleESService()
	exists, err := esService.IndexExists()
	if err != nil {
		global.ZapLog.Error("检查ES索引存在性失败", zap.Error(err))
		// 不退出程序，因为ES可能不是必需的
	} else if !exists {
		// 只在索引不存在时创建
		if err := esService.CreateIndex(); err != nil {
			global.ZapLog.Error("创建ES索引失败", zap.Error(err))
		} else {
			global.ZapLog.Info("ES索引创建成功")
		}
	} else {
		global.ZapLog.Info("ES索引已存在，跳过创建")
	}
}
