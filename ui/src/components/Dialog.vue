<script setup>

import { ref, watch, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  open: { type: Boolean, required: true },
  title: { type: String, default: '' },
  message: { type: String, default: '' },
  confirmLabel: { type: String },
  cancelLabel: { type: String },
  disabled: Boolean
})

const emit = defineEmits(['confirm', 'cancel', 'close'])

const root = ref(null)

function onKeydown(e) {
  if (e.key === 'Escape') {
    e.stopPropagation()
    emit('cancel')
    emit('close')
  }
}

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && root.value) {
      root.value.focus()
    }
  },
  { flush: 'post' }
)

onMounted(() => {
  window.addEventListener('keydown', onKeydown, true)
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown, true)
})

function onBackdropClick(e) {
  if (e.target === e.currentTarget) {
    emit('cancel')
    emit('close')
  }
}
</script>

<template>
  <teleport to="body">
    <div
      v-if="open"
      class="gd-backdrop"
      @click="onBackdropClick"
      aria-hidden="false">
      <div
        class="gd-panel"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="title ? 'gd-title' : undefined"
        :aria-describedby="message ? 'gd-desc' : undefined"
        tabindex="-1"
        ref="root">
        <header v-if="title" class="gd-header">
          <h2 id="gd-title">{{ title }}</h2>
        </header>

        <section class="gd-body">
          <p id="gd-desc">{{ message }}</p>
          <slot />
        </section>

        <footer class="gd-actions">
          <button v-if="cancelLabel" type="button" class="btn btn-ghost" :disabled="disabled" @click="() => { emit('cancel'); emit('close') }">
            {{ cancelLabel }}
          </button>

          <button v-if="confirmLabel" type="button" class="btn" :disabled="disabled" @click="() => { emit('confirm'); emit('close') }">
            {{ confirmLabel }}
          </button>
        </footer>
      </div>
    </div>
  </teleport>
</template>

<style scoped>

.gd-backdrop {
  position: fixed;
  inset: 0;
  display: grid;
  place-items: center;
  background: rgba(0, 0, 0, 0.35);
  z-index: 1000;
}

.gd-panel {
  min-width: 280px;
  max-width: 90vw;
  padding: 16px;
  background: #fff;
  color: #111;
  border: 1px solid #222;
  border-radius: 8px;
  outline: none;
  box-shadow: 0 8px 24px rgba(0,0,0,0.18);
}

.gd-header { margin-bottom: 8px; }
.gd-header h2 {
  margin: 0;
  font-size: 14px;
  line-height: 1.3;
  font-weight: 600;
}

.gd-body { 
    font-size: 14px;
    margin-bottom: 12px; 
}
.gd-body p { margin: 0; }

.gd-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

</style>
