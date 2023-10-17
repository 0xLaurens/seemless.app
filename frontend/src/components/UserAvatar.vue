<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import { useConnectedStore } from '@/stores/connected'
import type { User } from '@/models/user'

const connStore = useConnectedStore()
const props = defineProps<{ user: User }>()

const userStore = useUserStore()
const isUser = userStore.getUsername() === props.user.username
</script>

<template>
  <div
    class="user justify-center text-center mr-6 group"
    :class="{ 'cursor-pointer': !isUser, 'cursor-default': isUser }"
  >
    <div class="avatar placeholder mb-2">
      <div
        :class="{
          'group-hover:ring-4 ring-accent ring-offset-4 ring-offset-base-100': !isUser,
          'ring-4 ring-accent': connStore.getUserConnectionStatus(props.user.username)
        }"
        class="bg-neutral-focus text-neutral-content rounded-full w-24"
      >
        <span class="text-3xl">{{ user.username.charAt(0) }}</span>
      </div>
    </div>
    <h2>{{ user.username }} <span v-if="isUser">(YOU)</span></h2>
    <span class="hind text-accent text-lg font-medium">{{ user.device }}</span>
  </div>
</template>
