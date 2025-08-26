<template>
  <div class="image-management-page">
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
            <router-link to="/write" class="nav-item">
              <el-icon><Edit /></el-icon>
              写文章
            </router-link>
            <router-link to="/profile" class="nav-item">
              <el-icon><User /></el-icon>
              个人中心
            </router-link>
          </nav>
          <div class="user-section">
            <el-dropdown @command="handleUserCommand" trigger="click">
              <div class="user-info">
                <el-avatar :src="userStore.userInfo?.avatar" :size="32" />
                <span class="username">{{ userStore.userInfo?.username }}</span>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>
                    个人中心
                  </el-dropdown-item>
                  <el-dropdown-item command="logout" divided>
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </header>

    <!-- 主内容区域 -->
    <main class="main-content">
      <div class="container">
        <!-- 页面标题 -->
        <div class="page-header">
          <h1 class="page-title">图片管理</h1>
          <p class="page-subtitle">管理您上传的图片，可以重复使用而无需重新上传</p>
        </div>

        <!-- 工具栏 -->
        <div class="toolbar">
          <div class="toolbar-left">
            <el-button type="primary" @click="showUploadDialog = true">
              <el-icon><Upload /></el-icon>
              上传图片
            </el-button>
            <el-button 
              type="danger" 
              :disabled="selectedImages.length === 0"
              @click="handleBatchDelete"
            >
              <el-icon><Delete /></el-icon>
              批量删除 ({{ selectedImages.length }})
            </el-button>
          </div>
          
          <div class="toolbar-right">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索图片名称..."
              class="search-input"
              clearable
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            
            <el-select v-model="sortBy" placeholder="排序方式" class="sort-select" @change="handleSortChange">
              <el-option label="上传时间" value="createdAt" />
              <el-option label="文件大小" value="size" />
              <el-option label="文件名" value="filename" />
            </el-select>
            
            <el-button @click="toggleSortOrder" class="sort-order-btn">
              <el-icon v-if="sortOrder === 'desc'"><ArrowDown /></el-icon>
              <el-icon v-else><ArrowUp /></el-icon>
              {{ sortOrder === 'desc' ? '降序' : '升序' }}
            </el-button>
          </div>
        </div>

        <!-- 图片网格 -->
        <div v-if="loading" class="loading-container">
          <el-skeleton :rows="3" animated />
        </div>
        
        <div v-else-if="images.length === 0" class="empty-container">
          <el-empty description="暂无图片">
            <p class="empty-tip">您可以点击上方按钮上传图片</p>
            <el-button type="primary" @click="showUploadDialog = true">
              <el-icon><Upload /></el-icon>
              上传图片
            </el-button>
          </el-empty>
        </div>
        
        <div v-else class="image-grid">
          <div 
            v-for="image in images" 
            :key="image.id" 
            class="image-item"
            :class="{ selected: selectedImages.includes(image.id) }"
            @click="toggleImageSelection(image.id)"
          >
            <div class="image-wrapper">
              <img :src="image.url" :alt="image.original_name" class="image-preview" />
              <div class="image-overlay">
                <div class="overlay-actions">
                  <el-button 
                    type="primary" 
                    size="small" 
                    @click.stop="copyImageUrl(image)"
                    title="复制图片链接"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                  <el-button 
                    type="danger" 
                    size="small" 
                    @click.stop="handleDeleteImage(image)"
                    title="删除图片"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
              <div class="image-checkbox">
                <el-checkbox 
                  :model-value="selectedImages.includes(image.id)"
                  @change="(checked: any) => handleCheckboxChange(image.id, checked)"
                  @click.stop
                />
              </div>
            </div>
            <div class="image-info">
              <div class="image-name" :title="image.original_name">
                {{ image.original_name }}
              </div>
              <div class="image-meta">
                <span class="file-size">{{ formatFileSize(image.size) }}</span>
                <span class="upload-time">{{ formatDate(image.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="totalImages > 0" class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="totalImages"
            :page-sizes="[5, 10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </main>

    <!-- 上传图片对话框 -->
    <el-dialog
      v-model="showUploadDialog"
      title="上传图片"
      width="500px"
      :close-on-click-modal="false"
    >
      <div class="upload-container">
        <el-upload
          ref="uploadRef"
          :auto-upload="false"
          :on-change="handleFileChange"
          :limit="5"
          multiple
          drag
          accept="image/*"
          class="upload-area"
          v-model:file-list="fileList"
        >
          <el-icon class="upload-icon"><Upload /></el-icon>
          <div class="upload-text">
            <span>将图片拖拽到此处，或</span>
            <em>点击上传</em>
          </div>
          <template #tip>
            <div class="upload-tip">
              支持 JPG、PNG、GIF 格式，单个文件不超过 5MB
            </div>
          </template>
        </el-upload>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showUploadDialog = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleUpload"
            :loading="uploading"
            :disabled="fileList.length === 0"
          >
            开始上传
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { imageApi } from '@/api/image'
import { ElMessage, ElMessageBox, type UploadFile, type UploadInstance } from 'element-plus'
import { 
  House, Search, Edit, User, ArrowDown, SwitchButton,
  Upload, Delete, CopyDocument, ArrowUp
} from '@element-plus/icons-vue'
import type { ImageInfo, ImageListRequest } from '@/types/image'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const loading = ref(false)
const images = ref<ImageInfo[]>([])
const totalImages = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)
const searchKeyword = ref('')
const sortBy = ref('createdAt')
const sortOrder = ref('desc')
const selectedImages = ref<number[]>([])

// 上传相关
const showUploadDialog = ref(false)
const uploading = ref(false)
const uploadRef = ref<UploadInstance>()
const fileList = ref<UploadFile[]>([])

// 获取图片列表
const loadImages = async () => {
  loading.value = true
  try {
    const params: ImageListRequest = {
      page: currentPage.value,
      size: pageSize.value,
      sortBy: sortBy.value,
      sortOrder: sortOrder.value
    }
    
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    
    const response = await imageApi.getImageList(params)
    
    if (response.code === 0) {
      images.value = response.data.list
      totalImages.value = response.data.total
      // 调试信息
      console.log('图片列表数据:', response.data.list)
      if (response.data.list.length > 0) {
        console.log('第一张图片URL:', response.data.list[0].url)
      }
    } else {
      ElMessage.error(response.msg || '获取图片列表失败')
    }
  } catch (error) {
    console.error('获取图片列表失败:', error)
    ElMessage.error('获取图片列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  loadImages()
}

// 排序处理
const handleSortChange = () => {
  currentPage.value = 1
  loadImages()
}

// 切换排序方向
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
  currentPage.value = 1
  loadImages()
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadImages()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadImages()
}

// 图片选择处理
const toggleImageSelection = (imageId: number) => {
  const index = selectedImages.value.indexOf(imageId)
  if (index > -1) {
    selectedImages.value.splice(index, 1)
  } else {
    selectedImages.value.push(imageId)
  }
}

const handleCheckboxChange = (imageId: number, checked: any) => {
  if (checked) {
    if (!selectedImages.value.includes(imageId)) {
      selectedImages.value.push(imageId)
    }
  } else {
    const index = selectedImages.value.indexOf(imageId)
    if (index > -1) {
      selectedImages.value.splice(index, 1)
    }
  }
}

// 删除图片
const handleDeleteImage = async (image: ImageInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除图片 "${image.original_name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await imageApi.deleteImage({ image_id: image.id })
    
    if (response.code === 0) {
      ElMessage.success('删除成功')
      // 从选中列表中移除
      const index = selectedImages.value.indexOf(image.id)
      if (index > -1) {
        selectedImages.value.splice(index, 1)
      }
      // 重新加载列表
      await loadImages()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除图片失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedImages.value.length === 0) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedImages.value.length} 张图片吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await imageApi.deleteImages(selectedImages.value)
    
    if (response.code === 0) {
      // 全部删除成功
      ElMessage.success(response.msg || '批量删除成功')
      selectedImages.value = []
      await loadImages()
    } else {
      // 部分删除成功或全部失败
      if (response.data && response.data.successIds && response.data.successIds.length > 0) {
        // 部分删除成功，从选中列表中移除已删除的图片
        const successIds = response.data.successIds
        selectedImages.value = selectedImages.value.filter(id => !successIds.includes(id))
        
        // 显示部分成功消息
        ElMessage.warning(response.msg || '部分删除成功')
        
        // 重新加载列表以更新显示
        await loadImages()
      } else {
        // 全部删除失败
        ElMessage.error(response.msg || '批量删除失败')
      }
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

// 复制图片链接
const copyImageUrl = async (image: ImageInfo) => {
  try {
    await navigator.clipboard.writeText(image.url)
    ElMessage.success('图片链接已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

// 文件上传处理
const handleFileChange = (file: UploadFile) => {
  // 验证文件大小（5MB）
  const maxSize = 5 * 1024 * 1024
  if (file.size && file.size > maxSize) {
    ElMessage.error('文件大小不能超过 5MB')
    return false
  }
  
  // 验证文件类型
  if (!file.raw?.type?.startsWith('image/')) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  
  return true
}

// 上传图片
const handleUpload = async () => {
  if (fileList.value.length === 0) return
  
  uploading.value = true
  const uploadPromises = fileList.value.map(async (file) => {
    if (!file.raw) return
    
    try {
      const response = await imageApi.uploadImage(file.raw)
      if (response.code === 0) {
        return { success: true, filename: file.name }
      } else {
        return { success: false, filename: file.name, error: response.msg }
      }
    } catch (error) {
      return { success: false, filename: file.name, error: '上传失败' }
    }
  })
  
  try {
    const results = await Promise.all(uploadPromises)
    const successCount = results.filter(r => r?.success).length
    const failCount = results.length - successCount
    
    if (successCount > 0) {
      ElMessage.success(`成功上传 ${successCount} 张图片`)
    }
    if (failCount > 0) {
      ElMessage.error(`${failCount} 张图片上传失败`)
    }
    
    // 重置上传状态
    showUploadDialog.value = false
    fileList.value = []
    if (uploadRef.value) {
      uploadRef.value.clearFiles()
    }
    
    // 重新加载图片列表
    loadImages()
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

// 用户操作处理
const handleUserCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      userStore.logout()
      ElMessage.success('退出成功')
      router.push('/')
      break
  }
}

// 工具函数
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString: string): string => {
  return dayjs(dateString).format('YYYY-MM-DD HH:mm')
}

// 页面加载
onMounted(() => {
  loadImages()
})
</script>

<style scoped lang="scss">
.image-management-page {
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

.user-section {
  .user-info {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 0.5rem;
    transition: background-color 0.2s;
    
    &:hover {
      background: #f3f4f6;
    }
  }
  
  .username {
    font-weight: 500;
  }
}

.main-content {
  padding: 2rem 0;
}

.page-header {
  text-align: center;
  margin-bottom: 2rem;
  
  .page-title {
    font-size: 2rem;
    font-weight: 700;
    margin: 0 0 0.5rem 0;
    color: #1f2937;
  }
  
  .page-subtitle {
    color: #6b7280;
    margin: 0;
  }
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  gap: 1rem;
  
  .toolbar-left {
    display: flex;
    gap: 1rem;
  }
  
  .toolbar-right {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .search-input {
    width: 250px;
  }
  
  .sort-select {
    width: 120px;
  }
  
  .sort-order-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 80px;
  }
}

.loading-container {
  padding: 2rem;
}

.empty-container {
  text-align: center;
  padding: 4rem 2rem;
  
  .empty-tip {
    color: #6b7280;
    margin: 1rem 0 2rem 0;
  }
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.image-item {
  background: #fff;
  border-radius: 0.75rem;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
  cursor: pointer;
  
  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transform: translateY(-2px);
  }
  
  &.selected {
    border: 2px solid #3b82f6;
  }
}

.image-wrapper {
  position: relative;
  aspect-ratio: 16/9;
  overflow: hidden;
}

.image-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.2s;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
  
  .overlay-actions {
    display: flex;
    gap: 0.5rem;
  }
}

.image-item:hover .image-overlay {
  opacity: 1;
}

.image-checkbox {
  position: absolute;
  top: 0.5rem;
  left: 0.5rem;
  z-index: 10;
}

.image-info {
  padding: 1rem;
}

.image-name {
  font-weight: 500;
  margin-bottom: 0.5rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
  color: #6b7280;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.upload-container {
  .upload-area {
    width: 100%;
  }
  
  .upload-icon {
    font-size: 3rem;
    color: #9ca3af;
    margin-bottom: 1rem;
  }
  
  .upload-text {
    color: #6b7280;
    
    em {
      color: #3b82f6;
      font-style: normal;
    }
  }
  
  .upload-tip {
    color: #9ca3af;
    font-size: 0.875rem;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

// 响应式设计
@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
    
    .toolbar-left,
    .toolbar-right {
      justify-content: center;
    }
  }
  
  .image-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .search-input {
    width: 100%;
  }
}
</style>
