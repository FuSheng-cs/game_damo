# 天台十句----你需要在十句话内拯救一个绝望的女孩

> “在夜色中记录人生百态，却唯独找不到自己存在的意义。”

[![Go](https://img.shields.io/badge/Backend-Go%201.22-00ADD8?logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203-4FC08D?logo=vue.js)](https://vuejs.org/)
---

## 📖 关于游戏

《天台十句》是一款由 AI 驱动的沉浸式心理叙事文字冒险游戏。

在这个霓虹闪烁的深夜，你来到了天台。坐在围栏边缘的是“艾”——一个紫色内染发、指间夹着忽明忽暗香烟的女孩。她是一个极度虚无的独立摄影师，带着一种颓废的破碎感，仿佛随时都会从这里滑落。

你第一眼就被她深深迷住。你想要了解她，拯救她。

**你的目标很简单，也很残酷：**
你初始只有 **10 句话** 的机会与她交谈。
你要在这个倒计时结束前，通过自由的文本对话，试着走进她的内心，解开她的心结。

### 🎭 游戏特色

- **AI 实时对话**：艾的每一句回复都由 AI 实时生成，没有固定的对话树
- **动态好感度系统**：说出真正触动她的话，获得额外轮次
- **多结局**：死亡、消失、相识——每种结局都是一个故事
- **故事延续**：达成"相识"结局后，可以继续与艾聊天,解锁"日后谈"
- **存档系统**：随时保存进度，继续未完的对话
- **成就图鉴**：收集所有结局，解锁隐藏内容

---

## 🛠️ 技术栈

```
frontend/   Vue 3 + Vite + TypeScript + Tailwind CSS v4
backend/    Go 1.22 + Gin + 标准库 HTTP 客户端
```

### AI 支持的模型服务商

| 服务商 | 推荐模型 | 获取 Key |
|--------|----------|----------|
| 阿里云千问 | `qwen-plus` | [dashscope.aliyuncs.com](https://dashscope.aliyuncs.com) |
| 字节豆包 | `doubao-pro-4k` | [ark.cn-beijing.volces.com](https://ark.cn-beijing.volces.com) |
| OpenAI | `gpt-4o-mini` | [platform.openai.com](https://platform.openai.com) |
| 自定义中转 | 任意 | 填入 Base URL 即可 |

---

## 🚀 本地运行

### 环境要求

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)

### 1. 启动 Go 后端

```bash
cd backend
cp .env.example .env
# 编辑 .env，可选配置服务器端 API Key
go mod tidy
go run main.go
# 后端运行在 http://localhost:8080
```

### 2. 启动前端开发服务器

```bash
cd legacy_vue
npm install
# .env.local 已预置 VITE_BACKEND_URL=http://localhost:8080
npm run dev
# 前端运行在 http://localhost:5173
```

### 3. 配置 API Key

打开浏览器访问 `http://localhost:5173`，进入**游戏设置**，选择 AI 服务商并填入你的 API Key。

---

## 📦 项目结构

```
game_damo/
├── backend/                    # Go 后端
│   ├── main.go                 # 入口文件（Gin 服务器）
│   ├── config/config.go        # 配置加载
│   ├── llm/service.go          # LLM 调用层（Prompt 在此保护）
│   ├── handlers/game.go        # HTTP 路由处理器
│   ├── go.mod
│   └── .env.example            # 环境变量模板
│
├── legacy_vue/                 # Vue 3 前端
│   ├── src/
│   │   ├── views/
│   │   │   ├── StartView.vue   # 游戏主界面
│   │   │   ├── GameView.vue    # 游戏主场景
│   │   │   ├── SettingsView.vue# 设置（AI Key 配置）
│   │   │   ├── ChatAfterStoryView.vue  # 结局后续聊天
│   │   │   └── AchievementsView.vue    # 成就图鉴
│   │   ├── modules/
│   │   │   ├── LLMService.ts   # 前端 LLM 代理（调用后端）
│   │   │   ├── AudioManager.ts
│   │   │   ├── SaveSystem.ts
│   │   │   └── AchievementTracker.ts
│   │   └── store/
│   │       ├── gameStore.ts    # 游戏状态管理（Pinia）
│   │       └── settingsStore.ts
│   └── .env.example
│
├── docs/                       # 设计文档
└── README.md
```

---

## 🔐 隐私与安全

- **API Key 安全**：玩家的 Key 仅在请求时传输，不做持久化存储（服务器不记录）
- **本地存储**：玩家 Key 和游戏存档仅保存在玩家自己的浏览器 localStorage

---


## 📝 开发日志

- `v1.0.0` — 初始版本，Vue3 前端 + Qwen API
- `v1.1.0` — 添加 Go 后端，支持多 AI 服务商（Qwen / 豆包 / OpenAI）



---

*"你只有 10 句话的时间。但也许，我们都需要更多的时间来原谅自己。"*
