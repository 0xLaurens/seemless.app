import {defineStore} from 'pinia'
import type {Ref} from 'vue'
import {ref} from 'vue'
import {RequestTypes} from '@/models/request'
import {useUserStore} from '@/stores/user'
import {useToastStore} from '@/stores/toast'
import type {Message} from '@/models/message'
import {useConnStore} from '@/stores/connection'
import {ToastType} from '@/models/toast'
import router from '@/router'
import {useRetryStore} from '@/stores/retry'
import {useRoomStore} from "@/stores/room";

export const useWebsocketStore = defineStore('ws', () => {
    const user = useUserStore()
    const conn = useConnStore()
    const toast = useToastStore()
    const retry = useRetryStore()
    const room = useRoomStore()

    const ws: Ref<WebSocket | undefined> = ref()
    const intended_close: Ref<Boolean> = ref(false)

    function _setupWsListeners() {
        if (ws.value === undefined) {
            return
        }
        ws.value.onmessage = async (event) => _onMessage(event)
        ws.value.onclose = () => _onClose()
        ws.value.onopen = () => _onOpen()
        ws.value.onerror = () => _onError()
    }

    async function _onMessage(event: MessageEvent<any>) {
        const data = JSON.parse(event.data)
        console.log("WS:", data)
        switch (data.type) {
            case RequestTypes.Offer: {
                const answer = await conn.HandleRtcOffer(data)
                SendMessage(answer)
                break
            }
            case RequestTypes.NewIceCandidate: {
                await conn.HandleIceCandidate(data)
                break
            }
            case RequestTypes.Answer: {
                await conn.HandleRtcAnswer(data)
                break
            }
            case RequestTypes.Peers:
                if (!data.users) return
                for (const newUser of data.users) {
                    console.log(newUser.username, user.getUsername())
                    user.addUser(newUser)
                }
                break
            case RequestTypes.PeerJoined: {
                if (data.user.username === user.getUsername()) return
                user.addUser(data.user)
                const offer = await conn.CreateRtcOffer(data.user.username)
                if (offer) {
                    SendMessage(offer)
                }
                break;
            }
            case RequestTypes.PeerLeft:
                user.removeUser(data.user)
                break
            case RequestTypes.DuplicateUsername:
                toast.notify({
                    message: 'This username has already been taken in this room',
                    type: ToastType.Warning
                })
                await router.push({path: '/nick'})
                break
            case RequestTypes.DisplayName:
                user.setCurrentUser(data.user)
                user.setUsername(data.user.username)
                break
            case RequestTypes.RoomCreated:
                room.setRoomCode(data.roomCode)
                break
            case RequestTypes.RoomCodeInvalid:
                toast.notify({message: 'This room does not exist', type: ToastType.Warning})
                break
            case RequestTypes.RoomJoined:
                room.setRoomCode(data.roomCode)
                break
            default:
                console.log(`Unknown type ${data.type}`)
        }
    }

    function _onOpen() {
        retry.stop()
        toast.notify({message: 'Connected to room', type: ToastType.Success})
        // const payload = {
        //     type: 'Username',
        //     body: {
        //         username: user.getUsername()
        //     }
        // }
        //
        // ws.value?.send(JSON.stringify(payload))
    }

    function _onClose() {
        if (retry.isActive()) return
        toast.notify({message: 'Disconnected from room', type: ToastType.Warning})
        if (intended_close.value) return
        user.clearUsers()
        retry.start(() => {
            _reconnectWs()
        })
    }

    function _onError() {
        if (retry.isActive()) return
    }

    function _reconnectWs() {
        console.log('reconnecting')
        if (ws.value?.readyState !== ws.value?.CLOSED) {
            Close()
        }
        ws.value = undefined
        Open()
    }

    function SendMessage(msg: Message | undefined) {
        if (ws.value === undefined) {
            toast.notify({message: 'The connection is currently closed', type: ToastType.Warning})
            return
        }
        if (msg === undefined) {
            toast.notify({message: 'Failed to send message', type: ToastType.Warning})
            return
        }

        ws.value?.send(JSON.stringify(msg))
    }

    function isOpen(): boolean {
        return ws.value?.readyState === ws.value?.OPEN
    }

    function isClosed(): boolean {
        return ws.value?.readyState === ws.value?.CLOSED
    }

    function Open() {
        intended_close.value = false
        if (ws.value !== undefined && ws.value?.readyState == ws.value?.OPEN) return
        const WEBSOCKET_URL = import.meta.env.VITE_WS_ADDR ? import.meta.env.VITE_WS_ADDR : "wss://lightsail-container-service.mcs8tl3k5ku3q.eu-central-1.cs.amazonlightsail.com/ws"
        ws.value = new WebSocket(WEBSOCKET_URL)
        _setupWsListeners()
    }

    function Close() {
        intended_close.value = true
        ws.value?.close()
    }

    function GetConnection() {
        return ws.value
    }

    return {
        Open,
        Close,
        isOpen,
        isClosed,
        GetConnection,
        SendMessage
    }
})
