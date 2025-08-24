<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="900px"
    :close-on-click-modal="false"
    class="image-selector-dialog"
  >
    <div class="image-selector">
      <div class="image-grid">
        <div 
          v-for="image in images" 
          :key="image.id" 
          class="image-item"
          :class="{ selected: selectedImageId === image.id }"
          @click="selectedImageId = image.id"
        >
          <div class="image-container">
            <img :src="image.url" :alt="image.original_name" class="image-preview" />
            <div class="image-overlay">
              <div class="overlay-content">
                <el-icon class="check-icon"><Check /></el-icon>
                <span class="image-name">{{ image.original_name }}</span>
              </div>
            </div>
          </div>
          <div class="image-info">
            <span class="image-size">{{ formatFileSize(image.size) }}</span>
            <span class="image-dimensions" v-if="image.width && image.height">
              {{ image.width }}×{{ image.height }}
            </span>
          </div>
        </div>
      </div>
      <div v-if="images.length === 0" class="empty-state">
        <el-empty description="暂无图片">
          <div class="empty-content">
            <p class="empty-text">您还没有上传过图片</p>
            <p class="empty-hint">请先上传一些图片，然后就可以在这里选择使用了</p>
          </div>
        </el-empty>
      </div>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel" size="large">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleConfirm"
          :disabled="!selectedImageId"
          size="large"
        >
          {{ confirmText }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElIcon } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { imageApi } from '@/api/image'
import type { ImageInfo } from '@/types/image'

interface Props {
  modelValue: boolean
  title?: string
  confirmText?: string
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm', image: ImageInfo): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  title: '选择图片',
  confirmText: '选择此图片'
})

const emit = defineEmits<Emits>()

const visible = ref(false)
const images = ref<ImageInfo[]>([])
const selectedImageId = ref<number | null>(null)

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) {
    loadImages()
  }
})

// 监听 visible 变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
  if (!newVal) {
    selectedImageId.value = null
  }
})

// 加载用户图片
const loadImages = async () => {
  try {
    const response = await imageApi.getImageList({
      page: 1,
      size: 100,
      sortBy: 'createdAt',
      sortOrder: 'desc'
    })
    if (response.code === 0) {
      images.value = response.data.list
    }
  } catch (error) {
    console.error('加载用户图片失败:', error)
  }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// 确认选择
const handleConfirm = () => {
  if (!selectedImageId.value) return
  
  const selectedImage = images.value.find(img => img.id === selectedImageId.value)
  if (selectedImage) {
    emit('confirm', selectedImage)
    visible.value = false
  }
}

// 取消选择
const handleCancel = () => {
  emit('cancel')
  visible.value = false
}
</script>

<style scoped lang="scss">
.image-selector-dialog {
  :deep(.el-dialog__body) {
    padding: 0;
  }
  
  :deep(.el-dialog__header) {
    padding: 20px 24px 16px;
    border-bottom: 1px solid #f0f0f0;
    
    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
      color: #1f2937;
    }
  }
  
  :deep(.el-dialog__footer) {
    padding: 16px 24px 20px;
    border-top: 1px solid #f0f0f0;
  }
}

.image-selector {
  .image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
    max-height: 500px;
    overflow-y: auto;
    padding: 24px;
    background: #fafafa;
    
    &::-webkit-scrollbar {
      width: 6px;
    }
    
    &::-webkit-scrollbar-track {
      background: #f1f1f1;
      border-radius: 3px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: #c1c1c1;
      border-radius: 3px;
      
      &:hover {
        background: #a8a8a8;
      }
    }
  }
  
  .image-item {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    border: 2px solid transparent;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
      border-color: #e5e7eb;
    }
    
    &.selected {
      border-color: #3b82f6;
      box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1), 0 8px 25px rgba(0, 0, 0, 0.15);
      
      .image-overlay {
        opacity: 1;
        background: rgba(59, 130, 246, 0.9);
      }
      
      .check-icon {
        transform: scale(1);
        opacity: 1;
      }
    }
  }
  
  .image-container {
    position: relative;
    aspect-ratio: 1;
    overflow: hidden;
  }
  
  .image-preview {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }
  
  .image-item:hover .image-preview {
    transform: scale(1.05);
  }
  
  .image-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: all 0.3s ease;
    
    .overlay-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      color: white;
      text-align: center;
    }
    
    .check-icon {
      font-size: 32px;
      transform: scale(0.8);
      opacity: 0;
      transition: all 0.3s ease;
    }
    
    .image-name {
      font-size: 12px;
      font-weight: 500;
      max-width: 120px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
  
  .image-info {
    padding: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #f9fafb;
    border-top: 1px solid #f3f4f6;
    
    .image-size,
    .image-dimensions {
      font-size: 11px;
      color: #6b7280;
      font-weight: 500;
    }
  }
  
  .empty-state {
    padding: 60px 24px;
    text-align: center;
    
    .empty-content {
      margin-top: 16px;
    }
    
    .empty-text {
      font-size: 16px;
      color: #374151;
      margin-bottom: 8px;
      font-weight: 500;
    }
    
    .empty-hint {
      font-size: 14px;
      color: #6b7280;
      line-height: 1.5;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  
  .el-button {
    min-width: 80px;
    border-radius: 8px;
    font-weight: 500;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .image-selector-dialog {
    :deep(.el-dialog) {
      width: 95% !important;
      margin: 5vh auto;
    }
  }
  
  .image-selector {
    .image-grid {
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
      gap: 12px;
      padding: 16px;
      max-height: 400px;
    }
    
    .image-info {
      padding: 8px;
      
      .image-size,
      .image-dimensions {
        font-size: 10px;
      }
    }
  }
  
  .dialog-footer {
    flex-direction: column;
    
    .el-button {
      width: 100%;
    }
  }
}

@media (max-width: 480px) {
  .image-selector {
    .image-grid {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
      gap: 8px;
      padding: 12px;
    }
    
    .image-overlay {
      .check-icon {
        font-size: 24px;
      }
      
      .image-name {
        font-size: 10px;
        max-width: 80px;
      }
    }
  }
}
</style>
