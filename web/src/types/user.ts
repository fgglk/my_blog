// 用户信息接口
export interface UserInfo {
  id: number
  uuid: string
  username: string
  email: string
  nickname: string
  avatar: string
  bio?: string
  role: string
  status: number // 用户状态：0-禁用, 1-正常
  created_at: string
  updated_at: string
  last_login_at?: string
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
  keyword?: string
  status?: number
  sortBy?: string
  sortOrder?: string
}

// 用户列表响应
export interface UserListResponse {
  list: UserInfo[]
  total: number
  page: number
  pageSize: number
}

// 创建用户请求参数
export interface CreateUserRequest {
  username: string
  password: string
  nickname: string
  email: string
  role?: string
  bio?: string
  address?: string
} 