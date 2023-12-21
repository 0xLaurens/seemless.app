<script setup lang="ts">
import {ref} from "vue";
import QrIcon from "@/components/icons/QrIcon.vue";
import {useUserStore} from "@/stores/user";
import EditIcon from "@/components/icons/EditIcon.vue";
import CheckIcon from "@/components/icons/CheckIcon.vue";
import XIcon from "@/components/icons/XIcon.vue";
import SettingsIcon from "@/components/icons/SettingsIcon.vue";
import DeviceIcon from "@/components/icons/devices/DeviceIcon.vue";

const user = useUserStore()

const nameEditMode = ref(false)
const updatedName = ref("")

function editName() {
  updatedName.value = user.getUsername()
  nameEditMode.value = !nameEditMode.value
}

function updateName() {
  console.log(updatedName)
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
          <button class="btn btn-sm btn-square btn-ghost" @click="editName">
            <edit-icon/>
          </button>
        </div>
        <div class="flex" v-if="nameEditMode">
          <input v-model="updatedName" class="input input-sm w-48 md:w-auto flex-grow" @keydown.enter="updateName"/>
          <div class="flex space-x-3 ms-3">
            <button @click="updateName" class="btn btn-sm btn-success btn-circle">
              <check-icon/>
            </button>
            <button @click="editName" class="btn btn-sm btn-error btn-circle">
              <x-icon/>
            </button>
          </div>
        </div>
      </div>
      <div class="ml-auto flex items-center">
        <router-link to="/share" class="btn btn-square btn-ghost hidden md:flex">
          <qr-icon/>
        </router-link>
        <router-link to="/settings" class="btn btn-square btn-ghost hidden md:flex">
          <settings-icon/>
        </router-link>
      </div>
    </div>
  </div>
</template>