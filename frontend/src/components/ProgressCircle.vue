<script setup lang="ts">
import type {FileMessage} from "@/models/file";
import {computed} from "vue";

const props = defineProps<{
  file: FileMessage
}>()

const progress = computed(() => {
  return Math.round((props.file.progress / props.file.size * 100))
});
</script>

<template>
  <div class="radial-progress text-primary" :style="`--value:${progress};`" role="progressbar">{{ progress }}%</div>
</template>