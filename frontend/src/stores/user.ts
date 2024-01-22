import {defineStore} from 'pinia'
import {ref} from 'vue'
import type {Ref} from 'vue'
import type {User} from '@/models/user'
import {useConnStore} from "@/stores/connection";
import {useDownloadStore} from "@/stores/download";

export const useUserStore = defineStore('user', () => {
    const username = ref("")
    const users: Ref<User[]> = ref([])
    const usersMap: Ref<Map<string, boolean>> = ref(new Map())
    const currentUser: Ref<User | undefined> = ref(undefined)

    const conn = useConnStore()
    const download = useDownloadStore()

    function getUsername(): string {
        return username.value
    }

    function setUsername(name: string) {
        username.value = name
    }

    function getCurrentUser(): User | undefined {
        return currentUser.value
    }

    function setCurrentUser(user: User) {
        currentUser.value = user
    }

    function getUserByUsername(username: string) {
        return users.value.find((u) => u.username === username)
    }


    function addUser(user: User) {
        if (user.username === username.value) return
        if (usersMap.value.get(user.username)) return
        users.value.push(user)
        usersMap.value.set(user.username, true)
    }

    function removeUser(user: User) {
        conn.RemoveRtcConn(user.username)
        download.purgeOffers(user.username)
        users.value = users.value.filter((u) => u.username != user.username)
        usersMap.value.delete(user.username)
    }

    function updateUser(user: User) {
        const index = users.value.findIndex((u) => u.id === user.id)
        if (index === -1) return
        users.value[index] = user
    }

    function clearUsers() {
        users.value = []
        usersMap.value.clear()
    }

    return {
        getUsername,
        getUserByUsername,
        setUsername,
        addUser,
        removeUser,
        clearUsers,
        updateUser,
        getCurrentUser,
        setCurrentUser,
        users
    }
})
