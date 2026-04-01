<template>
  <div class="game-view min-h-screen w-full relative flex flex-col bg-black text-gray-100 overflow-hidden">
    <!-- Background (now using the character images directly as they contain the full scene) -->
    <div class="absolute inset-0 z-0">
      <img 
        :src="currentBg" 
        class="w-full h-full object-cover transition-opacity duration-1000"
        :class="gameStore.isEnding ? 'opacity-100' : 'opacity-80'"
        alt="Background"
      />
    </div>

    <!-- Character layer removed since the new images are full-scene compositions -->

    <!-- Top UI -->
    <div class="absolute top-0 left-0 right-0 z-30 p-6 flex justify-between items-start pointer-events-auto">
      <div class="flex gap-4">
        <div class="w-48 bg-black/50 p-2 rounded backdrop-blur-sm border border-gray-700">
          <div class="text-sm text-gray-300 mb-1 flex justify-between">
            <span>剩余机会</span>
            <span>{{ gameStore.roundCount }}/10</span>
          </div>
          <ProgressBar :value="gameStore.roundCount" :max="10" color="#ef4444" />
        </div>
        
        <div v-if="gameStore.hintCount > 0 && !gameStore.isEnding" class="bg-[#0a0515]/90 px-4 py-2 rounded-xl border border-gray-800 shadow-lg backdrop-blur-md flex items-center">
          <button 
            @click="handleHint" 
            :disabled="gameStore.isWaiting"
            class="text-sm text-yellow-500 hover:text-yellow-400 font-bold disabled:opacity-50 flex items-center gap-2 transition-colors"
          >
            <span class="text-lg">💡</span> 寻找线索 ({{ gameStore.hintCount }})
          </button>
        </div>
      </div>
      
      <div class="flex gap-3">
        <button @click="saveGame" class="bg-gray-800/80 hover:bg-gray-700 px-3 py-1.5 rounded text-sm border border-gray-600 backdrop-blur-sm transition-colors">保存</button>
        <button @click="goHome" class="bg-gray-800/80 hover:bg-gray-700 px-3 py-1.5 rounded text-sm border border-gray-600 backdrop-blur-sm transition-colors">退出</button>
      </div>
    </div>

    <!-- Dialog Box -->
    <div class="absolute bottom-0 left-0 right-0 z-30 p-4 md:p-8 pointer-events-auto bg-gradient-to-t from-black via-black/80 to-transparent pt-24">
      <div class="max-w-4xl mx-auto">
        <div class="bg-black/60 border border-gray-700/50 rounded-xl p-6 shadow-2xl backdrop-blur-md min-h-[160px] flex flex-col justify-end transition-all duration-500">
          
          <div v-if="gameStore.isWaiting" class="italic text-gray-500 text-base md:text-lg py-4">
            <TypewriterText :text="gameStore.waitingText" />
          </div>
          
          <div v-else-if="latestHint" class="text-base md:text-lg leading-relaxed text-yellow-100/90 mb-6 bg-yellow-900/20 p-4 rounded-xl border border-yellow-700/30">
            <div class="font-bold text-yellow-500 text-sm mb-2 flex justify-between">
              <span>【内心直觉】</span>
              <button @click="clearHint" class="text-gray-500 hover:text-white">✕</button>
            </div>
            <TypewriterText :text="latestHint" @complete="onTextComplete" :key="latestHint" />
          </div>
          
          <div v-else-if="latestMessage" class="text-lg md:text-xl leading-relaxed text-gray-100">
            <div class="font-pixel text-purple-400 text-lg mb-3">
              {{ latestMessage.role === 'assistant' ? (gameStore.isEnding ? '旁白' : '神秘女孩') : '你' }}
            </div>
            <TypewriterText :text="latestMessage.content" @complete="onTextComplete" :key="gameStore.messages.length" />
          </div>

          <!-- Input Area -->
          <div v-if="!gameStore.isWaiting && !gameStore.isEnding && textCompleted" class="mt-6 flex gap-4">
            <input 
              v-model="inputText" 
              @keyup.enter="handleSend"
              type="text" 
              class="flex-1 bg-gray-900/80 border border-gray-600 rounded-lg px-4 py-3 text-gray-100 focus:outline-none focus:border-purple-500 transition-colors shadow-inner"
              placeholder="对她说点什么..."
              autofocus
            />
            <button @click="handleSend" class="px-6 py-3 bg-purple-600 hover:bg-purple-500 text-white rounded-lg font-bold transition-all shadow-[0_0_15px_rgba(168,85,247,0.4)] hover:shadow-[0_0_20px_rgba(168,85,247,0.6)]">
              发送
            </button>
          </div>

          <div v-if="gameStore.isEnding && textCompleted" class="mt-6 flex justify-center gap-4">
            <button v-if="gameStore.endingType === 'end_acquaintance'" @click="goToChatAfter" class="px-8 py-3 bg-[#07c160] hover:bg-[#06ad56] text-white font-bold rounded-lg transition-colors shadow-[0_0_15px_rgba(7,193,96,0.4)]">
              添加联系人...
            </button>
            <button v-else @click="goHome" class="px-8 py-3 bg-gray-800 hover:bg-gray-700 text-white rounded-lg border border-gray-600 transition-colors">
              返回标题
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '@/store/gameStore'
import { audioManager } from '@/modules/AudioManager'
import { SaveSystem } from '@/modules/SaveSystem'
import { AchievementTracker } from '@/modules/AchievementTracker'
import TypewriterText from '@/components/TypewriterText.vue'
import ProgressBar from '@/components/ProgressBar.vue'

const router = useRouter()
const gameStore = useGameStore()

const inputText = ref('')
const textCompleted = ref(false)
const latestHint = ref<string | null>(null)

const latestMessage = computed(() => {
  if (gameStore.messages.length === 0) return null
  return gameStore.messages[gameStore.messages.length - 1]
})

const currentBg = computed(() => {
  if (gameStore.isEnding) {
    if (gameStore.endingType === 'end_death') return '/assets/images/cg_death_falling_16_9.webp'
    if (gameStore.endingType === 'end_true_release') return '/assets/images/cg_acquaintance_16_9.webp' // Reusing this image for the true ending fade out
    if (gameStore.endingType === 'end_self_deception') return '/assets/images/char_girl_sad.png'
    if (gameStore.endingType === 'end_eternal_cage') return '/assets/images/char_girl_smoke.png'
  }
  
  // Use the new full-scene character images based on the state
  if (gameStore.roundCount > 7) return '/assets/images/char_girl_smoke.png'
  if (gameStore.roundCount > 3) return '/assets/images/char_girl_normal.png'
  return '/assets/images/char_girl_sad.png'
})

const onTextComplete = () => {
  textCompleted.value = true
  if (gameStore.isEnding && gameStore.endingType) {
    AchievementTracker.unlock(gameStore.endingType)
  }
}

const handleSend = async () => {
  const text = inputText.value.trim()
  if (!text || gameStore.isWaiting || gameStore.isEnding) return
  
  audioManager.playSfx('click')
  inputText.value = ''
  textCompleted.value = false
  latestHint.value = null
  
  await gameStore.sendMessage(text)
}

const handleHint = async () => {
  if (gameStore.hintCount <= 0 || gameStore.isWaiting) return
  audioManager.playSfx('click')
  textCompleted.value = false
  const hint = await gameStore.requestHint()
  if (hint) {
    latestHint.value = hint
  }
}

const clearHint = () => {
  audioManager.playSfx('click')
  latestHint.value = null
  textCompleted.value = true // Ensure input shows up immediately after closing hint
}

const saveGame = () => {
  audioManager.playSfx('click')
  if (SaveSystem.save(1)) {
    alert('游戏已保存至栏位 1')
  } else {
    alert('保存失败')
  }
}

const goHome = () => {
  audioManager.playSfx('click')
  router.push('/')
}

const goToChatAfter = () => {
  audioManager.playSfx('click')
  router.push('/chat-after')
}

onMounted(() => {
  AchievementTracker.unlock('first_try')
  audioManager.playBgm('/assets/audio/bgm_rooftop.mp3')
})
</script>
