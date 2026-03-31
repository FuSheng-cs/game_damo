# 优化女主设定并增加真结局故事揭示 (Optimize Persona & True Ending Reveal)

## Summary

本次修改旨在提升女主“艾”的真实感与魅力，使其更像一个正常的、有血有肉的女孩，并调整真结局的流程：在达成拯救（释怀）结局后，增加一个独立的回忆录页面，向玩家完整揭示背后的故事。

## Current State Analysis

1. **女主设定**：当前 System Prompt (`src/modules/LLMService.ts`) 中，女主的设定过于侧重“厌世”、“颓废”和“防御机制”，导致其表现有时显得过于冰冷或像是在故意作对，缺乏正常女孩的柔软、生活化气息和互动主动性。
2. **结局揭示**：目前达成真结局（`end_true_release`）后，仅在 `GameView.vue` 中显示最终对话，随后玩家只能点击“返回标题”。故事的完整真相（雨夜、便利店、迟到、妹妹车祸、男主构建模拟器）并没有清晰地传达给玩家，玩家可能通关了仍一头雾水。

## Proposed Changes

### 1. 优化女主性格设定 (Prompt 修改)

* **文件**：`src/modules/LLMService.ts`

* **What/Why**：修改 System Prompt，增加“柔软与日常感”、“更直白生活化的语言”以及“互动主动性”。让她不再是一个纯粹的防卫机器，而是一个会偶尔展现脆弱、会用日常语言交流、甚至会对玩家产生好奇的普通女孩。这能大大提升角色的魅力和玩家的共情。

* **How**：

  * 在【表层人格】中加入：偶尔会流露出柔软、犹豫，使用更直白、生活化的词汇，减少谜语人。

  * 在【交流策略】中加入：增加互动主动性，当玩家表现友善时，偶尔会主动反问玩家的问题，展现出对外界的一丝好奇。

### 2. 增加独立的回忆录页面 (True Ending Reveal)

* **文件**：

  * 新增 `src/views/TrueEndingView.vue`

  * 修改 `src/router/index.ts`

  * 修改 `src/views/GameView.vue`

* **What/Why**：在玩家达成真结局并阅读完最后一句对话后，点击“继续”或“返回”按钮时，不直接返回首页，而是路由跳转到一个全新的、类似日记或独白的纯文本页面，向玩家完整讲述整个故事的来龙去脉。

* **How**：

  * 创建 `TrueEndingView.vue`，设计一个安静、极简的 UI（纯黑背景，白色打字机文字或渐显文字），以男主的口吻讲述雨夜的真相和构建这个模拟器的初衷。

  * 在 `router/index.ts` 中注册 `/true-ending` 路由。

  * 在 `GameView.vue` 中，当 `isEnding` 为 true 且 `endingType` 为 `end_true_release` 时，将原本的“返回标题”按钮修改为“阅读真相”或“醒来”，点击后跳转到 `/true-ending`。

  * `TrueEndingView.vue` 播放完毕后，提供返回首页的按钮。

## Assumptions & Decisions

* **故事文案**：回忆录的文案将基于 `docs/storyline_and_lore.md` 中的设定撰写，采用第一人称（男主视角），讲述他迟到的雨夜、妹妹的车祸、艾的坠楼，以及他为了原谅自己而创造了这个潜意识模拟器的真相。

* **路由逻辑**：真结局页面是一个独立的视图，与游戏主界面分离，以营造一种“大梦初醒”的氛围。

## Verification steps

1. 检查 `LLMService.ts` 的 Prompt 是否已正确更新。
2. 测试路由配置，确保能访问 `/true-ending` 页面。
3. 在 `GameStore` 中手动模拟触发真结局（修改 state），验证 `GameView` 按钮是否正确变为跳转真结局页面的按钮。
4. 进入真结局页面，检查排版、文字动画和返回首页功能是否正常。

