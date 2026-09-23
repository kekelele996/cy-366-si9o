import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, register as registerApi, fetchProfile, type UserItem } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<UserItem | null>(JSON.parse(localStorage.getItem('user') || 'null'))

  function setSession(t: string, u: UserItem) {
    token.value = t
    user.value = u
    localStorage.setItem('token', t)
    localStorage.setItem('user', JSON.stringify(u))
  }

  async function login(username: string, password: string) {
    const res = await loginApi({ username, password })
    setSession(res.token, res.user)
    return res.user
  }

  async function register(data: { username: string; password: string; nickname: string; phone?: string }) {
    const res = await registerApi(data)
    setSession(res.token, res.user)
    return res.user
  }

  async function fetchProfile() {
    if (!token.value) return null
    try {
      const u = await fetchProfile()
      user.value = u
      localStorage.setItem('user', JSON.stringify(u))
      return u
    } catch {
      return null
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return { token, user, login, register, fetchProfile, logout }
})
