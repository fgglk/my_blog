package flag

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"go.uber.org/zap"

	"server/global"
	"server/model/es"
	//"server/service"
)

// ImportESData 从JSON文件导入数据到Elasticsearch
func importEsData(filePath string) error {
	// 读取JSON文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		global.ZapLog.Error("读取文件失败", zap.String("file_path", filePath), zap.Error(err))
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析JSON数据到文章切片
	var articles []es.ArticleES
	if err := json.Unmarshal(data, &articles); err != nil {
		global.ZapLog.Error("解析JSON失败", zap.String("file_path", filePath), zap.Error(err))
		return fmt.Errorf("解析JSON失败: %v", err)
	}

	// 批量导入到Elasticsearch
	total := len(articles)
	if total == 0 {
		global.ZapLog.Info("没有数据需要导入", zap.String("file_path", filePath))
		return nil
	}

	indexName := (&es.ArticleES{}).IndexName()

	//esService := service.NewArticleESService()
	// 检查索引是否存在
	// exists, err := esService.IndexExists()
	// if err != nil {
	// 	global.ZapLog.Error("检查索引存在性失败", zap.String("index_name", indexName), zap.Error(err))
	// 	return fmt.Errorf("检查索引存在性失败: %v", err)
	// }

	// // 如果索引存在则删除
	// if exists {
	// 	if err := esService.DeleteIndex(); err != nil {
	// 		global.ZapLog.Error("删除索引失败", zap.String("index_name", indexName), zap.Error(err))
	// 		return fmt.Errorf("删除索引失败: %v", err)
	// 	}
	// }

	// 创建新索引
	// if err := esService.CreateIndex(); err != nil {
	// 	global.ZapLog.Error("创建索引失败", zap.String("index_name", indexName), zap.Error(err))
	// 	return fmt.Errorf("创建索引失败: %v", err)
	// }

	// 构建批量请求数据
	var request bulk.Request
	for _, article := range articles {
		// 将uint64类型的ID转换为string并获取指针
		idStr := strconv.FormatUint(article.ID, 10)
		// 为每条数据创建索引操作，指定文档的ID
		request = append(request, types.OperationContainer{Index: &types.IndexOperation{Id_: &idStr}})
		// 添加文档数据到请求
		request = append(request, article)
	}

	// 使用Elasticsearch客户端执行批量操作
	_, err = global.ES.Bulk().
		Request(&request).       // 提交请求数据
		Index(indexName).        // 指定索引名称
		Refresh(refresh.True).   // 强制刷新索引以使文档立即可见
		Do(context.Background()) // 执行请求
	if err != nil {
		global.ZapLog.Error("批量导入失败", zap.String("index_name", indexName), zap.Int("total_records", total), zap.Error(err))
		return fmt.Errorf("批量导入失败: %v", err)
	}

	global.ZapLog.Info("数据导入成功", zap.String("index_name", indexName), zap.Int("total_records", total))
	return nil
}
