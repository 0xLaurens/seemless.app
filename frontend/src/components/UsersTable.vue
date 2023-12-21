<script setup lang="ts">
import {useUserStore} from "@/stores/user";
import UserCard from "@/components/UserCard.vue";
import FileInput from "@/components/FileInput.vue";
import {type Ref, ref} from "vue";

const user = useUserStore()

const selectedUsers: Ref<Map<string, boolean>> = ref(new Map());

function addUser(selected: boolean, username: string) {
  if (selected) {
    selectedUsers.value.set(username, selected)
  } else {
    selectedUsers.value.delete(username)
  }
}

</script>

<template>
  <div class="card card-compact bg-base-200 dark:bg-base-300 mt-4 shadow-sm w-full h-full">
    <div class="card-body">
      <div class="card-title">
        <p>Users</p>
        <router-link v-if="user.users.length < 1" to="/share" class="btn btn-sm btn-outline">invite users</router-link>
      </div>
      <p v-if="user.users.length < 1">Nobody is connected to the same network.</p>
      {{ selectedUsers }}
      <div :key="u.id" id="users" v-for="u in user.users">
        <UserCard :user="u" @selected="addUser"/>
      </div>
      <div v-if="user.users.length > 0">
        <file-input :users="selectedUsers"/>
      </div>
    </div>
  </div>
</template>