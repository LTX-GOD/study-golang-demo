import { defineStore } from 'pinia'
import { ref } from 'vue'
import { productAPI, cartAPI } from '../services/api'
import type { AxiosResponse } from 'axios'

interface Product {
  _id: string
  product_name: string
  price: string
  rating: string
  image: string
}

export const useProductStore = defineStore('product', () => {
  const products = ref<Product[]>([])
  const loading = ref(false)
  const cart = ref<Product[]>([])

  // 获取所有商品
  const fetchProducts = async () => {
    loading.value = true
    try {
      const response: AxiosResponse<Product[]> = await productAPI.getAllProducts()
      products.value = response.data
      return { success: true, data: response.data }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '获取商品失败' 
      }
    } finally {
      loading.value = false
    }
  }

  // 搜索商品
  const searchProducts = async (name: string) => {
    loading.value = true
    try {
      const response: AxiosResponse<Product[]> = await productAPI.searchProducts(name)
      products.value = response.data
      return { success: true, data: response.data }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '搜索商品失败' 
      }
    } finally {
      loading.value = false
    }
  }

  // 添加商品（管理员）
  const addProduct = async (productData: {
    product_name: string
    price: string
    rating: string
    image: string
  }) => {
    try {
      const response = await productAPI.addProduct(productData)
      await fetchProducts() // 重新获取商品列表
      return { success: true, data: response.data }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.error || '添加商品失败' 
      }
    }
  }

  // 添加到购物车
  const addToCart = async (product: Product) => {
    try {
      // 获取当前用户信息
      const userStr = localStorage.getItem('user')
      if (!userStr) {
        return { success: false, message: '请先登录' }
      }
      
      const user = JSON.parse(userStr)
      const userId = user._id || user.user_id
      
      if (!userId) {
        return { success: false, message: '用户信息无效，请重新登录' }
      }

      // 检查商品是否已在本地购物车中
      const existingItem = cart.value.find(item => item._id === product._id)
      if (existingItem) {
        return { success: false, message: '商品已在购物车中' }
      }

      // 调用后端API添加到购物车
      await cartAPI.addToCart(product._id, userId)
      
      // 添加到本地购物车
      cart.value.push(product)
      localStorage.setItem('cart', JSON.stringify(cart.value))
      
      return { success: true, message: '已添加到购物车' }
    } catch (error: any) {
      console.error('添加到购物车失败:', error)
      return { 
        success: false, 
        message: error.response?.data?.error || '添加到购物车失败' 
      }
    }
  }

  // 从购物车移除
  const removeFromCart = (productId: string) => {
    cart.value = cart.value.filter(item => item._id !== productId)
    localStorage.setItem('cart', JSON.stringify(cart.value))
  }

  // 清空购物车
  const clearCart = () => {
    cart.value = []
    localStorage.removeItem('cart')
  }

  // 初始化购物车
  const initCart = () => {
    const storedCart = localStorage.getItem('cart')
    if (storedCart) {
      cart.value = JSON.parse(storedCart)
    }
  }

  // 计算购物车总价
  const cartTotal = () => {
    return cart.value.reduce((total, item) => {
      return total + parseFloat(item.price)
    }, 0)
  }

  return {
    products,
    loading,
    cart,
    fetchProducts,
    searchProducts,
    addProduct,
    addToCart,
    removeFromCart,
    clearCart,
    initCart,
    cartTotal
  }
})