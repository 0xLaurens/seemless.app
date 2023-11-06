<script setup lang="ts">
import DocumentIcon from "@/components/icons/DocumentIcon.vue";
import type {FileMessage} from "@/models/file";
import {formatFileSize} from "@/utils/filesize";

defineProps<{ file: FileMessage }>()
</script>

<template>
  <div class="card bg-base-100 shadow-xl w-full">
    <div class="card-body py-3 flex flex-row break-all">
      <div class="self-center">
        <document-icon/>
      </div>
      <div>
        <h2 class="card-title font-bold">{{ file.name || "example.txt" }}</h2>
        <p class="text-left font-normal">{{ formatFileSize(file.size) || "64Kb" }}</p>
      </div>
    </div>
  </div>
</template>