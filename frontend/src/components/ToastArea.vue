<template>
  <div class="toast hidden lg:flex toast-top toast-end z-50">
    <toast-popup v-bind:key="t.id" v-for="t in toast.toasts" :toast="t"/>
  </div>
</template>

<script setup lang="ts">
import {useToastStore} from '@/stores/toast'
import ToastPopup from '@/components/ToastPopup.vue'

const toast = useToastStore()
</script>
