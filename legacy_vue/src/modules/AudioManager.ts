import { Howl, Howler } from 'howler'
import { useSettingsStore } from '@/store/settingsStore'

class AudioManager {
  private bgm: Howl | null = null;
  private sfxMap: Record<string, Howl> = {};
  
  constructor() {}

  init() {
    const settingsStore = useSettingsStore()
    Howler.volume(1.0) // Master volume
    
    // Preload some SFX
    this.sfxMap['click'] = new Howl({
      src: ['/assets/audio/ui_click.wav'],
      volume: settingsStore.sfxVolume
    })
    this.sfxMap['typewriter'] = new Howl({
      src: ['/assets/audio/ui_typewriter.wav'],
      volume: settingsStore.textVolume,
      loop: true
    })
  }

  playBgm(src: string) {
    if (this.bgm) {
      this.bgm.stop()
    }
    const settingsStore = useSettingsStore()
    this.bgm = new Howl({
      src: [src],
      loop: true,
      volume: settingsStore.bgmVolume
    })
    this.bgm.play()
  }

  stopBgm() {
    if (this.bgm) {
      this.bgm.stop()
    }
  }

  playSfx(name: string) {
    const settingsStore = useSettingsStore()
    if (this.sfxMap[name]) {
      this.sfxMap[name].volume(settingsStore.sfxVolume)
      this.sfxMap[name].play()
    }
  }

  startTypewriter() {
    const settingsStore = useSettingsStore()
    if (this.sfxMap['typewriter'] && !this.sfxMap['typewriter'].playing()) {
      this.sfxMap['typewriter'].volume(settingsStore.textVolume)
      this.sfxMap['typewriter'].play()
    }
  }

  stopTypewriter() {
    if (this.sfxMap['typewriter']) {
      this.sfxMap['typewriter'].stop()
    }
  }

  updateVolumes() {
    const settingsStore = useSettingsStore()
    if (this.bgm) {
      this.bgm.volume(settingsStore.bgmVolume)
    }
    Object.values(this.sfxMap).forEach(sfx => {
      sfx.volume(settingsStore.sfxVolume)
    })
    if (this.sfxMap['typewriter']) {
      this.sfxMap['typewriter'].volume(settingsStore.textVolume)
    }
  }
}

export const audioManager = new AudioManager()
