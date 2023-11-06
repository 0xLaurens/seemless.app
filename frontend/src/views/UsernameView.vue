<script setup lang="ts">
import {useUserStore} from '@/stores/user'
import {ref} from 'vue'
import router from '@/router'
import BackIcon from "@/components/icons/BackIcon.vue";
import {useToastStore} from "@/stores/toast";
import {ToastType} from "@/models/toast";

const user = useUserStore()
const username = ref(user.getUsername())
const toast = useToastStore()

function cleanNickname(name: string) {
  name = name.trimStart()
  name = name.replace(/ /g, '_')
  username.value = name
}

function selectNickname() {
  if (!username.value) {
    return
  }

  if (username.value.length > 16) {
    toast.notify({message: "Max length of username exceeded", type: ToastType.Warning})
    return
  }

  user.setUsername(username.value)
  router.push({path: '/room/local'})
}
</script>

<template>
  <section class="flex h-screen items-center justify-center py-10">
    <div
        class="flex flex-col items-center max-w-3xl justify-between h-full md:h-auto md:justify-start w-full px-3 md:px-6">
      <div class="flex flex-row justify-between items-center w-full mb-10">
        <div>
          <router-link to="../" class="btn btn-md md:btn-sm btn-ghost font-black">
            <back-icon/>
          </router-link>
        </div>
        <div></div>
        <div></div>
      </div>
      <div>
        <h1 class="text-4xl font-black text-black dark:text-white capitalize mb-5">Choose a nickname</h1>
        <p class="text-lg md:text-xl font-medium mb-12 hind normal-case">
          Nicknames are used as a identifier to separate different users in a room. A room must
          always have users with unique nicknames.
        </p>
      </div>
      <div class="w-full">
        <form @submit="selectNickname" @submit.prevent="username.length > 16">
          <label class="text-xl md:text-2xl font-black text-black dark:text-white capitalize mb-5 hind break-words">
            Nickname<span v-if="username" class="hind font-black normal-case">: {{ username }}</span>
          </label>
          <input
              type="text"
              maxlength="16"
              :value="username"
              @input="cleanNickname(($event.target as HTMLInputElement).value)"
              placeholder="Your cool username"
              class="input input-accent w-full mb-5 font-regular hind"
          />
          <button :disabled="!username || username.length > 16" class="btn btn-accent w-full">Select Nickname</button>
        </form>
      </div>
    </div>
  </section>

  <!--  <div class="flex justify-start md:justify-center">-->
  <!--    <div class="max-w-6xl mx-auto align-middle justify-center px-4">-->
  <!--      <section class="flex flex-col h-screen justify-between md:justify-center">-->
  <!--        <div class="call-to-action md:mb-6">-->
  <!--          <div>-->
  <!--            <router-link to="../" class="btn btn-md md:btn-lg btn-ghost font-black mt-10 md:pt-0 mb-3">-->
  <!--              <back-icon/>-->
  <!--            </router-link>-->
  <!--          </div>-->
  <!--          <h1 class="max-w-2xl text-4xl font-black text-black dark:text-white capitalize mb-5">Choose a nickname</h1>-->
  <!--          <p class="max-w-xl text-lg md:text-xl font-medium mb-12 hind normal-case">-->
  <!--            Nicknames are used as a identifier to separate different users in a room. A room must-->
  <!--            always have users with unique nicknames.-->
  <!--          </p>-->
  <!--        </div>-->
  <!--        <div class="pb-10 md:pb-0">-->
  <!--          <form @submit="selectNickname">-->
  <!--            <label class="text-xl md:text-2xl font-black text-black dark:text-white capitalize mb-5 hind">-->
  <!--              Nickname<span v-if="username" class="hind font-black normal-case">: {{ username }}</span>-->
  <!--            </label>-->
  <!--            <input-->
  <!--                type="text"-->
  <!--                :value="username"-->
  <!--                @input="cleanNickname(($event.target as HTMLInputElement).value)"-->
  <!--                placeholder="Your cool username"-->
  <!--                class="input input-primary w-full mb-5 font-regular hind"-->
  <!--            />-->
  <!--            <button :disabled="!username" class="btn btn-primary w-full">Select Nickname</button>-->
  <!--          </form>-->
  <!--        </div>-->
  <!--      </section>-->
  <!--    </div>-->
  <!--  </div>-->
</template>
