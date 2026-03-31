const ACHIEVEMENT_KEY = 'damo_achievements'

export interface Achievement {
  id: string;
  name: string;
  description: string;
  unlockedAt?: number;
}

export const ACHIEVEMENTS: Achievement[] = [
  { id: 'end_death', name: '坠落', description: '未能触及她的内心，她选择了离开这个世界。' },
  { id: 'end_disappear', name: '消失', description: '平庸的交谈，她带着失望消失在夜色中。' },
  { id: 'end_acquaintance', name: '相识', description: '你真正触及了她的内心，拯救了她，并获得了她的联系方式。' },
  { id: 'first_try', name: '初次相遇', description: '第一次在天台遇见她。' }
]

export class AchievementTracker {
  static getUnlocked(): string[] {
    try {
      const data = localStorage.getItem(ACHIEVEMENT_KEY)
      return data ? JSON.parse(data) : []
    } catch {
      return []
    }
  }

  static unlock(id: string) {
    const unlocked = this.getUnlocked()
    if (!unlocked.includes(id)) {
      unlocked.push(id)
      localStorage.setItem(ACHIEVEMENT_KEY, JSON.stringify(unlocked))
      // You can also emit an event here to show a toast notification
      window.dispatchEvent(new CustomEvent('achievement-unlocked', { detail: id }))
    }
  }

  static getProgress(): { unlocked: number; total: number } {
    return {
      unlocked: this.getUnlocked().length,
      total: ACHIEVEMENTS.length
    }
  }
}
