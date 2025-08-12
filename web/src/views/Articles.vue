<template>
  <div class="articles-page">
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
          <h1 class="hero-title">探索所有文章</h1>
          <p class="hero-subtitle">发现更多精彩内容，拓展你的知识视野</p>
          <div class="hero-stats">
            <div class="stat-item">
              <el-icon><Document /></el-icon>
              <span>{{ articleStore.total }} 篇文章</span>
            </div>
            <div class="stat-item">
              <el-icon><View /></el-icon>
              <span>{{ formatNumber(totalViews) }} 总阅读量</span>
            </div>
            <div class="stat-item">
              <el-icon><ChatDotRound /></el-icon>
              <span>{{ formatNumber(totalComments) }} 总评论</span>
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
                <el-icon><Document /></el-icon>
                所有文章
              </h2>
              <div class="section-info">
                <span class="article-count">共 {{ articleStore.total }} 篇文章</span>
              </div>
            </div>
            
            <div class="article-list">
              <div v-if="articleStore.loading" class="loading">
                <el-skeleton :rows="3" animated />
                <el-skeleton :rows="3" animated />
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

            <!-- 分页 -->
            <div v-if="articleStore.total > 0" class="pagination">
              <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :page-sizes="[5, 10, 20, 50]"
                :total="articleStore.total"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-section">
            <h3>我的博客</h3>
            <p>分享技术、生活与思考，记录成长的每一步。在这里，我们一起探索技术的无限可能。</p>
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
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useArticleStore } from '@/stores/article'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  View, Promotion, ChatDotRound, House, Search, Edit, User, UserFilled,
  Document, Collection, Star
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Article } from '@/types/article'
import { getPlainTextSummary } from '@/utils/markdown'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const articleStore = useArticleStore()

const currentPage = ref(1)
const pageSize = ref(5) // 每页显示5篇文章

// 计算总阅读量和总评论数
const totalViews = computed(() => {
  return articleStore.articles.reduce((sum, article) => sum + article.view_count, 0)
})

const totalComments = computed(() => {
  return articleStore.articles.reduce((sum, article) => sum + article.comment_count, 0)
})

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

// 处理页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadArticles()
}

// 处理每页条数变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadArticles()
}

// 加载文章列表
const loadArticles = async () => {
  await articleStore.getArticles(currentPage.value, pageSize.value)
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

// 监听路由参数变化
watch(() => route.query.page, (newPage) => {
  if (newPage) {
    currentPage.value = parseInt(newPage as string) || 1
  }
}, { immediate: true })

onMounted(() => {
  // 从URL参数获取页码
  const pageParam = route.query.page
  if (pageParam) {
    currentPage.value = parseInt(pageParam as string) || 1
  }
  
  loadArticles()
})
</script>

<style lang="scss" scoped>
.articles-page {
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
  
  .section-info {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  
  .article-count {
    color: #666;
    font-size: 16px;
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
  height: 220px;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  }
}

.article-image {
  width: 380px;
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
  overflow: hidden;
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
  line-clamp: 2;
  -webkit-box-orient: vertical;
  height: 55px;
}

.article-summary {
  color: #666;
  line-height: 1.4;
  margin-bottom: 18px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  font-size: 16px;
  height: 45px;
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

// 分页
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

// 页脚
.footer {
  background: #1f2937;
  color: white;
  padding: 40px 0 20px;
  margin-top: 60px;
}

.footer-content {
  margin-bottom: 20px;
}

.footer-section {
  h3 {
    font-size: 18px;
    margin-bottom: 16px;
    color: #f9fafb;
  }
  
  p {
    color: #d1d5db;
    line-height: 1.6;
    margin: 0;
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
  .container {
    max-width: 95%;
    padding: 0 10px;
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
    line-clamp: 2;
  }
  
  .article-summary {
    height: auto;
    -webkit-line-clamp: 2;
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
  
  .nav {
    gap: 10px;
  }
  
  .nav-item {
    font-size: 14px;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style> 