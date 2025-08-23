<script setup>

import { defineProps, reactive, watch } from 'vue';

import { TheChessboard } from 'vue3-chessboard';
import 'vue3-chessboard/style.css';

const props = defineProps({
    inProgress: { type: Boolean },
})

const boardConfig = reactive({
    viewOnly: true,
    coordinates: true,
})

let boardAPI;

// The reactive API of the board is weird
watch(
  () => props.inProgress,
  (inProgress) => {
    boardConfig.viewOnly = !inProgress
  },
  { immediate: true },
)

</script>

<template>
  <TheChessboard 
    :key="props.inProgress" 
    :board-config="boardConfig" 
    @board-created="(api) => (boardAPI = api)"
    reactive-config 
    class="board" 
  />
</template>

<style scoped>

.board {
  border: 2px solid #000;
  border-radius: 4px;
  max-width: 80vh;
}

</style>
