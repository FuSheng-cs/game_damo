// LLMService - 前端 LLM 服务代理层
// 此文件只负责将请求转发给 Go 后端，不再直接调用任何 LLM 服务商。
// 系统 Prompt 已移至后端保护，前端不可见。

// LocalStorage 中存储玩家 LLM 配置的键名
const LS_KEYS = {
  PROVIDER: 'damo_llm_provider',
  API_KEY: 'damo_llm_api_key',
  MODEL: 'damo_llm_model',
  BASE_URL: 'damo_llm_base_url'
} as const

// 后端地址：生产环境留空（Nginx 同源代理），开发环境用 .env.local 注入
const BACKEND_URL = (import.meta.env.VITE_BACKEND_URL || '').replace(/\/$/, '')

// LLM 配置结构
export interface LLMConfig {
  provider: string   // openai | qwen | doubao
  apiKey: string
  model: string      // 可选，留空使用 Provider 默认模型
  baseUrl: string    // 可选，自定义 Base URL
}

// 读取玩家保存在本地的 LLM 配置
export function loadLLMConfig(): LLMConfig {
  return {
    provider: localStorage.getItem(LS_KEYS.PROVIDER) || 'qwen',
    apiKey: localStorage.getItem(LS_KEYS.API_KEY) || '',
    model: localStorage.getItem(LS_KEYS.MODEL) || '',
    baseUrl: localStorage.getItem(LS_KEYS.BASE_URL) || ''
  }
}

// 保存玩家的 LLM 配置到本地
export function saveLLMConfig(cfg: LLMConfig) {
  localStorage.setItem(LS_KEYS.PROVIDER, cfg.provider)
  localStorage.setItem(LS_KEYS.API_KEY, cfg.apiKey)
  localStorage.setItem(LS_KEYS.MODEL, cfg.model)
  localStorage.setItem(LS_KEYS.BASE_URL, cfg.baseUrl)
}

// 向后端发送请求的通用函数
async function callBackend(endpoint: string, body: object): Promise<string> {
  const url = `${BACKEND_URL}/api/${endpoint}`

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })

    const data = await response.json()

    if (data.error) {
      console.error(`[LLMService] Backend error on /${endpoint}:`, data.error)
      return data.error
    }

    return data.reply ?? ''
  } catch (e) {
    console.error(`[LLMService] Network error on /${endpoint}:`, e)
    return '网络连接失败，请检查后端服务是否启动。'
  }
}

export class LLMService {
  // 主游戏对话
  static async chat(userMessage: string, history: any[], loopContext: any): Promise<string> {
    const cfg = loadLLMConfig()

    return callBackend('chat', {
      history: history,
      user_message: userMessage,
      rounds_left: loopContext.roundsLeft,
      affection: loopContext.affection,
      provider: cfg.provider,
      api_key: cfg.apiKey,
      model: cfg.model,
      base_url: cfg.baseUrl
    })
  }

  // 故事结束后的续聊
  static async chatAfterStory(userMessage: string, history: any[]): Promise<string> {
    const cfg = loadLLMConfig()

    return callBackend('chat-after', {
      history: history,
      user_message: userMessage,
      provider: cfg.provider,
      api_key: cfg.apiKey,
      model: cfg.model,
      base_url: cfg.baseUrl
    })
  }

  // 获取游戏提示
  static async getHint(history: any[], _loopContext: any): Promise<string> {
    const cfg = loadLLMConfig()

    return callBackend('hint', {
      history: history,
      provider: cfg.provider,
      api_key: cfg.apiKey,
      model: cfg.model,
      base_url: cfg.baseUrl
    })
  }
}
