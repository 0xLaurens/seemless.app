<template>
  <div class="btm-nav btm-nav-md flex md:hidden bg-base-200 z-50">
    <router-link aria-label="home" active-class="text-primary active bg-base-200" to="/">
      <home-icon/>
    </router-link>
    <router-link aria-label="share" active-class="text-primary active bg-base-200" to="/share">
      <qr-icon/>
    </router-link>
    <router-link aria-label="settings" active-class="text-primary active bg-base-200" to="/settings">
      <settings-icon/>
    </router-link>
  </div>
</template>
<script setup lang="ts">
import HomeIcon from "@/components/icons/HomeIcon.vue";
import QrIcon from "@/components/icons/QrIcon.vue";
import SettingsIcon from "@/components/icons/SettingsIcon.vue";
</script>