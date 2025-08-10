import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

// 分类接口
export interface Category {
  id: number
  name: string
  slug: string
  article_count: number
}

// 分类API
export const categoryApi = {
  // 获取分类列表
  getCategoryList: (): Promise<ApiResponse<Category[]>> => {
    return request.get('/categories')
  },

  // 获取分类详情
  getCategory: (id: number): Promise<ApiResponse<Category>> => {
    return request.get(`/categories/${id}`)
  }
} 