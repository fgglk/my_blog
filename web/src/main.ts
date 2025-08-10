import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import 'element-plus/dist/index.css'
import './styles/index.scss'
import { useUserStore } from './stores/user'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

// 初始化用户状态
const userStore = useUserStore()
userStore.initUserState() 