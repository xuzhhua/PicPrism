<template>
  <AppNav @upload="showUpload = true" />
  <div class="layout">
    <AppSidebar />
    <ImageGrid />
    <DetailPanel />
  </div>
  <UploadZone v-if="showUpload" @close="showUpload = false" />
  <ImageLightbox v-if="gallery.lightboxOpen && gallery.selected" @close="gallery.lightboxOpen = false" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppNav from './components/AppNav.vue'
import AppSidebar from './components/AppSidebar.vue'
import ImageGrid from './components/ImageGrid.vue'
import DetailPanel from './components/DetailPanel.vue'
import UploadZone from './components/UploadZone.vue'
import ImageLightbox from './components/ImageLightbox.vue'
import { useGalleryStore } from './stores/gallery'

const showUpload = ref(false)
const gallery = useGalleryStore()

onMounted(() => {
  gallery.loadImages()
  gallery.loadTags()
})
</script>

<style>
.layout {
  display: flex;
  height: calc(100vh - 52px);
  overflow: hidden;
}
</style>
