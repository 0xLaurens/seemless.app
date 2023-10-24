<script setup lang="ts">
import {useUserStore} from '@/stores/user'
import type {User} from '@/models/user'
import {useConnStore} from '@/stores/connection'
import {useWebsocketStore} from '@/stores/websocket'

const props = defineProps<{ user: User }>()

const conn = useConnStore()
const ws = useWebsocketStore()

async function sendOffer(username: string) {
  const offer = await conn.CreateRtcOffer(username)
  ws.SendMessage(offer)
}

const userStore = useUserStore()
const isUser = userStore.getUsername() === props.user.username
</script>

<template>
  <div
      class="user justify-center text-center mr-3 md:mr-6 group w-36"
      :class="{ 'cursor-pointer': !isUser, 'cursor-default': isUser }"
      @click="!isUser ? sendOffer(user.username) : null"
  >
    <div class="avatar placeholder mb-2">
      <div
          :class="{
          'group-hover:ring-4 ring-accent ring-offset-4 ring-offset-base-100': !isUser,
          'ring-4 ring-accent': conn.GetUserConnection(user.username)
        }"
          class="bg-neutral-focus text-neutral-content rounded-full w-16 md:w-24"
      >
        <span class="text-3xl">{{ user.username.charAt(0) }}</span>
      </div>
    </div>
    <h2 class="break-words">{{ user.username }} <span v-if="isUser">(YOU)</span></h2>
    <span class="hind text-accent text-lg font-medium">{{ user.device }}</span>
  </div>
</template>
