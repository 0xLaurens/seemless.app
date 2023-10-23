import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'
import type { User } from '@/models/user'

export const useUserStore = defineStore('user', () => {
  const username = ref(localStorage.getItem('name') || '')
  const users: Ref<User[]> = ref([])

  function getUsername(): string {
    return username.value
  }

  function setUsername(name: string) {
    if (localStorage) {
      localStorage.setItem('name', name)
    }
    username.value = name
  }

  function initUsers(newUsers: User[]) {
    users.value = newUsers
  }

  function addUser(user: User) {
    users.value.push(user)
  }

  function removeUser(user: User) {
    users.value = users.value.filter((u) => u.username != user.username)
  }

  return {
    getUsername,
    setUsername,
    initUsers,
    addUser,
    removeUser,
    users
  }
})
