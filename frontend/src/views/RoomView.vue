<script setup lang="ts">
import {useRoute} from 'vue-router'
import QrIcon from '@/components/icons/QrIcon.vue'
import BackIcon from '@/components/icons/BackIcon.vue'
import {useUserStore} from '@/stores/user'
import {onMounted, onUnmounted} from 'vue'
import FileInput from '@/components/FileInput.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import ConfirmDownload from '@/components/ConfirmDownload.vue'
import {useWebsocketStore} from '@/stores/websocket'
import WsConnection from '@/components/WsConnection.vue'
import FileStatusCard from "@/components/FileStatusCard.vue";
import {useFileStore} from "@/stores/file";

const user = useUserStore()
const ws = useWebsocketStore();
const file = useFileStore();

const route = useRoute()
const id = route.params.id

onMounted(() => {
  ws.Open()
})

onUnmounted(() => {
  ws.Close()
})
</script>

<template>
  <confirm-download/>
  <section class="flex h-screen items-center justify-center py-10">
    <div class="flex flex-col items-center max-w-3xl w-full h-full justify-between px-3 md:px-0">
      <div class="flex flex-row justify-between items-center w-full">
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
      <div class="users-box flow-root">
        <div class="flex flex-wrap justify-center align-middle">
          <div :key="u.username" v-for="u in user.users">
            <user-avatar :user="u"/>
          </div>
          <p v-if="user.users.length < 1" class="break-words">Wait for other users to connect...</p>
        </div>
      </div>
      <div class="space-y-3" v-if="file.getCurrentOffer().value">
        <p class="text-bold text-xl">Files receiving.</p>
        <div v-bind:key="f.name" v-for="f in file.getCurrentOffer().value?.files" class="w-full">
          <file-status-card
              :target="file.getCurrentOffer().value?.target"
              :from="file.getCurrentOffer().value?.from"
              :file="f" class="w-full"
          />
        </div>
      </div>
      <div class="space-y-3" v-if="file.getSendOffer().value">
        <p class="text-bold text-xl">Files sending.</p>
        <div v-bind:key="f.name" v-for="f in file.getSendOffer().value?.files" class="w-full">
          <file-status-card
              :target="file.getSendOffer().value?.target"
              :from="file.getSendOffer().value?.from"
              :file="f" class="w-full"
          />
        </div>
      </div>
      <div>
        <div class="flex justify-center">
          <file-input/>
        </div>
        <div>
          <ws-connection/>
        </div>
      </div>
    </div>
  </section>
</template>
