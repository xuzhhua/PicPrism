<template>
  <header class="nav">
    <div class="nav-left">
      <span class="nav-kicker">self-hosted image desk</span>
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
  gap: 16px;
  min-height: 66px;
  padding: 12px 18px;
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  background: var(--bg-panel);
  box-shadow: var(--shadow-soft);
  backdrop-filter: blur(18px);
}
.nav-left {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.nav-kicker {
  font-size: 11px;
  line-height: 1.1;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--text-soft);
}
.nav-logo {
  font-weight: 700;
  font-size: 25px;
  line-height: 0.96;
  letter-spacing: -0.06em;
  color: var(--text);
}
.nav-right {
  display: flex;
  align-items: center;
  gap: 8px;
}
.icon-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 14px;
  border: 1px solid transparent;
  background: var(--bg-soft);
  color: var(--text-muted);
  cursor: pointer;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.08);
}
.icon-btn:hover {
  background: var(--bg-hover);
  border-color: var(--border);
  color: var(--text);
}
.icon-btn:active {
  transform: translateY(1px);
}
.token-missing {
  color: var(--danger);
}
.token-dot {
  position: absolute;
  top: 7px;
  right: 7px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--danger);
  border: 2px solid var(--bg-panel-strong);
}

/* Token 弹窗 */
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(17, 14, 11, 0.54);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.token-modal {
  width: min(420px, 100%);
  background: var(--bg-panel);
  border: 1px solid var(--border);
  border-radius: 24px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: var(--shadow-strong);
  backdrop-filter: blur(24px);
}
.token-modal-title {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--text);
  margin: 0;
}
.token-input {
  width: 100%;
  box-sizing: border-box;
  padding: 12px 14px;
  border: 1px solid var(--border);
  border-radius: 14px;
  background: var(--bg-soft);
  color: var(--text);
  font-size: 14px;
  outline: none;
}
.token-input:focus {
  border-color: var(--border-strong);
  box-shadow: var(--focus-ring);
}
.token-modal-footer {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}
.token-save-btn {
  padding: 10px 16px;
  border-radius: 14px;
  border: none;
  background: var(--accent);
  color: var(--accent-fg);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}
.token-cancel-btn {
  padding: 10px 16px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: var(--bg-soft);
  color: var(--text);
  font-size: 13px;
  cursor: pointer;
}
.token-cancel-btn:hover { background: var(--bg-hover); }

@media (max-width: 920px) {
  .nav {
    min-height: 60px;
    padding: 10px 14px;
    border-radius: 22px;
  }

  .nav-logo {
    font-size: 22px;
  }
}
</style>
