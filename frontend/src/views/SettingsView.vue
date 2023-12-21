<script setup lang="ts">
import {useSettingsStore} from "@/stores/settings";

const settings = useSettingsStore()

const localSettings = settings.getSettings()
</script>

<template>
  <div class="card card-normal bg-base-200 dark:bg-base-300 w-full self-center">
    <div class="card-body flex">
      <div class="card-title">
        <p>Settings</p>
      </div>
      <div class="form-control" v-for="setting of localSettings" :key="setting.name">
        <label class="label cursor-pointer">
          <input type="checkbox" class="toggle toggle-primary" @change="settings.saveSettings(localSettings)"
                 v-model="setting.value"/>
          <span class="label-text pl-3">{{ setting.description }}</span>
        </label>
      </div>
    </div>
  </div>

</template>