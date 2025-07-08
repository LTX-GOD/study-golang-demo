<template>
  <div class="products-page">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1>商品列表</h1>
          <div class="nav-buttons">
            <el-button @click="$router.push('/')">
              首页
            </el-button>
            <el-button @click="$router.push('/cart')">
              购物车 ({{ cart.length }})
            </el-button>
            <el-button @click="$router.push('/admin')">
              管理面板
            </el-button>
            <el-button type="danger" @click="handleLogout">
              退出登录
            </el-button>
          </div>
        </div>
      </el-header>
      
      <el-main>
        <!-- 搜索区域 -->
        <div class="search-section">
          <el-row :gutter="20">
            <el-col :span="16">
              <el-input
                v-model="searchQuery"
                placeholder="搜索商品名称"
                prefix-icon="Search"
                @keyup.enter="handleSearch"
              />
            </el-col>
            <el-col :span="4">
              <el-button type="primary" @click="handleSearch" :loading="loading">
                搜索
              </el-button>
            </el-col>
            <el-col :span="4">
              <el-button @click="handleRefresh" :loading="loading">
                刷新
              </el-button>
            </el-col>
          </el-row>
        </div>
        
        <!-- 商品列表 -->
        <div class="products-section">
          <el-row :gutter="20" v-loading="loading">
            <el-col :span="6" v-for="product in products" :key="product._id">
              <el-card class="product-card" shadow="hover">
                <div class="product-image">
                  <img :src="product.image || '/placeholder.jpg'" :alt="product.product_name" />
                </div>
                <div class="product-info">
                  <h3 class="product-name">{{ product.product_name }}</h3>
                  <div class="product-price">¥{{ product.price }}</div>
                  <div class="product-rating">
                    <el-rate
                      :model-value="parseFloat(product.rating)"
                      disabled
                      show-score
                      text-color="#ff9900"
                    />
                  </div>
                  <div class="product-actions">
                    <el-button
                      type="primary"
                      @click="handleAddToCart(product)"
                      :disabled="isInCart(product._id)"
                    >
                      {{ isInCart(product._id) ? '已在购物车' : '加入购物车' }}
                    </el-button>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>
          
          <!-- 空状态 -->
          <div v-if="!loading && products.length === 0" class="empty-state">
            <el-empty description="暂无商品数据">
              <el-button type="primary" @click="handleRefresh">
                重新加载
              </el-button>
            </el-empty>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useProductStore } from '../stores/product'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const productStore = useProductStore()

const searchQuery = ref('')

const products = computed(() => productStore.products)
const cart = computed(() => productStore.cart)
const loading = computed(() => productStore.loading)

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/')
}

const handleSearch = async () => {
  if (searchQuery.value.trim()) {
    const result = await productStore.searchProducts(searchQuery.value.trim())
    if (!result.success) {
      ElMessage.error(result.error)
    }
  } else {
    handleRefresh()
  }
}

const handleRefresh = async () => {
  searchQuery.value = ''
  const result = await productStore.fetchProducts()
  if (!result.success) {
    ElMessage.error(result.error)
  }
}

const handleAddToCart = async (product: any) => {
  const result = await productStore.addToCart(product)
  if (result.success) {
    ElMessage.success(result.message)
  } else {
    ElMessage.warning(result.message)
  }
}

const isInCart = (productId: string) => {
  return cart.value.some(item => item._id === productId)
}

onMounted(() => {
  handleRefresh()
})
</script>

<style scoped>
.products-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 20px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
}

.nav-buttons {
  display: flex;
  gap: 10px;
}

.search-section {
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.products-section {
  padding: 0 20px;
}

.product-card {
  margin-bottom: 20px;
  border-radius: 12px;
  overflow: hidden;
  transition: transform 0.3s ease;
}

.product-card:hover {
  transform: translateY(-5px);
}

.product-image {
  height: 200px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-info {
  padding: 15px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 10px 0;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-price {
  font-size: 20px;
  font-weight: bold;
  color: #e74c3c;
  margin-bottom: 10px;
}

.product-rating {
  margin-bottom: 15px;
}

.product-actions {
  text-align: center;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

:deep(.el-rate__text) {
  font-size: 12px;
}

:deep(.el-card__body) {
  padding: 0;
}
</style>