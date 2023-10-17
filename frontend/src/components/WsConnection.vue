<script setup lang="ts">
import { ref, watch } from 'vue'
import DisconnectedIcon from '@/components/icons/DisconnectedIcon.vue'
import ConnectedIcon from '@/components/icons/ConnectedIcon.vue'
import { useWebsocketStore } from '@/stores/websocket'

const hasConnection = ref(false)

const ws = useWebsocketStore()

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
  <div class="flex justify-center align-middle space-x-3">
    <disconnected-icon v-if="!hasConnection" />
    <connected-icon v-if="hasConnection" />
    <p class="text-2xl hind">
      {{ hasConnection ? 'Connected' : 'Offline' }}
    </p>
  </div>
</template>
