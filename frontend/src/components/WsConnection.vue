<script setup lang="ts">
import {ref, watch} from 'vue'
import DisconnectedIcon from '@/components/icons/DisconnectedIcon.vue'
import ConnectedIcon from '@/components/icons/ConnectedIcon.vue'
import {useWebsocketStore} from '@/stores/websocket'
import {useUserStore} from "@/stores/user";

const hasConnection = ref(false)

const ws = useWebsocketStore()
const user = useUserStore()

watch(
    () => ws.GetConnection(),
    (webSocket: WebSocket | undefined) => {
      webSocket?.addEventListener('open', () => (hasConnection.value = true))
      webSocket?.addEventListener('close', () => (hasConnection.value = false))
      webSocket?.addEventListener('error', () => (hasConnection.value = false))
    }
)
</script>
<template>
  <div class="flex justify-center align-middle items-center space-x-3">
    <disconnected-icon v-if="!hasConnection"/>
    <connected-icon v-if="hasConnection"/>
    <p class="text-2xl hind break-all">
      {{ hasConnection ? `Connected as: ${user.getUsername()}` : 'Connection Offline' }}
    </p>
  </div>
</template>
