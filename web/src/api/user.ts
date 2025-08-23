import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'
import type { 
  UserInfo, 
  LoginRequest, 
  RegisterRequest, 
  UpdateUserRequest, 
  ChangePasswordRequest,
  UserListRequest,
  UserListResponse
} from '@/types/user'

// 用户API
export const userApi = {
  // 用户注册
  register: (data: RegisterRequest): Promise<ApiResponse> => {
    return request.post('/users/register', data)
  },

  // 用户登录
  login: (data: LoginRequest): Promise<ApiResponse<{ token: string; user: UserInfo }>> => {
    return request.post('/users/login', data)
  },

  // 获取用户信息
  getUserInfo: (): Promise<ApiResponse<UserInfo>> => {
    return request.get('/users/info')
  },

  // 根据ID获取用户信息
  getUserById: (userId: number): Promise<ApiResponse<UserInfo>> => {
    return request.get(`/users/${userId}`)
  },

  // 更新用户信息
  updateUserInfo: (data: UpdateUserRequest): Promise<ApiResponse> => {
    return request.put('/users/update', data)
  },

  // 修改密码
  changePassword: (data: ChangePasswordRequest): Promise<ApiResponse> => {
    return request.put('/users/password', data)
  },

  // 删除用户
  deleteUser: (): Promise<ApiResponse> => {
    return request.delete('/users/delete')
  },

  // 获取用户列表
  getUserList: (params: UserListRequest): Promise<ApiResponse<UserListResponse>> => {
    return request.get('/users/list', { params })
  },

  // 获取验证码
  getCaptcha: (): Promise<ApiResponse<{ image: string; captcha_id: string }>> => {
    return request.get('/users/captcha')
  },

  // 发送邮箱验证码
  sendEmailCode: (email: string): Promise<ApiResponse> => {
    return request.get('/users/email/code', { params: { email } })
  },

  // 忘记密码
  forgotPassword: (data: { email: string; emailCode: string }): Promise<ApiResponse> => {
    return request.post('/users/forgot', data)
  },

  // 重置密码
  resetPassword: (data: { email: string; emailCode: string; newPassword: string }): Promise<ApiResponse> => {
    return request.post('/users/reset', data)
  },

  // 上传头像
  uploadAvatar: (file: File): Promise<ApiResponse<{ url: string }>> => {
    const formData = new FormData()
    formData.append('avatar', file) // 改回'avatar'，使用专门的头像上传API
    return request.post('/image/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 审核通过用户
  approveUser: (userId: string): Promise<ApiResponse> => {
    return request.put(`/users/${userId}/approve`)
  },

  // 拒绝用户
  rejectUser: (userId: string): Promise<ApiResponse> => {
    return request.put(`/users/${userId}/reject`)
  },

  // 删除用户（管理员）
  deleteUserById: (userId: string): Promise<ApiResponse> => {
    return request.delete(`/users/${userId}`)
  },

  // 更新用户状态
  updateUserStatus: (userId: string, status: string): Promise<ApiResponse> => {
    return request.put(`/users/${userId}/status`, { status })
  }
} 