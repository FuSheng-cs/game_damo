# DAMO — 你只有 10 句话

> 在霓虹灯闪烁的深夜，天台上坐着一个想要消失的女孩。  
> 你有 10 句话的时间，试着走进她的内心。

[![Go](https://img.shields.io/badge/Backend-Go%201.22-00ADD8?logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203-4FC08D?logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-purple)](LICENSE)

---

## 📖 关于游戏

DAMO 是一款 **AI 驱动的互动叙事游戏**。

你扮演一个深夜来到天台的陌生人，遇见了"艾"——一个带着颓废破碎感的独立摄影师，她正坐在围栏上，决意消失。

**你只有 10 句话。**

每一句话都会影响她的内心。普通的套话让她离你更远，真正触动她的话会让她停留更久，甚至改变结局。

### 🎭 游戏特色

- **AI 实时对话**：艾的每一句回复都由 AI 实时生成，没有固定的对话树
- **动态好感度系统**：说出真正触动她的话，获得额外轮次
- **多结局**：死亡、消失、相识——每种结局都是一个故事
- **故事延续**：达成"相识"结局后，可以继续与艾微信聊天
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
├── DEPLOYMENT_GUIDE.md         # 零基础部署教程
└── README.md
```

---

## 🔐 隐私与安全

- **Prompt 保护**：AI 角色设定、结局判定逻辑全部锁在 Go 后端，玩家无法通过 DevTools 看到
- **API Key 安全**：玩家的 Key 仅在请求时传输，不做持久化存储（服务器不记录）
- **本地存储**：玩家 Key 和游戏存档仅保存在玩家自己的浏览器 localStorage

---

## 🚢 部署

详见 [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) — 从零开始的完整部署教程，包含：
- 购买腾讯云/阿里云服务器
- 安装 Go、Nginx 环境
- 配置 systemd 自启动
- HTTPS 证书配置

---

## 📝 开发日志

- `v1.0.0` — 初始版本，Vue3 前端 + Qwen API
- `v1.1.0` — 添加 Go 后端，支持多 AI 服务商（Qwen / 豆包 / OpenAI）

---

## 📄 License

MIT License — 自由使用，但请保留原作者署名。

---

*"你只有 10 句话的时间。但也许，我们都需要更多的时间来原谅自己。"*
