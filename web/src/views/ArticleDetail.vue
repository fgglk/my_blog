<template>
  <div class="article-detail">
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

    <!-- 主要内容 -->
    <main class="main">
      <div class="container">
        <div v-if="articleStore.loading" class="loading">
          <el-skeleton :rows="10" animated />
        </div>
        
        <div v-else-if="!articleStore.currentArticle" class="empty">
          <el-empty description="文章不存在" />
        </div>
        
        <div v-else class="content">
          <!-- 登录提示 -->
          <div v-if="requireLoginToView" class="login-required">
            <el-empty description="此文章需要登录后才能查看">
              <el-button type="primary" @click="$router.push('/login')" size="large">
                <el-icon><User /></el-icon>
                立即登录
              </el-button>
            </el-empty>
          </div>
          
          <!-- 文章内容区域 -->
          <div v-else class="article-container">
            <!-- 文章主体 -->
            <article class="article">
              <div class="article-header">
                <!-- 分类和阅读时间 -->
                <div class="article-meta-top">
                  <div class="category-info">
                    <el-icon class="category-icon"><Document /></el-icon>
                    <span class="category-name">{{ articleStore.currentArticle.category.name }}</span>
                  </div>
                  <div class="read-time">
                    <el-icon class="time-icon"><Clock /></el-icon>
                    <span>{{ getReadTime(articleStore.currentArticle.content) }}分钟阅读</span>
                  </div>
                </div>
                
                <!-- 文章标题 -->
                <h1 class="article-title">{{ articleStore.currentArticle.title }}</h1>
                
                <!-- 作者信息 -->
                <div v-if="showAuthorSection" class="author-section">
                  <div class="author-info" @click="goToAuthorProfile" style="cursor: pointer;">
                    <el-avatar :size="40" class="author-avatar" :src="articleStore.currentArticle.author_avatar">
                      {{ articleStore.currentArticle.author_name.charAt(0) }}
                    </el-avatar>
                    <div class="author-details">
                      <span class="author-name">{{ articleStore.currentArticle.author_name }}</span>
                      <el-tag size="small" type="primary" class="author-tag">作者</el-tag>
                    </div>
                    <span class="publish-date">发布于 {{ formatDate(articleStore.currentArticle.created_at) }}</span>
                    <el-icon class="click-hint"><ArrowRight /></el-icon>
                  </div>
                </div>
                
                <!-- 文章标签 -->
                <div class="article-tags">
                  <el-tag 
                    v-for="tag in articleStore.currentArticle.tags" 
                    :key="tag.id"
                    size="small"
                    class="tag"
                    type="info"
                    effect="light"
                  >
                    #{{ tag.name }}
                  </el-tag>
                </div>
              </div>
              
              <!-- 登录检查 -->
              <div v-if="requireLoginToView" class="login-required">
                <div class="login-required-content">
                  <el-icon class="login-icon"><Lock /></el-icon>
                  <h3>需要登录才能查看</h3>
                  <p>此文章需要登录后才能查看完整内容</p>
                  <el-button type="primary" @click="goToLogin" size="large">
                    立即登录
                  </el-button>
                </div>
              </div>
              
              <div v-else class="article-content">
                <div class="markdown-content" v-html="renderedContent"></div>
              </div>
              
              <div class="article-actions">
                <div class="action-buttons">
                  <!-- 编辑按钮（仅作者或管理员可见） -->
                  <el-button 
                    v-if="canEditArticle"
                    type="success"
                    @click="goToEdit"
                    size="large"
                    class="action-btn edit-btn"
                  >
                    <el-icon><Edit /></el-icon>
                    编辑文章
                  </el-button>
                  
                  <el-button 
                    :type="articleStore.currentArticle.is_liked ? 'primary' : 'default'"
                    @click="handleLike"
                    :loading="likeLoading"
                    size="large"
                    class="action-btn like-btn"
                  >
                    <el-icon><Star /></el-icon>
                    {{ articleStore.currentArticle.is_liked ? '已点赞' : '点赞' }}
                    ({{ articleStore.currentArticle.like_count }})
                  </el-button>
                  
                  <el-button 
                    :type="articleStore.currentArticle.is_favorited ? 'warning' : 'default'"
                    @click="handleFavorite"
                    :loading="favoriteLoading"
                    size="large"
                    class="action-btn favorite-btn"
                  >
                    <el-icon><Collection /></el-icon>
                    {{ articleStore.currentArticle.is_favorited ? '已收藏' : '收藏' }}
                    ({{ articleStore.currentArticle.favorite_count }})
                  </el-button>

                  <el-button 
                    v-if="showShareButtons"
                    type="info"
                    @click="handleShare"
                    size="large"
                    class="action-btn share-btn"
                  >
                    <el-icon><Link /></el-icon>
                    分享
                  </el-button>
                </div>
              </div>
            </article>

            <!-- 侧边栏 -->
            <aside class="sidebar">
              <!-- 文章数据卡片 -->
              <div class="data-card">
                <h3 class="card-title">文章数据</h3>
                <div class="stats-grid">
                  <div class="stat-box">
                    <span class="stat-number primary">{{ articleStore.currentArticle.view_count }}</span>
                    <span class="stat-label">阅读数</span>
                  </div>
                  <div class="stat-box">
                    <span class="stat-number">{{ articleStore.currentArticle.comment_count }}</span>
                    <span class="stat-label">评论数</span>
                  </div>
                  <div class="stat-box">
                    <span class="stat-number">{{ articleStore.currentArticle.like_count }}</span>
                    <span class="stat-label">点赞数</span>
                  </div>
                  <div class="stat-box">
                    <span class="stat-number">{{ articleStore.currentArticle.favorite_count }}</span>
                    <span class="stat-label">收藏数</span>
                  </div>
                </div>
                <div class="article-info">
                  <div class="info-item">
                    <span class="info-label">发布时间</span>
                    <span class="info-value">{{ formatDateShort(articleStore.currentArticle.created_at) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="info-label">最后更新</span>
                    <span class="info-value">{{ formatDateShort(articleStore.currentArticle.updated_at) }}</span>
                  </div>
                  <div class="info-item">
                    <span class="info-label">文章字数</span>
                    <span class="info-value">{{ getWordCount(articleStore.currentArticle.content) }}字</span>
                  </div>
                </div>
              </div>

                      <!-- 目录导航 -->
        <div v-if="showTOC" class="toc-card">
          <h3 class="card-title">
            <el-icon><Document /></el-icon>
            文章目录
          </h3>
          <div class="toc-content">
            <div
              class="toc-item"
              v-for="(heading, index) in tocItems"
              :key="index"
              :class="`toc-level-${heading.level}`"
            >
              <a
                :href="`#${heading.id}`"
                class="toc-link"
                :class="{
                  active: activeHeading === heading.id,
                  [`level-${heading.level}`]: true
                }"
                @click.prevent="scrollToHeading(heading.id)"
              >
                <span class="toc-bullet">{{ '•'.repeat(heading.level) }}</span>
                <span class="toc-text">{{ heading.text }}</span>
              </a>
            </div>
          </div>
        </div>

        <!-- 移动端目录导航 (正常流布局) -->
        <div v-if="showTOC" class="toc-card mobile-toc">
          <h3 class="card-title">
            <el-icon><Document /></el-icon>
            文章目录
          </h3>
          <div class="toc-content">
            <div
              class="toc-item"
              v-for="(heading, index) in tocItems"
              :key="index"
              :class="`toc-level-${heading.level}`"
            >
              <a
                :href="`#${heading.id}`"
                class="toc-link"
                :class="{
                  active: activeHeading === heading.id,
                  [`level-${heading.level}`]: true
                }"
                @click.prevent="scrollToHeading(heading.id)"
              >
                <span class="toc-bullet">{{ '•'.repeat(heading.level) }}</span>
                <span class="toc-text">{{ heading.text }}</span>
              </a>
            </div>
          </div>
        </div>
            </aside>
          </div>

          <!-- 相关文章 -->
          <div class="related-articles-section" :class="{ 'empty-section': relatedArticles.length === 0 }">
            <div class="section-header">
              <h3 class="section-title">
                <el-icon><Link /></el-icon>
                相关文章
              </h3>
            </div>
            <div class="related-articles-list">
              <div v-if="relatedArticles.length === 0" class="empty">
                <div class="empty-content">
                  <el-icon class="empty-icon"><Box /></el-icon>
                  <p class="empty-text">暂无相关文章</p>
                </div>
              </div>
              <div v-else>
                <article 
                  v-for="article in relatedArticles" 
                  :key="article.id" 
                  class="related-article-card"
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
          </div>

          <!-- 评论区 -->
          <div v-if="showCommentsSection" class="comments-section">
            <div class="comments-header">
              <h3 class="comments-title">
                <el-icon><ChatDotRound /></el-icon>
                评论 ({{ articleStore.currentArticle.comment_count }})
              </h3>
            </div>
            
            <!-- 发表评论 -->
            <div v-if="articleReadingSettings.allowComments && userStore.userInfo" class="comment-form">
              <div class="comment-input-wrapper">
                <el-avatar :size="48" class="comment-avatar">
                  {{ userStore.userInfo?.username.charAt(0) || 'U' }}
                </el-avatar>
                <div class="comment-input-content">
                  <el-input
                    v-model="commentContent"
                    type="textarea"
                    :rows="3"
                    placeholder="写下你的评论..."
                    maxlength="500"
                    show-word-limit
                    class="comment-textarea"
                  />
                  <div class="comment-actions">
                    <el-button 
                      type="primary" 
                      @click="submitComment"
                      :loading="commentLoading"
                      :disabled="!commentContent.trim()"
                      size="small"
                      class="submit-btn"
                    >
                      发表评论
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
            
            <div v-else-if="articleReadingSettings.allowComments && !userStore.userInfo" class="login-prompt">
              <el-button type="primary" @click="$router.push('/login')" size="large" class="login-btn">
                <el-icon><User /></el-icon>
                登录后发表评论
              </el-button>
            </div>
            
            <div v-else-if="!articleReadingSettings.allowComments" class="comments-disabled">
              <el-empty description="评论功能已关闭">
                <p class="disabled-message">作者已关闭此文章的评论功能</p>
              </el-empty>
            </div>
            
            <!-- 评论列表 -->
            <div class="comments-list">
              <div v-if="comments.length === 0" class="no-comments">
                <el-empty description="暂无评论，快来抢沙发吧！" />
              </div>
              <div v-else>
                <CommentItem 
                  v-for="comment in comments" 
                  :key="comment.id"
                  :comment="comment"
                  :user-store="userStore"
                  :comment-api="commentApi"
                  :article-id="parseInt(route.params.id as string)"
                  :article-author-id="articleStore.currentArticle?.author_id"
                  @comment-updated="handleCommentUpdated"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useArticleStore } from '@/stores/article'
import { commentApi } from '@/api/comment'
import { articleApi } from '@/api/article'
import { ElMessage } from 'element-plus'
import { 
  View, ChatDotRound, Star, Collection, Document, Link, User, Clock, Box, Edit, Lock, ArrowRight
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'
import type { Comment } from '@/types/comment'
import type { Article, ReadingSettings } from '@/types/article'
import CommentItem from '@/components/CommentItem.vue'
import { getPlainTextSummary } from '@/utils/markdown'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const articleStore = useArticleStore()

const commentContent = ref('')
const commentLoading = ref(false)
const likeLoading = ref(false)
const favoriteLoading = ref(false)
const comments = ref<Comment[]>([])
const relatedArticles = ref<Article[]>([])
const activeHeading = ref('')

// Markdown渲染器
const md: MarkdownIt = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true, // 支持换行
  highlight: function (str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        const highlighted = hljs.highlight(str, { language: lang }).value;
        return `<pre class="hljs"><code class="language-${lang}">${highlighted}</code></pre>`;
      } catch (__) {}
    }
    // 如果没有语言或高亮失败，使用默认的代码块
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`;
  }
})

// 渲染文章内容
const renderedContent = computed(() => {
  if (!articleStore.currentArticle) return ''
  let html = md.render(articleStore.currentArticle.content)
  
  // 为标题添加ID属性
  html = html.replace(
    /<h([1-6])>(.*?)<\/h[1-6]>/g,
    (_match: string, level: string, content: string) => {
      // 生成标题ID
      const text = content.replace(/<[^>]*>/g, '') // 移除HTML标签
      const id = text
        .toLowerCase()
        .replace(/[^\w\u4e00-\u9fa5]+/g, '-') // 替换非字母数字字符为连字符
        .replace(/^-+|-+$/g, '') // 移除首尾连字符
        .substring(0, 50) // 限制长度
      
      return `<h${level} id="${id}">${content}</h${level}>`
    }
  )
  
  // 处理代码块，添加语言标识和行号
  html = html.replace(
    /<pre class="hljs"><code class="language-(\w+)">([\s\S]*?)<\/code><\/pre>/g,
    (_match: string, lang: string, code: string) => {
      const lines: string[] = code.split('\n')
      
      const codeWithLineNumbers = lines.map((line: string, index: number) => {
        return `<div class="code-line">
          <span class="line-number">${index + 1}</span>
          <span class="line-content">${line || ' '}</span>
        </div>`
      }).join('\n')
      
      const result = `
        <div class="code-block" data-language="${lang}">
          <div class="code-header">
            <span class="language-label">${lang}</span>
            <button class="copy-btn" onclick="copyCode(this)">复制代码</button>
          </div>
          <div class="code-content">
            ${codeWithLineNumbers}
          </div>
        </div>
      `
      return result
    }
  )
  
  // 处理没有语言标识的代码块
  html = html.replace(
    /<pre class="hljs"><code>([\s\S]*?)<\/code><\/pre>/g,
    (_match: string, code: string) => {
      const lines: string[] = code.split('\n')
      
      const codeWithLineNumbers = lines.map((line: string, index: number) => {
        return `<div class="code-line">
          <span class="line-number">${index + 1}</span>
          <span class="line-content">${line || ' '}</span>
        </div>`
      }).join('\n')
      
      return `
        <div class="code-block" data-language="text">
          <div class="code-header">
            <span class="language-label">text</span>
            <button class="copy-btn" onclick="copyCode(this)">复制代码</button>
          </div>
          <div class="code-content">
            ${codeWithLineNumbers}
          </div>
        </div>
      `
    }
  )
  
  return html
})

// 修复图片样式的方法
const fixImageStyles = () => {
  // 使用nextTick确保DOM已更新
  nextTick(() => {
    // 只选择文章内容中的图片，排除头像
    const images = document.querySelectorAll('.markdown-content img:not(.el-avatar img)')
    images.forEach((img: Element) => {
      const htmlImg = img as HTMLImageElement
      // 检查是否是头像组件内的图片
      if (!htmlImg.closest('.el-avatar')) {
        htmlImg.style.maxWidth = 'calc(100% - 40px)'
        htmlImg.style.margin = '20px 20px'
        htmlImg.style.display = 'block'
        htmlImg.style.borderRadius = '12px'
        htmlImg.style.boxShadow = '0 8px 25px rgba(0, 0, 0, 0.15)'
        htmlImg.style.transition = 'all 0.3s ease'
        htmlImg.style.objectFit = 'contain'
        htmlImg.style.height = 'auto'
        htmlImg.style.width = 'auto'
      }
    })
  })
}

// 解析阅读设置
const parseReadingSettings = (summary: string): ReadingSettings => {
  const defaultSettings: ReadingSettings = {
    allowComments: true,
    allowRepost: true,
    requireLogin: false,
    showAuthorInfo: true,
    enableTOC: true
  }
  
  if (!summary) {
    return defaultSettings
  }
  
  const settingsMatch = summary.match(/<!--READ_SETTINGS:({.*?})-->/)
  
  if (settingsMatch) {
    try {
      const settings = JSON.parse(settingsMatch[1])
      return {
        allowComments: settings.allowComments ?? true,
        allowRepost: settings.allowRepost ?? true,
        requireLogin: settings.requireLogin ?? false,
        showAuthorInfo: settings.showAuthorInfo ?? true,
        enableTOC: settings.enableTOC ?? true
      }
    } catch (error) {
      console.error('解析阅读设置失败:', error)
      return defaultSettings
    }
  }
  
  return defaultSettings
}

// 文章阅读设置
const articleReadingSettings = computed(() => {
  if (!articleStore.currentArticle?.summary) {
    return {
      allowComments: true,
      allowRepost: true,
      requireLogin: false,
      showAuthorInfo: true,
      enableTOC: true
    }
  }
  return parseReadingSettings(articleStore.currentArticle.summary)
})

// 是否显示评论区
const showCommentsSection = computed(() => {
  return articleReadingSettings.value.allowComments
})



// 是否需要登录才能查看
const requireLoginToView = computed(() => {
  return articleReadingSettings.value.requireLogin && !userStore.userInfo
})

// 是否显示分享按钮
const showShareButtons = computed(() => {
  return articleReadingSettings.value.allowRepost
})

// 是否显示作者信息
const showAuthorSection = computed(() => {
  return articleReadingSettings.value.showAuthorInfo
})

// 是否显示目录
const showTOC = computed(() => {
  return articleReadingSettings.value.enableTOC && tocItems.value.length > 0
})

// 生成目录
const tocItems = computed(() => {
  if (!articleStore.currentArticle) return []
  
  const headings: { id: string; text: string; level: number }[] = []
  const content = articleStore.currentArticle.content
  
  // 使用更健壮的方法：先移除所有代码块，再提取标题
  let processedContent = content
  
  // 移除所有代码块（包括 ``` 和 ~~~ 格式）
  processedContent = processedContent.replace(/```[\s\S]*?```/g, '')
  processedContent = processedContent.replace(/~~~[\s\S]*?~~~/g, '')
  
  // 移除缩进代码块（以4个空格或1个制表符开头的连续行）
  const lines = processedContent.split('\n')
  const filteredLines: string[] = []
  let inIndentedBlock = false
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const isIndented = line.match(/^( {4}|\t)/)
    
    if (isIndented) {
      if (!inIndentedBlock) {
        inIndentedBlock = true
      }
      // 跳过缩进行
      continue
    } else {
      if (inIndentedBlock) {
        inIndentedBlock = false
      }
      filteredLines.push(line)
    }
  }
  
  // 从过滤后的内容中提取标题
  filteredLines.forEach((line) => {
    const match = line.match(/^(#{1,6})\s+(.+)$/)
    if (match) {
      const level = match[1].length
      const text = match[2].trim()
      
      // 跳过空标题
      if (!text) {
        return
      }
      
      // 使用与标题ID相同的生成规则
      const id = text
        .toLowerCase()
        .replace(/[^\w\u4e00-\u9fa5]+/g, '-') // 替换非字母数字字符为连字符
        .replace(/^-+|-+$/g, '') // 移除首尾连字符
        .substring(0, 50) // 限制长度
      
      headings.push({ id, text, level })
    }
  })
  
  return headings
})

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 格式化日期（短）
const formatDateShort = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

// 获取阅读时间
const getReadTime = (content: string) => {
  const words = content.split(/\s+/).filter(word => word.length > 0)
  const averageReadingSpeed = 200 // 平均每分钟阅读200字
  const readTime = Math.ceil(words.length / averageReadingSpeed)
  return readTime > 0 ? readTime : 1
}

// 获取文章字数
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

// 加载文章详情
const loadArticle = async () => {
  const articleId = parseInt(route.params.id as string)
  if (isNaN(articleId)) {
    ElMessage.error('文章ID无效')
    router.push('/')
    return
  }
  
  await articleStore.getArticle(articleId)
  await loadComments(articleId)
  await loadRelatedArticles(articleId)
  
  // 文章加载完成后修复图片样式
  fixImageStyles()
}

// 跳转到作者个人中心
const goToAuthorProfile = () => {
  if (articleStore.currentArticle?.author_id) {
    router.push(`/profile/${articleStore.currentArticle.author_id}`)
  } else {
    ElMessage.warning('无法获取作者信息')
  }
}

// 加载相关文章
const loadRelatedArticles = async (articleId: number) => {
  try {
    const response = await articleApi.getRelatedArticles(articleId)
    if (response.code === 0) {
      relatedArticles.value = response.data.list || []
    } else {
      console.error('加载相关文章失败:', response.msg)
      relatedArticles.value = []
    }
  } catch (error) {
    console.error('加载相关文章失败:', error)
    relatedArticles.value = []
  }
}

// 加载评论
const loadComments = async (articleId: number) => {
  try {
    const response = await commentApi.getCommentList(articleId, { page: 1, size: 50 })
    if (response.code === 0) {
      comments.value = response.data.list || []
    } else {
      console.error('评论加载失败:', response.msg)
      comments.value = []
    }
  } catch (error) {
    console.error('加载评论失败:', error)
    comments.value = []
  }
}

// 处理评论更新（包括删除、创建等操作后）
const handleCommentUpdated = async () => {
  const articleId = parseInt(route.params.id as string)
  await loadComments(articleId)
  // 重新加载文章信息以更新评论数
  await articleStore.getArticle(articleId)
}

// 提交评论
const submitComment = async () => {
  // 检查用户是否已登录
  if (!userStore.userInfo) {
    if (userStore.getCurrentToken()) {
      try {
        await userStore.getUserInfo()
        if (!userStore.userInfo) {
          ElMessage.warning('登录已过期，请重新登录')
          router.push('/login')
          return
        }
      } catch (error) {
        ElMessage.warning('登录已过期，请重新登录')
        router.push('/login')
        return
      }
    } else {
      ElMessage.warning('请先登录')
      router.push('/login')
      return
    }
  }
  
  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }
  
  commentLoading.value = true
  try {
    const articleId = parseInt(route.params.id as string)
    
    const commentData = {
      content: commentContent.value,
      article_id: articleId
    }
    
    // 检查token是否存在
    const currentToken = userStore.getCurrentToken()
    if (!currentToken) {
      ElMessage.error('登录状态异常，请重新登录')
      router.push('/login')
      return
    }
    
    const response = await commentApi.createComment(commentData)
    
    if (response.code === 0) {
      ElMessage.success('评论发表成功')
      commentContent.value = ''
      await handleCommentUpdated()
    } else {
      ElMessage.error(response.msg || '评论发表失败')
    }
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(`评论发表失败: ${error.message}`)
    } else {
      ElMessage.error('评论发表失败')
    }
  } finally {
    commentLoading.value = false
  }
}

// 点赞
const handleLike = async () => {
  if (!userStore.userInfo) {
    ElMessage.warning('请先登录')
    return
  }
  
  likeLoading.value = true
  try {
    const articleId = parseInt(route.params.id as string)
    const result = await articleStore.toggleLike(articleId)
    if (result.success) {
      ElMessage.success('操作成功')
    } else {
      ElMessage.error(result.message || '操作失败')
    }
  } catch (error) {
    console.error('点赞失败:', error)
    ElMessage.error('操作失败')
  } finally {
    likeLoading.value = false
  }
}

// 收藏
const handleFavorite = async () => {
  if (!userStore.userInfo) {
    ElMessage.warning('请先登录')
    return
  }
  
  favoriteLoading.value = true
  try {
    const articleId = parseInt(route.params.id as string)
    const result = await articleStore.toggleFavorite(articleId)
    if (result.success) {
      ElMessage.success('操作成功')
    } else {
      ElMessage.error(result.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    favoriteLoading.value = false
  }
}

// 分享
const handleShare = () => {
  const url = window.location.href
  const title = articleStore.currentArticle?.title || '文章'
  
  if (navigator.share) {
    navigator.share({
      title,
      url
    })
  } else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(url).then(() => {
      ElMessage.success('链接已复制到剪贴板')
    }).catch(() => {
      ElMessage.error('复制失败')
    })
  }
}

      // 跳转到文章
      const goToArticle = (articleId: number) => {
        router.push(`/article/${articleId}`)
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

// 检查是否可以编辑文章
const canEditArticle = computed(() => {
  if (!userStore.userInfo || !articleStore.currentArticle) return false
  
  // 检查是否为文章作者
  const isAuthor = userStore.userInfo.id === articleStore.currentArticle.author_id
  
  // TODO: 如果有管理员角色判断，可以在这里添加
  // const isAdmin = userStore.userInfo.role === 'admin'
  
  return isAuthor
})

// 跳转到编辑页面
const goToEdit = () => {
  if (!articleStore.currentArticle) return
  router.push(`/write/${articleStore.currentArticle.id}`)
}

// 跳转到登录页面
const goToLogin = () => {
  router.push('/login')
}

// 退出登录
const handleLogout = async () => {
  userStore.logout()
  ElMessage.success('退出成功')
  router.push('/')
}

// 监听滚动更新目录
const handleScroll = () => {
  const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6')
  let currentHeading = ''
  
  headings.forEach((heading) => {
    const rect = heading.getBoundingClientRect()
    // 当标题顶部距离视窗顶部小于100px时，认为该标题是当前标题
    if (rect.top <= 100 && rect.bottom > 100) {
      currentHeading = heading.id
    }
  })
  
  activeHeading.value = currentHeading
  
  // 智能目录滚动逻辑
  const tocCard = document.querySelector('.toc-card') as HTMLElement
  if (!tocCard) return
  
  const dataCard = document.querySelector('.data-card') as HTMLElement
  if (!dataCard) return
  
  const articleContainer = document.querySelector('.article-container') as HTMLElement
  if (!articleContainer) return
  
  // 计算数据卡片的底部位置
  const dataCardRect = dataCard.getBoundingClientRect()
  const dataCardBottom = dataCardRect.bottom + window.scrollY
  
  // 计算视口高度
  const viewportHeight = window.innerHeight
  
  // 当前滚动位置
  const currentScrollY = window.scrollY
  
  // 计算目录应该开始显示的位置（数据卡片底部 + 30px间距）
  const tocStartPosition = dataCardBottom + 30
  
  // 如果页面还没滚动到数据卡片底部，隐藏目录
  if (currentScrollY + 140 < tocStartPosition) {
    tocCard.style.opacity = '0'
    tocCard.style.pointerEvents = 'none'
    return
  } else {
    tocCard.style.opacity = '1'
    tocCard.style.pointerEvents = 'auto'
  }
  
  // 获取文章内容区域（不包括评论区）
  const articleContent = document.querySelector('.article-content') as HTMLElement
  if (!articleContent) return
  
  const articleContentRect = articleContent.getBoundingClientRect()
  const articleContentBottom = articleContentRect.bottom
                  
                  // 当文章内容底部接近视口底部时，开始让目录跟随滚动
  // 这样目录就不会显示在评论区
  if (articleContentBottom <= viewportHeight) {
    // 文章内容已经滚动完毕，隐藏目录
    tocCard.style.opacity = '0'
    tocCard.style.pointerEvents = 'none'
  } else {
    // 文章内容还在视口中，保持目录固定显示
    tocCard.classList.remove('follow-scroll')
    tocCard.style.transform = 'translateY(0)'
  }
}



// 监听渲染内容变化
watch(renderedContent, () => {
  fixImageStyles()
  // 重新初始化滚动监听
  nextTick(() => {
    // 重新绑定滚动事件监听器
    window.removeEventListener('scroll', handleScroll)
    window.addEventListener('scroll', handleScroll)
    // 立即执行一次滚动检测
    handleScroll()
  })
}, { flush: 'post' })

    // 监听路由参数变化，重新加载文章数据
    watch(() => route.params.id, async (newId, oldId) => {
      if (newId !== oldId) {
        await loadArticle()
      }
    }, { immediate: false })

// 复制代码功能
const copyCode = (button: HTMLElement) => {
  const codeBlock = button.closest('.code-block')
  if (!codeBlock) return
  
  const lineContents = codeBlock.querySelectorAll('.line-content')
  if (!lineContents.length) return
  
  const text = Array.from(lineContents).map(el => el.textContent || '').join('\n')
  
  navigator.clipboard.writeText(text).then(() => {
    // 临时改变按钮文本
    const originalText = button.textContent
    button.textContent = '已复制!'
    button.classList.add('copied')
    
    setTimeout(() => {
      button.textContent = originalText
      button.classList.remove('copied')
    }, 2000)
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 将复制函数添加到全局作用域
if (typeof window !== 'undefined') {
  (window as any).copyCode = copyCode
}

// 组件挂载时加载数据
onMounted(async () => {
  // 确保用户状态已初始化
  if (userStore.getCurrentToken() && !userStore.userInfo) {
    await userStore.getUserInfo()
  }

  await loadArticle()
  
  // 监听滚动事件以更新目录
  window.addEventListener('scroll', handleScroll)
  
  // 延时修复图片样式，确保DOM完全渲染
  setTimeout(() => {
    fixImageStyles()
  }, 500)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 滚动到指定标题
const scrollToHeading = (id: string) => {
  const element = document.getElementById(id)
  if (element) {
    // 使用平滑滚动
    element.scrollIntoView({ 
      behavior: 'smooth',
      block: 'start'
    })
    
    // 更新当前活动标题
    activeHeading.value = id
    
    // 添加临时高亮效果
    element.style.transition = 'background-color 0.3s ease'
    element.style.backgroundColor = 'rgba(59, 130, 246, 0.1)'
    
    setTimeout(() => {
      element.style.backgroundColor = ''
    }, 2000)
  } else {
    console.warn(`标题元素未找到: ${id}`)
  }
}

// 添加全局复制代码函数
;(window as any).copyCode = function(button: HTMLElement) {
  const codeBlock = button.closest('.code-block')
  if (!codeBlock) return
  
  const codeContent = codeBlock.querySelector('.code-content')
  if (!codeContent) return
  
  // 获取所有代码行的内容
  const lines = Array.from(codeContent.querySelectorAll('.line-content'))
  const code = lines.map(line => line.textContent || '').join('\n')
  
  // 复制到剪贴板
  navigator.clipboard.writeText(code).then(() => {
    // 显示复制成功提示
    const originalText = button.textContent
    button.textContent = '已复制!'
    button.style.background = '#238636'
    
    setTimeout(() => {
      button.textContent = originalText
      button.style.background = '#21262d'
    }, 2000)
  }).catch(err => {
    console.error('复制失败:', err)
    ElMessage.error('复制失败')
  })
}
</script>

<style lang="scss" scoped>


.article-detail {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
}

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
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
  
  &.router-link-active {
    color: #3b82f6;
    background: rgba(59, 130, 246, 0.1);
    font-weight: bold;
  }
}

.main {
  padding: 20px 0;
}

.loading, .empty {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 40px;
  text-align: center;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin: 20px auto;
  max-width: 800px;
}

.content {
  max-width: 1600px;
  margin: 0 auto;
  padding: 0 15px;
}

// 文章容器
.article-container {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 240px;
  gap: 30px;
  align-items: start;
  max-width: 100%;
  
      @media (max-width: 1200px) {
      grid-template-columns: minmax(0, 1fr) 220px;
      gap: 25px;
    }
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 25px;
  }
}

// 文章主体
.article {
    background: #ffffff;
    border-radius: 8px;
    padding: 100px 140px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e1e4e8;
    position: relative;
    min-height: 600px;
    transition: all 0.3s ease;
    max-width: 100%;
    min-width: 0;
  
  // 全局图片样式重置 - 排除头像
  img:not(.el-avatar img) {
    max-width: calc(100% - 40px) !important;
    margin: 20px 20px !important;
    display: block !important;
    border-radius: 12px;
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
    transition: all 0.3s ease;
    object-fit: contain !important;
    height: auto !important;
  }
  
  @media (max-width: 768px) {
    padding: 60px 50px;
  }
  
  &:hover {
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
    transform: translateY(-2px);
  }
  

  
  h1.article-title {
    font-size: 52px;
    font-weight: 800;
    color: #1a1a1a;
    margin: 0 0 80px 0;
    line-height: 1.3;
    text-align: center;
    word-wrap: break-word;
    overflow-wrap: break-word;
    border-bottom: 3px solid #e5e7eb;
    padding-bottom: 30px;
    
    @media (max-width: 768px) {
      font-size: 36px;
      margin-bottom: 50px;
      text-align: left;
    }
  }
}

.article-header {
  margin-bottom: 50px;
}



.article-meta-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(59, 130, 246, 0.1);
}

.category-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #3b82f6;
  font-size: 16px;
  font-weight: 600;
  background: rgba(59, 130, 246, 0.05);
  padding: 8px 12px;
  border-radius: 15px;
  border: 1px solid rgba(59, 130, 246, 0.1);
}

.category-icon {
  font-size: 18px;
}

.read-time {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
  background: rgba(107, 114, 128, 0.05);
  padding: 8px 12px;
  border-radius: 15px;
  border: 1px solid rgba(107, 114, 128, 0.1);
}

.time-icon {
  font-size: 16px;
}

.author-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 25px;
  padding-top: 25px;
  border-top: 1px solid rgba(59, 130, 246, 0.1);
}

.author-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(59, 130, 246, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(59, 130, 246, 0.1);
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(59, 130, 246, 0.08);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
  }
}

.author-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 24px;
  font-weight: bold;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  border: 3px solid rgba(255, 255, 255, 0.8);
}

.author-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-name {
  font-size: 20px;
  font-weight: 700;
  color: #1a1a1a;
}

.author-tag {
  font-size: 12px;
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
  padding: 4px 10px;
  border-radius: 10px;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.publish-date {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.click-hint {
  font-size: 16px;
  color: #3b82f6;
  opacity: 0.7;
  transition: all 0.3s ease;
  margin-left: auto;
}

.author-info:hover .click-hint {
  opacity: 1;
  transform: translateX(2px);
}



.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 25px;
}

.tag {
  font-size: 14px;
  font-weight: 600;
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
  padding: 8px 16px;
  border-radius: 20px;
  border: 2px solid rgba(59, 130, 246, 0.2);
  transition: all 0.3s ease;
  cursor: pointer;

  &:hover {
    background: rgba(59, 130, 246, 0.2);
    border-color: rgba(59, 130, 246, 0.4);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
  }
  }
  
.article-content {
  font-size: 22px;
  line-height: 2.2;
  color: #1a1a1a;
  word-break: break-word;
  overflow-wrap: break-word;
  word-wrap: break-word;
  max-width: 100%;
  padding: 0 40px;
  margin-top: 50px;
  
  .markdown-content {
    font-size: 22px;
    line-height: 2.2;
    color: #1a1a1a;
    word-break: break-word;
    overflow-wrap: break-word;
    word-wrap: break-word;
    max-width: 90ch;
    margin: 0 auto;
    letter-spacing: 0;
    text-align: left;
    
    // 确保代码块不受markdown内容样式影响
    pre, code {
      word-break: normal !important;
      overflow-wrap: normal !important;
      word-wrap: normal !important;
      white-space: pre !important;
    }
    
    // 代码块样式 - 深色主题
    .code-block {
      background: #0d1117 !important;
      border: 1px solid #30363d !important;
      border-radius: 8px !important;
      margin: 20px 0 !important;
      overflow: hidden !important;
      
      .code-header {
        background: #161b22 !important;
        border-bottom: 1px solid #30363d !important;
        padding: 8px 16px !important;
        display: flex !important;
        justify-content: space-between !important;
        align-items: center !important;
        
        .language-label {
          color: #f0f6fc !important;
          font-size: 12px !important;
          font-weight: 500 !important;
          text-transform: uppercase !important;
        }
        
        .copy-btn {
          background: #21262d !important;
          color: #f0f6fc !important;
          border: 1px solid #30363d !important;
          border-radius: 4px !important;
          padding: 4px 8px !important;
          font-size: 12px !important;
          cursor: pointer !important;
          
          &:hover {
            background: #30363d !important;
          }
        }
      }
      
      .code-content {
        background: #0d1117 !important;
        color: #f0f6fc !important;
        padding: 16px !important;
        overflow-x: auto !important;
        overflow-y: hidden !important;
        font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace !important;
        font-size: 14px !important;
        line-height: 1.5 !important;
        
        &::-webkit-scrollbar {
          height: 8px !important;
        }
        
        &::-webkit-scrollbar-track {
          background: #161b22 !important;
          border-radius: 4px !important;
        }
        
        &::-webkit-scrollbar-thumb {
          background: #484f58 !important;
          border-radius: 4px !important;
          
          &:hover {
            background: #6e7681 !important;
          }
        }
        
        .code-line {
          display: flex !important;
          align-items: flex-start !important;
          min-height: 21px !important;
          
          .line-number {
            color: #6e7681 !important;
            margin-right: 16px !important;
            min-width: 40px !important;
            text-align: right !important;
            user-select: none !important;
            font-size: 12px !important;
            line-height: 1.5 !important;
          }
          
          .line-content {
            color: #f0f6fc !important;
            white-space: pre !important;
            word-wrap: normal !important;
            word-break: normal !important;
            flex: 1 !important;
            font-size: 14px !important;
            line-height: 1.5 !important;
          }
        }
      }
    }
    
    // 传统代码块样式（备用）
    pre {
      background: #0d1117 !important;
      color: #f0f6fc !important;
      border: 1px solid #30363d !important;
      border-radius: 6px !important;
      padding: 16px !important;
      margin: 16px 0 !important;
      overflow-x: auto !important;
      overflow-y: hidden !important;
      white-space: pre !important;
      word-wrap: normal !important;
      word-break: normal !important;
      font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace !important;
      font-size: 14px !important;
      line-height: 1.5 !important;
      
      code {
        background: transparent !important;
        color: inherit !important;
        padding: 0 !important;
        border-radius: 0 !important;
        border: none !important;
        font-family: inherit !important;
        font-size: inherit !important;
        white-space: pre !important;
      }
    }
    
    // 行内代码样式
    code:not(pre code) {
      background: #f6f8fa !important;
      color: #d73a49 !important;
      padding: 2px 4px !important;
      border-radius: 3px !important;
      font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace !important;
      font-size: 0.9em !important;
      border: 1px solid #e1e4e8 !important;
      white-space: pre !important;
      word-wrap: normal !important;
      word-break: normal !important;
    }
        
    // 标题样式优化
    h1, h2, h3, h4, h5, h6 {
      margin-top: 80px;
      margin-bottom: 40px;
      color: #1a1a1a;
      font-weight: 700;
      line-height: 1.5;
      position: relative;
      scroll-margin-top: 100px;
      text-rendering: optimizeLegibility;
      -webkit-font-smoothing: antialiased;
      
      &:first-child {
        margin-top: 0;
      }
    }
    
    h1 {
      font-size: 40px;
      margin-bottom: 45px;
      border-bottom: 3px solid #e5e7eb;
      padding-bottom: 20px;
    }
    
    h2 {
      font-size: 36px;
      margin-bottom: 40px;
      border-bottom: 3px solid #e5e7eb;
      padding-bottom: 18px;
    }
    
    h3 {
      font-size: 32px;
      margin-bottom: 35px;
    }
    
    h4 {
      font-size: 28px;
      margin-bottom: 30px;
    }
    
    h5 {
      font-size: 24px;
      margin-bottom: 25px;
    }
    
    h6 {
      font-size: 20px;
      margin-bottom: 20px;
    }
    
    // 段落样式优化
    p {
      margin-bottom: 45px;
      line-height: 2.3;
      color: #1a1a1a;
      word-wrap: break-word;
      overflow-wrap: break-word;
      font-size: 21px;
      text-align: justify;
      letter-spacing: 0.8px;
      text-indent: 0;
      max-width: 100%;
      text-rendering: optimizeLegibility;
      -webkit-font-smoothing: antialiased;
      
      // 首段特殊样式
      &:first-of-type {
        font-size: 22px;
        color: #4b5563;
        font-weight: 500;
        line-height: 2.4;
        margin-bottom: 50px;
        padding: 30px;
        background: #f8f9fa;
        border-radius: 8px;
        border-left: 4px solid #3b82f6;
      }
    }
    
    // 列表样式优化
    ul, ol {
      margin: 40px 0 40px 50px;
      padding-left: 40px;
      
      li {
        margin-bottom: 25px;
        line-height: 2.1;
        color: #1a1a1a;
        font-size: 20px;
        padding-left: 10px;
        text-rendering: optimizeLegibility;
        -webkit-font-smoothing: antialiased;
      }
    }
    
    // 引用块样式优化
    blockquote {
      margin: 50px 0;
      padding: 30px 35px;
      background: #f8f9fa;
      border-left: 5px solid #3b82f6;
      border-radius: 8px;
      color: #374151;
      font-style: italic;
      font-size: 20px;
      line-height: 2.1;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      position: relative;
      text-rendering: optimizeLegibility;
      -webkit-font-smoothing: antialiased;
      
      &::before {
        content: '"';
        position: absolute;
        top: -5px;
        left: 15px;
        font-size: 40px;
        color: #3b82f6;
        opacity: 0.3;
        font-family: serif;
      }
      
      p {
        margin: 0;
        color: #4b5563;
        font-weight: 500;
      }
    }
    
    // 表格样式优化
    table {
      width: 100%;
      border-collapse: collapse;
      margin: 50px 0;
      font-size: 20px;
      border: 1px solid #e1e4e8;
      background: white;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      text-rendering: optimizeLegibility;
      -webkit-font-smoothing: antialiased;
    }
    
    th, td {
      border: 1px solid #e5e7eb;
      padding: 25px 30px;
      text-align: left;
      vertical-align: middle;
      text-rendering: optimizeLegibility;
      -webkit-font-smoothing: antialiased;
    }
    
    th {
      background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
      color: #fff;
      font-weight: 700;
      font-size: 20px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      text-rendering: optimizeLegibility;
      -webkit-font-smoothing: antialiased;
    }
    
    td {
      background: #fff;
      color: #374151;
      font-size: 19px;
      line-height: 1.8;
      text-rendering: optimizeLegibility;
      -webkit-font-smoothing: antialiased;
    }
    
    tr:nth-child(even) td {
      background: #f9fafb;
    }
    
    tr:hover td {
      background: rgba(59, 130, 246, 0.05);
    }
    
    // 分割线样式
    hr {
      border: none;
      border-top: 3px dashed rgba(59, 130, 246, 0.4);
      margin: 50px 0;
      position: relative;
      
      &::after {
        content: '✦';
        position: absolute;
        top: -10px;
        left: 50%;
        transform: translateX(-50%);
        background: white;
        padding: 0 15px;
        color: #3b82f6;
        font-size: 18px;
      }
    }
    
    // 强调文本样式
    strong {
      font-weight: 700;
      color: #1a1a1a;
      background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
      padding: 2px 6px;
      border-radius: 4px;
    }
    
    em {
      font-style: italic;
      color: #6b7280;
      font-weight: 500;
    }
    
    // 链接样式优化
    a {
      color: #3b82f6;
      text-decoration: none;
      border-bottom: 2px solid transparent;
      transition: all 0.3s ease;
      font-weight: 600;
      padding: 2px 4px;
      border-radius: 4px;
      
      &:hover {
        color: #1d4ed8;
        background: rgba(59, 130, 246, 0.1);
        border-bottom-color: #3b82f6;
        transform: translateY(-1px);
      }
    }
    
    // 强制重置文章图片样式，排除头像
    img:not(.el-avatar img) {
      max-width: calc(100% - 40px) !important;
      width: auto !important;
      height: auto !important;
      border-radius: 12px;
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
      margin: 20px 20px !important;
      transition: all 0.3s ease;
      display: block !important;
      box-sizing: border-box !important;
      object-fit: contain;
      
      &:hover {
        transform: scale(1.02);
        box-shadow: 0 12px 35px rgba(0, 0, 0, 0.2);
      }
      
      // 移动端响应式处理，确保两侧空白对称
      @media (max-width: 768px) {
        max-width: 100% !important;
        margin: 15px 0 !important;
      }
      
      @media (max-width: 480px) {
        max-width: 100% !important;
        margin: 10px 0 !important;
      }
    }
  }
}

.article-actions {
  border-top: 2px solid rgba(59, 130, 246, 0.1);
  padding-top: 40px;
  
  .action-buttons {
    display: flex;
    gap: 20px;
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .action-btn {
    border-radius: 30px;
    padding: 15px 30px;
    font-weight: 600;
    font-size: 16px;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 140px;
    justify-content: center;
    
    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 8px 25px rgba(59, 130, 246, 0.3);
    }
    
    .el-icon {
      font-size: 18px;
    }
  }

  .like-btn {
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
    border: 2px solid rgba(59, 130, 246, 0.3);
    color: #3b82f6;

    &:hover {
      background: linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(6, 182, 212, 0.2) 100%);
      border-color: rgba(59, 130, 246, 0.5);
      color: #1d4ed8;
    }
    
    &.el-button--primary {
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
      border-color: #3b82f6;
      color: #fff;
      
      &:hover {
        background: linear-gradient(135deg, #1d4ed8 0%, #0891b2 100%);
        border-color: #1d4ed8;
      }
    }
  }

  .favorite-btn {
    background: linear-gradient(135deg, rgba(255, 193, 7, 0.1) 0%, rgba(245, 158, 11, 0.1) 100%);
    border: 2px solid rgba(255, 193, 7, 0.3);
    color: #f59e0b;

    &:hover {
      background: linear-gradient(135deg, rgba(255, 193, 7, 0.2) 0%, rgba(245, 158, 11, 0.2) 100%);
      border-color: rgba(255, 193, 7, 0.5);
      color: #d97706;
    }
    
    &.el-button--warning {
      background: linear-gradient(135deg, #f59e0b 0%, #f97316 100%);
      border-color: #f59e0b;
      color: #fff;
      
      &:hover {
        background: linear-gradient(135deg, #d97706 0%, #ea580c 100%);
        border-color: #d97706;
      }
    }
  }

  .edit-btn {
    background: linear-gradient(135deg, rgba(34, 197, 94, 0.1) 0%, rgba(22, 163, 74, 0.1) 100%);
    border: 2px solid rgba(34, 197, 94, 0.3);
    color: #22c55e;

    &:hover {
      background: linear-gradient(135deg, rgba(34, 197, 94, 0.2) 0%, rgba(22, 163, 74, 0.2) 100%);
      border-color: rgba(34, 197, 94, 0.5);
      color: #16a34a;
    }
    
    &.el-button--success {
      background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
      border-color: #22c55e;
      color: #fff;
      
      &:hover {
        background: linear-gradient(135deg, #16a34a 0%, #15803d 100%);
        border-color: #16a34a;
      }
    }
  }

  .share-btn {
    background: linear-gradient(135deg, rgba(139, 92, 246, 0.1) 0%, rgba(168, 85, 247, 0.1) 100%);
    border: 2px solid rgba(139, 92, 246, 0.3);
    color: #8b5cf6;

    &:hover {
      background: linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(168, 85, 247, 0.2) 100%);
      border-color: rgba(139, 92, 246, 0.5);
      color: #7c3aed;
    }
  }
}

// 侧边栏
.sidebar {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.data-card, .toc-card {
  background: rgba(255, 255, 255, 0.98);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 50%, #8b5cf6 100%);
  }
}

// 目录卡片特殊样式 - 智能固定定位
.toc-card {
  position: fixed;
  top: 140px; // 给数据卡片留出空间
  right: calc((100vw - 1600px) / 2 + 15px); // 计算侧边栏位置
  width: 240px;
  max-height: calc(100vh - 160px);
  overflow-y: auto;
  z-index: 100;
  opacity: 0; // 初始隐藏
  pointer-events: none; // 初始不可交互
  transition: all 0.3s ease;
  
  // 当屏幕较小时调整位置
  @media (max-width: 1630px) {
    right: 15px;
  }
  
  @media (max-width: 1200px) {
    right: calc((100vw - 100%) / 2 + 15px);
    width: 220px;
    top: 120px;
  }
  
  // 自定义滚动条样式
  &::-webkit-scrollbar {
    width: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.05);
    border-radius: 3px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(59, 130, 246, 0.3);
    border-radius: 3px;
    
    &:hover {
      background: rgba(59, 130, 246, 0.5);
    }
  }
  
  // 移动端隐藏固定目录，显示正常流布局
  @media (max-width: 900px) {
    display: none;
  }
  
  // 当需要跟随滚动时的样式
  &.follow-scroll {
    position: fixed !important;
    top: auto !important;
    right: calc((100vw - 1600px) / 2 + 15px) !important;
    
    @media (max-width: 1630px) {
      right: 15px !important;
    }
    
    @media (max-width: 1200px) {
      right: calc((100vw - 100%) / 2 + 15px) !important;
      width: 220px !important;
    }
  }
}

// 移动端目录样式
.mobile-toc {
  display: none;
  position: static !important;
  top: auto !important;
  right: auto !important;
  width: auto !important;
  max-height: none !important;
  overflow-y: visible !important;
  
  @media (max-width: 900px) {
    display: block;
  }
}

// 数据卡片需要相对定位
.data-card {
  position: relative;
  overflow: hidden;
}



.card-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 25px;
  
  .el-icon {
    font-size: 24px;
    color: #3b82f6;
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 25px;
}

.stat-box {
  text-align: center;
  padding: 15px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.05) 0%, rgba(6, 182, 212, 0.05) 100%);
  border: 1px solid rgba(59, 130, 246, 0.1);
  transition: all 0.3s ease;
  
  &:hover {
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(59, 130, 246, 0.2);
    border-color: rgba(59, 130, 246, 0.2);
  }
}

.stat-number {
  font-size: 24px;
  font-weight: 800;
  color: #1a1a1a;
  line-height: 1;
  display: block;
  margin-bottom: 5px;
  
  &.primary {
    color: #3b82f6;
  }
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 600;
}

.article-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 25px;
  padding-top: 20px;
  border-top: 1px solid rgba(59, 130, 246, 0.1);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 15px;
  color: #374151;
  font-weight: 500;
}

.info-label {
  color: #6b7280;
}

.info-value {
  color: #1a1a1a;
}

.toc-content {
  .toc-item {
    margin-bottom: 8px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    // 根据层级设置缩进和样式
    &.toc-level-1 {
      margin-left: 0;
    }
    
    &.toc-level-2 {
      margin-left: 12px;
    }
    
    &.toc-level-3 {
      margin-left: 24px;
    }
    
    &.toc-level-4 {
      margin-left: 36px;
    }
    
    &.toc-level-5 {
      margin-left: 48px;
    }
    
    &.toc-level-6 {
      margin-left: 60px;
    }
  }
  
  .toc-link {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    color: #6b7280;
    text-decoration: none;
    border-radius: 8px;
    transition: all 0.3s ease;
    line-height: 1.4;
    border: 1px solid transparent;
    position: relative;
    
    // 层级样式区分
    &.level-1 {
      font-size: 15px;
      font-weight: 600;
      padding: 10px 14px;
      
      .toc-bullet {
        color: #3b82f6;
        font-weight: bold;
      }
    }
    
    &.level-2 {
      font-size: 14px;
      font-weight: 500;
      
      .toc-bullet {
        color: #06b6d4;
      }
    }
    
    &.level-3 {
      font-size: 13px;
      font-weight: 400;
      
      .toc-bullet {
        color: #8b5cf6;
      }
    }
    
    &.level-4 {
      font-size: 12px;
      font-weight: 400;
      
      .toc-bullet {
        color: #f59e0b;
      }
    }
    
    &.level-5, &.level-6 {
      font-size: 11px;
      font-weight: 400;
      
      .toc-bullet {
        color: #10b981;
      }
    }
    
    .toc-bullet {
      font-size: 8px;
      line-height: 1;
      min-width: 20px;
      text-align: center;
      transition: all 0.3s ease;
    }
    
    .toc-text {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    
    &:hover {
      color: #3b82f6;
      background: rgba(59, 130, 246, 0.05);
      border-color: rgba(59, 130, 246, 0.1);
      transform: translateX(4px);
      
      .toc-bullet {
        color: #3b82f6;
        transform: scale(1.2);
      }
    }
    
    &.active {
      color: #3b82f6;
      background: rgba(59, 130, 246, 0.1);
      border-color: rgba(59, 130, 246, 0.2);
      font-weight: 600;
      transform: translateX(4px);
      box-shadow: 0 2px 8px rgba(59, 130, 246, 0.15);
      
      .toc-bullet {
        color: #3b82f6;
        transform: scale(1.3);
      }
      
      .toc-text {
        color: #3b82f6;
      }
    }
  }
}

.related-list {
  .related-item {
    padding: 20px;
    margin-bottom: 15px;
    background: rgba(102, 126, 234, 0.03);
    border-radius: 16px;
    border: 1px solid rgba(102, 126, 234, 0.1);
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      background: rgba(102, 126, 234, 0.08);
      border-color: rgba(102, 126, 234, 0.2);
      transform: translateY(-2px);
      box-shadow: 0 4px 15px rgba(102, 126, 234, 0.15);
    }
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .related-title {
      font-size: 16px;
      font-weight: 600;
      color: #1a1a1a;
      margin-bottom: 8px;
      line-height: 1.4;
    }
    
    .related-meta {
      display: flex;
      gap: 15px;
      font-size: 13px;
      color: #6b7280;
      
      span {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
  }
}

/* 相关文章样式 */
.related-articles-section {
  margin-top: 60px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  
  /* 当没有相关文章时，缩小模块 */
  &.empty-section {
    padding: 20px;
    margin-top: 30px;
    
    .section-header {
      margin-bottom: 15px;
      
      .section-title {
        font-size: 20px;
        
        .el-icon {
          font-size: 18px;
        }
      }
    }
  }
  
  .section-header {
    margin-bottom: 30px;
    
    .section-title {
      display: flex;
      align-items: center;
      gap: 12px;
      font-size: 28px;
      font-weight: 700;
      color: #1a202c;
      margin: 0;
      
      .el-icon {
        font-size: 24px;
        color: #3b82f6;
      }
    }
  }
  
  .related-articles-list {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 20px;
    
    .empty {
      grid-column: 1 / -1;
      text-align: center;
      padding: 20px;
      min-height: 100px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(255, 255, 255, 0.6);
      border-radius: 12px;
      border: 1px dashed rgba(59, 130, 246, 0.3);
      
      .empty-content {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12px;
        
        .empty-icon {
          font-size: 32px;
          color: #cbd5e0;
        }
        
        .empty-text {
          margin: 0;
          color: #718096;
          font-size: 14px;
        }
      }
    }
    
    .related-article-card {
      background: rgba(255, 255, 255, 0.9);
      border-radius: 16px;
      overflow: hidden;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      cursor: pointer;
      border: 1px solid rgba(255, 255, 255, 0.3);
      
      &:hover {
        transform: translateY(-8px);
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
        border-color: rgba(59, 130, 246, 0.3);
      }
      
      .article-image {
        height: 160px;
        overflow: hidden;
        position: relative;
        
        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.3s ease;
        }
        
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
        }
      }
      
      &:hover .article-image {
        img {
          transform: scale(1.05);
        }
        
        &::before {
          opacity: 1;
        }
      }
      
      .article-content {
        padding: 20px;
        
        .article-meta {
          display: flex;
          align-items: center;
          gap: 15px;
          margin-bottom: 12px;
          flex-wrap: wrap;
          
          .el-tag {
            border-radius: 8px;
            font-weight: 600;
            font-size: 12px;
            padding: 6px 12px;
            box-shadow: 0 2px 8px rgba(59, 130, 246, 0.2);
          }
          
          .article-date {
            font-size: 13px;
            color: #64748b;
            font-weight: 500;
          }
          
                      .article-stats {
              display: flex;
              gap: 8px;
              margin-left: auto;
              
              .stat-item {
                display: flex;
                align-items: center;
                gap: 3px;
                font-size: 11px;
                color: #64748b;
                font-weight: 500;
                
                .el-icon {
                  font-size: 12px;
                  color: #94a3b8;
                }
              }
            }
        }
        
                .article-title {
          font-size: 18px;
          font-weight: 700;
          color: #1a202c;
          margin: 0 0 10px 0;
          line-height: 1.4;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
          /* Standard properties for compatibility */
          display: -moz-box;
          -moz-box-orient: vertical;
          display: box;
          box-orient: vertical;
          line-clamp: 2;
        }

        .article-summary {
          color: #64748b;
          line-height: 1.6;
          margin: 0 0 15px 0;
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
          /* Standard properties for compatibility */
          display: -moz-box;
          -moz-box-orient: vertical;
          display: box;
          box-orient: vertical;
          line-clamp: 3;
        }
        
        .article-footer {
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 12px;
          
                      .author {
              display: flex;
              align-items: center;
              gap: 6px;
              
              .el-avatar {
                border: 2px solid rgba(255, 255, 255, 0.8);
                box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
              }
              
              span {
                font-size: 13px;
                font-weight: 600;
                color: #374151;
              }
            }
          
          .read-more {
            color: #3b82f6;
            font-weight: 600;
            font-size: 13px;
            text-decoration: none;
            transition: all 0.3s ease;
            padding: 6px 12px;
            border-radius: 6px;
            background: rgba(59, 130, 246, 0.1);
            cursor: pointer;
            
            &:hover {
              background: rgba(59, 130, 246, 0.2);
              transform: translateX(4px);
            }
          }
        }
      }
    }
  }
}

// 评论区
.comments-section {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 60px;
  margin-top: 40px;
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

.comments-header {
  margin-bottom: 30px;
  
  .comments-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 20px;
    font-weight: bold;
    color: #333;
    
    .el-icon {
      color: #3b82f6;
    }
  }
}

.comment-form {
  margin-bottom: 40px;
  
  .comment-input-wrapper {
    display: flex;
    gap: 20px;
    align-items: flex-start;
    padding: 25px;
    background: rgba(59, 130, 246, 0.05);
    border-radius: 15px;
    border: 1px solid rgba(59, 130, 246, 0.1);
  }
  
  .comment-input-content {
    flex: 1;
  }
  
  .comment-textarea {
    margin-bottom: 15px;
    
    :deep(.el-textarea__inner) {
      border-radius: 12px;
      border: 1px solid rgba(59, 130, 246, 0.2);
      background: rgba(255, 255, 255, 0.8);
      backdrop-filter: blur(10px);
      font-size: 16px;
      line-height: 1.6;
      padding: 15px;
      transition: all 0.3s ease;
      
      &:focus {
        border-color: #3b82f6;
        box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
      }
    }
  }
  
  .comment-actions {
    display: flex;
    justify-content: flex-end;
  }

  .submit-btn {
    border-radius: 25px;
    padding: 12px 24px;
    font-weight: 500;
    transition: all 0.3s ease;
    background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
    border: none;
    color: #fff;
    box-shadow: 0 4px 10px rgba(59, 130, 246, 0.3);
    cursor: pointer;

    &:hover {
      background: linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%);
      box-shadow: 0 6px 15px rgba(59, 130, 246, 0.4);
    }
  }

  .login-btn {
    border-radius: 25px;
    padding: 12px 24px;
    font-weight: 500;
    transition: all 0.3s ease;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    color: #fff;
    box-shadow: 0 4px 10px rgba(102, 126, 234, 0.3);

    &:hover {
      background: linear-gradient(135deg, #764ba2 0%, #667eea 100%);
      box-shadow: 0 6px 15px rgba(102, 126, 234, 0.4);
    }
  }
}

.login-prompt {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 20px;
  border: 2px dashed rgba(102, 126, 234, 0.3);
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }
}

  .comments-list {
    .no-comments {
      text-align: center;
      padding: 40px 0;
    }
  }




// 评论禁用样式
.comments-disabled {
  text-align: center;
  padding: 40px 20px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.3);
  
  .el-empty {
    .el-empty__description {
      font-size: 16px;
      color: #6b7280;
      margin-bottom: 15px;
    }
  }
  
  .disabled-message {
    font-size: 14px;
    color: #9ca3af;
    margin: 0;
  }
}


// 响应式设计
@media (max-width: 1200px) {
  .content {
    max-width: 100%;
    padding: 0 15px;
  }
  
  .article-container {
    grid-template-columns: 1fr 220px;
    gap: 25px;
  }
}

@media (max-width: 1024px) {
  .content {
    padding: 0 10px;
  }
  
  .article-container {
    grid-template-columns: 1fr 200px;
    gap: 20px;
  }
}

@media (max-width: 900px) {
  .article-container {
    grid-template-columns: 1fr;
    gap: 25px;
  }
  
  .sidebar {
    order: -1;
  }
}

// 登录检查样式
.login-required {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border-radius: 16px;
  margin: 30px 0;
  
  .login-required-content {
    text-align: center;
    padding: 40px;
    
    .login-icon {
      font-size: 64px;
      color: #64748b;
      margin-bottom: 20px;
    }
    
    h3 {
      font-size: 24px;
      font-weight: 600;
      color: #1e293b;
      margin-bottom: 12px;
    }
    
    p {
      font-size: 16px;
      color: #64748b;
      margin-bottom: 24px;
      line-height: 1.6;
    }
  }
}

@media (max-width: 768px) {
  .content {
    padding: 0 8px;
  }
  
  .article {
    padding: 35px 25px;
  }
  
  .article-content {
    padding: 0 15px;
  }
  
  .article-title {
    font-size: 28px;
  }
  
  .article-meta-top {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .category-info, .read-time {
    width: 100%;
    justify-content: center;
  }

  .author-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }



  .data-card {
    padding: 25px;
  }

  .stats-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .stat-box {
    padding: 12px;
  }

  .stat-number {
    font-size: 22px;
  }

  .article-info {
    padding-top: 15px;
  }

  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .info-label {
    width: 100%;
    text-align: center;
  }

  .comments-section {
    padding: 30px 20px;
  }
  
  .comment-form {
    .comment-input-wrapper {
      padding: 15px;
      gap: 12px;
    }
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: center;
    gap: 10px;
  }
  
  .action-btn {
    width: 100%;
    max-width: 250px;
  }
  
  .related-articles-section {
    padding: 25px 20px;
    
    .section-title {
      font-size: 24px;
    }
    
    .related-articles-list {
      grid-template-columns: 1fr;
    }
  }
  
  // 目录移动端优化
  .toc-card {
    padding: 20px;
    
    .toc-content {
      max-height: 300px;
    }
    
    .toc-link {
      padding: 6px 10px;
      font-size: 13px;
      
      &.level-1 {
        font-size: 14px;
        padding: 8px 12px;
      }
      
      &.level-2 {
        font-size: 13px;
      }
      
      &.level-3, &.level-4, &.level-5, &.level-6 {
        font-size: 12px;
      }
    }
  }
}

@media (max-width: 480px) {
  .content {
    padding: 0 5px;
  }
  
  .article {
    padding: 30px 20px;
  }
  
  .article-content {
    padding: 0 10px;
  }
  
  .article-title {
    font-size: 28px;
  }
  
  .data-card {
    padding: 25px;
  }

  .stats-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .stat-box {
    padding: 12px;
  }

  .stat-number {
    font-size: 22px;
  }

  .comments-section {
    padding: 25px;
  }
  
  .comment-form {
    .comment-input-wrapper {
      padding: 15px;
      flex-direction: column;
      align-items: center;
      text-align: center;
    }
    
    .comment-avatar {
      margin-bottom: 10px;
    }
  }
  
  .action-buttons {
    gap: 8px;
  }
  
  .action-btn {
    max-width: 100%;
    font-size: 14px;
    padding: 10px 20px;
  }
  
  .related-articles-section {
    padding: 20px;
    margin-top: 40px;
    
    .section-title {
      font-size: 22px;
    }
    
    .related-articles-list {
      gap: 15px;
    }
    
    .related-article-card {
      .article-image {
        height: 150px;
      }
      
      .article-content {
        padding: 15px;
      }
      
      .article-title {
        font-size: 16px;
      }
      
      .article-summary {
        font-size: 14px;
      }
      
      .article-meta {
        .article-stats {
          display: none;
        }
      }
    }
  }
}
</style>