<template>
  <div class="user-management-page">
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
            <router-link to="/write" class="nav-item">写文章</router-link>
            <router-link to="/profile" class="nav-item">个人中心</router-link>
            <el-button @click="handleLogout" link>退出</el-button>
          </nav>
        </div>
      </div>
    </header>

    <main class="main">
      <div class="container">
        <div class="management-container">
          <!-- 页面标题 -->
          <div class="page-header">
            <h1 class="page-title">用户管理</h1>
            <p class="page-subtitle">管理系统中的所有用户账户</p>
          </div>

          <!-- 搜索和筛选区域 -->
          <div class="search-section">
            <div class="search-box">
              <el-input
                v-model="searchKeyword"
                placeholder="搜索用户名、邮箱或昵称"
                clearable
                @keyup.enter="handleSearch"
                @clear="handleSearch"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-button type="primary" @click="handleSearch">
                <el-icon><Search /></el-icon>
                搜索
              </el-button>
            </div>
          </div>

          <!-- 用户统计 -->
          <div class="stats-section">
            <div class="stats-grid">
              <div class="stat-card">
                <div class="stat-icon">
                  <el-icon><User /></el-icon>
                </div>
                <div class="stat-content">
                  <span class="stat-number">{{ totalUsers }}</span>
                  <span class="stat-label">总用户数</span>
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">
                  <el-icon><UserFilled /></el-icon>
                </div>
                <div class="stat-content">
                  <span class="stat-number">{{ adminUsers }}</span>
                  <span class="stat-label">管理员</span>
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">
                  <el-icon><Clock /></el-icon>
                </div>
                <div class="stat-content">
                  <span class="stat-number">{{ recentUsers }}</span>
                  <span class="stat-label">本月新增</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 用户列表 -->
          <div class="users-section">
            <div class="section-header">
              <h2 class="section-title">用户列表</h2>
              <div class="section-actions">
                <el-button type="primary" @click="refreshUsers">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </div>
            </div>

            <!-- 加载状态 -->
            <div v-if="loading" class="loading-section">
              <el-skeleton :rows="5" animated />
            </div>

            <!-- 用户列表 -->
            <div v-else-if="users.length > 0" class="users-list">
              <div class="user-card" v-for="user in users" :key="user.id">
                <div class="user-avatar">
                  <el-avatar :size="60" :src="user.avatar">
                    {{ user.nickname?.charAt(0) || user.username?.charAt(0) || 'U' }}
                  </el-avatar>
                                     <div v-if="user.role === 'admin'" class="admin-badge">
                     <el-icon><Star /></el-icon>
                   </div>
                </div>
                
                <div class="user-info">
                  <div class="user-header">
                    <h3 class="user-name">{{ user.nickname || user.username }}</h3>
                    <div class="user-role">
                      <el-tag :type="user.role === 'admin' ? 'danger' : 'info'" size="small">
                        {{ user.role === 'admin' ? '管理员' : '普通用户' }}
                      </el-tag>
                    </div>
                  </div>
                  
                  <div class="user-details">
                    <p class="user-username">@{{ user.username }}</p>
                    <p class="user-email">{{ user.email }}</p>
                    <p v-if="user.bio" class="user-bio">{{ user.bio }}</p>
                  </div>
                  
                  <div class="user-meta">
                    <span class="meta-item">
                      <el-icon><Calendar /></el-icon>
                      注册时间：{{ formatDate(user.createdAt) }}
                    </span>
                    <span class="meta-item">
                      <el-icon><Clock /></el-icon>
                      最后更新：{{ formatDate(user.updatedAt) }}
                    </span>
                  </div>
                </div>
                
                <div class="user-actions">
                  <el-button type="primary" size="small" @click="viewUserDetail(user)">
                    <el-icon><View /></el-icon>
                    查看详情
                  </el-button>
                  <el-button 
                    v-if="user.role !== 'admin' || currentUser.id !== user.id"
                    type="warning" 
                    size="small" 
                    @click="toggleUserRole(user)"
                  >
                    <el-icon><Setting /></el-icon>
                    {{ user.role === 'admin' ? '取消管理员' : '设为管理员' }}
                  </el-button>
                  <el-button 
                    v-if="currentUser.id !== user.id"
                    type="danger" 
                    size="small" 
                    @click="deleteUser(user)"
                  >
                    <el-icon><Delete /></el-icon>
                    删除用户
                  </el-button>
                </div>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-else class="empty-section">
              <el-empty description="暂无用户数据">
                <el-button type="primary" @click="refreshUsers">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </el-empty>
            </div>

            <!-- 分页 -->
            <div v-if="totalUsers > 0" class="pagination-section">
              <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :total="totalUsers"
                :page-sizes="[10, 20, 50, 100]"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 用户详情对话框 -->
    <el-dialog v-model="showUserDetail" title="用户详情" width="600px">
      <div v-if="selectedUser" class="user-detail">
        <div class="detail-header">
          <el-avatar :size="80" :src="selectedUser.avatar">
            {{ selectedUser.nickname?.charAt(0) || selectedUser.username?.charAt(0) || 'U' }}
          </el-avatar>
          <div class="detail-info">
            <h3>{{ selectedUser.nickname || selectedUser.username }}</h3>
            <p>@{{ selectedUser.username }}</p>
            <el-tag :type="selectedUser.role === 'admin' ? 'danger' : 'info'">
              {{ selectedUser.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </div>
        </div>
        
        <div class="detail-content">
          <div class="detail-item">
            <label>邮箱：</label>
            <span>{{ selectedUser.email }}</span>
          </div>
          <div class="detail-item">
            <label>个人简介：</label>
            <span>{{ selectedUser.bio || '暂无简介' }}</span>
          </div>
          <div class="detail-item">
            <label>注册时间：</label>
            <span>{{ formatDate(selectedUser.createdAt) }}</span>
          </div>
          <div class="detail-item">
            <label>最后更新：</label>
            <span>{{ formatDate(selectedUser.updatedAt) }}</span>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Search, User, UserFilled, Clock, Refresh, 
  Calendar, View, Setting, Delete, Star
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { UserInfo, UserListResponse } from '@/types/user'

const router = useRouter()
const userStore = useUserStore()

// 当前用户信息
const currentUser = computed(() => userStore.userInfo)

// 响应式数据
const loading = ref(false)
const users = ref<UserInfo[]>([])
const totalUsers = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')

// 对话框状态
const showUserDetail = ref(false)
const selectedUser = ref<UserInfo | null>(null)

// 计算属性
const adminUsers = computed(() => {
  return users.value.filter(user => user.role === 'admin').length
})

const recentUsers = computed(() => {
  const now = dayjs()
  const thisMonth = now.startOf('month')
  return users.value.filter(user => {
    return dayjs(user.createdAt).isAfter(thisMonth)
  }).length
})

// 获取用户列表
const loadUsers = async () => {
  console.log('开始加载用户列表...')
  loading.value = true
  try {
    console.log('调用用户列表API，参数:', { page: currentPage.value, size: pageSize.value })
    const response = await userApi.getUserList({
      page: currentPage.value,
      size: pageSize.value
    })
    
    console.log('用户列表API响应:', response)
    
    if (response.code === 0) {
      const data = response.data as UserListResponse
      users.value = data.list
      totalUsers.value = data.total
      console.log('用户列表加载成功，用户数量:', data.list.length)
    } else {
      console.error('获取用户列表失败:', response.msg)
      ElMessage.error(response.msg || '获取用户列表失败')
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索用户
const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
}

// 刷新用户列表
const refreshUsers = () => {
  loadUsers()
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadUsers()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadUsers()
}

// 查看用户详情
const viewUserDetail = (user: UserInfo) => {
  selectedUser.value = user
  showUserDetail.value = true
}

// 切换用户角色
const toggleUserRole = async (user: UserInfo) => {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  const action = newRole === 'admin' ? '设为管理员' : '取消管理员权限'
  
  try {
    await ElMessageBox.confirm(
      `确定要${action}用户 "${user.nickname || user.username}" 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // TODO: 调用后端API切换用户角色
    ElMessage.success(`${action}成功`)
    loadUsers() // 重新加载用户列表
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

// 删除用户
const deleteUser = async (user: UserInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.nickname || user.username}" 吗？此操作不可恢复！`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // TODO: 调用后端API删除用户
    ElMessage.success('删除用户成功')
    loadUsers() // 重新加载用户列表
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 退出登录
const handleLogout = () => {
  userStore.logout()
  ElMessage.success('退出成功')
  router.push('/')
}

// 页面加载时获取用户列表
onMounted(() => {
  console.log('用户管理页面已加载')
  console.log('当前用户信息:', currentUser.value)
  loadUsers()
})
</script>

<style lang="scss" scoped>
.user-management-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
}

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.08);
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
}

.main {
  padding: 20px 0;
}

.management-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
  
  .page-title {
    font-size: 32px;
    font-weight: bold;
    color: #333;
    margin-bottom: 10px;
  }
  
  .page-subtitle {
    color: #666;
    font-size: 16px;
  }
}

.search-section {
  margin-bottom: 30px;
  
  .search-box {
    display: flex;
    gap: 15px;
    max-width: 500px;
    margin: 0 auto;
    
    .el-input {
      flex: 1;
    }
  }
}

.stats-section {
  margin-bottom: 30px;
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
  }
  
  .stat-card {
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    border-radius: 15px;
    padding: 25px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.3);
    display: flex;
    align-items: center;
    gap: 15px;
    
    .stat-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 50px;
      height: 50px;
      border-radius: 50%;
      background: linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%);
      color: #fff;
      
      .el-icon {
        font-size: 24px;
      }
    }
    
    .stat-content {
      display: flex;
      flex-direction: column;
      gap: 5px;
      
      .stat-number {
        font-size: 28px;
        font-weight: bold;
        color: #333;
        line-height: 1;
      }
      
      .stat-label {
        font-size: 14px;
        color: #666;
      }
    }
  }
}

.users-section {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  
  .section-title {
    font-size: 24px;
    font-weight: bold;
    color: #333;
  }
}

.loading-section {
  padding: 40px 0;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 20px;
  border-radius: 15px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  }
}

.user-avatar {
  position: relative;
  
  .admin-badge {
    position: absolute;
    top: -5px;
    right: -5px;
    background: #f59e0b;
    color: #fff;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
  }
}

.user-info {
  flex: 1;
  
  .user-header {
    display: flex;
    align-items: center;
    gap: 15px;
    margin-bottom: 10px;
    
    .user-name {
      font-size: 18px;
      font-weight: bold;
      color: #333;
      margin: 0;
    }
  }
  
  .user-details {
    margin-bottom: 10px;
    
    .user-username {
      color: #3b82f6;
      font-size: 14px;
      margin: 0 0 5px 0;
    }
    
    .user-email {
      color: #666;
      font-size: 14px;
      margin: 0 0 5px 0;
    }
    
    .user-bio {
      color: #666;
      font-size: 14px;
      margin: 0;
      font-style: italic;
    }
  }
  
  .user-meta {
    display: flex;
    gap: 20px;
    
    .meta-item {
      display: flex;
      align-items: center;
      gap: 5px;
      color: #666;
      font-size: 12px;
      
      .el-icon {
        font-size: 14px;
      }
    }
  }
}

.user-actions {
  display: flex;
  gap: 10px;
}

.empty-section {
  padding: 60px 0;
  text-align: center;
}

.pagination-section {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

// 用户详情对话框样式
.user-detail {
  .detail-header {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 25px;
    padding-bottom: 20px;
    border-bottom: 1px solid #e2e8f0;
    
    .detail-info {
      h3 {
        font-size: 20px;
        font-weight: bold;
        color: #333;
        margin: 0 0 5px 0;
      }
      
      p {
        color: #666;
        margin: 0 0 10px 0;
      }
    }
  }
  
  .detail-content {
    .detail-item {
      display: flex;
      margin-bottom: 15px;
      
      label {
        font-weight: bold;
        color: #333;
        width: 100px;
        flex-shrink: 0;
      }
      
      span {
        color: #666;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .management-container {
    padding: 0 10px;
  }
  
  .page-header .page-title {
    font-size: 24px;
  }
  
  .search-box {
    flex-direction: column;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .user-card {
    flex-direction: column;
    text-align: center;
    
    .user-actions {
      justify-content: center;
      flex-wrap: wrap;
    }
  }
  
  .user-meta {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
