import type {FileMessage, FileOffer} from '@/models/file'
import {FileSetup, FileStatus} from '@/models/file'
import {v4 as uuidv4} from 'uuid';
import type {Ref} from 'vue'
import {ref} from 'vue'
import {defineStore} from 'pinia'
import {useConnStore} from '@/stores/connection'
import {useUserStore} from '@/stores/user'
import StreamSaver from "streamsaver";

export const useFileStore = defineStore('file', () => {
    const conn = useConnStore()
    const user = useUserStore()
    const offeredFiles: Ref<Map<string, FileList>> = ref(new Map())
    const stream: Ref<WritableStream<Uint8Array> | undefined> = ref(undefined)
    const writer: Ref<WritableStreamDefaultWriter<Uint8Array> | undefined> = ref(undefined)
    const accSize = ref(0)
    const receiveOffer: Ref<FileOffer | undefined> = ref(undefined)
    const sendOffer: Ref<FileOffer | undefined> = ref(undefined)
    StreamSaver.mitm = import.meta.env.VITE_MITM_URL ? `${import.meta.env.VITE_MITM_URL}/mitm.html?version=2.0.0` : `https://jimmywarting.github.io/StreamSaver.js/mitm.html?version=2.0.0`
    console.log(StreamSaver.mitm)

    const fileProgress: Ref<Map<string, FileOffer>> = ref(new Map())


    function setSendOffer(offer: FileOffer) {
        sendOffer.value = offer;
    }

    function getSendOffer() {
        return sendOffer
    }

    function setFileProgress(offer: FileOffer) {
        fileProgress.value.set(offer.target, offer)
    }

    function getFileProgress(): Ref<Map<string, FileOffer>> {
        return fileProgress
    }

    function setReceiveOffer(offer: FileOffer) {
        receiveOffer.value = offer
    }

    function getReceiveOffer(): Ref<FileOffer | undefined> {
        return receiveOffer
    }

    async function buildFile(chunk: ArrayBuffer) {
        if (!receiveOffer.value) {
            console.warn("current file not setup")
            return
        }

        const file = receiveOffer.value.files[receiveOffer.value.current]
        const connection = conn.GetUserConnection(receiveOffer.value.from)

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
        receiveOffer.value.files[receiveOffer.value.current] = file;
        receiveOffer.value.status = FileSetup.DownloadProgress

        connection?.dc?.send(JSON.stringify(receiveOffer.value))

        if (accSize.value == file.size) {
            console.log(receiveOffer.value?.files)
            console.log("close writer", accSize.value)
            await writer.value?.close()

            console.log(receiveOffer.value)

            file.status = FileStatus.Complete
            receiveOffer.value.files[receiveOffer.value.current] = file;

            stream.value = undefined
            writer.value = undefined
            accSize.value = 0

            const offer = receiveOffer.value;
            console.log("RequestNext", offer.current, offer.files.length, offer.current + 1 < offer.files.length)
            if (offer.current + 1 >= offer.files.length) {
                console.log("No next file")
                return
            }
            offer.current += 1
            offer.status = FileSetup.RequestNext
            receiveOffer.value = offer

            setReceiveOffer(offer)

            connection?.dc?.send(JSON.stringify(offer))
        }
    }

    function addOfferedFiles(id: string, files: FileList) {
        offeredFiles.value.set(id, files)
    }

    function getOfferedFiles(id: string): FileList | undefined {
        return offeredFiles.value.get(id)
    }

    function removeOfferedFile(id: string) {
        offeredFiles.value.delete(id)
    }

    function filesToFileMessage(files: FileList): FileMessage[] {
        const messages: FileMessage[] = []
        for (const file of files) {
            console.log(file)
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

    async function sendFilesOffer(files: FileList | undefined, targets: string[]) {
        if (!files) return

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

        for (const target of targets) {
            const connection = conn.GetUserConnection(target)
            if (!connection) continue
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
        if (!connection || !connection.dc) {
            return
        }

        const stream = file.stream()
        const reader = stream.getReader()


        if (offer.current == 0) {
            offer.status = FileSetup.LatestOffer
            connection.dc?.send(JSON.stringify(offer))
        }

        const CHUNK_SIZE = 16384 //16 KiB chrome being chrome

        connection.dc.bufferedAmountLowThreshold = 1048576 // 1mb

        const readChunk = async () => {
            const {value} = await reader.read();
            if (!value) return
            let buf = value
            while (buf.byteLength) {

                while (connection.dc!.bufferedAmount > connection.dc!.bufferedAmountLowThreshold) {
                    await new Promise(resolve => {
                        connection.dc!.onbufferedamountlow = resolve;
                    });
                }
                const chunk = buf.slice(0, CHUNK_SIZE);
                buf = buf.slice(CHUNK_SIZE)
                connection.dc?.send(chunk);
            }

            await readChunk()
        }
        await readChunk()
    }

    return {
        setFileProgress,
        getFileProgress,

        setSendOffer,
        getSendOffer,
        setReceiveOffer,
        getReceiveOffer,
        addOfferedFiles,
        getOfferedFiles,
        removeOfferedFile,
        sendFile,
        buildFile,
        respondToFileOffer,
        sendFilesOffer,
    }
})
