<script setup lang="ts">
import { useRoute } from 'vue-router'
import QrIcon from '@/components/icons/QrIcon.vue'
import PlaneIcon from '@/components/icons/PlaneIcon.vue'
import BackIcon from '@/components/icons/BackIcon.vue'
import { useUserStore } from '@/stores/user'
import { onMounted, onUnmounted, ref } from 'vue'
import type { Ref } from 'vue'
import { RequestTypes } from '@/models/request'
import WsConnection from '@/components/WsConnection.vue'
import type { Message } from '@/models/message'
import { useRtcStore } from '@/stores/rtc'
import FileInput from '@/components/FileInput.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { useConnectedStore } from '@/stores/connected'
import { useRetryStore } from '@/stores/retry'

const user = useUserStore()
const conn = useConnectedStore()
const route = useRoute()
const id = route.params.id

const rtc = useRtcStore()
const retry = useRetryStore()
let pc: RTCPeerConnection
const ws: Ref<WebSocket | undefined> = ref()
const users = ref([])

onMounted(() => {
  openWS()
  pc = rtc.getPeerConnection()
  pc.ondatachannel = (dce) => rtc.setDatachannel(dce.channel)
  pc.onicecandidate = (ev) => {
    if (!ev.candidate) return
    rtc.setLocalFragment(ev.candidate.usernameFragment)
    console.log(conn.getIceTargets())
    // for (const user of conn.getConnectedUsers()) {
    console.log('ICE USER')
    sendMessage(RequestTypes.NewIceCandidate, undefined, undefined, JSON.stringify(ev.candidate))
    // }
  }
})

function sendMessage(type: string, target?: string, sdp?: string, candidate?: string) {
  console.log('SENT:', type)
  let message: Message = {
    type: type,
    target: target,
    sdp: sdp,
    candidate: candidate
  }

  message.from = user.getUsername()
  ws.value?.send(JSON.stringify(message))
}

async function sendOffer(username: string) {
  if (username === user.getUsername()) {
    return
  }

  let offer = await rtc.createOffer()
  conn.createUserConnection(username, true)
  sendMessage('Offer', username, offer.sdp)
}

function openWS() {
  ws.value = new WebSocket(import.meta.env.VITE_WS_ADDR)
  ws.value.onmessage = async (event) => {
    let data = JSON.parse(event.data)
    switch (data.type) {
      case RequestTypes.Offer: {
        let answer = await rtc.createOfferAnswer(data)
        sendMessage('Answer', data.from, answer.sdp)
        break
      }
      case RequestTypes.NewIceCandidate: {
        console.log(data.from)
        await rtc.handleIceCandidate(JSON.parse(data.candidate))
        break
      }
      case RequestTypes.Answer: {
        await rtc.handleAnswer(data)
        break
      }
      case RequestTypes.Peers:
        users.value = []
        for (let user of data.users) {
          if (user.username) users.value.push(user.username)
        }
        break
      case RequestTypes.PeerJoined:
        users.value.push(data?.username)
        break
      case RequestTypes.PeerLeft:
        users.value = users.value.filter((u) => u != data.username)
        conn.removeUser(data.username)
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
  ws.value.onclose = () => {
    console.log('CONNECTION CLOSED')
    retry.start(reconnectWs)
  }
  ws.value.onopen = () => {
    console.log('CONNECTION OPENED')
    retry.stop()
    let payload = {
      type: 'Username',
      body: {
        username: user.getUsername()
      }
    }

    ws.value?.send(JSON.stringify(payload))
  }
}

function reconnectWs() {
  if (ws.value?.readyState !== ws.value?.CLOSED) {
    ws.value?.close()
  }
  openWS()
  console.log('Reconnecting...')
}

onUnmounted(() => {
  console.log('UNMOUNTED')
  ws.value?.close()
  for (const u of conn.getIceTargets()) {
    conn.removeUser(u)
  }
  rtc.close()
})
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
            <div
              :key="u"
              v-for="u in users"
              @click="u !== user.getUsername() ? sendOffer(u) : null"
            >
              <user-avatar :user="u" />
            </div>
          </div>
        </div>

        <div>
          <ws-connection :ws="ws" />
        </div>

        <div class="relative mb-24">
          <div class="flex justify-center flex-wrap space-x-6 pt-5">
            <file-input />
          </div>
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
