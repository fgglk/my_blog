# 我的博客前端

基于 Vue 3 + TypeScript + Element Plus 的现代化博客前端项目。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - JavaScript 的超集，提供类型安全
- **Vite** - 现代前端构建工具
- **Pinia** - Vue 3 的状态管理库
- **Vue Router** - Vue.js 官方路由管理器
- **Element Plus** - 基于 Vue 3 的 UI 组件库
- **Axios** - HTTP 客户端
- **Sass** - CSS 预处理器

## 功能特性

- 🏠 **首页** - 文章列表展示，支持分页
- 🔍 **搜索** - 文章搜索功能
- 📝 **写文章** - Markdown 编辑器，支持草稿保存
- 👤 **用户系统** - 登录、注册、个人中心
- 💬 **评论系统** - 文章评论功能
- ❤️ **互动功能** - 点赞、收藏
- 📱 **响应式设计** - 支持移动端

## 项目结构

```
web/
├── src/
│   ├── api/           # API 请求
│   ├── components/    # 公共组件
│   ├── router/        # 路由配置
│   ├── stores/        # 状态管理
│   ├── styles/        # 样式文件
│   ├── types/         # TypeScript 类型定义
│   ├── utils/         # 工具函数
│   ├── views/         # 页面组件
│   ├── App.vue        # 根组件
│   └── main.ts        # 入口文件
├── public/            # 静态资源
├── index.html         # HTML 模板
├── package.json       # 项目配置
├── vite.config.ts     # Vite 配置
└── tsconfig.json      # TypeScript 配置
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

项目将在 `http://localhost:3000` 启动。

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 环境配置

确保后端服务器运行在 `http://localhost:8080`，或者修改 `vite.config.ts` 中的代理配置：

```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080', // 修改为你的后端地址
      changeOrigin: true
    }
  }
}
```

## 主要页面

### 首页 (`/`)
- 文章列表展示
- 分类和标签筛选
- 分页功能

### 搜索页 (`/search`)
- 文章搜索
- 搜索结果展示

### 登录页 (`/login`)
- 用户登录
- 验证码功能

### 注册页 (`/register`)
- 用户注册
- 邮箱验证码

### 文章详情页 (`/article/:id`)
- 文章内容展示
- Markdown 渲染
- 评论系统
- 点赞收藏

### 写文章页 (`/write`)
- Markdown 编辑器
- 分类标签选择
- 草稿保存

### 个人中心 (`/profile`)
- 用户信息展示
- 资料编辑
- 密码修改

## API 接口

项目使用 RESTful API 与后端通信，主要接口包括：

- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `POST /api/articles` - 创建文章
- `PUT /api/articles/:id` - 更新文章
- `DELETE /api/articles/:id` - 删除文章
- `POST /api/users/login` - 用户登录
- `POST /api/users/register` - 用户注册
- `GET /api/users/info` - 获取用户信息

## 开发指南

### 添加新页面

1. 在 `src/views/` 目录下创建新的 Vue 组件
2. 在 `src/router/index.ts` 中添加路由配置
3. 在 `src/types/` 中添加相关类型定义

### 添加新 API

1. 在 `src/api/` 目录下创建或修改 API 文件
2. 在 `src/types/` 中添加请求和响应类型
3. 在 `src/stores/` 中添加状态管理

### 样式开发

项目使用 Sass 预处理器，支持嵌套语法和变量：

```scss
// 使用变量
$primary-color: #409eff;

.button {
  background: $primary-color;
  
  &:hover {
    background: darken($primary-color, 10%);
  }
}
```

## 部署

### 构建项目

```bash
npm run build
```

### 部署到服务器

将 `dist` 目录下的文件部署到 Web 服务器，如 Nginx：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 注意事项

1. 确保后端 API 服务器正在运行
2. 检查 CORS 配置是否正确
3. 确保 JWT token 在请求头中正确传递
4. 移动端适配可能需要进一步优化

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 许可证

MIT License 