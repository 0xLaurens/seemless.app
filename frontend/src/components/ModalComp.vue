<script setup lang="ts">
import { useModalStore } from '@/stores/modal'

const modal = useModalStore()
</script>

<template>
  <div
    v-if="modal.getActiveState()"
    class="fixed inset-0 z-50 flex items-center justify-center bg-gray-400/25 backdrop-filter backdrop-blur"
    :aria-hidden="modal.getActiveState()"
    @click="modal.close()"
  >
    <dialog @click.stop class="card h-[50%] max-w-xl bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title text-3xl font-black">{{ modal.getTitle() || 'Title' }}</h2>
        <p class="py-4 hind text-xl">{{ modal.getMessage() || 'description' }}</p>
        <div class="flex justify-between">
          <button @click="modal.close()" class="btn w-full btn-lg btn-primary capitalize">
            {{ modal.getBtnText() }}
          </button>
        </div>
      </div>
    </dialog>
  </div>
</template>
