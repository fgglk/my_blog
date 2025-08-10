// 评论接口
export interface Comment {
  id: number
  content: string
  article_id: number
  user_id: number
  user_name: string
  user_avatar: string // 修改字段名以匹配后端响应
  parent_id?: number
  parent_user_name?: string
  children?: Comment[]
  created_at: string
  updated_at: string
}

// 创建评论请求参数
export interface CreateCommentRequest {
  content: string
  article_id: number
  parent_id?: number
}

// 评论列表响应
export interface CommentListResponse {
  list: Comment[]
  total: number
  page: number
  size: number
} 