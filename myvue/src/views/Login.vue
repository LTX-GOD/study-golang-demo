<template>
  <div class="login-page">
    <el-container>
      <el-main>
        <div class="login-container">
          <el-card class="login-card">
            <template #header>
              <div class="card-header">
                <h2>用户登录</h2>
              </div>
            </template>
            
            <el-form
              ref="loginFormRef"
              :model="loginForm"
              :rules="rules"
              label-width="80px"
              @submit.prevent="handleLogin"
            >
              <el-form-item label="邮箱" prop="email">
                <el-input
                  v-model="loginForm.email"
                  type="email"
                  placeholder="请输入邮箱"
                  prefix-icon="Message"
                />
              </el-form-item>
              
              <el-form-item label="密码" prop="password">
                <el-input
                  v-model="loginForm.password"
                  type="password"
                  placeholder="请输入密码"
                  prefix-icon="Lock"
                  show-password
                />
              </el-form-item>
              
              <el-form-item>
                <el-button
                  type="primary"
                  :loading="loading"
                  @click="handleLogin"
                  style="width: 100%"
                >
                  登录
                </el-button>
              </el-form-item>
              
              <el-form-item>
                <div class="form-footer">
                  <span>还没有账号？</span>
                  <el-button type="text" @click="$router.push('/register')">
                    立即注册
                  </el-button>
                </div>
              </el-form-item>
            </el-form>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const loginFormRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  email: '',
  password: ''
})

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 30, message: '密码长度应在6-30位之间', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const result = await userStore.login(loginForm)
        if (result.success) {
          ElMessage.success('登录成功')
          router.push('/products')
        } else {
          ElMessage.error(result.error)
        }
      } catch (error) {
        ElMessage.error('登录失败，请重试')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.login-card {
  width: 400px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0;
  color: #333;
  font-weight: 600;
}

.form-footer {
  text-align: center;
  width: 100%;
}

.form-footer span {
  color: #666;
  margin-right: 8px;
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