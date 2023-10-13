<template>
  <div
    @click="download.close()"
    class="fixed inset-0 z-50 flex items-center justify-center bg-gray-400/25 backdrop-filter backdrop-blur"
    v-if="download.popup"
    :aria-hidden="download.popup"
  >
    <div @click.stop class="card w-96 bg-neutral text-neutral-content">
      <div class="card-body items-center text-center">
        <h2 class="card-title text-white font-black text-2xl pb-3">
          {{ download.activeDownload?.from || 'Username' }} would like to send you
        </h2>
        <p class="text-lg pb-6 break-all">
          {{ download.activeDownload?.file.name || 'Filename' }}
        </p>
        <div class="card-actions justify-end">
          <button class="btn btn-outline" @click="download.removeDownload(download.activeDownload)">
            Deny
          </button>
          <a
            class="btn btn-primary"
            @click="download.removeDownload(download.activeDownload)"
            :href="download.url"
            download
          >
            Accept
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDownloadStore } from '@/stores/download'

const download = useDownloadStore()
</script>
