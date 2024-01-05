<template>
  <dialog
      @mouseleave="store.unfreeze(toast.id)"
      @mouseover.once="store.freeze(toast.timeoutId)"
      :class="toast.type.toString()"
      class="alert cursor-pointer"
  >
    <toast-icon :type="toast.type"/>
    <span>{{ toast.message }}</span>
  </dialog>
</template>

<script setup lang="ts">
import type {Toast} from '@/models/toast'
import ToastIcon from '@/components/icons/ToastIcon.vue'
import {useToastStore} from '@/stores/toast'

const store = useToastStore()

defineProps<{
  toast: Toast
}>()
</script>
