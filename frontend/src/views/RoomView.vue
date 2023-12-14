<script setup lang="ts">
import {useRoute} from 'vue-router'
import QrIcon from '@/components/icons/QrIcon.vue'
import BackIcon from '@/components/icons/BackIcon.vue'
import {useUserStore} from '@/stores/user'
import {onMounted, onUnmounted, type Ref, ref} from 'vue'
import FileInput from '@/components/FileInput.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import ConfirmDownload from '@/components/ConfirmDownload.vue'
import {useWebsocketStore} from '@/stores/websocket'
import WsConnection from '@/components/WsConnection.vue'
import FileStatusCard from "@/components/FileStatusCard.vue";
import {useFileStore} from "@/stores/file";
import type {Message} from "@/models/message";
import {RequestTypes} from "@/models/request";
import {useRoomStore} from "@/stores/room";
import type {RoomCode} from "@/models/room";

const user = useUserStore()
const ws = useWebsocketStore();
const file = useFileStore();
const room = useRoomStore();

const route = useRoute()
const id = route.params.id
const code: Ref<RoomCode | undefined> = ref()

onMounted(() => {
  ws.Open()
})

onUnmounted(() => {
  ws.Close()
})
// function createPublicRoom() {
//   const msg: Message = {
//     type: RequestTypes.PublicRoomCreate
//   }
//   ws.SendMessage(msg)
// }
//
// function joinPublicRoom() {
//   if (!code.value) {
//     return
//   }
//
//   const msg: Message = {
//     type: RequestTypes.PublicRoomJoin,
//     roomCode: code.value
//   }
//   ws.SendMessage(msg)
// }
</script>

<template>
  <confirm-download/>
  <section class="items-center py-10 my-auto px-3">
    <div class="flex flex-row justify-between mx-auto items-center max-w-6xl pb-10">
      <div>
        <router-link to="../nick" class="btn btn-md md:btn-lg btn-ghost font-black">
          <back-icon/>
        </router-link>
      </div>
      <div>
        <h1 class="text-2xl lg:text-4xl font-black text-center text-black dark:text-white capitalize">
          Room: {{ id }}
        </h1>
      </div>
      <div>
        <button class="btn btn-md md:btn-lg btn-accent btn-outline">
          <qr-icon/>
        </button>
      </div>
    </div>


<!--    <button class="btn btn-primary" @click="createPublicRoom">CREATE PUBLIC ROOM</button>-->
<!--    <span v-if="room.getRoomCode()">connected to {{room.getRoomCode()}}</span>-->
<!--    <br/>-->
<!--    <br/>-->
<!--    <input class="input input-primary" type="text" v-model="code"/>-->
<!--    <button type="submit" class="btn btn-primary" @click="joinPublicRoom">JOIN PUBLIC ROOM</button>-->


    <div class="flex flex-1 flex-row w-full px-3 md:pl-10 space-x-3 md:h-[60vh] items-center justify-between">
      <div class="items-center w-full "
           :class="(file.getReceiveOffer()?.value == undefined && file.getSendOffer()?.value  == undefined) ? '' : 'max-w-2xl'">
        <div class="users-box flow-root">
          <div class="flex flex-wrap justify-center align-middle">
            <div :key="u.username" v-for="u in user.users">
              <user-avatar :user="u"/>
            </div>
            <p v-if="user.users.length < 1" class="break-words">Wait for other users to connect...</p>
          </div>
        </div>
      </div>
      <div class="files card bg-neutral pt-3 pb-6 max-w-2xl w-full max-h-[45vh]"
           :class="(file.getReceiveOffer()?.value == undefined && file.getSendOffer()?.value  == undefined) ? 'hidden' : 'hidden md:flex md:flex-col'">
        <h2 class="text-2xl pl-3">Files.</h2>
        <div class="divider my-0"></div>
        <div class="w-full max-h-96 overflow-y-scroll">
          <div v-bind:key="f.name" v-for="f in file.getReceiveOffer().value?.files" class="w-full">
            <file-status-card
                :target="file.getReceiveOffer().value?.target"
                :from="file.getReceiveOffer().value?.from"
                :file="f" class="w-full"
            />
          </div>
          <div v-bind:key="f.name" v-for="f in file.getSendOffer().value?.files" class="w-full">
            <file-status-card
                :target="file.getSendOffer().value?.target"
                :from="file.getSendOffer().value?.from"
                :file="f" class="w-full"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="pt-6">
      <div class="flex justify-center">
        <file-input/>
      </div>
      <div>
        <ws-connection/>
      </div>
    </div>
  </section>
  <div class="files flex flex-col md:hidden w-full space-y-3 px-3">
    <p class="text-bold text-xl">Files.</p>
    <div v-bind:key="f.name" v-for="f in file.getReceiveOffer().value?.files" class="w-full">
      <file-status-card
          :target="file.getReceiveOffer().value?.target"
          :from="file.getReceiveOffer().value?.from"
          :file="f" class="w-full"
      />
    </div>
    <div v-bind:key="f.name" v-for="f in file.getSendOffer().value?.files" class="w-full">
      <file-status-card
          :target="file.getSendOffer().value?.target"
          :from="file.getSendOffer().value?.from"
          :file="f" class="w-full"
      />
    </div>
  </div>

</template>
