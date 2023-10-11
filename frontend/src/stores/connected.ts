import { defineStore } from 'pinia'
import { ref } from 'vue'

// Store for keeping track of the users you are connected to
export const useConnectedStore = defineStore('connected', () => {
  const connected = ref(new Map<string, boolean>())
  const ice = ref(new Set<string>())

  function createUserConnection(username: string, status: boolean) {
    connected.value.set(username, status)
  }

  function addIceTarget(username: string) {
    ice.value.add(username)
  }

  function getUserConnectionStatus(username: string): boolean | undefined {
    return connected.value.get(username)
  }

  function removeUser(username: string) {
    connected.value.delete(username)
    ice.value.delete(username)
  }

  function getConnectedUsers() {
    return Array.from(connected.value.keys())
  }

  function getIceTargets() {
    return Array.from(ice.value.keys())
  }

  function clearConnectedUsers() {
    connected.value.clear()
  }

  return {
    addIceTarget,
    getIceTargets,
    createUserConnection,
    getUserConnectionStatus,
    removeUser,
    getConnectedUsers,
    clearConnectedUsers
  }
})
