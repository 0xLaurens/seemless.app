<script setup lang="ts">
import { useUserStore } from '@/stores/user'

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
      <div class="rooms flex justify-between">
        <div class="card bg-base-300 w-96 shadow-2xl py-6 ml-0 m-12">
          <button class="btn btn-primary btn-circle btn-sm self-end mr-3">?</button>
          <div class="card-body flex items-center h-full self-center text-center">
            <h2 class="card-title text-white text-2xl font-extrabold">Local Network</h2>
            <p class="hind text-base normal-case text-gray-500">Share files on local network</p>
          </div>
        </div>

        <div class="card bg-base-300 w-96 shadow-2xl py-6 ml-0 m-12">
          <button class="btn btn-primary btn-circle btn-sm self-end mr-3">?</button>
          <div class="card-body flex items-center text-center">
            <h2 class="card-title text-white text-2xl font-extrabold">Outside Local Network</h2>
            <p class="hind text-base normal-case text-gray-500">
              Share files with people on other connections
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
