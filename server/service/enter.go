package service

// 导出所有service实例，供API层直接调用
type ServiceGroup struct {
	ArticleService
	UserService
	CommentService
	PageService
	CategoryService
	TagService
}

var ServiceGroups = new(ServiceGroup)
