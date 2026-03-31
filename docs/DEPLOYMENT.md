# 部署与性能说明 (Deployment & Performance)

## 部署说明
本项目基于 Vite 构建，构建后的产物为纯静态文件（HTML/CSS/JS/Assets），可部署在任何静态文件服务器或 CDN 上（如 Vercel, Netlify, GitHub Pages, Nginx 等）。

### 构建步骤
1. 安装依赖: `npm install`
2. 执行构建: `npm run build`
3. 产物路径: `/dist` 目录下的所有文件。
4. 部署: 将 `/dist` 目录上传至服务器的 Web 根目录即可。

由于配置了 PWA（`vite-plugin-pwa`），服务器需要配置 HTTPS 以确保 Service Worker 正常注册和离线可用。

## 性能预算与指标
- **首屏渲染时间**：目标 <= 1.5s
- **Lighthouse 性能分**：目标 >= 90
- **首次加载体积**：核心 JS/CSS 和首屏 UI 必须 <= 300KB（Gzip 后）。由于所有非首屏背景图和音频均按需或延迟加载，首屏开销已被最小化。
- **图片策略**：推荐全部转为 WebP 格式，保留 PNG 作为 Fallback。所有背景切片大小按移动端和PC端提供两套，通过 `<picture>` 或 CSS 媒体查询响应式加载。

## 无障碍声明 (Accessibility - WCAG 2.2 AA)
本项目在开发中致力于达到 WCAG 2.2 AA 级别：
1. **对比度**：确保文本与背景颜色的对比度达到 4.5:1 以上（主色调为深色背景+高亮文字）。
2. **键盘导航**：游戏内所有可交互元素（选项、按钮、设置滑块）均可通过 `Tab` 键聚焦，通过 `Enter`/`Space` 触发。
3. **屏幕阅读器**：在设置中提供“屏幕阅读器友好模式”。打字机动画文本将添加 `aria-live="polite"` 标签，关键视觉信息使用 `aria-label` 进行描述。
4. **焦点顺序**：合理的 DOM 结构保证了焦点从左到右、从上到下的自然移动。
