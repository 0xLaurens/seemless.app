<script setup lang="ts">
import {type User, UserTransfer} from "@/models/user";
import DeviceIcon from "@/components/icons/devices/DeviceIcon.vue";
import {useDownloadStore} from "@/stores/download";
import {onUnmounted, ref} from "vue";

const props = defineProps<{
  user: User
}>()

const selected = ref(false)

const emit = defineEmits(['selected', 'removeUser'])

const download = useDownloadStore()

onUnmounted(() => {
  emit("removeUser", props.user.username)
})
</script>

<template>
  <div class="flex items-center hover:bg-base-100 hover:cursor-pointer px-3 py-0.5 rounded-2xl"
       @click="selected = !selected; $emit('selected', selected, user.username)">
    <div class="avatar placeholder">
      <div class="w-12 rounded items-center p-1">
        <device-icon :user="user"/>
      </div>
    </div>
    <div class="ms-3">
      <p class="font-bold">{{ user.username }}</p>
      <p>{{ user.device }}</p>
      <p v-if="download.activeOffer">{{ UserTransfer.Transfer }}</p>
    </div>
    <div class="ml-auto flex items-center space-x-3">
      <input type="checkbox" :aria-label="`send to ${user.username}`" class="checkbox checkbox-primary"
             @change="$emit('selected', selected, user.username)"
             v-model="selected"/>
    </div>
  </div>
</template>