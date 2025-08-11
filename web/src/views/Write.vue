<template>
  <div class="write-page">
    <!-- 顶部导航栏 -->
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
            <router-link to="/search" class="nav-item">
              <el-icon><Search /></el-icon>
              搜索
            </router-link>
            <router-link to="/write" class="nav-item active">
              <el-icon><Edit /></el-icon>
              写文章
            </router-link>
          </nav>
          <div class="user-section">
            <el-icon class="notification-icon"><Bell /></el-icon>
            <div class="user-info">
              <el-avatar :src="userStore.userInfo?.avatar" :size="32" />
              <span class="username">{{ userStore.userInfo?.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="main">
      <div class="container">
        <div class="content-layout">
          <!-- 左侧主要内容 -->
          <div class="main-content">
            <!-- 页面标题 -->
            <div class="page-header">
              <h1 class="page-title">{{ isEdit ? '编辑文章' : '发布新文章' }}</h1>
              <p class="page-subtitle">{{ isEdit ? '修改并完善你的文章内容' : '创作并分享你的知识与见解' }}</p>
            </div>

            <!-- 文章表单 -->
            <el-form
              ref="articleFormRef"
              :model="articleForm"
              :rules="articleRules"
              class="article-form"
            >
              <!-- 文章标题 -->
              <div class="form-section">
                <label class="form-label">文章标题</label>
                <el-input
                  v-model="articleForm.title"
                  placeholder="给你的文章起一个吸引人的标题..."
                  size="large"
                  class="title-input"
                />
              </div>

              <!-- 分类和标签 -->
              <div class="form-section">
                <div class="form-row">
                  <div class="form-group">
                    <label class="form-label">文章分类</label>
                    <el-select
                      v-model="articleForm.categoryId"
                      placeholder="选择分类"
                      size="large"
                      class="category-select"
                    >
                      <el-option
                        v-for="category in categories"
                        :key="category.id"
                        :label="category.name"
                        :value="category.id"
                      />
                    </el-select>
                  </div>
                  <div class="form-group">
                    <label class="form-label">文章标签</label>
                    <el-input
                      v-model="articleForm.tags"
                      placeholder="输入标签,用逗号分隔多个标签"
                      size="large"
                      class="tags-input"
                    />
                  </div>
                </div>
              </div>

              <!-- 富文本编辑器 -->
              <div class="form-section">
                <div class="editor-container">
                  <div class="editor-toolbar">
                    <el-button-group>
                      <el-tooltip content="粗体文本 (Ctrl+B)" placement="top">
                        <el-button @click="insertMarkdown('**', '**')" size="small" class="toolbar-btn">
                          <strong>B</strong>
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="斜体文本 (Ctrl+I)" placement="top">
                        <el-button @click="insertMarkdown('*', '*')" size="small" class="toolbar-btn">
                          <em>I</em>
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="引用文本" placement="top">
                        <el-button @click="insertMarkdown('> ', '')" size="small" class="toolbar-btn">
                          "
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="代码块" placement="top">
                        <el-button @click="insertMarkdown('`', '`')" size="small" class="toolbar-btn">
                          &lt;/&gt;
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="插入链接" placement="top">
                        <el-button @click="insertMarkdown('[', '](url)')" size="small" class="toolbar-btn">
                          <el-icon><Link /></el-icon>
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="插入图片" placement="top">
                        <el-button @click="showImageUploadDialog" size="small" class="toolbar-btn">
                          <el-icon><Picture /></el-icon>
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="有序列表" placement="top">
                        <el-button @click="insertMarkdown('1. ', '')" size="small" class="toolbar-btn">
                          <el-icon><Message /></el-icon>
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="无序列表" placement="top">
                        <el-button @click="insertMarkdown('- ', '')" size="small" class="toolbar-btn">
                          <el-icon><Collection /></el-icon>
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="二级标题" placement="top">
                        <el-button @click="insertMarkdown('## ', '')" size="small" class="toolbar-btn">
                          H
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="撤销 (Ctrl+Z)" placement="top">
                        <el-button @click="undo" size="small" class="toolbar-btn">
                          <el-icon><Refresh /></el-icon>
                        </el-button>
                      </el-tooltip>
                      <el-tooltip content="重做 (Ctrl+Y)" placement="top">
                        <el-button @click="redo" size="small" class="toolbar-btn">
                          <el-icon><Refresh /></el-icon>
                        </el-button>
                      </el-tooltip>
                    </el-button-group>
                  </div>
                  <div class="editor-content">
                    <el-input
                      v-model="articleForm.content"
                      type="textarea"
                      :rows="20"
                      placeholder="请输入文章内容，支持 Markdown 格式..."
                      class="content-textarea"
                    />
                    <div class="editor-footer">
                      <span class="word-count">{{ getWordCount(articleForm.content) }}字</span>
                      <div v-if="articleForm.content.length < 10" class="validation-error">
                        <el-icon><InfoFilled /></el-icon>
                        文章内容至少10个字符
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 文章封面图 -->
              <div class="form-section">
                <label class="form-label">文章封面图</label>
                                   <div class="cover-upload">
                     <el-upload
                       class="upload-area"
                       drag
                       action="#"
                       :auto-upload="false"
                       :on-change="handleCoverChange"
                       :show-file-list="false"
                     >
                       <div v-if="!articleForm.coverImage" class="upload-content">
                         <el-icon class="upload-icon"><Upload /></el-icon>
                         <p>点击上传或拖放图片到此处</p>
                       </div>
                       <div v-else class="uploaded-image">
                         <img :src="articleForm.coverImage" alt="封面图" />
                         <div class="image-overlay">
                           <el-button type="primary" size="small" @click="removeCoverImage">更换图片</el-button>
                         </div>
                       </div>
                     </el-upload>
                     <p class="upload-tip">支持 JPG, PNG, GIF (最大10MB)</p>
                   </div>
              </div>
            </el-form>
          </div>

          <!-- 右侧边栏 -->
          <div class="sidebar">
            <!-- 发布设置 -->
            <div class="sidebar-card">
              <h3 class="card-title">发布设置</h3>
              <p class="card-subtitle">发布状态</p>
              <div class="radio-group">
                <el-radio v-model="publishStatus" label="publish">立即发布</el-radio>
                <el-radio v-model="publishStatus" label="draft">保存为草稿</el-radio>
              </div>
            </div>

            <!-- 阅读设置 -->
            <div class="sidebar-card">
              <h3 class="card-title">阅读设置</h3>
              <p class="card-subtitle">控制文章的阅读权限和交互</p>
              <div class="reading-settings">
                <div class="setting-item">
                  <el-switch
                    v-model="readingSettings.allowComments"
                    active-text="允许评论"
                    size="large"
                  />
                </div>
                <div class="setting-item">
                  <el-switch
                    v-model="readingSettings.allowRepost"
                    active-text="允许转载"
                    size="large"
                  />
                </div>
                <div class="setting-item">
                  <el-switch
                    v-model="readingSettings.requireLogin"
                    active-text="需要登录"
                    size="large"
                  />
                </div>
              </div>
            </div>

            <!-- 统计信息 -->
            <div class="sidebar-card">
              <div class="stats-item">
                <span>预计阅读时间</span>
                <span class="stats-value">{{ getReadTime(articleForm.content) }}分钟</span>
              </div>
              <div class="stats-item">
                <span>字数统计</span>
                <span class="stats-value">{{ getWordCount(articleForm.content) }}</span>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="sidebar-card">
              <div class="action-buttons-group">
                <el-button type="primary" @click="handleSubmit" :loading="publishing" class="action-btn">
                  <el-icon><Promotion /></el-icon>
                  {{ getSubmitButtonText() }}
                </el-button>
                <el-button @click="previewArticle" class="action-btn">
                  <el-icon><View /></el-icon>
                  预览文章
                </el-button>
              </div>
            </div>

            <!-- 写作提示 -->
            <div class="sidebar-card">
              <h3 class="card-title">写作提示</h3>
              <ul class="tips-list">
                <li>
                  <el-icon class="tip-icon"><Check /></el-icon>
                  使用清晰的标题结构,帮助读者快速理解文章内容
                </li>
                <li>
                  <el-icon class="tip-icon"><Check /></el-icon>
                  添加相关代码示例可以提高文章的实用性
                </li>
                <li>
                  <el-icon class="tip-icon"><Check /></el-icon>
                  适当使用图片可以让文章更具吸引力
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 底部提示 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-hint">
            <el-icon class="hint-icon"><InfoFilled /></el-icon>
            请输入文章内容,支持Markdown格式
          </div>
          <div class="footer-copyright">
            © 2023 我的博客,保留所有权利.
          </div>
        </div>
      </div>
    </footer>
  </div>

  <!-- 图片上传对话框 -->
  <el-dialog
    v-model="imageUploadDialogVisible"
    title="插入图片"
    width="500px"
    :close-on-click-modal="false"
  >
    <div class="image-upload-dialog">
      <div class="upload-section">
        <el-upload
          class="upload-area"
          drag
          action="#"
          :auto-upload="false"
          :on-change="handleImageUpload"
          :show-file-list="false"
          accept="image/*"
        >
          <div class="upload-content">
            <el-icon class="upload-icon"><Upload /></el-icon>
            <p>点击上传或拖放图片到此处</p>
            <p class="upload-tip">支持 JPG, PNG, GIF (最大10MB)</p>
          </div>
        </el-upload>
      </div>
      
      <div v-if="uploadedImageUrl" class="preview-section">
        <h4>图片预览</h4>
        <div class="image-preview">
          <img :src="uploadedImageUrl" alt="预览图片" />
        </div>
        <div class="image-info">
          <el-input
            v-model="imageAltText"
            placeholder="图片描述（可选）"
            size="small"
          />
        </div>
      </div>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="imageUploadDialogVisible = false">取消</el-button>
        <el-button 
          type="primary" 
          @click="insertImageToContent"
          :disabled="!uploadedImageUrl"
        >
          插入图片
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 预览对话框 -->
  <el-dialog
    v-model="previewDialogVisible"
    :title="articleForm.title || '文章预览'"
    width="80%"
    :close-on-click-modal="false"
    class="preview-dialog"
  >
    <div class="preview-container">
      <!-- 文章信息 -->
      <div class="preview-header">
        <h1 class="preview-title">{{ articleForm.title }}</h1>
        <div class="preview-meta">
          <span class="meta-item">
            <el-icon><User /></el-icon>
            {{ userStore.userInfo?.username }}
          </span>
          <span class="meta-item">
            <el-icon><Clock /></el-icon>
            {{ new Date().toLocaleDateString() }}
          </span>
          <span class="meta-item">
            <el-icon><View /></el-icon>
            预计阅读 {{ getReadTime(articleForm.content) }} 分钟
          </span>
          <span class="meta-item">
            <el-icon><Document /></el-icon>
            {{ getWordCount(articleForm.content) }} 字
          </span>
        </div>
      </div>

      <!-- 文章内容 -->
      <div class="preview-content" v-html="previewContent"></div>
    </div>

    <template #footer>
      <div class="preview-footer">
        <el-button @click="previewDialogVisible = false">关闭预览</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="publishing">
          <el-icon><Promotion /></el-icon>
          {{ getSubmitButtonText() }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useArticleStore } from '@/stores/article'
import { articleApi } from '@/api/article'
import { categoryApi } from '@/api/category'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { Category, ReadingSettings } from '@/types/article'
import MarkdownIt from 'markdown-it'
import {
  House, Search, Edit, Bell, ArrowDown, Link, Message, Collection,
  Refresh, InfoFilled, Upload, Promotion, View,
  Check, Picture, User, Clock, Document
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const articleStore = useArticleStore()

const articleFormRef = ref<FormInstance>()
const saving = ref(false)
const publishing = ref(false)

// Markdown渲染器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

// 判断是否为编辑模式
const isEdit = computed(() => {
  return !!route.params.id
})

const articleForm = reactive({
  title: '',
  content: '',
  categoryId: undefined as number | undefined,
  tags: '',
  coverImage: null as string | null
})

const publishStatus = ref('draft')

// 阅读设置
const readingSettings = reactive<ReadingSettings>({
  allowComments: true,
  allowRepost: true,
  requireLogin: false
})

const articleRules: FormRules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 1, max: 100, message: '标题长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' },
    { min: 10, message: '文章内容至少 10 个字符', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ]
}

// 分类数据
const categories = ref<Category[]>([])

// 加载分类数据
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategoryList()
    if (response.code === 0) {
      categories.value = response.data
    }
  } catch (error) {
    console.error('加载分类失败:', error)
    ElMessage.error('加载分类失败')
  }
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

// 阅读时间估算
const getReadTime = (content: string) => {
  const wordCount = getWordCount(content)
  const wordsPerMinute = 200
  const readTime = Math.ceil(wordCount / wordsPerMinute)
  return readTime || 1
}

// 解析阅读设置
const parseReadingSettings = (summary: string): ReadingSettings => {
  const defaultSettings: ReadingSettings = {
    allowComments: true,
    allowRepost: true,
    requireLogin: false
  }
  
  if (!summary) return defaultSettings
  
  const settingsMatch = summary.match(/<!--READ_SETTINGS:({.*?})-->/)
  if (settingsMatch) {
    try {
      const settings = JSON.parse(settingsMatch[1])
      return {
        allowComments: settings.allowComments ?? true,
        allowRepost: settings.allowRepost ?? true,
        requireLogin: settings.requireLogin ?? false
      }
    } catch (error) {
      console.error('解析阅读设置失败:', error)
      return defaultSettings
    }
  }
  
  return defaultSettings
}



// 编码阅读设置到摘要
const encodeReadingSettings = (summary: string, settings: ReadingSettings): string => {
  const settingsJson = JSON.stringify(settings)
  const settingsTag = `<!--READ_SETTINGS:${settingsJson}-->`
  
  // 移除现有的设置标签
  const cleanSummary = summary.replace(/<!--READ_SETTINGS:.*?-->/, '').trim()
  
  // 添加新的设置标签
  return cleanSummary ? `${cleanSummary}\n\n${settingsTag}` : settingsTag
}

// 插入 Markdown 语法
const insertMarkdown = (before: string, after: string) => {
  const textarea = document.querySelector('.content-textarea textarea') as HTMLTextAreaElement
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = articleForm.content.substring(start, end)
  
  const newText = articleForm.content.substring(0, start) + 
                  before + selectedText + after + 
                  articleForm.content.substring(end)
  
  articleForm.content = newText
  
  // 设置光标位置
  setTimeout(() => {
    textarea.focus()
    if (before === '' && after.includes('![')) {
      // 对于图片插入，将光标放在alt文本位置
      const imageStart = start + before.length + 2 // 跳过 ![ 和 ]
      textarea.setSelectionRange(imageStart, imageStart)
    } else {
      textarea.setSelectionRange(start + before.length, end + before.length)
    }
  }, 0)
}

// 键盘快捷键处理
const handleKeydown = (event: KeyboardEvent) => {
  const textarea = event.target as HTMLTextAreaElement
  if (!textarea || !textarea.classList.contains('content-textarea')) return
  
  // Ctrl+B: 粗体
  if (event.ctrlKey && event.key === 'b') {
    event.preventDefault()
    insertMarkdown('**', '**')
  }
  // Ctrl+I: 斜体
  else if (event.ctrlKey && event.key === 'i') {
    event.preventDefault()
    insertMarkdown('*', '*')
  }
  // Ctrl+Z: 撤销
  else if (event.ctrlKey && event.key === 'z') {
    event.preventDefault()
    undo()
  }
  // Ctrl+Y: 重做
  else if (event.ctrlKey && event.key === 'y') {
    event.preventDefault()
    redo()
  }
}

// 撤销
const undo = () => {
  // 实现撤销功能
  ElMessage.info('撤销功能开发中')
}

// 重做
const redo = () => {
  // 实现重做功能
  ElMessage.info('重做功能开发中')
}

// 处理封面图上传
const handleCoverChange = async (file: any) => {
  try {
    // 调用图片上传API
    const response = await articleApi.uploadImage(file.raw)
    if (response.code === 0) {
      articleForm.coverImage = response.data.url
      ElMessage.success('封面图上传成功')
    } else {
      ElMessage.error(response.msg || '封面图上传失败')
    }
  } catch (error) {
    console.error('封面图上传失败:', error)
    ElMessage.error('封面图上传失败')
  }
}

// 移除封面图
const removeCoverImage = () => {
  articleForm.coverImage = null
  ElMessage.info('封面图已移除')
}

// 图片上传对话框相关状态
const imageUploadDialogVisible = ref(false)
const uploadedImageUrl = ref<string | null>(null)
const imageAltText = ref<string | null>(null)

// 显示图片上传对话框
const showImageUploadDialog = () => {
  imageUploadDialogVisible.value = true
  uploadedImageUrl.value = null
  imageAltText.value = null
}

// 处理图片上传
const handleImageUpload = async (file: any) => {
  try {
    // 调用图片上传API
    const response = await articleApi.uploadImage(file.raw)
    if (response.code === 0) {
      uploadedImageUrl.value = response.data.url
      ElMessage.success('图片上传成功')
    } else {
      ElMessage.error(response.msg || '图片上传失败')
    }
  } catch (error) {
    console.error('图片上传失败:', error)
    ElMessage.error('图片上传失败')
  }
}

// 插入图片到文章内容
const insertImageToContent = () => {
  if (!uploadedImageUrl.value) {
    ElMessage.warning('请先上传图片')
    return
  }

  const altText = imageAltText.value || '图片'
  const markdownImage = `![${altText}](${uploadedImageUrl.value})`

  insertMarkdown('', markdownImage)
  imageUploadDialogVisible.value = false
  uploadedImageUrl.value = null
  imageAltText.value = null
}

// 获取提交按钮文本
const getSubmitButtonText = () => {
  if (isEdit.value) {
    return '更新文章'
  }
  return publishStatus.value === 'publish' ? '发布文章' : '保存草稿'
}

// 统一提交处理函数
const handleSubmit = async () => {
  if (publishStatus.value === 'draft') {
    await saveDraft()
  } else {
    await publishArticle()
  }
}

// 保存草稿
const saveDraft = async () => {
  if (!articleFormRef.value) return
  
  try {
    await articleFormRef.value.validate()
    saving.value = true
    
         // 处理标签字符串，转换为标签名称数组
         const tagNames = articleForm.tags
           .split(/[,，]/) // 支持英文逗号和中文逗号
           .map(tag => tag.trim())
           .filter(tag => tag.length > 0)
         
         // 编码阅读设置到摘要字段
         const summaryWithSettings = encodeReadingSettings('', readingSettings)
         
         const articleData = {
           title: articleForm.title.trim(),
           content: articleForm.content,
           summary: summaryWithSettings, // 包含阅读设置的摘要
           category_id: articleForm.categoryId || 3, // 使用选中的分类ID，默认为3（技术）
           tag_names: tagNames, // 发送标签名称数组
           cover_image: articleForm.coverImage, // 添加封面图片
           status: 0 // 0 表示草稿，1 表示已发布
         }
         
         // 调试信息
         
    
    let result
    if (isEdit.value) {
      const articleId = parseInt(route.params.id as string)
      result = await articleStore.updateArticle(articleId, articleData)
    } else {
      result = await articleStore.createArticle(articleData)
    }
    
    if (result.success) {
      ElMessage.success('草稿保存成功')
      // 草稿保存成功后不跳转，用户可以继续编辑
    } else {
      ElMessage.error(result.message || '保存失败')
    }
  } catch (error) {
    console.error('保存草稿失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 发布文章
const publishArticle = async () => {
  if (!articleFormRef.value) return
  
  try {
    await articleFormRef.value.validate()
    publishing.value = true
    
    // 检查用户登录状态
    if (!userStore.userInfo) {
      ElMessage.error('请先登录')
      return
    }
    
         // 处理标签字符串，转换为标签名称数组
     const tagNames = articleForm.tags
       .split(/[,，]/) // 支持英文逗号和中文逗号
       .map(tag => tag.trim())
       .filter(tag => tag.length > 0)
     
     // 编码阅读设置到摘要字段
     const summaryWithSettings = encodeReadingSettings('', readingSettings)
     
     const articleData = {
       title: articleForm.title.trim(),
       content: articleForm.content,
       summary: summaryWithSettings, // 包含阅读设置的摘要
       category_id: articleForm.categoryId || 3, // 使用选中的分类ID，默认为3（技术）
       tag_names: tagNames, // 发送标签名称数组
       cover_image: articleForm.coverImage, // 添加封面图片
       status: 1 // 发布文章时总是设置为已发布状态
     }
     
     // 调试信息
     
    
    let result
    if (isEdit.value) {
      const articleId = parseInt(route.params.id as string)
      result = await articleStore.updateArticle(articleId, articleData)
    } else {
      result = await articleStore.createArticle(articleData)
    }
    
    if (result.success) {
      ElMessage.success(isEdit.value ? '文章更新成功' : '文章发布成功')
      router.push('/')
    } else {
      ElMessage.error(result.message || '操作失败')
    }
  } catch (error) {
    console.error('发布文章失败:', error)
    ElMessage.error('操作失败')
  } finally {
    publishing.value = false
  }
}

// 预览相关状态
const previewDialogVisible = ref(false)
const previewContent = ref('')

// 预览文章
const previewArticle = () => {
  if (!articleForm.title.trim()) {
    ElMessage.warning('请先输入文章标题')
    return
  }
  if (!articleForm.content.trim()) {
    ElMessage.warning('请先输入文章内容')
    return
  }
  
  // 渲染Markdown内容
  previewContent.value = md.render(articleForm.content)
  previewDialogVisible.value = true
}

// 加载文章数据（编辑模式）
const loadArticle = async () => {
  if (!isEdit.value) return
  
  const articleId = parseInt(route.params.id as string)
  if (isNaN(articleId)) {
    ElMessage.error('文章ID无效')
    router.push('/write')
    return
  }
  
  try {
    await articleStore.getArticle(articleId)
    const article = articleStore.currentArticle
    
    if (article) {
      articleForm.title = article.title
      articleForm.content = article.content
      articleForm.categoryId = article.category_id
      articleForm.tags = article.tags.map(tag => tag.name).join(',')
      articleForm.coverImage = article.cover_image || null
      publishStatus.value = article.is_published ? 'publish' : 'draft'
      
      // 解析阅读设置
      if (article.summary) {
        const settings = parseReadingSettings(article.summary)
        readingSettings.allowComments = settings.allowComments
        readingSettings.allowRepost = settings.allowRepost
        readingSettings.requireLogin = settings.requireLogin
      }
    }
  } catch (error) {
    ElMessage.error('加载文章失败')
    router.push('/write')
  }
}

// 在组件挂载时加载分类
onMounted(() => {
  loadCategories()
  loadArticle()
  document.addEventListener('keydown', handleKeydown)
})

// 在组件卸载时移除事件监听
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style lang="scss" scoped>
.write-page {
  min-height: 100vh;
  background-color: #f8f9fa;
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
    margin: 0;
  }
}

.nav {
  display: flex;
  align-items: center;
  gap: 30px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 5px;
  text-decoration: none;
  color: #666;
  font-size: 16px;
  transition: color 0.3s ease;
  
  &:hover {
    color: #409eff;
  }
  
  &.active {
    color: #409eff;
    font-weight: bold;
  }
}

.user-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.notification-icon {
  font-size: 20px;
  color: #666;
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  
  .username {
    font-size: 14px;
    color: #333;
  }
}

.main {
  padding: 30px 0;
}

.content-layout {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 30px;
  align-items: start;
}

.main-content {
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.page-header {
  margin-bottom: 40px;
  
  .page-title {
    font-size: 32px;
    font-weight: bold;
    color: #333;
    margin: 0 0 10px 0;
  }
  
  .page-subtitle {
    font-size: 16px;
    color: #666;
    margin: 0;
  }
}

.form-section {
  margin-bottom: 30px;
}

.form-label {
  display: block;
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 10px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 20px;
}

.title-input {
  font-size: 18px;
  font-weight: 500;
}

.category-select,
.tags-input {
  width: 100%;
}

.editor-container {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

.editor-toolbar {
  padding: 12px;
  background: #f8f9fa;
  border-bottom: 1px solid #e4e7ed;
  
  .toolbar-btn {
    margin-right: 5px;
  }
}

.editor-content {
  position: relative;
}

.content-textarea {
  border: none;
  resize: none;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  
  :deep(.el-textarea__inner) {
    border: none;
    padding: 20px;
    min-height: 400px;
  }
}

.editor-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background: #f8f9fa;
  border-top: 1px solid #e4e7ed;
  
  .word-count {
    font-size: 14px;
    color: #666;
  }
  
  .validation-error {
    display: flex;
    align-items: center;
    gap: 5px;
    color: #f56c6c;
    font-size: 14px;
  }
}

.cover-upload {
  .upload-area {
    width: 100%;
  }
  
  .upload-content {
    padding: 40px;
    text-align: center;
    
    .upload-icon {
      font-size: 48px;
      color: #c0c4cc;
      margin-bottom: 10px;
    }
    
    p {
      color: #666;
      margin: 0;
    }
  }
  
  .uploaded-image {
    position: relative;
    width: 100%;
    height: 200px;
    border-radius: 8px;
    overflow: hidden;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    
    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.5);
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: opacity 0.3s ease;
      
      &:hover {
        opacity: 1;
      }
    }
  }
  
  .upload-tip {
    margin-top: 10px;
    font-size: 14px;
    color: #999;
    text-align: center;
  }
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  
  .card-title {
    font-size: 18px;
    font-weight: 600;
    color: #333;
    margin: 0 0 10px 0;
  }
  
  .card-subtitle {
    font-size: 14px;
    color: #666;
    margin: 0 0 15px 0;
  }
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.reading-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
  
  .setting-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 0;
    border-bottom: 1px solid #f0f2f5;
    
    &:last-child {
      border-bottom: none;
    }
    
    :deep(.el-switch) {
      .el-switch__label {
        font-size: 14px;
        color: #374151 !important;
        font-weight: 500 !important;
        transition: all 0.3s ease !important;
      }
      
      .el-switch__core {
        background-color: #e5e7eb !important;
        border-color: #d1d5db !important;
        transition: all 0.3s ease !important;
        width: 50px !important;
        height: 24px !important;
        
        &.is-checked {
          background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
          border-color: #3b82f6 !important;
          box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2), 0 2px 8px rgba(59, 130, 246, 0.3) !important;
        }
        
        .el-switch__action {
          transition: all 0.3s ease !important;
          width: 20px !important;
          height: 20px !important;
          background-color: #ffffff !important;
          border-radius: 50% !important;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1) !important;
          
          &.is-checked {
            background-color: #ffffff !important;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15) !important;
            transform: translateX(26px) !important;
          }
        }
      }
      
      &.is-checked {
        .el-switch__label {
          color: #3b82f6 !important;
          font-weight: 600 !important;
          text-shadow: 0 0 1px rgba(59, 130, 246, 0.1) !important;
        }
      }
      
      &:hover {
        .el-switch__core {
          &.is-checked {
            background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%) !important;
            box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.3), 0 4px 12px rgba(59, 130, 246, 0.4) !important;
          }
        }
      }
    }
    
        // 额外的强制样式，确保开关按钮变蓝
    :deep(.el-switch.is-checked) {
      .el-switch__core {
        background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
        border-color: #3b82f6 !important;
      }
      
      .el-switch__action {
        background-color: #ffffff !important;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15) !important;
      }
    }
    
    // 全局覆盖 Element Plus 的开关样式
    :deep(.el-switch) {
      &.is-checked {
        .el-switch__core {
          background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
          border-color: #3b82f6 !important;
        }
        
        .el-switch__action {
          background-color: #ffffff !important;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15) !important;
        }
      }
    }
  }
}

.switch-group {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.switch-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #333;
}

.stats-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 14px;
  color: #333;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  .stats-value {
    font-weight: 600;
    color: #409eff;
  }
}

.action-buttons-group {
  display: flex;
  gap: 12px;
  
  .action-btn {
    flex: 1;
    margin-bottom: 0;
    transition: all 0.3s ease;
    
    &.el-button--primary {
      background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
      border-color: #3b82f6;
      box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
      
      &:hover {
        background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
        border-color: #2563eb;
        box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
        transform: translateY(-1px);
      }
      
      &:active {
        transform: translateY(0);
        box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
      }
    }
    
    &.el-button--default {
      border-color: #d1d5db;
      color: #374151;
      
      &:hover {
        border-color: #3b82f6;
        color: #3b82f6;
        background-color: rgba(59, 130, 246, 0.05);
      }
    }
  }
}



.tips-list {
  list-style: none;
  padding: 0;
  margin: 0;
  
  li {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 12px;
    font-size: 14px;
    color: #666;
    line-height: 1.5;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .tip-icon {
      color: #67c23a;
      margin-top: 2px;
      flex-shrink: 0;
    }
  }
}

.footer {
  margin-top: 40px;
  padding: 20px 0;
  border-top: 1px solid #e4e7ed;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
  
  .hint-icon {
    color: #909399;
  }
}

.footer-copyright {
  font-size: 14px;
  color: #999;
}

.image-upload-dialog {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.upload-section {
  .upload-area {
    width: 100%;
  }
  
  .upload-content {
    padding: 40px;
    text-align: center;
    
    .upload-icon {
      font-size: 48px;
      color: #c0c4cc;
      margin-bottom: 10px;
    }
    
    p {
      color: #666;
      margin: 0;
    }
  }
}

.preview-section {
  .image-preview {
    width: 100%;
    height: 150px;
    border-radius: 8px;
    overflow: hidden;
    margin-bottom: 10px;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: contain;
    }
  }
  
  .image-info {
    .el-input {
      width: 100%;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

// 预览对话框样式
.preview-dialog {
  .el-dialog__body {
    padding: 0;
    max-height: 70vh;
    overflow-y: auto;
  }
}

.preview-container {
  padding: 30px;
  background: #fff;
}

.preview-header {
  border-bottom: 2px solid #f0f2f5;
  padding-bottom: 20px;
  margin-bottom: 30px;
}

.preview-title {
  font-size: 28px;
  font-weight: bold;
  color: #1a1a1a;
  line-height: 1.3;
  margin: 0 0 15px 0;
}

.preview-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  color: #666;
  font-size: 14px;

  .meta-item {
    display: flex;
    align-items: center;
    gap: 5px;
    
    .el-icon {
      font-size: 16px;
      color: #3b82f6;
    }
  }
}

.preview-content {
  // 复用文章详情页的样式
  line-height: 1.8;
  color: #333;
  font-size: 16px;
  
  h1, h2, h3, h4, h5, h6 {
    color: #1a1a1a;
    font-weight: 600;
    line-height: 1.4;
    margin: 25px 0 15px 0;
    
    &:first-child {
      margin-top: 0;
    }
  }
  
  h1 { font-size: 24px; }
  h2 { font-size: 22px; }
  h3 { font-size: 20px; }
  h4 { font-size: 18px; }
  h5 { font-size: 16px; }
  h6 { font-size: 14px; }
  
  p {
    margin: 15px 0;
    line-height: 1.8;
    
    &:first-child {
      margin-top: 0;
    }
    
    &:last-child {
      margin-bottom: 0;
    }
  }
  
  img {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
    margin: 20px 0;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  blockquote {
    border-left: 4px solid #3b82f6;
    background: rgba(59, 130, 246, 0.05);
    margin: 20px 0;
    padding: 15px 20px;
    border-radius: 0 8px 8px 0;
    
    p {
      margin: 0;
      color: #555;
      font-style: italic;
    }
  }
  
  code {
    background: #f1f3f4;
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    color: #e74c3c;
  }
  
  pre {
    background: #f8f9fa;
    border: 1px solid #e9ecef;
    border-radius: 8px;
    padding: 20px;
    margin: 20px 0;
    overflow-x: auto;
    
    code {
      background: none;
      padding: 0;
      color: #333;
      font-size: 14px;
    }
  }
  
  ul, ol {
    margin: 15px 0;
    padding-left: 25px;
    
    li {
      margin: 8px 0;
      line-height: 1.6;
    }
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
    margin: 20px 0;
    
    th, td {
      border: 1px solid #e5e7eb;
      padding: 12px 15px;
      text-align: left;
    }
    
    th {
      background: #f9fafb;
      font-weight: 600;
      color: #374151;
    }
  }
  
  a {
    color: #3b82f6;
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
}

.preview-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 30px;
  border-top: 1px solid #f0f2f5;
  background: #fafbfc;
}

// 预览对话框移动端适配
@media (max-width: 768px) {
  .preview-dialog {
    .el-dialog {
      width: 95% !important;
      margin: 20px auto !important;
    }
    
    .el-dialog__body {
      max-height: 60vh;
    }
  }
  
  .preview-container {
    padding: 20px;
  }
  
  .preview-title {
    font-size: 24px;
  }
  
  .preview-meta {
    gap: 15px;
    font-size: 13px;
    
    .meta-item {
      gap: 4px;
      
      .el-icon {
        font-size: 14px;
      }
    }
  }
  
  .preview-content {
    font-size: 15px;
    
    h1 { font-size: 22px; }
    h2 { font-size: 20px; }
    h3 { font-size: 18px; }
    h4 { font-size: 16px; }
    h5 { font-size: 15px; }
    h6 { font-size: 14px; }
  }
  
  .preview-footer {
    padding: 15px 20px;
    flex-direction: column;
    gap: 10px;
    
    .el-button {
      width: 100%;
    }
  }
}

@media (max-width: 1200px) {
  .content-layout {
    grid-template-columns: 1fr;
  }
  
  .sidebar {
    order: -1;
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 20px;
  }
  
  .page-header {
    margin-bottom: 30px;
    
    .page-title {
      font-size: 24px;
    }
  }
  
  .form-row {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .header-content {
    flex-direction: column;
    gap: 15px;
  }
  
  .nav {
    gap: 20px;
  }
  
  .user-section {
    gap: 15px;
  }
  
  // 确保按钮组在移动端正确显示
  .action-buttons-group {
    flex-direction: column;
    
    .action-btn {
      width: 100%;
      transition: all 0.3s ease;
      
      &.el-button--primary {
        background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
        border-color: #3b82f6;
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
        
        &:hover {
          background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
          border-color: #2563eb;
          box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
          transform: translateY(-1px);
        }
        
        &:active {
          transform: translateY(0);
          box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
        }
      }
      
      &.el-button--default {
        border-color: #d1d5db;
        color: #374151;
        
        &:hover {
          border-color: #3b82f6;
          color: #3b82f6;
          background-color: rgba(59, 130, 246, 0.05);
        }
      }
    }
  }
}

// 全局覆盖 Element Plus 开关样式
:deep(.el-switch.is-checked) {
  .el-switch__core {
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
    border-color: #3b82f6 !important;
  }
  
  .el-switch__action {
    background-color: #ffffff !important;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15) !important;
  }
}

:deep(.el-switch) {
  &.is-checked {
    .el-switch__core {
      background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
      border-color: #3b82f6 !important;
    }
    
    .el-switch__action {
      background-color: #ffffff !important;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15) !important;
    }
  }
}
</style> 