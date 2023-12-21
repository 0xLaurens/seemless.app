<script setup lang="ts">
import ToastArea from '@/components/ToastArea.vue'
import MobileNav from "@/components/MobileNav.vue";
import NavigationBar from "@/components/NavigationBar.vue";
import DownloadPopup from "@/components/DownloadPopup.vue";
</script>

<template>
  <ToastArea/>
  <DownloadPopup/>
  <section class="max-w-7xl mx-auto flex flex-col min-h-screen px-2 xl:px-0">
    <NavigationBar/>
    <RouterView/>
  </section>
  <MobileNav/>
</template>
