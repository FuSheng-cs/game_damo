<template>
  <div class="typewriter-text" @click="skip">
    <span :class="{'typewriter-cursor': isTyping}">{{ displayedText }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  text: string;
  speed?: number; // ms per char
}>()

const emit = defineEmits<{
  (e: 'complete'): void
}>()

const displayedText = ref('')
const isTyping = ref(false)
let timer: ReturnType<typeof setInterval> | null = null
let currentIndex = 0

const startTyping = () => {
  stopTyping()
  displayedText.value = ''
  currentIndex = 0
  isTyping.value = true
  
  timer = setInterval(() => {
    if (currentIndex < props.text.length) {
      displayedText.value += props.text[currentIndex]
      currentIndex++
    } else {
      completeTyping()
    }
  }, props.speed || 50)
}

const stopTyping = () => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

const completeTyping = () => {
  stopTyping()
  isTyping.value = false
  displayedText.value = props.text
  emit('complete')
}

const skip = () => {
  if (isTyping.value) {
    completeTyping()
  }
}

watch(() => props.text, () => {
  startTyping()
})

onMounted(() => {
  startTyping()
})

onUnmounted(() => {
  stopTyping()
})
</script>

<style scoped>
.typewriter-text {
  cursor: pointer;
  white-space: pre-wrap;
}
</style>
