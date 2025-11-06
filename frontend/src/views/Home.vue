<template>
  <div class="h-full min-h-0 flex flex-col items-center justify-center px-4 py-4 sm:py-8 bg-gradient-to-br from-gray-50 to-white relative">
    <div class="max-w-4xl w-full text-center">
      <!-- Кнопки Top - прибиты к верху -->
      <div class="absolute top-4 left-0 right-0 flex flex-row gap-2 justify-center items-center">
        <button
          @click="loadTopWeekly"
          :disabled="loading"
          class="px-4 py-1.5 text-xs bg-white border border-gray-300 text-apple-dark rounded-full font-medium hover:bg-gray-50 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-apple-dark focus:ring-offset-2 shadow-sm"
        >
          Top за неделю
        </button>
        <button
          @click="loadTopAllTime"
          :disabled="loading"
          class="px-4 py-1.5 text-xs bg-white border border-gray-300 text-apple-dark rounded-full font-medium hover:bg-gray-50 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-apple-dark focus:ring-offset-2 shadow-sm"
        >
          Top за всё время
        </button>
      </div>

      <!-- Цитата -->
      <div 
        v-if="quote" 
        class="mb-8 transition-all duration-500 ease-in-out"
        :class="{ 'opacity-0': loading }"
      >
        <blockquote 
          @dblclick="handleDoubleClick"
          :class="['font-light leading-tight text-apple-dark mb-8 text-balance cursor-pointer select-none', quoteClass]"
        >
          "{{ quote.text }}"
        </blockquote>
        <p class="text-xl md:text-2xl text-gray-600 font-light">
          — {{ quote.author }}
        </p>
      </div>

      <!-- Загрузка -->
      <div v-if="loading" class="mb-8">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-gray-300 border-t-apple-dark"></div>
      </div>

      <!-- Ошибка -->
      <div v-if="error" class="mb-8 text-red-500">
        <p class="text-lg">{{ error }}</p>
      </div>

      <!-- Кнопки действий - в одну строку на всех устройствах -->
      <div class="flex flex-row gap-3 sm:gap-4 justify-center items-center">
        <!-- Кнопка лайка -->
        <button
          @click="handleLike"
          :disabled="loading || !quote"
          :class="[
            'px-4 sm:px-6 py-2.5 sm:py-3 rounded-full font-medium text-sm sm:text-base transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-offset-2',
            isLiked 
              ? 'bg-red-50 text-red-500 border-2 border-red-200 hover:bg-red-100 hover:border-red-300 shadow-sm' 
              : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 shadow-sm'
          ]"
        >
          <span class="flex items-center gap-1.5 sm:gap-2">
            <span 
              :class="[
                'text-lg sm:text-xl transition-transform duration-300',
                isLiked ? 'scale-110' : 'scale-100',
                likeAnimating ? 'animate-pulse' : ''
              ]"
            >
              ❤️
            </span>
            <span>{{ quote?.likes_count || 0 }}</span>
          </span>
        </button>

        <!-- Кнопка обновления -->
        <button
          @click="loadRandomQuote"
          :disabled="loading"
          class="px-5 sm:px-8 py-2.5 sm:py-3 bg-apple-dark text-white rounded-full font-medium text-sm sm:text-base hover:bg-gray-800 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-apple-dark focus:ring-offset-2 shadow-sm"
        >
          {{ loading ? 'Загрузка...' : 'Следующая цитата' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { quotesApi, type Quote } from '@/api/client'

const quote = ref<Quote | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const likeAnimating = ref(false)

// Проверка, лайкнута ли текущая цитата (из ответа сервера)
const isLiked = computed(() => {
  return quote.value?.is_liked || false
})

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

const loadQuote = async (quoteLoader: () => Promise<Quote>) => {
  loading.value = true
  error.value = null
  
  try {
    // Сервер автоматически возвращает is_liked в ответе
    quote.value = await quoteLoader()
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

const loadRandomQuote = () => loadQuote(quotesApi.getRandom)
const loadTopWeekly = () => loadQuote(quotesApi.getTopWeekly)
const loadTopAllTime = () => loadQuote(quotesApi.getTopAllTime)

const handleLike = async () => {
  if (!quote.value || loading.value || isLiked.value) return

  // Ставим лайк
  try {
    likeAnimating.value = true
    // Сервер возвращает обновленную цитату с is_liked = true
    const updatedQuote = await quotesApi.like(quote.value.id)
    quote.value = updatedQuote

    // Анимация
    setTimeout(() => {
      likeAnimating.value = false
    }, 600)
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string }; status?: number }; message?: string }
    const errorMessage = error?.response?.data?.error || error?.message || 'Не удалось поставить лайк'
    
    // Если ошибка "уже лайкнуто", обновляем цитату для получения актуального состояния
    if (errorMessage.includes('уже') || errorMessage.includes('already')) {
      try {
        const updatedQuote = await quotesApi.getById(quote.value.id)
        quote.value = updatedQuote
      } catch {
        // Игнорируем ошибку обновления
      }
    } else {
      error.value = errorMessage
    }
    likeAnimating.value = false
  }
}

const handleDoubleClick = () => {
  if (!quote.value || loading.value) return
  handleLike()
}

onMounted(() => {
  loadRandomQuote()
})
</script>

