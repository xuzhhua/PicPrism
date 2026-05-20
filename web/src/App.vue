<template>
  <a class="skip-link" href="#main-content">跳到内容</a>
  <div class="app-shell">
    <AppNav @upload="showUpload = true" />
    <main id="main-content" class="layout">
      <AppSidebar class="app-sidebar" />
      <ImageGrid class="app-gallery" />
      <DetailPanel class="app-detail" />
    </main>
  </div>
  <UploadZone v-if="showUpload" @close="showUpload = false" />
  <ImageLightbox v-if="gallery.lightboxOpen && gallery.selected" @close="gallery.lightboxOpen = false" />
  <AppToast />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppNav from './components/AppNav.vue'
import AppSidebar from './components/AppSidebar.vue'
import ImageGrid from './components/ImageGrid.vue'
import DetailPanel from './components/DetailPanel.vue'
import UploadZone from './components/UploadZone.vue'
import ImageLightbox from './components/ImageLightbox.vue'
import AppToast from './components/AppToast.vue'
import { useGalleryStore } from './stores/gallery'

const showUpload = ref(false)
const gallery = useGalleryStore()

onMounted(() => {
  gallery.loadImages()
  gallery.loadTags()
})
</script>

<style>
.app-shell {
  min-height: 100dvh;
  padding: 18px;
}

.layout {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr) minmax(290px, 340px);
  grid-template-areas: "sidebar gallery detail";
  gap: 18px;
  min-height: calc(100dvh - 102px);
  padding-top: 18px;
  align-items: start;
}

.app-sidebar {
  grid-area: sidebar;
}

.app-gallery {
  grid-area: gallery;
}

.app-detail {
  grid-area: detail;
}

@media (max-width: 1180px) {
  .layout {
    grid-template-columns: 208px minmax(0, 1fr);
    grid-template-areas:
      "sidebar gallery"
      "detail detail";
  }
}

@media (max-width: 920px) {
  .app-shell {
    padding: 12px;
  }

  .layout {
    grid-template-columns: minmax(0, 1fr);
    grid-template-areas:
      "gallery"
      "detail"
      "sidebar";
    gap: 14px;
    min-height: auto;
    padding-top: 14px;
  }
}
</style>
