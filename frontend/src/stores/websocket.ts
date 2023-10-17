import { defineStore } from 'pinia'
import type { Ref } from 'vue'
import { ref } from 'vue'
import { RequestTypes } from '@/models/request'
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast'
import type { Message } from '@/models/message'
import { useConnStore } from '@/stores/connection'
import { ToastType } from '@/models/toast'
import router from '@/router'

export const useWebsocketStore = defineStore('ws', () => {
  const user = useUserStore()
  const conn = useConnStore()
  const toast = useToastStore()

  const ws: Ref<WebSocket | undefined> = ref()

  function _setupWsListeners() {
    if (ws.value === undefined) {
      return
    }
    ws.value.onmessage = async (event) => _onMessage(event)
    ws.value.onclose = () => _onClose()
    ws.value.onopen = () => _onOpen()
  }

  async function _onMessage(event: MessageEvent<any>) {
    const data = JSON.parse(event.data)
    switch (data.type) {
      case RequestTypes.Offer: {
        const answer = await conn.HandleRtcOffer(data)
        SendMessage(answer)
        break
      }
      case RequestTypes.NewIceCandidate: {
        // if (data.from === user.getUsername()) return
        // if (!connected.userInIceTargets(data.from)) return
        // await rtc.handleIceCandidate(JSON.parse(data.candidate))
        break
      }
      case RequestTypes.Answer: {
        await conn.HandleRtcAnswer(data)
        break
      }
      case RequestTypes.Peers:
        user.initUsers(data.users)
        break
      case RequestTypes.PeerJoined:
        user.addUser(data.user)
        break
      case RequestTypes.PeerLeft:
        user.removeUser(data.user)
        break
      case RequestTypes.UsernamePrompt:
      case RequestTypes.Username:
        console.log(data.message)
        break
      case RequestTypes.DuplicateUsername:
        toast.notify({
          message: 'This username has already been taken in this room',
          type: ToastType.Warning
        })
        await router.push({ path: '/nick' })
        break
      default:
        console.log(`Unknown type ${data.type}`)
    }
  }

  function _onOpen() {
    const payload = {
      type: 'Username',
      body: {
        username: user.getUsername()
      }
    }

    ws.value?.send(JSON.stringify(payload))
  }

  function _onClose() {}

  function SendMessage(msg: Message | undefined) {
    if (ws.value === undefined) {
      toast.notify({ message: 'The connection is currently closed', type: ToastType.Warning })
      return
    }
    if (msg === undefined) {
      toast.notify({ message: 'Failed to send message', type: ToastType.Warning })
      return
    }

    ws.value?.send(JSON.stringify(msg))
  }

  function Open() {
    if (ws.value !== undefined) return
    ws.value = new WebSocket(import.meta.env.VITE_WS_ADDR)
    _setupWsListeners()
  }

  function Close() {
    ws.value?.close()
  }

  function GetConnection() {
    return ws.value
  }

  return {
    Open,
    Close,
    GetConnection,
    SendMessage
  }
})
