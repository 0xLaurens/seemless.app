<script setup lang="ts">
</script>

<template>
  <div class="card card-normal bg-base-200 max-w-2xl dark:bg-base-300 self-center">
    <div class="card-body flex">
      <div class="card-title">
        <p>About</p>
      </div>
      <div>
        <p>Seemless.app is a cross platform app for transferring files between devices.
          You might not always have a
          cable or external
          storage available to exchange files from one device to another.</p>
      </div>
      <div>
        <p class="font-bold">How it works.</p>
        <p>Seemless.app helps users connect to each other directly using WebRTC. The server has no knowledge of the
          files being sent due to the connection being peer-to-peer. </p>
      </div>
      <div>
        <p class="font-bold">How to use.</p>
        <div class="px-8">
          <ol class="list-decimal">
            <li>Open Seemless.app site on one device (you are on it right now!)</li>
            <li>If you are on the same wifi network open the Seemless.app site on the other device. Scan the QR code or
              use the join URL if the user is on another network.
            </li>
            <li>Select the devices you want to transfer files to</li>
            <li>Press the transfer button</li>
            <li>Accept the files on the other device.</li>
            <li>Keep your devices open until the transfer is complete</li>
            <li>Success!</li>
          </ol>
        </div>
      </div>
      <div>
        <p class="font-bold">Troubleshooting</p>
        <div class="px-8">
          <ol class="list-decimal">
            <li>Chrome and Firefox browsers are recommended, other browsers might not support all the right web api's
            </li>
            <li>Install the latest version of your browser</li>
            <li>Make sure you are connected on the same network or are connected to the same public room</li>
            <li>Reload the page if you don't see other devices</li>
          </ol>
        </div>
      </div>
    </div>
  </div>
</template>