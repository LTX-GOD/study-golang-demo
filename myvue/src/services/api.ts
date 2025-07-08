import axios from 'axios'
import type { AxiosResponse } from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8000', // 后端服务地址
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.token = token
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 用户相关API
export const userAPI = {
  // 用户注册
  register: (userData: {
    name: string
    email: string
    password: string
    phone: string
  }) => api.post('/user/signup', userData),

  // 用户登录
  login: (credentials: {
    email: string
    password: string
  }) => api.post('/user/login', credentials)
}

// 商品相关API
export const productAPI = {
  // 获取所有商品
  getAllProducts: () => api.get('/users/productview'),

  // 根据名称搜索商品
  searchProducts: (name: string) => api.get(`/users/search?name=${name}`),

  // 添加商品（管理员）
  addProduct: (productData: {
    product_name: string
    price: string
    rating: string
    image: string
  }) => api.post('/admin/addproduct', productData)
}

// 购物车相关API
export const cartAPI = {
  // 添加商品到购物车
  addToCart: (productId: string, userId: string) => 
    api.get(`/addtocart?id=${productId}&userID=${userId}`),

  // 获取购物车商品
  getCartItems: (userId: string) => 
    api.get(`/listcart?id=${userId}`),

  // 从购物车移除商品
  removeFromCart: (productId: string, userId: string) => 
    api.get(`/removeitem?id=${productId}&userID=${userId}`),

  // 清空购物车（购买）
  buyFromCart: (userId: string) => 
    api.get(`/cartcheckout?userID=${userId}`)
}

export default api