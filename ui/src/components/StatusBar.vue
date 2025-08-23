<script setup>

import { VBtn } from 'vuetensils'

const props = defineProps({
  playerName: String,
  timeMs: Number,
  hideControls: Boolean,
  inProgress: { type: Boolean, required: true },
  hasDrawOffer: Boolean,
  disabled: Boolean,
  showMsBelow: { type: Number, default: 20000 }
})

const emit = defineEmits(['start', 'resign', 'offer-draw', 'accept-draw', 'decline-draw'])

function formatTime(ms) {
  if (ms == null || ms < 0) return '00:00'
  const totalSeconds = Math.floor(ms / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  const base = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
  if (ms < props.showMsBelow) {
    const millis = ms % 1000
    return `${base}.${String(millis).padStart(3, '0')}`
  }
  return base
}

</script>

<template>
  <div class="status-bar" role="group" aria-label="Player status bar">
    <div v-if="!inProgress" class="sb-row">
      <VBtn class="btn" :disabled="disabled" @click="emit('start')" aria-label="Start game">
        Start game
      </VBtn>
    </div>
    <div v-else class="sb-row">
      <span class="sb-name">{{ playerName }}</span>
      <span :class="{ urgent: !hideControls && (timeMs || 0) < showMsBelow }">
        {{ formatTime(timeMs) }}
      </span>

      <VBtn class="btn btn-ghost" v-if="!hideControls" :disabled="disabled" @click="emit('resign')">
        Resign
      </VBtn>

      <VBtn class="btn" v-if="!hideControls && !hasDrawOffer" :disabled="disabled" @click="emit('offer-draw')">
        Offer draw
      </VBtn>
      <VBtn class="btn btn-ghost urgent" v-else-if="!hideControls && hasDrawOffer" :disabled="disabled" @click="emit('accept-draw')">
        Accept draw
      </VBtn>
      <VBtn class="btn urgent" v-if="hasDrawOffer" :disabled="disabled" @click="emit('decline-draw')">
        Decline draw
      </VBtn>
    </div>
  </div>
</template>

<style scoped>

.status-bar {
  display: flex;
  justify-content: center;
  width: 100%;
}

.sb-name {
    font-weight: bold;
}

.sb-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 6px 0;
  font-size: 14px;
  line-height: 1.2;
}

.urgent { 
    animation: pulse 1s steps(2, end) infinite; 
}

@keyframes pulse { 
    50% { opacity: 0.55; } 
}

</style>

