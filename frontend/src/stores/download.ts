import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'
import type { Download } from '@/models/download'

export const useDownloadStore = defineStore('download', () => {
  const downloads: Ref<Map<string, Download>> = ref(new Map())
  const popup: Ref<boolean> = ref(false)

  function open() {
    popup.value = true
  }

  function close() {
    popup.value = false
  }

  function getDownload(): Download {
    return Array.from(downloads.value.values())[0]
  }

  function addDownload(download: Download) {
    downloads.value.set(download.file.name, download)
    open()
  }

  function removeDownload(filename: string) {
    downloads.value.delete(filename)
  }

  return {
    popup,
    open,
    close,
    addDownload,
    removeDownload,
    getDownload
  }
})
