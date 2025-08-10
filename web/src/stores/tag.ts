import { defineStore } from 'pinia'
import { ref } from 'vue'
import { tagApi } from '@/api/tag'
import type { Tag } from '@/types/article'

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([])
  const loading = ref(false)

  // 加载标签列表
  const loadTags = async () => {
    loading.value = true
    try {
      const response = await tagApi.getTagList()
      if (response.code === 0) {
        tags.value = response.data
      }
    } catch (error) {
      console.error('获取标签列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 刷新标签数据
  const refreshTags = async () => {
    await loadTags()
  }

  return {
    tags,
    loading,
    loadTags,
    refreshTags
  }
})
