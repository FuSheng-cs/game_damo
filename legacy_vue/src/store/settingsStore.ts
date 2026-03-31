import { defineStore } from 'pinia'

export interface SettingsState {
  fontSize: number; // in px
  lineHeight: number;
  themeColor: string;
  isDarkMode: boolean;
  bgmVolume: number; // 0 to 1
  sfxVolume: number; // 0 to 1
  textVolume: number; // 0 to 1
  isScreenReaderMode: boolean;
}

export const useSettingsStore = defineStore('settings', {
  state: (): SettingsState => ({
    fontSize: 16,
    lineHeight: 1.5,
    themeColor: '#a855f7',
    isDarkMode: true,
    bgmVolume: 0.5,
    sfxVolume: 0.8,
    textVolume: 0.5,
    isScreenReaderMode: false
  }),
  actions: {
    updateSettings(newSettings: Partial<SettingsState>) {
      Object.assign(this, newSettings)
      this.applyTheme()
    },
    applyTheme() {
      const root = document.documentElement
      root.style.setProperty('--font-size-base', `${this.fontSize}px`)
      root.style.setProperty('--line-height-base', `${this.lineHeight}`)
      root.style.setProperty('--theme-color', this.themeColor)
      
      if (this.isDarkMode) {
        root.classList.add('dark')
      } else {
        root.classList.remove('dark')
      }
    }
  }
})
