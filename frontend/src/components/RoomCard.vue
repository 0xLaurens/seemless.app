<script setup lang="ts">
import { useModalStore } from '@/stores/modal'

const modal = useModalStore()
defineProps({
  title: String,
  description: String,
  modalTitle: String,
  modalMessage: String,
  modalBtnText: String
})
</script>

<template>
  <div
    class="card bg-base-300 w-96 shadow-2xl ml-0 m-12 py-6 hover:shadow-primary hover:cursor-pointer"
  >
    <button
      class="absolute btn btn-primary btn-circle btn-sm self-end mr-3"
      @click="modal.trigger({ title: modalTitle, message: modalMessage, btnText: modalBtnText })"
    >
      ?
    </button>
    <div class="card-body flex items-center text-center">
      <h2 class="card-title text-white text-2xl font-extrabold">{{ title || 'Title' }}</h2>
      <p class="hind text-base normal-case text-gray-500 mb-4">
        {{ description || 'Description' }}
      </p>
      <button class="btn w-full btn-primary">Connect</button>
    </div>
  </div>
</template>
