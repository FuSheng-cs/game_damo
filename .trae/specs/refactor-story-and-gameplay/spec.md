# 重构故事与游戏机制 Spec

## Why
用户希望重构当前的游戏故事线和底层设定。摒弃之前“潜意识模拟器”、“男主车祸迟到”的设定，转而塑造一个极具破碎感、神秘感、带有藤本树《再见绘梨》与《炎拳》精神内核的独立女性摄影师角色。同时引入基于好感度的动态回合奖励机制，并在达成好结局后新增一个类似微信聊天的后日谈页面。

## What Changes
- **重构角色与背景设定**：
  - 剔除之前所有的“轮回、潜意识模拟器、妹妹车祸、男主迟到”设定。
  - **新设定**：女孩是个无父无母的孤儿，独立摄影师，在夜色中记录人生百态。她目前极度虚无（类似《炎拳》的精神内核），坐在天台抽烟想要轻生。玩家第一眼被她迷住，想要了解并拯救她。她极具魅力、破碎感、神秘感和一丝奇幻色彩（类似《再见绘梨》）。
- **修改游戏机制**：
  - 基础仍是 10 句话的机会。
  - **动态回合奖励**：当玩家的话触动到她内心或非常有意思时，由 AI 判定好感度增加 5 点，并在前端页面显示。每次好感度增加，额外获得 1 次对话机会。
- **重构结局**：
  - 改回三个基础结局：**死亡**、**消失**、**相识**。
  - 【相识】极难达到，必须真正触及她的内心。
- **新增后日谈页面（相识之后）**：
  - 在达成【相识】结局后，不再是之前那个男主自白的“阅读真相”页面。
  - 新增一个类似微信聊天界面的视图（`ChatAfterStoryView.vue`），用于达成好结局后与女孩的日常聊天。
- **BREAKING**: 原有的 `damo_meta_data` 轮回机制逻辑将被废弃或大幅简化。真结局页面 `TrueEndingView.vue` 将被替换。

## Impact
- Affected specs: `docs/storyline_and_lore.md`, `docs/prompts_and_settings.md`, `README.md`
- Affected code: `src/modules/LLMService.ts`, `src/store/gameStore.ts`, `src/views/GameView.vue`, `src/router/index.ts`
- Added code: `src/views/ChatAfterStoryView.vue`

## ADDED Requirements
### Requirement: 动态回合奖励与好感度前端展示
The system SHALL provide a mechanism where the AI determines if a user's message is touching or interesting. If yes, it outputs a specific tag (e.g., `[好感度+5]`), which the game store parses to increase affection by 5 and add 1 round to the remaining count. The affection bar must be visible on the frontend.

#### Scenario: Success case
- **WHEN** user inputs a deeply empathetic message
- **THEN** AI responds with `[好感度+5]`, remaining rounds increase by 1, affection bar increases, and AI replies favorably.

### Requirement: 微信聊天后日谈页面
The system SHALL provide a new route `/chat-after` that simulates a mobile chat interface for casual conversation with the girl after achieving the '相识' ending.

## MODIFIED Requirements
### Requirement: 结局判定体系
**Reason**: User requested to revert to 3 endings and remove the simulation lore.
**Migration**: Update the `LLMService.ts` prompt to only output `[结局:死亡]`, `[结局:消失]`, or `[结局:相识]`. Update `gameStore.ts` to parse only these three and route appropriately.

## REMOVED Requirements
### Requirement: 潜意识模拟器与轮回记忆
**Reason**: Conflicts with the new independent character lore requested by the user.
**Migration**: Remove the Déjà vu prompt logic, the previous summary injection, and the `TrueEndingView.vue` narrative.