<script setup lang="ts">
import UserBar from "@/components/UserBar.vue";
import {onDeactivated, onMounted} from "vue";
import {useWebsocketStore} from "@/stores/websocket";
import UsersTable from "@/components/UsersTable.vue";
import IncomingTable from "@/components/IncomingTable.vue";
import OutgoingTable from "@/components/OutgoingTable.vue";

const ws = useWebsocketStore()

onMounted(() => {
  ws.Open()
})

onDeactivated(() => {
  ws.Close()
})
</script>

<template>
  <section id="me">
    <UserBar/>
  </section>

  <section id="users">
    <UsersTable/>
  </section>

  <section id="incoming">
    <IncomingTable/>
  </section>

  <section id="outgoing">
    <OutgoingTable/>
  </section>
</template>