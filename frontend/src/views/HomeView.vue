<script setup lang="ts">
  let imageSrc = "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg/440px-Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg"
</script>

<template>
  <div class="grid grid-cols-12 gap-12 mx-28 my-48 content-center">
    <div class="col-span-12 md:col-span-6">
      <div class="flex flex-col justify-center items-center md:items-start h-full">
        <h1 class="text-4xl md:text-5xl lg:text-6xl font-black text-white capitalize mb-5">Share files directly from anywhere to anywhere</h1>
        <p class="text-lg md:text-xl font-medium mb-3">Send files of any size directly to another device without ever storing it online.</p>
        <a href="#" class="btn btn-primary btn-lg btn-wide">
          Share Now
        </a>
      </div>
    </div>
    <div class="col-span-12 md:col-span-6">
      <img :src="imageSrc" alt="Hero Image" class="w-full h-auto max-w-2xl rounded-2xl">
    </div>
    <div class="col-span-12">
      <div>Blazingly Fast</div>
      <div>No File Size Limit</div>
      <div></div>
      <div></div>
    </div>
  </div>
</template>
