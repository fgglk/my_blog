import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Article, ArticleListResponse } from '@/types/article'
import type { LikeResponse, FavoriteResponse } from '@/types/api'
import { articleApi } from '@/api/article'

export const useArticleStore = defineStore('article', () => {
  const articles = ref<Article[]>([])
  const currentArticle = ref<Article | null>(null)
  const loading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 获取文章列表
  const getArticles = async (page = 1, size = 10, append = false) => {
    loading.value = true
    try {
      const response = await articleApi.getArticleList({ page, size })
      if (response.code === 0) {
        const data = response.data as ArticleListResponse
        if (append) {
          // 追加模式：将新文章添加到现有列表
          articles.value.push(...data.list)
        } else {
          // 替换模式：替换整个列表
          articles.value = data.list
        }
        total.value = data.total
        currentPage.value = data.page
        pageSize.value = data.size
      }
    } catch (error) {
      console.error('获取文章列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 获取文章详情
  const getArticle = async (id: number) => {
    loading.value = true
    try {
      const response = await articleApi.getArticle(id)
      if (response.code === 0) {
        currentArticle.value = response.data

      }
    } catch (error) {
      console.error('获取文章详情失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 搜索文章
  const searchArticles = async (keyword: string, page = 1, size = 10) => {
    loading.value = true
    try {
      const response = await articleApi.searchArticles({ keyword, page, size })
      if (response.code === 0) {
        const data = response.data as ArticleListResponse
        articles.value = data.list
        total.value = data.total
        currentPage.value = data.page
        pageSize.value = data.size
      }
    } catch (error) {
      console.error('搜索文章失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 创建文章
  const createArticle = async (articleData: any) => {
    try {
      const response = await articleApi.createArticle(articleData)
      if (response.code === 0) {
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '创建文章失败' }
    }
  }

  // 更新文章
  const updateArticle = async (id: number, articleData: any) => {
    try {
      const response = await articleApi.updateArticle(id, articleData)
      if (response.code === 0) {
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '更新文章失败' }
    }
  }

  // 删除文章
  const deleteArticle = async (id: number) => {
    try {
      const response = await articleApi.deleteArticle(id)
      if (response.code === 0) {
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '删除文章失败' }
    }
  }

  // 点赞/取消点赞
  const toggleLike = async (articleId: number) => {
    try {
      const response = await articleApi.toggleLike(articleId)
      if (response.code === 0) {
        // 根据后端返回的实际状态更新前端状态
        const isLiked = (response.data as LikeResponse)?.liked || false
        if (currentArticle.value && currentArticle.value.id === articleId) {
          const wasLiked = currentArticle.value.is_liked
          currentArticle.value.is_liked = isLiked
          // 根据状态变化调整点赞数
          if (isLiked && !wasLiked) {
            currentArticle.value.like_count += 1
          } else if (!isLiked && wasLiked) {
            currentArticle.value.like_count -= 1
          }
        }
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '操作失败' }
    }
  }

  // 收藏/取消收藏
  const toggleFavorite = async (articleId: number) => {
    try {
      const response = await articleApi.toggleFavorite(articleId)
      if (response.code === 0) {
        // 根据后端返回的实际状态更新前端状态
        const isFavorited = (response.data as FavoriteResponse)?.favorited || false
        if (currentArticle.value && currentArticle.value.id === articleId) {
          const wasFavorited = currentArticle.value.is_favorited
          currentArticle.value.is_favorited = isFavorited
          // 根据状态变化调整收藏数
          if (isFavorited && !wasFavorited) {
            currentArticle.value.favorite_count += 1
          } else if (!isFavorited && wasFavorited) {
            currentArticle.value.favorite_count -= 1
          }
        }
        return { success: true }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      return { success: false, message: '操作失败' }
    }
  }

  // 获取用户文章列表
  const getUserArticles = async (page = 1, size = 10) => {
    loading.value = true
    try {
      const response = await articleApi.getUserArticles(page, size)
      if (response.code === 0) {
        const data = response.data as ArticleListResponse
        return { success: true, data }
      } else {
        return { success: false, message: response.msg }
      }
    } catch (error) {
      console.error('获取用户文章失败:', error)
      return { success: false, message: '获取用户文章失败' }
    } finally {
      loading.value = false
    }
  }

  return {
    articles,
    currentArticle,
    loading,
    total,
    currentPage,
    pageSize,
    getArticles,
    getArticle,
    searchArticles,
    createArticle,
    updateArticle,
    deleteArticle,
    toggleLike,
    toggleFavorite,
    getUserArticles
  }
}) 