// 用户信息接口
export interface UserInfo {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  bio?: string
  role: string
  createdAt: string
  updatedAt: string
}

// 登录请求参数
export interface LoginRequest {
  username: string
  password: string
  captcha?: string
}

// 注册请求参数
export interface RegisterRequest {
  username: string
  password: string
  email: string
  nickname: string
  captcha?: string
  emailCode?: string
}

// 更新用户信息请求参数
export interface UpdateUserRequest {
  nickname?: string
  email?: string
  avatar?: string
  bio?: string
}

// 修改密码请求参数
export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

// 用户列表请求参数
export interface UserListRequest {
  page: number
  size: number
}

// 用户列表响应
export interface UserListResponse {
  list: UserInfo[]
  total: number
  page: number
  pageSize: number
} 