import {defineStore} from 'pinia'
import type {Ref} from 'vue'
import {ref} from 'vue'
import type {Download} from '@/models/download'
import type {FileOffer} from "@/models/file";
import {FileSetup} from "@/models/file";
import {useFileStore} from "@/stores/file";

export const useDownloadStore = defineStore('download', () => {
    const downloads: Ref<Map<string, Download>> = ref(new Map())
    const popup: Ref<boolean> = ref(false)
    const offers: Ref<FileOffer[]> = ref([])
    const activeOffer: Ref<FileOffer | undefined> = ref(undefined)
    const url: Ref<string | undefined> = ref()
    const activeDownload: Ref<Download | undefined> = ref()

    const file = useFileStore()

    function open() {
        popup.value = true
    }

    function close() {
        popup.value = false
    }

    function addOffer(offer: FileOffer) {
        offers.value.push(offer)
        setLatestOffer()
        open()
    }

    function acceptOffer(offer: FileOffer | undefined) {
        close()
        if (!offer) return
        removeOffer(offer)
        file.respondToFileOffer(offer, FileSetup.Accept)
    }

    function denyOffer(offer: FileOffer | undefined) {
        close()
        if (!offer) return
        removeOffer(offer)
        file.respondToFileOffer(offer, FileSetup.Deny)
    }

    function removeOffer(offer: FileOffer | undefined) {
        offers.value.filter((o) => o === offer)
        setLatestOffer()
    }

    function setLatestOffer() {
        activeOffer.value = offers.value.length > 1
            ? undefined
            : offers.value[0]
        console.log(activeOffer.value)
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
        removeDownload,
        activeOffer,
        addOffer,
        denyOffer,
        acceptOffer,
    }
})
