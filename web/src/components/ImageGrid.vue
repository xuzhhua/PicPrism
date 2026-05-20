<template>
  <div class="gallery">
    <div class="gallery-toolbar">
      <!-- 多选操作栏 -->
      <template v-if="gallery.multiSelected.size > 0">
        <span class="gallery-count">已选 {{ gallery.multiSelected.size }} 张</span>
        <div class="toolbar-actions">
          <button class="action-btn" @click="selectAll">全选</button>
          <button class="action-btn danger" :disabled="deleting" @click="deleteSelected">
            {{ deleting ? '删除中...' : '删除选中' }}
          </button>
          <button class="action-btn" @click="clearSelection">取消</button>
        </div>
      </template>
      <!-- 普通工具栏 -->
      <template v-else>
        <span class="gallery-count">{{ gallery.total }} 张图片</span>
        <div class="toolbar-actions">
          <span class="hint-text">Ctrl+点击 多选</span>
          <select class="sort-select" v-model="sort" @change="onSortChange">
            <option value="newest">最新</option>
            <option value="oldest">最早</option>
            <option value="largest">最大</option>
          </select>
        </div>
      </template>
    </div>

    <div v-if="gallery.loading" class="state-panel state-loading" aria-live="polite" aria-busy="true">
      <div class="state-copy">
        <p class="state-kicker">loading</p>
        <h2 class="state-title">正在整理图片墙</h2>
        <p class="state-text">正在拉取图片、标签和当前筛选结果。</p>
      </div>
      <div class="loading-stack" aria-hidden="true">
        <span v-for="item in loadingCards" :key="item" class="loading-card" />
      </div>
    </div>
    <div v-else-if="!gallery.images.length" class="state-panel state-empty">
      <div class="state-copy">
        <p class="state-kicker">empty view</p>
        <h2 class="state-title">{{ emptyTitle }}</h2>
        <p class="state-text">{{ emptyDescription }}</p>
      </div>
      <div class="state-actions">
        <span class="state-chip">右上角可上传图片</span>
        <span v-if="gallery.activeTag" class="state-chip subtle">当前筛选：{{ gallery.activeTag }}</span>
      </div>
    </div>

    <div v-else class="masonry">
      <div
        v-for="img in gallery.images"
        :key="img.id"
        class="masonry-item"
        :class="{
          selected: gallery.selected?.id === img.id,
          'multi-checked': gallery.multiSelected.has(img.id),
          'in-multi-mode': gallery.multiSelected.size > 0,
        }"
        @click="onClickImage(img, $event)"
        @mouseleave="closeCopyMenu"
      >
        <img
          :src="img.thumb_url || img.webp_url || img.url"
          :alt="img.filename"
          loading="lazy"
          class="masonry-img"
        />
        <!-- 勾选框（多选模式下显示，或 hover 显示） -->
        <div class="check-wrap" @click.stop="toggleMulti(img.id)">
          <span class="checkbox" :class="{ checked: gallery.multiSelected.has(img.id) }">
            <svg v-if="gallery.multiSelected.has(img.id)" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
          </span>
        </div>
        <!-- 悬浮操作图标 -->
        <div class="card-icon-row" @click.stop>
          <div class="copy-wrap">
            <button class="card-icon-btn" @click.stop="toggleCopyMenu(img.id)" title="复制链接">
              <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
            </button>
            <div v-if="copyMenuFor === img.id" class="copy-menu">
              <button class="copy-menu-item" @click.stop="copyLink(img, 'url')">原图</button>
              <button class="copy-menu-item" @click.stop="copyLink(img, 'webp_url')">优化版</button>
              <button class="copy-menu-item" @click.stop="copyLink(img, 'markdown')">Markdown</button>
              <button class="copy-menu-item" @click.stop="copyLink(img, 'html')">HTML</button>
              <button class="copy-menu-item" @click.stop="copyLink(img, 'bbcode')">BBCode</button>
            </div>
          </div>
          <button class="card-icon-btn danger" @click.stop="deleteOne(img)" title="删除">
            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4h6v2"/></svg>
          </button>
        </div>
        <div class="grid-overlay">
          <span class="grid-filename">{{ img.filename }}</span>
        </div>
      </div>
    </div>

    <div v-if="gallery.totalPages > 1" class="pagination">
      <button class="page-btn" :disabled="gallery.page <= 1" @click="gallery.setPage(gallery.page - 1)">上一页</button>
      <span class="page-info">{{ gallery.page }} / {{ gallery.totalPages }}</span>
      <button class="page-btn" :disabled="gallery.page >= gallery.totalPages" @click="gallery.setPage(gallery.page + 1)">下一页</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useGalleryStore } from '../stores/gallery'
import { useToastStore } from '../stores/toast'
import type { ImageItem } from '../api/client'

const gallery = useGalleryStore()
const toast = useToastStore()
const sort = ref(gallery.sort)
const deleting = ref(false)
const copyMenuFor = ref<string | null>(null)
const loadingCards = [1, 2, 3, 4, 5, 6]

const emptyTitle = computed(() => {
  if (gallery.activeTag) {
    return `标签“${gallery.activeTag}”下还没有图片`
  }
  return '图库还没有内容'
})

const emptyDescription = computed(() => {
  if (gallery.activeTag) {
    return '切换到其他标签，或者先给图片打上这个标签。'
  }
  return '上传第一张图片后，这里会显示原图、优化版和各类外链格式。'
})

function onClickImage(img: ImageItem, e: MouseEvent) {
  if (e.ctrlKey || e.metaKey || gallery.multiSelected.size > 0) {
    toggleMulti(img.id)
    return
  }
  gallery.selected = img
}

function toggleCopyMenu(id: string) {
  copyMenuFor.value = copyMenuFor.value === id ? null : id
}

function closeCopyMenu() {
  copyMenuFor.value = null
}

async function copyLink(img: ImageItem, type: 'url' | 'webp_url' | 'markdown' | 'html' | 'bbcode') {
  const text =
    type === 'url' ? img.url
    : type === 'webp_url' ? (img.webp_url || img.url)
    : type === 'markdown' ? img.markdown
    : type === 'html' ? img.html
    : img.bbcode
  try {
    await navigator.clipboard.writeText(text)
  } catch {
    const el = document.createElement('textarea')
    el.value = text
    el.style.position = 'fixed'
    el.style.opacity = '0'
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
  }
  copyMenuFor.value = null
}

async function deleteOne(img: ImageItem) {
  if (!confirm(`确定删除「${img.filename}」？`)) return
  try {
    await gallery.deleteImage(img.id)
  } catch (e: any) {
    if (e?.status === 401) {
      toast.push('未授权：请先在右上角配置 API Token')
    } else {
      toast.push(e?.message ?? '删除失败')
    }
  }
}

function toggleMulti(id: string) {
  const s = gallery.multiSelected
  if (s.has(id)) {
    s.delete(id)
    // 触发响应式更新
    gallery.multiSelected = new Set(s)
  } else {
    gallery.multiSelected = new Set(s.add(id))
  }
}

function selectAll() {
  gallery.multiSelected = new Set(gallery.images.map((i) => i.id))
}

function clearSelection() {
  gallery.multiSelected = new Set()
}

async function deleteSelected() {
  if (!gallery.multiSelected.size) return
  if (!confirm(`确定删除选中的 ${gallery.multiSelected.size} 张图片？`)) return
  deleting.value = true
  try {
    await gallery.deleteMultiple([...gallery.multiSelected])
  } catch (e: any) {
    if (e?.status === 401) {
      toast.push('未授权：请先在右上角配置 API Token')
    } else {
      toast.push(e?.message ?? '删除失败')
    }
  } finally {
    deleting.value = false
  }
}

function onSortChange() {
  gallery.sort = sort.value
  gallery.page = 1
  gallery.loadImages()
}
</script>

<style scoped>
.gallery {
  min-width: 0;
  min-height: 0;
  padding: 20px;
  overflow-y: auto;
  border: 1px solid var(--border);
  border-radius: calc(var(--radius-lg) + 2px);
  background: var(--bg-panel);
  box-shadow: var(--shadow-soft);
  backdrop-filter: blur(18px);
}

.gallery-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 20px;
  min-height: 48px;
  padding: 14px 16px;
  border: 1px solid var(--border);
  border-radius: 20px;
  background: var(--bg-elevated);
}

.gallery-count {
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 0.01em;
  color: var(--text-muted);
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.hint-text {
  font-size: 12px;
  color: var(--text-soft);
}

.sort-select {
  min-height: 38px;
  padding: 0 40px 0 12px;
  border-radius: 14px;
  border: 1px solid var(--border);
  appearance: none;
  -webkit-appearance: none;
  background-color: var(--bg-soft);
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='14' height='14' viewBox='0 0 14 14' fill='none'%3E%3Cpath d='M3.25 5.5L7 9.25L10.75 5.5' stroke='%23b1a79b' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
  background-position: right 12px center;
  background-repeat: no-repeat;
  background-size: 14px 14px;
  color: var(--text);
  font-size: 13px;
  cursor: pointer;
  outline: none;
}

.sort-select option {
  background: var(--bg-panel-strong);
  color: var(--text);
}

.action-btn {
  min-height: 38px;
  padding: 0 14px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: var(--bg-soft);
  color: var(--text);
  font-size: 13px;
  cursor: pointer;
}

.action-btn:hover:not(:disabled) {
  background: var(--bg-hover);
  border-color: var(--border-strong);
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.action-btn.danger {
  background: var(--danger);
  border-color: var(--danger);
  color: #fff8f6;
}

.action-btn.danger:hover:not(:disabled) {
  background: #7a4038;
}

.state-panel {
  display: grid;
  gap: 20px;
  padding: 28px;
  border: 1px solid var(--border);
  border-radius: 28px;
  background: linear-gradient(180deg, var(--bg-soft), transparent);
}

.state-copy {
  max-width: 38rem;
}

.state-kicker {
  margin: 0 0 10px;
  color: var(--text-soft);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.state-title {
  margin: 0 0 8px;
  color: var(--text);
  font-size: 30px;
  line-height: 0.98;
  letter-spacing: -0.05em;
}

.state-text {
  margin: 0;
  color: var(--text-muted);
  font-size: 14px;
  line-height: 1.6;
}

.state-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.state-chip {
  display: inline-flex;
  align-items: center;
  min-height: 38px;
  padding: 0 14px;
  border: 1px solid var(--border-strong);
  border-radius: 999px;
  background: var(--bg-panel-strong);
  color: var(--text);
  font-size: 12px;
}

.state-chip.subtle {
  color: var(--text-muted);
  background: transparent;
}

.loading-stack {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.loading-card {
  display: block;
  min-height: 120px;
  border-radius: 22px;
  background: linear-gradient(110deg, var(--bg-soft) 10%, var(--bg-panel-strong) 28%, var(--bg-soft) 42%);
  background-size: 200% 100%;
  animation: shimmer 1.2s linear infinite;
}

.masonry {
  columns: 4 220px;
  column-gap: 16px;
}

.masonry-item {
  position: relative;
  break-inside: avoid;
  margin-bottom: 16px;
  border: 1px solid transparent;
  border-radius: 24px;
  overflow: visible;
  background: var(--bg-hover);
  box-shadow: 0 10px 24px rgba(49, 40, 31, 0.08);
  cursor: pointer;
  transition: border-color 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
}

.masonry-item:hover {
  transform: translateY(-3px);
  border-color: var(--border);
  box-shadow: 0 22px 38px rgba(49, 40, 31, 0.16);
}

.masonry-item.in-multi-mode {
  cursor: pointer;
}

.masonry-item:hover .grid-overlay,
.masonry-item.selected .grid-overlay {
  opacity: 1;
}

.card-icon-row {
  position: absolute;
  right: 12px;
  bottom: 18px;
  z-index: 4;
  display: flex;
  gap: 8px;
  opacity: 0;
  transform: translateY(6px);
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.masonry-item:hover .card-icon-row {
  opacity: 1;
  transform: translateY(0);
}

.masonry-item.in-multi-mode:hover .card-icon-row {
  opacity: 0;
}

.card-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 1px solid rgba(255, 248, 239, 0.18);
  border-radius: 12px;
  background: rgba(20, 17, 14, 0.68);
  color: #fff8f1;
  cursor: pointer;
  backdrop-filter: blur(12px);
}

.card-icon-btn:hover {
  background: rgba(20, 17, 14, 0.84);
}

.card-icon-btn.danger {
  border-color: rgba(197, 141, 132, 0.35);
  background: rgba(138, 75, 66, 0.72);
}

.card-icon-btn.danger:hover {
  background: rgba(138, 75, 66, 0.92);
}

.copy-wrap {
  position: relative;
}

.copy-menu {
  position: absolute;
  right: 0;
  bottom: calc(100% + 10px);
  z-index: 10;
  display: flex;
  min-width: 126px;
  flex-direction: column;
  gap: 4px;
  padding: 6px;
  border: 1px solid rgba(255, 248, 239, 0.12);
  border-radius: 16px;
  background: rgba(24, 21, 19, 0.94);
  box-shadow: 0 18px 34px rgba(0, 0, 0, 0.32);
  backdrop-filter: blur(16px);
}

.copy-menu-item {
  padding: 8px 10px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: #fff7ef;
  font-size: 12px;
  text-align: left;
  white-space: nowrap;
  cursor: pointer;
}

.copy-menu-item:hover {
  background: rgba(255, 255, 255, 0.09);
}

.masonry-item.selected {
  border-color: var(--border-strong);
  box-shadow: 0 0 0 2px rgba(33, 29, 24, 0.12), 0 22px 38px rgba(49, 40, 31, 0.16);
}

.masonry-item.multi-checked {
  border-color: var(--border-strong);
}

.masonry-item.multi-checked .masonry-img {
  opacity: 0.78;
}

.masonry-img {
  display: block;
  width: 100%;
  height: auto;
  border-radius: 24px;
  transition: opacity 0.2s ease;
}

.check-wrap {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 2;
  opacity: 0;
  transition: opacity 0.18s ease;
}

.masonry-item:hover .check-wrap,
.masonry-item.multi-checked .check-wrap,
.masonry-item.in-multi-mode .check-wrap {
  opacity: 1;
}

.checkbox {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: 1px solid rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  background: rgba(20, 17, 14, 0.58);
  color: #fff;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.24);
}

.checkbox.checked {
  background: var(--accent);
  border-color: var(--accent);
}

.grid-overlay {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  padding: 22px 14px 14px;
  border-radius: 0 0 24px 24px;
  background: linear-gradient(180deg, transparent, rgba(20, 17, 14, 0.82));
  opacity: 0;
  transition: opacity 0.18s ease;
}

.grid-filename {
  display: block;
  overflow: hidden;
  color: #fffaf3;
  font-size: 12px;
  letter-spacing: -0.01em;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-top: 24px;
  padding-bottom: 24px;
}

.page-btn {
  min-height: 38px;
  padding: 0 14px;
  border: 1px solid var(--border);
  border-radius: 14px;
  background: var(--bg-soft);
  color: var(--text);
  cursor: pointer;
}

.page-btn:hover:not(:disabled) {
  background: var(--bg-hover);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  color: var(--text-muted);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

@media (max-width: 920px) {
  .gallery {
    overflow: visible;
    padding: 16px;
  }

  .gallery-toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .toolbar-actions {
    justify-content: flex-start;
  }

  .masonry {
    columns: 2 180px;
  }

  .state-panel {
    padding: 20px;
  }

  .state-title {
    font-size: 24px;
  }

  .loading-stack {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@keyframes shimmer {
  from {
    background-position: 100% 0;
  }
  to {
    background-position: -100% 0;
  }
}
</style>
