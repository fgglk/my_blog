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

  // 批量删除图片 - 改进错误处理
  deleteImages: async (imageIds: number[]): Promise<ApiResponse> => {
    const results = []
    const successIds = []
    const failedIds = []
    
    // 逐个删除，避免 Promise.all 的快速失败问题
    for (const id of imageIds) {
      try {
        const result: any = await request.delete(`/image/delete/${id}`)
        // request工具已经通过响应拦截器处理，直接返回data部分
        if (result.code === 0) {
          successIds.push(id)
          results.push({ id, success: true })
        } else {
          failedIds.push(id)
          results.push({ id, success: false, error: result.msg })
        }
      } catch (error) {
        failedIds.push(id)
        results.push({ id, success: false, error: '网络错误' })
      }
    }
    
    // 返回详细的结果
    if (failedIds.length === 0) {
      return { code: 0, msg: `成功删除 ${successIds.length} 张图片`, data: { successIds, failedIds } }
    } else if (successIds.length === 0) {
      return { code: 1, msg: `删除失败，${failedIds.length} 张图片删除失败`, data: { successIds, failedIds } }
    } else {
      return { 
        code: 1, 
        msg: `部分删除成功，${successIds.length} 张删除成功，${failedIds.length} 张删除失败`, 
        data: { successIds, failedIds } 
      }
    }
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
