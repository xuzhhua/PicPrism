<template>
  <Teleport to="body">
    <div class="toast-container">
      <transition-group name="toast">
        <div
          v-for="t in toast.toasts"
          :key="t.id"
          class="toast-item"
          :class="t.type"
          @click="toast.remove(t.id)"
        >
          <svg v-if="t.type === 'error'" xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          <svg v-else-if="t.type === 'success'" xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
          <span>{{ t.message }}</span>
        </div>
      </transition-group>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useToastStore } from '../stores/toast'
const toast = useToastStore()
</script>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  z-index: 9999;
  pointer-events: none;
}
.toast-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  border-radius: 8px;
  font-size: 13px;
  color: #fff;
  cursor: pointer;
  pointer-events: all;
  box-shadow: 0 4px 20px rgba(0,0,0,0.35);
  max-width: 380px;
  word-break: break-word;
  backdrop-filter: blur(6px);
}
.toast-item.error   { background: rgba(180, 40, 30, 0.92); }
.toast-item.success { background: rgba(30, 130, 70, 0.92); }
.toast-item.info    { background: rgba(40, 80, 160, 0.92); }

.toast-enter-active, .toast-leave-active { transition: all 0.22s ease; }
.toast-enter-from { opacity: 0; transform: translateY(12px); }
.toast-leave-to   { opacity: 0; transform: translateY(12px); }
</style>
