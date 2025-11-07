import axios from 'axios'

// Получаем URL API из переменной окружения
// Если VITE_API_URL не задан, используем относительный путь
// Для локальной разработки можно задать http://localhost:8080 или http://localhost:8081
const API_URL = import.meta.env.VITE_API_URL || ''

// Создаем API client
export const apiClient = axios.create({
  baseURL: API_URL ? `${API_URL}/api` : '/api',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 секунд таймаут
})

// Обработчик ошибок для логирования
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', {
      message: error.message,
      url: error.config?.url,
      status: error.response?.status,
      data: error.response?.data,
    })
    return Promise.reject(error)
  }
)

// Интерфейсы для типизации
export interface Quote {
  id: string
  text: string
  author: string
  likes_count: number
  is_liked: boolean // Информация о том, лайкнул ли текущий пользователь эту цитату
  created_at: string
  updated_at: string
}

export interface CreateQuoteRequest {
  text: string
  author: string
}

export interface UpdateQuoteRequest {
  text?: string
  author?: string
}

export interface PaginatedQuotesResponse {
  quotes: Quote[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// API методы
export const quotesApi = {
  // Получить случайную цитату
  getRandom: async (): Promise<Quote> => {
    const response = await apiClient.get<Quote>('/quotes/random')
    return response.data
  },

  // Получить все цитаты с пагинацией
  getAll: async (page: number = 1, pageSize: number = 10, search?: string): Promise<PaginatedQuotesResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    })
    if (search) {
      params.append('search', search)
    }
    const response = await apiClient.get<PaginatedQuotesResponse>(`/quotes?${params.toString()}`)
    return response.data
  },

  // Получить цитату по ID
  getById: async (id: string): Promise<Quote> => {
    const response = await apiClient.get<Quote>(`/quotes/${id}`)
    return response.data
  },

  // Создать новую цитату
  create: async (data: CreateQuoteRequest): Promise<Quote> => {
    const response = await apiClient.post<Quote>('/quotes', data)
    return response.data
  },

  // Обновить цитату
  update: async (id: string, data: UpdateQuoteRequest): Promise<Quote> => {
    const response = await apiClient.put<Quote>(`/quotes/${id}`, data)
    return response.data
  },

  // Удалить цитату
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/quotes/${id}`)
  },

  // Поставить лайк цитате
  like: async (id: string): Promise<Quote> => {
    const response = await apiClient.put<Quote>(`/quotes/${id}/like`)
    return response.data
  },

  // Получить топ цитату за неделю
  getTopWeekly: async (): Promise<Quote> => {
    const response = await apiClient.get<Quote>('/quotes/top/weekly')
    return response.data
  },

  // Получить топ цитату за всё время
  getTopAllTime: async (): Promise<Quote> => {
    const response = await apiClient.get<Quote>('/quotes/top/alltime')
    return response.data
  },

  // Сбросить все лайки
  resetLikes: async (): Promise<void> => {
    await apiClient.delete('/quotes/likes/reset')
  },
}

