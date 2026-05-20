import request from './request'

export interface LoginPhoneParams {
  phone: string
  code: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: {
    id: number
    phone: string
    role: string
    nickname: string
    avatar_url: string
  }
}

export interface RefreshTokenParams {
  refresh_token: string
}

export const authApi = {
  loginByPhone: (params: LoginPhoneParams) =>
    request.post<unknown, LoginResponse>('/auth/login/phone', params),

  sendCode: (phone: string) =>
    request.post('/auth/sms/code', { phone }),

  refreshToken: (params: RefreshTokenParams) =>
    request.post<unknown, LoginResponse>('/auth/refresh', params),

  logout: () =>
    request.post('/auth/logout')
}
