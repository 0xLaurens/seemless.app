<template>
  <div
    @click="download.close()"
    class="fixed inset-0 z-50 flex items-center justify-center bg-gray-400/25 backdrop-filter backdrop-blur"
    v-if="download.popup"
    :aria-hidden="download.popup"
  >
    <div class="">
      <div class="card w-96 bg-neutral text-neutral-content">
        <div class="card-body items-center text-center">
          <h2 class="card-title text-white font-black text-2xl pb-3">
            Username would like to send you
          </h2>
          <p class="text-lg pb-6">filename.png</p>
          <div class="card-actions justify-end">
            <button class="btn btn-outline">Deny</button>
            <button class="btn btn-primary">Accept</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDownloadStore } from '@/stores/download'

const download = useDownloadStore()
</script>
