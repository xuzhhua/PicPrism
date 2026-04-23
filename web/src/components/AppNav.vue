<template>
  <header class="nav">
    <div class="nav-left">
      <span class="nav-logo">PicPrism</span>
    </div>
    <div class="nav-right">
      <!-- Token 未设置时显示警告徽章 -->
      <button class="icon-btn token-btn" :class="{ 'token-missing': !auth.token }" title="设置 API Token" @click="showToken = true">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
        <span v-if="!auth.token" class="token-dot" />
      </button>
      <button class="icon-btn" :title="isDark ? '切换亮色' : '切换暗色'" @click="theme.toggle()">
        <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
      </button>
      <button class="icon-btn" title="上传图片" @click="$emit('upload')">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>
      </button>
    </div>
  </header>

  <!-- Token 设置弹窗 -->
  <div v-if="showToken" class="overlay" @click.self="showToken = false">
    <div class="token-modal">
      <p class="token-modal-title">API Token</p>
      <input
        v-model="tokenInput"
        class="token-input"
        type="password"
        placeholder="输入 PICPRISM_TOKEN 的值"
        @keydown.enter="saveToken"
      />
      <div class="token-modal-footer">
        <button class="token-save-btn" @click="saveToken">保存</button>
        <button class="token-cancel-btn" @click="showToken = false">取消</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useThemeStore } from '../stores/theme'
import { useAuthStore } from '../stores/auth'

defineEmits<{ upload: [] }>()

const theme = useThemeStore()
const auth = useAuthStore()
const isDark = computed(() => theme.isDark)

const showToken = ref(false)
const tokenInput = ref(auth.token)

function saveToken() {
  auth.setToken(tokenInput.value.trim())
  showToken.value = false
}
</script>

<style scoped>
.nav {
  position: sticky;
  top: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 52px;
  background: var(--bg-panel);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(8px);
}
.nav-logo {
  font-weight: 700;
  font-size: 16px;
  letter-spacing: -0.3px;
  color: var(--text);
}
.nav-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
.icon-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: var(--radius);
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}
.icon-btn:hover {
  background: var(--bg-hover);
  color: var(--text);
}
.token-missing {
  color: #f59e0b;
}
.token-dot {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #ef4444;
  border: 1.5px solid var(--bg-panel);
}

/* Token 弹窗 */
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
}
.token-modal {
  width: 360px;
  background: var(--bg-panel);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.token-modal-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
  margin: 0;
}
.token-input {
  width: 100%;
  box-sizing: border-box;
  padding: 8px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--text);
  font-size: 13px;
  outline: none;
}
.token-input:focus { border-color: var(--accent); }
.token-modal-footer {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}
.token-save-btn {
  padding: 6px 16px;
  border-radius: var(--radius);
  border: none;
  background: var(--accent);
  color: var(--accent-fg);
  font-size: 13px;
  cursor: pointer;
}
.token-cancel-btn {
  padding: 6px 16px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: transparent;
  color: var(--text);
  font-size: 13px;
  cursor: pointer;
}
.token-cancel-btn:hover { background: var(--bg-hover); }
</style>
