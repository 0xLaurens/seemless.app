<script setup lang="ts">
import DocumentIcon from "@/components/icons/DocumentIcon.vue";
import type {FileMessage} from "@/models/file";
import {useUserStore} from "@/stores/user";
import {computed} from "vue";

const props = defineProps<{
  from: string | undefined,
  target: string | undefined,
  file: FileMessage
}>()

const progress = computed(() => {
  return Math.round((props.file.progress / props.file.size * 100))
});

const user = useUserStore()
</script>

<template>
  <div class="card bg-neutral w-full">
    <div class="flex py-1 md:py-3 px-3 items-center">
      <document-icon class="h-9 w-9 mr-3"/>
      <div class="w-full">
        <div class="flex justify-between">
          <div class="flex space-x-3">
            <p class="font-bold truncate max-w-sm">{{ file.name }}</p>
            <p class="hidden lg:flex">{{ target == user.getUsername() ? `from ${from}` : `to ${target}` }}</p>
          </div>
          <p class="hidden lg:flex">{{ progress }}%</p>
        </div>
        <progress class="progress progress-accent w-full" :value="Math.round((file.progress / file.size * 100))"
                  max="100"></progress>
        <div class="flex space-x-3 justify-between text-sm lg:hidden">
          <p>{{ target == user.getUsername() ? `from ${from}` : `to ${target}` }}</p>
          <p>{{ progress }}%</p>
        </div>
      </div>
      <div class="ml-3 items-center">
        <svg v-if="progress == 100.0" xmlns="http://www.w3.org/2000/svg"
             fill="none" viewBox="0 0 24 24"
             stroke-width="1.5" stroke="currentColor"
             class="w-8 h-auto text-accent">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-if="progress < 100" xmlns="http://www.w3.org/2000/svg" fill="none"
             viewBox="0 0 24 24" stroke-width="1.5"
             stroke="currentColor"
             class="w-6 h-auto">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"/>
        </svg>
      </div>
    </div>
  </div>
</template>