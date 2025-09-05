import axios from 'axios'
import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: process.env.NODE_ENV === 'production' ? 'https://zjy456.cn/api' : '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

        // 请求拦截器
        request.interceptors.request.use(
          (config: InternalAxiosRequestConfig) => {
            // 添加token到请求头
            const token = localStorage.getItem('token')
            
            if (token && config.headers) {
              config.headers.Authorization = `Bearer ${token}`
            }
            return config
          },
          (error) => {
            return Promise.reject(error)
          }
        )

// 响应拦截器
request.interceptors.response.use(
          (response: AxiosResponse) => {
          const { data } = response
          
          // 如果响应码不是0，说明有错误，但不自动显示错误消息
          if (data.code !== 0) {
            // 如果是401未授权，跳转到登录页
            if (data.code === 401) {
              localStorage.removeItem('token')
              router.push('/login')
            }
            
            return Promise.reject(new Error(data.msg || '请求失败'))
          }
          
          return data
        },
          (error) => {
          let message = '网络错误'
          
          if (error.response) {
            const { status, data } = error.response
            
            switch (status) {
              case 400:
                message = data.msg || '请求参数错误'
                break
              case 401:
                message = '未授权，请重新登录'
                localStorage.removeItem('token')
                router.push('/login')
                break
              case 403:
                message = '拒绝访问'
                break
              case 404:
                message = '请求地址不存在'
                break
              case 500:
                message = '服务器内部错误'
                break
              default:
                message = data.msg || '请求失败'
            }
          } else if (error.request) {
            message = '网络连接失败'
          }
          
          ElMessage.error(message)
          return Promise.reject(error)
        }
)

export default request 