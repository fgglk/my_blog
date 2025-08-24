import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'
import type { 
  ImageInfo, 
  ImageListRequest, 
  ImageListResponse, 
  DeleteImageRequest,
  UploadImageResponse
} from '@/types/image'

// 图片管理API
export const imageApi = {
  // 获取用户图片列表
  getImageList: (params: ImageListRequest): Promise<ApiResponse<ImageListResponse>> => {
    return request.get('/image/list', { params })
  },

  // 删除图片
  deleteImage: (data: DeleteImageRequest): Promise<ApiResponse> => {
    return request.delete(`/image/delete/${data.image_id}`)
  },

  // 批量删除图片 - 后端没有批量删除接口，需要逐个删除
  deleteImages: async (imageIds: number[]): Promise<ApiResponse> => {
    const promises = imageIds.map(id => 
      request.delete(`/image/delete/${id}`)
    )
    const results = await Promise.all(promises)
    // 检查是否所有删除都成功
    const failedCount = results.filter(r => r.data.code !== 0).length
    if (failedCount > 0) {
      return { code: 1, msg: `${failedCount} 张图片删除失败`, data: null }
    }
    return { code: 0, msg: '批量删除成功', data: null }
  },

  // 上传图片
  uploadImage: (file: File): Promise<ApiResponse<UploadImageResponse>> => {
    const formData = new FormData()
    formData.append('image', file)
    return request.post('/image/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取图片详情 - 后端没有单独的详情接口，使用show接口
  getImageDetail: (imageId: number): Promise<ApiResponse<ImageInfo>> => {
    return request.get(`/image/show/${imageId}`)
  }
}
