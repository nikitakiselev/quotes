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
      <div class="mb-8 min-h-[300px] flex items-center justify-center">
        <Transition name="quote" mode="out-in">
          <div 
            v-if="quote" 
            :key="quote.id"
            class="w-full"
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
        </Transition>
      </div>

      <!-- Загрузка - поверх всей страницы -->
      <Transition name="loader">
        <div
          v-if="showLoader"
          class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-80 backdrop-blur-sm"
        >
          <div class="animate-spin rounded-full h-16 w-16 border-4 border-gray-300 border-t-apple-dark"></div>
        </div>
      </Transition>

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
          @touchstart="handleTouchStart"
          @touchend="handleTouchEnd"
          :disabled="loading"
          :class="[
            'px-5 sm:px-8 py-2.5 sm:py-3 rounded-full font-medium text-sm sm:text-base text-white transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-apple-dark focus:ring-offset-2 shadow-sm flex items-center gap-2 select-none touch-manipulation',
            isSpacePressed || isTouching ? 'bg-gray-700 scale-95' : 'bg-apple-dark hover:bg-gray-800 active:bg-gray-700 active:scale-95'
          ]"
        >
          <span>Следующая цитата</span>
          <!-- Kbd компонент - только для десктопа -->
          <kbd class="hidden md:inline-flex items-center justify-center bg-white border border-gray-300 rounded shadow-sm">
            <svg
              width="45"
              height="20"
              viewBox="0 0 122.88 54.99"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              class="text-gray-600"
            >
              <path
                fill-rule="evenodd"
                clip-rule="evenodd"
                d="M11.61,0h99.65c3.19,0,6.09,1.31,8.2,3.41c2.11,2.11,3.42,5.02,3.42,8.2v31.77c0,3.18-1.31,6.09-3.42,8.19 c-2.11,2.11-5.01,3.42-8.19,3.42H11.61c-3.18,0-6.09-1.31-8.2-3.42c-2.1-2.1-3.41-5-3.41-8.2V11.61c0-3.2,1.31-6.1,3.41-8.2 C5.51,1.31,8.42,0,11.61,0L11.61,0z M29.34,26.85l4.06-0.27c0.09,0.69,0.27,1.23,0.54,1.59c0.44,0.59,1.07,0.89,1.89,0.89 c0.61,0,1.08-0.15,1.41-0.45c0.33-0.3,0.5-0.65,0.5-1.05c0-0.38-0.15-0.72-0.47-1.02c-0.31-0.3-1.04-0.58-2.18-0.85 c-1.87-0.44-3.21-1.03-4-1.77c-0.8-0.74-1.21-1.67-1.21-2.82c0-0.75,0.21-1.46,0.62-2.12c0.41-0.67,1.03-1.19,1.86-1.58 c0.83-0.38,1.96-0.57,3.4-0.57c1.77,0,3.11,0.35,4.04,1.05c0.93,0.69,1.48,1.8,1.65,3.32l-4.02,0.25 c-0.11-0.66-0.33-1.15-0.67-1.45c-0.34-0.3-0.82-0.45-1.42-0.45c-0.5,0-0.87,0.11-1.12,0.33c-0.25,0.22-0.38,0.49-0.38,0.81 c0,0.23,0.1,0.44,0.31,0.63c0.2,0.19,0.67,0.37,1.43,0.54c1.87,0.43,3.2,0.86,4.01,1.29c0.81,0.44,1.4,0.97,1.76,1.62 c0.37,0.64,0.55,1.36,0.55,2.16c0,0.94-0.24,1.8-0.73,2.59c-0.49,0.79-1.17,1.39-2.06,1.8c-0.88,0.41-1.99,0.61-3.32,0.61 c-2.35,0-3.98-0.48-4.88-1.44C29.98,29.54,29.46,28.32,29.34,26.85L29.34,26.85z M43.79,35.71V21.08h3.6v1.56 c0.5-0.66,0.95-1.11,1.37-1.34c0.56-0.31,1.18-0.47,1.86-0.47c1.34,0,2.37,0.54,3.11,1.63c0.73,1.09,1.1,2.43,1.1,4.03 c0,1.77-0.4,3.11-1.2,4.04c-0.8,0.93-1.81,1.39-3.03,1.39c-0.59,0-1.13-0.11-1.62-0.32c-0.49-0.21-0.92-0.53-1.31-0.95v5.04H43.79 L43.79,35.71z M47.64,26.42c0,0.84,0.17,1.47,0.5,1.88c0.33,0.41,0.75,0.61,1.26,0.61c0.44,0,0.82-0.19,1.12-0.59 c0.3-0.39,0.45-1.06,0.45-1.99c0-0.86-0.16-1.5-0.47-1.9c-0.32-0.41-0.7-0.61-1.15-0.61c-0.49,0-0.89,0.2-1.22,0.61 C47.81,24.84,47.64,25.5,47.64,26.42L47.64,26.42z M60.17,24.52l-3.67-0.41c0.14-0.68,0.34-1.22,0.6-1.6 c0.26-0.39,0.63-0.73,1.13-1.01c0.35-0.21,0.83-0.37,1.45-0.48c0.61-0.11,1.28-0.17,2-0.17c1.15,0,2.07,0.07,2.77,0.2 c0.7,0.14,1.28,0.42,1.74,0.85c0.33,0.3,0.59,0.73,0.77,1.27c0.19,0.55,0.28,1.07,0.28,1.57v4.68c0,0.5,0.03,0.89,0.09,1.17 c0.06,0.28,0.19,0.64,0.39,1.08h-3.59c-0.14-0.27-0.24-0.48-0.28-0.62c-0.04-0.14-0.09-0.36-0.13-0.67c-0.5,0.51-1,0.87-1.5,1.09 c-0.68,0.29-1.47,0.44-2.37,0.44c-1.19,0-2.1-0.29-2.72-0.88c-0.62-0.59-0.93-1.31-0.93-2.17c0-0.81,0.22-1.47,0.67-1.99 c0.45-0.52,1.27-0.91,2.47-1.16c1.44-0.31,2.37-0.52,2.8-0.65c0.43-0.12,0.88-0.28,1.36-0.48c0-0.5-0.1-0.85-0.29-1.05 c-0.19-0.2-0.53-0.3-1.02-0.3c-0.63,0-1.1,0.11-1.41,0.32C60.51,23.74,60.32,24.05,60.17,24.52L60.17,24.52z M63.49,26.64 c-0.53,0.2-1.08,0.38-1.65,0.53c-0.78,0.22-1.27,0.44-1.48,0.65c-0.22,0.22-0.32,0.47-0.32,0.76c0,0.32,0.11,0.59,0.32,0.79 c0.21,0.2,0.52,0.31,0.93,0.31c0.43,0,0.83-0.11,1.2-0.33c0.37-0.22,0.63-0.49,0.78-0.81c0.15-0.32,0.23-0.74,0.23-1.25V26.64 L63.49,26.64z M77.01,27.54l3.65,0.43c-0.2,0.81-0.53,1.5-0.99,2.1c-0.46,0.59-1.05,1.05-1.76,1.38c-0.71,0.33-1.62,0.49-2.72,0.49 c-1.06,0-1.95-0.11-2.66-0.31c-0.7-0.21-1.31-0.55-1.82-1.02c-0.51-0.47-0.91-1.02-1.2-1.65c-0.29-0.63-0.43-1.47-0.43-2.52 c0-1.09,0.18-2,0.53-2.72c0.26-0.53,0.61-1.01,1.06-1.43c0.45-0.42,0.9-0.74,1.38-0.94c0.75-0.33,1.71-0.49,2.88-0.49 c1.64,0,2.88,0.31,3.74,0.93c0.86,0.62,1.46,1.52,1.81,2.71l-3.61,0.51c-0.11-0.45-0.32-0.79-0.62-1.02 c-0.3-0.23-0.7-0.34-1.21-0.34c-0.63,0-1.15,0.24-1.54,0.72c-0.39,0.48-0.59,1.21-0.59,2.19c0,0.87,0.2,1.53,0.59,1.98 c0.39,0.45,0.88,0.68,1.49,0.68c0.5,0,0.93-0.14,1.27-0.41C76.58,28.52,76.84,28.1,77.01,27.54L77.01,27.54z M93.54,27.4h-7.7 c0.07,0.65,0.24,1.14,0.5,1.46c0.37,0.46,0.86,0.69,1.45,0.69c0.38,0,0.74-0.1,1.08-0.3c0.21-0.13,0.43-0.35,0.67-0.66l3.78,0.37 c-0.58,1.06-1.28,1.83-2.1,2.29c-0.82,0.46-1.99,0.69-3.52,0.69c-1.33,0-2.37-0.2-3.13-0.6c-0.76-0.39-1.39-1.02-1.89-1.89 c-0.5-0.86-0.75-1.88-0.75-3.04c0-1.66,0.5-3,1.51-4.03c1-1.02,2.39-1.54,4.16-1.54c1.44,0,2.56,0.23,3.4,0.69 c0.83,0.46,1.46,1.12,1.9,2c0.43,0.87,0.65,2.01,0.65,3.41V27.4L93.54,27.4z M89.63,25.46c-0.07-0.79-0.28-1.35-0.6-1.7 c-0.33-0.34-0.75-0.51-1.28-0.51c-0.61,0-1.1,0.26-1.47,0.77c-0.23,0.32-0.38,0.8-0.44,1.43H89.63L89.63,25.46z M2.5,39.07 c0.43,1.63,1.29,3.09,2.45,4.25c1.71,1.71,4.07,2.77,6.66,2.77h99.66c0.1,0,0.21,0,0.31-0.01l0.64-0.04 c2.21-0.22,4.21-1.22,5.71-2.72c1.16-1.16,2.02-2.62,2.45-4.24V11.61c0-2.49-1.03-4.78-2.68-6.43c-1.65-1.65-3.93-2.68-6.43-2.68 H11.61c-2.5,0-4.78,1.03-6.43,2.68C3.53,6.83,2.5,9.11,2.5,11.61V39.07L2.5,39.07z"
                fill="currentColor"
              />
            </svg>
          </kbd>
        </button>
      </div>
    </div>

    <!-- Кнопка помощи (хоткеи) -->
    <button
      @click="showHotkeysModal = true"
      class="fixed bottom-6 right-6 w-10 h-10 bg-white border border-gray-300 rounded-full shadow-lg hover:bg-gray-50 transition-colors duration-200 flex items-center justify-center text-gray-600 hover:text-gray-800 focus:outline-none focus:ring-2 focus:ring-apple-dark focus:ring-offset-2 z-40"
      title="Горячие клавиши"
    >
      <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="12" cy="12" r="10"></circle>
        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
        <line x1="12" y1="17" x2="12.01" y2="17"></line>
      </svg>
    </button>

    <!-- Модальное окно с хоткеями -->
    <Transition name="modal">
      <div
        v-if="showHotkeysModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click="showHotkeysModal = false"
      >
        <!-- Размытый фон -->
        <div class="absolute inset-0 modal-backdrop"></div>
        
        <!-- Модальное окно -->
        <div
          class="relative bg-white rounded-2xl shadow-2xl max-w-md w-full p-8"
          @click.stop
        >
          <!-- Заголовок -->
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-2xl font-light text-apple-dark">Горячие клавиши</h2>
            <button
              @click="showHotkeysModal = false"
              class="text-gray-400 hover:text-gray-600 transition-colors focus:outline-none focus:ring-2 focus:ring-apple-dark rounded-full p-1"
            >
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <!-- Список хоткеев -->
          <div class="space-y-4">
            <div class="flex items-center justify-between py-3 border-b border-gray-100">
              <div class="flex items-center gap-3">
                <kbd class="px-3 py-1.5 bg-gray-100 border border-gray-300 rounded text-sm font-medium text-gray-700">
                  Space
                </kbd>
                <span class="text-gray-700">Следующая цитата</span>
              </div>
            </div>

            <div class="flex items-center justify-between py-3 border-b border-gray-100">
              <div class="flex items-center gap-3">
                <kbd class="px-3 py-1.5 bg-gray-100 border border-gray-300 rounded text-sm font-medium text-gray-700">
                  F
                </kbd>
                <span class="text-gray-700">Поставить лайк</span>
              </div>
            </div>
          </div>

          <!-- Подсказка -->
          <p class="mt-6 text-sm text-gray-500 text-center">
            Нажмите вне окна или ESC для закрытия
          </p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { quotesApi, type Quote } from '@/api/client'

const route = useRoute()
const router = useRouter()

const quote = ref<Quote | null>(null)
const loading = ref(false)
const showLoader = ref(false) // Показывать loader только если загрузка длится > 0.8s
const error = ref<string | null>(null)
const likeAnimating = ref(false)
const isUpdatingUrl = ref(false) // Флаг для предотвращения повторной загрузки при программном обновлении URL
const isSpacePressed = ref(false) // Флаг для отслеживания нажатия пробела
const isTouching = ref(false) // Флаг для отслеживания тач-событий
const showHotkeysModal = ref(false) // Флаг для отображения модального окна с хоткеями
let loaderTimer: ReturnType<typeof setTimeout> | null = null

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

const loadQuote = async (quoteLoader: () => Promise<Quote>, updateUrl = true) => {
  loading.value = true
  error.value = null
  showLoader.value = false
  
  // Очищаем предыдущий таймер, если есть
  if (loaderTimer) {
    clearTimeout(loaderTimer)
    loaderTimer = null
  }
  
  // Показываем loader только если загрузка длится больше 0.8 секунд
  loaderTimer = setTimeout(() => {
    if (loading.value) {
      showLoader.value = true
    }
  }, 800)
  
  try {
    // Сервер автоматически возвращает is_liked в ответе
    const loadedQuote = await quoteLoader()
    quote.value = loadedQuote
    
    // Обновляем URL в адресной строке, чтобы он указывал на ID цитаты
    // Но только если ID в URL отличается от ID загруженной цитаты
    if (updateUrl && loadedQuote && route.params.id !== loadedQuote.id) {
      isUpdatingUrl.value = true
      router.replace({ name: 'Quote', params: { id: loadedQuote.id } }).finally(() => {
        // Небольшая задержка, чтобы watch успел обработать изменение
        setTimeout(() => {
          isUpdatingUrl.value = false
        }, 100)
      })
    }
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
    showLoader.value = false
    
    // Очищаем таймер
    if (loaderTimer) {
      clearTimeout(loaderTimer)
      loaderTimer = null
    }
  }
}

// Загрузка цитаты по ID из URL
const loadQuoteById = async (id: string) => {
  // Проверяем, что это действительно другая цитата
  if (quote.value?.id === id) {
    return // Цитата уже загружена
  }
  await loadQuote(() => quotesApi.getById(id), false) // Не обновляем URL, так как он уже правильный
}

const loadRandomQuote = () => {
  // Очищаем ID из URL при загрузке случайной цитаты
  if (route.params.id) {
    isUpdatingUrl.value = true
    router.replace('/').finally(() => {
      setTimeout(() => {
        isUpdatingUrl.value = false
      }, 100)
    })
  }
  loadQuote(quotesApi.getRandom)
}
const loadTopWeekly = async () => {
  // Очищаем ID из URL перед загрузкой
  if (route.params.id) {
    isUpdatingUrl.value = true
    await router.replace('/')
    setTimeout(() => {
      isUpdatingUrl.value = false
    }, 100)
  }
  // Обновляем URL для топ-цитат
  await loadQuote(quotesApi.getTopWeekly, true)
}
const loadTopAllTime = async () => {
  // Очищаем ID из URL перед загрузкой
  if (route.params.id) {
    isUpdatingUrl.value = true
    await router.replace('/')
    setTimeout(() => {
      isUpdatingUrl.value = false
    }, 100)
  }
  // Обновляем URL для топ-цитат
  await loadQuote(quotesApi.getTopAllTime, true)
}

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

// Обработчики тач-событий для кнопки
const handleTouchStart = () => {
  if (!loading.value) {
    isTouching.value = true
  }
}

const handleTouchEnd = () => {
  isTouching.value = false
}

// Обработчик нажатия клавиши пробела (для визуальной обратной связи)
const handleKeyDown = (event: KeyboardEvent) => {
  // Игнорируем, если пользователь вводит текст в поле ввода
  if (
    event.target instanceof HTMLInputElement ||
    event.target instanceof HTMLTextAreaElement ||
    event.target instanceof HTMLSelectElement
  ) {
    return
  }

  // Обрабатываем пробел
  if (event.code === 'Space' || event.key === ' ') {
    // Предотвращаем прокрутку страницы при нажатии пробела
    event.preventDefault()
    
    // Устанавливаем флаг нажатия для визуальной обратной связи
    if (!isSpacePressed.value && !loading.value) {
      isSpacePressed.value = true
    }
  }

  // Обрабатываем клавишу F для лайка
  if (event.code === 'KeyF' || event.key === 'f' || event.key === 'F') {
    // Предотвращаем стандартное поведение
    event.preventDefault()
    
    // Ставим лайк, если цитата есть и не идет загрузка
    if (quote.value && !loading.value && !isLiked.value) {
      handleLike()
    }
  }

  // Обрабатываем клавишу Escape для закрытия модального окна
  if (event.code === 'Escape' || event.key === 'Escape') {
    if (showHotkeysModal.value) {
      showHotkeysModal.value = false
    }
  }
}

// Обработчик отпускания клавиши пробела (для загрузки цитаты)
const handleKeyUp = (event: KeyboardEvent) => {
  // Игнорируем, если пользователь вводит текст в поле ввода
  if (
    event.target instanceof HTMLInputElement ||
    event.target instanceof HTMLTextAreaElement ||
    event.target instanceof HTMLSelectElement
  ) {
    return
  }

  // Обрабатываем только пробел
  if (event.code === 'Space' || event.key === ' ') {
    // Предотвращаем прокрутку страницы при отпускании пробела
    event.preventDefault()
    
    // Сбрасываем флаг нажатия
    isSpacePressed.value = false
    
    // Загружаем случайную цитату только если не идет загрузка
    if (!loading.value) {
      loadRandomQuote()
    }
  }
}

// Загружаем цитату при монтировании или изменении ID в URL
onMounted(() => {
  const quoteId = route.params.id as string | undefined
  if (quoteId) {
    loadQuoteById(quoteId)
  } else {
    loadRandomQuote()
  }

  // Добавляем обработчики нажатия и отпускания клавиш
  window.addEventListener('keydown', handleKeyDown)
  window.addEventListener('keyup', handleKeyUp)
})

// Удаляем обработчики при размонтировании
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('keyup', handleKeyUp)
})

// Отслеживаем изменения ID в URL
watch(() => route.params.id, (newId, oldId) => {
  // Пропускаем загрузку, если URL обновляется программно
  if (isUpdatingUrl.value) {
    return
  }
  
  // Загружаем цитату только если ID действительно изменился
  if (newId && typeof newId === 'string' && newId !== oldId) {
    loadQuoteById(newId)
  } else if (!newId && !quote.value) {
    loadRandomQuote()
  }
})

</script>

<style scoped>
/* Анимация появления loader */
.loader-enter-active {
  transition: opacity 0.2s ease-out;
}

.loader-leave-active {
  transition: opacity 0.2s ease-in;
}

.loader-enter-from,
.loader-leave-to {
  opacity: 0;
}

/* Анимация появления цитаты */
.quote-enter-active {
  transition: opacity 0.5s ease-out, transform 0.5s ease-out;
}

.quote-leave-active {
  transition: opacity 0.3s ease-in, transform 0.3s ease-in;
}

.quote-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.quote-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Анимация модального окна */
.modal-enter-active {
  transition: opacity 0.3s ease-out;
}

.modal-leave-active {
  transition: opacity 0.2s ease-in;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .bg-white,
.modal-leave-active .bg-white {
  transition: transform 0.3s ease-out, opacity 0.3s ease-out;
}

.modal-enter-from .bg-white,
.modal-leave-to .bg-white {
  opacity: 0;
  transform: scale(0.95) translateY(-10px);
}

/* Плавное размытие и затемнение фона */
.modal-backdrop {
  background-color: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.modal-enter-active .modal-backdrop {
  transition: background-color 0.25s ease-out, backdrop-filter 0.25s ease-out, -webkit-backdrop-filter 0.25s ease-out;
  transition-delay: 0s;
}

.modal-leave-active .modal-backdrop {
  transition: background-color 0.2s ease-in, backdrop-filter 0.2s ease-in, -webkit-backdrop-filter 0.2s ease-in;
}

.modal-enter-from .modal-backdrop {
  background-color: rgba(0, 0, 0, 0);
  backdrop-filter: blur(0px);
  -webkit-backdrop-filter: blur(0px);
}

.modal-leave-to .modal-backdrop {
  background-color: rgba(0, 0, 0, 0);
  backdrop-filter: blur(0px);
  -webkit-backdrop-filter: blur(0px);
}

.modal-enter-active .bg-white {
  transition-delay: 0.05s;
}

/* Улучшение отзывчивости на тач-устройствах */
.touch-manipulation {
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
  -webkit-touch-callout: none;
}
</style>

