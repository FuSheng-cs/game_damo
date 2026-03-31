<template>
  <div class="achievements-view min-h-screen bg-gray-900 p-8">
    <div class="max-w-3xl mx-auto">
      <div class="flex justify-between items-end border-b border-gray-700 pb-4 mb-8">
        <h2 class="text-3xl font-pixel text-purple-400">成就图鉴</h2>
        <div class="text-gray-400">完成度：{{ progress.unlocked }} / {{ progress.total }}</div>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div 
          v-for="ach in ACHIEVEMENTS" 
          :key="ach.id"
          class="p-4 rounded-lg border flex items-start gap-4 transition-all duration-300"
          :class="isUnlocked(ach.id) ? 'bg-gray-800 border-purple-500/50' : 'bg-gray-900 border-gray-800 opacity-50 grayscale'"
        >
          <div class="w-12 h-12 bg-gray-700 rounded flex items-center justify-center text-2xl">
            <span v-if="isUnlocked(ach.id)">🏆</span>
            <span v-else>🔒</span>
          </div>
          <div>
            <h4 class="font-bold text-lg" :class="isUnlocked(ach.id) ? 'text-purple-300' : 'text-gray-500'">
              {{ isUnlocked(ach.id) ? ach.name : '???' }}
            </h4>
            <p class="text-sm text-gray-400 mt-1">
              {{ isUnlocked(ach.id) ? ach.description : '尚未解锁此成就' }}
            </p>
          </div>
        </div>
      </div>
      
      <div class="mt-10 flex justify-center">
        <button @click="goBack" class="py-2 px-8 bg-gray-800 hover:bg-gray-700 rounded border border-gray-700">返回标题</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { AchievementTracker, ACHIEVEMENTS } from '@/modules/AchievementTracker'
import { audioManager } from '@/modules/AudioManager'

const router = useRouter()
const unlockedIds = ref<string[]>([])
const progress = ref({ unlocked: 0, total: 0 })

onMounted(() => {
  unlockedIds.value = AchievementTracker.getUnlocked()
  progress.value = AchievementTracker.getProgress()
})

const isUnlocked = (id: string) => unlockedIds.value.includes(id)

const goBack = () => {
  audioManager.playSfx('click')
  router.push('/')
}
</script>
