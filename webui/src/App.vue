<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div id="app-root" class="vh-100 vw-100 overflow-hidden">
    <RouterView />
  </div>
</template>

<style>
/* reset per occupare l'intero schermo senza scrollbar duplicate  */
html, body, #app, #app-root {
  height: 100%;
  width: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden;
}
</style>