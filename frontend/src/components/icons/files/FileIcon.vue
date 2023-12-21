<script setup lang="ts">
import {type FileMessage} from "@/models/file";
import ImageIcon from "@/components/icons/files/ImageIcon.vue";
import MusicIcon from "@/components/icons/files/MusicIcon.vue";
import DocumentIcon from "@/components/icons/files/DocumentIcon.vue";

defineProps<{ file: FileMessage | undefined }>()
</script>

<template>
  <image-icon class="h-fit invert dark:invert-0 w-auto text-white" v-if="file?.mime.startsWith('image/')"/>
  <music-icon class="h-fit invert dark:invert-0 w-auto text-white" v-else-if="file?.mime.startsWith('audio/')"/>
  <document-icon v-else class="h-fit invert dark:invert-0 w-auto text-white"/>
</template>