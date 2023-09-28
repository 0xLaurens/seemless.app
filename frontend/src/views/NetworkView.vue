<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import RoomCard from '@/components/RoomCard.vue'

const user = useUserStore()
</script>

<template>
  <div
    class="features grid grid-cols-12 gap-8 content-center mx-48 font-bold capitalize text-2xl py-48"
  >
    <div class="col-span-12">
      <div class="call-to-action mb-20">
        <h1 class="max-w-2xl text-4xl font-black text-white capitalize mb-5">Network</h1>
        <p class="max-w-xl text-lg md:text-xl font-medium mb-12 hind normal-case">
          Welcome {{ user.getUsername() || 'name' }}, please select a network you would like to be
          discoverable in.
        </p>
      </div>
      <div class="rooms flex">
        <RoomCard
          modal-title="Local Network"
          modal-message="The Local Network is based on IP ranges. People with the same public IP address will be connected to each other."
          modal-btn-text="I understand"
          title="Local Network"
          description="Share Files on local network"
        />
        <RoomCard
          modal-title="Outside Network"
          modal-message="If you want to connect to members outside of your network connection you will have to use the outside local network option. The server will use a stun server to find routing to each other"
          modal-btn-text="I understand"
          title="Outside Network"
          description="Share Files with people outside of your network"
        />
      </div>
    </div>
  </div>
</template>
