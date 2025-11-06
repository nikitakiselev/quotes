<template>
  <div class="min-h-screen bg-apple-gray py-8 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- Заголовок -->
      <div class="mb-8 flex justify-between items-center">
        <div>
          <h1 class="text-4xl font-light text-apple-dark mb-2">Административная панель</h1>
        </div>
        <button
          @click="handleLogout"
          class="px-4 py-2 text-sm bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300 transition-colors"
        >
          Выйти
        </button>
      </div>

      <!-- Форма добавления/редактирования -->
      <div class="bg-white rounded-2xl shadow-sm p-6 mb-8">
        <h2 class="text-2xl font-light mb-6 text-apple-dark">
          {{ editingQuote ? 'Редактировать цитату' : 'Добавить новую цитату' }}
        </h2>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Текст цитаты</label>
            <textarea
              v-model="form.text"
              required
              rows="4"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-apple-dark focus:border-transparent"
              placeholder="Введите текст цитаты..."
            ></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Автор</label>
            <input
              v-model="form.author"
              required
              type="text"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-apple-dark focus:border-transparent"
              placeholder="Введите имя автора..."
            />
          </div>
          <div class="flex gap-4">
            <button
              type="submit"
              :disabled="submitting"
              class="px-6 py-2 bg-apple-dark text-white rounded-lg font-medium hover:bg-gray-800 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ submitting ? 'Сохранение...' : (editingQuote ? 'Сохранить' : 'Добавить') }}
            </button>
            <button
              v-if="editingQuote"
              type="button"
              @click="cancelEdit"
              class="px-6 py-2 bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300 transition-colors"
            >
              Отмена
            </button>
          </div>
        </form>
      </div>

      <!-- Поиск -->
      <div class="bg-white rounded-2xl shadow-sm p-6 mb-8">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Поиск по тексту или автору..."
          class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-apple-dark focus:border-transparent"
          @input="handleSearch"
        />
      </div>

      <!-- Таблица цитат -->
      <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-apple-gray">
              <tr>
                <th class="px-6 py-4 text-left text-sm font-medium text-gray-700">Текст</th>
                <th class="px-6 py-4 text-left text-sm font-medium text-gray-700">Автор</th>
                <th class="px-6 py-4 text-left text-sm font-medium text-gray-700">Дата создания</th>
                <th class="px-6 py-4 text-right text-sm font-medium text-gray-700">Действия</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loading" class="text-center">
                <td colspan="4" class="px-6 py-12">
                  <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-gray-300 border-t-apple-dark"></div>
                </td>
              </tr>
              <tr v-else-if="quotes.length === 0" class="text-center">
                <td colspan="4" class="px-6 py-12 text-gray-500">
                  Цитаты не найдены
                </td>
              </tr>
              <tr v-else v-for="quote in quotes" :key="quote.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 text-sm text-gray-900 max-w-md">
                  <div class="truncate" :title="quote.text">
                    {{ quote.text }}
                  </div>
                </td>
                <td class="px-6 py-4 text-sm text-gray-700">{{ quote.author }}</td>
                <td class="px-6 py-4 text-sm text-gray-500">
                  {{ new Date(quote.created_at).toLocaleDateString('ru-RU') }}
                </td>
                <td class="px-6 py-4 text-right">
                  <div class="flex justify-end gap-2">
                    <button
                      @click="editQuote(quote)"
                      class="px-4 py-2 text-sm text-blue-600 hover:text-blue-800 transition-colors"
                    >
                      Редактировать
                    </button>
                    <button
                      @click="deleteQuote(quote.id)"
                      class="px-4 py-2 text-sm text-red-600 hover:text-red-800 transition-colors"
                    >
                      Удалить
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Пагинация -->
        <div v-if="totalPages > 1" class="px-6 py-4 bg-apple-gray border-t border-gray-200 flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Страница {{ currentPage }} из {{ totalPages }} (всего {{ total }})
          </div>
          <div class="flex gap-2">
            <button
              @click="changePage(currentPage - 1)"
              :disabled="currentPage === 1"
              class="px-4 py-2 text-sm bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Назад
            </button>
            <button
              @click="changePage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="px-4 py-2 text-sm bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Вперед
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { quotesApi, type Quote, type CreateQuoteRequest, type UpdateQuoteRequest } from '@/api/client'

const router = useRouter()
const quotes = ref<Quote[]>([])
const loading = ref(false)
const submitting = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const totalPages = ref(0)
const searchQuery = ref('')
const editingQuote = ref<Quote | null>(null)

const form = ref<CreateQuoteRequest>({
  text: '',
  author: '',
})

let searchTimeout: ReturnType<typeof setTimeout>

const handleLogout = () => {
  localStorage.removeItem('admin_authenticated')
  router.push('/login')
}

const loadQuotes = async () => {
  loading.value = true
  try {
    const response = await quotesApi.getAll(currentPage.value, pageSize.value, searchQuery.value || undefined)
    quotes.value = response.quotes
    total.value = response.total
    totalPages.value = response.total_pages
  } catch (err) {
    console.error('Error loading quotes:', err)
    await Swal.fire({
      icon: 'error',
      title: 'Ошибка',
      text: 'Не удалось загрузить цитаты',
      confirmButtonText: 'OK',
    })
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (editingQuote.value) {
      await quotesApi.update(editingQuote.value.id, form.value as UpdateQuoteRequest)
      await Swal.fire({
        icon: 'success',
        title: 'Успешно!',
        text: 'Цитата обновлена',
        timer: 2000,
        showConfirmButton: false,
      })
    } else {
      await quotesApi.create(form.value)
      await Swal.fire({
        icon: 'success',
        title: 'Успешно!',
        text: 'Цитата добавлена',
        timer: 2000,
        showConfirmButton: false,
      })
    }
    form.value = { text: '', author: '' }
    editingQuote.value = null
    await loadQuotes()
  } catch (err) {
    console.error('Error saving quote:', err)
    await Swal.fire({
      icon: 'error',
      title: 'Ошибка',
      text: 'Не удалось сохранить цитату',
      confirmButtonText: 'OK',
    })
  } finally {
    submitting.value = false
  }
}

const editQuote = (quote: Quote) => {
  editingQuote.value = quote
  form.value = {
    text: quote.text,
    author: quote.author,
  }
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const cancelEdit = () => {
  editingQuote.value = null
  form.value = { text: '', author: '' }
}

const deleteQuote = async (id: string) => {
  const result = await Swal.fire({
    icon: 'warning',
    title: 'Удаление цитаты',
    text: 'Вы уверены, что хотите удалить эту цитату?',
    showCancelButton: true,
    confirmButtonColor: '#d33',
    cancelButtonColor: '#3085d6',
    confirmButtonText: 'Да, удалить',
    cancelButtonText: 'Отмена',
  })

  if (!result.isConfirmed) {
    return
  }

  try {
    await quotesApi.delete(id)
    await Swal.fire({
      icon: 'success',
      title: 'Удалено!',
      text: 'Цитата успешно удалена',
      timer: 2000,
      showConfirmButton: false,
    })
    await loadQuotes()
  } catch (err) {
    console.error('Error deleting quote:', err)
    await Swal.fire({
      icon: 'error',
      title: 'Ошибка',
      text: 'Не удалось удалить цитату',
      confirmButtonText: 'OK',
    })
  }
}

const changePage = (page: number) => {
  currentPage.value = page
  loadQuotes()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    loadQuotes()
  }, 500)
}

onMounted(() => {
  loadQuotes()
})
</script>

