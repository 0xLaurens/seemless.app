<script setup lang="ts">
import InfoIcon from "@/components/icons/InfoIcon.vue";

</script>

<template>
  <nav class="navbar w-full">
    <router-link aria-label="home" class="btn btn-ghost text-xl normal-case" to="/">Seemless.app</router-link>
    <div class="ml-auto space-x-3">
      <router-link aria-label="about" to="/about" class="btn btn-square btn-ghost">
        <info-icon class="!w-8 !h-8"/>
      </router-link>
    </div>
  </nav>
</template>
