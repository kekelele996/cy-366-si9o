<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="brand">🎮 电竞馆上机管理系统</h1>
      <van-tabs v-model:active="tab">
        <van-tab title="登录" name="login">
          <van-form @submit="onLogin">
            <van-cell-group inset>
              <van-field v-model="loginForm.username" name="username" label="用户名" placeholder="请输入用户名" :rules="[{ required: true, message: '请输入用户名' }]" />
              <van-field v-model="loginForm.password" type="password" name="password" label="密码" placeholder="请输入密码" :rules="[{ required: true, message: '请输入密码' }]" />
            </van-cell-group>
            <div class="submit-btn"><van-button round block type="primary" native-type="submit" :loading="loading">登 录</van-button></div>
          </van-form>
        </van-tab>
        <van-tab title="注册" name="register">
          <van-form @submit="onRegister">
            <van-cell-group inset>
              <van-field v-model="regForm.username" name="username" label="用户名" placeholder="3-32位" :rules="[{ required: true, message: '请输入用户名' }]" />
              <van-field v-model="regForm.nickname" name="nickname" label="昵称" placeholder="请输入昵称" :rules="[{ required: true, message: '请输入昵称' }]" />
              <van-field v-model="regForm.phone" name="phone" label="手机号" placeholder="选填" />
              <van-field v-model="regForm.password" type="password" name="password" label="密码" placeholder="至少6位" :rules="[{ required: true, message: '请输入密码' }]" />
            </van-cell-group>
            <div class="submit-btn"><van-button round block type="primary" native-type="submit" :loading="loading">注 册</van-button></div>
          </van-form>
        </van-tab>
      </van-tabs>
      <p class="tips">演示账号：admin / admin123456（管理员），member / member123456（会员）</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showSuccessToast } from 'vant'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const tab = ref('login')
const loading = ref(false)
const loginForm = reactive({ username: '', password: '' })
const regForm = reactive({ username: '', nickname: '', phone: '', password: '' })

async function onLogin() {
  loading.value = true
  try {
    await authStore.login(loginForm.username, loginForm.password)
    showSuccessToast('登录成功')
    router.push('/dashboard')
  } finally {
    loading.value = false
  }
}

async function onRegister() {
  loading.value = true
  try {
    await authStore.register({ username: regForm.username, nickname: regForm.nickname, phone: regForm.phone, password: regForm.password })
    showSuccessToast('注册成功')
    router.push('/dashboard')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page { min-height: 100vh; background: linear-gradient(135deg, #1989fa 0%, #07c160 100%); display: flex; align-items: center; justify-content: center; padding: 20px; }
.login-card { width: 100%; max-width: 420px; background: #fff; border-radius: 16px; padding: 24px 12px 16px; box-shadow: 0 8px 24px rgba(0,0,0,0.15); }
.brand { text-align: center; font-size: 20px; margin: 0 0 20px; }
.submit-btn { margin: 16px 16px 0; }
.tips { text-align: center; color: #969799; font-size: 12px; margin-top: 16px; }
</style>
