import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'Login', component: () => import('@/pages/Login.vue') },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/pages/Dashboard.vue') },
      { path: 'stations', name: 'Stations', component: () => import('@/pages/Stations.vue') },
      { path: 'recharge', name: 'Recharge', component: () => import('@/pages/Recharge.vue') },
      { path: 'reservations', name: 'Reservations', component: () => import('@/pages/Reservations.vue') },
      { path: 'sessions', name: 'Sessions', component: () => import('@/pages/Sessions.vue') },
      { path: 'tournaments', name: 'Tournaments', component: () => import('@/pages/Tournaments.vue') },
      { path: 'audits', name: 'Audits', component: () => import('@/pages/Audits.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (to.name !== 'Login' && !token) {
    return { name: 'Login' }
  }
  if (to.name === 'Login' && token) {
    return { name: 'Dashboard' }
  }
  return true
})

export default router
