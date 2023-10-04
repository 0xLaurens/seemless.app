<script setup lang="ts">
import { useRtcStore } from '@/stores/rtc'
import { Ref, ref } from 'vue'

let rtc = useRtcStore()
let files: Ref<File[]> = ref([])

function setFiles(event) {
  files.value = []
  for (const file of event.target.files) {
    files.value.push(file)
  }
}
</script>

<template>
  <div class="flex space-y-3">
    <form @submit.prevent="files.length < 1" @submit="rtc.sendFiles(files)">
      <input
        v-on:change="setFiles"
        type="file"
        class="file-input file-input-bordered w-full max-w-xs"
      />
      <button :disabled="files.length < 1" type="submit" class="btn btn-primary">Send</button>
    </form>
  </div>
  <!--  Temporary download button-->
  <a v-if="rtc.blobURL.length > 0" target="_blank" :href="rtc.blobURL">Download Test</a>
</template>
