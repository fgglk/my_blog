import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/article/:id',
    name: 'ArticleDetail',
    component: () => import('@/views/ArticleDetail.vue'),
    meta: { title: '文章详情' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/profile/:id',
    name: 'UserProfile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '用户主页' }
  },
  {
    path: '/write',
    name: 'Write',
    component: () => import('@/views/Write.vue'),
    meta: { title: '写文章', requiresAuth: true }
  },
  {
    path: '/write/:id',
    name: 'EditArticle',
    component: () => import('@/views/Write.vue'),
    meta: { title: '编辑文章', requiresAuth: true }
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/views/Search.vue'),
    meta: { title: '搜索' }
  },
  {
    path: '/articles',
    name: 'Articles',
    component: () => import('@/views/Articles.vue'),
    meta: { title: '所有文章' }
  },
  {
    path: '/favorites',
    name: 'Favorites',
    component: () => import('@/views/Favorites.vue'),
    meta: { title: '我的收藏', requiresAuth: true }
  },
  {
    path: '/user-management',
    name: 'UserManagement',
    component: () => import('@/views/UserManagement.vue'),
    meta: { title: '用户管理', requiresAuth: true, requiresAdmin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  console.log('路由守卫触发，目标路径:', to.path)
  console.log('路由元信息:', to.meta)
  
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 我的博客` : '我的博客'
  
  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('token')
    console.log('需要认证，token存在:', !!token)
    if (!token) {
      console.log('跳转到登录页面')
      next('/login')
      return
    }
  }
  
  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin) {
    console.log('需要管理员权限')
    const token = localStorage.getItem('token')
    if (!token) {
      console.log('没有token，跳转到登录页面')
      next('/login')
      return
    }
    
    const userInfo = localStorage.getItem('userInfo')
    console.log('localStorage中的用户信息:', userInfo)
    if (!userInfo) {
      // 如果没有用户信息，先跳转到登录页面
      console.log('没有用户信息，跳转到登录页面')
      next('/login')
      return
    }
    
    try {
      const user = JSON.parse(userInfo)
      console.log('解析后的用户信息:', user)
      if (!user || user.role !== 'admin') {
        // 如果不是管理员，跳转到个人中心
        console.log('不是管理员，跳转到个人中心')
        next('/profile')
        return
      }
      console.log('管理员权限验证通过')
    } catch (error) {
      console.error('解析用户信息失败:', error)
      // 解析失败，清除可能损坏的数据并跳转到登录页面
      localStorage.removeItem('userInfo')
      next('/login')
      return
    }
  }
  
  console.log('路由守卫通过，允许访问:', to.path)
  next()
})

export default router 