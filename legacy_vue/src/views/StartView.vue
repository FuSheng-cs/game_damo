<template>
  <div class="start-view min-h-screen w-full flex flex-col items-center justify-center relative overflow-hidden font-sans">
    <!-- Background Image with slow pan -->
    <div 
      class="absolute inset-0 bg-cover bg-center bg-no-repeat scale-105 animate-slow-pan z-0"
      style="background-image: url('/assets/images/bg_rooftop_night_1920.webp');"
    ></div>
    
    <!-- Overlay Gradient for text readability -->
    <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-black/60 to-black/95 z-0 pointer-events-none"></div>

    <!-- Rain VFX layer -->
    <div class="absolute inset-0 bg-[url('/assets/images/vfx_rain_sprite.webp')] opacity-20 mix-blend-screen pointer-events-none z-0 bg-repeat animation-rain"></div>
    
    <div class="z-10 relative w-full max-w-6xl px-6 flex flex-col md:flex-row items-center justify-between gap-12 mt-8 md:mt-16">
      
      <!-- Left Column: Title and Intro -->
      <div class="flex-1 text-left flex flex-col items-center md:items-start text-center md:text-left">
        <h1 class="text-6xl md:text-8xl font-black text-transparent bg-clip-text bg-gradient-to-b from-purple-100 to-purple-500 mb-6 tracking-widest title-glow" style="font-family: 'STZhongsong', 'SimSun', serif;">
          天台十句
        </h1>
        <div class="h-1 w-24 bg-purple-500 mb-8 rounded-full shadow-[0_0_10px_rgba(168,85,247,0.8)] hidden md:block"></div>
        <p class="text-gray-300 mb-10 text-lg md:text-2xl font-light italic leading-relaxed text-shadow-sm">
          "你只有 10 句话的时间。<br/>
          <span class="text-purple-300 font-medium">但也许，我们都需要更多的时间来原谅自己。</span>"
        </p>

        <!-- Game Intro Panel -->
        <div class="w-full max-w-lg bg-black/40 p-6 md:p-8 rounded-2xl border border-purple-500/30 backdrop-blur-md text-left shadow-2xl transition-all duration-500 hover:border-purple-400/60 hover:bg-black/50 hover:shadow-[0_0_40px_rgba(168,85,247,0.25)]">
          <h2 class="text-xl font-bold text-purple-300 mb-4 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            关于游戏
          </h2>
          <div class="text-gray-300 text-sm md:text-base space-y-4 leading-relaxed">
            <p>在这个霓虹闪烁的深夜，你来到了天台。坐在围栏边缘的是"艾"——一个极度厌世、带着颓废破碎感的女孩。</p>
            <p><strong>你的目标：</strong>你只有 <span class="text-purple-400 font-bold text-lg">10 句话</span> 的机会与她交谈。试着走进她的内心，把她从边缘拉回来。</p>
            <ul class="list-none space-y-2 mt-3 text-gray-400 bg-black/20 p-4 rounded-xl border border-white/5">
              <li class="flex items-start gap-3">
                <span class="text-purple-500 mt-1 text-xs">◆</span> 
                <span>倾听大于说教，普通的套话只会让她离你远去。</span>
              </li>
              <li class="flex items-start gap-3">
                <span class="text-purple-500 mt-1 text-xs">◆</span> 
                <span>你的每一次尝试都会在她的潜意识中留下痕迹。</span>
              </li>
              <li class="flex items-start gap-3">
                <span class="text-purple-500 mt-1 text-xs">◆</span> 
                <span>建议在设置中配置 <strong class="text-purple-300 mx-1">千问 API Key</strong> 以获得完整 AI 对话体验（支持千问/豆包/OpenAI）。</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Right Column: Main Actions -->
      <div class="flex flex-col gap-5 w-full max-w-xs shrink-0 mt-8 md:mt-0">
        <button @click="startGame" class="btn-glass primary-btn group">
          <span class="relative z-10 flex items-center justify-center gap-2 text-xl font-bold tracking-[0.2em]">
            开始游戏
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 group-hover:translate-x-1.5 transition-transform duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </span>
          <div class="absolute inset-0 bg-gradient-to-r from-purple-600/80 to-indigo-600/80 opacity-0 group-hover:opacity-100 transition-opacity duration-300 -z-0"></div>
        </button>
        
        <button @click="loadGame" class="btn-glass group">
          <span class="relative z-10 tracking-[0.2em] font-medium text-lg">读取存档</span>
          <div class="absolute inset-0 bg-white/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 -z-0"></div>
        </button>
        
        <button @click="goToAchievements" class="btn-glass group">
          <span class="relative z-10 tracking-[0.2em] font-medium text-lg">成就图鉴</span>
          <div class="absolute inset-0 bg-white/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 -z-0"></div>
        </button>
        
        <button @click="goToSettings" class="btn-glass group">
          <span class="relative z-10 tracking-[0.2em] font-medium text-lg">游戏设置</span>
          <div class="absolute inset-0 bg-white/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 -z-0"></div>
        </button>
      </div>

    </div>
    
    <!-- Footer -->
    <div class="absolute bottom-6 text-xs text-gray-500/70 z-10 tracking-[0.15em] uppercase flex items-center gap-4 font-mono">
      <span>v1.1.0</span>
      <span class="w-1 h-1 rounded-full bg-gray-500/70"></span>
      <span>AI-Driven Narrative Experience</span>
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

.title-glow {
  text-shadow: 0 0 20px rgba(168, 85, 247, 0.5), 0 0 40px rgba(168, 85, 247, 0.3), 0 4px 8px rgba(0,0,0,0.8);
}

.text-shadow-sm {
  text-shadow: 0 2px 4px rgba(0,0,0,0.9);
}

.btn-glass {
  @apply relative overflow-hidden py-4 px-8 bg-black/40 border border-gray-500/30 text-gray-200 rounded-xl transition-all duration-300 backdrop-blur-md shadow-lg hover:shadow-[0_8px_25px_rgba(0,0,0,0.5)] hover:border-gray-300/50 hover:-translate-y-1 flex items-center justify-center;
}

.primary-btn {
  @apply bg-purple-900/50 border-purple-500/60 text-white shadow-[0_0_20px_rgba(168,85,247,0.4)] hover:shadow-[0_0_35px_rgba(168,85,247,0.6)] hover:border-purple-300;
}

@keyframes slow-pan {
  0% { transform: scale(1.05) translate(0, 0); }
  50% { transform: scale(1.1) translate(-1%, -1%); }
  100% { transform: scale(1.05) translate(0, 0); }
}

.animate-slow-pan {
  animation: slow-pan 40s ease-in-out infinite;
}

@keyframes rain-fall {
  0% { background-position: 0 0; }
  100% { background-position: 20% 100%; }
}

.animation-rain {
  animation: rain-fall 2s linear infinite;
  background-size: 400px 400px;
}
</style>