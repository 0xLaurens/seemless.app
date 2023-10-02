<script setup lang="ts">
import { useRoute } from 'vue-router'
import QrIcon from '@/components/icons/QrIcon.vue'
import PlaneIcon from '@/components/icons/PlaneIcon.vue'
import BackIcon from '@/components/icons/BackIcon.vue'
import { useUserStore } from '@/stores/user'
import { onMounted, onUnmounted, ref } from 'vue'
import { RequestTypes } from '@/models/request'
import WsConnection from '@/components/WsConnection.vue'
import type { Message } from '@/models/message'
// import { createOffer, createOfferAnswer, handleAnswer } from '@/helpers/rtc'
import type { SessionDescriptionMessage } from '@/models/sdm'

const user = useUserStore()
const route = useRoute()
const id = route.params.id

const rtc = new RTCPeerConnection()
// const ws = new WebSocket('ws://192.168.14.249:3000/ws')
const ws = new WebSocket('ws://127.0.0.1:3000/ws')
let localFragment: string | null

let users = ref([])
const dataChannel = rtc.createDataChannel('test')

rtc.ondatachannel = (dc) => {
  dc.channel.onmessage = (event) => {
    console.log(event.data)
  }
}

dataChannel.onerror = (error) => {
  console.log('Data Channel Error:', error)
}

dataChannel.onmessage = (event) => {
  console.log(event)
  console.log('Got Data Channel Message:', event.data)
}

dataChannel.onopen = () => {
  console.log('DATA CHANNEL OPEN')
  dataChannel.send('Hello World!')
}

dataChannel.onclose = () => {
  console.log('The Data Channel is Closed')
}

onMounted(() => {
  rtc.onicecandidate = (ev) => {
    if (ev.candidate) {
      localFragment = ev.candidate.usernameFragment
      sendMessage(RequestTypes.NewIceCandidate, undefined, undefined, JSON.stringify(ev.candidate))
    }
  }
  rtc.onicecandidateerror = (ev) => console.log(ev)
})

ws.onmessage = async (event) => {
  let data = JSON.parse(event.data)
  switch (data.type) {
    case RequestTypes.Offer: {
      await createOfferAnswer(data)
      break
    }
    case RequestTypes.NewIceCandidate: {
      await handleIceCandidate(JSON.parse(data.candidate))
      break
    }
    case RequestTypes.Answer: {
      await handleAnswer(data)
      break
    }
    case RequestTypes.Peers:
      for (let user of data.users) {
        users.value.push(user.username)
      }
      break
    case RequestTypes.PeerJoined:
      users.value.push(data.username)
      break
    case RequestTypes.PeerLeft:
      users.value = users.value.filter((u) => u != data.username)
      break
    case RequestTypes.UsernamePrompt:
    case RequestTypes.Username:
    case RequestTypes.DuplicateUsername:
      console.log(data.message)
      break
    default:
      console.log(`Unknown type ${data.type}`)
  }
}

function sendMessage(type: string, target?: string, sdp?: string, candidate?: string) {
  console.log('SENT:', type)
  let message: Message = {
    type: type,
    target: target,
    sdp: sdp,
    candidate: candidate
  }

  if (message.sdp) {
    message.from = user.getUsername()
  }
  ws.send(JSON.stringify(message))
}

async function sendOffer(username: string) {
  let offer = await createOffer()
  console.log(offer)
  sendMessage('Offer', username, offer.sdp)
}

ws.onopen = () => {
  let payload = {
    type: 'Username',
    body: {
      username: user.getUsername()
    }
  }

  ws.send(JSON.stringify(payload))
}

ws.onerror = (e) => {
  console.log(e)
}

onUnmounted(() => {
  ws.close()
})

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
  sendMessage('Answer', offer.from, answer.sdp)
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
  if (candidate == null || candidate.usernameFragment === localFragment) {
    return
  }
  await rtc.addIceCandidate(candidate)
}

// async function handleIceCandidate(
//   rtc: RTCPeerConnection,
//   data: { type: string; candidate: string }
// ) {
//   const candidate = JSON.parse(data.candidate)
//   await rtc.addIceCandidate(candidate)
// }
</script>

<template>
  <div class="pt-10 pb-4 sm:pt-16 lg:overflow-hidden lg:pt-8 lg:pb-4">
    <div class="mx-auto max-w-7xl lg:px-8">
      <div class="w-full px-4 sm:px-6 sm:text-center lg:px-0">
        <div class="z-0 relative lg:pt-12">
          <div class="flex justify-between align-middle">
            <router-link to="../network" class="btn btn-md md:btn-lg btn-ghost font-black">
              <back-icon />
            </router-link>
            <h1 class="text-4xl font-black text-center text-white capitalize mb-24">
              Room: {{ id }}
            </h1>
            <button class="btn btn-md md:btn-lg btn-accent btn-outline">
              <qr-icon />
            </button>
          </div>
        </div>

        <div class="users-box flow-root mb-24">
          <div class="flex flex-wrap justify-center align-middle">
            <div :key="user" v-for="user in users" @click="sendOffer(user)">
              <div class="user justify-center text-center mr-6">
                <div class="avatar placeholder mb-2">
                  <div
                    class="bg-neutral-focus text-neutral-content rounded-full w-24 hover:ring-4 ring-accent ring-offset-4 ring-offset-base-100"
                  >
                    <span class="text-3xl">{{ user[0] }}</span>
                  </div>
                </div>
                <h2>{{ user }}</h2>
                <span class="hind text-accent text-lg font-medium">Android</span>
              </div>
            </div>
          </div>
        </div>

        <div>
          <ws-connection :ws="ws" />
        </div>

        <div class="relative mb-24">
          <div class="flex justify-center flex-wrap space-x-6 pt-5">
            <button class="btn btn-md md:btn-lg btn-primary">Upload File</button>
            <button class="btn btn-md md:btn-lg btn-primary">
              Broadcast
              <plane-icon />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
