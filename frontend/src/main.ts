import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import './style.css'
import App from './App.vue'
import Home from './views/Home.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/quote/:id', name: 'Quote', component: Home },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const app = createApp(App)
app.use(router)
app.mount('#app')

