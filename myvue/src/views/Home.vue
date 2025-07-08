<template>
  <div class="home">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1>电商系统</h1>
          <div class="nav-buttons">
            <el-button v-if="!isLoggedIn" type="primary" @click="$router.push('/login')">
              登录
            </el-button>
            <el-button v-if="!isLoggedIn" @click="$router.push('/register')">
              注册
            </el-button>
            <div v-if="isLoggedIn" class="user-info">
              <span>欢迎，{{ user?.name }}</span>
              <el-button type="primary" @click="$router.push('/products')">
                商品列表
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
        </div>
      </el-header>
      
      <el-main>
        <div class="welcome-section">
          <h2>欢迎来到我们的电商平台</h2>
          <p>这里有最优质的商品和最贴心的服务</p>
          
          <div v-if="!isLoggedIn" class="action-buttons">
            <el-button size="large" type="primary" @click="$router.push('/register')">
              立即注册
            </el-button>
            <el-button size="large" @click="$router.push('/login')">
              已有账号？登录
            </el-button>
          </div>
          
          <div v-if="isLoggedIn" class="action-buttons">
            <el-button size="large" type="primary" @click="$router.push('/products')">
              开始购物
            </el-button>
          </div>
        </div>
        
        <div class="features">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-card class="feature-card">
                <h3>优质商品</h3>
                <p>精选优质商品，品质保证</p>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card class="feature-card">
                <h3>快速配送</h3>
                <p>全国范围内快速配送服务</p>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card class="feature-card">
                <h3>贴心服务</h3>
                <p>7x24小时客户服务支持</p>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { useProductStore } from '../stores/product'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const productStore = useProductStore()
const router = useRouter()

const isLoggedIn = computed(() => userStore.isLoggedIn)
const user = computed(() => userStore.user)
const cart = computed(() => productStore.cart)

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/')
}

onMounted(() => {
  userStore.initUser()
  productStore.initCart()
})
</script>

<style scoped>
.home {
  min-height: 100vh;
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
  align-items: center;
}

.user-info {
  display: flex;
  gap: 10px;
  align-items: center;
}

.user-info span {
  margin-right: 10px;
  font-weight: bold;
}

.welcome-section {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  margin-bottom: 40px;
}

.welcome-section h2 {
  font-size: 36px;
  margin-bottom: 20px;
  color: #333;
}

.welcome-section p {
  font-size: 18px;
  color: #666;
  margin-bottom: 30px;
}

.action-buttons {
  display: flex;
  gap: 20px;
  justify-content: center;
}

.features {
  padding: 0 20px 40px;
}

.feature-card {
  text-align: center;
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.feature-card h3 {
  color: #409eff;
  margin-bottom: 15px;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
}
</style>