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

    <div v-if="gallery.loading" class="state-msg">加载中...</div>
    <div v-else-if="!gallery.images.length" class="state-msg">暂无图片</div>

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
import { ref } from 'vue'
import { useGalleryStore } from '../stores/gallery'
import { useToastStore } from '../stores/toast'
import type { ImageItem } from '../api/client'

const gallery = useGalleryStore()
const toast = useToastStore()
const sort = ref(gallery.sort)
const deleting = ref(false)
const copyMenuFor = ref<string | null>(null)

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
  flex: 1;
  padding: 16px 20px;
  overflow-y: auto;
  min-width: 0;
}
.gallery-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  min-height: 32px;
}
.gallery-count {
  font-size: 13px;
  color: var(--text-muted);
}
.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
.hint-text {
  font-size: 12px;
  color: var(--text-muted);
}
.sort-select {
  font-size: 13px;
  padding: 4px 8px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--bg-panel);
  color: var(--text);
  cursor: pointer;
  outline: none;
}
.action-btn {
  font-size: 13px;
  padding: 5px 12px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--bg-panel);
  color: var(--text);
  cursor: pointer;
  transition: background 0.12s;
}
.action-btn:hover:not(:disabled) { background: var(--bg-hover); }
.action-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.action-btn.danger {
  background: #c0392b;
  border-color: #c0392b;
  color: #fff;
}
.action-btn.danger:hover:not(:disabled) { background: #a93226; }
.state-msg {
  text-align: center;
  padding: 80px 0;
  color: var(--text-muted);
  font-size: 14px;
}
.masonry {
  columns: 4 180px;
  column-gap: 10px;
}
.masonry-item {
  break-inside: avoid;
  margin-bottom: 10px;
  position: relative;
  border-radius: var(--radius);
  overflow: visible;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.15s, transform 0.15s, box-shadow 0.15s;
  background: var(--bg-hover);
}
.masonry-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0,0,0,0.3);
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
  bottom: 28px;
  right: 6px;
  z-index: 4;
  display: flex;
  gap: 5px;
  opacity: 0;
  transition: opacity 0.15s;
}
.masonry-item:hover .card-icon-row {
  opacity: 1;
}
.masonry-item.in-multi-mode:hover .card-icon-row {
  opacity: 0;
}
.card-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 6px;
  border: 1px solid rgba(255,255,255,0.3);
  background: rgba(0,0,0,0.6);
  color: #fff;
  cursor: pointer;
  backdrop-filter: blur(4px);
  transition: background 0.12s;
}
.card-icon-btn:hover { background: rgba(0,0,0,0.82); }
.card-icon-btn.danger {
  border-color: rgba(255,80,80,0.5);
  background: rgba(192,57,43,0.65);
}
.card-icon-btn.danger:hover { background: rgba(192,57,43,0.9); }
.copy-wrap {
  position: relative;
}
.copy-menu {
  position: absolute;
  bottom: calc(100% + 6px);
  right: 0;
  background: rgba(18,18,18,0.96);
  border: 1px solid rgba(255,255,255,0.13);
  border-radius: 7px;
  padding: 4px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 106px;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 18px rgba(0,0,0,0.55);
  z-index: 10;
}
.copy-menu-item {
  padding: 5px 10px;
  border: none;
  background: transparent;
  color: #fff;
  font-size: 12px;
  text-align: left;
  cursor: pointer;
  border-radius: 4px;
  white-space: nowrap;
  transition: background 0.1s;
}
.copy-menu-item:hover { background: rgba(255,255,255,0.13); }
.masonry-item.selected {
  border-color: #3b82f6;
}
.masonry-item.multi-checked {
  border-color: #3b82f6;
}
.masonry-item.multi-checked .masonry-img {
  opacity: 0.75;
}
.masonry-img {
  width: 100%;
  height: auto;
  display: block;
  border-radius: var(--radius);
  transition: opacity 0.15s;
}
/* 勾选框 */
.check-wrap {
  position: absolute;
  top: 6px;
  left: 6px;
  z-index: 2;
  opacity: 0;
  transition: opacity 0.15s;
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
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid rgba(255,255,255,0.95);
  background: rgba(0,0,0,0.45);
  color: #fff;
  transition: background 0.15s, border-color 0.15s;
  box-shadow: 0 1px 6px rgba(0,0,0,0.6);
}
.checkbox.checked {
  background: #3b82f6;
  border-color: #3b82f6;
  color: #fff;
}
.grid-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 6px 8px;
  background: linear-gradient(transparent, rgba(0,0,0,0.6));
  opacity: 0;
  transition: opacity 0.15s;
}
.grid-filename {
  font-size: 11px;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
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
  padding: 6px 14px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--bg-panel);
  color: var(--text);
  font-size: 13px;
  cursor: pointer;
  transition: background 0.12s;
}
.page-btn:hover:not(:disabled) { background: var(--bg-hover); }
.page-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.page-info { font-size: 13px; color: var(--text-muted); }
</style>
