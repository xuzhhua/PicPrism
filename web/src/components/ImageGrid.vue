<template>
  <div class="gallery">
    <div class="gallery-toolbar">
      <span class="gallery-count">{{ gallery.total }} 张图片</span>
      <select class="sort-select" v-model="sort" @change="onSortChange">
        <option value="newest">最新</option>
        <option value="oldest">最早</option>
        <option value="largest">最大</option>
      </select>
    </div>

    <div v-if="gallery.loading" class="state-msg">加载中...</div>
    <div v-else-if="!gallery.images.length" class="state-msg">暂无图片</div>

    <div v-else class="grid">
      <div
        v-for="img in gallery.images"
        :key="img.id"
        class="grid-item"
        :class="{ selected: gallery.selected?.id === img.id }"
        @click="gallery.selected = img"
      >
        <img
          :src="img.thumb_url || img.webp_url || img.url"
          :alt="img.filename"
          loading="lazy"
          class="grid-img"
        />
        <div class="grid-overlay">
          <span class="grid-filename">{{ img.filename }}</span>
        </div>
      </div>
    </div>

    <div v-if="gallery.totalPages > 1" class="pagination">
      <button
        class="page-btn"
        :disabled="gallery.page <= 1"
        @click="gallery.setPage(gallery.page - 1)"
      >上一页</button>
      <span class="page-info">{{ gallery.page }} / {{ gallery.totalPages }}</span>
      <button
        class="page-btn"
        :disabled="gallery.page >= gallery.totalPages"
        @click="gallery.setPage(gallery.page + 1)"
      >下一页</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useGalleryStore } from '../stores/gallery'

const gallery = useGalleryStore()
const sort = ref(gallery.sort)

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
}
.gallery-count {
  font-size: 13px;
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
.state-msg {
  text-align: center;
  padding: 80px 0;
  color: var(--text-muted);
  font-size: 14px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 8px;
}
.grid-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: var(--radius);
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.15s;
  background: var(--bg-hover);
}
.grid-item:hover .grid-overlay,
.grid-item.selected .grid-overlay {
  opacity: 1;
}
.grid-item.selected {
  border-color: var(--accent);
}
.grid-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
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
.page-btn:hover:not(:disabled) {
  background: var(--bg-hover);
}
.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.page-info {
  font-size: 13px;
  color: var(--text-muted);
}
</style>
