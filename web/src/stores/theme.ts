import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usePreferredDark } from '@vueuse/core'

export const useThemeStore = defineStore('theme', () => {
  const prefersDark = usePreferredDark()
  const isDark = ref(
    localStorage.getItem('picprism_theme') === 'dark' ||
    (localStorage.getItem('picprism_theme') === null && prefersDark.value)
  )

  function apply() {
    document.documentElement.classList.toggle('dark', isDark.value)
  }

  function toggle() {
    isDark.value = !isDark.value
    localStorage.setItem('picprism_theme', isDark.value ? 'dark' : 'light')
    apply()
  }

  apply()

  return { isDark, toggle }
})
