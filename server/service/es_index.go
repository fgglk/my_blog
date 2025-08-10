package service

import (
	"context"
	"fmt"
	"server/global"
	"server/model/es"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

type ArticleESService struct {
	client *elasticsearch.TypedClient
	index  string
}

func NewArticleESService() *ArticleESService {
	return &ArticleESService{
		client: global.ES,
		index:  (&es.ArticleES{}).IndexName(),
	}
}

func (s *ArticleESService) DeleteIndex() error {
	global.ZapLog.Info("开始删除索引", zap.String("index", s.index))

	res, err := s.client.Indices.Delete(s.index).Do(context.Background())
	if err != nil {
		return fmt.Errorf("删除索引请求失败: %v", err)
	}

	if !res.Acknowledged {
		return fmt.Errorf("删除索引未被确认")
	}

	global.ZapLog.Info("索引删除成功", zap.String("index", s.index))
	return nil
}

func (s *ArticleESService) CreateIndex() error {
	exists, err := s.IndexExists()
	if err != nil {
		return fmt.Errorf("检查索引存在性失败: %v", err)
	}
	if exists {
		global.ZapLog.Info("索引已存在，自动删除并重新创建", zap.String("index", s.index))

		// 调用独立的删除索引方法
		if err := s.DeleteIndex(); err != nil {
			return fmt.Errorf("删除索引失败: %v", err)
		}
	}

	mapping := es.GetMapping()

	res, err := s.client.Indices.Create(s.index).Mappings(mapping).Do(context.Background())
	if err != nil {
		return fmt.Errorf("创建索引请求失败: %v", err)
	}

	if !res.Acknowledged {
		return fmt.Errorf("索引创建未被确认")
	}

	global.ZapLog.Info("索引创建成功", zap.String("index", s.index))
	return nil
}

func (s *ArticleESService) IndexExists() (bool, error) {
	return s.client.Indices.Exists(s.index).Do(context.Background())
}
