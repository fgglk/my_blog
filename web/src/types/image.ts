// 图片信息接口
export interface ImageInfo {
  id: number
  uuid: string
  user_id: number
  filename: string
  original_name: string
  url: string
  size: number
  mime_type: string
  width?: number
  height?: number
  created_at: string
  updated_at: string
}

// 图片列表请求参数
export interface ImageListRequest {
  page: number
  size: number
  keyword?: string
  sortBy?: string
  sortOrder?: string
}

// 图片列表响应
export interface ImageListResponse {
  list: ImageInfo[]
  total: number
  page: number
  pageSize: number
}

// 删除图片请求
export interface DeleteImageRequest {
  image_id: number
}

// 上传图片响应
export interface UploadImageResponse {
  id: number
  url: string
  filename: string
  size: number
}
