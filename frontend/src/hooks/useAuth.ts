import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { USER_ROLE } from '@/constants'

export function useAuth() {
  const authStore = useAuthStore()
  const isLogin = computed(() => authStore.token !== '')
  const user = computed(() => authStore.user)
  const isAdmin = computed(() => authStore.user?.role === USER_ROLE.ADMIN)
  const isStaffOrAdmin = computed(() => authStore.user?.role === USER_ROLE.ADMIN || authStore.user?.role === USER_ROLE.STAFF)
  return { isLogin, user, isAdmin, isStaffOrAdmin, login: authStore.login, logout: authStore.logout, fetchProfile: authStore.fetchProfile }
}
