package service

import (
	"server/global"
	"server/model/database"
	"go.uber.org/zap"
)

type TagService struct{}

// GetTagListWithCount 获取标签列表及其文章数量
func (s *TagService) GetTagListWithCount() ([]database.TagWithCount, error) {
	var tags []database.TagWithCount
	
	// 使用原生SQL查询标签及其文章数量，排除软删除的标签
	query := `
		SELECT 
			t.id,
			t.name,
			t.slug,
			t.count,
			t.created_at,
			t.updated_at,
			COUNT(DISTINCT at.article_id) as article_count
		FROM tags t
		LEFT JOIN article_tags at ON t.id = at.tag_id
		LEFT JOIN articles a ON at.article_id = a.id AND a.status = 1
		WHERE t.deleted_at IS NULL
		GROUP BY t.id, t.name, t.slug, t.count, t.created_at, t.updated_at
		ORDER BY t.count DESC, t.id ASC
	`
	
	err := global.DB.Raw(query).Scan(&tags).Error
	return tags, err
}

// GetTagByID 根据ID获取标签
func (s *TagService) GetTagByID(id uint) (database.Tag, error) {
	var tag database.Tag
	err := global.DB.Where("id = ?", id).First(&tag).Error
	return tag, err
}

// CleanupOrphanTags 批量清理不再被任何文章使用的孤儿标签
func (s *TagService) CleanupOrphanTags() (int64, error) {
	// 查找所有孤儿标签（没有被任何文章使用的标签），排除已软删除的标签
	var orphanTagIDs []uint
	query := `
		SELECT t.id 
		FROM tags t 
		LEFT JOIN article_tags at ON t.id = at.tag_id 
		WHERE t.deleted_at IS NULL AND at.tag_id IS NULL
	`
	
	if err := global.DB.Raw(query).Pluck("id", &orphanTagIDs).Error; err != nil {
		return 0, err
	}
	
	if len(orphanTagIDs) == 0 {
		return 0, nil
	}
	
	// 删除孤儿标签
	result := global.DB.Delete(&database.Tag{}, orphanTagIDs)
	if result.Error != nil {
		return 0, result.Error
	}
	
	global.ZapLog.Info("批量清理孤儿标签完成", 
		zap.Int("清理数量", len(orphanTagIDs)),
		zap.Any("标签IDs", orphanTagIDs))
	
	return result.RowsAffected, nil
} 