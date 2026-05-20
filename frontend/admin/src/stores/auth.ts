import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type LoginPhoneParams } from '@/api/auth'

const ACCESS_TOKEN_KEY = 'lx_access_token'
const REFRESH_TOKEN_KEY = 'lx_refresh_token'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string>(localStorage.getItem(ACCESS_TOKEN_KEY) || '')
  const refreshTokenValue = ref<string>(localStorage.getItem(REFRESH_TOKEN_KEY) || '')
  const user = ref<{
    id: number
    phone: string
    role: string
    nickname: string
    avatar_url: string
  } | null>(null)

  const isLoggedIn = computed(() => !!accessToken.value)

  function saveTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshTokenValue.value = refresh
    localStorage.setItem(ACCESS_TOKEN_KEY, access)
    localStorage.setItem(REFRESH_TOKEN_KEY, refresh)
  }

  async function login(params: LoginPhoneParams) {
    const res = await authApi.loginByPhone(params)
    saveTokens(res.access_token, res.refresh_token)
    user.value = res.user
  }

  async function refreshToken(): Promise<string> {
    const res = await authApi.refreshToken({ refresh_token: refreshTokenValue.value })
    saveTokens(res.access_token, res.refresh_token)
    return res.access_token
  }

  function logout() {
    accessToken.value = ''
    refreshTokenValue.value = ''
    user.value = null
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
  }

  return {
    accessToken,
    user,
    isLoggedIn,
    login,
    refreshToken,
    logout
  }
})
