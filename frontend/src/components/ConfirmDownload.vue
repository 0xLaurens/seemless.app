<template>
  <div
      @click="download.close()"
      class="fixed inset-0 z-50 flex items-center justify-center bg-gray-400/25 backdrop-filter backdrop-blur"
      v-if="download.popup"
      :aria-hidden="download.popup"
  >
    <div @click.stop class="card w-96 bg-neutral text-neutral-content">
      <div class="card-body items-center text-center">
        <h2 class="card-title text-black text-left dark:text-white font-bold text-2xl">
          {{ download.activeOffer?.from || 'Username' }}
        </h2>
        <h2 class="card-title text-black dark:text-white font-medium text-xl">
          would like to send you {{ download.activeOffer?.files.length }}
          file(s)
        </h2>
        <div class="w-full space-y-1 max-h-60 overflow-y-scroll my-3">
          <file-preview :file="file" :key="file.name" v-for="file in download.activeOffer?.files"/>
        </div>
        <img
            class="pb-6 h-32 w-auto"
            v-if="download.activeDownload?.mime.startsWith('image')"
            :src="download.url"
            alt="preview"
        />
        <div class="card-actions justify-end">
          <button class="btn btn-outline" @click="download.denyOffer(download.activeOffer)">
            Deny
          </button>
          <a
              class="btn btn-accent"
              @click="download.acceptOffer(download.activeOffer)"
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
import {useDownloadStore} from '@/stores/download'
import FilePreview from "@/components/FilePreview.vue";

const download = useDownloadStore()
</script>
