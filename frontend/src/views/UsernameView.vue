<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import { ref } from 'vue'
import router from '@/router'

const user = useUserStore()
const username = ref('')

function cleanNickname(name: string) {
  name = name.trimStart()
  name = name.replace(/ /g, '_')
  username.value = name
}

function selectNickname() {
  if (!username.value) {
    return
  }

  user.setUsername(username.value)
  router.push({ path: '/' })
}
</script>

<template>
  <div
    class="features grid grid-cols-12 gap-8 content-center mx-48 font-bold capitalize text-2xl py-48"
  >
    <div class="col-span-12 content-center mb-6 max-w-2xl mx-auto">
      <div class="call-to-action mb-20">
        <h1 class="max-w-2xl text-4xl font-black text-white capitalize mb-5">Choose a nickname</h1>
        <p class="max-w-xl text-lg md:text-xl font-medium mb-12 hind normal-case">
          Nicknames are used as a identifier to separate different users in a room. A room must
          always have users with unique nicknames.
        </p>
      </div>
      <div class="prompt w-full">
        <h2 class="text-2xl font-black text-white capitalize mb-5 hind">
          Nickname<span v-if="username" class="hind font-black normal-case">: {{ username }}</span>
        </h2>
        <input
          type="text"
          :value="username"
          @input="cleanNickname(($event.target as HTMLInputElement).value)"
          placeholder="Your cool username"
          class="input input-primary w-full mb-5 font-regular hind"
        />
        <button :disabled="!username" @click="selectNickname" class="btn btn-primary w-full">
          Select Nickname
        </button>
      </div>
    </div>
  </div>
</template>
