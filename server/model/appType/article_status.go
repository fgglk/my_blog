package appType

import "fmt"

// ArticleStatus 文章状态类型
type ArticleStatus uint8

// 文章状态常量
const (
	StatusDraft     ArticleStatus = iota // 0 - 草稿
	StatusPublished                      // 1 - 已发布
	StatusArchived                       // 2 - 已归档
)

// String 将状态转换为字符串描述
func (s ArticleStatus) String() string {
	switch s {
	case StatusDraft:
		return "draft"
	case StatusPublished:
		return "published"
	case StatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

// FromString 将字符串转换为状态枚举
func (s *ArticleStatus) FromString(status string) error {
	switch status {
	case "draft":
		*s = StatusDraft
	case "published":
		*s = StatusPublished
	case "archived":
		*s = StatusArchived
	default:
		return fmt.Errorf("invalid article status: %s", status)
	}
	return nil
}
