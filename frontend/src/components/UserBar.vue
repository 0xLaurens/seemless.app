<script setup lang="ts">
import {ref} from "vue";
import QrIcon from "@/components/icons/QrIcon.vue";
import {useUserStore} from "@/stores/user";
import EditIcon from "@/components/icons/EditIcon.vue";
import CheckIcon from "@/components/icons/CheckIcon.vue";
import XIcon from "@/components/icons/XIcon.vue";
import SettingsIcon from "@/components/icons/SettingsIcon.vue";
import DeviceIcon from "@/components/icons/devices/DeviceIcon.vue";
import {useWebsocketStore} from "@/stores/websocket";
import {RequestTypes} from "@/models/request";

const user = useUserStore()
const ws = useWebsocketStore()

const nameEditMode = ref(false)
const updatedName = ref("")

function editName() {
  updatedName.value = user.getUsername()
  nameEditMode.value = !nameEditMode.value
}

function updateName() {
  ws.SendMessage({type: RequestTypes.ChangeDisplayName, displayName: updatedName.value})
  editName()
}

</script>

<template>
  <div class="card card-compact bg-base-200 dark:bg-base-300">
    <div class="card-body flex flex-row min-w-xs">
      <div class="flex items-center justify-end">
        <div class="avatar rounded placeholder">
          <div class="w-10 md:w-12  items-center p-1">
            <device-icon :user="user.getCurrentUser()"/>
          </div>
        </div>
        <div class="flex items-center md:ps-1" v-if="!nameEditMode">
          <p class="font-medium text-xl">{{ user.getUsername() || 'username' }}</p>
          <button aria-label="edit name" class="btn btn-sm btn-square btn-ghost" @click="editName">
            <edit-icon/>
          </button>
        </div>
        <div class="flex" v-if="nameEditMode">
          <input aria-label="username" v-model="updatedName" class="input input-sm w-48 md:w-auto flex-grow"
                 @keydown.enter="updateName"/>
          <div class="flex space-x-3 ms-3">
            <button aria-label="confirm change" @click="updateName" class="btn btn-sm btn-success btn-circle">
              <check-icon/>
            </button>
            <button aria-label="cancel change" @click="editName" class="btn btn-sm btn-error btn-circle">
              <x-icon/>
            </button>
          </div>
        </div>
      </div>
      <div class="ml-auto flex items-center">
        <router-link aria-label="share" to="/share" class="btn btn-square btn-ghost hidden md:flex">
          <qr-icon/>
        </router-link>
        <router-link aria-label="settings" to="/settings" class="btn btn-square btn-ghost hidden md:flex">
          <settings-icon/>
        </router-link>
      </div>
    </div>
  </div>
</template>