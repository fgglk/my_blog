package appType

// CommentStatusType 评论状态类型
type CommentStatusType uint8

// 评论状态常量
const (
	CommentStatusPending  CommentStatusType = iota // 0-待审核
	CommentStatusApproved                          // 1-已发布
	CommentStatusRejected                          // 2-已拒绝
)

// String 实现String接口，便于日志和JSON序列化
func (s CommentStatusType) String() string {
	return [...]string{"pending", "approved", "rejected"}[s]
}
