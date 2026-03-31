# 项目清理记录 (Cleanup Changelog)

在将游戏核心机制从“预设选项树（SceneRenderer）”全面升级为“AI原生对话驱动（LLMService）”，以及优化美术资源结构后，对项目进行了深度清理。

## 删除的文件与目录

### 1. 废弃的核心代码模块
*   `src/modules/SceneRenderer.ts`：原本用于处理固定节点、选项（Choice）和好感度分支的硬编码剧情引擎。在引入阿里云千问大模型后，剧情分支完全由 `LLMService.ts` 动态接管，此文件已无用。

### 2. 无用的初始模板与占位组件
*   `src/components/HelloWorld.vue`：Vite 初始模板自带的组件。
*   `src/assets/hero.png`
*   `src/assets/vue.svg`

### 3. 冗余的美术资源与废弃的特效方案
*   `public/assets/images/cigarette_full.png`
*   `public/assets/images/jump_1.png`
*   `public/assets/images/jump_2.png`
*   `public/assets/images/jump_3.png`
*   `public/assets/images/smoke_exhale.png`
*   `public/assets/images/smoke_inhale.png`
    *(上述资源为早期测试单帧动画时的残留，现已由新版的 CG 图和 CSS 特效替代。)*

### 4. 过期的设计文档
由于游戏故事线、机制和美术需求已全面重构并整合到了新的文档中，以下早期的草稿与设定文档已被删除以防混淆：
*   `docs/game_design.md`（旧版游戏设计，被 `storyline_and_lore.md` 替代）
*   `docs/technical_plan.md`（旧版技术方案，由于架构变更已不再准确）
*   `docs/resource_list.csv`（旧版资源清单，被 `art_assets_requirements.md` 替代）

## 保留与新增的核心结构
当前项目处于非常干净的状态，核心由以下部分组成：
*   **交互引擎**：`LLMService.ts` (负责与大模型通信和 Prompt 注入)
*   **状态管理**：`gameStore.ts` (管理对话历史、轮回次数、等待状态)
*   **存档系统**：`SaveSystem.ts` (防篡改的本地存储)
*   **成就系统**：`AchievementTracker.ts` (记录多结局的解锁状态)
*   **视图层**：`GameView.vue` (负责全屏 CG 渲染与毛玻璃对话框交互)
*   **最新设定文档**：位于 `docs/` 目录下的 `storyline_and_lore.md`、`art_assets_requirements.md` 和 `prompts_and_settings.md`。
