<template>
  <div class="admin-page">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1>管理面板</h1>
          <div class="nav-buttons">
            <el-button @click="$router.push('/')">
              首页
            </el-button>
            <el-button @click="$router.push('/products')">
              商品列表
            </el-button>
            <el-button @click="$router.push('/cart')">
              购物车
            </el-button>
            <el-button type="danger" @click="handleLogout">
              退出登录
            </el-button>
          </div>
        </div>
      </el-header>
      
      <el-main>
        <div class="admin-content">
          <!-- 添加商品表单 -->
          <el-card class="add-product-card">
            <template #header>
              <div class="card-header">
                <h2>添加新商品</h2>
              </div>
            </template>
            
            <el-form
              ref="productFormRef"
              :model="productForm"
              :rules="rules"
              label-width="100px"
              @submit.prevent="handleAddProduct"
            >
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="商品名称" prop="product_name">
                    <el-input
                      v-model="productForm.product_name"
                      placeholder="请输入商品名称"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="商品价格" prop="price">
                    <el-input
                      v-model="productForm.price"
                      placeholder="请输入商品价格"
                      type="number"
                      step="0.01"
                    >
                      <template #prepend>¥</template>
                    </el-input>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="商品评分" prop="rating">
                    <el-rate
                      v-model="ratingValue"
                      @change="updateRating"
                      show-score
                      text-color="#ff9900"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="商品图片" prop="image">
                    <el-input
                      v-model="productForm.image"
                      placeholder="请输入图片URL"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-form-item>
                <div class="form-actions">
                  <el-button @click="resetForm">
                    重置
                  </el-button>
                  <el-button
                    type="primary"
                    :loading="loading"
                    @click="handleAddProduct"
                  >
                    添加商品
                  </el-button>
                </div>
              </el-form-item>
            </el-form>
          </el-card>
          
          <!-- 商品预览 -->
          <el-card v-if="productForm.product_name" class="preview-card">
            <template #header>
              <div class="card-header">
                <h2>商品预览</h2>
              </div>
            </template>
            
            <div class="product-preview">
              <el-row :gutter="20">
                <el-col :span="8">
                  <div class="preview-image">
                    <img 
                      :src="productForm.image || '/placeholder.jpg'" 
                      :alt="productForm.product_name"
                      @error="handleImageError"
                    />
                  </div>
                </el-col>
                <el-col :span="16">
                  <div class="preview-info">
                    <h3 class="preview-name">{{ productForm.product_name }}</h3>
                    <div class="preview-price">¥{{ productForm.price }}</div>
                    <div class="preview-rating">
                      <el-rate
                        v-model="ratingValue"
                        disabled
                        show-score
                        text-color="#ff9900"
                      />
                    </div>
                    <div class="preview-actions">
                      <el-button type="primary" disabled>
                        加入购物车
                      </el-button>
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </el-card>
          
          <!-- 统计信息 -->
          <el-card class="stats-card">
            <template #header>
              <div class="card-header">
                <h2>统计信息</h2>
              </div>
            </template>
            
            <el-row :gutter="20">
              <el-col :span="8">
                <div class="stat-item">
                  <div class="stat-value">{{ products.length }}</div>
                  <div class="stat-label">商品总数</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="stat-item">
                  <div class="stat-value">{{ cart.length }}</div>
                  <div class="stat-label">购物车商品</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="stat-item">
                  <div class="stat-value">{{ user?.name || 'N/A' }}</div>
                  <div class="stat-label">当前用户</div>
                </div>
              </el-col>
            </el-row>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useProductStore } from '../stores/product'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const productStore = useProductStore()
const productFormRef = ref<FormInstance>()
const loading = ref(false)
const ratingValue = ref(5)

const productForm = reactive({
  product_name: '',
  price: '',
  rating: '5',
  image: ''
})

const products = computed(() => productStore.products)
const cart = computed(() => productStore.cart)
const user = computed(() => userStore.user)

const rules: FormRules = {
  product_name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '商品名称长度应在2-50位之间', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入商品价格', trigger: 'blur' },
    { pattern: /^\d+(\.\d{1,2})?$/, message: '请输入正确的价格格式', trigger: 'blur' }
  ],
  rating: [
    { required: true, message: '请选择商品评分', trigger: 'change' }
  ],
  image: [
    { required: true, message: '请输入图片URL', trigger: 'blur' },
    { type: 'url', message: '请输入正确的URL格式', trigger: 'blur' }
  ]
}

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/')
}

const updateRating = (value: number) => {
  productForm.rating = value.toString()
}

const handleAddProduct = async () => {
  if (!productFormRef.value) return
  
  await productFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const result = await productStore.addProduct(productForm)
        if (result.success) {
          ElMessage.success('商品添加成功')
          resetForm()
        } else {
          ElMessage.error(result.error)
        }
      } catch (error) {
        ElMessage.error('添加商品失败，请重试')
      } finally {
        loading.value = false
      }
    }
  })
}

const resetForm = () => {
  if (productFormRef.value) {
    productFormRef.value.resetFields()
  }
  Object.assign(productForm, {
    product_name: '',
    price: '',
    rating: '5',
    image: ''
  })
  ratingValue.value = 5
}

const handleImageError = (event: Event) => {
  const target = event.target as HTMLImageElement
  target.src = '/placeholder.jpg'
}

onMounted(() => {
  productStore.fetchProducts()
})
</script>

<style scoped>
.admin-page {
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

.admin-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.add-product-card,
.preview-card,
.stats-card {
  margin-bottom: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0;
  color: #333;
  font-weight: 600;
}

.form-actions {
  text-align: center;
  margin-top: 20px;
}

.product-preview {
  padding: 20px;
}

.preview-image {
  width: 100%;
  height: 200px;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
  border: 2px dashed #ddd;
}

.preview-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-info {
  padding-left: 20px;
}

.preview-name {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 15px 0;
  color: #333;
}

.preview-price {
  font-size: 24px;
  font-weight: bold;
  color: #e74c3c;
  margin-bottom: 15px;
}

.preview-rating {
  margin-bottom: 20px;
}

.preview-actions {
  margin-top: 20px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

:deep(.el-card__header) {
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
}
</style>