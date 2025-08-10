<template>
  <div class="comment-item">
    <div class="comment-main">
      <div class="comment-header">
        <div class="comment-user">
          <el-avatar :size="32" :src="comment.user_avatar">
            {{ comment.user_name.charAt(0) }}
          </el-avatar>
          <div class="user-info">
            <div class="username">{{ comment.user_name }}</div>
            <div class="comment-time">{{ formatDate(comment.created_at) }}</div>
            <div v-if="comment.parent_user_name" class="reply-info">
              回复 @{{ comment.parent_user_name }}
            </div>
          </div>
        </div>
        <div class="comment-actions">
          <el-button 
            link
            size="small"
            @click="showReplyInput(comment.id, comment.user_name)"
            class="reply-btn"
          >
            <el-icon><ChatDotRound /></el-icon>
            回复
          </el-button>
          <!-- 删除按钮 - 只有评论作者可以删除 -->
          <el-button 
            v-if="canDelete"
            link
            size="small"
            @click="handleDelete"
            class="delete-btn"
          >
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>
      <div class="comment-content">
        {{ comment.content }}
      </div>
      
      <!-- 回复输入框 -->
      <div v-if="replyToComment === comment.id" class="reply-form">
        <div class="reply-input-wrapper">
          <el-avatar :size="24" class="reply-avatar" :src="userStore.userInfo?.avatar">
            {{ userStore.userInfo?.username.charAt(0) || 'U' }}
          </el-avatar>
          <div class="reply-input-content">
            <el-input
              v-model="replyContent"
              type="textarea"
              :rows="2"
              :placeholder="`回复 @${replyToUserName}...`"
              maxlength="300"
              show-word-limit
              class="reply-textarea"
            />
            <div class="reply-actions">
              <el-button 
                link
                size="small"
                @click="cancelReply"
              >
                取消
              </el-button>
              <el-button 
                type="primary" 
                size="small"
                @click="submitReply(comment.id)"
                :loading="replyLoading"
                :disabled="!replyContent.trim()"
              >
                发表回复
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 子评论 -->
    <div v-if="comment.children && comment.children.length > 0" class="comment-replies">
      <div class="replies-toggle" @click="toggleReplies">
        <el-icon class="toggle-icon" :class="{ expanded: isExpanded }">
          <ArrowDown />
        </el-icon>
        <span>{{ isExpanded ? '收起回复' : `展开回复(${getTotalRepliesCount(comment.children)})` }}</span>
      </div>
      
      <div v-if="isExpanded" class="replies-container">
        <div 
          v-for="child in comment.children" 
          :key="child.id"
          class="reply-item"
        >
          <!-- 递归渲染子评论 -->
          <CommentItem 
            :comment="child"
            :user-store="userStore"
            :comment-api="commentApi"
            :article-id="articleId"
            :article-author-id="articleAuthorId"
            @comment-updated="$emit('comment-updated')"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChatDotRound, ArrowDown, Delete } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Comment } from '@/types/comment'
import type { commentApi } from '@/api/comment'
import type { useUserStore } from '@/stores/user'

interface Props {
  comment: Comment
  userStore: ReturnType<typeof useUserStore>
  commentApi: typeof commentApi
  articleId: number
  articleAuthorId?: number // 文章作者ID
}

interface Emits {
  (e: 'comment-updated'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 回复相关变量
const replyToComment = ref<number | null>(null)
const replyToUserName = ref('')
const replyContent = ref('')
const replyLoading = ref(false)

// 展开/折叠状态
const isExpanded = ref(true)

// 删除相关
const deleteLoading = ref(false)

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 计算总回复数（包括嵌套回复）
const getTotalRepliesCount = (children: any[]): number => {
  let total = children.length
  for (const child of children) {
    if (child.children && child.children.length > 0) {
      total += getTotalRepliesCount(child.children)
    }
  }
  return total
}

// 显示回复输入框
const showReplyInput = (commentId: number, userName: string) => {
  replyToComment.value = commentId
  replyToUserName.value = userName
  replyContent.value = ''
}

// 取消回复
const cancelReply = () => {
  replyToComment.value = null
  replyToUserName.value = ''
  replyContent.value = ''
}

// 提交回复
const submitReply = async (commentId: number) => {
  if (!replyContent.value.trim()) return
  
  replyLoading.value = true
  try {
    const response = await props.commentApi.createComment({
      content: replyContent.value,
      article_id: props.articleId,
      parent_id: commentId
    })
    
    if (response.code === 0) {
      ElMessage.success('回复发表成功')
      replyContent.value = ''
      replyToComment.value = null
      replyToUserName.value = ''
      emit('comment-updated')
    } else {
      ElMessage.error(response.msg || '回复发表失败')
    }
  } catch (error) {
    ElMessage.error('回复发表失败')
  } finally {
    replyLoading.value = false
  }
}

// 检查是否可以删除评论（管理员、文章作者或评论作者可以删除）
const canDelete = computed(() => {
  if (!props.userStore.userInfo) return false
  
  const currentUserId = props.userStore.userInfo.id
  const userRole = props.userStore.userInfo.role
  
  // 管理员可以删除任何评论
  if (userRole === 'admin') return true
  
  // 文章作者可以删除任何评论
  if (props.articleAuthorId && currentUserId === props.articleAuthorId) return true
  
  // 评论作者可以删除自己的评论
  if (props.comment.user_id === currentUserId) return true
  
  return false
})

// 删除评论
const handleDelete = async () => {
  if (!props.userStore.userInfo) {
    ElMessage.warning('请先登录')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      '确定要删除这条评论吗？删除后无法恢复。',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    deleteLoading.value = true
    const response = await props.commentApi.deleteComment(props.comment.id)
    
    if (response.code === 0) {
      ElMessage.success('评论删除成功')
      emit('comment-updated')
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  } finally {
    deleteLoading.value = false
  }
}

// 展开/折叠评论
const toggleReplies = () => {
  isExpanded.value = !isExpanded.value
}
</script>

<style lang="scss" scoped>
.comment-item {
  margin-bottom: 20px;
  
  .comment-main {
    background: rgba(255, 255, 255, 0.8);
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  }
  
  .comment-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;
  }
  
  .comment-user {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
    
    .user-info {
      .username {
        font-weight: bold;
        color: #333;
        margin-bottom: 4px;
      }
      
      .comment-time {
        font-size: 12px;
        color: #999;
      }
      
      .reply-info {
        font-size: 12px;
        color: #667eea;
        margin-top: 2px;
        font-style: italic;
      }
    }
  }
  
  .comment-content {
    color: #333;
    line-height: 1.6;
    margin-left: 44px;
  }
  
  .comment-actions {
    display: flex;
    gap: 8px;
    
    .reply-btn {
      color: #667eea;
      font-size: 12px;
      padding: 4px 8px;
      
      &:hover {
        color: #5a6fd8;
        background: rgba(102, 126, 234, 0.1);
      }
      
      .el-icon {
        margin-right: 4px;
      }
    }
    
    .delete-btn {
      color: #f56c6c;
      font-size: 12px;
      padding: 4px 8px;
      
      &:hover {
        color: #e74c3c;
        background: rgba(245, 108, 108, 0.1);
      }
      
      .el-icon {
        margin-right: 4px;
      }
    }
  }
  
  .reply-form {
    margin-top: 15px;
    margin-left: 44px;
    
    .reply-input-wrapper {
      display: flex;
      gap: 12px;
      
      .reply-input-content {
        flex: 1;
        
        .reply-textarea {
          margin-bottom: 10px;
        }
        
        .reply-actions {
          display: flex;
          gap: 10px;
          justify-content: flex-end;
        }
      }
    }
  }
  
  .comment-replies {
    margin-top: 15px;
    margin-left: 44px;
    
    .replies-toggle {
      display: flex;
      align-items: center;
      gap: 6px;
      color: #667eea;
      font-size: 14px;
      cursor: pointer;
      padding: 8px 0;
      transition: all 0.3s ease;
      
      &:hover {
        color: #5a6fd8;
      }
      
      .toggle-icon {
        transition: transform 0.3s ease;
        
        &.expanded {
          transform: rotate(180deg);
        }
      }
    }
    
    .replies-container {
      margin-top: 10px;
      
      .reply-item {
        margin-bottom: 15px;
        padding-left: 15px;
        border-left: 2px solid rgba(102, 126, 234, 0.1);
        
        &:last-child {
          margin-bottom: 0;
        }
      }
    }
  }
}

// 嵌套回复样式
.comment-item .comment-item {
  .comment-main {
    background: rgba(255, 255, 255, 0.6);
    padding: 15px;
    
    .comment-user {
      .user-info {
        .username {
          font-size: 14px;
        }
        
        .comment-time {
          font-size: 11px;
        }
        
        .reply-info {
          font-size: 11px;
        }
      }
    }
    
    .comment-content {
      font-size: 14px;
      margin-left: 36px;
    }
    
    .reply-form {
      margin-left: 36px;
    }
  }
  
  .comment-replies {
    margin-left: 36px;
    
    .replies-container {
      .reply-item {
        padding-left: 12px;
        border-left: 1px dashed rgba(102, 126, 234, 0.15);
      }
    }
  }
}

// 更深层嵌套
.comment-item .comment-item .comment-item {
  .comment-main {
    background: rgba(255, 255, 255, 0.4);
    padding: 12px;
    
    .comment-user {
      .user-info {
        .username {
          font-size: 13px;
        }
        
        .comment-time {
          font-size: 10px;
        }
        
        .reply-info {
          font-size: 10px;
        }
      }
    }
    
    .comment-content {
      font-size: 13px;
      margin-left: 28px;
    }
    
    .reply-form {
      margin-left: 28px;
    }
  }
  
  .comment-replies {
    margin-left: 28px;
    
    .replies-container {
      .reply-item {
        padding-left: 10px;
        border-left: 1px dotted rgba(102, 126, 234, 0.1);
      }
    }
  }
}
</style> 