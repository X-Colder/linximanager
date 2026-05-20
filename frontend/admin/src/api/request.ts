import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig
} from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PaginatedData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 刷新 token 锁：防止并发刷新
let isRefreshing = false
let pendingRequests: Array<(token: string) => void> = []

const request: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' }
})

// 请求拦截：注入 Authorization
request.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const auth = useAuthStore()
  if (auth.accessToken) {
    config.headers.Authorization = `Bearer ${auth.accessToken}`
  }
  return config
})

// 响应拦截：统一处理错误 + Token 刷新
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { code, message, data } = response.data
    if (code !== 0) {
      ElMessage.error(message || '请求失败')
      return Promise.reject(new Error(message))
    }
    return data as unknown as AxiosResponse
  },
  async (error) => {
    const originalConfig: AxiosRequestConfig & { _retry?: boolean } =
      error.config

    // Access Token 过期，尝试刷新
    if (error.response?.status === 401 && !originalConfig._retry) {
      if (isRefreshing) {
        // 等待刷新完成后重试
        return new Promise((resolve) => {
          pendingRequests.push((token: string) => {
            originalConfig.headers = {
              ...originalConfig.headers,
              Authorization: `Bearer ${token}`
            }
            resolve(request(originalConfig))
          })
        })
      }

      originalConfig._retry = true
      isRefreshing = true

      try {
        const auth = useAuthStore()
        const newToken = await auth.refreshToken()
        pendingRequests.forEach((cb) => cb(newToken))
        pendingRequests = []
        return request(originalConfig)
      } catch {
        pendingRequests = []
        const auth = useAuthStore()
        auth.logout()
        return Promise.reject(error)
      } finally {
        isRefreshing = false
      }
    }

    // 网络错误或服务端错误提示
    if (!error.response) {
      ElMessage.error('网络连接失败，请检查网络')
    } else if (error.response.status >= 500) {
      ElMessage.error('服务器错误，请稍后重试')
    } else if (error.response.status === 403) {
      ElMessage.error('无权限访问')
    } else if (error.response.status === 404) {
      ElMessage.error('资源不存在')
    }

    return Promise.reject(error)
  }
)

export default request
