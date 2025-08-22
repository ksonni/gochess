
<script setup>

import { TheChessboard } from 'vue3-chessboard';
import 'vue3-chessboard/style.css';

</script>

<template>
  <TheChessboard class="board" />
</template>

<style scoped>

.board {
  border: 2px solid #000;
  border-radius: 4px;
  max-width: 80vh;
}

</style>
