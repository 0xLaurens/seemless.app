<script setup lang="ts">
let imageSrc =
    'https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg/440px-Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg'
</script>

<template>
  <div class="max-w-6xl mx-auto px-4">
    <section
        class="pt-4 pb-24 md:py-24 md:mb-12 flex flex-col md:flex-row justify-center items-center md:items-start space-y-20 md:space-y-0 md:space-x-6 lg:space-x-12 text-center md:text-left">
      <div class="w-full md:w-7/12 md:mr-8">
        <h1 class="text-5xl lg:text-6xl font-extrabold tracking-tight mb-6 mt-6 dark:text-white">Share Files Directly
          <span class="text-transparent bg-clip-text bg-gradient-to-b from-primary-focus to-primary whitespace-nowrap">to any device</span>
        </h1>
        <p class="mb-4 text-lg opacity-90">Share your files peer-to-peer from any device without ever storing it
          online.</p>
        <div class="inline-flex flex-col justify-center items-center my-5 gap-2">
          <router-link class="btn btn-primary btn-wide" to="/nick">Share Now</router-link>
        </div>
      </div>
      <div class="relative w-full md:w-5/12 flex flex-col justify-start items-start gap-16 md:gap-24">
        <img :src="imageSrc" class="rounded-lg shadow-2xl" alt="tin can telephone"/>
      </div>
    </section>
    <section>
      <div
          class="features grid grid-cols-12 gap-8 content-center font-bold capitalize text-2xl mb-48"
      >
        <div class="left col-span-12 md:col-span-7 space-y-3">
          <h2 class="hind">🔥 Blazingly Fast</h2>
          <h2 class="hind">♾️ No File Size Limit</h2>
        </div>
        <div class="right col-span-12 md:col-span-5 space-y-3 gap-16 md:gap-24">
          <h2 class="hind">🔄 Peer-to-Peer</h2>
          <h2 class="hind">🔒 End-to-end encryption</h2>
        </div>
      </div>
    </section>
  </div>
</template>
