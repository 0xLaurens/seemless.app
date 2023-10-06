import { defineStore } from 'pinia'
import type { SessionDescriptionMessage } from '@/models/sdm'
import { ref } from 'vue'
import type { FileMessage } from '@/models/file'
import { FileStatus } from '@/models/file'
import { useUserStore } from '@/stores/user'
import { DcStatus } from '@/models/datachannel'
import { useConnectedStore } from '@/stores/connected'

export const useRtcStore = defineStore('rtc', () => {
  const CHUNK_SIZE = 65536 //64 KiB
  const blobURL = ref('')
  const test_buf: any[] = []
  const rtc = new RTCPeerConnection()
  const user = useUserStore()
  const conn = useConnectedStore()
  let dc = createDatachannel('files')
  let localFragment: string | null

  rtc.ondatachannel = (dc) => setDatachannel(dc.channel)

  function getPeerConnection() {
    return rtc
  }

  function close() {
    if (dc.readyState === 'open') {
      dc.send(JSON.stringify({ username: user.getUsername(), status: DcStatus.ClientClose }))
    }

    rtc.close()
  }

  function createDatachannel(name: string): RTCDataChannel {
    const dataChannel = rtc.createDataChannel(name)
    addDcListeners(dataChannel)
    return dataChannel
  }

  async function addDcListeners(dc: RTCDataChannel) {
    dc.onopen = () =>
      dc.send(JSON.stringify({ username: user.getUsername(), status: DcStatus.ClientHello }))
    dc.onmessage = (ev) => {
      if (typeof ev.data === 'object') {
        recvFile(undefined, ev.data)
        return
      }

      if (!isJSON(ev.data)) return

      const fileOrDcMessage = JSON.parse(ev.data)

      if (fileOrDcMessage.username && fileOrDcMessage.status) {
        const connected = fileOrDcMessage.status === DcStatus.ClientHello
        conn.createUserConnection(fileOrDcMessage.username, connected)
        return
      }

      recvFile(fileOrDcMessage, undefined)
    }
    dc.onclose = (ev) => {
      dc.send(JSON.stringify({ username: user.getUsername(), status: DcStatus.ClientClose }))
    }
  }

  function isJSON(str: string) {
    try {
      JSON.parse(str)
    } catch (e) {
      return false
    }

    return true
  }

  function setDatachannel(datachannel: RTCDataChannel) {
    dc = datachannel
    addDcListeners(dc)
  }

  function getDatachannel(): RTCDataChannel {
    return dc
  }

  async function sendFiles(files: File[]) {
    if (dc.readyState !== 'open') {
      console.log("COULDN'T SEND FILES DATACHANNEL IS CLOSED")
      return
    }

    for (const file of files) {
      await sendFile(file)
    }
  }

  async function recvFile(file: FileMessage | undefined, chunk: string | undefined) {
    if (!chunk && !file) {
      return
    }

    if (file?.status === FileStatus.Complete) {
      console.log('TRANSFER COMPLETE')
      const blob = new Blob(test_buf, { type: file.name })
      const blobFile = new File([blob], file.name)
      blobURL.value = URL.createObjectURL(blobFile)
      return
    }

    if (chunk) {
      test_buf.push(chunk)
    }
  }

  async function sendFile(file: File) {
    const pl: FileMessage = { mime: file.type, name: file.name, status: FileStatus.Busy }
    let buf = await file.arrayBuffer()
    while (buf.byteLength) {
      const chunk = buf.slice(0, CHUNK_SIZE)
      buf = buf.slice(CHUNK_SIZE, buf.byteLength)
      dc.send(chunk)
    }
    pl.status = FileStatus.Complete
    dc.send(JSON.stringify(pl))
  }

  function setLocalFragment(fragment: string | null) {
    localFragment = fragment
  }

  function getLocalFragment(): string | null {
    return localFragment
  }

  async function createOffer(): Promise<RTCSessionDescriptionInit> {
    console.log('CREATE OFFER')
    const offer = await rtc.createOffer()
    await rtc.setLocalDescription(offer)
    return offer
  }

  async function createOfferAnswer(offer: SessionDescriptionMessage) {
    console.log('HANDLE OFFER ANSWER')
    offer.type = 'offer'
    await rtc.setRemoteDescription(offer).catch(console.error)
    const answer = await rtc.createAnswer()
    await rtc.setLocalDescription(answer)
    return answer
    // sendMessage('Answer', offer.from, answer.sdp)
  }

  async function handleAnswer(answer: SessionDescriptionMessage) {
    console.log('HANDLE ANSWER')
    answer.type = 'answer'
    if (rtc.signalingState === 'stable') {
      return
    }

    await rtc.setRemoteDescription(answer).catch(console.error)
  }

  async function handleIceCandidate(candidate: RTCIceCandidateInit | null) {
    console.log('HANDLE ICE CANDIDATE', getLocalFragment())
    if (getLocalFragment() === null) return

    if (candidate == null || candidate.usernameFragment === getLocalFragment()) {
      return
    }
    await rtc.addIceCandidate(candidate)
  }

  return {
    blobURL,
    close,
    getPeerConnection,
    setDatachannel,
    setLocalFragment,
    createOffer,
    createOfferAnswer,
    handleAnswer,
    handleIceCandidate,
    sendFiles
  }
})
