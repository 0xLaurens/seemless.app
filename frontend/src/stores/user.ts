import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const username = ref('')

  function getUsername(): string {
    return username
  }

  function setUsername(name: string) {
    username.value = name
  }

  return {
    getUsername,
    setUsername
  }
})
