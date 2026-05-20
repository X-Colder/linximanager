// 基于 uni.request 封装，支持 JWT、错误处理、loading
const BASE_URL = 'https://api.linxi-shop.com/api/v1'

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: Record<string, unknown>
  loading?: boolean
  loadingText?: string
}

interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

function getToken(): string {
  return uni.getStorageSync('access_token') || ''
}

function request<T = unknown>(options: RequestOptions): Promise<T> {
  const { url, method = 'GET', data, loading = false, loadingText = '加载中...' } = options

  if (loading) {
    uni.showLoading({ title: loadingText, mask: true })
  }

  return new Promise((resolve, reject) => {
    const token = getToken()
    uni.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      success: (res) => {
        const body = res.data as ApiResponse<T>
        if (res.statusCode === 200 && body.code === 0) {
          resolve(body.data)
        } else if (res.statusCode === 401) {
          uni.showToast({ title: '登录已过期', icon: 'error' })
          setTimeout(() => {
            uni.reLaunch({ url: '/pages/login/index' })
          }, 1500)
          reject(new Error('Unauthorized'))
        } else {
          uni.showToast({ title: body.message || '请求失败', icon: 'error' })
          reject(new Error(body.message))
        }
      },
      fail: () => {
        uni.showToast({ title: '网络连接失败', icon: 'error' })
        reject(new Error('Network Error'))
      },
      complete: () => {
        if (loading) uni.hideLoading()
      }
    })
  })
}

export const http = {
  get: <T = unknown>(url: string, data?: Record<string, unknown>, loading?: boolean) =>
    request<T>({ url, method: 'GET', data, loading }),

  post: <T = unknown>(url: string, data?: Record<string, unknown>, loading?: boolean) =>
    request<T>({ url, method: 'POST', data, loading, loadingText: '提交中...' }),

  put: <T = unknown>(url: string, data?: Record<string, unknown>, loading?: boolean) =>
    request<T>({ url, method: 'PUT', data, loading }),

  delete: <T = unknown>(url: string, data?: Record<string, unknown>) =>
    request<T>({ url, method: 'DELETE', data })
}

export default http
