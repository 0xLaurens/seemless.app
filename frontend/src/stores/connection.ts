import { defineStore } from 'pinia'
import type { Ref } from 'vue'
import { ref } from 'vue'
import type { Connection } from '@/models/connection'
import { useToastStore } from '@/stores/toast'
import { ToastType } from '@/models/toast'
import type { SessionDescriptionMessage } from '@/models/sdm'
import { useUserStore } from '@/stores/user'
import { RequestTypes } from '@/models/request'
import type { Message } from '@/models/message'

export const useConnStore = defineStore('conn', () => {
  const user = useUserStore()
  const toast = useToastStore()

  const DATACHANNEL_NAME = 'files'
  const conn: Ref<Map<string, Connection>> = ref(new Map())

  /////////////////////////////
  ///  RTC Event Listeners  ///
  /////////////////////////////
  function _setupRtcConnEventListeners(connection: Connection) {
    connection.pc.ondatachannel = (dc) => {
      connection.dc = dc.channel
      _setupDatachannelEventListeners(connection)
    }
    connection.pc.onicecandidate = (ice) => {
      toast.notify({ message: `ICE candidate ${ice.type}`, type: ToastType.Info })
    }
  }

  function _setupDatachannelEventListeners(connection: Connection) {
    if (!connection.dc) return
    connection.dc.onopen = () =>
      toast.notify({ message: 'Datachannel Opened', type: ToastType.Success })
    connection.dc.onmessage = (m) => console.log(m)
    connection.dc.onclose = () =>
      toast.notify({ message: 'Datachannel Closed :(', type: ToastType.Warning })
  }

  // ICE helper functions
  function _addIceCandidate(ice: RTCIceCandidate) {}

  function _removeIceCandidate() {}

  /////////////
  ///  RTC  ///
  /////////////
  function _setupRtcConn(username: string) {
    const connection: Connection = {
      pc: new RTCPeerConnection(),
      username: username
    }
    connection.dc = connection.pc.createDataChannel(DATACHANNEL_NAME)
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
      toast.notify({ message: 'Fatal error in connection', type: ToastType.Error })
    const connection = conn.value.get(answer.from)
    answer.type = 'answer'
    await connection?.pc.setRemoteDescription(answer)
    toast.notify({ message: 'Received: Answer!', type: ToastType.Success })
  }

  return {
    CreateRtcOffer,
    HandleRtcOffer,
    HandleRtcAnswer
  }
})
