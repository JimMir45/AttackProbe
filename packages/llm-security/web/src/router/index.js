import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/target'
  },
  {
    path: '/target',
    name: 'Target',
    component: () => import('../views/Target.vue')
  },
  {
    path: '/testcase',
    name: 'TestCase',
    component: () => import('../views/TestCase.vue')
  },
  {
    path: '/task',
    name: 'Task',
    component: () => import('../views/Task.vue')
  },
  {
    path: '/task/:id',
    name: 'TaskDetail',
    component: () => import('../views/TaskDetail.vue')
  },
  {
    path: '/report',
    name: 'Report',
    component: () => import('../views/Report.vue')
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/Settings.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
