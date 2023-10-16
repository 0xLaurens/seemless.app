import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const username = ref(localStorage.getItem('name') || '')

  function getUsername(): string {
    return username.value
  }

  function setUsername(name: string) {
    if (localStorage) {
      localStorage.setItem('name', name)
    }
    username.value = name
  }

  return {
    getUsername,
    setUsername
  }
})
