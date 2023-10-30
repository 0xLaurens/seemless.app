<script setup lang="ts">
let imageSrc =
    'https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg/440px-Fig_7_Le_Telephone_by_T_du_Moncel_Paris_1880_%28Large%29.jpg'
</script>

<template>
  <section class="flex h-screen items-center justify-center py-10">
    <!--    <div class="container h-full px-4 md:px-6 lg:flex lg:items-center lg:justify-between">-->
    <div
        class="container flex flex-col lg:flex-row items-center max-w-6xl px-4 md:px-6 lg:items-center justify-between h-full md:h-auto md:justify-start w-full">
      <div class="lg:w-1/2 space-y-4">
        <h1 class="text-4xl text-black dark:text-white font-bold tracking-tighter sm:text-5xl xl:text-6xl"> Share Files
          <span class="text-accent">Easily</span> and <span class="text-accent">Securely</span></h1>
        <p class="text-lg">With our platform, you can share files with anyone, anywhere, at
          any
          time in
          just a few clicks.
        </p>
        <router-link class="btn btn-md font-normal btn-accent normal-case" to="/nick">Start Sharing</router-link>
      </div>
      <div class="lg:w-1/2">
        <img :src="imageSrc"
             class="mx-auto aspect-content rounded-xl object-cover sm:w-full lg:order-last"
             alt="tin can telephone"/>
      </div>
    </div>
  </section>
</template>
