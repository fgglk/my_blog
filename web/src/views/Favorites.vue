<template>
  <div class="favorites-page">
    <div class="container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon class="title-icon"><Star /></el-icon>
          我的收藏
        </h1>
        <p class="page-subtitle">您收藏的所有文章</p>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>

      <!-- 空状态 -->
      <div v-else-if="favorites.length === 0" class="empty-state">
        <el-empty description="暂无收藏文章">
          <el-button type="primary" @click="$router.push('/articles')">
            去发现文章
          </el-button>
        </el-empty>
      </div>

      <!-- 收藏列表 -->
      <div v-else class="favorites-list">
        <div class="list-header">
          <span class="total-count">共 {{ total }} 篇收藏文章</span>
          <div class="sort-controls">
            <el-select v-model="sortBy" placeholder="排序方式" @change="loadFavorites">
              <el-option label="收藏时间" value="created_at" />
              <el-option label="发布时间" value="article_created_at" />
              <el-option label="阅读量" value="view_count" />
            </el-select>
          </div>
        </div>

        <!-- 文章列表 -->
        <div class="articles-container">
          <article 
            v-for="favorite in favorites" 
            :key="favorite.id" 
            class="article-card"
            @click="goToArticle(favorite.article_id)"
          >
            <!-- 文章配图 -->
            <div class="article-image">
              <img :src="getArticleImage(favorite.article)" alt="文章配图" />
            </div>
            
            <!-- 文章内容 -->
            <div class="article-content">
              <!-- 文章元信息 -->
              <div class="article-meta">
                <el-tag size="small" type="primary">{{ getCategoryName(favorite.article.category_id) }}</el-tag>
                <span class="article-date">{{ formatDate(favorite.article.created_at) }}</span>
                <div class="article-stats">
                  <span class="stat-item">
                    <el-icon><View /></el-icon>
                    {{ favorite.article.view_count }} 浏览
                  </span>
                  <span class="stat-item">
                    <el-icon><ChatDotRound /></el-icon>
                    {{ favorite.article.comment_count }} 评论
                  </span>
                  <span class="stat-item">
                    <el-icon><Star /></el-icon>
                    {{ favorite.article.like_count }} 点赞
                  </span>
                  <span class="stat-item">
                    <el-icon><Collection /></el-icon>
                    {{ favorite.article.favorite_count }} 收藏
                  </span>
                </div>
              </div>
              
              <!-- 文章标题 -->
              <h3 class="article-title">{{ favorite.article.title }}</h3>
              
              <!-- 文章摘要 -->
              <p class="article-summary">{{ getArticleSummary(favorite.article.content) }}</p>
              
              <!-- 文章底部 -->
              <div class="article-footer">
                <div class="author">
                  <el-avatar :size="24" :src="favorite.article.author_avatar">
                    {{ getAuthorInitial(favorite.article.author_name) }}
                  </el-avatar>
                  <span>{{ favorite.article.author_name }}</span>
                </div>
                <div class="favorite-info">
                  <span class="favorite-time">
                    <el-icon><Clock /></el-icon>
                    收藏于 {{ formatDate(favorite.created_at) }}
                  </span>
                  <el-button 
                    type="danger" 
                    size="small" 
                    @click.stop="removeFavorite(favorite.id, favorite.article_id)"
                    :loading="removingId === favorite.id"
                    class="unfavorite-btn"
                  >
                    <el-icon><Star /></el-icon>
                    取消收藏
                  </el-button>
                </div>
              </div>
            </div>
          </article>
        </div>

        <!-- 分页 -->
        <div class="pagination-container" v-if="total > pageSize">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Star, View, ChatDotRound, Clock, Collection } from '@element-plus/icons-vue'
import { articleApi } from '@/api/article'
import { categoryApi, type Category } from '@/api/category'
import type { Favorite } from '@/types/api'
import { getPlainTextSummary } from '@/utils/markdown'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const favorites = ref<Favorite[]>([])
const categories = ref<Category[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const sortBy = ref('created_at')
const removingId = ref<number | null>(null)

// 获取分类名称
const getCategoryName = (categoryId: number) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category?.name || '未分类'
}

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN')
}

// 获取作者首字母
const getAuthorInitial = (authorName: string) => {
  if (!authorName) return 'U'
  return authorName.charAt(0).toUpperCase()
}

// 获取文章摘要
const getArticleSummary = (content: string) => {
  if (!content) return '暂无摘要'
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

// 加载收藏列表
const loadFavorites = async () => {
  loading.value = true
  try {
    const response = await articleApi.getFavorites({
      page: currentPage.value,
      size: pageSize.value,
      sort: sortBy.value
    })
    if (response.code === 0) {
      favorites.value = response.data.list
      total.value = response.data.total
    } else {
      ElMessage.error(response.msg || '加载收藏列表失败')
    }
  } catch (error) {
    ElMessage.error('加载收藏列表失败')
    console.error('加载收藏列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategoryList()
    if (response.code === 0) {
      categories.value = response.data
    }
  } catch (error) {
    console.error('加载分类列表失败:', error)
  }
}

// 取消收藏
const removeFavorite = async (favoriteId: number, articleId: number) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏这篇文章吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    removingId.value = favoriteId
    // 使用现有的收藏切换API
    const response = await articleApi.toggleFavorite(articleId)
    if (response.code === 0) {
      ElMessage.success('取消收藏成功')
      // 重新加载列表
      await loadFavorites()
    } else {
      ElMessage.error(response.msg || '取消收藏失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('取消收藏失败')
      console.error('取消收藏失败:', error)
    }
  } finally {
    removingId.value = null
  }
}

// 跳转到文章详情
const goToArticle = (articleId: number) => {
  router.push(`/article/${articleId}`)
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadFavorites()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadFavorites()
}

// 页面初始化
onMounted(() => {
  loadCategories()
  loadFavorites()
})
</script>

<style scoped lang="scss">
.favorites-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 50%, #cbd5e1 100%);
  padding: 40px 0;
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="%23ffffff" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="%23ffffff" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="%23ffffff" opacity="0.1"/><circle cx="10" cy="60" r="0.5" fill="%23ffffff" opacity="0.1"/><circle cx="90" cy="40" r="0.5" fill="%23ffffff" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
    pointer-events: none;
  }
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  position: relative;
  z-index: 1;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
  
  .page-title {
    font-size: 36px;
    font-weight: 700;
    color: #1e293b;
    margin-bottom: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 16px;
    background: linear-gradient(135deg, #1e293b 0%, #3b82f6 50%, #8b5cf6 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    
    .title-icon {
      color: #f59e0b;
      font-size: 32px;
      filter: drop-shadow(0 2px 4px rgba(245, 158, 11, 0.3));
    }
  }
  
  .page-subtitle {
    color: #64748b;
    font-size: 18px;
    font-weight: 400;
    background: rgba(100, 116, 139, 0.1);
    padding: 8px 24px;
    border-radius: 20px;
    display: inline-block;
  }
}

.loading-container {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.empty-state {
  background: white;
  border-radius: 8px;
  padding: 60px 20px;
  text-align: center;
}

.favorites-list {
  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
    padding: 20px 24px;
    background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
    border-radius: 16px;
    border: 1px solid rgba(0, 0, 0, 0.05);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    
    .total-count {
      font-size: 18px;
      color: #1e293b;
      font-weight: 600;
      background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }
}

.articles-container {
  margin-bottom: 30px;
}

.article-card {
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 40px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  height: 260px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  
  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.12);
    border-color: rgba(59, 130, 246, 0.2);
  }
}

.article-image {
  width: 380px;
  flex-shrink: 0;
  overflow: hidden;
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(147, 51, 234, 0.1) 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
    z-index: 1;
  }
  
  &:hover::before {
    opacity: 1;
  }
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.4s ease;
  }
  
  &:hover img {
    transform: scale(1.05);
  }
}

.article-content {
  padding: 28px;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
}

.article-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  
  .article-date {
    color: #94a3b8;
    font-size: 13px;
    font-weight: 500;
    background: rgba(148, 163, 184, 0.1);
    padding: 4px 8px;
    border-radius: 6px;
  }
  
  .article-stats {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .stat-item {
      display: flex;
      align-items: center;
      gap: 4px;
      color: #64748b;
      font-size: 12px;
      font-weight: 500;
      padding: 4px 8px;
      border-radius: 6px;
      background: rgba(100, 116, 139, 0.08);
      transition: all 0.2s ease;
      
      &:hover {
        background: rgba(59, 130, 246, 0.1);
        color: #3b82f6;
      }
      
      .el-icon {
        font-size: 14px;
      }
    }
  }
}

.article-title {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 20px;
  line-height: 1.4;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  height: 68px;
  display: -moz-box;
  -moz-box-orient: vertical;
  display: box;
  box-orient: vertical;
  line-clamp: 2;
  background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.article-summary {
  color: #64748b;
  line-height: 1.6;
  margin-bottom: 24px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  font-size: 16px;
  height: 76px;
  display: -moz-box;
  -moz-box-orient: vertical;
  display: box;
  box-orient: vertical;
  line-clamp: 3;
  font-weight: 400;
}

.article-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.author {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #475569;
  font-size: 13px;
  font-weight: 500;
  padding: 6px 12px;
  border-radius: 8px;
  background: rgba(71, 85, 105, 0.05);
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(59, 130, 246, 0.1);
    color: #3b82f6;
  }
}

.favorite-info {
  display: flex;
  align-items: center;
  gap: 20px;
  
  .favorite-time {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #94a3b8;
    font-size: 13px;
    font-weight: 500;
    padding: 6px 12px;
    border-radius: 8px;
    background: rgba(148, 163, 184, 0.08);
    transition: all 0.2s ease;
    
    &:hover {
      background: rgba(34, 197, 94, 0.1);
      color: #22c55e;
    }
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

@media (max-width: 768px) {
  .favorites-page {
    padding: 20px 0;
  }
  
  .container {
    padding: 0 16px;
  }
  
  .page-header {
    margin-bottom: 30px;
    
    .page-title {
      font-size: 28px;
      gap: 12px;
      
      .title-icon {
        font-size: 24px;
      }
    }
    
    .page-subtitle {
      font-size: 16px;
      padding: 6px 20px;
    }
  }
  
  .article-card {
    flex-direction: column;
    height: auto;
    margin-bottom: 30px;
    
    .article-image {
      width: 100%;
      height: 220px;
    }
    
    .article-content {
      padding: 24px;
    }
    
    .article-title {
      font-size: 20px;
      height: 56px;
      margin-bottom: 16px;
    }
    
    .article-summary {
      font-size: 15px;
      height: 72px;
      margin-bottom: 20px;
    }
  }
  
  .list-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
    padding: 16px 20px;
    
    .total-count {
      font-size: 16px;
    }
  }
  
  .article-meta {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
    
    .article-stats {
      gap: 12px;
      
      .stat-item {
        font-size: 11px;
        padding: 3px 6px;
      }
    }
  }
  
  .article-footer {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
    padding-top: 12px;
  }
  
  .favorite-info {
    width: 100%;
    justify-content: space-between;
    gap: 12px;
    
    .favorite-time {
      font-size: 12px;
      padding: 4px 8px;
    }
  }
}

.unfavorite-btn {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  border: none;
  border-radius: 8px;
  padding: 8px 16px;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.2);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
    background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
  }
  
  .el-icon {
    margin-right: 4px;
  }
}
</style>
