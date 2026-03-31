import { defineStore } from 'pinia'
import { LLMService } from '@/modules/LLMService'

export interface Message {
  role: 'user' | 'assistant';
  content: string;
}

export interface GameState {
  roundCount: number;
  hintCount: number;
  affection: number;
  messages: Message[];
  isWaiting: boolean;
  waitingText: string;
  isEnding: boolean;
  endingType: string | null;
}

const WAITING_TEXTS = [
  "她吐了一口烟圈……",
  "风太大了，她没听清……",
  "她在看着你的眼睛发呆……",
  "霓虹灯在她脸上闪烁……",
  "她轻轻弹了弹烟灰……"
];

export const useGameStore = defineStore('game', {
  state: (): GameState => ({
    roundCount: 10,
    hintCount: 3,
    affection: 0,
    messages: [
      { role: 'assistant', content: '夜晚的天台，微风吹过。霓虹灯的色彩在她的头发上跳跃。她坐在围栏上，指间的香烟忽明忽暗。你记得她，那个在夜色中游荡的摄影师。' }
    ],
    isWaiting: false,
    waitingText: '',
    isEnding: false,
    endingType: null
  }),
  actions: {

    async requestHint() {
      if (this.hintCount <= 0 || this.isWaiting || this.isEnding) return null;
      
      this.isWaiting = true;
      this.waitingText = "你在脑海中寻找线索……";
      this.hintCount -= 1;
      
      const hint = await LLMService.getHint(this.messages, {
        roundsLeft: this.roundCount,
        affection: this.affection
      });
      
      this.isWaiting = false;
      return hint;
    },
    async sendMessage(userText: string) {
      if (this.roundCount <= 0 || this.isEnding) return;
      
      this.messages.push({ role: 'user', content: userText });
      this.roundCount -= 1;
      
      this.isWaiting = true;
      this.waitingText = WAITING_TEXTS[Math.floor(Math.random() * WAITING_TEXTS.length)];
      
      const reply = await LLMService.chat(userText, this.messages.slice(0, -1), {
        roundsLeft: this.roundCount,
        affection: this.affection
      });
      
      this.isWaiting = false;
      
      let finalReply = reply;
      
      // Check for affection boost
      if (reply.includes('[好感度+5]')) {
        this.affection += 5;
        this.roundCount += 1;
        finalReply = finalReply.replace(/\[好感度\+5\]/g, '').trim();
      }

      // Match the exact ending tags defined in LLMService.ts
      const endingMatch = finalReply.match(/\[结局:(死亡|消失|相识)\]/);
      
      if (endingMatch) {
        this.isEnding = true;
        const typeStr = endingMatch[1];
        if (typeStr === '死亡') this.endingType = 'end_death';
        else if (typeStr === '消失') this.endingType = 'end_disappear';
        else if (typeStr === '相识') this.endingType = 'end_acquaintance';
        else this.endingType = 'end_disappear';
        
        finalReply = finalReply.replace(/\[结局:.*?\]/, '').trim();
      } else if (this.roundCount <= 0) {
        this.isEnding = true;
        this.endingType = 'end_disappear';
      }
      
      if (finalReply) {
        this.messages.push({ role: 'assistant', content: finalReply });
      }
    },
    resetGame() {
      this.roundCount = 10;
      this.hintCount = 3;
      this.affection = 0;
      this.messages = [
        { role: 'assistant', content: '夜晚的天台，微风吹过。霓虹灯的色彩在她的头发上跳跃。她坐在围栏上，指间的香烟忽明忽暗。你记得她，那个在夜色中游荡的摄影师。' }
      ];
      this.isWaiting = false;
      this.isEnding = false;
      this.endingType = null;
    },
    loadState(state: any) {
      this.roundCount = state.roundCount;
      this.hintCount = state.hintCount ?? 3;
      this.affection = state.affection ?? 0;
      this.messages = state.messages || [];
      this.isEnding = state.isEnding || false;
      this.endingType = state.endingType || null;
    }
  }
})
