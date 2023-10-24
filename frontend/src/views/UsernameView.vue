<script setup lang="ts">
import {useUserStore} from '@/stores/user'
import {ref} from 'vue'
import router from '@/router'

const user = useUserStore()
const username = ref(user.getUsername())

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
  router.push({path: '/room/local'})
}
</script>

<template>
  <div class="flex justify-start md:justify-center">
    <div class="max-w-6xl mx-auto align-middle justify-center px-4">
      <section class="flex flex-col h-screen justify-between md:justify-center">
        <div class="call-to-action pt-10 md:pt-0 md:mb-6">
          <h1 class="max-w-2xl text-4xl font-black text-black dark:text-white capitalize mb-5">Choose a nickname</h1>
          <p class="max-w-xl text-lg md:text-xl font-medium mb-12 hind normal-case">
            Nicknames are used as a identifier to separate different users in a room. A room must
            always have users with unique nicknames.
          </p>
        </div>
        <div class="pb-10 md:pb-0">
          <form @submit="selectNickname">
            <label class="text-xl md:text-2xl font-black text-black dark:text-white capitalize mb-5 hind">
              Nickname<span v-if="username" class="hind font-black normal-case">: {{ username }}</span>
            </label>
            <input
                type="text"
                :value="username"
                @input="cleanNickname(($event.target as HTMLInputElement).value)"
                placeholder="Your cool username"
                class="input input-primary w-full mb-5 font-regular hind"
            />
            <button :disabled="!username" class="btn btn-primary w-full">Select Nickname</button>
          </form>
        </div>
      </section>
    </div>
  </div>
</template>
