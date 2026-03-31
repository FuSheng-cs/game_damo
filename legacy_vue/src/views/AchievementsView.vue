<template>
  <div class="achievements-view min-h-screen bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-gray-900 via-gray-950 to-black p-8 relative overflow-hidden">
    <!-- Glow effects -->
    <div class="absolute top-[-20%] left-[20%] w-96 h-96 bg-purple-900/10 rounded-full blur-[100px] pointer-events-none"></div>
    <div class="max-w-3xl mx-auto relative z-10">
      <div class="flex justify-between items-end border-b border-gray-700/80 pb-4 mb-8">
        <h2 class="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 tracking-wide">成就图鉴</h2>
        <div class="text-gray-400 font-mono text-sm">完成度：{{ progress.unlocked }} / {{ progress.total }}</div>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div 
          v-for="ach in ACHIEVEMENTS" 
          :key="ach.id"
          class="p-5 rounded-2xl border flex items-start gap-4 transition-all duration-300 backdrop-blur-sm shadow-lg hover:-translate-y-1"
          :class="isUnlocked(ach.id) ? 'bg-gray-800/80 border-purple-500/40 hover:shadow-[0_0_20px_rgba(168,85,247,0.3)]' : 'bg-gray-900/50 border-gray-800 opacity-60 grayscale'"
        >
          <div class="w-14 h-14 bg-gray-900/80 rounded-xl flex items-center justify-center text-3xl shadow-inner border border-gray-700/50 flex-shrink-0">
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
      
      <div class="mt-12 flex justify-center">
        <button @click="goBack" class="py-2.5 px-8 bg-gradient-to-r from-gray-700 to-gray-600 hover:from-gray-600 hover:to-gray-500 rounded-xl transition-all duration-300 shadow-lg text-white font-medium hover:-translate-y-0.5 border border-gray-600">返回标题</button>
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
