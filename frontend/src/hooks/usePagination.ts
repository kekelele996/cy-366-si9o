import { reactive, ref } from 'vue'

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export function usePagination<T = any>(fetcher: (page: number, pageSize: number) => Promise<PageData<T>>, pageSize = 10) {
  const list = ref<T[]>([]) as any
  const total = ref(0)
  const loading = ref(false)
  const page = reactive({ current: 1, size: pageSize })

  async function load() {
    loading.value = true
    try {
      const data = await fetcher(page.current, page.size)
      list.value = data.list || []
      total.value = data.total || 0
    } finally {
      loading.value = false
    }
  }

  function onPageChange(p: number) {
    page.current = p
    load()
  }

  return { list, total, loading, page, load, onPageChange }
}
