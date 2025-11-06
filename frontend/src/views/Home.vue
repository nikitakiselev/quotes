<template>
  <div class="h-full min-h-0 flex items-center justify-center px-4 py-4 sm:py-8 bg-gradient-to-br from-gray-50 to-white">
    <div class="max-w-4xl w-full text-center">
      <!-- Цитата -->
      <div 
        v-if="quote" 
        class="mb-12 transition-all duration-500 ease-in-out"
        :class="{ 'opacity-0': loading }"
      >
        <blockquote 
          :class="['font-light leading-tight text-apple-dark mb-8 text-balance', quoteClass]"
        >
          "{{ quote.text }}"
        </blockquote>
        <p class="text-xl md:text-2xl text-gray-600 font-light">
          — {{ quote.author }}
        </p>
      </div>

      <!-- Загрузка -->
      <div v-if="loading" class="mb-12">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-gray-300 border-t-apple-dark"></div>
      </div>

      <!-- Ошибка -->
      <div v-if="error" class="mb-12 text-red-500">
        <p class="text-lg">{{ error }}</p>
      </div>

      <!-- Кнопка обновления -->
      <button
        @click="loadRandomQuote"
        :disabled="loading"
        class="px-8 py-4 bg-apple-dark text-white rounded-full font-medium text-lg hover:bg-gray-800 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-apple-dark focus:ring-offset-2"
      >
        {{ loading ? 'Загрузка...' : 'Обновить цитату' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { quotesApi, type Quote } from '@/api/client'

const quote = ref<Quote | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

// Автоматическое уменьшение шрифта для длинных цитат
const quoteClass = computed(() => {
  if (!quote.value) return 'text-4xl md:text-5xl lg:text-6xl'
  
  const textLength = quote.value.text.length
  
  if (textLength > 200) {
    return 'text-2xl md:text-3xl lg:text-4xl'
  } else if (textLength > 100) {
    return 'text-3xl md:text-4xl lg:text-5xl'
  }
  
  return 'text-4xl md:text-5xl lg:text-6xl'
})

const loadRandomQuote = async () => {
  loading.value = true
  error.value = null
  
  try {
    quote.value = await quotesApi.getRandom()
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string }; status?: number }; message?: string }
    const errorMessage = error?.response?.data?.error || error?.message || 'Неизвестная ошибка'
    error.value = `Не удалось загрузить цитату: ${errorMessage}. Попробуйте еще раз.`
    console.error('Error loading quote:', {
      error: err,
      message: error?.message,
      response: error?.response?.data,
      status: error?.response?.status,
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRandomQuote()
})
</script>

