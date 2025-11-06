<template>
  <div class="min-h-screen flex items-center justify-center px-4 py-8 bg-gradient-to-br from-gray-50 to-white">
    <div class="max-w-md w-full">
      <div class="bg-white rounded-2xl shadow-lg p-8">
        <h1 class="text-3xl font-light text-apple-dark mb-2 text-center">Административная панель</h1>
        <p class="text-gray-600 text-center mb-8">Введите пароль для доступа</p>
        
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Пароль</label>
            <input
              v-model="password"
              type="password"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-apple-dark focus:border-transparent"
              placeholder="Введите пароль"
              autofocus
            />
          </div>
          
          <div v-if="error" class="text-red-500 text-sm text-center">
            {{ error }}
          </div>
          
          <button
            type="submit"
            :disabled="loading"
            class="w-full px-6 py-3 bg-apple-dark text-white rounded-lg font-medium hover:bg-gray-800 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? 'Вход...' : 'Войти' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const password = ref('')
const error = ref('')
const loading = ref(false)

const ADMIN_PASSWORD = import.meta.env.VITE_ADMIN_PASSWORD || 'admin'

const handleLogin = () => {
  error.value = ''
  loading.value = true
  
  // Простая проверка пароля
  if (password.value === ADMIN_PASSWORD) {
    localStorage.setItem('admin_authenticated', 'true')
    router.push('/admin')
  } else {
    error.value = 'Неверный пароль'
    password.value = ''
  }
  
  loading.value = false
}
</script>

