<script setup lang="ts">
import {useRoute} from 'vue-router'
import QrIcon from '@/components/icons/QrIcon.vue'
import PlaneIcon from '@/components/icons/PlaneIcon.vue'
import BackIcon from '@/components/icons/BackIcon.vue'
import {useUserStore} from '@/stores/user'
import {onMounted, onUnmounted} from 'vue'
import FileInput from '@/components/FileInput.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import ConfirmDownload from '@/components/ConfirmDownload.vue'
import {useWebsocketStore} from '@/stores/websocket'
import WsConnection from '@/components/WsConnection.vue'

const user = useUserStore()
const ws = useWebsocketStore()

const route = useRoute()
const id = route.params.id

onMounted(() => {
  ws.Open()
})

onUnmounted(() => {
  ws.Close()
})
</script>

<template>
  <confirm-download/>
  <div class="pt-10 pb-4 sm:pt-16 lg:overflow-hidden lg:pt-8 lg:pb-4">
    <div class="mx-auto max-w-7xl lg:px-8">
      <div class="w-full px-4 sm:px-6 sm:text-center lg:px-0">
        <div class="relative lg:pt-12">
          <div class="flex justify-between align-middle">
            <router-link to="../network" class="btn btn-md md:btn-lg btn-ghost font-black">
              <back-icon/>
            </router-link>
            <h1 class="text-4xl font-black text-center text-black dark:text-white capitalize mb-24">
              Room: {{ id }}
            </h1>
            <button class="btn btn-md md:btn-lg btn-accent btn-outline">
              <qr-icon/>
            </button>
          </div>
        </div>

        <div class="users-box flow-root mb-24">
          <div class="flex flex-wrap justify-center align-middle">
            <div :key="u.username" v-for="u in user.users">
              <user-avatar :user="u"/>
            </div>
          </div>
        </div>

        <div>
          <ws-connection/>
        </div>

        <div class="relative mb-24">
          <div class="flex justify-center flex-wrap space-x-6 pt-5">
            <file-input/>
          </div>
        </div>

        <div class="relative mb-24">
          <div class="flex justify-center flex-wrap space-x-6 pt-5">
            <button class="btn btn-md md:btn-lg btn-primary">Upload File</button>
            <button class="btn btn-md md:btn-lg btn-primary">
              Broadcast
              <plane-icon/>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
