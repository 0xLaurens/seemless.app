<script setup lang="ts">
import UserBar from "@/components/UserBar.vue";
import {onDeactivated, onMounted} from "vue";
import {useWebsocketStore} from "@/stores/websocket";
import UsersTable from "@/components/UsersTable.vue";
import IncomingTable from "@/components/IncomingTable.vue";
import OutgoingTable from "@/components/OutgoingTable.vue";
import {useRoute} from "vue-router";
import type {Message} from "@/models/message";
import {RequestTypes} from "@/models/request";

const ws = useWebsocketStore()

const route = useRoute()

onMounted(() => {
  ws.Open()

  if (!route.params.code) {
    return
  }

  setTimeout(() => {
    console.log("sending join request")
    if (typeof route.params.code !== "string") return

    const message: Message = {
      type: RequestTypes.RoomJoin,
      roomCode: route.params.code
    }
    ws.SendMessage(message)
  }, 300)
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