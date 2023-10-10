import { defineStore } from 'pinia'
import { ref } from 'vue'

// Store for keeping track of the users you are connected to
export const useConnectedStore = defineStore('connected', () => {
  const connected = ref(new Map<string, boolean>())

  function createUserConnection(username: string, status: boolean) {
    connected.value.set(username, status)
  }

  function getUserConnectionStatus(username: string): boolean | undefined {
    return connected.value.get(username)
  }

  function removeUser(username: string) {
    connected.value.delete(username)
  }

  function getConnectedUsers() {
    return Array.from(connected.value.keys())
  }

  return {
    createUserConnection,
    getUserConnectionStatus,
    removeUser,
    getConnectedUsers
  }
})
