<template>
  <div class="start-view min-h-screen w-full flex flex-col items-center justify-center relative overflow-hidden">
    <!-- Mock Background -->
    <div class="absolute inset-0 bg-gradient-to-b from-gray-900 to-black z-0"></div>
    
    <div class="z-10 text-center relative w-full max-w-4xl px-4">
      <h1 class="text-6xl font-pixel text-purple-400 mb-6 tracking-wider drop-shadow-[0_0_10px_rgba(168,85,247,0.8)]">DAMO</h1>
      <p class="text-gray-400 mb-12 text-lg italic">"你只有 10 句话的时间。但也许，我们都需要更多的时间来原谅自己。"</p>
      
      <div class="flex flex-col md:flex-row gap-8 items-center justify-center">
        <!-- Main Actions -->
        <div class="flex flex-col gap-4 w-64">
          <button @click="startGame" class="btn-primary">开始游戏</button>
          <button @click="loadGame" class="btn-secondary">读取存档</button>
          <button @click="goToAchievements" class="btn-secondary">成就图鉴</button>
          <button @click="goToSettings" class="btn-secondary">游戏设置</button>
        </div>
        
        <!-- Game Intro Panel -->
        <div class="w-full md:w-96 bg-gray-800/50 p-6 rounded-lg border border-gray-700 backdrop-blur-sm text-left">
          <h2 class="text-xl font-pixel text-purple-300 mb-4 border-b border-gray-600 pb-2">关于游戏</h2>
          <div class="text-gray-300 text-sm space-y-3 leading-relaxed">
            <p>在这个霓虹闪烁的深夜，你来到了天台。坐在围栏边缘的是"艾"——一个极度厌世、带着颓废破碎感的女孩。</p>
            <p><strong>目标：</strong>你只有 <span class="text-red-400 font-bold">10 句话</span> 的机会与她交谈。试着走进她的内心，把她从边缘拉回来。</p>
            <p><strong>注意：</strong></p>
            <ul class="list-disc pl-5 text-gray-400 space-y-1">
              <li>倾听大于说教，普通的套话只会让她离你远去。</li>
              <li>你的每一次尝试都会在她的潜意识中留下痕迹。</li>
              <li>在"游戏设置"中配置 <span class="text-purple-400">AI API Key</span> 以获得完整的 AI 对话体验（支持 千问/豆包/OpenAI）。</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
    
    <div class="absolute bottom-4 text-xs text-gray-600 z-10">
      v1.0.0 | AI-Driven Narrative Experience
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useGameStore } from '@/store/gameStore'
import { audioManager } from '@/modules/AudioManager'

const router = useRouter()
const gameStore = useGameStore()

const startGame = () => {
  audioManager.playSfx('click')
  gameStore.resetGame()
  router.push('/game')
}

const loadGame = () => {
  audioManager.playSfx('click')
  // For simplicity, we could implement a load menu here, or just load slot 1
  // We will load slot 1 if it exists, otherwise alert
  import('@/modules/SaveSystem').then(({ SaveSystem }) => {
    if (SaveSystem.load(1)) {
      router.push('/game')
    } else {
      alert('未找到有效存档。')
    }
  })
}

const goToSettings = () => {
  audioManager.playSfx('click')
  router.push('/settings')
}

const goToAchievements = () => {
  audioManager.playSfx('click')
  router.push('/achievements')
}
</script>

<style scoped>
@reference "tailwindcss";
.btn-primary {
  @apply py-3 px-6 bg-purple-600 hover:bg-purple-500 text-white rounded font-bold transition-colors shadow-[0_0_15px_rgba(168,85,247,0.5)];
}
.btn-secondary {
  @apply py-3 px-6 bg-gray-800 hover:bg-gray-700 text-gray-200 border border-gray-600 rounded transition-colors;
}
</style>