import type { FileMessage } from '@/models/file'
import { FileStatus } from '@/models/file'
import type { Download } from '@/models/download'
import type { Ref } from 'vue'
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useDownloadStore } from '@/stores/download'
import { useConnStore } from '@/stores/connection'
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast'
import { ToastType } from '@/models/toast'

export const useFileStore = defineStore('file', () => {
  const downloadStore = useDownloadStore()
  const conn = useConnStore()
  const user = useUserStore()
  const toast = useToastStore()
  const CHUNK_SIZE = 65536 //64 KiB

  const buf: Ref<string[]> = ref([])

  function buildFile(file: FileMessage | undefined, chunk: string | undefined) {
    if (!chunk && !file) {
      return
    }

    if (file?.status === FileStatus.Complete) {
      console.log('TRANSFER COMPLETE')
      const blob = new Blob(buf.value, { type: file.name })
      const load: Download = { from: file.from, file: new File([blob], file.name), mime: file.mime }
      downloadStore.addDownload(load)
      console.log(load)
      buf.value = []
      return
    }

    if (chunk) {
      buf.value.push(chunk)
    }
  }

  async function sendFiles(files: File[]) {
    const connections = conn.GetConnections()
    if (connections === undefined) return

    if (connections.length < 1) {
      toast.notify({ message: 'Not connected to anyone', type: ToastType.Warning })
      return
    }

    for (const file of files) {
      await sendFile(file)
    }
  }

  async function sendFile(file: File) {
    const connections = conn.GetConnections()
    if (connections === undefined) return

    const pl: FileMessage = {
      mime: file.type,
      name: file.name,
      status: FileStatus.Busy,
      from: user.getUsername()
    }
    let buf = await file.arrayBuffer()
    while (buf.byteLength) {
      const chunk = buf.slice(0, CHUNK_SIZE)
      buf = buf.slice(CHUNK_SIZE, buf.byteLength)

      for (const connection of connections) {
        connection.dc?.send(chunk)
      }
    }
    pl.status = FileStatus.Complete

    for (const connection of connections) {
      connection.dc?.send(JSON.stringify(pl))
    }
  }

  return {
    sendFile,
    sendFiles,
    buildFile
  }
})
