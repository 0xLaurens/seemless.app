<script setup lang="ts">
import {ref} from 'vue'
import type {Ref} from 'vue'
import {useFileStore} from '@/stores/file'

const fileStore = useFileStore()
let files: Ref<File[]> = ref([])

function setFiles(event: Event) {
  const target = event.target as HTMLInputElement
  if (!target.files) return
  files.value = []

  for (const file of target.files) {
    files.value.push(file)
  }
}
</script>

<template>
  <div class="flex pb-6">
    <form @submit.prevent="files.length < 1" @submit="fileStore.sendFiles(files)">
      <input
          v-on:change="setFiles"
          type="file"
          class="file-input file-input-bordered file-input-primary w-full md:w-auto max-w-md md:mr-10 mb-3 md:mb-0"
      />
      <button :disabled="files.length < 1" type="submit" class="btn w-full md:btn-wide btn-primary">Send</button>
    </form>
  </div>
</template>
