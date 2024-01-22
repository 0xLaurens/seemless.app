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
import {FileSetup} from "@/models/file";
import {useDownloadStore} from "@/stores/download";

export const useConnStore = defineStore('conn', () => {
    const user = useUserStore()
    const download = useDownloadStore()
    const toast = useToastStore()
    const ws = useWebsocketStore()
    const file = useFileStore()

    const makingOffer = ref(false)

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
                target: connection.userId,
                type: RequestTypes.NewIceCandidate
            }
            ws.SendMessage(message)
        }
        connection.pc.onnegotiationneeded = async () => {
            try {
                makingOffer.value = true;
                const offer = await CreateRtcOffer(connection.userId)
                ws.SendMessage(offer)
            } catch (err) {
                console.error(err);
            } finally {
                makingOffer.value = false;
            }
        };
        connection.pc.oniceconnectionstatechange = () => {
            if (connection.pc.iceConnectionState == "failed") {
                connection.pc.restartIce()
            }
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
            const offer = JSON.parse(ev.data)
            switch (offer.status) {
                case FileSetup.RequestNext: {
                    const files = file.getOfferedFiles(offer.id)
                    if (files == undefined) return

                    if (offer.current < files.length) {
                        await file.sendFile(files[offer.current], offer)
                    }
                    break
                }
                case FileSetup.LatestOffer: {
                    file.setReceiveOffer(offer)
                    break
                }
                case FileSetup.DownloadProgress: {
                    file.setFileProgress(offer)
                    break;
                }
                case FileSetup.Offer: {
                    download.addOffer(offer)
                    break
                }
                case FileSetup.AcceptOffer: {
                    file.setSendOffer(offer)
                    file.setFileProgress(offer)
                    const files = file.getOfferedFiles(offer.id)
                    if (files == undefined) {
                        console.log("something went wrong!")
                        return
                    }

                    await file.sendFile(files[0], offer)
                    break
                }
                case FileSetup.DenyOffer: {
                    console.log("offer was denied :(")
                }


            }
        }
        connection.dc.onclose = () => {
            // toast.notify({
            //     message: `Connection lost to ${connection.username}`,
            //     type: ToastType.Warning
            // })
            const userToRemove = user.getUserByUsername(connection.userId)
            if (userToRemove) {
                user.removeUser(userToRemove)
            }
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

    function RemoveRtcConn(username: string) {
        const connection = conn.value.get(username)
        if (!connection) return

        if (connection.dc) connection.dc.close()
        connection.pc.close()

        conn.value.delete(username)
    }

    /////////////
    ///  RTC  ///
    /////////////
    function _setupRtcConn(userId: string) {
        const connection: Connection = {
            pc: new RTCPeerConnection({iceServers: [{urls: 'stun:stun.l.google.com:19302'}]}),
            userId: userId
        }
        _setupRtcConnEventListeners(connection)
        conn.value.set(userId, connection)
    }


    async function CreateRtcOffer(userId: string): Promise<Message | undefined> {
        if (!conn.value.has(userId)) _setupRtcConn(userId)

        const connection = conn.value.get(userId)
        if (connection == undefined) {
            toast.notify({
                message: 'Something went wrong creating the connection',
                type: ToastType.Warning
            })
            return undefined
        }

        const target = user.getUserByUsername(userId)
        connection.dc = connection.pc.createDataChannel(target!.id)
        connection.dc.binaryType = "arraybuffer"
        _setupDatachannelEventListeners(connection)
        _setupRtcConnEventListeners(connection)

        const offer = await connection.pc.createOffer()
        await connection.pc.setLocalDescription(offer)
        return {
            from: user.getUsername(),
            sdp: offer.sdp,
            target: userId,
            type: RequestTypes.Offer
        }
    }

    async function HandleRtcOffer(offer: SessionDescriptionMessage): Promise<Message | undefined> {
        _setupRtcConn(offer.from)
        console.log("handling offer", offer)
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
        RemoveRtcConn,
        CreateRtcOffer,
        HandleRtcOffer,
        HandleRtcAnswer,
        HandleIceCandidate,
        GetConnections,
        conn,
        GetUserConnection
    }
})
