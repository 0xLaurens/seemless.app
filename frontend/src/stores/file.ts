import type {FileMessage} from '@/models/file'
import {FileStatus} from '@/models/file'
import type {Download} from '@/models/download'
import type {Ref} from 'vue'
import {ref} from 'vue'
import {defineStore} from 'pinia'
import {useDownloadStore} from '@/stores/download'
import {useConnStore} from '@/stores/connection'
import {useUserStore} from '@/stores/user'
import {useToastStore} from '@/stores/toast'
import {ToastType} from '@/models/toast'
import StreamSaver from "streamsaver";

export const useFileStore = defineStore('file', () => {
    const downloadStore = useDownloadStore()
    const conn = useConnStore()
    const user = useUserStore()
    const toast = useToastStore()
    const stream: Ref<WritableStream<Uint8Array> | undefined> = ref(undefined)
    const writer: Ref<WritableStreamDefaultWriter<Uint8Array> | undefined> = ref(undefined)
    const accSize = ref(0)
    const filesize = 51346
    console.log(StreamSaver.mitm)
    StreamSaver.mitm = `${import.meta.env.VITE_MITM_URL}/mitm.html?version=2.0.0`

    async function buildFile(file: FileMessage | undefined, chunk: ArrayBuffer | undefined) {
        if (!chunk && !file) {
            return;
        }
        console.log(accSize.value)

        if (!stream.value) {
            console.log("stream.value")
            stream.value = StreamSaver.createWriteStream("test3.jpg", {size: filesize})
            writer.value = stream.value.getWriter()
            console.log(writer.value)
        }

        if (chunk == undefined) return
        console.log(chunk)
        const buffer = new Uint8Array(chunk)
        console.log(buffer)
        await writer.value?.write(buffer)

        accSize.value += buffer.length

        if (accSize.value == filesize) {
            console.log("close writer", accSize.value)
            await writer.value?.close()
            stream.value = undefined
            writer.value = undefined
            accSize.value = 0
        }
    }

    async function sendFiles(files: File[]) {
        const connections = conn.GetConnections()
        if (connections === undefined) return

        if (connections.length < 1) {
            toast.notify({message: 'Not connected to anyone', type: ToastType.Warning})
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
            size: file.size,
            status: FileStatus.Busy,
            from: user.getUsername()
        }
        const stream = file.stream()
        const reader = stream.getReader()

        const readChunk = async () => {
            const {done, value} = await reader.read();

            if (done) {
                pl.status = FileStatus.Complete
                for (const connection of connections) {
                    connection.dc?.send(JSON.stringify(pl))
                }
                return;
            }

            for (const connection of connections) {
                console.log(value)
                connection.dc?.send(value)
            }

            await readChunk()
        }

        await readChunk()
    }

    return {
        sendFile,
        sendFiles,
        buildFile
    }
})
