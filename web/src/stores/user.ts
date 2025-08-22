import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { UserInfo } from '@/types/user'
import { userApi } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)
  const token = ref<string>('')

  // 从localStorage获取token
  const getTokenFromStorage = () => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      token.value = storedToken
    }
    return storedToken
  }

  // 初始化时自动恢复用户状态
  const initUserState = async () => {
    getTokenFromStorage() // Call to update token.value
    
    // 尝试从localStorage恢复用户信息
    const storedUserInfo = localStorage.getItem('userInfo')
    if (storedUserInfo) {
      try {
        const user = JSON.parse(storedUserInfo)
        userInfo.value = user
      } catch (error) {
        console.error('解析存储的用户信息失败:', error)
        localStorage.removeItem('userInfo')
      }
    }
    
    if (token.value) {
      await getUserInfo()
    }
  }

  // 设置用户信息
  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
    // 同时保存到localStorage，供路由守卫使用
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  // 设置token
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 清除用户信息
  const clearUserInfo = () => {
    userInfo.value = null
    token.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  // 获取当前token（实时从localStorage获取）
  const getCurrentToken = () => {
    return localStorage.getItem('token') || token.value
  }

  // 登录
  const login = async (username: string, password: string, captcha?: string, captchaId?: string) => {
    try {
      const loginData: any = { username, password }
      if (captcha && captchaId) {
        loginData.captcha_code = captcha
        loginData.captcha_id = captchaId
      }
      
      const response = await userApi.login(loginData)
      if (response.code === 0) {
        setToken(response.data.token)
        setUserInfo(response.data.user)
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '登录失败' }
    }
  }

  // 注册
  const register = async (userData: any) => {
    try {
      const response = await userApi.register(userData)
      if (response.code === 0) {
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '注册失败' }
    }
  }

  // 获取用户信息
  const getUserInfo = async () => {
    if (!token.value) {
      return
    }
    
    try {
      const response = await userApi.getUserInfo()
      
      if (response.code === 0) {
        setUserInfo(response.data)
      } else {
        // 如果获取用户信息失败，清除token和用户信息
        clearUserInfo()
      }
    } catch (error) {
      // 如果发生异常，清除token和用户信息
      clearUserInfo()
    }
  }

  // 登出
  const logout = () => {
    clearUserInfo()
  }

  // 更新用户信息
  const updateUserInfo = async (userData: any) => {
    try {
      const response = await userApi.updateUserInfo(userData)
      if (response.code === 0) {
        setUserInfo(response.data)
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '更新失败' }
    }
  }

  // 修改密码
  const changePassword = async (passwordData: any) => {
    try {
      const response = await userApi.changePassword(passwordData)
      if (response.code === 0) {
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '修改失败' }
    }
  }

  // 上传头像
  const uploadAvatar = async (file: File) => {
    try {
      const response = await userApi.uploadAvatar(file)
      if (response.code === 0 && response.data) {
        // 更新用户头像
        if (userInfo.value) {
          userInfo.value.avatar = response.data.url
        }
        return { success: true, data: response.data }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '上传失败' }
    }
  }

  return {
    userInfo,
    token,
    login,
    register,
    getUserInfo,
    logout,
    updateUserInfo,
    changePassword,
    uploadAvatar,
    setUserInfo,
    setToken,
    clearUserInfo,
    initUserState,
    getCurrentToken
  }
}) 