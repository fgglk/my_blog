// 阅读设置接口
export interface ReadingSettings {
  allowComments: boolean
  allowRepost: boolean
  requireLogin: boolean
}

// 文章接口
export interface Article {
  id: number
  title: string
  content: string
  summary?: string // 用于存储阅读设置
  category_id: number
  category: Category
  tags: Tag[]
  author_id: number
  author_name: string
  author_avatar?: string
  cover_image?: string
  view_count: number
  like_count: number
  comment_count: number
  favorite_count: number
  is_published: boolean
  is_liked: boolean
  is_favorited: boolean
  status?: string
  created_at: string
  updated_at: string
}

// 分类接口
export interface Category {
  id: number
  name: string
}

// 标签接口
export interface Tag {
  id: number
  name: string
}

// 文章列表响应
export interface ArticleListResponse {
  list: Article[]
  total: number
  page: number
  size: number
  totalPage: number
}

// 创建文章请求参数
export interface CreateArticleRequest {
  title: string
  content: string
  summary?: string // 用于存储阅读设置
  category_id: number
  tags: number[]
  status: number
}

// 更新文章请求参数
export interface UpdateArticleRequest {
  title?: string
  content?: string
  categoryId?: number
  tagIds?: number[]
  isPublished?: boolean
}

// 文章列表查询参数
export interface ArticleListParams {
  page: number
  size: number
  categoryId?: number
  tagId?: number
  keyword?: string
}

// 搜索文章参数
export interface SearchArticlesParams {
  keyword: string
  page: number
  size: number
} 