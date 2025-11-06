import axios from 'axios'

// Используем тот же API, что и основной сайт
// В production это будет относительный путь, в dev можно указать полный URL
const API_URL = import.meta.env.VITE_API_URL || ''

export const apiClient = axios.create({
  baseURL: API_URL ? `${API_URL}/api` : '/api',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
})

// Интерфейсы для типизации
export interface Quote {
  id: string
  text: string
  author: string
  likes_count: number
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

  getById: async (id: string): Promise<Quote> => {
    const response = await apiClient.get<Quote>(`/quotes/${id}`)
    return response.data
  },

  create: async (data: CreateQuoteRequest): Promise<Quote> => {
    const response = await apiClient.post<Quote>('/quotes', data)
    return response.data
  },

  update: async (id: string, data: UpdateQuoteRequest): Promise<Quote> => {
    const response = await apiClient.put<Quote>(`/quotes/${id}`, data)
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/quotes/${id}`)
  },

  resetLikes: async (): Promise<void> => {
    await apiClient.delete('/quotes/likes/reset')
  },
}

