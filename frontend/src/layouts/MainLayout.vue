<template>
  <div class="main-layout">
    <van-nav-bar :title="currentTitle" left-arrow @click-left="onBack" @click-right="onLogout">
      <template #right>
        <van-icon name="logout" size="20" />
      </template>
    </van-nav-bar>
    <div class="layout-body">
      <router-view />
    </div>
    <van-tabbar route active-color="#1989fa">
      <van-tabbar-item replace to="/dashboard" icon="wap-home">看板</van-tabbar-item>
      <van-tabbar-item replace to="/stations" icon="apps-o">机位</van-tabbar-item>
      <van-tabbar-item replace to="/recharge" icon="gold-coin-o">充值</van-tabbar-item>
      <van-tabbar-item replace to="/reservations" icon="calendar-o">预约</van-tabbar-item>
      <van-tabbar-item replace to="/tournaments" icon="trophy-o">赛事</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showConfirmDialog } from 'vant'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const titles: Record<string, string> = {
  dashboard: '运营看板',
  stations: '机位管理',
  recharge: '充值与时长包',
  reservations: '机位预约',
  sessions: '上机记录与排行榜',
  tournaments: '赛事管理',
  audits: '操作审计',
}

const currentTitle = computed(() => titles[route.name as string] || '电竞馆')

function onBack() {
  if (route.path !== '/dashboard') {
    router.push('/dashboard')
  }
}

async function onLogout() {
  try {
    await showConfirmDialog({ title: '提示', message: '确定退出登录吗？' })
    authStore.logout()
    router.push('/login')
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  padding-bottom: 50px;
}
.layout-body {
  padding: 12px 12px 0;
}
</style>
