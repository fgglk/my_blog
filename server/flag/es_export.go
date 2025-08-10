package flag

import (
	"context"
	"encoding/json"
	"os"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/zap"

	"server/global"
)

func exportEsData(path string) error {
	indexName := "articles"
	ctx := context.Background()
	var scrollID string
	first := true

	// 创建输出文件
	file, err := os.Create(path)
	if err != nil {
		global.ZapLog.Error("创建导出文件失败", zap.Error(err))
		return err
	}
	defer file.Close()

	// 写入JSON数组开始标记
	if _, err := file.WriteString("[\n"); err != nil {
		global.ZapLog.Error("写入文件失败", zap.Error(err))
		return err
	}
	defer func() {
		// 写入JSON数组结束标记
		if _, err := file.WriteString("\n]"); err != nil {
			global.ZapLog.Error("写入文件失败", zap.Error(err))
		}
	}()

	// 1. 初始搜索请求 - 使用TypedClient的类型安全API
	pageSize := 1000
	res, err := global.ES.Search().
		Index(indexName).
		Query(&types.Query{MatchAll: &types.MatchAllQuery{}}).
		Size(pageSize).
		Scroll("1m").
		Do(ctx)

	if err != nil {
		global.ZapLog.Error("初始搜索请求失败", zap.Error(err))
		return err
	}

	// 提取Scroll ID
	scrollID = *res.ScrollId_

	// 处理第一批结果
	if len(res.Hits.Hits) > 0 {
		if err := writeHitsToFile(file, res.Hits.Hits, &first); err != nil {
			return err
		}
	}

	// 2. 滚动请求 - 使用类型安全的Scroll API
	for {
		res, err := global.ES.Scroll().
			ScrollId(scrollID).
			Scroll("1m").
			Do(ctx)

		if err != nil {
			global.ZapLog.Error("滚动请求失败", zap.Error(err))
			break
		}

		if len(res.Hits.Hits) == 0 {
			break // 没有更多结果
		}

		if err := writeHitsToFile(file, res.Hits.Hits, &first); err != nil {
			break
		}
	}

	// 3. 清理滚动上下文 - 使用类型安全的ClearScroll API
	_, err = global.ES.ClearScroll().
		ScrollId(scrollID).
		Do(ctx)
	if err != nil {
		global.ZapLog.Warn("清理滚动上下文失败", zap.Error(err))
	}

	global.ZapLog.Info("Elasticsearch数据导出完成", zap.String("path", path))
	return nil
}

// 辅助函数：将搜索结果写入文件
func writeHitsToFile(file *os.File, hits []types.Hit, first *bool) error {
	for _, hit := range hits {
		// 直接访问_source字段，无需类型断言
		data, err := json.Marshal(hit.Source_)
		if err != nil {
			global.ZapLog.Error("JSON序列化失败", zap.Error(err))
			continue
		}

		// 处理JSON数组逗号分隔
		if *first {
			*first = false
		} else {
			if _, err := file.WriteString(",\n"); err != nil {
				global.ZapLog.Error("写入文件失败", zap.Error(err))
				return err
			}
		}

		// 写入数据
		if _, err := file.Write(data); err != nil {
			global.ZapLog.Error("写入文件失败", zap.Error(err))
			return err
		}
	}
	return nil
}
