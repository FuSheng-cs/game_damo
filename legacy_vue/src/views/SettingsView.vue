<template>
  <div class="settings-view min-h-screen bg-gray-900 p-8">
    <div class="max-w-2xl mx-auto bg-gray-800 p-8 rounded-lg shadow-lg">
      <h2 class="text-3xl font-pixel text-purple-400 mb-8">设置 (Settings)</h2>
      
      <div class="space-y-6">
        <!-- API Settings -->
        <div class="space-y-4">
          <h3 class="text-xl border-b border-gray-700 pb-2 text-purple-300">AI 模型设置</h3>
          
          <div class="flex flex-col gap-2">
            <label class="text-gray-300">阿里云 千问大模型 API Key (可选)</label>
            <input 
              type="password" 
              v-model="apiKey" 
              @change="saveApiKey"
              placeholder="sk-xxxxxxxxxxxxxxxxxxxxxxxx"
              class="w-full bg-gray-900 border border-gray-600 rounded p-2 text-gray-100 focus:border-purple-500 outline-none"
            >
            <p class="text-xs text-gray-500">留空则使用模拟对话模式。API Key 将保存在本地浏览器中。</p>
          </div>
        </div>

        <!-- Audio -->
        <div class="space-y-4">
          <h3 class="text-xl border-b border-gray-700 pb-2">音频</h3>
          
          <div class="flex items-center justify-between">
            <label>主音量</label>
            <input type="range" min="0" max="1" step="0.1" v-model.number="settings.bgmVolume" @change="updateSettings" class="w-1/2">
          </div>
          
          <div class="flex items-center justify-between">
            <label>音效音量</label>
            <input type="range" min="0" max="1" step="0.1" v-model.number="settings.sfxVolume" @change="updateSettings" class="w-1/2">
          </div>
          
          <div class="flex items-center justify-between">
            <label>文字音效音量</label>
            <input type="range" min="0" max="1" step="0.1" v-model.number="settings.textVolume" @change="updateSettings" class="w-1/2">
          </div>
        </div>

        <!-- Display -->
        <div class="space-y-4">
          <h3 class="text-xl border-b border-gray-700 pb-2">显示</h3>
          
          <div class="flex items-center justify-between">
            <label>字体大小 ({{ settings.fontSize }}px)</label>
            <input type="range" min="14" max="24" step="1" v-model.number="settings.fontSize" @change="updateSettings" class="w-1/2">
          </div>
          
          <div class="flex items-center justify-between">
            <label>行高 ({{ settings.lineHeight }})</label>
            <input type="range" min="1.2" max="2.0" step="0.1" v-model.number="settings.lineHeight" @change="updateSettings" class="w-1/2">
          </div>
          
          <div class="flex items-center justify-between">
            <label>屏幕阅读器友好模式</label>
            <input type="checkbox" v-model="settings.isScreenReaderMode" @change="updateSettings" class="w-5 h-5">
          </div>
        </div>
      </div>
      
      <div class="mt-10 flex justify-end">
        <button @click="goBack" class="py-2 px-6 bg-gray-700 hover:bg-gray-600 rounded">返回</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSettingsStore } from '@/store/settingsStore'
import { audioManager } from '@/modules/AudioManager'

const router = useRouter()
const settings = useSettingsStore()

const apiKey = ref('')

onMounted(() => {
  apiKey.value = localStorage.getItem('damo_qwen_api_key') || ''
})

const saveApiKey = () => {
  localStorage.setItem('damo_qwen_api_key', apiKey.value)
}

const updateSettings = () => {
  audioManager.playSfx('click')
  settings.applyTheme()
  audioManager.updateVolumes()
}

const goBack = () => {
  audioManager.playSfx('click')
  router.back()
}
</script>
