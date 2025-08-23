<script setup>

import { ref } from 'vue';

import StatusBar from '@/components/StatusBar.vue';
import Dialog from '@/components/Dialog.vue';
import Board from '@/components/Board.vue';

import '@/styles.css';

const confirmResign = ref(false);
const confirmOffer = ref(false);
const startPrompt = ref(false);

const url = ref('');

const showResult = ref(false);
const result = ref({
    title: 'Victory',
    message: 'Congratulations, you won by checkmate!',
});

const inProgress = ref(false);

function resign() {
    console.error('Resign not implemented!');
}

function offerDraw() {
    console.error('Draw offer not implemented!');
}

function rematch() {
    console.error('Rmatch not implemented!');
}

function start() {
    console.error('Start not implemented!');
    inProgress.value = true;
    startPrompt.value = true;
}

async function copyUrl() {
    await navigator.clipboard.writeText(url.value);
}

</script>

<template>
  <div class="container">
    <StatusBar 
        isOpponent 
        playerName="Opponent" 
        :inProgress="inProgress" 
    />

    <Board :inProgress="inProgress" />

    <StatusBar 
        :timeMs="60000" 
        @resign="confirmResign = true"
        @offer-draw="confirmOffer = true"
        @start="start()"
        playerName="Me" 
        :inProgress="inProgress" 
    />

    <Dialog
        :open="confirmResign"
        title="Resign"
        message="Are you sure you would like to resign?"
        confirm-label="Resign"
        cancel-label="Cancel"
        @confirm="resign()"
        @cancel="confirmResign = false"
        @close="confirmResign = false" 
    />

    <Dialog
        :open="confirmOffer"
        title="Offer draw"
        message="Are you sure you would like to offer a draw?"
        confirm-label="Offer"
        cancel-label="Cancel"
        @confirm="offerDraw()"
        @cancel="confirmOffer = false"
        @close="confirmOffer = false" 
    />

    <Dialog
        :open="startPrompt"
        title="Start game"
        confirm-label="Copy link"
        cancel-label="Close"
        @confirm="copyUrl()"
        @cancel="startPrompt = false"
        @close="startPrompt = false">
        <div>
            Please share
            <a :href="url" role="button" @click.prevent="copyUrl()">this link</a>
            with your opponent.
        </div>
    </Dialog>

    <Dialog
        :open="showResult"
        :title="result.title"
        :message="result.message"
        confirm-label="Rematch"
        cancel-label="Close"
        @confirm="rematch()"
        @cancel="showResult = false"
        @close="showResult = false" 
    />
  </div>
</template>

<style scoped>

.container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 6px;
  height: 100vh;
}

</style>

