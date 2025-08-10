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
            <router-link to="/" class="nav-item">首页</router-link>
            <router-link to="/search" class="nav-item">搜索</router-link>
            <template v-if="userStore.userInfo">
              <router-link to="/write" class="nav-item">写文章</router-link>
              <router-link to="/profile" class="nav-item">个人中心</router-link>
              <el-button @click="handleLogout" link>退出</el-button>
            </template>
            <template v-else>
              <router-link to="/login" class="nav-item">登录</router-link>
              <router-link to="/register" class="nav-item">注册</router-link>
            </template>
          </nav>
        </div>
      </div>
    </header>

    <!-- 搜索区域 -->
    <div class="search-section">
      <div class="container">
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章..."
            size="large"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch" :loading="articleStore.loading">
                搜索
              </el-button>
            </template>
          </el-input>
        </div>
      </div>
    </div>

    <!-- 搜索结果 -->
    <main class="main">
      <div class="container">
        <div v-if="searchKeyword && !articleStore.loading" class="search-results">
          <div class="results-header">
            <h2>搜索结果</h2>
            <p>找到 {{ articleStore.total }} 篇文章</p>
          </div>
          
          <div v-if="articleStore.articles.length === 0" class="no-results">
            <el-empty description="没有找到相关文章" />
          </div>
          
          <div v-else class="article-list">
            <article 
              v-for="article in articleStore.articles" 
              :key="article.id" 
              class="article-card"
              @click="goToArticle(article.id)"
            >
              <div class="article-header">
                <h3 class="article-title">{{ article.title }}</h3>
                <div class="article-meta">
                  <span class="author">{{ article.authorName }}</span>
                  <span class="date">{{ formatDate(article.createdAt) }}</span>
                  <span class="category">{{ article.category.name }}</span>
                </div>
              </div>
              <div class="article-content">
                <p>{{ getArticleSummary(article.content) }}</p>
              </div>
              <div class="article-footer">
                <div class="tags">
                  <el-tag 
                    v-for="tag in article.tags" 
                    :key="tag.id" 
                    size="small" 
                    class="tag"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
                <div class="stats">
                  <span class="stat">{{ article.viewCount }} 阅读</span>
                  <span class="stat">{{ article.likeCount }} 点赞</span>
                  <span class="stat">{{ article.commentCount }} 评论</span>
                </div>
              </div>
            </article>
          </div>

          <!-- 分页 -->
          <div v-if="articleStore.total > 0" class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="articleStore.total"
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
                <h3>关键词搜索</h3>
                <p>输入文章标题、内容中的关键词进行搜索</p>
              </div>
              <div class="tip-item">
                <h3>作者搜索</h3>
                <p>输入作者名称查找该作者的文章</p>
              </div>
              <div class="tip-item">
                <h3>标签搜索</h3>
                <p>输入标签名称查找相关文章</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useArticleStore } from '@/stores/article'
import { ElMessage } from 'element-plus'
import { 
  House, Search, Edit, User, UserFilled, Promotion, Document, View, 
  ChatDotRound, Star, Collection, Refresh, Filter
} from '@element-plus/icons-vue'
import type { Article } from '@/types/article'
import { getPlainTextSummary } from '@/utils/markdown'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const articleStore = useArticleStore()

const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 获取文章摘要
const getArticleSummary = (content: string) => {
  return getPlainTextSummary(content, 140)
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
const performSearch = () => {
  if (!searchKeyword.value.trim()) return
  
  articleStore.searchArticles(
    searchKeyword.value,
    currentPage.value,
    pageSize.value
  )
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

// 退出登录
const handleLogout = () => {
  userStore.logout()
  ElMessage.success('退出成功')
  router.push('/')
}

onMounted(() => {
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
  background-color: #f5f5f5;
}

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

.logo-link {
  text-decoration: none;
  
  h1 {
    color: #409eff;
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
  transition: color 0.3s ease;
  
  &:hover {
    color: #409eff;
  }
  
  &.router-link-active {
    color: #409eff;
    font-weight: bold;
  }
}

.search-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px 0;
}

.search-box {
  max-width: 600px;
  margin: 0 auto;
}

.main {
  padding: 20px 0;
}

.search-results {
  .results-header {
    text-align: center;
    margin-bottom: 30px;
    
    h2 {
      font-size: 24px;
      color: #333;
      margin-bottom: 8px;
    }
    
    p {
      color: #666;
      font-size: 16px;
    }
  }
  
  .no-results {
    text-align: center;
    padding: 40px 0;
  }
}

.article-list {
  .article-card {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    }
  }
}

.article-header {
  margin-bottom: 15px;
}

.article-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 10px;
  line-height: 1.4;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  font-size: 14px;
  color: #666;
  
  .author {
    color: #409eff;
  }
}

.article-content {
  margin-bottom: 15px;
  
  p {
    color: #666;
    line-height: 1.6;
  }
}

.article-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tags {
  display: flex;
  gap: 8px;
}

.stats {
  display: flex;
  gap: 15px;
  
  .stat {
    color: #999;
    font-size: 14px;
  }
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

.search-tips {
  .tips-content {
    background: #fff;
    border-radius: 8px;
    padding: 40px;
    text-align: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    
    h2 {
      font-size: 24px;
      color: #333;
      margin-bottom: 30px;
    }
  }
  
  .tips-list {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 30px;
    max-width: 800px;
    margin: 0 auto;
  }
  
  .tip-item {
    text-align: center;
    
    h3 {
      font-size: 18px;
      color: #409eff;
      margin-bottom: 10px;
    }
    
    p {
      color: #666;
      line-height: 1.6;
    }
  }
}

@media (max-width: 768px) {
  .search-section {
    padding: 30px 0;
  }
  
  .search-box {
    padding: 0 20px;
  }
  
  .tips-list {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .article-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style> 