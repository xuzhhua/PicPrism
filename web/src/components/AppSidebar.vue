<template>
  <aside class="sidebar">
    <div class="sidebar-top">
      <p class="sidebar-eyebrow">browse</p>
      <p class="sidebar-total"><span class="sidebar-total-num mono">{{ gallery.total }}</span> images</p>
    </div>
    <div class="sidebar-section">
      <div class="tag-list">
        <div
          class="tag-item"
          :class="{ active: activeTag === '' }"
          @click="gallery.setTag('')"
        >
          <span class="tag-name">全部图片</span>
          <span class="tag-count">{{ gallery.total }}</span>
        </div>
      </div>
    </div>
    <div class="sidebar-section" v-if="gallery.tags.length">
      <div class="sidebar-label">标签</div>
      <div class="tag-list">
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
  min-width: 0;
  padding: 18px 12px;
  overflow-y: auto;
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  background: var(--bg-panel);
  box-shadow: var(--shadow-soft);
  min-height: 0;
  position: sticky;
  top: 84px;
  backdrop-filter: blur(18px);
}
.sidebar-top {
  padding: 4px 8px 16px;
  margin-bottom: 8px;
  border-bottom: 1px solid var(--border);
}
.sidebar-eyebrow {
  margin: 0 0 6px;
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--text-soft);
}
.sidebar-total {
  margin: 0;
  font-size: 22px;
  line-height: 1;
  letter-spacing: -0.05em;
  color: var(--text);
}
.sidebar-total-num {
  font-size: 26px;
}
.sidebar-section {
  margin-bottom: 18px;
}

.tag-list {
  display: grid;
  gap: 8px;
}

.sidebar-label {
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-soft);
  padding: 0 8px 8px;
  margin-bottom: 2px;
}
.tag-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 14px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-muted);
  border: 1px solid transparent;
  background: transparent;
}
.tag-item:hover {
  background: var(--bg-hover);
  border-color: var(--border);
  color: var(--text);
}
.tag-item.active {
  background: var(--bg-elevated);
  border-color: var(--border-strong);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.06);
  color: var(--text);
  font-weight: 600;
}
.tag-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tag-count {
  font-size: 11px;
  color: var(--text-soft);
  margin-left: 8px;
  flex-shrink: 0;
  font-variant-numeric: tabular-nums;
}

@media (max-width: 920px) {
  .sidebar {
    position: static;
    top: auto;
    overflow: visible;
    padding: 16px 14px;
  }

  .sidebar-top {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 12px;
    padding-bottom: 14px;
  }

  .sidebar-total {
    font-size: 18px;
    white-space: nowrap;
  }

  .sidebar-total-num {
    font-size: 21px;
  }

  .sidebar-section {
    margin-bottom: 14px;
  }

  .tag-list {
    display: flex;
    gap: 10px;
    overflow-x: auto;
    padding-bottom: 4px;
    scroll-snap-type: x proximity;
  }

  .tag-item {
    flex: 0 0 auto;
    min-width: max-content;
    scroll-snap-align: start;
  }
}
</style>
