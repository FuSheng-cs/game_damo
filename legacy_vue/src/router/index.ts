import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import StartView from '@/views/StartView.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Start',
    component: StartView
  },
  {
    path: '/game',
    name: 'Game',
    component: () => import('@/views/GameView.vue')
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue')
  },
  {
    path: '/achievements',
    name: 'Achievements',
    component: () => import('@/views/AchievementsView.vue')
  },
  {
    path: '/chat-after',
    name: 'ChatAfter',
    component: () => import('@/views/ChatAfterStoryView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
