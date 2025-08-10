import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'
import type {CreateCommentRequest, CommentListResponse } from '@/types/comment'

// 评论API
export const commentApi = {
  // 获取文章评论列表
  getCommentList: (articleId: number, params: { page: number; size: number }): Promise<ApiResponse<CommentListResponse>> => {
    return request.get('/comments', { 
      params: {
        ...params,
        article_id: articleId
      }
    })
  },

            // 创建评论
          createComment: (data: CreateCommentRequest): Promise<ApiResponse> => {
            return request.post('/comments', {
              content: data.content,
              article_id: data.article_id,
              parent_id: data.parent_id
            })
          },

  // 删除评论
  deleteComment: (id: number): Promise<ApiResponse> => {
    return request.delete(`/comments/${id}`)
  }
} 