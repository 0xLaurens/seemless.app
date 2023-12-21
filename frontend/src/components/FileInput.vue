<script setup lang="ts">
import {ref} from 'vue'
import type {Ref} from 'vue'
import {useFileStore} from '@/stores/file'
import PlaneIcon from "@/components/icons/PlaneIcon.vue";

const fileStore = useFileStore()
const files: Ref<FileList | undefined> = ref()

const props = defineProps<{
  users: Map<string, boolean>
}>()

function handleFiles(event: Event) {
  const target = event.target as HTMLInputElement
  if (!target.files) return;
  files.value = target.files
  fileStore.sendFilesOffer(files.value, [...props.users.keys()])
}


</script>

<template>
  <div class="flex pb-6">
    <button :disabled="users.size < 1" class="btn btn-primary w-full">
      <label for="fileInput"
             class="btn btn-ghost w-full">
        Transfer
        <plane-icon/>
        <input id="fileInput" type="file" class="hidden"
               @change="handleFiles($event)"/>
      </label>
    </button>
  </div>
</template>
