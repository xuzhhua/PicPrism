<template>
  <aside class="sidebar">
    <div class="sidebar-section">
      <div
        class="tag-item"
        :class="{ active: activeTag === '' }"
        @click="gallery.setTag('')"
      >
        <span class="tag-name">全部图片</span>
        <span class="tag-count">{{ gallery.total }}</span>
      </div>
    </div>
    <div class="sidebar-section" v-if="gallery.tags.length">
      <div class="sidebar-label">标签</div>
      <div
        v-for="tag in gallery.tags"
        :key="tag.id"
        class="tag-item"
        :class="{ active: activeTag === tag.name }"
        @click="gallery.setTag(tag.name)"
      >
        <span class="tag-name">{{ tag.name }}</span>
        <span class="tag-count">{{ tag.count }}</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useGalleryStore } from '../stores/gallery'

const gallery = useGalleryStore()
const activeTag = computed(() => gallery.activeTag)
</script>

<style scoped>
.sidebar {
  width: 200px;
  flex-shrink: 0;
  padding: 12px 8px;
  overflow-y: auto;
  border-right: 1px solid var(--border);
  height: calc(100vh - 52px);
  position: sticky;
  top: 52px;
}
.sidebar-section {
  margin-bottom: 16px;
}
.sidebar-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-muted);
  padding: 4px 8px;
  margin-bottom: 4px;
}
.tag-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13.5px;
  color: var(--text-muted);
  transition: background 0.12s, color 0.12s;
}
.tag-item:hover {
  background: var(--bg-hover);
  color: var(--text);
}
.tag-item.active {
  background: var(--bg-hover);
  color: var(--text);
  font-weight: 500;
}
.tag-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tag-count {
  font-size: 11px;
  color: var(--text-muted);
  margin-left: 4px;
  flex-shrink: 0;
}
</style>
