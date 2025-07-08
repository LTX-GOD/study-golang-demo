import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userAPI } from '../services/api'
import type { AxiosResponse } from 'axios'

interface User {
  _id: string
  name: string
  email: string
  phone: string
  user_id: string
  usercart: any[]
  address: any[]
  order: any[]
}

interface LoginResponse {
  user: User
  token: string
  refreshToken: string
}

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const isLoggedIn = computed(() => !!token.value)

  // 登录
  const login = async (credentials: { email: string; password: string }) => {
    try {
      const response: AxiosResponse<LoginResponse> = await userAPI.login(credentials)
      const { user: userData, token: userToken, refreshToken } = response.data
      
      user.value = userData
      token.value = userToken
      
      localStorage.setItem('token', userToken)
      localStorage.setItem('refreshToken', refreshToken)
      localStorage.setItem('user', JSON.stringify(userData))
      
      return { success: true, data: response.data }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '登录失败' 
      }
    }
  }

  // 注册
  const register = async (userData: {
    name: string
    email: string
    password: string
    phone: string
  }) => {
    try {
      const response = await userAPI.register(userData)
      return { success: true, data: response.data }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '注册失败' 
      }
    }
  }

  // 登出
  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

  // 初始化用户信息
  const initUser = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      user.value = JSON.parse(storedUser)
    }
  }

  return {
    user,
    token,
    isLoggedIn,
    login,
    register,
    logout,
    initUser
  }
})