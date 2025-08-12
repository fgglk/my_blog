package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"server/global"
	"server/model/appType"
	"server/model/database"
	"server/model/es"
	"server/model/request"

	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ArticleService 文章服务结构体
type ArticleService struct{}

// handleArticleCategory 处理文章分类
func (s *ArticleService) handleArticleCategory(tx *gorm.DB, categoryID uint) error {
	// 检查分类是否存在
	var category database.Category
	if err := tx.Where("id = ?", categoryID).First(&category).Error; err != nil {
		return errors.New("分类不存在: " + err.Error())
	}
	return nil
}

// handleArticleCategoryByName 按名称处理文章分类
func (s *ArticleService) handleArticleCategoryByName(tx *gorm.DB, categoryName string) (uint, error) {
	var category database.Category
	// 查找分类
	if err := tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
		return 0, errors.New("分类不存在: " + err.Error())
	}
	return category.ID, nil
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(req request.ArticleCreateRequest) (database.Article, error) {
	// 初始化文章结构体
	article := database.Article{
		BaseModelWithStatus: database.BaseModelWithStatus{
			Status: req.Status,
		},
		Title:        req.Title,
		Content:      req.Content,
		Summary:      req.Summary,
		CoverImage:   req.CoverImage, // 添加封面图片
		AuthorID:     req.AuthorID,
		ViewCount:    req.ViewCount,
		CommentCount: req.CommentCount,
		LikeCount:    req.LikeCount,
		// 生成唯一slug
		Slug: generateSlug(req.Title),
	}



	// 使用事务确保数据一致性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 处理分类
	var categoryID uint
	if req.CategoryID > 0 {
		// 如果提供了分类ID，验证分类是否存在
		categoryID = req.CategoryID
		if err := s.handleArticleCategory(tx, categoryID); err != nil {
			tx.Rollback()
			return article, err
		}
	} else {
		// 默认使用技术分类（ID为3，因为数据库中已有ID为3的分类）
		categoryID = 3
		if err := s.handleArticleCategory(tx, categoryID); err != nil {
			tx.Rollback()
			return article, err
		}
	}
	article.CategoryID = categoryID

	// 保存文章基本信息，使用 Create 方法
	if err := tx.Create(&article).Error; err != nil {
		tx.Rollback()
		return article, err
	}
	
	// 如果状态为0，显式更新状态字段以避免默认值
	if req.Status == 0 {
		if err := tx.Model(&article).Update("status", 0).Error; err != nil {
			tx.Rollback()
			return article, err
		}
	}

	// 处理标签关联
	if len(req.Tags) > 0 || len(req.TagNames) > 0 {
		if err := s.handleArticleTags(tx, article.ID, req.Tags, req.TagNames); err != nil {
			tx.Rollback()
			return article, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return article, err
	}

	// 事务提交后再同步到ES
	if article.Status == 1 {
		// 异步同步到ES
		go s.SyncArticleToES(article.ID)
	}

	// 重新查询文章以获取完整关联数据
	if err := global.DB.Preload("Category").Preload("Tags").Preload("Author").Where("id = ?", article.ID).First(&article).Error; err != nil {
		return article, err
	}

	return article, nil
}

// GetArticleList 分页获取文章列表
func (s *ArticleService) GetArticleList(req request.ArticleQueryRequest) ([]database.Article, int64, error) {
	var articles []database.Article
	var total int64
	// 默认查询已发布文章，管理员可以查看所有状态
	query := global.DB.Model(&database.Article{})

	// 条件筛选
	// 修改状态筛选逻辑，只有当Status > 0时才使用传入的状态
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	} else {
		query = query.Where("status = 1") // 默认只查询已发布文章
	}

	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.AuthorID > 0 {
		query = query.Where("author_id = ?", req.AuthorID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	// 验证分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 || req.Size > 100 {
		req.Size = 10 // 设置默认每页条数
	}
	offset := (req.Page - 1) * req.Size
	// 修复链式调用语法错误
	if err := query.Preload("Category").Preload("Author").Preload("Tags").
		Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetArticleByID 根据ID获取文章详情
func (s *ArticleService) GetArticleByID(id uint, isAdmin bool, currentUserID uint) (database.Article, error) {
	var article database.Article
	if err := global.DB.Preload("Category").Preload("Author").Preload("Tags").Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return article, errors.New("文章不存在")
		}
		return article, err
	}

	// 权限检查：未发布文章只能作者或管理员查看
	if article.Status != 1 && !isAdmin && article.AuthorID != currentUserID {
		return article, errors.New("无权访问此文章")
	}

	// 异步增加阅读量
	go s.IncrementViewCount(id)

	return article, nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(req request.ArticleUpdateRequest, userID uint, isAdmin bool) (database.Article, error) {
	var article database.Article
	if err := global.DB.Where("id = ?", req.ID).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return article, errors.New("文章不存在")
		}
		return article, err
	}

	// 权限检查
	if !isAdmin && article.AuthorID != userID {
		return article, errors.New("无权修改此文章")
	}

	// 使用事务更新文章
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查分类ID是否存在（如果提供了）
	if req.CategoryID > 0 {
		if err := s.handleArticleCategory(tx, req.CategoryID); err != nil {
			tx.Rollback()
			return article, err
		}
	}

	// 更新基本信息
	updateData := map[string]interface{}{}

	// 只有当字段不为空时才更新
	if req.Title != "" {
		updateData["Title"] = req.Title
	}
	if req.Content != "" {
		updateData["Content"] = req.Content
	}
	if req.Summary != "" {
		updateData["Summary"] = req.Summary
	}

	// 如果提供了封面图片，更新封面图片
	if req.CoverImage != "" {
		updateData["CoverImage"] = req.CoverImage
	}

	// 只有当CategoryID大于0时才更新分类ID
	if req.CategoryID > 0 {
		updateData["CategoryID"] = req.CategoryID
	}

	// 如果状态变更，更新状态
	statusChanged := false
	if article.Status != req.Status {
		updateData["Status"] = req.Status
		article.Status = req.Status
		statusChanged = true
	}

	if err := tx.Model(&article).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return article, err
	}

	// 更新标签关联（只有当明确提供了标签数据时才更新）
	if len(req.Tags) > 0 || len(req.TagNames) > 0 {
		// 获取旧的标签ID（在删除关联之前）
		var oldTagIDs []uint
		if err := tx.Model(&database.ArticleTag{}).Where("article_id = ?", article.ID).Pluck("tag_id", &oldTagIDs).Error; err != nil {
			tx.Rollback()
			return article, err
		}

		// 先删除旧关联（使用 Unscoped 确保硬删除）
		result := tx.Unscoped().Where("article_id = ?", article.ID).Delete(&database.ArticleTag{})
		if result.Error != nil {
			global.ZapLog.Error("删除文章标签关联失败", zap.Error(result.Error), zap.Uint("articleID", article.ID))
			tx.Rollback()
			return article, result.Error
		}
		global.ZapLog.Info("删除文章标签关联", zap.Uint("articleID", article.ID), zap.Int64("deletedCount", result.RowsAffected))

		// 添加新关联
		if err := s.handleArticleTags(tx, article.ID, req.Tags, req.TagNames); err != nil {
			tx.Rollback()
			return article, err
		}

		// 清理不再被任何文章使用的旧标签
		if err := s.cleanupOrphanTags(tx, oldTagIDs); err != nil {
			tx.Rollback()
			return article, err
		}
	}

	// 如果变为发布状态，同步到ES
	if statusChanged && article.Status == 1 {
		go s.SyncArticleToES(article.ID)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return article, err
	}

	// 重新加载完整文章数据
	global.DB.Preload("Category").Preload("Tags").First(&article, article.ID)

	return article, nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(articleID uint, userID uint, isAdmin bool) error {
	var article database.Article
	if err := global.DB.Where("id = ?", articleID).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 权限检查
	if !isAdmin && article.AuthorID != userID {
		return errors.New("无权删除此文章")
	}

	// 使用事务删除文章及关联数据
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取该文章使用的标签ID（在删除关联之前）
	var tagIDs []uint
	if err := tx.Model(&database.ArticleTag{}).Where("article_id = ?", articleID).Pluck("tag_id", &tagIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除标签关联
	if err := tx.Where("article_id = ?", articleID).Delete(&database.ArticleTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除文章记录
	if err := tx.Delete(&article).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 清理不再被任何文章使用的标签
	if err := s.cleanupOrphanTags(tx, tagIDs); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// 异步从ES删除文章
	go s.DeleteArticleFromES(articleID)

	return nil
}

// cleanupOrphanTags 清理不再被任何文章使用的标签
func (s *ArticleService) cleanupOrphanTags(tx *gorm.DB, tagIDs []uint) error {
	for _, tagID := range tagIDs {
		// 检查该标签是否还被其他文章使用
		var count int64
		if err := tx.Model(&database.ArticleTag{}).Where("tag_id = ?", tagID).Count(&count).Error; err != nil {
			return err
		}

		// 如果没有文章使用该标签，则删除该标签
		if count == 0 {
			if err := tx.Delete(&database.Tag{}, tagID).Error; err != nil {
				return err
			}
			global.ZapLog.Info("清理孤儿标签", zap.Uint("tagID", tagID))
		}
	}
	return nil
}

// handleArticleTags 处理文章与标签的关联关系（支持ID和名称）
func (s *ArticleService) handleArticleTags(tx *gorm.DB, articleID uint, tagIDs []uint, tagNames []string) error {
	// 合并标签ID和通过名称获取的ID
	var allTagIDs []uint

	// 处理标签ID
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			var tag database.Tag
			if err := tx.Where("id = ?", tagID).First(&tag).Error; err != nil {
				return errors.New("标签不存在: " + err.Error())
			}
			allTagIDs = append(allTagIDs, tagID)
		}
	}

	// 处理标签名称
	if len(tagNames) > 0 {
		for _, tagName := range tagNames {
			var tag database.Tag
			// 查找标签
			if err := tx.Where("name = ?", tagName).First(&tag).Error; err != nil {
				// 如果不存在，自动创建
				if errors.Is(err, gorm.ErrRecordNotFound) {
					tag = database.Tag{
						Name: tagName,
						Slug: generateSlug(tagName), // 为标签生成slug
					}
					if err := tx.Create(&tag).Error; err != nil {
						return errors.New("创建标签失败: " + err.Error())
					}
				} else {
					return errors.New("查询标签失败: " + err.Error())
				}
			}
			allTagIDs = append(allTagIDs, tag.ID)
		}
	}

	// 去重标签ID
	uniqueTagIDs := make([]uint, 0, len(allTagIDs))
	seen := make(map[uint]bool)
	for _, id := range allTagIDs {
		if !seen[id] {
			seen[id] = true
			uniqueTagIDs = append(uniqueTagIDs, id)
		}
	}

	// 批量创建文章标签关联
	var articleTags []database.ArticleTag
	for _, tagID := range uniqueTagIDs {
		articleTags = append(articleTags, database.ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
		})
	}

	if len(articleTags) > 0 {
		// 使用 ON DUPLICATE KEY IGNORE 避免主键冲突
		if err := tx.CreateInBatches(&articleTags, 100).Error; err != nil {
			// 如果批量插入失败，尝试逐个插入并忽略重复项
			for _, articleTag := range articleTags {
				if err := tx.Where("article_id = ? AND tag_id = ?", articleTag.ArticleID, articleTag.TagID).
					FirstOrCreate(&articleTag).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SyncArticleToES 同步文章到Elasticsearch
func (s *ArticleService) SyncArticleToES(articleID uint) error {
	// 获取完整文章信息
	var article database.Article
	if err := global.DB.Preload("Category").Preload("Tags").Where("id = ?", articleID).First(&article).Error; err != nil {
		return err
	}

	// 转换为ES文章结构
	esArticle := es.ArticleES{
		ID:            uint64(article.ID),
		Title:         article.Title,
		Content:       article.Content,
		Summary:       article.Summary,
		CategoryID:    article.CategoryID,
		UserID:        article.AuthorID,
		Status:        appType.ArticleStatus(article.Status),
		ViewCount:     article.ViewCount,
		LikeCount:     article.LikeCount,
		CommentCount:  article.CommentCount,  // 添加评论数字段
		FavoriteCount: article.FavoriteCount, // 假设需要从其他表获取收藏量
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
	}

	// 提取标签名称
	for _, tag := range article.Tags {
		esArticle.Tags = append(esArticle.Tags, tag.Name)
	}

	// 同步到ES
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := global.ES.Index(esArticle.IndexName()).
		Id(strconv.FormatUint(esArticle.ID, 10)).
		Document(esArticle).
		Do(ctx)

	if err != nil {
		global.ZapLog.Error("同步文章到ES失败", zap.Uint("articleID", articleID), zap.Error(err))
		return err
	}

	global.ZapLog.Info("文章同步到ES成功", zap.Uint("articleID", articleID))
	return nil
}

// DeleteArticleFromES 从Elasticsearch删除文章
func (s *ArticleService) DeleteArticleFromES(articleID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建删除请求
	deleteRequest := global.ES.Delete(
		(&es.ArticleES{}).IndexName(),
		strconv.FormatUint(uint64(articleID), 10),
	)

	// 执行删除操作
	_, err := deleteRequest.Do(ctx)

	if err != nil {
		global.ZapLog.Error("从ES删除文章失败", zap.Uint("articleID", articleID), zap.Error(err))
		return err
	}

	global.ZapLog.Info("文章从ES删除成功", zap.Uint("articleID", articleID))
	return nil
}

// IncrementViewCount 增加阅读量
func (s *ArticleService) IncrementViewCount(articleID uint) error {
	// 使用上下文设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用乐观锁更新阅读量
	result := global.DB.WithContext(ctx).Model(&database.Article{}).
		Where("id = ?", articleID).
		Update("view_count", gorm.Expr("view_count + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("更新阅读量失败，文章不存在或已被删除")
	}

	return nil
}

// generateSlug 基于标题生成唯一slug
func generateSlug(title string) string {
	// 基本slug生成（转换为小写，替换空格为连字符）
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	// 移除特殊字符（可以根据需要添加更多处理）
	// 添加时间戳确保唯一性
	slug = slug + "-" + time.Now().Format("20060102150405") + "-" + uuid.New().String()[:8]
	return slug
}

// SyncArticleStatsToES 同步文章统计数据到ES
func (s *ArticleService) SyncArticleStatsToES(articleID uint) error {
	// 获取完整文章信息（包括统计数据）
	var article database.Article
	if err := global.DB.Preload("Category").Preload("Tags").Where("id = ?", articleID).First(&article).Error; err != nil {
		return err
	}

	// 转换为ES文章结构
	esArticle := es.ArticleES{
		ID:            uint64(article.ID),
		Title:         article.Title,
		Content:       article.Content,
		Summary:       article.Summary,
		CategoryID:    article.CategoryID,
		UserID:        article.AuthorID,
		Status:        appType.ArticleStatus(article.Status),
		ViewCount:     article.ViewCount,
		LikeCount:     article.LikeCount,
		CommentCount:  article.CommentCount,
		FavoriteCount: article.FavoriteCount,
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
	}

	// 提取标签名称
	for _, tag := range article.Tags {
		esArticle.Tags = append(esArticle.Tags, tag.Name)
	}

	// 同步到ES（使用Index而不是Update，这样如果文档不存在会创建，存在会更新）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := global.ES.Index(esArticle.IndexName()).
		Id(strconv.FormatUint(esArticle.ID, 10)).
		Document(esArticle).
		Do(ctx)

	if err != nil {
		global.ZapLog.Error("同步文章统计数据到ES失败", zap.Uint("articleID", articleID), zap.Error(err))
		return err
	}

	global.ZapLog.Info("文章统计数据同步到ES成功", zap.Uint("articleID", articleID))
	return nil
}

// ToggleLike 切换文章点赞状态
func (s *ArticleService) ToggleLike(articleID uint, userID uint) (bool, error) {
	// 使用事务确保数据一致性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查文章是否存在
	var article database.Article
	if err := tx.Where("id = ?", articleID).First(&article).Error; err != nil {
		tx.Rollback()
		return false, errors.New("文章不存在")
	}

	// 检查是否已经点赞
	var like database.Like
	exists := tx.Where("article_id = ? AND user_id = ?", articleID, userID).First(&like).Error == nil

	if exists {
		// 取消点赞
		if err := tx.Delete(&like).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		// 减少点赞数
		article.LikeCount--
	} else {
		// 新增点赞
		like = database.Like{
			ArticleID: articleID,
			UserID:    userID,
		}
		if err := tx.Create(&like).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		// 增加点赞数
		article.LikeCount++
	}

	// 更新文章点赞数
	if err := tx.Model(&article).Update("like_count", article.LikeCount).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 异步同步到ES
	go s.SyncArticleStatsToES(articleID)

	return !exists, nil // 返回是否点赞成功
}

// ToggleFavorite 切换文章收藏状态
func (s *ArticleService) ToggleFavorite(articleID uint, userID uint) (bool, error) {
	// 使用事务确保数据一致性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查文章是否存在
	var article database.Article
	if err := tx.Where("id = ?", articleID).First(&article).Error; err != nil {
		tx.Rollback()
		return false, errors.New("文章不存在")
	}

	// 检查是否已经收藏
	var favorite database.Favorite
	exists := tx.Where("article_id = ? AND user_id = ?", articleID, userID).First(&favorite).Error == nil

	if exists {
		// 取消收藏
		if err := tx.Delete(&favorite).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		// 减少收藏数
		article.FavoriteCount--
	} else {
		// 新增收藏
		favorite = database.Favorite{
			ArticleID: articleID,
			UserID:    userID,
		}
		if err := tx.Create(&favorite).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		// 增加收藏数
		article.FavoriteCount++
	}

	// 更新文章收藏数
	if err := tx.Model(&article).Update("favorite_count", article.FavoriteCount).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 异步同步到ES
	go s.SyncArticleStatsToES(articleID)

	return !exists, nil // 返回是否收藏成功
}

func (s *ArticleService) SearchArticles(req request.SearchArticleRequest) (es.ArticleSearchResult, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 10
	}

	// 构建查询条件 (使用types.Query结构体)
	boolQuery := types.BoolQuery{}

	// 必须匹配已发布状态
	boolQuery.Must = append(boolQuery.Must, types.Query{Term: map[string]types.TermQuery{"status": {Value: appType.StatusPublished}}})
	boost1 := float32(3.0)
	boost2 := float32(2.0)
	boost3 := float32(1.5)
	// 关键词搜索 (多字段)
	if req.Keyword != "" {
		boolQuery.Should = append(boolQuery.Should,
			types.Query{Match: map[string]types.MatchQuery{"title": {Query: req.Keyword, Boost: &boost1}}},
			types.Query{Match: map[string]types.MatchQuery{"content": {Query: req.Keyword}}},
			types.Query{Match: map[string]types.MatchQuery{"summary": {Query: req.Keyword, Boost: &boost2}}},
			types.Query{Match: map[string]types.MatchQuery{"tags": {Query: req.Keyword, Boost: &boost3}}},
		)
		boolQuery.MinimumShouldMatch = "1"
	}

	// 分类筛选
	if req.CategoryID > 0 {
		boolQuery.Filter = append(boolQuery.Filter, types.Query{Term: map[string]types.TermQuery{"category_id": {Value: req.CategoryID}}})
	}

	// 标签筛选
	if req.Tag != "" {
		boolQuery.Filter = append(boolQuery.Filter, types.Query{Term: map[string]types.TermQuery{"tags": {Value: req.Tag}}})
	}

	// 设置排序字段和顺序
	from := (req.Page - 1) * req.Size
	size := req.Size

	// 声明排序字段
	var sortField string
	switch req.Sort {
	case "time":
		sortField = "created_at"
	case "view":
		sortField = "view_count"
	case "comment":
		sortField = "comment_count"
	case "like":
		sortField = "like_count"
	default:
		sortField = "created_at"
	}

	// 使用枚举设置排序方向
	var sortDirection string
	sortDirection = "desc"
	if req.Order == "asc" {
		sortDirection = "asc"
	}

	// 构建排序字符串
	sortStr := fmt.Sprintf("%s:%s", sortField, sortDirection)

	// 执行查询
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建完整查询
	searchQuery := global.ES.Search().
		Index((&es.ArticleES{}).IndexName()).
		Sort(sortStr).
		Request(&search.Request{
			Query: &types.Query{Bool: &boolQuery}, // 使用types.Query类型
			From:  &from,                          // 传递指针
			Size:  &size,                          // 传递指针
		}).
		SourceIncludes_("id", "title", "content", "summary", "category_id", "tags", "user_id", "status", "view_count", "comment_count", "like_count", "favorite_count", "created_at", "updated_at")

	// 执行查询并处理响应
	resp, err := searchQuery.Do(ctx)
	if err != nil {
		global.ZapLog.Error("Elasticsearch查询错误",
			zap.Error(err),
			zap.String("query", fmt.Sprintf("%+v", searchQuery)),
			zap.Stack("stack"))
		return es.ArticleSearchResult{}, err
	}

	// 解析结果
	result := es.ArticleSearchResult{
		Total:     resp.Hits.Total.Value,
		Page:      req.Page,
		Size:      req.Size,
		TotalPage: (int(resp.Hits.Total.Value) + req.Size - 1) / req.Size,
	}

	// 处理查询结果
	result.Articles = make([]es.ArticleES, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var article es.ArticleES
		if err := json.Unmarshal(hit.Source_, &article); err != nil {
			global.ZapLog.Warn("解析文章数据失败", zap.Error(err), zap.String("hit_id", *hit.Id_))
			continue
		}
		result.Articles = append(result.Articles, article)
	}

	return result, nil
}

// GetWebsiteStats 获取网站统计数据
func (s *ArticleService) GetWebsiteStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 文章总数（已发布）
	var articleCount int64
	if err := global.DB.Model(&database.Article{}).Where("status = ?", 1).Count(&articleCount).Error; err != nil {
		return nil, err
	}
	stats["articleCount"] = int(articleCount)

	// 总阅读量
	var totalViews int64
	if err := global.DB.Model(&database.Article{}).Where("status = ?", 1).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews).Error; err != nil {
		return nil, err
	}
	stats["viewCount"] = int(totalViews)

	// 总评论数
	var totalComments int64
	if err := global.DB.Model(&database.Article{}).Where("status = ?", 1).Select("COALESCE(SUM(comment_count), 0)").Scan(&totalComments).Error; err != nil {
		return nil, err
	}
	stats["commentCount"] = int(totalComments)

	// 总点赞数
	var totalLikes int64
	if err := global.DB.Model(&database.Article{}).Where("status = ?", 1).Select("COALESCE(SUM(like_count), 0)").Scan(&totalLikes).Error; err != nil {
		return nil, err
	}
	stats["likeCount"] = int(totalLikes)

	// 总收藏数
	var totalFavorites int64
	if err := global.DB.Model(&database.Article{}).Where("status = ?", 1).Select("COALESCE(SUM(favorite_count), 0)").Scan(&totalFavorites).Error; err != nil {
		return nil, err
	}
	stats["favoriteCount"] = int(totalFavorites)

	return stats, nil
}

// GetUserArticles 获取用户文章列表
func (s *ArticleService) GetUserArticles(userID uint, page, size int) ([]database.Article, int64, error) {
	var articles []database.Article
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	if err := global.DB.Model(&database.Article{}).Where("author_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询文章列表，预加载关联数据
	if err := global.DB.
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetRelatedArticles 获取相关文章
// 规则：1. 只显示同分类文章 2. 最多显示4篇 3. 优先显示同标签文章 4. 不足4篇按实际数量显示
func (s *ArticleService) GetRelatedArticles(articleID uint) ([]database.Article, error) {
	// 首先获取当前文章信息
	var currentArticle database.Article
	if err := global.DB.Preload("Category").Preload("Tags").Where("id = ? AND status = ?", articleID, 1).First(&currentArticle).Error; err != nil {
		return nil, err
	}

	// 获取同分类的所有文章（排除当前文章）
	var sameCategoryArticles []database.Article
	if err := global.DB.
		Preload("Category").
		Preload("Tags").
		Preload("Author").
		Where("category_id = ? AND id != ? AND status = ?", currentArticle.CategoryID, articleID, 1).
		Order("created_at DESC").
		Find(&sameCategoryArticles).Error; err != nil {
		return nil, err
	}

	// 如果没有同分类文章，返回空数组
	if len(sameCategoryArticles) == 0 {
		return []database.Article{}, nil
	}

	// 如果同分类文章少于等于4篇，直接返回
	if len(sameCategoryArticles) <= 4 {
		return sameCategoryArticles, nil
	}

	// 如果同分类文章多于4篇，优先显示同标签的文章
	var relatedArticles []database.Article
	var otherArticles []database.Article

	// 获取当前文章的标签ID列表
	currentTagIDs := make(map[uint]bool)
	for _, tag := range currentArticle.Tags {
		currentTagIDs[tag.ID] = true
	}

	// 分类文章：有相同标签的优先
	for _, article := range sameCategoryArticles {
		hasCommonTag := false
		for _, tag := range article.Tags {
			if currentTagIDs[tag.ID] {
				hasCommonTag = true
				break
			}
		}
		
		if hasCommonTag {
			relatedArticles = append(relatedArticles, article)
		} else {
			otherArticles = append(otherArticles, article)
		}
	}

	// 组合结果：优先显示有相同标签的文章，然后补充其他文章
	result := make([]database.Article, 0, 4)
	
	// 先添加有相同标签的文章（最多4篇）
	for i := 0; i < len(relatedArticles) && len(result) < 4; i++ {
		result = append(result, relatedArticles[i])
	}
	
	// 如果还不够4篇，从其他文章中补充
	for i := 0; i < len(otherArticles) && len(result) < 4; i++ {
		result = append(result, otherArticles[i])
	}

	return result, nil
}
