<script setup lang="ts">
import MobileIcon from "@/components/icons/devices/MobileIcon.vue";
import ComputerIcon from "@/components/icons/devices/ComputerIcon.vue";
import type {User} from "@/models/user";

defineProps<{ user: User | undefined }>()
</script>

<template>
  <mobile-icon
      v-if="user?.device.toLowerCase().includes('android') || user?.device.toLowerCase().includes('ios')"
      class="h-fit invert dark:invert-0 w-auto text-white font-black"/>
  <computer-icon class="h-fit invert dark:invert-0 w-auto text-white font-black" v-else/>
</template>