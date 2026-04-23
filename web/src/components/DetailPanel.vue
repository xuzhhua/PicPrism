<template>
  <aside v-if="img" class="detail">
    <div class="detail-header">
      <span class="detail-title">详情</span>
      <button class="icon-btn" @click="gallery.selected = null">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>

    <div class="detail-preview" @click="gallery.lightboxOpen = true">
      <img :src="img.webp_url || img.url" :alt="img.filename" class="preview-img preview-clickable" />
    </div>

    <div class="detail-meta">
      <div class="meta-row">
        <span class="meta-label">文件名</span>
        <span class="meta-val">{{ img.filename }}</span>
      </div>
      <div class="meta-row" v-if="img.width">
        <span class="meta-label">尺寸</span>
        <span class="meta-val">{{ img.width }} × {{ img.height }}</span>
      </div>
      <div class="meta-row">
        <span class="meta-label">大小</span>
        <span class="meta-val">{{ formatSize(img.size) }}</span>
      </div>
      <div class="meta-row">
        <span class="meta-label">格式</span>
        <span class="meta-val">{{ img.mime_type }}</span>
      </div>
      <div class="meta-row">
        <span class="meta-label">上传时间</span>
        <span class="meta-val">{{ formatDate(img.created_at) }}</span>
      </div>
    </div>

    <!-- 标签编辑 -->
    <div class="detail-section">
      <div class="section-label">标签</div>
      <div class="tags-wrap">
        <span
          v-for="tag in editTags"
          :key="tag"
          class="tag-chip"
        >
          {{ tag }}
          <button class="chip-del" @click="removeTag(tag)">×</button>
        </span>
        <input
          v-if="addingTag"
          ref="tagInputRef"
          v-model="newTag"
          class="tag-input"
          placeholder="标签名"
          @keydown.enter="confirmTag"
          @keydown.esc="addingTag = false"
          @blur="confirmTag"
        />
        <button v-else class="add-tag-btn" @click="startAddTag">+ 添加</button>
      </div>
      <button
        v-if="tagsChanged"
        class="save-btn"
        :disabled="saving"
        @click="saveTags"
      >{{ saving ? '保存中...' : '保存标签' }}</button>
    </div>

    <!-- 链接 -->
    <div class="detail-section">
      <div class="section-label">链接</div>
      <div class="link-item">
        <span class="link-type">原图</span>
        <input readonly :value="img.url" class="link-input" />
        <button class="copy-btn" @click="copy(img.url)">{{ copied === img.url ? '✓' : '复制' }}</button>
      </div>
      <div class="link-item">
        <span class="link-type">优化版</span>
        <input readonly :value="img.webp_url" class="link-input" />
        <button class="copy-btn" @click="copy(img.webp_url)">{{ copied === img.webp_url ? '✓' : '复制' }}</button>
      </div>
      <div class="link-item">
        <span class="link-type">Markdown</span>
        <input readonly :value="img.markdown" class="link-input" />
        <button class="copy-btn" @click="copy(img.markdown)">{{ copied === img.markdown ? '✓' : '复制' }}</button>
      </div>
      <div class="link-item">
        <span class="link-type">HTML</span>
        <input readonly :value="img.html" class="link-input" />
        <button class="copy-btn" @click="copy(img.html)">{{ copied === img.html ? '✓' : '复制' }}</button>
      </div>
      <div class="link-item">
        <span class="link-type">BBCode</span>
        <input readonly :value="img.bbcode" class="link-input" />
        <button class="copy-btn" @click="copy(img.bbcode)">{{ copied === img.bbcode ? '✓' : '复制' }}</button>
      </div>
    </div>

    <!-- 删除 -->
    <div class="detail-section">
      <button class="del-btn" @click="handleDelete">删除图片</button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useGalleryStore } from '../stores/gallery'

const gallery = useGalleryStore()
const img = computed(() => gallery.selected)

const editTags = ref<string[]>([])
const newTag = ref('')
const addingTag = ref(false)
const tagInputRef = ref<HTMLInputElement | null>(null)
const saving = ref(false)
const copied = ref('')

const tagsChanged = computed(() => {
  if (!img.value) return false
  const orig = [...(img.value.tags ?? [])].sort().join(',')
  const curr = [...editTags.value].sort().join(',')
  return orig !== curr
})

watch(img, (v) => {
  editTags.value = [...(v?.tags ?? [])]
  addingTag.value = false
  newTag.value = ''
}, { immediate: true })

async function startAddTag() {
  addingTag.value = true
  await nextTick()
  tagInputRef.value?.focus()
}

function confirmTag() {
  const t = newTag.value.trim()
  if (t && !editTags.value.includes(t)) editTags.value.push(t)
  newTag.value = ''
  addingTag.value = false
}

function removeTag(tag: string) {
  editTags.value = editTags.value.filter((t) => t !== tag)
}

async function saveTags() {
  if (!img.value) return
  saving.value = true
  try {
    await gallery.updateTags(img.value.id, editTags.value)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  if (!img.value) return
  if (!confirm(`确定删除 ${img.value.filename}？`)) return
  await gallery.deleteImage(img.value.id)
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  copied.value = text
  setTimeout(() => { if (copied.value === text) copied.value = '' }, 1500)
}

function formatSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1048576).toFixed(1) + ' MB'
}

function formatDate(s: string) {
  return new Date(s).toLocaleString('zh-CN', { hour12: false })
}
</script>

<style scoped>
.detail {
  width: 280px;
  flex-shrink: 0;
  border-left: 1px solid var(--border);
  overflow-y: auto;
  height: calc(100vh - 52px);
  position: sticky;
  top: 52px;
  padding-bottom: 24px;
}
.detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px 10px;
  border-bottom: 1px solid var(--border);
}
.detail-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text);
}
.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
}
.icon-btn:hover { background: var(--bg-hover); color: var(--text); }
.detail-preview {
  padding: 16px;
  background: var(--bg-hover);
  border-bottom: 1px solid var(--border);
}
.preview-img {
  width: 100%;
  max-height: 200px;
  object-fit: contain;
  display: block;
  border-radius: 6px;
}
.preview-clickable {
  cursor: zoom-in;
  transition: opacity 0.15s;
}
.preview-clickable:hover { opacity: 0.85; }
.detail-meta {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
}
.meta-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  font-size: 12.5px;
  padding: 3px 0;
}
.meta-label { color: var(--text-muted); }
.meta-val { color: var(--text); font-weight: 500; max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.detail-section {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
}
.section-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-muted);
  margin-bottom: 8px;
}
.tags-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}
.tag-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background: var(--bg-hover);
  border: 1px solid var(--border);
  border-radius: 100px;
  font-size: 12px;
  color: var(--text);
}
.chip-del {
  border: none;
  background: none;
  cursor: pointer;
  color: var(--text-muted);
  font-size: 14px;
  line-height: 1;
  padding: 0;
}
.tag-input {
  padding: 2px 8px;
  border-radius: 100px;
  border: 1px solid var(--accent);
  background: var(--bg-panel);
  color: var(--text);
  font-size: 12px;
  outline: none;
  width: 80px;
}
.add-tag-btn {
  font-size: 12px;
  color: var(--text-muted);
  background: none;
  border: 1px dashed var(--border);
  border-radius: 100px;
  padding: 2px 10px;
  cursor: pointer;
}
.add-tag-btn:hover { color: var(--text); border-color: var(--accent); }
.save-btn {
  margin-top: 8px;
  width: 100%;
  padding: 6px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--accent);
  color: var(--accent-fg);
  font-size: 13px;
  cursor: pointer;
}
.save-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.link-item {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}
.link-type {
  font-size: 11px;
  color: var(--text-muted);
  width: 52px;
  flex-shrink: 0;
}
.link-input {
  flex: 1;
  font-size: 11px;
  padding: 4px 6px;
  border-radius: 5px;
  border: 1px solid var(--border);
  background: var(--bg-hover);
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}
.copy-btn {
  font-size: 11px;
  padding: 4px 8px;
  border-radius: 5px;
  border: 1px solid var(--border);
  background: var(--bg-panel);
  color: var(--text-muted);
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.12s, color 0.12s;
}
.copy-btn:hover { background: var(--bg-hover); color: var(--text); }
.del-btn {
  width: 100%;
  padding: 7px;
  border-radius: 6px;
  border: 1px solid #ef4444;
  background: transparent;
  color: #ef4444;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.12s;
}
.del-btn:hover { background: rgba(239,68,68,0.08); }
</style>
