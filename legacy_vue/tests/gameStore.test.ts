import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGameStore } from '../src/store/gameStore'

describe('Game Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with correct default state', () => {
    const store = useGameStore()
    expect(store.roundCount).toBe(10)
    expect(store.messages.length).toBe(1)
    expect(store.isWaiting).toBe(false)
  })

  it('resetGame restores default state', () => {
    const store = useGameStore()
    store.roundCount = 5
    store.resetGame()
    
    expect(store.roundCount).toBe(10)
    expect(store.messages.length).toBe(1)
  })
})
