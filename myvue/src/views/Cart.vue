<template>
  <div class="cart-page">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1>购物车</h1>
          <div class="nav-buttons">
            <el-button @click="$router.push('/')">
              首页
            </el-button>
            <el-button @click="$router.push('/products')">
              继续购物
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
        <div class="cart-content">
          <!-- 购物车商品列表 -->
          <div v-if="cart.length > 0" class="cart-items">
            <el-card class="cart-summary">
              <div class="summary-header">
                <h2>购物车商品 ({{ cart.length }}件)</h2>
                <el-button type="danger" @click="handleClearCart">
                  清空购物车
                </el-button>
              </div>
            </el-card>
            
            <el-card v-for="item in cart" :key="item._id" class="cart-item">
              <el-row :gutter="20" align="middle">
                <el-col :span="4">
                  <div class="item-image">
                    <img :src="item.image || '/placeholder.jpg'" :alt="item.product_name" />
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="item-info">
                    <h3 class="item-name">{{ item.product_name }}</h3>
                    <div class="item-rating">
                      <el-rate
                        :model-value="parseFloat(item.rating)"
                        disabled
                        show-score
                        text-color="#ff9900"
                        size="small"
                      />
                    </div>
                  </div>
                </el-col>
                <el-col :span="4">
                  <div class="item-price">
                    <span class="price">¥{{ item.price }}</span>
                  </div>
                </el-col>
                <el-col :span="4">
                  <div class="item-quantity">
                    <el-input-number
                      v-model="quantities[item._id]"
                      :min="1"
                      :max="10"
                      size="small"
                      @change="updateQuantity(item._id, $event)"
                    />
                  </div>
                </el-col>
                <el-col :span="4">
                  <div class="item-actions">
                    <el-button
                      type="danger"
                      size="small"
                      @click="handleRemoveFromCart(item._id)"
                    >
                      移除
                    </el-button>
                  </div>
                </el-col>
              </el-row>
            </el-card>
            
            <!-- 结算区域 -->
            <el-card class="checkout-section">
              <div class="checkout-content">
                <div class="total-info">
                  <div class="total-items">
                    <span>商品总数：{{ totalItems }}件</span>
                  </div>
                  <div class="total-price">
                    <span class="label">总计：</span>
                    <span class="price">¥{{ totalPrice.toFixed(2) }}</span>
                  </div>
                </div>
                <div class="checkout-actions">
                  <el-button size="large" @click="$router.push('/products')">
                    继续购物
                  </el-button>
                  <el-button type="primary" size="large" @click="handleCheckout">
                    立即结算
                  </el-button>
                </div>
              </div>
            </el-card>
          </div>
          
          <!-- 空购物车状态 -->
          <div v-else class="empty-cart">
            <el-empty description="购物车是空的">
              <el-button type="primary" @click="$router.push('/products')">
                去购物
              </el-button>
            </el-empty>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useProductStore } from '../stores/product'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const productStore = useProductStore()

const quantities = reactive<Record<string, number>>({})

const cart = computed(() => productStore.cart)

const totalItems = computed(() => {
  return Object.values(quantities).reduce((sum, qty) => sum + qty, 0)
})

const totalPrice = computed(() => {
  return cart.value.reduce((total, item) => {
    const quantity = quantities[item._id] || 1
    return total + (parseFloat(item.price) * quantity)
  }, 0)
})

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/')
}

const handleRemoveFromCart = async (productId: string) => {
  try {
    await ElMessageBox.confirm(
      '确定要从购物车中移除这个商品吗？',
      '确认移除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    productStore.removeFromCart(productId)
    delete quantities[productId]
    ElMessage.success('商品已移除')
  } catch {
    // 用户取消操作
  }
}

const handleClearCart = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清空购物车吗？',
      '确认清空',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    productStore.clearCart()
    Object.keys(quantities).forEach(key => {
      delete quantities[key]
    })
    ElMessage.success('购物车已清空')
  } catch {
    // 用户取消操作
  }
}

const updateQuantity = (productId: string, quantity: number) => {
  quantities[productId] = quantity
}

const handleCheckout = () => {
  ElMessage.success(`结算成功！总金额：¥${totalPrice.value.toFixed(2)}`)
  productStore.clearCart()
  Object.keys(quantities).forEach(key => {
    delete quantities[key]
  })
  router.push('/products')
}

const initQuantities = () => {
  cart.value.forEach(item => {
    if (!quantities[item._id]) {
      quantities[item._id] = 1
    }
  })
}

onMounted(() => {
  initQuantities()
})
</script>

<style scoped>
.cart-page {
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

.cart-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.cart-summary {
  margin-bottom: 20px;
}

.summary-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-header h2 {
  margin: 0;
  color: #333;
}

.cart-item {
  margin-bottom: 15px;
  border-radius: 8px;
}

.item-image {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-info {
  padding-left: 10px;
}

.item-name {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 10px 0;
  color: #333;
}

.item-rating {
  margin-bottom: 5px;
}

.item-price {
  text-align: center;
}

.item-price .price {
  font-size: 18px;
  font-weight: bold;
  color: #e74c3c;
}

.item-quantity {
  text-align: center;
}

.item-actions {
  text-align: center;
}

.checkout-section {
  margin-top: 30px;
  border-radius: 8px;
}

.checkout-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.total-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.total-items {
  font-size: 14px;
  color: #666;
}

.total-price {
  font-size: 20px;
  font-weight: bold;
}

.total-price .label {
  color: #333;
  margin-right: 10px;
}

.total-price .price {
  color: #e74c3c;
  font-size: 24px;
}

.checkout-actions {
  display: flex;
  gap: 15px;
}

.empty-cart {
  text-align: center;
  padding: 100px 20px;
}

:deep(.el-rate__text) {
  font-size: 12px;
}
</style>