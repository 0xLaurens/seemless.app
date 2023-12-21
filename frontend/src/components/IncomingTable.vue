<script setup lang="ts">
import FileCard from "@/components/FileCard.vue";
import {useFileStore} from "@/stores/file";
import {useDownloadStore} from "@/stores/download";

const file = useFileStore()
</script>

<template>
  <div class="card card-compact bg-base-200 dark:bg-base-300 mt-3 shadow-sm w-full h-full"
       v-if="file.getReceiveOffer().value != undefined">
    <div class="card-body">
      <div class="card-title">
        <p>Incoming</p>
      </div>
      <template v-for="f in file.getReceiveOffer().value?.files" :key="f.name">
        <FileCard :file="f" :incoming="true" progress :from="file.getReceiveOffer().value?.from"/>
      </template>
    </div>
  </div>
</template>
