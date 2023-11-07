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

export const useFileStore = defineStore('file', () => {
    const conn = useConnStore()
    const user = useUserStore()
    const toast = useToastStore()
    const offeredFiles: Ref<Map<string, File[]>> = ref(new Map())
    const stream: Ref<WritableStream<Uint8Array> | undefined> = ref(undefined)
    const writer: Ref<WritableStreamDefaultWriter<Uint8Array> | undefined> = ref(undefined)
    const accSize = ref(0)
    const currentOffer: Ref<FileOffer | undefined> = ref(undefined)
    StreamSaver.mitm = `${import.meta.env.VITE_MITM_URL}/mitm.html?version=2.0.0`


    function setCurrentOffer(offer: FileOffer) {
        currentOffer.value = offer
    }

    function getCurrentOffer(): Ref<FileOffer | undefined> {
        return currentOffer
    }

    async function buildFile(chunk: ArrayBuffer) {
        if (!currentOffer.value) {
            console.warn("current file not setup")
            return
        }

        const file = currentOffer.value.files[currentOffer.value.current]

        if (!stream.value) {
            console.log("stream.value")
            stream.value = StreamSaver.createWriteStream(file.name, {size: file.size})
            writer.value = stream.value.getWriter()
            console.log(writer.value)
        }

        const buffer = new Uint8Array(chunk)
        await writer.value?.write(buffer)

        accSize.value += buffer.length

        // update file progress
        file.progress = accSize.value;
        currentOffer.value.files[currentOffer.value.current] = file;

        if (accSize.value == file.size) {
            console.log(currentOffer.value?.files)
            console.log("close writer", accSize.value)
            await writer.value?.close()

            console.log(currentOffer.value)
            const connection = conn.GetUserConnection(currentOffer.value.from)
            if (!connection) return

            file.status = FileStatus.Complete
            currentOffer.value.files[currentOffer.value.current] = file;

            stream.value = undefined
            writer.value = undefined
            accSize.value = 0

            const offer = currentOffer.value;
            console.log("RequestNext", offer.current, offer.files.length, offer.current + 1 < offer.files.length)
            if (offer.current + 1 >= offer.files.length) {
                console.log("No next file")
                return
            }
            offer.current += 1
            offer.status = FileSetup.RequestNext
            currentOffer.value = offer

            setCurrentOffer(offer)

            connection.dc?.send(JSON.stringify(offer))
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
                mime: file.type,
                name: file.name,
                size: file.size,
                progress: 0,
                status: FileStatus.Init,
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

        const fileMessages = filesToFileMessage(files);

        const offer: FileOffer = {
            id: uuidv4(),
            status: FileSetup.Offer,
            files: fileMessages,
            current: 0,
            target: "target",
            from: user.getUsername()
        }

        addOfferedFiles(offer.id, files)

        for (const connection of connections) {
            offer.target = connection.username
            connection.dc?.send(JSON.stringify(offer))
        }
    }

    function respondToFileOffer(offer: FileOffer, status: FileSetup) {
        const connection = conn.GetUserConnection(offer.from)
        if (!connection) return

        offer.status = status

        connection.dc?.send(JSON.stringify(offer))
    }

    async function sendFile(file: File, offer: FileOffer) {
        const connection = conn.GetUserConnection(offer.target)
        console.log("send File", offer)
        if (!connection) {
            toast.notify({message: `No longer connected to ${offer.target}`, type: ToastType.Warning})
            return
        }

        const stream = file.stream()
        const reader = stream.getReader()


        if (offer.current == 0) {
            offer.status = FileSetup.LatestOffer
            connection.dc?.send(JSON.stringify(offer))
        }
        const readChunk = async () => {
            const {value} = await reader.read();
            if (!value) return
            connection.dc?.send(value)
            await readChunk()
        }

        await readChunk()
    }

    return {
        setCurrentOffer,
        getCurrentOffer,
        addOfferedFiles,
        getOfferedFiles,
        removeOfferedFile,
        sendFile,
        buildFile,
        respondToFileOffer,
        sendFilesOffer,
    }
})
