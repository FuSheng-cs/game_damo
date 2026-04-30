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
    store.waitingText = 'waiting'
    store.resetGame()
    
    expect(store.roundCount).toBe(10)
    expect(store.messages.length).toBe(1)
    expect(store.waitingText).toBe('')
  })

  it('loadState restores only stable state', () => {
    const store = useGameStore()
    store.isWaiting = true
    store.waitingText = 'waiting'

    store.loadState({
      roundCount: 4,
      hintCount: 2,
      affection: 20,
      messages: [{ role: 'user', content: 'hello' }],
      isEnding: true,
      endingType: 'end_acquaintance'
    })

    expect(store.roundCount).toBe(4)
    expect(store.hintCount).toBe(2)
    expect(store.affection).toBe(20)
    expect(store.isWaiting).toBe(false)
    expect(store.waitingText).toBe('')
  })
})
