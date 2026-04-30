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
        
        <!-- 好感度条 -->
        <div class="w-48 bg-black/50 p-2 rounded backdrop-blur-sm border border-gray-700">
          <div class="text-sm text-gray-300 mb-1 flex justify-between">
            <span>心意相通</span>
            <span>{{ gameStore.affection }}/30</span>
          </div>
          <ProgressBar :value="gameStore.affection" :max="30" color="#a855f7" />
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
        <button
          @click="openSaveSlots"
          :disabled="gameStore.isWaiting"
          class="bg-gray-800/80 hover:bg-gray-700 px-3 py-1.5 rounded text-sm border border-gray-600 backdrop-blur-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-gray-800/80"
        >
          保存
        </button>
        <button @click="goHome" class="bg-gray-800/80 hover:bg-gray-700 px-3 py-1.5 rounded text-sm border border-gray-600 backdrop-blur-sm transition-colors">退出</button>
      </div>
    </div>

    <!-- Save Slots -->
    <div
      v-if="showSaveSlots"
      class="absolute inset-0 z-40 flex items-center justify-center bg-black/55 px-4 backdrop-blur-sm pointer-events-auto"
      @click.self="closeSaveSlots"
    >
      <div class="w-full max-w-sm rounded-lg border border-purple-400/30 bg-[#090711]/95 p-5 shadow-[0_0_32px_rgba(168,85,247,0.22)]">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-base font-bold text-purple-100">选择保存栏位</h2>
          <button
            type="button"
            class="text-xl leading-none text-gray-500 transition-colors hover:text-white"
            aria-label="关闭"
            @click="closeSaveSlots"
          >
            ×
          </button>
        </div>

        <div class="grid gap-3">
          <button
            v-for="slotId in SAVE_SLOT_IDS"
            :key="slotId"
            type="button"
            class="rounded border border-gray-700 bg-black/45 px-4 py-3 text-left transition-colors hover:border-purple-400/70 hover:bg-purple-950/35"
            @click="saveToSlot(slotId)"
          >
            <span class="block text-sm font-bold text-gray-100">栏位 {{ slotId }}</span>
            <span class="mt-1 block text-xs text-gray-400">{{ getSlotStatus(slotId) }}</span>
          </button>
        </div>
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
            <div class="font-pixel text-purple-400 text-lg mb-3 flex items-center gap-3">
              <span>{{ latestMessage.role === 'assistant' ? (gameStore.isEnding ? '旁白' : '神秘女孩') : '你' }}</span>
              
              <!-- 情绪指示器 -->
              <transition name="fade" mode="out-in">
                <span v-if="gameStore.currentEmotion && latestMessage.role === 'assistant'" :class="emotionStyle" class="text-sm px-3 py-0.5 rounded-full animate-pulse">
                  {{ emotionLabel }}
                </span>
              </transition>
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
import { SaveSystem, type SaveSlot } from '@/modules/SaveSystem'
import { AchievementTracker } from '@/modules/AchievementTracker'
import TypewriterText from '@/components/TypewriterText.vue'
import ProgressBar from '@/components/ProgressBar.vue'

const router = useRouter()
const gameStore = useGameStore()
const SAVE_SLOT_IDS = [1, 2, 3] as const

const inputText = ref('')
const textCompleted = ref(false)
const latestHint = ref<string | null>(null)
const showSaveSlots = ref(false)
const saveSlots = ref<SaveSlot[]>([])

const latestMessage = computed(() => {
  if (gameStore.messages.length === 0) return null
  return gameStore.messages[gameStore.messages.length - 1]
})

const saveSlotMap = computed(() => new Map(saveSlots.value.map((slot) => [slot.id, slot])))

const currentBg = computed(() => {
  if (gameStore.isEnding) {
    if (gameStore.endingType === 'end_death') return '/assets/images/cg_end_fall.png'
    if (gameStore.endingType === 'end_disappear') return '/assets/images/cg_end_disappear.png'
    if (gameStore.endingType === 'end_true_release') return '/assets/images/cg_acquaintance_16_9.webp' // Reusing this image for the true ending fade out
    if (gameStore.endingType === 'end_self_deception') return '/assets/images/char_girl_sad.png'
    if (gameStore.endingType === 'end_eternal_cage') return '/assets/images/char_girl_smoke.png'
  }
  
  // Use the new full-scene character images based on the state
  if (gameStore.roundCount > 7) return '/assets/images/char_girl_smoke.png'
  if (gameStore.roundCount > 3) return '/assets/images/char_girl_normal.png'
  return '/assets/images/char_girl_sad.png'
})

// 情绪标签和样式
const emotionLabel = computed(() => {
  switch (gameStore.currentEmotion) {
    case '刺痛': return '💔 心中一刺'
    case '惊讶': return '✨ 微微惊讶'
    case '柔软': return '💕 心中一软'
    case '好奇': return '🤔 产生好奇'
    default: return ''
  }
})

const emotionStyle = computed(() => {
  switch (gameStore.currentEmotion) {
    case '刺痛': return 'bg-red-900/60 text-red-200 border border-red-500/40'
    case '惊讶': return 'bg-yellow-900/60 text-yellow-200 border border-yellow-500/40'
    case '柔软': return 'bg-pink-900/60 text-pink-200 border border-pink-500/40'
    case '好奇': return 'bg-blue-900/60 text-blue-200 border border-blue-500/40'
    default: return ''
  }
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

const formatSaveTime = (timestamp: number) =>
  new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(timestamp))

const getSlotStatus = (slotId: number) => {
  const slot = saveSlotMap.value.get(slotId)
  return slot ? `已有存档：${formatSaveTime(slot.timestamp)}` : '空栏位'
}

const refreshSaveSlots = () => {
  saveSlots.value = SaveSystem.getSlots()
}

const openSaveSlots = () => {
  audioManager.playSfx('click')
  if (gameStore.isWaiting) {
    alert('正在等待回应，暂时不能保存。')
    return
  }

  refreshSaveSlots()
  showSaveSlots.value = true
}

const closeSaveSlots = () => {
  showSaveSlots.value = false
}

const saveToSlot = (slotId: number) => {
  audioManager.playSfx('click')
  const existingSlot = saveSlotMap.value.get(slotId)
  if (existingSlot && !confirm(`栏位 ${slotId} 已有存档，是否覆盖？`)) return

  if (SaveSystem.save(slotId)) {
    refreshSaveSlots()
    showSaveSlots.value = false
    alert(`游戏已保存至栏位 ${slotId}`)
  } else {
    alert(gameStore.isWaiting ? '正在等待回应，暂时不能保存。' : '保存失败')
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
