<script setup>

import { TheChessboard } from 'vue3-chessboard';
import 'vue3-chessboard/style.css';

</script>

<template>
  <div class="board-container">
    <TheChessboard />
  </div>
</template>

<style scoped>
.board-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100vh;
}
</style>

