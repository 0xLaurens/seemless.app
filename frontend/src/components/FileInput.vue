<script setup lang="ts">
import {ref} from 'vue'
import type {Ref} from 'vue'
import {useFileStore} from '@/stores/file'

const fileStore = useFileStore()
const files: Ref<FileList | undefined> = ref()

function handleFiles(event: Event) {
  const target = event.target as HTMLInputElement
  if (!target.files) return;
  files.value = target.files
}

</script>

<template>
  <div class="flex pb-6">
    <form @submit.prevent="!files || files.length < 1" @submit="fileStore.sendFilesOffer(files)">
      <input
          type="file"
          @change="handleFiles($event)"
          multiple
          class="file-input file-input-bordered file-input-accent w-full md:w-auto max-w-md md:mr-10 mb-3 md:mb-0"
      />
      <button :disabled="!files || files.length < 1" type="submit" class="btn w-full md:btn-wide btn-primary">Send
      </button>
    </form>
  </div>
</template>
