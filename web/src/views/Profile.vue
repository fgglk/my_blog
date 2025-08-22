<template>
    <div class="profile-page">
      <header class="header">
        <div class="container">
          <div class="header-content">
            <div class="logo">
              <router-link to="/" class="logo-link">
                <h1>我的博客</h1>
              </router-link>
            </div>
            <nav class="nav">
              <router-link to="/" class="nav-item">首页</router-link>
              <router-link to="/search" class="nav-item">搜索</router-link>
              <router-link to="/write" class="nav-item">写文章</router-link>
              <router-link to="/profile" class="nav-item active">个人中心</router-link>
              <el-button @click="handleLogout" link>退出</el-button>
            </nav>
          </div>
        </div>
      </header>
  
      <main class="main">
        <div class="container">
          <div class="profile-container">
            <!-- 左侧边栏 -->
            <aside class="sidebar">
              <!-- 个人信息卡片 -->
              <div class="profile-card">
                <div class="profile-header">
                  <div class="avatar-section">
                    <div class="avatar-wrapper">
                      <el-avatar 
                        :size="120" 
                        :src="avatarUrlWithCacheBust"
                        :key="`avatar-${avatarKey}-${forceRefreshAvatar}`"
                        class="user-avatar"
                        @error="handleAvatarError"
                      >
                        {{ displayUserInfo?.nickname?.charAt(0) || 'U' }}
                      </el-avatar>
                      <div v-if="!isViewingOtherUser" class="avatar-overlay" @click="triggerAvatarUpload">
                        <el-icon><Camera /></el-icon>
                        <span>更换头像</span>
                      </div>
                      <input
                        ref="avatarInput"
                        type="file"
                        accept="image/*"
                        style="display: none"
                        @change="handleAvatarChange"
                      />
                    </div>
                  </div>
                  <div class="user-info">
                    <h2 class="user-name">{{ displayUserInfo?.nickname || '未设置昵称' }}</h2>
                    <p class="user-username">@{{ displayUserInfo?.username }}</p>
                    <p class="user-email">{{ displayUserInfo?.email }}</p>
                    <div class="user-bio">
                      <p>{{ displayUserInfo?.bio || '这个人很懒，什么都没有留下...' }}</p>
                    </div>
                  </div>
                </div>
                <div v-if="!isViewingOtherUser" class="profile-actions">
                  <el-button @click="showEditDialog = true" type="primary" size="large">
                    <el-icon><Edit /></el-icon>
                    编辑资料
                  </el-button>
                  <el-button @click="showPasswordDialog = true" size="large">
                    <el-icon><Lock /></el-icon>
                    修改密码
                  </el-button>
                </div>
              </div>
  
              <!-- 内容统计卡片 -->
              <div class="stats-card">
                <h3 class="stats-title">内容统计</h3>
                <div class="stats-grid">
                  <div class="stats-item">
                    <div class="stats-icon">
                      <el-icon><Document /></el-icon>
                    </div>
                    <div class="stats-info">
                      <span class="stats-number">{{ userStats.articleCount }}</span>
                      <span class="stats-label">文章</span>
                    </div>
                  </div>
                  <div class="stats-item">
                    <div class="stats-icon">
                      <el-icon><View /></el-icon>
                    </div>
                    <div class="stats-info">
                      <span class="stats-number">{{ userStats.totalViews }}</span>
                      <span class="stats-label">阅读</span>
                    </div>
                  </div>
                  <div class="stats-item">
                    <div class="stats-icon">
                      <el-icon><Star /></el-icon>
                    </div>
                    <div class="stats-info">
                      <span class="stats-number">{{ userStats.totalLikes }}</span>
                      <span class="stats-label">点赞</span>
                    </div>
                  </div>
                  <div class="stats-item">
                    <div class="stats-icon">
                      <el-icon><ChatDotRound /></el-icon>
                    </div>
                    <div class="stats-info">
                      <span class="stats-number">{{ userStats.totalComments }}</span>
                      <span class="stats-label">评论</span>
                    </div>
                  </div>
                </div>
              </div>
  
              <!-- 快捷操作卡片 -->
              <div v-if="!isViewingOtherUser" class="quick-actions-card">
                <h3 class="quick-actions-title">快捷操作</h3>
                <div class="quick-actions-list">
                  <div class="quick-action-item" @click="$router.push('/write')">
                    <div class="action-icon">
                      <el-icon><Document /></el-icon>
                    </div>
                    <span class="action-text">发布新文章</span>
                    <el-icon class="arrow-icon"><ArrowRight /></el-icon>
                  </div>
                  <div class="quick-action-item" @click="$router.push('/favorites')">
                    <div class="action-icon">
                      <el-icon><Collection /></el-icon>
                    </div>
                    <span class="action-text">我的收藏</span>
                    <el-icon class="arrow-icon"><ArrowRight /></el-icon>
                  </div>

                </div>
              </div>
            </aside>
  
            <!-- 右侧主内容区域 -->
            <main class="main-content">
              <!-- 文章列表标题 -->
              <div class="articles-header">
                <h2 class="articles-title">{{ isViewingOtherUser ? '他的文章' : '我的文章' }}</h2>
                <div v-if="!isViewingOtherUser" class="articles-controls">
                  <el-button type="primary" @click="$router.push('/write')" size="small">
                    <el-icon><Plus /></el-icon>
                    写文章
                  </el-button>
                </div>
              </div>

              <!-- 文章类型选项卡 -->
              <div class="articles-tabs">
                <el-tabs v-model="activeTab" @tab-change="handleTabChange">
                  <el-tab-pane label="已发布" name="published">
                    <template #label>
                      <span class="tab-label">
                        <el-icon><Document /></el-icon>
                        已发布 ({{ publishedCount }})
                      </span>
                    </template>
                  </el-tab-pane>
                  <el-tab-pane v-if="!isViewingOtherUser" label="草稿箱" name="draft">
                    <template #label>
                      <span class="tab-label">
                        <el-icon><EditPen /></el-icon>
                        草稿箱 ({{ draftCount }})
                      </span>
                    </template>
                  </el-tab-pane>
                </el-tabs>
              </div>
  
              <!-- 文章列表 -->
              <div v-if="userArticles.length === 0" class="empty-articles">
                <el-empty description="暂无更多文章">
                  <p class="empty-tip">您可以点击下方按钮发布新文章</p>
                  <el-button type="primary" @click="$router.push('/write')">
                    <el-icon><Document /></el-icon>
                    发布新文章
                  </el-button>
                </el-empty>
              </div>
              <div v-else class="articles-list">
                <article 
                  v-for="article in userArticles" 
                  :key="article.id"
                  class="article-card"
                  @click="goToArticle(article.id)"
                >
                  <div class="article-image">
                    <img :src="getArticleImage(article)" alt="文章配图" />
                  </div>
                  <div class="article-content">
                    <div class="article-meta">
                      <el-tag size="small" :type="(article.status === '1' || article.is_published) ? 'primary' : 'warning'">
                        {{ (article.status === '1' || article.is_published) ? (article.category?.name || '未分类') : '草稿' }}
                      </el-tag>
                      <span class="article-date">{{ formatDate(article.created_at) }}</span>
                      <div v-if="article.status === '1' || article.is_published" class="article-stats">
                        <span class="stat-item">
                          <el-icon><View /></el-icon>
                          {{ article.view_count }} 浏览
                        </span>
                        <span class="stat-item">
                          <el-icon><ChatDotRound /></el-icon>
                          {{ article.comment_count }} 评论
                        </span>
                        <span class="stat-item">
                          <el-icon><Star /></el-icon>
                          {{ article.like_count }} 点赞
                        </span>
                        <span class="stat-item">
                          <el-icon><Collection /></el-icon>
                          {{ article.favorite_count }} 收藏
                        </span>
                      </div>
                      <div v-else class="draft-info">
                        <span class="draft-text">草稿 · {{ formatDate(article.updated_at) }} 保存</span>
                      </div>
                    </div>
                    <h3 class="article-title">{{ article.title || '未命名草稿' }}</h3>
                    <p class="article-summary">{{ getArticleSummary(article.content) }}</p>
                    <div class="article-footer">
                      <div class="article-actions">
                        <span 
                          v-if="article.status === '1' || article.is_published" 
                          class="read-more"
                          @click.stop="goToArticle(article.id)"
                        >
                          阅读全文 →
                        </span>
                        <span v-else class="draft-placeholder">{{ getWordCount(article.content) }} 字</span>
                        <div v-if="!isViewingOtherUser" class="action-buttons">
                          <el-button type="success" size="small" @click.stop="goToEdit(article.id)">
                            <el-icon><Edit /></el-icon>
                          </el-button>
                          <el-button type="danger" size="small" @click.stop="deleteArticle(article.id)">
                            <el-icon><Delete /></el-icon>
                          </el-button>
                        </div>
                      </div>
                    </div>
                  </div>
                </article>
              </div>
            </main>
          </div>
        </div>
      </main>
  
      <!-- 编辑资料对话框 -->
      <el-dialog v-model="showEditDialog" title="编辑资料" width="600px" class="edit-dialog">
        <el-form :model="editForm" label-width="80px" class="edit-form">
          <el-form-item label="昵称">
            <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="editForm.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="个人简介">
            <el-input 
              v-model="editForm.bio" 
              type="textarea" 
              :rows="4"
              placeholder="这个人很懒，什么都没有留下..."
              maxlength="200"
              show-word-limit
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="updateProfile" :loading="updateLoading">
            保存
          </el-button>
        </template>
      </el-dialog>
  
      <!-- 修改密码对话框 -->
      <el-dialog v-model="showPasswordDialog" title="修改密码" width="500px">
        <el-form :model="passwordForm" label-width="100px">
          <el-form-item label="原密码">
            <el-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入原密码" />
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" />
          </el-form-item>
          <el-form-item label="确认密码">
            <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showPasswordDialog = false">取消</el-button>
          <el-button type="primary" @click="changePassword" :loading="passwordLoading">
            修改
          </el-button>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useArticleStore } from '@/stores/article'
import { userApi } from '@/api/user'
import { articleApi } from '@/api/article'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Camera, Edit, Lock, Document, View, Star, ChatDotRound, Plus, 
  ArrowRight, Delete, Collection, EditPen
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Article } from '@/types/article'
import { getPlainTextSummary } from '@/utils/markdown'
  
  const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const articleStore = useArticleStore()

// 当前查看的用户信息（用于查看其他用户时）
const currentViewingUser = ref<any>(null)

// 判断是否正在查看其他用户
const isViewingOtherUser = computed(() => {
  return route.params.id !== undefined
})

// 当前用户ID
const currentUserId = computed(() => {
  if (isViewingOtherUser.value) {
    return parseInt(route.params.id as string)
  }
  return userStore.userInfo?.id
})

// 显示的用户信息
const displayUserInfo = computed(() => {
  if (isViewingOtherUser.value && currentViewingUser.value) {
    return currentViewingUser.value
  }
  return userStore.userInfo
})
  
  const showEditDialog = ref(false)
  const showPasswordDialog = ref(false)
  const updateLoading = ref(false)
  const passwordLoading = ref(false)
  const avatarInput = ref<HTMLInputElement>()

const avatarKey = ref(0)
const activeTab = ref('published')

  // 用户统计信息
  const userStats = reactive({
    articleCount: 0,
    draftCount: 0,
    totalViews: 0,
    totalLikes: 0,
    totalComments: 0
  })
  
  // 用户文章列表
  const userArticles = ref<Article[]>([])
  const allUserArticles = ref<Article[]>([])

  // 计算发布和草稿数量
  const publishedCount = computed(() => {
    return allUserArticles.value.filter(article => article.status === '1' || article.is_published).length
  })

  const draftCount = computed(() => {
    return allUserArticles.value.filter(article => article.status === '0' || (!article.is_published && article.status !== '1')).length
  })
  
  const editForm = reactive({
    nickname: userStore.userInfo?.nickname || '',
    email: userStore.userInfo?.email || '',
    bio: userStore.userInfo?.bio || '这个人很懒，什么都没有留下...'
  })
  
  const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 头像URL缓存破坏计算属性
const avatarUrlWithCacheBust = computed(() => {
  const avatarUrl = displayUserInfo.value?.avatar
  if (!avatarUrl) return ''
  
  // 使用更简单的缓存破坏机制
  const timestamp = Date.now()
  return `${avatarUrl}?v=${timestamp}`
})

  // 强制刷新头像的响应式变量
  const forceRefreshAvatar = ref(0)
  
  // 处理头像加载错误
  const handleAvatarError = (event: Event) => {
    const target = event.target as HTMLImageElement
    console.error('头像加载失败:', target.src)
  }

// 触发头像上传
  const triggerAvatarUpload = () => {
    avatarInput.value?.click()
  }
  
  // 处理头像上传
  const handleAvatarChange = async (event: Event) => {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    if (!file) return
  
    try {
      console.log('开始上传头像...')
      const result = await userStore.uploadAvatar(file)
      console.log('头像上传结果:', result)
      if (result.success) {
        ElMessage.success('头像上传成功')
        console.log('当前用户头像URL:', userStore.userInfo?.avatar)
        
        // 强制刷新头像显示 - 通过更新key强制重新渲染
        avatarKey.value = Date.now()
        forceRefreshAvatar.value = Date.now()
        
        // 清空文件输入框，允许重复选择同一文件
        if (target) {
          target.value = ''
        }
      } else {
        ElMessage.error(result.message || '头像上传失败')
      }
    } catch (error) {
      console.error('头像上传错误:', error)
      ElMessage.error('头像上传失败')
    }
  }
  
  // 更新个人资料
  const updateProfile = async () => {
    if (!editForm.nickname.trim()) {
      ElMessage.warning('请输入昵称')
      return
    }
  
    updateLoading.value = true
    try {
      const result = await userStore.updateUserInfo(editForm)
      if (result.success) {
        ElMessage.success('资料更新成功')
        showEditDialog.value = false
        // 重新加载用户信息
        await userStore.getUserInfo()
      } else {
        ElMessage.error(result.message || '更新失败')
      }
    } catch (error) {
      ElMessage.error('更新失败')
    } finally {
      updateLoading.value = false
    }
  }
  
  // 修改密码
  const changePassword = async () => {
    if (!passwordForm.oldPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
      ElMessage.warning('请填写完整信息')
      return
    }
  
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      ElMessage.error('两次输入的密码不一致')
      return
    }
  
    if (passwordForm.newPassword.length < 6) {
      ElMessage.error('新密码长度不能少于6位')
      return
    }
  
    passwordLoading.value = true
    try {
      const result = await userStore.changePassword(passwordForm)
      if (result.success) {
        ElMessage.success('密码修改成功')
        showPasswordDialog.value = false
        // 清空表单
        passwordForm.oldPassword = ''
        passwordForm.newPassword = ''
        passwordForm.confirmPassword = ''
      } else {
        ElMessage.error(result.message || '修改失败')
      }
    } catch (error) {
      ElMessage.error('修改失败')
    } finally {
      passwordLoading.value = false
    }
  }
  
  // 删除文章
  const deleteArticle = async (articleId: number) => {
    try {
      await ElMessageBox.confirm('确定要删除这篇文章吗？', '确认删除', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      
      const result = await articleStore.deleteArticle(articleId)
      if (result.success) {
        ElMessage.success('文章删除成功')
        await loadUserArticles()
      } else {
        ElMessage.error(result.message || '删除失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('删除失败')
      }
    }
  }
  
  // 获取文章摘要
  const getArticleSummary = (content: string) => {
    return getPlainTextSummary(content, 140)
  }

  // 字数统计
  const getWordCount = (content: string) => {
    if (!content) return 0
    // 移除HTML标签
    const cleanContent = content.replace(/<[^>]*>/g, '')
    // 移除Markdown语法
    const markdownFree = cleanContent.replace(/[#*`~\[\]()!]/g, '')
    // 统计中文字符和英文单词
    const chineseChars = (markdownFree.match(/[\u4e00-\u9fa5]/g) || []).length
    const englishWords = markdownFree.split(/\s+/).filter(word => /[a-zA-Z]/.test(word)).length
    return chineseChars + englishWords
  }
  
  // 获取文章配图
  const getArticleImage = (article: Article) => {
    // 如果文章有封面图片，使用封面图片
    if (article.cover_image) {
      return article.cover_image
    }
    // 否则使用占位图片
    return `https://picsum.photos/300/200?random=${article.id}`
  }
  


  // 根据选项卡过滤文章
  const filterArticlesByTab = () => {
    if (activeTab.value === 'published') {
      userArticles.value = allUserArticles.value.filter(article => article.status === '1' || article.is_published)
    } else if (activeTab.value === 'draft') {
      userArticles.value = allUserArticles.value.filter(article => article.status === '0' || (!article.is_published && article.status !== '1'))
    }
  }

  // 选项卡切换处理
  const handleTabChange = (tabName: string) => {
    activeTab.value = tabName
    filterArticlesByTab()
  }
  
  // 跳转到文章详情
  const goToArticle = (articleId: number) => {
    router.push(`/article/${articleId}`)
  }
  
  // 跳转到编辑文章
  const goToEdit = (articleId: number) => {
    router.push(`/write/${articleId}`)
  }
  
  // 格式化日期
  const formatDate = (date: string) => {
    return dayjs(date).format('YYYY-MM-DD')
  }
  
  // 退出登录
  const handleLogout = () => {
    userStore.logout()
    ElMessage.success('退出成功')
    router.push('/')
  }

  // 加载用户信息
  const loadUserInfo = async (userId: number) => {
    try {
      if (isViewingOtherUser.value) {
        const response = await userApi.getUserById(userId)
        if (response.code === 0) {
          currentViewingUser.value = response.data
          console.log('加载其他用户信息成功:', response.data)
        } else {
          ElMessage.error('获取用户信息失败')
        }
      } else {
        currentViewingUser.value = null
      }
    } catch (error) {
      console.error('加载用户信息失败:', error)
      ElMessage.error('加载用户信息失败')
    }
  }

  // 加载用户文章
  const loadUserArticles = async () => {
    try {
      console.log('开始加载用户文章...')
      let response
      if (isViewingOtherUser.value) {
        // 获取其他用户的所有文章（不分页）
        response = await articleApi.getArticlesByUserId(currentUserId.value!, 1, 1000)
        console.log('获取其他用户文章响应:', response)
        if (response.code === 0 && response.data) {
          allUserArticles.value = response.data.list || []
          console.log('其他用户文章列表:', allUserArticles.value)
          filterArticlesByTab()
          // 对于其他用户，只统计已发布的文章
          const publishedArticles = allUserArticles.value.filter(article => article.status === '1' || article.is_published)
          userStats.articleCount = publishedArticles.length
          userStats.draftCount = 0 // 其他用户看不到草稿
          userStats.totalViews = publishedArticles.reduce((sum, article) => sum + article.view_count, 0)
          userStats.totalLikes = publishedArticles.reduce((sum, article) => sum + article.like_count, 0)
          userStats.totalComments = publishedArticles.reduce((sum, article) => sum + article.comment_count, 0)
          console.log('其他用户统计信息:', userStats)
        } else {
          console.log('获取其他用户文章失败:', response.msg)
        }
      } else {
        // 获取当前用户的所有文章（直接调用API，不分页）
        response = await articleApi.getUserArticles(1, 1000)
        console.log('获取当前用户文章响应:', response)
        if (response.code === 0 && response.data) {
          allUserArticles.value = response.data.list || []
          console.log('当前用户文章列表:', allUserArticles.value)
          console.log('当前用户文章总数:', response.data.total)
          filterArticlesByTab()
          // 对于当前用户，分别统计已发布和草稿文章
          const publishedArticles = allUserArticles.value.filter(article => article.status === '1' || article.is_published)
          const draftArticles = allUserArticles.value.filter(article => article.status === '0' || !article.is_published)
          userStats.articleCount = publishedArticles.length
          userStats.draftCount = draftArticles.length
          userStats.totalViews = publishedArticles.reduce((sum, article) => sum + article.view_count, 0)
          userStats.totalLikes = publishedArticles.reduce((sum, article) => sum + article.like_count, 0)
          userStats.totalComments = publishedArticles.reduce((sum, article) => sum + article.comment_count, 0)
          console.log('当前用户统计信息:', userStats)
          console.log('发布文章数:', publishedArticles.length)
          console.log('草稿文章数:', draftArticles.length)
        } else {
          console.log('获取当前用户文章失败:', response.msg)
        }
      }
    } catch (error) {
      console.error('加载用户文章失败:', error)
    }
  }

  // 监听路由参数变化
  watch(() => route.params.id, async (newId, oldId) => {
    if (newId !== oldId) {
      await loadUserInfo(currentUserId.value!)
      await loadUserArticles()
    }
  }, { immediate: false })
  
  onMounted(async () => {
    await loadUserInfo(currentUserId.value!)
    await loadUserArticles()
  })
  </script>
  
  <style lang="scss" scoped>
  .profile-page {
    min-height: 100vh;
    background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
    position: relative;
  }
  
  .header {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    box-shadow: 0 2px 20px rgba(0, 0, 0, 0.08);
    position: sticky;
    top: 0;
    z-index: 100;
  }
  
  .header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 15px 0;
  }
  
  .logo-link {
    text-decoration: none;
    
    h1 {
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
      -webkit-background-clip: text;
      background-clip: text;
      -webkit-text-fill-color: transparent;
      font-size: 24px;
      font-weight: bold;
    }
  }
  
  .nav {
    display: flex;
    align-items: center;
    gap: 20px;
  }
  
  .nav-item {
    text-decoration: none;
    color: #333;
    font-size: 16px;
    transition: all 0.3s ease;
    padding: 8px 16px;
    border-radius: 20px;
    
    &:hover {
      color: #3b82f6;
      background: rgba(59, 130, 246, 0.1);
    }
    
    &.active {
      color: #3b82f6;
      background: rgba(59, 130, 246, 0.1);
      font-weight: bold;
    }
  }
  
  .main {
    padding: 20px 0;
  }
  
  .profile-container {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 5px;
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 20px;
  align-items: start;
}
  
  // 左侧边栏
  .sidebar {
    display: flex;
    flex-direction: column;
    gap: 25px;
  }
  
  // 个人信息卡片
  .profile-card {
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    padding: 30px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.3);
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 3px;
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
    }
  }
  
  .profile-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    margin-bottom: 25px;
  }
  
  .avatar-section {
    margin-bottom: 20px;
    
    .avatar-wrapper {
      position: relative;
      display: inline-block;
      
      .user-avatar {
        background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
        color: #fff;
        font-size: 48px;
        font-weight: bold;
        box-shadow: 0 4px 15px rgba(59, 130, 246, 0.3);
        border: 4px solid rgba(255, 255, 255, 0.8);
      }
      
      .avatar-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.6);
        color: #fff;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        opacity: 0;
        transition: all 0.3s ease;
        cursor: pointer;
        
        .el-icon {
          font-size: 24px;
          margin-bottom: 5px;
        }
        
        span {
          font-size: 12px;
        }
      }
      
      &:hover .avatar-overlay {
        opacity: 1;
      }
    }
  }
  
  .user-info {
    width: 100%;
    
    .user-name {
      font-size: 24px;
      font-weight: 700;
      color: #333;
      margin-bottom: 8px;
    }
    
    .user-username {
      font-size: 14px;
      color: #3b82f6;
      margin-bottom: 8px;
      font-weight: 500;
    }
    
    .user-email {
      font-size: 14px;
      color: #666;
      margin-bottom: 15px;
    }
    
    .user-bio {
      p {
        font-size: 14px;
        color: #666;
        line-height: 1.6;
        font-style: italic;
        background: rgba(59, 130, 246, 0.05);
        padding: 12px;
        border-radius: 8px;
        border-left: 3px solid #3b82f6;
      }
    }
  }
  
  .profile-actions {
    display: flex;
    flex-direction: column;
    gap: 12px;
    
    .el-button {
      border-radius: 25px;
      padding: 12px 24px;
      font-weight: 500;
      transition: all 0.3s ease;
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 15px rgba(59, 130, 246, 0.3);
      }
    }
  }
  
  // 统计卡片
  .stats-card {
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    padding: 25px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.3);
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 3px;
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
    }
  }
  
  .stats-title {
    font-size: 18px;
    font-weight: bold;
    color: #333;
    margin-bottom: 20px;
    text-align: center;
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
  }
  
  .stats-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 15px;
    border-radius: 12px;
    background: rgba(59, 130, 246, 0.05);
    border: 1px solid rgba(59, 130, 246, 0.1);
    transition: all 0.3s ease;
    
    &:hover {
      background: rgba(59, 130, 246, 0.1);
      transform: translateY(-2px);
      box-shadow: 0 4px 15px rgba(59, 130, 246, 0.2);
    }
    
    &:nth-child(1) .stats-icon {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }
    
    &:nth-child(2) .stats-icon {
      background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    }
    
    &:nth-child(3) .stats-icon {
      background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
    }
    
    &:nth-child(4) .stats-icon {
      background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
    }
    
    .stats-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 40px;
      height: 40px;
      border-radius: 50%;
      color: #fff;
      
      .el-icon {
        font-size: 18px;
      }
    }
    
    .stats-info {
      display: flex;
      flex-direction: column;
      gap: 2px;
    }
    
    .stats-number {
      font-size: 20px;
      font-weight: bold;
      color: #333;
      line-height: 1;
    }
    
    .stats-label {
      font-size: 12px;
      color: #666;
      font-weight: 500;
    }
  }
  
  // 快捷操作卡片
  .quick-actions-card {
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    padding: 25px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.3);
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 3px;
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
    }
  }
  
  .quick-actions-title {
    font-size: 18px;
    font-weight: bold;
    color: #333;
    margin-bottom: 20px;
    text-align: center;
  }
  
  .quick-actions-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .quick-action-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 15px;
    border-radius: 12px;
    background: rgba(59, 130, 246, 0.05);
    border: 1px solid rgba(59, 130, 246, 0.1);
    cursor: pointer;
    transition: all 0.3s ease;
    position: relative;
    
    &:hover {
      background: rgba(59, 130, 246, 0.1);
      transform: translateX(5px);
      box-shadow: 0 4px 15px rgba(59, 130, 246, 0.2);
    }
    
    .action-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 35px;
      height: 35px;
      border-radius: 50%;
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
      color: #fff;
      
      .el-icon {
        font-size: 16px;
      }
    }
    
    .action-text {
      flex: 1;
      font-size: 14px;
      color: #333;
      font-weight: 500;
    }
    
    .arrow-icon {
      color: #3b82f6;
      font-size: 16px;
    }
    

  }
  
  // 右侧主内容区域
  .main-content {
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    padding: 30px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.3);
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 3px;
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
    }
  }
  
  .articles-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    .articles-title {
      font-size: 24px;
      font-weight: bold;
      color: #333;
    }
    
    .articles-controls {
      display: flex;
      align-items: center;
      gap: 15px;
    }
  }

  .articles-tabs {
    margin-bottom: 30px;
    
    .tab-label {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .el-icon {
        font-size: 16px;
      }
    }
    
    :deep(.el-tabs__header) {
      margin-bottom: 20px;
    }
    
    :deep(.el-tabs__nav-wrap::after) {
      height: 1px;
      background-color: #e4e7ed;
    }
    
    :deep(.el-tabs__active-bar) {
      background-color: #3b82f6;
    }
    
    :deep(.el-tabs__item.is-active) {
      color: #3b82f6;
    }
  }
  
  .empty-articles {
    text-align: center;
    padding: 60px 0;
    
    .empty-tip {
      color: #666;
      margin: 15px 0;
      font-size: 14px;
    }
  }
  
  .articles-list {
  .article-card {
    background: #fff;
    border-radius: 12px;
    overflow: hidden;
    margin-bottom: 20px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
    cursor: pointer;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    display: flex;
    height: 160px; // 固定高度
    width: 100%;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
    }
  }
  
  .article-image {
    width: 200px; // 缩小图片宽度
    flex-shrink: 0;
    overflow: hidden;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
  
  .article-content {
    padding: 16px;
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    overflow: hidden; // 防止内容溢出
  }
  
  .article-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
    
    .article-date {
      color: #666;
      font-size: 12px;
    }
    
    .article-stats {
      display: flex;
      align-items: center;
      gap: 12px;
      
      .stat-item {
        display: flex;
        align-items: center;
        gap: 2px;
        color: #666;
        font-size: 11px;
        
        .el-icon {
          font-size: 12px;
        }
      }
    }
  }
  
  .article-title {
    font-size: 16px;
    font-weight: bold;
    color: #333;
    margin-bottom: 8px;
    line-height: 1.3;
    overflow: hidden;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    height: 42px; // 固定标题高度
    /* Standard properties for compatibility */
    display: -moz-box;
    -moz-box-orient: vertical;
    display: box;
    box-orient: vertical;
    line-clamp: 2;
  }
  
  .article-summary {
    color: #666;
    line-height: 1.4;
    margin-bottom: 12px;
    overflow: hidden;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    font-size: 13px;
    height: 36px; // 固定摘要高度
    /* Standard properties for compatibility */
    display: -moz-box;
    -moz-box-orient: vertical;
    display: box;
    box-orient: vertical;
    line-clamp: 2;
  }
  
  .article-footer {
    display: flex;
    align-items: center;
    justify-content: flex-end; // 右对齐
    margin-top: auto; // 推到底部
  }
  
  .article-actions {
    display: flex;
    align-items: center;
    gap: 12px; // 增加间距
  }
  
  .action-buttons {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .draft-info {
    color: #f59e0b;
    font-size: 12px;
    font-style: italic;
  }

  .draft-text {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .draft-placeholder {
    color: #9ca3af;
    font-size: 12px;
    font-style: italic;
  }
  
  .read-more {
    color: #3b82f6;
    text-decoration: none;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    
    &:hover {
      text-decoration: underline;
      color: #2563eb;
    }
  }
}

// 响应式设计
@media (max-width: 1200px) {
  .profile-container {
    max-width: 100%;
    padding: 0;
    grid-template-columns: 320px 1fr;
    gap: 25px;
  }
}

@media (max-width: 1024px) {
  .profile-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .sidebar {
    order: -1;
  }
}

@media (max-width: 768px) {
  .profile-container {
    padding: 0 10px;
  }
  
  .profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 15px;
  }
  
  .profile-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .stats-item {
    flex-direction: row;
    justify-content: center;
    gap: 10px;
  }
  
  .articles-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .article-card {
    flex-direction: column;
    height: 200px; // 移动端稍微高一点
  }
  
  .article-image {
    width: 100%;
    height: 120px;
  }
  
  .article-content {
    flex: 1;
    padding: 12px;
  }
  
  .article-title {
    font-size: 16px;
    height: auto;
    -webkit-line-clamp: 2;
    /* Standard properties for compatibility */
    line-clamp: 2;
  }
  
  .article-summary {
    height: auto;
    -webkit-line-clamp: 2;
    /* Standard properties for compatibility */
    line-clamp: 2;
    font-size: 13px;
  }
  
  .article-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .article-stats {
    flex-wrap: wrap;
    gap: 8px;
  }
}

@media (max-width: 480px) {
  .profile-card,
  .stats-card,
  .quick-actions-card,
  .main-content {
    padding: 20px;
  }
  
  .stats-grid {
    gap: 10px;
  }
  
  .stats-item {
    padding: 12px;
    
    .stats-icon {
      width: 35px;
      height: 35px;
      
      .el-icon {
        font-size: 16px;
      }
    }
    
    .stats-number {
      font-size: 18px;
    }
  }
  
  .quick-action-item {
    padding: 10px 12px;
    
    .action-icon {
      width: 30px;
      height: 30px;
      
      .el-icon {
        font-size: 14px;
      }
    }
    
    .action-text {
      font-size: 13px;
    }
  }
  
  .profile-container {
    padding: 0 5px;
  }
  
  .articles-header {
    .articles-title {
      font-size: 20px;
    }
  }
  
  .article-card {
    margin-bottom: 15px;
    height: 180px;
  }
  
  .article-image {
    height: 100px;
  }
  
  .article-content {
    padding: 10px;
  }
  
  .article-title {
    font-size: 14px;
    margin-bottom: 6px;
  }
  
  .article-summary {
    font-size: 12px;
    margin-bottom: 8px;
  }
}
</style>