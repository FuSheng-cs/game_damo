import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGameStore } from '../src/store/gameStore'
import { SaveSystem } from '../src/modules/SaveSystem'

describe('Save System', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('saves and loads state successfully', () => {
    const store = useGameStore()
    store.roundCount = 5
    store.affection = 15
    
    const success = SaveSystem.save(1)
    expect(success).toBe(true)
    
    // Reset store
    store.resetGame()
    expect(store.roundCount).toBe(10)
    
    // Load
    const loadSuccess = SaveSystem.load(1)
    expect(loadSuccess).toBe(true)
    expect(store.roundCount).toBe(5)
    expect(store.affection).toBe(15)
  })

  it('does not save while waiting for an async response', () => {
    const store = useGameStore()
    store.roundCount = 5
    store.isWaiting = true

    const success = SaveSystem.save(1)

    expect(success).toBe(false)
    expect(localStorage.getItem('damo_save_1')).toBeNull()
  })

  it('supports three independent save slots', () => {
    const store = useGameStore()
    store.roundCount = 9
    expect(SaveSystem.save(1)).toBe(true)

    store.roundCount = 6
    store.affection = 5
    expect(SaveSystem.save(2)).toBe(true)

    store.roundCount = 3
    store.affection = 10
    expect(SaveSystem.save(3)).toBe(true)

    const slots = SaveSystem.getSlots()
    expect(slots.map((slot) => slot.id)).toEqual([1, 2, 3])

    store.resetGame()
    expect(SaveSystem.load(2)).toBe(true)
    expect(store.roundCount).toBe(6)
    expect(store.affection).toBe(5)
  })

  it('fails to load if data is tampered', () => {
    const store = useGameStore()
    SaveSystem.save(1)
    
    // Tamper data
    const savedStr = localStorage.getItem('damo_save_1')
    if (savedStr) {
      const saved = JSON.parse(savedStr)
      // Base64 decode, modify, encode
      const jsonStr = decodeURIComponent(atob(saved.data))
      const payload = JSON.parse(jsonStr)
      payload.state.roundCount = 999 // Cheat
      saved.data = btoa(encodeURIComponent(JSON.stringify(payload)))
      localStorage.setItem('damo_save_1', JSON.stringify(saved))
    }
    
    // Try to load
    const loadSuccess = SaveSystem.load(1)
    expect(loadSuccess).toBe(false)
  })
})
