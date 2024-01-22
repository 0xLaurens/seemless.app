<script setup lang="ts">
import QrcodeVue from "qrcode.vue";
import PlaneIcon from "@/components/icons/PlaneIcon.vue";
import {onMounted, ref} from "vue";
import {useRoomStore} from "@/stores/room";
import {useWebsocketStore} from "@/stores/websocket";
import type {Message} from "@/models/message";
import {RequestTypes} from "@/models/request";
import {useToastStore} from "@/stores/toast";
import {ToastType} from "@/models/toast";
import router from "@/router";

const room = useRoomStore()
const ws = useWebsocketStore()
const toast = useToastStore()

const canShare = ref(false)
const baseUrl = ref("")

const roomCode = ref("")

onMounted(() => {
  baseUrl.value = location.host

  try {
    canShare.value = navigator.canShare({
      title: "seemless.app",
      text: "hello",
    })
  } catch (error: any) {
    console.log(typeof error)
    // only available on https error can be ignored in development
    console.warn(error.message) // navigator.canShare is not a function
  }
})

function shareUrl() {
  navigator.share({
    title: "seemless.app",
    text: `${location.host}/c/${room.getRoomCode()}`,
  })
}

function copyToClipboard() {
  navigator.clipboard.writeText(`${location.host}/c/${room.getRoomCode()}`)
  toast.notify({message: "Copy successful!", type: ToastType.Info})
}

function joinRoom() {
  const message: Message = {
    type: RequestTypes.RoomJoin,
    roomCode: roomCode.value
  }
  ws.SendMessage(message)
  router.push({path: "/"})
}

</script>

<template>
  <div class="card card-normal bg-base-200 dark:bg-base-300 max-w-xl self-center">
    <div class="card-body flex items-center">
      <div>
        <p class="text-center break-words">Open this page on another device to share files
          across the
          network:</p>
      </div>
      <div>
        <div class="flex items-center w-full justify-center mb-6">
          <div class="p-3 bg-white rounded">
            <qrcode-vue class="h-32 lg:!h-56 !w-auto" :value="`${baseUrl}/c/${room.getRoomCode()}`" level="H"/>
          </div>
        </div>
      </div>
      <div>
        <div class="join">
          <input readonly class="input input-sm w-56 sm:w-auto md:input-md join-item"
                 :value="`${baseUrl}/c/${room.getRoomCode()}`"/>
          <button aria-label="copy to clipboard" @click="copyToClipboard"
                  class="btn btn-outline btn-sm md:btn-md join-item">
            <span class="hidden md:flex">Copy</span>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                 stroke="currentColor" class="w-6 h-6">
              <path stroke-linecap="round" stroke-linejoin="round"
                    d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184"/>
            </svg>
          </button>
          <button aria-label="share menu" v-if="canShare" @click="shareUrl"
                  class="btn btn-outline btn-sm md:btn-md join-item">
            <plane-icon/>
          </button>
        </div>
      </div>
      <div>
        <form @submit.prevent="roomCode.length < 5 || roomCode.length > 5" @submit="joinRoom">
          <div class="join">
            <input class="input input-sm w-56 sm:w-auto md:input-md join-item" placeholder="Enter Room Code"
                   v-model="roomCode" @input="roomCode = roomCode.toUpperCase()"
                   maxlength="5">
            <button :disabled="roomCode.length < 5 || roomCode.length > 5" aria-label="join room"
                    class="btn btn-primary btn-sm md:btn-md join-item" type="submit">
              Join Room
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>