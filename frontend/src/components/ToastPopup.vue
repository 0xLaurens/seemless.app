<template>
  <div :class="toast.type.toString()" class="alert">
    <toast-icon :type="toast.type" />
    <span>{{ toast.message }}</span>
  </div>
</template>

<script setup lang="ts">
import type { Toast } from '@/models/toast'
import ToastIcon from '@/components/icons/ToastIcon.vue'

defineProps<{
  toast: Toast
}>()
</script>
