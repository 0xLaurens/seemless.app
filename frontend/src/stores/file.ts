import type {FileMessage, FileOffer} from '@/models/file'
import {FileSetup, FileStatus} from '@/models/file'
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
    StreamSaver.mitm = `${import.meta.env.VITE_MITM_URL}/mitm.html?version=2.0.0`

    async function buildFile(chunk: ArrayBuffer) {
        if (chunk) {
            return;
        }

        if (!stream.value) {
            console.log("stream.value")
            stream.value = StreamSaver.createWriteStream("test3.jpg", {size: filesize})
            writer.value = stream.value.getWriter()
            console.log(writer.value)
        }

        const buffer = new Uint8Array(chunk)
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

    function filesToFileMessage(files: File[]): FileMessage[] {
        const messages: FileMessage[] = []
        for (const file of files) {
            const msg: FileMessage = {
                mime: file.type, name: file.name, size: file.size, status: FileStatus.Init
            }
            messages.push(msg)
        }
        return messages
    }

    async function sendFilesOffer(files: File[]) {
        const connections = conn.GetConnections()
        if (connections === undefined) return
        if (connections.length < 1) {
            toast.notify({message: 'Not connected to anyone', type: ToastType.Warning})
            return
        }

        const fileMessages = filesToFileMessage(files)

        const offer: FileOffer = {
            status: FileSetup.Offer,
            files: fileMessages,
            from: user.getUsername()
        }

        for (const connection of connections) {
            connection.dc?.send(JSON.stringify(offer))
        }
    }

    function respondToFileOffer(offer: FileOffer, status: FileSetup) {
        const connection = conn.GetUserConnection(offer.from)
        if (!connection) return

        offer.status = status
        connection.dc?.send(JSON.stringify(offer))
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
            status: FileStatus.Busy,
            mime: file.type,
            name: file.name,
            size: file.size,
        }
        const stream = file.stream()
        const reader = stream.getReader()

        const readChunk = async () => {
            const {done, value} = await reader.read();

            if (done) {
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
        buildFile,
        respondToFileOffer,
        sendFilesOffer,
    }
})
