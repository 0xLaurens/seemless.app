import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { SessionDescriptionMessage } from '@/models/sdm'

export const useRtcStore = defineStore('rtc', () => {
  let rtc: RTCPeerConnection
  const wsLocal = ref()
  const localFragment = ref()

  function setup(ws: WebSocket) {
    wsLocal.value = ws
    rtc = new RTCPeerConnection()
  }

  function getPeerConnection() {
    return rtc
  }

  function getLocalFragment() {
    return localFragment.value
  }

  function setLocalFragment(fragment: string | null) {
    localFragment.value = fragment
  }

  async function createOffer(): Promise<RTCSessionDescriptionInit> {
    console.log('CREATE OFFER')
    const offer = await rtc.createOffer()
    await rtc.setLocalDescription(offer)
    return offer
  }

  async function createOfferAnswer(
    offer: SessionDescriptionMessage
  ): Promise<RTCSessionDescriptionInit> {
    console.log('HANDLE OFFER ANSWER')
    await rtc.setRemoteDescription(offer)
    const answer = await rtc.createAnswer()
    await rtc.setLocalDescription(answer)
    return answer
  }

  async function handleAnswer(answer: SessionDescriptionMessage) {
    console.log('HANDLE ANSWER', rtc.connectionState)
    answer.type = 'answer'
    if (rtc.signalingState === 'stable') {
      return
    }

    await rtc.setRemoteDescription(answer)
  }

  async function handleIceCandidate(data: { type: string; candidate: string }) {
    const candidate = JSON.parse(data.candidate)
    await rtc.addIceCandidate(candidate)
  }

  function closePeerConnection() {
    rtc.close()
  }

  return {
    getPeerConnection,
    setup,
    closePeerConnection,
    createOffer,
    createOfferAnswer,
    handleIceCandidate,
    handleAnswer,
    setLocalFragment,
    getLocalFragment
  }
})
