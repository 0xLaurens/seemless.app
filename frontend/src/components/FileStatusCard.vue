<script setup lang="ts">
import {ref} from "vue";

const progress = ref(0)
const receiving = ref(false)
const username = ref("username");

import DocumentIcon from "@/components/icons/DocumentIcon.vue";
</script>

<template>
  <div class="card bg-neutral">
    <div class="flex py-2 px-3 items-center">
      <document-icon class="h-9 w-auto mr-3"/>
      <div class="w-full">
        <div class="flex justify-between align-bottom">
          <div class="flex space-x-2">
            <div class="flex">
              <p class="font-bold">Files.txt</p>
            </div>
            <p>{{ receiving ? "from" : "to" }} {{ username }}</p>
          </div>
          <p>{{ progress }}%</p>
        </div>
        <progress class="progress progress-accent w-full" :value="progress" max="100"></progress>
      </div>
      <div class="p-3 items-center">
        <svg v-if="progress == 100" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
             stroke-width="1.5" stroke="currentColor"
             class="w-8 h-auto text-accent">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-if="progress < 100" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
             stroke="currentColor"
             class="w-8 h-auto ">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
      </div>
    </div>
  </div>
</template>