<template>
  <div class="home">
    <!-- 导航栏 -->
    <header class="header">
      <div class="container">
        <div class="header-content">
          <div class="logo">
            <el-icon class="logo-icon"><Promotion /></el-icon>
            <h1>我的博客</h1>
          </div>
          <nav class="nav">
            <router-link to="/" class="nav-item">
              <el-icon><House /></el-icon>
              首页
            </router-link>
            <router-link to="/search" class="nav-item">
              <el-icon><Search /></el-icon>
              搜索
            </router-link>
            <template v-if="userStore.userInfo">
              <router-link to="/write" class="nav-item">
                <el-icon><Edit /></el-icon>
                写文章
              </router-link>
              <router-link to="/profile" class="nav-item">
                <el-icon><User /></el-icon>
                个人中心
              </router-link>
              <el-button @click="handleLogout" link>退出</el-button>
            </template>
            <template v-else>
              <router-link to="/login" class="nav-item">
                <el-icon><User /></el-icon>
                登录
              </router-link>
              <router-link to="/register" class="nav-item">
                <el-icon><UserFilled /></el-icon>
                注册
              </router-link>
            </template>
          </nav>
        </div>
      </div>
    </header>

    <!-- Hero区域 -->
    <section class="hero">
      <div class="container">
        <div class="hero-content">
          <h1 class="hero-title">探索思想的无限可能</h1>
          <p class="hero-subtitle">分享技术、生活与思考，记录成长的每一步</p>
                     <div class="hero-stats">
             <div class="stat-item">
               <el-icon><Document /></el-icon>
               <span>{{ stats.articleCount }} 篇文章</span>
             </div>
             <div class="stat-item">
               <el-icon><View /></el-icon>
               <span>{{ formatNumber(stats.viewCount) }} 阅读量</span>
             </div>
             <div class="stat-item">
               <el-icon><ChatDotRound /></el-icon>
               <span>{{ stats.commentCount }} 评论</span>
             </div>
             <div class="stat-item">
               <el-icon><Star /></el-icon>
               <span>{{ stats.likeCount }} 点赞</span>
             </div>
             <div class="stat-item">
               <el-icon><Collection /></el-icon>
               <span>{{ stats.favoriteCount }} 收藏</span>
             </div>
           </div>
        </div>
      </div>
    </section>

    <!-- 主要内容 -->
    <main class="main">
      <div class="container">
                 <div class="content">
           <!-- 文章列表 -->
           <div class="article-section">
             <div class="section-header">
               <h2 class="section-title">
                 <el-icon><Promotion /></el-icon>
                 最新文章
               </h2>
               <router-link to="/articles" class="view-all">
                 查看全部 →
               </router-link>
             </div>
             
             <div class="article-list">
               <div v-if="articleStore.loading" class="loading">
                 <el-skeleton :rows="3" animated />
                 <el-skeleton :rows="3" animated />
                 <el-skeleton :rows="3" animated />
               </div>
               <div v-else-if="articleStore.articles.length === 0" class="empty">
                 <el-empty description="暂无文章" />
               </div>
               <div v-else>
                 <article 
                   v-for="article in articleStore.articles" 
                   :key="article.id" 
                   class="article-card"
                   @click="goToArticle(article.id)"
                 >
                   <div class="article-image">
                     <img :src="getArticleImage(article)" alt="文章配图" />
                   </div>
                   <div class="article-content">
                                           <div class="article-meta">
                        <el-tag size="small" type="primary">{{ article.category.name }}</el-tag>
                        <span class="article-date">{{ formatDate(article.created_at) }}</span>
                        <div class="article-stats">
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
                      </div>
                     <h3 class="article-title">{{ article.title }}</h3>
                     <p class="article-summary">{{ getArticleSummary(article.content) }}</p>
                     <div class="article-footer">
                                               <div class="author">
                          <el-avatar :size="24" :src="article.author_avatar">
                            {{ article.author_name?.charAt(0) || 'U' }}
                          </el-avatar>
                          <span>{{ article.author_name }}</span>
                        </div>
                       <span class="read-more" @click.stop="goToArticle(article.id)">
                         阅读全文 →
                       </span>
                     </div>
                   </div>
                 </article>
               </div>
             </div>

             <!-- 加载更多 -->
             <div v-if="articleStore.articles.length > 0 && hasMore" class="load-more">
               <el-button type="primary" :loading="articleStore.loading" @click="loadMore">
                 <el-icon><Refresh /></el-icon>
                 加载更多
               </el-button>
             </div>
           </div>

           <!-- 分类和标签区域 -->
           <div class="categories-tags-section">
             <div class="categories-tags-grid">
               <!-- 文章分类 -->
               <div class="sidebar-card">
                 <h3 class="sidebar-title">
                   <el-icon><Document /></el-icon>
                   文章分类
                 </h3>
                 <div class="category-list">
                   <div 
                     v-for="category in categories" 
                     :key="category.id"
                     class="category-item"
                     @click="filterByCategory(category.id)"
                   >
                     <span class="category-name">{{ category.name }}</span>
                     <span class="category-count">{{ category.article_count }}</span>
                   </div>
                 </div>
               </div>

               <!-- 热门标签 -->
               <div class="sidebar-card">
                 <div class="sidebar-header">
                   <h3 class="sidebar-title">
                     <el-icon><Collection /></el-icon>
                     热门标签
                   </h3>
                   <el-button 
                     size="small" 
                     :icon="Refresh" 
                     @click="tagStore.refreshTags()"
                     :loading="tagStore.loading"
                     circle
                     class="refresh-btn"
                   />
                 </div>
                 <div class="tag-list">
                   <el-tag 
                     v-for="tag in tagStore.tags" 
                     :key="tag.id"
                     size="small"
                     class="tag-item"
                     @click="filterByTag(tag.id)"
                   >
                     #{{ tag.name }}
                   </el-tag>
                 </div>
               </div>
             </div>
           </div>
         </div>
      </div>
    </main>

    <!-- 订阅区域 -->
    <section class="subscribe-section">
      <div class="container">
        <div class="subscribe-content">
          <div class="subscribe-info">
            <h2>订阅我的博客</h2>
            <p>获取最新文章和技术分享，直接发送到您的邮箱</p>
            <div class="subscribe-form">
              <el-input
                v-model="subscribeEmail"
                placeholder="您的邮箱地址"
                size="large"
                class="subscribe-input"
              />
              <el-button type="primary" size="large" @click="handleSubscribe">
                立即订阅
              </el-button>
            </div>
            <p class="privacy-note">我们尊重您的隐私，不会向第三方分享您的信息</p>
          </div>
          <div class="subscribe-icon">
            <el-icon><Message /></el-icon>
          </div>
        </div>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-section">
            <h3>我的博客</h3>
            <p>分享技术、生活与思考，记录成长的每一步。在这里，我们一起探索技术的无限可能。</p>
            <div class="social-links">
              <el-button link size="small">
                <el-icon><Link /></el-icon>
              </el-button>
              <el-button link size="small">
                <el-icon><ChatDotRound /></el-icon>
              </el-button>
              <el-button link size="small">
                <el-icon><Message /></el-icon>
              </el-button>
            </div>
          </div>
          
          <div class="footer-section">
            <h3>快速链接</h3>
            <ul>
              <li><router-link to="/">首页</router-link></li>
              <li><router-link to="/articles">所有文章</router-link></li>
              <li><router-link to="/search">搜索</router-link></li>
              <li><router-link to="/profile">个人中心</router-link></li>
              <li><router-link to="/write">写文章</router-link></li>
            </ul>
          </div>
          
          <div class="footer-section">
            <h3>功能特色</h3>
            <ul>
              <li><a href="#" @click.prevent>Markdown 编辑器</a></li>
              <li><a href="#" @click.prevent>文章分类</a></li>
              <li><a href="#" @click.prevent>标签系统</a></li>
              <li><a href="#" @click.prevent>评论互动</a></li>
              <li><a href="#" @click.prevent>点赞收藏</a></li>
            </ul>
          </div>
          
          <div class="footer-section">
            <h3>关于博客</h3>
            <div class="contact-info">
              <p>
                <el-icon><Document /></el-icon>
                基于 Vue 3 + Go 开发
              </p>
              <p>
                <el-icon><Star /></el-icon>
                开源项目，欢迎贡献
              </p>
            </div>
          </div>
        </div>
        
        <div class="footer-bottom">
          <p>© 2025 我的博客，保留所有权利。</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useArticleStore } from '@/stores/article'
import { useTagStore } from '@/stores/tag'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  View, Promotion, ChatDotRound, House, Search, Edit, User, UserFilled,
  Document, Refresh, Collection, Message, Link, Star
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Article } from '@/types/article'
import { categoryApi, type Category } from '@/api/category'
import { getPlainTextSummary } from '@/utils/markdown'

const router = useRouter()
const userStore = useUserStore()
const articleStore = useArticleStore()
const tagStore = useTagStore()

const currentPage = ref(1)
const pageSize = ref(4) // 修改为4篇文章
const subscribeEmail = ref('')
const hasMore = ref(true) // 添加是否有更多文章的标识

// 统计数据
const stats = ref({
  articleCount: 0,
  viewCount: 0,
  commentCount: 0,
  likeCount: 0,
  favoriteCount: 0
})

const categories = ref<Category[]>([])

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

// 格式化数字
const formatNumber = (num: number) => {
  if (num >= 10000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 获取文章摘要
const getArticleSummary = (content: string) => {
  return getPlainTextSummary(content, 140)
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

// 跳转到文章详情
const goToArticle = (id: number) => {
  router.push(`/article/${id}`)
}

// 加载更多文章
const loadMore = async () => {
  if (!hasMore.value || articleStore.loading) return
  
  currentPage.value++
  try {
    await articleStore.getArticles(currentPage.value, pageSize.value, true) // 使用追加模式
    
    // 检查是否还有更多文章
    hasMore.value = articleStore.articles.length < articleStore.total
  } catch (error) {
    console.error('加载更多文章失败:', error)
    currentPage.value-- // 恢复页码
  }
}

// 加载文章列表
const loadArticles = async () => {
  currentPage.value = 1
  hasMore.value = true
  await articleStore.getArticles(currentPage.value, pageSize.value)
  
  // 检查是否还有更多文章
  if (articleStore.total > pageSize.value) {
    hasMore.value = true
  } else {
    hasMore.value = false
  }
}

// 获取统计数据
const loadStats = async () => {
  try {
    // 调用后端API获取真实统计数据
    const response = await fetch('/api/articles/stats')
    if (!response.ok) {
      throw new Error('获取统计数据失败')
    }
    
    const data = await response.json()
    if (data.code === 0) {
      const realStats = data.data
      
      // 添加数字动画效果
      animateNumber(stats.value.articleCount, realStats.articleCount, (value) => {
        stats.value.articleCount = value
      })
      
      animateNumber(stats.value.viewCount, realStats.viewCount, (value) => {
        stats.value.viewCount = value
      })
      
      animateNumber(stats.value.commentCount, realStats.commentCount, (value) => {
        stats.value.commentCount = value
      })
      
      animateNumber(stats.value.likeCount, realStats.likeCount, (value) => {
        stats.value.likeCount = value
      })
      
      animateNumber(stats.value.favoriteCount, realStats.favoriteCount, (value) => {
        stats.value.favoriteCount = value
      })
    } else {
      throw new Error(data.message || '获取统计数据失败')
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    // 如果API调用失败，使用默认数据
    const fallbackStats = {
      articleCount: 0,
      viewCount: 0,
      commentCount: 0,
      likeCount: 0,
      favoriteCount: 0
    }
    
    stats.value = fallbackStats
  }
}

// 数字动画函数
const animateNumber = (start: number, end: number, callback: (value: number) => void) => {
  const duration = 2000
  const startTime = Date.now()
  
  const animate = () => {
    const elapsed = Date.now() - startTime
    const progress = Math.min(elapsed / duration, 1)
    
    // 使用缓动函数
    const easeOutQuart = 1 - Math.pow(1 - progress, 4)
    const current = Math.floor(start + (end - start) * easeOutQuart)
    
    callback(current)
    
    if (progress < 1) {
      requestAnimationFrame(animate)
    }
  }
  
  animate()
}

// 加载分类数据
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategoryList()
    if (response.code === 0) {
      categories.value = response.data
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 移除本地的 loadTags 方法，使用 tagStore

// 按分类筛选
const filterByCategory = (categoryId: number) => {
  // TODO: 实现分类筛选
  console.log('筛选分类:', categoryId)
}

// 按标签筛选
const filterByTag = (tagId: number) => {
  // TODO: 实现标签筛选
  console.log('筛选标签:', tagId)
}

// 处理订阅
const handleSubscribe = () => {
  if (!subscribeEmail.value) {
    ElMessage.warning('请输入邮箱地址')
    return
  }
  
  // TODO: 实现订阅功能
  ElMessage.success('订阅成功！')
  subscribeEmail.value = ''
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    userStore.logout()
    ElMessage.success('退出成功')
    router.push('/')
  } catch {
    // 用户取消
  }
}

// 定时器引用
let refreshTimer: NodeJS.Timeout | null = null

onMounted(() => {
  loadArticles()
  loadStats()
  loadCategories()
  tagStore.loadTags()
  
  // 设置定时器，每30秒刷新一次标签数据
  refreshTimer = setInterval(() => {
    tagStore.refreshTags()
  }, 30000)
})

onUnmounted(() => {
  // 清理定时器
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
})
</script>

<style lang="scss" scoped>
.home {
  min-height: 100vh;
  background-color: #f8fafc;
}

.container {
  max-width: 80%;
  margin: 0 auto;
  padding: 0 10px;
}

// 导航栏
.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .logo-icon {
    font-size: 24px;
    color: #409eff;
  }
  
  h1 {
    color: #409eff;
    font-size: 24px;
    font-weight: bold;
    margin: 0;
  }
}

.nav {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  color: #333;
  font-size: 16px;
  transition: color 0.3s ease;
  
  &:hover {
    color: #409eff;
  }
  
  &.router-link-active {
    color: #409eff;
    font-weight: bold;
  }
}

// Hero区域
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  color: white;
  padding: 80px 0;
  text-align: center;
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(45deg, rgba(255,255,255,0.1) 0%, transparent 50%, rgba(255,255,255,0.1) 100%);
    animation: shimmer 3s ease-in-out infinite;
  }
  
  &::after {
    content: '';
    position: absolute;
    top: -50%;
    right: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(255,255,255,0.05) 0%, transparent 70%);
    animation: float 6s ease-in-out infinite;
  }
}

@keyframes shimmer {
  0%, 100% { transform: translateX(-100%); }
  50% { transform: translateX(100%); }
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) rotate(0deg); }
  33% { transform: translate(30px, -30px) rotate(120deg); }
  66% { transform: translate(-20px, 20px) rotate(240deg); }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
}

.hero-title {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 20px;
  line-height: 1.2;
  animation: fadeInUp 1s ease-out;
}

.hero-subtitle {
  font-size: 20px;
  margin-bottom: 40px;
  opacity: 0.9;
  animation: fadeInUp 1s ease-out 0.3s both;
}

.hero-stats {
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: nowrap;
  margin-top: 40px;
  animation: fadeInUp 1s ease-out 0.6s both;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  transition: all 0.3s ease;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 25px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  
  &:hover {
    transform: translateY(-2px);
    background: rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
  }
  
  .el-icon {
    font-size: 20px;
    animation: pulse 2s ease-in-out infinite;
  }
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

// 主要内容
.main {
  padding: 60px 0;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

// 文章区域
.article-section {
  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 30px;
  }
  
  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 24px;
    color: #333;
    margin: 0;
  }
  
  .view-all {
    color: #409eff;
    text-decoration: none;
    font-size: 16px;
    
    &:hover {
      text-decoration: underline;
    }
  }
}

.article-list {
  .loading {
    .el-skeleton {
      margin-bottom: 30px;
    }
  }
  
  .empty {
    text-align: center;
    padding: 60px 0;
  }
}

.article-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 30px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: flex;
  height: 220px; // 从200px增加到220px
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  }
}

.article-image {
  width: 380px; // 从350px增加到380px
  flex-shrink: 0;
  overflow: hidden;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.article-content {
  padding: 24px;
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
  font-size: 21px;
  font-weight: bold;
  color: #333;
  margin-bottom: 16px;
  line-height: 1.3;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  height: 55px; // 增加标题高度
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
          margin-bottom: 18px;
          overflow: hidden;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          font-size: 16px;
          height: 45px; // 增加摘要高度
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
  justify-content: space-between;
}

.author {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 12px;
}

.read-more {
  color: #409eff;
  text-decoration: none;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  
  &:hover {
    text-decoration: underline;
    color: #337ecc;
  }
}

.load-more {
  text-align: center;
  margin-top: 40px;
}

// 分类和标签区域
.categories-tags-section {
  margin-top: 40px;
  
  .categories-tags-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 30px;
  }
  
  .sidebar-card {
    background: #fff;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .sidebar-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
    
    .refresh-btn {
      --el-button-size: 24px;
      --el-button-padding: 4px;
    }
  }
  
  .sidebar-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 18px;
    color: #333;
    margin: 0;
    
    .el-icon {
      color: #409eff;
    }
  }
}

.category-list {
  .category-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 0;
    cursor: pointer;
    color: #666;
    transition: color 0.3s ease;
    border-bottom: 1px solid #f0f0f0;
    
    &:last-child {
      border-bottom: none;
    }
    
    &:hover {
      color: #409eff;
    }
    
    .category-count {
      background: #f0f9ff;
      color: #409eff;
      padding: 2px 8px;
      border-radius: 12px;
      font-size: 12px;
    }
  }
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  
  .tag-item {
    cursor: pointer;
  }
}

// 订阅区域
.subscribe-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 60px 0;
}

.subscribe-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 800px;
  margin: 0 auto;
}

.subscribe-info {
  flex: 1;
  
  h2 {
    font-size: 32px;
    margin-bottom: 16px;
  }
  
  p {
    font-size: 18px;
    margin-bottom: 24px;
    opacity: 0.9;
  }
}

.subscribe-form {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  
  .subscribe-input {
    flex: 1;
  }
}

.privacy-note {
  font-size: 14px;
  opacity: 0.8;
}

.subscribe-icon {
  font-size: 80px;
  opacity: 0.3;
}

// 页脚
.footer {
  background: #1f2937;
  color: white;
  padding: 60px 0 20px;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 40px;
  margin-bottom: 40px;
}

.footer-section {
  h3 {
    font-size: 18px;
    margin-bottom: 20px;
    color: #f9fafb;
  }
  
  p {
    color: #d1d5db;
    line-height: 1.6;
    margin-bottom: 16px;
  }
  
  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    
    li {
      margin-bottom: 8px;
      
      a {
        color: #d1d5db;
        text-decoration: none;
        transition: color 0.3s ease;
        
        &:hover {
          color: #409eff;
        }
      }
    }
  }
}

.social-links {
  display: flex;
  gap: 12px;
  
  .el-button {
    color: #d1d5db;
    
    &:hover {
      color: #409eff;
    }
  }
}

.contact-info {
  p {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }
}

.footer-bottom {
  border-top: 1px solid #374151;
  padding-top: 20px;
  text-align: center;
  
  p {
    color: #9ca3af;
    margin: 0;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .categories-tags-section {
    .categories-tags-grid {
      grid-template-columns: 1fr;
      gap: 20px;
    }
  }
  
  .article-card {
    flex-direction: column;
    height: 260px;
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
  
  .hero-title {
    font-size: 32px;
  }
  
  .hero-stats {
    gap: 10px;
    flex-wrap: wrap;
  }
  
  .stat-item {
    font-size: 14px;
    padding: 8px 12px;
    
    .el-icon {
      font-size: 16px;
    }
  }
  
  .subscribe-content {
    flex-direction: column;
    text-align: center;
    gap: 30px;
  }
  
  .subscribe-form {
    flex-direction: column;
  }
  
  .footer-content {
    grid-template-columns: 1fr;
    gap: 30px;
  }
  
  .nav {
    gap: 10px;
  }
  
  .nav-item {
    font-size: 14px;
  }
}
</style> 

