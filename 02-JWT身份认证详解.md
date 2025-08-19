# JWT 身份认证详解 - 从原理到实现

## 问题背景

在博客项目中，用户登录后需要保持登录状态，后续请求需要携带身份凭证。传统的 Session 方式在前后端分离架构中存在问题，因此采用 JWT (JSON Web Token) 进行身份认证。

## 什么是 JWT？

### JWT 基本概念

JWT (JSON Web Token) 是一种开放标准 (RFC 7519)，用于在各方之间安全地传输信息作为 JSON 对象。JWT 可以被验证和信任，因为它是数字签名的。

### JWT 的组成结构

JWT 由三部分组成，用点 (.) 分隔：

```
Header.Payload.Signature
```

#### 1. Header (头部)

包含令牌的类型和使用的签名算法：

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

#### 2. Payload (载荷)

包含声明 (claims)，声明是关于实体和其他数据的声明：

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022,
  "exp": 1516242622
}
```

**标准声明：**
- `iss` (issuer): 签发者
- `sub` (subject): 主题
- `aud` (audience): 受众
- `exp` (expiration time): 过期时间
- `nbf` (not before): 生效时间
- `iat` (issued at): 签发时间
- `jti` (JWT ID): JWT ID

#### 3. Signature (签名)

用于验证消息在传输过程中没有被篡改：

```
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret
)
```

## JWT 的优势和劣势

### 优势

1. **无状态**: 服务器不需要存储会话信息
2. **可扩展性**: 适合分布式系统
3. **跨域友好**: 可以在不同域名间使用
4. **标准化**: 遵循 RFC 7519 标准

### 劣势

1. **无法撤销**: 一旦签发，在过期前无法撤销
2. **存储问题**: 客户端存储可能存在安全风险
3. **大小限制**: 每次请求都会携带完整 token
4. **性能开销**: 每次请求都需要验证签名

## 后端实现

### 1. JWT 工具函数

```go
// utils/jwt.go
package utils

import (
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

type JWTConfig struct {
    SecretKey string        `mapstructure:"secret_key"`
    Expire    time.Duration `mapstructure:"expire"`
}

var jwtConfig JWTConfig

// 初始化 JWT 配置
func InitJWT(config JWTConfig) {
    jwtConfig = config
}

// 生成 JWT Token
func GenerateToken(userID uint, username, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.Expire)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "blog-system",
            Subject:   "user-token",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtConfig.SecretKey))
}

// 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtConfig.SecretKey), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### 2. JWT 中间件

```go
// middleware/jwt.go
package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "server/utils"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从请求头获取 token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": 401,
                "message": "Authorization header is required",
            })
            c.Abort()
            return
        }
        
        // 检查 Bearer 前缀
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": 401,
                "message": "Authorization header format must be Bearer {token}",
            })
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // 解析 token
        claims, err := utils.ParseToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": 401,
                "message": "Invalid token: " + err.Error(),
            })
            c.Abort()
            return
        }
        
        // 将用户信息存储到上下文中
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}
```

## 前端实现

### 1. 请求拦截器

```typescript
// utils/request.ts
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request: AxiosInstance = axios.create({
    baseURL: '/api',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
})

// 请求拦截器
request.interceptors.request.use(
    (config: AxiosRequestConfig) => {
        const token = localStorage.getItem('token')
        if (token && config.headers) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    (response: AxiosResponse) => {
        const { code, message, data } = response.data
        
        if (code === 200) {
            return data
        } else {
            ElMessage.error(message || '请求失败')
            return Promise.reject(new Error(message || '请求失败'))
        }
    },
    (error) => {
        if (error.response) {
            const { status, data } = error.response
            
            switch (status) {
                case 401:
                    // Token 过期或无效
                    localStorage.removeItem('token')
                    localStorage.removeItem('userInfo')
                    router.push('/login')
                    ElMessage.error('登录已过期，请重新登录')
                    break
                case 403:
                    ElMessage.error('没有权限访问')
                    break
                case 404:
                    ElMessage.error('请求的资源不存在')
                    break
                case 500:
                    ElMessage.error('服务器内部错误')
                    break
                default:
                    ElMessage.error(data.message || '请求失败')
            }
        } else {
            ElMessage.error('网络错误，请检查网络连接')
        }
        
        return Promise.reject(error)
    }
)

export default request
```

### 2. 用户状态管理

```typescript
// stores/user.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, getUserInfo } from '@/api/user'
import type { UserInfo, LoginForm } from '@/types/user'

export const useUserStore = defineStore('user', () => {
    const token = ref<string>(localStorage.getItem('token') || '')
    const userInfo = ref<UserInfo | null>(null)
    
    // 计算属性
    const isLoggedIn = computed(() => !!token.value)
    const hasRole = computed(() => (role: string) => userInfo.value?.role === role)
    
    // 登录
    const loginAction = async (loginForm: LoginForm) => {
        try {
            const response = await login(loginForm)
            const { token: newToken, user } = response
            
            // 保存 token 和用户信息
            token.value = newToken
            userInfo.value = user
            localStorage.setItem('token', newToken)
            localStorage.setItem('userInfo', JSON.stringify(user))
            
            return response
        } catch (error) {
            throw error
        }
    }
    
    // 登出
    const logout = () => {
        token.value = ''
        userInfo.value = null
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
    }
    
    return {
        token,
        userInfo,
        isLoggedIn,
        hasRole,
        loginAction,
        logout,
    }
})
```

## 安全考虑

### 1. Token 存储安全

```typescript
// 使用 httpOnly cookie 存储 token（推荐）
// 或者使用 sessionStorage（页面关闭后清除）
sessionStorage.setItem('token', token)

// 避免在 localStorage 中存储敏感信息
// 如果必须使用 localStorage，考虑加密存储
const encryptedToken = btoa(token) // 简单 base64 编码
localStorage.setItem('token', encryptedToken)
```

### 2. Token 过期处理

```typescript
// 自动刷新 token
const refreshTokenIfNeeded = async () => {
    const token = localStorage.getItem('token')
    if (!token) return
    
    try {
        const decoded = jwtDecode(token)
        const now = Date.now() / 1000
        
        // 如果 token 将在 5 分钟内过期，则刷新
        if (decoded.exp && decoded.exp - now < 300) {
            const response = await refreshToken()
            localStorage.setItem('token', response.token)
        }
    } catch (error) {
        console.error('Failed to refresh token:', error)
    }
}

// 定期检查 token
setInterval(refreshTokenIfNeeded, 60000) // 每分钟检查一次
```

## 总结

JWT 是现代 Web 应用中广泛使用的身份认证方案，特别适合前后端分离架构。通过合理的实现和配置，可以构建安全、高效的身份认证系统。

**关键要点：**
- 理解 JWT 的结构和原理
- 正确实现 token 的生成和验证
- 合理处理 token 的存储和传输
- 注意安全性和性能优化
- 完善的错误处理和用户体验

---

*本文详细介绍了 JWT 身份认证的原理、实现和最佳实践，希望对您的开发工作有所帮助。*
