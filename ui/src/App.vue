<script setup>

import StatusBar from '@/components/StatusBar.vue';
import Dialog from '@/components/Dialog.vue';
import Board from '@/components/Board.vue';

import '@/styles.css';

</script>

<template>
  <div class="container">
    <StatusBar hideControls playerName="Opponent" inProgress />

    <Board />

    <StatusBar :timeMs="60000" playerName="Me" inProgress />
  </div>
</template>

<style scoped>

.container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 12px;
  height: 100vh;
}

</style>

