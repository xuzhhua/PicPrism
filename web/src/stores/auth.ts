import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('picprism_token') ?? '')

  function setToken(t: string) {
    token.value = t
    localStorage.setItem('picprism_token', t)
  }

  function clearToken() {
    token.value = ''
    localStorage.removeItem('picprism_token')
  }

  return { token, setToken, clearToken }
})
