<script setup lang="ts">
let imageSrc =
  'https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg/440px-Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg'
</script>

<template>
  <div
    class="features grid grid-cols-12 gap-8 content-center mx-48 font-bold capitalize text-2xl py-48"
  >
    <div class="left col-span-7">
      <h1 class="max-w-2xl text-4xl md:text-5xl lg:text-6xl font-black text-white capitalize mb-5">
        Share files directly from anywhere to anywhere
      </h1>
      <p class="max-w-xl text-lg md:text-xl font-medium mb-12 hind">
        Send files of any size directly to another device without ever storing it online.
      </p>
      <router-link to="/nick" class="btn btn-primary btn-md btn-wide"> Share Now</router-link>
    </div>
    <div class="right col-span-5 space-y-3 mx-auto">
      <img :src="imageSrc" class="rounded-lg shadow-2xl" alt="tin can telephone" />
    </div>
  </div>
  <div
    class="features grid grid-cols-12 gap-8 content-center mx-48 font-bold capitalize text-2xl mb-48"
  >
    <div class="left col-span-7 space-y-3">
      <h2 class="hind">🔥 Blazingly Fast</h2>
      <h2 class="hind">♾️ No File Size Limit</h2>
    </div>
    <div class="right col-span-5 space-y-3">
      <h2 class="hind">🔄 Peer-to-Peer</h2>
      <h2 class="hind">🔒 End-to-end encryption</h2>
    </div>
  </div>
</template>
