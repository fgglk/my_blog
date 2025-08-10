<template>
  <div class="login-page">
    <!-- 背景动画元素 -->
    <div class="background-animation">
      <div class="floating-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
        <div class="shape shape-5"></div>
      </div>
    </div>

    <div class="login-container">
      <div class="login-card">
        <!-- 装饰性头部 -->
        <div class="login-header">
          <div class="logo-section">
            <div class="logo-icon">
              <el-icon size="32"><User /></el-icon>
            </div>
            <h2 class="gradient-text">欢迎回来</h2>
            <p class="subtitle">登录您的博客账户</p>
          </div>
        </div>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          @submit.prevent="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              size="large"
              class="custom-input"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              class="custom-input"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="captcha">
            <div class="captcha-container">
              <el-input
                v-model="loginForm.captcha"
                placeholder="请输入验证码"
                size="large"
                class="custom-input captcha-input"
              >
                <template #prefix>
                  <el-icon><Key /></el-icon>
                </template>
              </el-input>
              <div class="captcha-wrapper">
                <img
                  v-if="captchaUrl"
                  :src="captchaUrl"
                  alt="验证码"
                  class="captcha-image"
                  @click="refreshCaptcha"
                  @error="handleImageError"
                  @load="handleImageLoad"
                />
                <div v-else class="captcha-placeholder" @click="refreshCaptcha">
                  <el-icon><Refresh /></el-icon>
                  <span>点击刷新</span>
                </div>
              </div>
            </div>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              size="large"
              class="login-button"
              :loading="loading"
              @click="handleLogin"
            >
              <el-icon v-if="!loading"><Right /></el-icon>
              {{ loading ? '登录中...' : '立即登录' }}
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-footer">
          <div class="links-section">
            <router-link to="/register" class="link-item">
              <el-icon><UserFilled /></el-icon>
              <span>还没有账号？立即注册</span>
            </router-link>
            <router-link to="/forgot-password" class="link-item">
              <el-icon><QuestionFilled /></el-icon>
              <span>忘记密码？</span>
            </router-link>
          </div>
        </div>

        <!-- 装饰性底部 -->
        <div class="decoration-bottom">
          <div class="decoration-line"></div>
          <div class="decoration-dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user'
import { ElMessage } from 'element-plus'
import { 
  User, Lock, Key, Right, UserFilled, QuestionFilled, Refresh 
} from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const captchaUrl = ref('')
const captchaId = ref('')

const loginForm = reactive({
  username: '',
  password: '',
  captcha: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 1, max: 20, message: '用户名长度在 1 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为 6 位', trigger: 'blur' }
  ]
}

// 获取验证码
const getCaptcha = async () => {
  try {
    console.log('开始获取验证码...')
    const response = await userApi.getCaptcha()
    console.log('验证码API响应:', response)
    
    if (response.code === 0) {
      captchaUrl.value = response.data.image
      captchaId.value = response.data.captcha_id
      console.log('验证码URL:', captchaUrl.value)
      console.log('验证码ID:', captchaId.value)
    } else {
      console.error('获取验证码失败:', response.msg)
      ElMessage.error('获取验证码失败，请重试')
    }
  } catch (error) {
    console.error('获取验证码失败:', error)
    ElMessage.error('获取验证码失败，请重试')
  }
}

// 刷新验证码
const refreshCaptcha = () => {
  getCaptcha()
  loginForm.captcha = ''
}

// 处理图片加载错误
const handleImageError = (event: Event) => {
  console.error('验证码图片加载失败:', event)
  ElMessage.error('验证码图片加载失败，请刷新重试')
}

// 处理图片加载成功
const handleImageLoad = () => {
  console.log('验证码图片加载成功')
}

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    console.log('登录表单数据:', {
      username: loginForm.username,
      password: loginForm.password,
      captcha: loginForm.captcha,
      captchaId: captchaId.value
    })
    
    const result = await userStore.login(loginForm.username, loginForm.password, loginForm.captcha, captchaId.value)
    
    if (result.success) {
      ElMessage.success('登录成功')
      router.push('/')
    } else {
      ElMessage.error(result.message || '登录失败')
      refreshCaptcha()
    }
  } catch (error) {
    console.error('登录失败:', error)
    ElMessage.error('登录失败')
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getCaptcha()
})
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

// 背景动画
.background-animation {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1;
}

.floating-shapes {
  position: relative;
  width: 100%;
  height: 100%;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
  
  &.shape-1 {
    width: 80px;
    height: 80px;
    top: 20%;
    left: 10%;
    animation-delay: 0s;
  }
  
  &.shape-2 {
    width: 120px;
    height: 120px;
    top: 60%;
    right: 15%;
    animation-delay: 2s;
  }
  
  &.shape-3 {
    width: 60px;
    height: 60px;
    bottom: 20%;
    left: 20%;
    animation-delay: 4s;
  }
  
  &.shape-4 {
    width: 100px;
    height: 100px;
    top: 30%;
    right: 30%;
    animation-delay: 1s;
  }
  
  &.shape-5 {
    width: 40px;
    height: 40px;
    bottom: 40%;
    right: 10%;
    animation-delay: 3s;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

.login-container {
  width: 100%;
  max-width: 450px;
  position: relative;
  z-index: 2;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 50px 40px;
  box-shadow: 
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2);
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(90deg, #667eea, #764ba2, #f093fb);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-section {
  .logo-icon {
    width: 80px;
    height: 80px;
    background: linear-gradient(135deg, #667eea, #764ba2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    color: white;
    box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
    animation: pulse 2s ease-in-out infinite;
  }
  
  .gradient-text {
    font-size: 32px;
    font-weight: bold;
    background: linear-gradient(135deg, #667eea, #764ba2);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    margin-bottom: 8px;
  }
  
  .subtitle {
    color: #666;
    font-size: 16px;
    margin: 0;
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.login-form {
  .el-form-item {
    margin-bottom: 24px;
  }
  
  .custom-input {
    :deep(.el-input__wrapper) {
      background: rgba(255, 255, 255, 0.8);
      border: 2px solid rgba(102, 126, 234, 0.1);
      border-radius: 12px;
      transition: all 0.3s ease;
      
      &:hover {
        border-color: rgba(102, 126, 234, 0.3);
      }
      
      &.is-focus {
        border-color: #667eea;
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
      }
    }
    
    :deep(.el-input__inner) {
      height: 48px;
      font-size: 16px;
    }
    
    :deep(.el-input__prefix) {
      color: #667eea;
      font-size: 18px;
    }
  }
  
  .login-button {
    width: 100%;
    height: 52px;
    font-size: 18px;
    font-weight: bold;
    border-radius: 12px;
    background: linear-gradient(135deg, #667eea, #764ba2);
    border: none;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
}

.captcha-container {
  display: flex;
  gap: 12px;
  align-items: center;
  
  .captcha-input {
    flex: 1;
  }
  
  .captcha-wrapper {
    width: 120px;
    height: 48px;
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      transform: scale(1.02);
    }
  }
  
  .captcha-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border: 2px solid rgba(102, 126, 234, 0.1);
    border-radius: 12px;
    background: #fff;
    cursor: pointer;
  }
  
  .captcha-placeholder {
    width: 100%;
    height: 100%;
    background: rgba(102, 126, 234, 0.1);
    border: 2px dashed rgba(102, 126, 234, 0.3);
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #667eea;
    font-size: 12px;
    
    .el-icon {
      font-size: 16px;
      margin-bottom: 2px;
    }
  }
}

.login-footer {
  margin-top: 30px;
  
  .links-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .link-item {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #667eea;
    text-decoration: none;
    font-size: 14px;
    padding: 8px 12px;
    border-radius: 8px;
    transition: all 0.3s ease;
    
    &:hover {
      background: rgba(102, 126, 234, 0.1);
      transform: translateX(5px);
    }
    
    .el-icon {
      font-size: 16px;
    }
  }
}

.decoration-bottom {
  margin-top: 30px;
  text-align: center;
  
  .decoration-line {
    height: 2px;
    background: linear-gradient(90deg, transparent, #667eea, transparent);
    margin-bottom: 15px;
  }
  
  .decoration-dots {
    display: flex;
    justify-content: center;
    gap: 8px;
    
    span {
      width: 6px;
      height: 6px;
      border-radius: 50%;
      background: #667eea;
      animation: dot-pulse 2s ease-in-out infinite;
      
      &:nth-child(2) {
        animation-delay: 0.3s;
      }
      
      &:nth-child(3) {
        animation-delay: 0.6s;
      }
    }
  }
}

@keyframes dot-pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .login-card {
    padding: 40px 30px;
  }
  
  .logo-section {
    .logo-icon {
      width: 60px;
      height: 60px;
    }
    
    .gradient-text {
      font-size: 28px;
    }
  }
  
  .captcha-container {
    flex-direction: column;
    gap: 8px;
    
    .captcha-wrapper {
      width: 100%;
      height: 40px;
    }
  }
}

@media (max-width: 480px) {
  .login-card {
    padding: 30px 20px;
  }
  
  .logo-section {
    .gradient-text {
      font-size: 24px;
    }
    
    .subtitle {
      font-size: 14px;
    }
  }
  
  .login-form {
    .custom-input :deep(.el-input__inner) {
      height: 44px;
      font-size: 14px;
    }
    
    .login-button {
      height: 48px;
      font-size: 16px;
    }
  }
}
</style> 