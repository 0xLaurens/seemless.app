<script setup lang="ts">
import {ref, watch} from "vue";
import type {Ref} from "vue";

const props = defineProps<{
  from: string | undefined,
  target: string | undefined,
  file: FileMessage
}>()

const user = useUserStore()
const progress: Ref<number> = ref(0);
watch(props.file, (file) => {
  progress.value = parseFloat((file.progress / file.size * 100).toFixed(0));
})

import DocumentIcon from "@/components/icons/DocumentIcon.vue";
import type {FileMessage} from "@/models/file";
import {useUserStore} from "@/stores/user";
</script>

<template>
  <div class="card bg-neutral">
    <div class="flex py-2 px-3 items-center">
      <document-icon class="h-9 w-auto mr-3"/>
      <div class="w-full">
        <div class="flex justify-between align-bottom">
          <div class="flex space-x-2">
            <div class="flex">
              <p class="font-bold">{{
                  file.name.length < 25 ? file.name : file.name.slice(0, 10) + "..." + file.name.slice(file.name.length - 11, file.name.length)
                }}</p>
            </div>
            <p>{{ target == user.getUsername() ? `from ${from}` : `to ${target}` }}</p>
          </div>
          <p>{{ progress }}%</p>
        </div>
        <progress class="progress progress-accent w-full" :value="progress" max="100"></progress>
      </div>
      <div class="p-3 items-center">
        <svg v-if="progress == 100.0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
             stroke-width="1.5" stroke="currentColor"
             class="w-8 h-auto text-accent">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <!--        <svg v-if="progress < 100" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"-->
        <!--             stroke="currentColor"-->
        <!--             class="w-8 h-auto ">-->
        <!--          <path stroke-linecap="round" stroke-linejoin="round"-->
        <!--                d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>-->
        <!--        </svg>-->
        <svg v-if="progress < 100" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
             stroke="currentColor"
             class="w-6 h-auto">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"/>
        </svg>
      </div>
    </div>
  </div>
</template>