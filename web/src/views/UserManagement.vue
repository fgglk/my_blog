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
            <div class="header-content">
              <div class="title-section">
                <h1 class="page-title">用户管理</h1>
                <p class="page-subtitle">管理系统用户账户、权限和角色</p>
              </div>
              
              <!-- 全局操作栏 -->
              <div class="global-actions">
                                                 <el-select v-model="filterStatus" placeholder="所有用户" class="filter-select" @change="handleFilterChange">
                  <el-option label="所有用户" value="all" />
                  <el-option label="活跃用户" value="active" />
                  <el-option label="已禁用用户" value="disabled" />
                </el-select>
                
                <el-button class="filter-btn">
                  <el-icon><Filter /></el-icon>
                  筛选
                </el-button>
                
                <el-button class="export-btn">
                  <el-icon><Download /></el-icon>
                  导出
                </el-button>
                
                <el-button type="primary" class="add-user-btn">
                  <el-icon><Plus /></el-icon>
                  + 添加用户
                </el-button>
              </div>
            </div>
          </div>

          <!-- 统计卡片区域 -->
          <div class="stats-section">
            <div class="stats-grid">
              <div class="stat-card">
                <div class="stat-content">
                  <div class="stat-main">
                    <span class="stat-number">{{ totalUsers }}</span>
                    <div class="stat-trend">
                      <el-icon class="trend-icon up"><ArrowUp /></el-icon>
                      <span class="trend-text">8.2% 较上月</span>
                    </div>
                  </div>
                  <span class="stat-label">总用户数</span>
                </div>
                <div class="stat-icon">
                  <el-icon><User /></el-icon>
                </div>
              </div>

              <div class="stat-card">
                <div class="stat-content">
                  <div class="stat-main">
                    <span class="stat-number">{{ activeUsers }}</span>
                    <div class="stat-trend">
                      <el-icon class="trend-icon up"><ArrowUp /></el-icon>
                      <span class="trend-text">12.5% 较上月</span>
                    </div>
                  </div>
                  <span class="stat-label">活跃用户</span>
                </div>
                <div class="stat-icon active">
                  <el-icon><UserFilled /></el-icon>
                </div>
              </div>

              <div class="stat-card">
                <div class="stat-content">
                  <div class="stat-main">
                    <span class="stat-number">{{ adminUsers }}</span>
                    <div class="stat-trend">
                      <el-icon class="trend-icon down"><ArrowDown /></el-icon>
                      <span class="trend-text">2.1% 较上月</span>
                    </div>
                  </div>
                  <span class="stat-label">管理员</span>
                </div>
                <div class="stat-icon admin">
                  <el-icon><Setting /></el-icon>
                </div>
              </div>


            </div>
          </div>

          <!-- 用户列表区域 -->
          <div class="users-section">
            <div class="section-header">
              <h2 class="section-title">用户列表</h2>
              <div class="table-controls">
                                 <el-select v-model="sortBy" placeholder="最近添加" class="sort-select" @change="handleSortChange">
                   <el-option label="最近添加" value="recent" />
                   <el-option label="用户名" value="username" />
                   <el-option label="注册时间" value="createdAt" />
                   <el-option label="最后登录" value="lastLogin" />
                 </el-select>
                
                <el-select v-model="pageSize" placeholder="10条/页" class="page-size-select">
                  <el-option label="10条/页" :value="10" />
                  <el-option label="20条/页" :value="20" />
                  <el-option label="50条/页" :value="50" />
                  <el-option label="100条/页" :value="100" />
                </el-select>
                
                <el-input
                  v-model="searchKeyword"
                  placeholder="搜索用户..."
                  class="search-input"
                  @keyup.enter="handleSearch"
                  @clear="handleSearch"
                >
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
              </div>
            </div>

            <!-- 用户表格 -->
            <div class="users-table">
              <el-table
                v-loading="loading"
                :data="users"
                style="width: 100%"
                @selection-change="handleSelectionChange"
              >
                <el-table-column type="selection" width="55" />
                
                <el-table-column label="用户信息" min-width="300">
                  <template #default="{ row }">
                    <div class="user-info-cell">
                      <el-avatar :size="40" :src="row.avatar">
                        {{ row.nickname?.charAt(0) || row.username?.charAt(0) || 'U' }}
                      </el-avatar>
                      <div class="user-details">
                        <div class="user-name">{{ row.nickname || row.username }}</div>
                        <div class="user-email">{{ row.email }}</div>
                                                 <div class="user-id">ID: {{ row.uuid }}</div>
                      </div>
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column label="角色" width="120">
                  <template #default="{ row }">
                    <el-tag 
                      :type="getRoleType(row.role)" 
                      size="small"
                      class="role-tag"
                    >
                      {{ getRoleLabel(row.role) }}
                    </el-tag>
                  </template>
                </el-table-column>
                
                <el-table-column label="状态" width="120">
                  <template #default="{ row }">
                    <div class="status-cell">
                      <span class="status-dot" :class="getStatusClass(row)"></span>
                      <span class="status-text">{{ getStatusLabel(row) }}</span>
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column label="最后登录" width="180">
                  <template #default="{ row }">
                    <span class="last-login">
                      {{ row.last_login_at ? formatDate(row.last_login_at) : '-' }}
                    </span>
                  </template>
                </el-table-column>
                
                <el-table-column label="操作" width="200" fixed="right">
                  <template #default="{ row }">
                    <div class="action-buttons">
                      <el-button 
                        type="primary" 
                        size="small" 
                        circle
                        @click="viewUserDetail(row)"
                        title="查看详情"
                      >
                        <el-icon><View /></el-icon>
                      </el-button>
                      
                                             <el-button 
                         type="warning" 
                         size="small" 
                         circle
                         @click="editUser()"
                         title="编辑用户"
                       >
                        <el-icon><Edit /></el-icon>
                      </el-button>
                      
                      <el-button 
                        v-if="row.status === 0"
                        type="success" 
                        size="small" 
                        circle
                        @click="approveUser(row)"
                        title="启用用户"
                      >
                        <el-icon><Check /></el-icon>
                      </el-button>
                      
                      <el-button 
                        v-if="row.status === 1"
                        type="danger" 
                        size="small" 
                        circle
                        @click="rejectUser(row)"
                        title="禁用用户"
                      >
                        <el-icon><Close /></el-icon>
                      </el-button>
                      
                      <el-button 
                        v-if="currentUser?.id !== row.id"
                        type="danger" 
                        size="small" 
                        circle
                        @click="deleteUser(row)"
                        title="删除用户"
                      >
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>

            <!-- 分页 -->
            <div class="pagination-section">
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
            <el-tag :type="getRoleType(selectedUser.role)">
              {{ getRoleLabel(selectedUser.role) }}
            </el-tag>
          </div>
        </div>
        
        <div class="detail-content">
          <div class="detail-item">
            <label>用户ID：</label>
            <span>{{ selectedUser.uuid }}</span>
          </div>
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
            <span>{{ formatDate(selectedUser.created_at) }}</span>
          </div>
          <div class="detail-item">
            <label>最后更新：</label>
            <span>{{ formatDate(selectedUser.updated_at) }}</span>
          </div>
          <div class="detail-item">
            <label>最后登录：</label>
            <span>{{ selectedUser.last_login_at ? formatDate(selectedUser.last_login_at) : '从未登录' }}</span>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Search, User, UserFilled, 
  View, Setting, Delete,
  Filter, Download, Plus, ArrowUp, ArrowDown,
  Edit, Check, Close
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
const pageSize = ref(10)
const searchKeyword = ref('')
const filterStatus = ref('all')
const sortBy = ref('recent')
const selectedRows = ref<UserInfo[]>([])

// 对话框状态
const showUserDetail = ref(false)
const selectedUser = ref<UserInfo | null>(null)

// 计算属性
const adminUsers = computed(() => {
  return users.value.filter(user => user.role === 'admin').length
})

const activeUsers = computed(() => {
  return users.value.filter(user => user.status === 1).length
})



// 获取用户列表
const loadUsers = async () => {
  console.log('开始加载用户列表...')
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      size: pageSize.value
    }
    
    // 添加搜索关键词
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    
    // 添加状态筛选
    if (filterStatus.value !== 'all') {
      // 将前端筛选值转换为后端期望的数字
      switch (filterStatus.value) {
        case 'active':
          params.status = 1
          break
        case 'disabled':
          params.status = 0
          break
        default:
          params.status = filterStatus.value
      }
    }
    
    // 添加排序
    if (sortBy.value !== 'recent') {
      params.sortBy = sortBy.value
    }
    
    console.log('调用用户列表API，参数:', params)
    const response = await userApi.getUserList(params)
    
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

// 监听筛选状态变化
const handleFilterChange = () => {
  currentPage.value = 1
  loadUsers()
}

// 监听排序变化
const handleSortChange = () => {
  currentPage.value = 1
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

// 表格选择处理
const handleSelectionChange = (selection: UserInfo[]) => {
  selectedRows.value = selection
}

// 查看用户详情
const viewUserDetail = (user: UserInfo) => {
  selectedUser.value = user
  showUserDetail.value = true
}

// 编辑用户
const editUser = () => {
  ElMessage.info('编辑用户功能开发中...')
}

// 启用用户
const approveUser = async (user: UserInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要启用用户 "${user.nickname || user.username}" 吗？`,
      '确认启用',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
      }
    )
    
    loading.value = true
    const response = await userApi.approveUser(user.uuid)
    
    if (response.code === 0) {
      ElMessage.success('启用用户成功')
      loadUsers()
    } else {
      ElMessage.error(response.msg || '启用失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('启用用户失败:', error)
      ElMessage.error('启用失败')
    }
  } finally {
    loading.value = false
  }
}

// 禁用用户
const rejectUser = async (user: UserInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要禁用用户 "${user.nickname || user.username}" 吗？`,
      '确认禁用',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    loading.value = true
    const response = await userApi.rejectUser(user.uuid)
    
    if (response.code === 0) {
      ElMessage.success('禁用用户成功')
      loadUsers()
    } else {
      ElMessage.error(response.msg || '禁用失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('禁用用户失败:', error)
      ElMessage.error('禁用失败')
    }
  } finally {
    loading.value = false
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
    
    loading.value = true
    const response = await userApi.deleteUserById(user.id)
    
    if (response.code === 0) {
      ElMessage.success('删除用户成功')
      loadUsers()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error('删除失败')
    }
  } finally {
    loading.value = false
  }
}

// 获取角色类型
const getRoleType = (role: string) => {
  switch (role) {
    case 'admin':
      return 'danger'
    case 'editor':
      return 'warning'
    default:
      return 'info'
  }
}

// 获取角色标签
const getRoleLabel = (role: string) => {
  switch (role) {
    case 'admin':
      return '管理员'
    case 'editor':
      return '编辑'
    default:
      return '普通用户'
  }
}

// 获取状态样式类
const getStatusClass = (user: UserInfo) => {
  switch (user.status) {
    case 1:
      return 'active'
    case 0:
      return 'disabled'
    default:
      return 'active'
  }
}

// 获取状态标签
const getStatusLabel = (user: UserInfo) => {
  switch (user.status) {
    case 1:
      return '活跃'
    case 0:
      return '已禁用'
    default:
      return '活跃'
  }
}

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
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
  background: #f8fafc;
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
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-header {
  margin-bottom: 30px;
  
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 20px;
  }
  
  .title-section {
    .page-title {
      font-size: 28px;
      font-weight: bold;
      color: #1f2937;
      margin: 0 0 8px 0;
    }
    
    .page-subtitle {
      color: #6b7280;
      font-size: 16px;
      margin: 0;
    }
  }
  
  .global-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .filter-select {
      width: 120px;
    }
    
    .filter-btn, .export-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 16px;
      border: 1px solid #d1d5db;
      background: #fff;
      color: #374151;
      
      &:hover {
        background: #f9fafb;
        border-color: #9ca3af;
      }
    }
    
    .add-user-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 16px;
      background: #3b82f6;
      border: none;
      
      &:hover {
        background: #2563eb;
      }
    }
  }
}

.stats-section {
  margin-bottom: 30px;
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 20px;
  }
  
  .stat-card {
    background: #fff;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e5e7eb;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 4px;
      background: linear-gradient(90deg, #3b82f6, #06b6d4);
    }
    
    .stat-content {
      flex: 1;
      
      .stat-main {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 8px;
        
        .stat-number {
          font-size: 32px;
          font-weight: bold;
          color: #1f2937;
          line-height: 1;
        }
        
        .stat-trend {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 14px;
          
          .trend-icon {
            font-size: 16px;
            
            &.up {
              color: #10b981;
            }
            
            &.down {
              color: #ef4444;
            }
          }
          
          .trend-text {
            color: #6b7280;
          }
        }
      }
      
      .stat-label {
        font-size: 14px;
        color: #6b7280;
      }
    }
    
    .stat-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 48px;
      height: 48px;
      border-radius: 12px;
      background: #3b82f6;
      color: #fff;
      
      &.active {
        background: #10b981;
      }
      
      &.admin {
        background: #3b82f6;
      }
      
      &.pending {
        background: #f59e0b;
      }
      
      .el-icon {
        font-size: 20px;
      }
    }
  }
}

.users-section {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #e5e7eb;
  
  .section-title {
    font-size: 20px;
    font-weight: bold;
    color: #1f2937;
    margin: 0;
  }
  
  .table-controls {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .sort-select, .page-size-select {
      width: 120px;
    }
    
    .search-input {
      width: 200px;
    }
  }
}

.users-table {
  .user-info-cell {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .user-details {
      .user-name {
        font-weight: 600;
        color: #1f2937;
        margin-bottom: 4px;
      }
      
      .user-email {
        font-size: 14px;
        color: #6b7280;
        margin-bottom: 2px;
      }
      
      .user-id {
        font-size: 12px;
        color: #9ca3af;
      }
    }
  }
  
  .role-tag {
    font-weight: 500;
  }
  
  .status-cell {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .status-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      
      &.active {
        background: #10b981;
      }
      
      &.pending {
        background: #f59e0b;
      }
      
      &.disabled {
        background: #9ca3af;
      }
    }
    
    .status-text {
      font-size: 14px;
      color: #374151;
    }
  }
  
  .last-login {
    font-size: 14px;
    color: #6b7280;
  }
  
  .action-buttons {
    display: flex;
    gap: 8px;
  }
}

.pagination-section {
  padding: 20px 24px;
  border-top: 1px solid #e5e7eb;
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
    border-bottom: 1px solid #e5e7eb;
    
    .detail-info {
      h3 {
        font-size: 20px;
        font-weight: bold;
        color: #1f2937;
        margin: 0 0 5px 0;
      }
      
      p {
        color: #6b7280;
        margin: 0 0 10px 0;
      }
    }
  }
  
  .detail-content {
    .detail-item {
      display: flex;
      margin-bottom: 15px;
      
      label {
        font-weight: 600;
        color: #374151;
        width: 100px;
        flex-shrink: 0;
      }
      
      span {
        color: #6b7280;
      }
    }
  }
}

// 响应式设计
@media (max-width: 1024px) {
  .page-header .header-content {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
  
  .global-actions {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  }
}

@media (max-width: 768px) {
  .management-container {
    padding: 0 10px;
  }
  
  .page-header .page-title {
    font-size: 24px;
  }
  
  .global-actions {
    flex-direction: column;
    align-items: stretch;
    
    .filter-select, .filter-btn, .export-btn, .add-user-btn {
      width: 100%;
    }
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .section-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
    
    .table-controls {
      flex-wrap: wrap;
      gap: 8px;
      
      .sort-select, .page-size-select, .search-input {
        width: 100%;
      }
    }
  }
  
  .action-buttons {
    flex-wrap: wrap;
  }
}
</style>
