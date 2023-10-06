<script setup lang="ts">
import { useUserStore } from '@/stores/user'

const props = defineProps({ user: String })
const connected: boolean = false
const userStore = useUserStore()
const isUser = userStore.getUsername() === props.user
</script>

<template>
  <div class="user justify-center text-center mr-6">
    <div class="avatar placeholder mb-2">
      <div
        :class="{
          'hover:ring-8 ring-accent ring-offset-4 ring-offset-base-100 cursor-pointer': !isUser,
          'ring-8 ring-accent': connected
        }"
        class="bg-neutral-focus text-neutral-content rounded-full w-24"
      >
        <span class="text-3xl">{{ props.user[0] }}</span>
      </div>
    </div>
    <h2>{{ props.user }} <span v-if="isUser">(YOU)</span></h2>
    <span class="hind text-accent text-lg font-medium">Android</span>
  </div>
</template>
