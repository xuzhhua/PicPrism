export interface ImageItem {
  id: string
  filename: string
  ext: string
  size: number
  width: number
  height: number
  mime_type: string
  created_at: string
  tags: string[]
  url: string
  webp_url: string
  thumb_url: string
  markdown: string
  html: string
  bbcode: string
}

export interface ListResponse {
  items: ImageItem[]
  total: number
  page: number
  limit: number
}

export interface TagItem {
  id: number
  name: string
  count: number
}

function getToken(): string {
  return localStorage.getItem('picprism_token') ?? ''
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string>),
  }
  if (token) headers['Authorization'] = `Bearer ${token}`
  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }

  const res = await fetch(path, { ...options, headers })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? res.statusText)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export const api = {
  listImages(params: { tag?: string; page?: number; limit?: number; sort?: string }) {
    const q = new URLSearchParams()
    if (params.tag) q.set('tag', params.tag)
    if (params.page) q.set('page', String(params.page))
    if (params.limit) q.set('limit', String(params.limit))
    if (params.sort) q.set('sort', params.sort)
    return request<ListResponse>(`/api/v1/images?${q}`)
  },

  getImage(id: string) {
    return request<ImageItem>(`/api/v1/images/${id}`)
  },

  uploadImage(file: File, tags: string[]) {
    const form = new FormData()
    form.append('file', file)
    if (tags.length) form.append('tags', tags.join(','))
    return request<ImageItem>('/api/v1/images', { method: 'POST', body: form })
  },

  deleteImage(id: string) {
    return request<void>(`/api/v1/images/${id}`, { method: 'DELETE' })
  },

  updateTags(id: string, tags: string[]) {
    return request<ImageItem>(`/api/v1/images/${id}/tags`, {
      method: 'PUT',
      body: JSON.stringify({ tags }),
    })
  },

  listTags() {
    return request<TagItem[]>('/api/v1/tags')
  },
}
