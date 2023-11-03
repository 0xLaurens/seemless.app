import type {FileMessage, FileOffer} from '@/models/file'
import {FileSetup, FileStatus} from '@/models/file'
import {v4 as uuidv4} from 'uuid';
import type {Ref} from 'vue'
import {ref} from 'vue'
import {defineStore} from 'pinia'
import {useConnStore} from '@/stores/connection'
import {useUserStore} from '@/stores/user'
import {useToastStore} from '@/stores/toast'
import {ToastType} from '@/models/toast'
import StreamSaver from "streamsaver";
import type {Connection} from "@/models/connection";

export const useFileStore = defineStore('file', () => {
    const conn = useConnStore()
    const user = useUserStore()
    const toast = useToastStore()
    const offeredFiles: Ref<Map<string, File[]>> = ref(new Map())
    const stream: Ref<WritableStream<Uint8Array> | undefined> = ref(undefined)
    const writer: Ref<WritableStreamDefaultWriter<Uint8Array> | undefined> = ref(undefined)
    const accSize = ref(0)
    const currentFile: Ref<FileMessage | undefined> = ref(undefined)
    StreamSaver.mitm = `${import.meta.env.VITE_MITM_URL}/mitm.html?version=2.0.0`


    function setCurrentFile(file: FileMessage) {
        currentFile.value = file
    }

    async function buildFile(chunk: ArrayBuffer) {
        if (!currentFile.value) {
            console.warn("current file not setup")
            return
        }
        if (!stream.value) {
            console.log("stream.value")
            stream.value = StreamSaver.createWriteStream(currentFile.value.name, {size: currentFile.value.size})
            writer.value = stream.value.getWriter()
            console.log(writer.value)
        }

        const buffer = new Uint8Array(chunk)
        await writer.value?.write(buffer)

        accSize.value += buffer.length

        if (accSize.value == currentFile.value.size) {
            console.log("close writer", accSize.value)
            await writer.value?.close()
            stream.value = undefined
            writer.value = undefined
            accSize.value = 0
        }
    }

    function addOfferedFiles(id: string, files: File[]) {
        offeredFiles.value.set(id, files)
    }

    function getOfferedFiles(id: string): File[] | undefined {
        return offeredFiles.value.get(id)
    }

    function removeOfferedFile(id: string) {
        offeredFiles.value.delete(id)
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
            id: uuidv4(),
            status: FileSetup.Offer,
            files: fileMessages,
            from: user.getUsername()
        }

        addOfferedFiles(offer.id, files)

        for (const connection of connections) {
            connection.dc?.send(JSON.stringify(offer))
        }
    }

    function respondToFileOffer(offer: FileOffer, status: FileSetup) {
        const connection = conn.GetUserConnection(offer.from)
        if (!connection) return

        offer.status = status
        offer.from = user.getUsername()

        connection.dc?.send(JSON.stringify(offer))
    }

    // async function sendFiles(files: File[], username: string) {
    //     const connection = conn.GetUserConnection(username)
    //     if (!connection) {
    //         toast.notify({message: `No longer connected to ${username}`, type: ToastType.Warning})
    //         return
    //     }
    //
    //     for (const file of files) {
    //         await sendFile(file, connection)
    //     }
    // }

    async function sendFile(file: File, username: string) {
        const connection = conn.GetUserConnection(username)
        if (!connection) {
            toast.notify({message: `No longer connected to ${username}`, type: ToastType.Warning})
            return
        }
        const pl: FileMessage = {
            status: FileStatus.Init,
            mime: file.type,
            name: file.name,
            size: file.size,
        }
        const stream = file.stream()
        const reader = stream.getReader()

        connection.dc?.send(JSON.stringify(pl))
        const readChunk = async () => {
            const {done, value} = await reader.read();

            if (done) {
                pl.status = FileStatus.Complete
                connection.dc?.send(JSON.stringify(pl))
                return;
            }

            connection.dc?.send(value)
            await readChunk()
        }

        await readChunk()
    }

    return {
        setCurrentFile,
        addOfferedFiles,
        getOfferedFiles,
        removeOfferedFile,
        sendFile,
        buildFile,
        respondToFileOffer,
        sendFilesOffer,
    }
})
