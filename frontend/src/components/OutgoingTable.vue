<script setup lang="ts">
import FileCard from "@/components/FileCard.vue";
import {useFileStore} from "@/stores/file";

const file = useFileStore()
</script>

<template>
  <div class="card card-compact bg-base-200 dark:bg-base-300 mt-3 shadow-xl w-full h-full"
       v-if="file.getFileProgress().value != undefined && file.getFileProgress().value.size > 1">
    <div class="card-body">
      <div class="card-title">
        <p>Outgoing</p>
      </div>
      <template v-for="progress in file.getFileProgress().value.values()">
        <template :key="f.name" v-for="f in progress.files">
          <FileCard :file="f" :incoming="false" progress :target="progress.target"/>
        </template>
      </template>
    </div>
  </div>
</template>