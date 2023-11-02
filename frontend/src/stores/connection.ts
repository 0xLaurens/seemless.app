import {defineStore} from 'pinia'
import type {Ref} from 'vue'
import {ref} from 'vue'
import type {Connection} from '@/models/connection'
import {useToastStore} from '@/stores/toast'
import {ToastType} from '@/models/toast'
import type {SessionDescriptionMessage} from '@/models/sdm'
import {useUserStore} from '@/stores/user'
import {RequestTypes} from '@/models/request'
import type {Message} from '@/models/message'
import {useWebsocketStore} from '@/stores/websocket'
import {isJSON} from '@/utils/json'
import {useFileStore} from '@/stores/file'
import {FileSetup, FileStatus} from "@/models/file";
import {useDownloadStore} from "@/stores/download";

export const useConnStore = defineStore('conn', () => {
    const user = useUserStore()
    const download = useDownloadStore()
    const toast = useToastStore()
    const ws = useWebsocketStore()
    const file = useFileStore()

    const DATACHANNEL_NAME = 'files'
    const conn: Ref<Map<string, Connection>> = ref(new Map())

    /////////////////////////////
    ///  RTC Event Listeners  ///
    /////////////////////////////
    function _setupRtcConnEventListeners(connection: Connection) {
        connection.pc.ondatachannel = (dc) => {
            connection.dc = dc.channel
            connection.dc.binaryType = "arraybuffer"
            _setupDatachannelEventListeners(connection)
        }
        connection.pc.onicecandidate = (ice) => {
            const message: Message = {
                candidate: ice.candidate,
                from: user.getUsername(),
                target: connection.username,
                type: RequestTypes.NewIceCandidate
            }
            ws.SendMessage(message)
        }
    }

    function _setupDatachannelEventListeners(connection: Connection) {
        if (!connection.dc) return
        connection.dc.onmessage = async (ev) => {
            if (typeof ev.data === 'object') {
                await file.buildFile(ev.data)
                return
            }

            if (!isJSON(ev.data)) return
            const message = JSON.parse(ev.data)
            if (message.status == FileStatus.Complete) {
                console.log(FileStatus.Complete)
            }

            if (message.status == FileSetup.Offer) {
                console.log(message)
                download.addOffer(message)
            }
        }
        connection.dc.onclose = () => {
            toast.notify({
                message: `Connection lost to ${connection.username}`,
                type: ToastType.Warning
            })
            conn.value.delete(connection.username)
        }
        connection.dc.onerror = (err) => console.error(err)
    }

    // ICE helper functions
    async function _addIceCandidate(
        ice: RTCIceCandidateInit | undefined | null,
        from: string | undefined
    ) {
        if (ice === null) return
        if (from === undefined) return
        const connection = conn.value.get(from)
        if (connection === undefined) return

        await connection.pc.addIceCandidate(ice)
    }

    /////////////
    ///  RTC  ///
    /////////////
    function _setupRtcConn(username: string) {
        const connection: Connection = {
            pc: new RTCPeerConnection(),
            username: username
        }
        _setupRtcConnEventListeners(connection)
        conn.value.set(username, connection)
    }

    async function CreateRtcOffer(username: string): Promise<Message | undefined> {
        if (!conn.value.has(username)) _setupRtcConn(username)
        const connection = conn.value.get(username)
        if (connection == undefined) {
            toast.notify({
                message: 'Something went wrong creating the connection',
                type: ToastType.Warning
            })
            return undefined
        }

        connection.dc = connection.pc.createDataChannel(DATACHANNEL_NAME)
        connection.dc.binaryType = "arraybuffer"
        _setupDatachannelEventListeners(connection)
        _setupRtcConnEventListeners(connection)

        const offer = await connection.pc.createOffer()
        await connection.pc.setLocalDescription(offer)
        return {
            from: user.getUsername(),
            sdp: offer.sdp,
            target: username,
            type: RequestTypes.Offer
        }
    }

    async function HandleRtcOffer(offer: SessionDescriptionMessage): Promise<Message | undefined> {
        _setupRtcConn(offer.from)
        const connection = conn.value.get(offer.from)
        if (connection == undefined) {
            toast.notify({
                message: 'Something went wrong creating the connection',
                type: ToastType.Warning
            })
            return undefined
        }
        offer.type = 'offer'
        await connection.pc.setRemoteDescription(offer).catch(console.error)
        const answer = await connection.pc.createAnswer()
        await connection.pc.setLocalDescription(answer)
        return {
            from: user.getUsername(),
            sdp: answer.sdp,
            target: offer.from,
            type: RequestTypes.Answer
        }
    }

    async function HandleRtcAnswer(answer: SessionDescriptionMessage) {
        if (!conn.value.has(answer.from))
            toast.notify({message: 'Fatal error in connection', type: ToastType.Error})
        const connection = conn.value.get(answer.from)
        answer.type = 'answer'
        await connection?.pc.setRemoteDescription(answer)
        toast.notify({message: 'Received: Answer!', type: ToastType.Success})
    }

    async function HandleIceCandidate(message: Message | undefined) {
        if (message === undefined) return

        await _addIceCandidate(message.candidate, message.from)
    }

    function GetConnections() {
        return Array.from(conn.value.values())
    }

    function GetUserConnection(target: string) {
        return conn.value.get(target)
    }

    return {
        CreateRtcOffer,
        HandleRtcOffer,
        HandleRtcAnswer,
        HandleIceCandidate,
        GetConnections,
        conn,
        GetUserConnection
    }
})
