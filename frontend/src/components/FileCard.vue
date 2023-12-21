<script setup lang="ts">
import {formatFileSize} from "@/utils/filesize";
import type {FileMessage} from "@/models/file";
import FileIcon from "@/components/icons/files/FileIcon.vue";
import ProgressCircle from "@/components/ProgressCircle.vue";

defineProps<{ file: FileMessage | undefined, incoming: boolean, progress: boolean, target?: string, from?: string }>()
</script>

<template>
  <div class="flex items-center" v-if="file">
    <div class="avatar placeholder">
      <div class="w-10 rounded items-center p-1">
        <file-icon :file="file"/>
      </div>
    </div>
    <div class="ms-3 text-black dark:text-white">
      <p class="font-bold break-all">{{ file?.name || "filename.jpeg" }} </p>
      <p class="font-normal" v-if="target">to {{ target }}</p>
      <p class="font-normal" v-if="from">from {{ from }}</p>
      <p>{{ formatFileSize(file?.size) || "658kb" }} • {{ file?.mime || "jpeg" }} </p>
    </div>
    <div class="ml-auto flex items-center space-x-3" v-if="progress">
      <ProgressCircle :file="file"/>
    </div>
  </div>
</template>