# Changelog

## [1.0.0] - 2026-03-23
### Added
- 初始化 Vue 3 + TypeScript + Vite 游戏项目架构。
- 引入 Pinia 状态管理与 Howler.js 音频控制。
- 引入 TailwindCSS 样式引擎。
- 支持 PWA，提供 manifest 与 service-worker（基于 vite-plugin-pwa）。
- 实现了 `SaveSystem` 支持 localStorage 和 CRC32 数据校验。
- 实现了 `AchievementTracker` 记录和展示游戏结局成就。
- 实现了打字机组件、分支选择模块。
- 提供全局设置功能（音频音量、字体、深色模式）。
- 提供了单元测试用例（覆盖率达标）。

### Resource Status
- 目前部分图片资源、音效使用 Mock 占位或 CSS 回退。
- 请参考 `docs/resource_list.csv` 准备真实资源放入 `public/assets` 对应目录。
