import request from '@/utils/request'
import type { ApiResponse, LikeResponse, FavoriteResponse, FavoriteListResponse } from '@/types/api'
import type { 
  Article, 
  ArticleListResponse, 
  CreateArticleRequest, 
  UpdateArticleRequest,
  ArticleListParams,
  SearchArticlesParams
} from '@/types/article'

// 文章API
export const articleApi = {
  // 获取文章列表
  getArticleList: (params: ArticleListParams): Promise<ApiResponse<ArticleListResponse>> => {
    return request.get('/articles', { params })
  },

  // 获取文章详情
  getArticle: (id: number): Promise<ApiResponse<Article>> => {
    return request.get(`/articles/${id}`)
  },

  // 搜索文章
  searchArticles: (params: SearchArticlesParams): Promise<ApiResponse<ArticleListResponse>> => {
    return request.get('/articles/search', { params })
  },

  // 创建文章
  createArticle: (data: CreateArticleRequest): Promise<ApiResponse> => {
    return request.post('/articles', data)
  },

  // 更新文章
  updateArticle: (id: number, data: UpdateArticleRequest): Promise<ApiResponse> => {
    return request.put(`/articles/${id}`, data)
  },

  // 删除文章
  deleteArticle: (id: number): Promise<ApiResponse> => {
    return request.delete(`/articles/${id}`)
  },

  // 点赞/取消点赞
  toggleLike: (articleId: number): Promise<ApiResponse<LikeResponse>> => {
    return request.post('/articles/like', { article_id: articleId })
  },

  // 收藏/取消收藏
  toggleFavorite: (articleId: number): Promise<ApiResponse<FavoriteResponse>> => {
    return request.post('/articles/favorite', { article_id: articleId })
  },

  // 上传图片
  uploadImage: (file: File): Promise<ApiResponse<{ id: number; url: string }>> => {
    const formData = new FormData()
    formData.append('image', file)
    return request.post('/image/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取用户文章列表
  getUserArticles: (page = 1, size = 10): Promise<ApiResponse<ArticleListResponse>> => {
    return request.get('/articles/my', { params: { page, size } })
  },

  // 根据用户ID获取文章列表
  getArticlesByUserId: (userId: number, page = 1, size = 10): Promise<ApiResponse<ArticleListResponse>> => {
    return request.get(`/articles/user/${userId}`, { params: { page, size } })
  },

  // 获取相关文章
  getRelatedArticles: (articleId: number): Promise<ApiResponse<ArticleListResponse>> => {
    return request.get(`/articles/${articleId}/related`)
  },

  // 获取用户收藏列表
  getFavorites: (params: { page?: number; size?: number; sort?: string }): Promise<ApiResponse<FavoriteListResponse>> => {
    return request.get('/articles/favorites', { params })
  }
} 