<script setup lang="ts">
  let imageSrc = "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg/440px-Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg"
</script>

<template>
  <div class="hero min-h-screen bg-base-200">
    <div class="hero-content flex-col lg:flex-row-reverse">
      <img :src="imageSrc" class="rounded-lg shadow-2xl"  alt="tin can telephone"/>
      <div>
        <h1 class="max-w-2xl text-4xl md:text-5xl lg:text-6xl font-black text-white capitalize mb-5">Share files directly from anywhere to anywhere</h1>
        <p class="max-w-xl text-lg md:text-xl font-medium mb-12">Send files of any size directly to another device without ever storing it online.</p>
        <a href="#" class="btn btn-primary btn-md btn-wide">
          Share Now
        </a>
      </div>
    </div>
  </div>
  </template>
