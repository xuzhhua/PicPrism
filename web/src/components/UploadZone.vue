<template>
  <div class="overlay" @click.self="$emit('close')">
    <div class="modal">
      <div class="modal-header">
        <span class="modal-title">上传图片</span>
        <button class="icon-btn" @click="$emit('close')">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div
        class="drop-zone"
        :class="{ dragging }"
        @dragover.prevent="dragging = true"
        @dragleave="dragging = false"
        @drop.prevent="onDrop"
        @click="fileInput?.click()"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="drop-icon"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>
        <p class="drop-text">拖放图片到此处，或点击选择文件</p>
        <p class="drop-hint">支持 JPG、PNG、GIF、WebP、BMP、TIFF、SVG</p>
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          multiple
          class="hidden"
          @change="onFileChange"
        />
      </div>

      <!-- 标签输入 -->
      <div class="tag-row">
        <label class="tag-label">标签（可选）</label>
        <input
          v-model="tagInput"
          class="tag-text-input"
          placeholder="多个标签用逗号分隔"
        />
      </div>

      <!-- 上传队列 -->
      <div v-if="queue.length" class="queue">
        <div v-for="item in queue" :key="item.name" class="queue-item">
          <span class="q-name">{{ item.name }}</span>
          <span class="q-status" :class="item.status">
            {{ statusText(item.status) }}
            <span v-if="item.errorMsg" class="q-error-msg">: {{ item.errorMsg }}</span>
          </span>
        </div>
      </div>

      <div class="modal-footer">
        <button
          class="upload-btn"
          :disabled="!pendingFiles.length || uploading"
          @click="startUpload"
        >
          {{ uploading ? '上传中...' : `上传 ${pendingFiles.length} 张` }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { api } from '../api/client'
import { useGalleryStore } from '../stores/gallery'

defineEmits<{ close: [] }>()

const gallery = useGalleryStore()
const fileInput = ref<HTMLInputElement | null>(null)
const dragging = ref(false)
const tagInput = ref('')
const uploading = ref(false)

interface QueueItem {
  name: string
  file: File
  status: 'pending' | 'uploading' | 'done' | 'error'
  errorMsg?: string
}

const queue = ref<QueueItem[]>([])
const pendingFiles = computed(() => queue.value.filter((i) => i.status === 'pending'))

function addFiles(files: FileList | File[]) {
  for (const f of Array.from(files)) {
    if (!queue.value.find((q) => q.name === f.name)) {
      queue.value.push({ name: f.name, file: f, status: 'pending' })
    }
  }
}

function onDrop(e: DragEvent) {
  dragging.value = false
  if (e.dataTransfer?.files) addFiles(e.dataTransfer.files)
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files) addFiles(input.files)
  input.value = ''
}

async function startUpload() {
  uploading.value = true
  const tags = tagInput.value.split(',').map((t) => t.trim()).filter(Boolean)
  for (const item of queue.value) {
    if (item.status !== 'pending') continue
    item.status = 'uploading'
    try {
      await api.uploadImage(item.file, tags)
      item.status = 'done'
    } catch (err) {
      item.status = 'error'
      item.errorMsg = err instanceof Error ? err.message : String(err)
    }
  }
  uploading.value = false
  await gallery.loadImages()
  await gallery.loadTags()
}

function statusText(s: QueueItem['status']) {
  return { pending: '等待', uploading: '上传中', done: '完成', error: '失败' }[s]
}
</script>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal {
  width: 480px;
  max-width: 95vw;
  background: var(--bg-panel);
  border-radius: 12px;
  border: 1px solid var(--border);
  overflow: hidden;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
}
.modal-title { font-size: 14px; font-weight: 600; }
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
.drop-zone {
  margin: 16px;
  border: 2px dashed var(--border);
  border-radius: 10px;
  padding: 36px 20px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;
}
.drop-zone:hover, .drop-zone.dragging {
  border-color: var(--accent);
  background: var(--bg-hover);
}
.drop-icon { margin: 0 auto 12px; color: var(--text-muted); }
.drop-text { font-size: 14px; color: var(--text); margin: 0 0 4px; }
.drop-hint { font-size: 12px; color: var(--text-muted); margin: 0; }
.hidden { display: none; }
.tag-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px 12px;
}
.tag-label { font-size: 13px; color: var(--text-muted); white-space: nowrap; }
.tag-text-input {
  flex: 1;
  padding: 6px 10px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--bg-hover);
  color: var(--text);
  font-size: 13px;
  outline: none;
}
.queue {
  overflow-y: auto;
  max-height: 180px;
  margin: 0 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
}
.queue-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 7px 12px;
  font-size: 12.5px;
  border-bottom: 1px solid var(--border);
}
.queue-item:last-child { border-bottom: none; }
.q-name { color: var(--text); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.q-status { font-size: 11px; margin-left: 8px; flex-shrink: 0; }
.q-status.pending { color: var(--text-muted); }
.q-status.uploading { color: #3b82f6; }
.q-status.done { color: #22c55e; }
.q-status.error { color: #ef4444; }
.q-error-msg { opacity: 0.85; }
.modal-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  margin-top: auto;
}
.upload-btn {
  width: 100%;
  padding: 8px;
  border-radius: 7px;
  border: none;
  background: var(--accent);
  color: var(--accent-fg);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.15s;
}
.upload-btn:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
