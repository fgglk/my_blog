import type { Article } from './article'

// API响应接口
export interface ApiResponse<T = any> {
  code: number
  data: T
  msg: string
}

// 分页参数
export interface PageParams {
  page: number
  size: number
}

// 分页响应
export interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  size: number
  totalPage: number
}

// 点赞响应
export interface LikeResponse {
  liked: boolean
}

// 收藏响应
export interface FavoriteResponse {
  favorited: boolean
}

// 收藏记录
export interface Favorite {
  id: number
  article_id: number
  user_id: number
  created_at: string
  updated_at: string
  article: Article
}

// 收藏列表响应
export interface FavoriteListResponse {
  list: Favorite[]
  total: number
  page: number
  size: number
  totalPage: number
} 