<template>
  <dialog
      @click="download.close()"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/10 backdrop-filter backdrop-blur"
      v-if="download.popup"
      :aria-hidden="download.popup"
  >
    <div @click.stop class="card w-96 bg-base-100 dark:bg-neutral text-neutral-content">
      <div class="card-body items-center text-center">
        <h2 class="card-title text-black text-left dark:text-white font-bold text-2xl">
          {{ download.activeOffer?.from || 'Username' }}
        </h2>
        <h2 class="card-title text-black dark:text-white font-medium text-xl">
          would like to send you {{ download.activeOffer?.files.length }}
          file(s)
        </h2>
        <div class="w-full space-y-1 max-h-60 overflow-y-scroll my-3">
          <file-card :file="file" :progress="false" :key="file.name" v-for="file in download.activeOffer?.files"
                     incoming/>
        </div>
        <img
            class="pb-6 h-32 w-auto"
            v-if="download.activeDownload?.mime.startsWith('image')"
            :src="download.url"
            alt="preview"
        />
        <div class="card-actions justify-end">
          <button aria-label="deny files" class="btn btn-outline" @click="download.denyOffer(download.activeOffer)">
            Deny
          </button>
          <a
              aria-label="accept files"
              class="btn btn-primary"
              @click="download.acceptOffer(download.activeOffer)"
              :href="download.url"
              download
          >
            Accept
          </a>
        </div>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import {useDownloadStore} from '@/stores/download'
import FileCard from "@/components/FileCard.vue";

const download = useDownloadStore()
</script>
