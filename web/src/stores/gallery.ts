import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api, type ImageItem, type TagItem } from '../api/client'

export const useGalleryStore = defineStore('gallery', () => {
  const images = ref<ImageItem[]>([])
  const total = ref(0)
  const page = ref(1)
  const limit = ref(40)
  const activeTag = ref('')
  const sort = ref('newest')
  const loading = ref(false)
  const tags = ref<TagItem[]>([])
  const selected = ref<ImageItem | null>(null)
  const lightboxOpen = ref(false)
  const multiSelected = ref<Set<string>>(new Set())

  const totalPages = computed(() => Math.ceil(total.value / limit.value))

  async function loadImages() {
    loading.value = true
    try {
      const res = await api.listImages({
        tag: activeTag.value || undefined,
        page: page.value,
        limit: limit.value,
        sort: sort.value,
      })
      images.value = res.items ?? []
      total.value = res.total
    } finally {
      loading.value = false
    }
  }

  async function loadTags() {
    const res = await api.listTags()
    tags.value = res ?? []
  }

  async function deleteImage(id: string) {
    await api.deleteImage(id)
    images.value = images.value.filter((i) => i.id !== id)
    total.value--
    if (selected.value?.id === id) selected.value = null
    multiSelected.value.delete(id)
    await loadTags()
  }

  async function deleteMultiple(ids: string[]) {
    await Promise.all(ids.map((id) => api.deleteImage(id)))
    const idSet = new Set(ids)
    images.value = images.value.filter((i) => !idSet.has(i.id))
    total.value -= ids.length
    if (selected.value && idSet.has(selected.value.id)) selected.value = null
    ids.forEach((id) => multiSelected.value.delete(id))
    await loadTags()
  }

  async function updateTags(id: string, newTags: string[]) {
    const updated = await api.updateTags(id, newTags)
    const idx = images.value.findIndex((i) => i.id === id)
    if (idx !== -1) images.value[idx] = updated
    if (selected.value?.id === id) selected.value = updated
    await loadTags()
  }

  function setTag(tag: string) {
    activeTag.value = tag
    page.value = 1
    loadImages()
  }

  function setPage(p: number) {
    page.value = p
    loadImages()
  }

  return {
    images, total, page, limit, activeTag, sort, loading,
    tags, selected, lightboxOpen, multiSelected, totalPages,
    loadImages, loadTags, deleteImage, deleteMultiple, updateTags, setTag, setPage,
  }
})
