import { defineStore } from 'pinia'
import type { SessionDescriptionMessage } from '@/models/sdm'

export const useRtcStore = defineStore('rtc', () => {
  const rtc = new RTCPeerConnection()
  let dc = createDatachannel('test')
  let localFragment: string | null

  rtc.ondatachannel = (dc) => setDatachannel(dc.channel)

  function getPeerConnection() {
    return rtc
  }

  function createDatachannel(name: string): RTCDataChannel {
    const dataChannel = rtc.createDataChannel(name)
    addDcListeners(dataChannel)
    return dataChannel
  }

  function addDcListeners(dc: RTCDataChannel) {
    dc.onopen = () => dc.send('Great Success!')
    dc.onmessage = (ev) => console.log(ev.data)
    dc.onclose = (ev) => console.log(ev)
  }

  function setDatachannel(datachannel: RTCDataChannel) {
    dc = datachannel
    addDcListeners(dc)
  }

  function getDatachannel(): RTCDataChannel {
    return dc
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
    getPeerConnection,
    setDatachannel,
    setLocalFragment,
    createOffer,
    createOfferAnswer,
    handleAnswer,
    handleIceCandidate
  }
})
