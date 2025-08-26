<template>
  <div class="search-page">
    <!-- 导航栏 -->
    <header class="header">
      <div class="container">
        <div class="header-content">
          <div class="logo">
            <router-link to="/" class="logo-link">
              <h1>我的博客</h1>
            </router-link>
          </div>
          <nav class="nav">
            <router-link to="/" class="nav-item">
              <el-icon><House /></el-icon>
              首页
            </router-link>
            <router-link to="/search" class="nav-item active">
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

    <!-- 搜索区域 -->
    <div class="search-section">
      <div class="container">
        <div class="search-box">
          <h1 class="search-title">搜索文章</h1>
          <p class="search-subtitle">发现精彩内容，探索知识世界</p>
          <div class="search-input-wrapper">
            <el-input
              v-model="searchKeyword"
              placeholder="输入关键词搜索文章..."
              size="large"
              clearable
              @keyup.enter="handleSearch"
              class="search-input"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
              <template #append>
                <el-button @click="handleSearch" :loading="loading" type="primary">
                  搜索
                </el-button>
              </template>
            </el-input>
          </div>
          
          <!-- 搜索过滤器 -->
          <div class="search-filters">
            <div class="filter-group">
              <el-select v-model="selectedCategory" placeholder="选择分类" clearable @change="handleSearch">
                <el-option
                  v-for="category in categories"
                  :key="category.id"
                  :label="category.name"
                  :value="category.id"
                />
              </el-select>
            </div>
            <div class="filter-group">
              <el-select v-model="sortBy" placeholder="排序方式" @change="handleSearch">
                <el-option label="最新发布" value="time" />
                <el-option label="最多阅读" value="view" />
                <el-option label="最多评论" value="comment" />
                <el-option label="最多点赞" value="like" />
              </el-select>
            </div>
            <div class="filter-group">
              <el-select v-model="sortOrder" placeholder="排序顺序" @change="handleSearch">
                <el-option label="降序" value="desc" />
                <el-option label="升序" value="asc" />
              </el-select>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 搜索结果 -->
    <main class="main">
      <div class="container">
        <div v-if="searchKeyword && !loading" class="search-results">
          <div class="results-header">
            <h2>搜索结果</h2>
            <p>找到 {{ searchResult.total }} 篇文章</p>
          </div>
          
          <div v-if="searchResult.articles.length === 0" class="no-results">
            <el-empty description="没有找到相关文章">
              <template #description>
                <p>没有找到与 "{{ searchKeyword }}" 相关的文章</p>
                <p>试试其他关键词或调整搜索条件</p>
              </template>
            </el-empty>
          </div>
          
          <div v-else class="article-list">
            <article 
              v-for="article in searchResult.articles" 
              :key="article.id" 
              class="article-card"
              @click="goToArticle(article.id)"
            >
              <div class="article-image">
                <img :src="getArticleImage(article)" alt="文章配图" />
              </div>
              <div class="article-content">
                <div class="article-meta">
                  <el-tag size="small" type="primary">{{ getCategoryName(article.category_id) }}</el-tag>
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
                      {{ article.favorite_count || 0 }} 收藏
                    </span>
                  </div>
                </div>
                <h3 class="article-title">{{ article.title }}</h3>
                <p class="article-summary">{{ getArticleSummary(article.content || article.summary) }}</p>
                <div class="article-footer">
                  <div class="author">
                    <el-avatar :size="24" :src="article.author_avatar">
                      {{ getAuthorInitial(article.author_name || article.author_nickname) }}
                    </el-avatar>
                    <span>{{ article.author_nickname || article.author_name }}</span>
                  </div>
                  <span class="read-more" @click.stop="goToArticle(article.id)">
                    阅读全文 →
                  </span>
                </div>
              </div>
            </article>
          </div>

          <!-- 分页 -->
          <div v-if="searchResult.total > 0" class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="searchResult.total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
        
        <div v-else-if="!searchKeyword" class="search-tips">
          <div class="tips-content">
            <h2>搜索提示</h2>
            <div class="tips-list">
              <div class="tip-item">
                <div class="tip-icon">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="tip-content">
                  <h3>关键词搜索</h3>
                  <p>输入文章标题、内容中的关键词进行搜索</p>
                </div>
              </div>
              <div class="tip-item">
                <div class="tip-icon">
                  <el-icon><User /></el-icon>
                </div>
                <div class="tip-content">
                  <h3>作者搜索</h3>
                  <p>输入作者名称查找该作者的文章</p>
                </div>
              </div>
              <div class="tip-item">
                <div class="tip-icon">
                  <el-icon><Collection /></el-icon>
                </div>
                <div class="tip-content">
                  <h3>标签搜索</h3>
                  <p>输入标签名称查找相关文章</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { 
  Search, Document, User, UserFilled, Collection, View, Star, ChatDotRound, House, Edit
} from '@element-plus/icons-vue'
import { getPlainTextSummary } from '@/utils/markdown'
import dayjs from 'dayjs'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 搜索相关状态
const searchKeyword = ref('')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const selectedCategory = ref<number | null>(null)
const sortBy = ref('time')
const sortOrder = ref('desc')

// 分类列表
const categories = ref<any[]>([])

// 搜索结果
const searchResult = reactive({
  articles: [] as any[],
  total: 0,
  page: 1,
  size: 10,
  totalPage: 0
})

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

// 获取文章摘要
const getArticleSummary = (content: string) => {
  return getPlainTextSummary(content, 140)
}

// 获取文章配图
const getArticleImage = (article: any) => {
  // 如果文章有封面图片，使用封面图片
  if (article.cover_image) {
    return article.cover_image
  }
  // 否则使用占位图片
  return `https://picsum.photos/300/200?random=${article.id}`
}

// 获取分类名称
const getCategoryName = (categoryId: number) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : '未分类'
}

// 获取作者首字母
const getAuthorInitial = (authorName: string) => {
  if (!authorName) return 'U'
  return authorName.charAt(0).toUpperCase()
}

// 跳转到文章详情
const goToArticle = (id: number) => {
  router.push(`/article/${id}`)
}

// 处理搜索
const handleSearch = () => {
  if (!searchKeyword.value.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }
  
  currentPage.value = 1
  performSearch()
}

// 执行搜索
const performSearch = async () => {
  if (!searchKeyword.value.trim()) return
  
  loading.value = true
  try {
    const params: any = {
      keyword: searchKeyword.value,
      page: currentPage.value,
      size: pageSize.value,
      sort: sortBy.value,
      order: sortOrder.value
    }
    
    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }
    
    // 添加调试日志
    console.log('搜索参数:', params)
    
    const response = await request.get('/articles/search', { params })
    
    searchResult.articles = response.data.articles || []
    searchResult.total = response.data.total || 0
    searchResult.page = response.data.page || 1
    searchResult.size = response.data.size || 10
    searchResult.totalPage = response.data.total_page || 0
  } catch (error) {
    console.error('搜索失败:', error)
    ElMessage.error('搜索失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  performSearch()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  performSearch()
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const response = await request.get('/categories')
    categories.value = response.data || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 退出登录
const handleLogout = () => {
  userStore.logout()
  ElMessage.success('退出成功')
  router.push('/')
}

onMounted(async () => {
  await loadCategories()
  
  // 如果URL中有搜索参数，自动执行搜索
  const keyword = route.query.keyword as string
  if (keyword) {
    searchKeyword.value = keyword
    handleSearch()
  }
})
</script>

<style lang="scss" scoped>
.search-page {
  min-height: 100vh;
  background: #f8fafc;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 0;
}

.logo-link {
  text-decoration: none;
  color: inherit;
  
  h1 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
  }
}

.nav {
  display: flex;
  align-items: center;
  gap: 2rem;
  
  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    text-decoration: none;
    color: #6b7280;
    font-weight: 500;
    transition: color 0.2s;
    
    &:hover {
      color: #3b82f6;
    }
  }
}

.search-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60px 0;
  text-align: center;
}

.search-box {
  max-width: 800px;
  margin: 0 auto;
}

.search-title {
  color: #fff;
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 10px;
}

.search-subtitle {
  color: rgba(255, 255, 255, 0.8);
  font-size: 18px;
  margin-bottom: 30px;
}

.search-input-wrapper {
  margin-bottom: 20px;
}

.search-input {
  :deep(.el-input__wrapper) {
    border-radius: 50px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  }
  
  :deep(.el-input__inner) {
    font-size: 16px;
    padding: 15px 20px;
  }
  
  :deep(.el-button) {
    border-radius: 0 50px 50px 0;
    padding: 15px 30px;
    font-size: 16px;
  }
}

.search-filters {
  display: flex;
  gap: 15px;
  justify-content: center;
  flex-wrap: wrap;
  
  .filter-group {
    :deep(.el-select) {
      .el-input__wrapper {
        border-radius: 25px;
        background: rgba(255, 255, 255, 0.9);
        backdrop-filter: blur(10px);
      }
    }
  }
}

.main {
  padding: 2rem 0;
}

.search-results {
  .results-header {
    text-align: center;
    margin-bottom: 2rem;
    
    h2 {
      font-size: 2rem;
      font-weight: 700;
      margin: 0 0 0.5rem 0;
      color: #1f2937;
    }
    
    p {
      color: #6b7280;
      margin: 0;
    }
  }
  
  .no-results {
    text-align: center;
    padding: 60px 0;
  }
}

.article-list {
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
    -webkit-box-orient: vertical;
    height: 55px;
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
    height: 45px;
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
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.search-tips {
  .tips-content {
    background: #fff;
    border-radius: 0.75rem;
    padding: 3rem 2rem;
    text-align: center;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    
    h2 {
      font-size: 2rem;
      font-weight: 700;
      margin: 0 0 2rem 0;
      color: #1f2937;
    }
  }
  
  .tips-list {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
    max-width: 900px;
    margin: 0 auto;
  }
  
  .tip-item {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    padding: 1.5rem;
    border-radius: 0.75rem;
    background: #fff;
    border: 1px solid #e5e7eb;
    transition: all 0.2s;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
    
    .tip-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 3rem;
      height: 3rem;
      border-radius: 50%;
      background: #3b82f6;
      color: #fff;
      font-size: 1.25rem;
    }
    
    .tip-content {
      text-align: left;
      
      h3 {
        font-size: 1.125rem;
        color: #1f2937;
        margin-bottom: 0.5rem;
        font-weight: 600;
      }
      
      p {
        color: #6b7280;
        line-height: 1.6;
        margin: 0;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .search-section {
    padding: 40px 0;
  }
  
  .search-title {
    font-size: 28px;
  }
  
  .search-subtitle {
    font-size: 16px;
  }
  
  .search-filters {
    flex-direction: column;
    align-items: center;
  }
  
  .article-card {
    flex-direction: column !important;
    height: 260px !important;
  }
  
  .article-image {
    width: 100% !important;
    height: 120px !important;
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
  
  .tips-list {
    grid-template-columns: 1fr !important;
    gap: 1rem !important;
  }
  
  .tip-item {
    flex-direction: column;
    text-align: center;
    
    .tip-content {
      text-align: center !important;
    }
  }
}
</style> 