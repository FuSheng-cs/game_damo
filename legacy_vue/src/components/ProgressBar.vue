<template>
  <div class="progress-bar-container w-full h-4 bg-gray-800 rounded-full overflow-hidden border border-gray-600">
    <div 
      class="progress-fill h-full transition-all duration-500 ease-in-out"
      :style="{ width: `${percentage}%`, backgroundColor: color }"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  value: number;
  max: number;
  color?: string;
}>()

const percentage = computed(() => {
  return Math.min(Math.max((props.value / props.max) * 100, 0), 100)
})

const color = computed(() => props.color || 'var(--theme-color)')
</script>
