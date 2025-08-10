import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

// 标签接口
export interface Tag {
  id: number
  name: string
  slug: string
  article_count: number
}

// 标签API
export const tagApi = {
  // 获取标签列表
  getTagList: (): Promise<ApiResponse<Tag[]>> => {
    return request.get('/tags')
  },

  // 获取标签详情
  getTag: (id: number): Promise<ApiResponse<Tag>> => {
    return request.get(`/tags/${id}`)
  }
} 