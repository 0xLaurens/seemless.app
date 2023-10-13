import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'
import type { Download } from '@/models/download'

export const useDownloadStore = defineStore('download', () => {
  const downloads: Ref<Map<string, Download>> = ref(new Map())
  const popup: Ref<boolean> = ref(false)
  const url: Ref<string | undefined> = ref()
  const activeDownload: Ref<Download | undefined> = ref()

  function open() {
    popup.value = true
  }

  function close() {
    popup.value = false
  }

  function setLatestDownload() {
    const latest =
      Array.from(downloads.value.values()).length > 1
        ? undefined
        : Array.from(downloads.value.values())[0]
    url.value = latest === undefined ? latest : URL.createObjectURL(latest.file)
    activeDownload.value = latest
  }

  function addDownload(download: Download) {
    downloads.value.set(download.file.name, download)
    setLatestDownload()
    open()
  }

  function removeDownload(download: Download | undefined) {
    if (download === undefined) return
    downloads.value.delete(download.file.name)
    setLatestDownload()
    close()
  }

  return {
    popup,
    url,
    activeDownload,
    open,
    close,
    addDownload,
    removeDownload
  }
})
