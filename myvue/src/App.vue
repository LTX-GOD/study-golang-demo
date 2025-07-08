<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from './stores/user'
import { useProductStore } from './stores/product'

const userStore = useUserStore()
const productStore = useProductStore()

onMounted(() => {
  // 初始化用户和购物车状态
  userStore.initUser()
  productStore.initCart()
})
</script>

<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style>
#app {
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #f5f5f5;
}

.el-container {
  min-height: 100vh;
}

.el-header {
  padding: 0;
}

.el-main {
  padding: 20px;
}
</style>
