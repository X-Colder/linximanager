// C端 API 请求封装
const BASE_URL = 'https://api.linxi-shop.com/api/v1'

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: Record<string, unknown>
  loading?: boolean
}

interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

function getToken(): string {
  return uni.getStorageSync('consumer_token') || ''
}

function request<T = unknown>(options: RequestOptions): Promise<T> {
  const { url, method = 'GET', data, loading = false } = options

  if (loading) {
    uni.showLoading({ title: '加载中...', mask: true })
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        Authorization: getToken() ? `Bearer ${getToken()}` : ''
      },
      success: (res) => {
        const body = res.data as ApiResponse<T>
        if (res.statusCode === 200 && body.code === 0) {
          resolve(body.data)
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
    request<T>({ url, method: 'POST', data, loading }),

  put: <T = unknown>(url: string, data?: Record<string, unknown>) =>
    request<T>({ url, method: 'PUT', data }),

  delete: <T = unknown>(url: string) =>
    request<T>({ url, method: 'DELETE' })
}

export default http
