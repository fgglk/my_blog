package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"server/global"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"
)

type ArticleApi struct{}

var articleService = service.ServiceGroups.ArticleService

// @Summary 创建文章
// @Description 创建新文章，需要认证
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.ArticleCreateRequest true "文章创建信息"
// @Success 200 {object} response.Response{data=response.ArticleResponse}
// @Router /api/articles [post]
func (a *ArticleApi) CreateArticle(c *gin.Context) {
	var req request.ArticleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			response.FailWithMessage("参数错误: "+utils.TranslateValidationError(validateErr), c)
			return
		}
		response.FailWithMessage("请求参数格式错误", c)

		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		fmt.Printf("获取用户ID失败: %v\n", err)
		response.NoAuth(err.Error(), c)
		return
	}

	fmt.Printf("文章创建API: 用户ID=%d\n", userID)

	// 将用户ID设置到请求对象中
	req.AuthorID = userID

	// 调用服务层创建文章
	article, err := articleService.CreateArticle(req)
	if err != nil {
		response.FailWithMessage("创建文章失败: "+err.Error(), c)
		return
	}

	// 直接使用article对象中的关联数据，无需额外查询
	// 注意：这需要确保ArticleService.CreateArticle方法已经预加载了相关关联数据
	response.OkWithData(response.ToArticleResponse(article, article.Category, article.Tags, article.Author.Username, userID), c)
}

// @Summary 获取文章列表
// @Description 分页获取文章列表，支持标题、分类、标签筛选
// @Tags article
// @Accept json
// @Produce json
// @Param title query string false "文章标题模糊查询"
// @Param category_id query int false "分类ID"
// @Param tag_id query int false "标签ID"
// @Param page query int true "页码"
// @Param size query int true "每页条数"
// @Success 200 {object} response.Response{data=response.ArticleListResponse}
// @Router /api/articles [get]
func (a *ArticleApi) GetArticleList(c *gin.Context) {
	var req request.ArticleQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 添加默认值处理
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 10 // 设置默认每页10条
	}

	// 调用服务层获取文章列表
	articles, total, err := articleService.GetArticleList(req)
	if err != nil {
		response.FailWithMessage("获取文章列表失败: "+err.Error(), c)
		return
	}

	// 转换为响应模型
	var articleResponses []response.ArticleResponse
	for _, article := range articles {
		// 修复：将转换结果添加到切片，而非直接发送响应
		articleResponses = append(articleResponses, response.ToArticleResponse(article, article.Category, article.Tags, article.Author.Username, 0))
	}

	// 构造分页响应
	listResponse := response.ArticleListResponse{
		List:      articleResponses,
		Total:     total,
		Page:      req.Page,
		Size:      req.Size,
		TotalPage: (int(total) + req.Size - 1) / req.Size,
	}

	response.OkWithData(listResponse, c)
}

// @Summary 获取文章详情
// @Description 根据ID获取文章详细信息
// @Tags article
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{data=response.ArticleResponse}
// @Router /api/articles/{id} [get]
func (a *ArticleApi) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的文章ID", c)
		return
	}

	// 尝试获取用户ID，如果失败则设置为0（未登录用户）
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		// 处理未登录或登录过期的情况
		currentUserID = 0
	}

	isAdmin := false
	if currentUserID != 0 {
		isAdmin = utils.IsAdmin(currentUserID)
	}

	// 调用方法时传入完整参数
	article, err := articleService.GetArticleByID(uint(id), isAdmin, currentUserID)
	if err != nil {
		response.FailWithMessage("获取文章失败: "+err.Error(), c)
		return
	}

	// 获取关联数据
	articleResponse := response.ToArticleResponse(article, article.Category, article.Tags, article.Author.Username, currentUserID)

	response.OkWithData(articleResponse, c)
}

// @Summary 更新文章
// @Description 更新文章信息，需要认证和作者权限
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Param data body request.ArticleUpdateRequest true "文章更新信息"
// @Success 200 {object} response.Response{data=response.ArticleResponse}
// @Router /api/articles/{id} [put]
func (a *ArticleApi) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	// 使用utils.StringToUint替代直接转换
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage("无效的文章ID", c)
		return
	}

	var req request.ArticleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			response.FailWithMessage("参数错误: "+utils.TranslateValidationError(validateErr), c)
			return
		}
		response.FailWithMessage("请求参数格式错误", c)
		return
	}
	req.ID = id

	// 使用utils.GetUserID获取用户ID并处理错误
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	// 使用utils.IsAdmin检查管理员权限
	isAdmin := utils.IsAdmin(currentUserID)

	// 调用方法时传入完整参数
	article, err := articleService.GetArticleByID(id, isAdmin, currentUserID)
	if err != nil {
		response.FailWithMessage("文章不存在", c)
		return
	}
	if article.AuthorID != currentUserID {
		response.Forbidden("没有权限修改此文章", c)
		return
	}

	// 调用服务层更新文章
	updatedArticle, err := articleService.UpdateArticle(req, currentUserID, isAdmin)
	if err != nil {
		response.FailWithMessage("更新文章失败: "+err.Error(), c)
		return
	}

	response.OkWithData(response.ToArticleResponse(updatedArticle, updatedArticle.Category, updatedArticle.Tags, updatedArticle.Author.Username, currentUserID), c)
}

// @Summary 删除文章
// @Description 删除文章，需要认证和作者权限
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/articles/{id} [delete]
func (a *ArticleApi) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	// 使用utils.StringToUint替代直接转换
	id, err := utils.StringToUint(idStr)
	if err != nil {
		response.FailWithMessage("无效的文章ID", c)
		return
	}

	// 使用utils.GetUserID获取用户ID并处理错误
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	// 使用utils.IsAdmin检查管理员权限
	isAdmin := utils.IsAdmin(currentUserID)

	// 调用方法时传入完整参数
	article, err := articleService.GetArticleByID(id, isAdmin, currentUserID)
	if err != nil {
		response.FailWithMessage("文章不存在", c)
		return
	}
	if article.AuthorID != currentUserID {
		response.Forbidden("没有权限删除此文章", c)
		return
	}

	// 调用服务层删除文章
	if err := articleService.DeleteArticle(id, currentUserID, isAdmin); err != nil {
		response.FailWithMessage("删除文章失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("文章删除成功", c)
}

// @Summary 搜索文章
// @Description 全文搜索文章，支持关键词和分类筛选
// @Tags article
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param category_id query int false "分类ID"
// @Param page query int false "页码"
// @Param size query int false "每页条数"
// @Success 200 {object} response.Response{data=es.ArticleSearchResult}
// @Router /api/articles/search [get]
func (a *ArticleApi) SearchArticles(c *gin.Context) {
	var req request.SearchArticleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 调用服务层搜索文章
	result, err := articleService.SearchArticles(req)
	if err != nil {
		response.FailWithMessage("搜索失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// @Summary 切换文章点赞状态
// @Description 点赞或取消点赞文章，需要认证
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.ToggleLikeRequest true "文章ID"
// @Success 200 {object} response.Response{data=map[string]bool}
// @Router /api/articles/toggle-like [post]
func (a *ArticleApi) ToggleLike(c *gin.Context) {
	var req request.ToggleLikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			response.FailWithMessage("参数错误: "+utils.TranslateValidationError(validateErr), c)
			return
		}
		response.FailWithMessage("请求参数格式错误", c)
		return
	}

	// 使用utils.GetUserID获取用户ID并处理错误
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	// 调用服务层切换点赞状态
	liked, err := articleService.ToggleLike(req.ArticleID, currentUserID)
	if err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}

	response.OkWithData(map[string]bool{"liked": liked}, c)
}

// @Summary 切换文章收藏状态
// @Description 收藏或取消收藏文章，需要认证
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.ToggleFavoriteRequest true "文章ID"
// @Success 200 {object} response.Response{data=map[string]bool}
// @Router /api/articles/toggle-favorite [post]
func (a *ArticleApi) ToggleFavorite(c *gin.Context) {
	var req request.ToggleFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			response.FailWithMessage("参数错误: "+utils.TranslateValidationError(validateErr), c)
			return
		}
		response.FailWithMessage("请求参数格式错误", c)
		return
	}

	// 使用utils.GetUserID获取用户ID并处理错误
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	// 调用服务层切换收藏状态
	favorited, err := articleService.ToggleFavorite(req.ArticleID, currentUserID)
	if err != nil {
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}

	response.OkWithData(map[string]bool{"favorited": favorited}, c)
}

// @Summary 获取用户收藏列表
// @Description 获取当前用户的收藏文章列表
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param size query int false "每页条数"
// @Param sort query string false "排序方式"
// @Success 200 {object} response.Response{data=response.FavoriteListResponse}
// @Router /api/favorites [get]
func (a *ArticleApi) GetUserFavorites(c *gin.Context) {
	// 获取当前用户ID
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	sort := c.DefaultQuery("sort", "created_at")

	// 验证参数
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 调用服务层获取收藏列表
	favorites, total, err := articleService.GetUserFavorites(currentUserID, page, size, sort)
	if err != nil {
		response.FailWithMessage("获取收藏列表失败: "+err.Error(), c)
		return
	}

	// 添加调试日志
	global.ZapLog.Info("获取收藏列表",
		zap.Uint("userID", currentUserID),
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("sort", sort),
		zap.Int64("total", total),
		zap.Int("count", len(favorites)),
	)

	// 调试第一个收藏记录的作者信息
	if len(favorites) > 0 {
		firstFavorite := favorites[0]
		global.ZapLog.Info("第一个收藏记录的作者信息",
			zap.Uint("articleID", firstFavorite.ArticleID),
			zap.String("authorUsername", firstFavorite.Article.Author.Username),
			zap.String("authorNickname", firstFavorite.Article.Author.Nickname),
			zap.String("authorAvatar", firstFavorite.Article.Author.Avatar),
		)
	}

	// 转换收藏记录为响应格式
	var favoriteResponses []response.FavoriteResponse
	for _, favorite := range favorites {
		favoriteResponses = append(favoriteResponses, response.ToFavoriteResponse(favorite, currentUserID))
	}

	// 构建响应数据
	result := response.FavoriteListResponse{
		List:      favoriteResponses,
		Total:     total,
		Page:      page,
		Size:      size,
		TotalPage: (int(total) + size - 1) / size,
	}

	response.OkWithData(result, c)
}

// @Summary 取消收藏
// @Description 取消收藏指定文章
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "收藏记录ID"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/favorites/{id} [delete]
func (a *ArticleApi) RemoveFavorite(c *gin.Context) {
	// 获取收藏记录ID
	favoriteIDStr := c.Param("id")
	favoriteID, err := strconv.ParseUint(favoriteIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的收藏记录ID", c)
		return
	}

	// 获取当前用户ID
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	// 调用服务层取消收藏
	if err := articleService.RemoveFavorite(uint(favoriteID), currentUserID); err != nil {
		response.FailWithMessage("取消收藏失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("取消收藏成功", c)
}

// @Summary 获取网站统计数据
// @Description 获取网站的文章数量、阅读量、评论数、点赞数、收藏数等统计数据
// @Tags article
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/articles/stats [get]
func (a *ArticleApi) GetWebsiteStats(c *gin.Context) {
	stats, err := articleService.GetWebsiteStats()
	if err != nil {
		response.FailWithMessage("获取统计数据失败: "+err.Error(), c)
		return
	}

	response.OkWithData(stats, c)
}

// @Summary 获取用户文章列表
// @Description 获取当前登录用户的文章列表，需要认证
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码，默认为1"
// @Param size query int false "每页条数，默认为10"
// @Success 200 {object} response.Response{data=response.ArticleListResponse}
// @Router /api/articles/my [get]
func (a *ArticleApi) GetUserArticles(c *gin.Context) {
	// 获取当前用户ID
	userID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 验证参数
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	// 允许更大的size用于获取用户所有文章，但限制最大值防止性能问题
	if size > 1000 {
		size = 1000
	}

	// 调用服务层获取用户文章
	articles, total, err := articleService.GetUserArticles(userID, page, size)
	if err != nil {
		response.FailWithMessage("获取用户文章失败: "+err.Error(), c)
		return
	}

	// 转换为响应模型
	var articleResponses []response.ArticleResponse
	for _, article := range articles {
		articleResponses = append(articleResponses, response.ToArticleResponse(article, article.Category, article.Tags, article.Author.Username, userID))
	}

	// 构建响应数据
	result := response.ArticleListResponse{
		List:  articleResponses,
		Total: total,
		Page:  page,
		Size:  size,
	}

	response.OkWithData(result, c)
}

// GetArticlesByUserID 根据用户ID获取文章列表
// @Description 根据用户ID获取该用户的文章列表，公开接口
// @Tags article
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param page query int false "页码，默认为1"
// @Param size query int false "每页条数，默认为10"
// @Success 200 {object} response.Response{data=response.ArticleListResponse}
// @Router /api/articles/user/{user_id} [get]
func (a *ArticleApi) GetArticlesByUserID(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的用户ID", c)
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 验证参数
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	// 允许更大的size用于获取用户所有文章，但限制最大值防止性能问题
	if size > 1000 {
		size = 1000
	}

	// 调用服务层获取用户文章（只获取已发布的文章）
	articles, total, err := articleService.GetArticlesByUserID(uint(userID), page, size)
	if err != nil {
		response.FailWithMessage("获取用户文章失败: "+err.Error(), c)
		return
	}

	// 转换为响应模型
	var articleResponses []response.ArticleResponse
	for _, article := range articles {
		articleResponses = append(articleResponses, response.ToArticleResponse(article, article.Category, article.Tags, article.Author.Username, 0))
	}

	// 构建响应数据
	result := response.ArticleListResponse{
		List:  articleResponses,
		Total: total,
		Page:  page,
		Size:  size,
	}

	response.OkWithData(result, c)
}

// @Summary 获取相关文章
// @Description 根据当前文章获取相关文章，优先显示同分类同标签的文章
// @Tags article
// @Accept json
// @Produce json
// @Param id path int true "当前文章ID"
// @Success 200 {object} response.Response{data=response.ArticleListResponse}
// @Router /api/articles/{id}/related [get]
func (a *ArticleApi) GetRelatedArticles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的文章ID", c)
		return
	}

	// 调用服务层获取相关文章
	articles, err := articleService.GetRelatedArticles(uint(id))
	if err != nil {
		response.FailWithMessage("获取相关文章失败: "+err.Error(), c)
		return
	}

	// 转换为响应模型
	var articleResponses []response.ArticleResponse
	for _, article := range articles {
		articleResponses = append(articleResponses, response.ToArticleResponse(article, article.Category, article.Tags, article.Author.Username, 0))
	}

	// 构建响应数据
	result := response.ArticleListResponse{
		List:  articleResponses,
		Total: int64(len(articleResponses)),
		Page:  1,
		Size:  len(articleResponses),
	}

	response.OkWithData(result, c)
}

// @Summary 同步所有已发布文章到ES
// @Description 将数据库中所有已发布的文章同步到Elasticsearch
// @Tags article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{msg=string}
// @Router /api/articles/sync-es [post]
func (a *ArticleApi) SyncAllArticlesToES(c *gin.Context) {
	// 检查管理员权限
	currentUserID, err := utils.GetUserID(c)
	if err != nil {
		response.NoAuth(err.Error(), c)
		return
	}

	if !utils.IsAdmin(currentUserID) {
		response.Forbidden("需要管理员权限", c)
		return
	}

	// 异步执行同步任务
	go func() {
		if err := articleService.SyncAllPublishedArticlesToES(); err != nil {
			global.ZapLog.Error("同步所有文章到ES失败", zap.Error(err))
		}
	}()

	response.OkWithMessage("文章同步任务已启动，请查看日志了解进度", c)
}
