# DAMO - 零基础上线部署教程

> **适合人群**：没有服务器经验、第一次部署网站的同学  
> **目标**：把 DAMO 游戏部署到云服务器，让其他人通过域名或 IP 访问

---

## 目录

1. [准备工作](#1-准备工作)
2. [购买并配置服务器](#2-购买并配置服务器)
3. [连接服务器](#3-连接服务器)
4. [安装服务器环境](#4-安装服务器环境)
5. [部署 Go 后端](#5-部署-go-后端)
6. [构建并部署前端](#6-构建并部署前端)
7. [配置 Nginx](#7-配置-nginx)
8. [配置 HTTPS（推荐）](#8-配置-https推荐)
9. [设置开机自启](#9-设置开机自启)
10. [常见问题排查](#10-常见问题排查)

---

## 1. 准备工作

你需要：
- **一个 API Key**（任选一个）：
  - 阿里云千问：[dashscope.aliyuncs.com](https://dashscope.aliyuncs.com) → 控制台 → API Key 管理
  - 字节豆包：[ark.cn-beijing.volces.com](https://ark.cn-beijing.volces.com) → 火山方舟控制台
  - OpenAI：[platform.openai.com](https://platform.openai.com) → API Keys
- **本地安装**（你的电脑上）：
  - [Go 1.22+](https://go.dev/dl/)
  - [Node.js 20+](https://nodejs.org/)
  - [Git](https://git-scm.com/)

---

## 2. 购买并配置服务器

### 推荐配置
| 配置 | 说明 |
|------|------|
| 操作系统 | **Ubuntu 22.04 LTS**（强烈推荐，最稳定） |
| CPU | 1核 或 2核 |
| 内存 | 1GB 或 2GB（够用） |
| 带宽 | 1Mbps 以上 |
| 地域 | 国内（上海/北京/广州）访问快 |

### 购买步骤（以腾讯云为例）
1. 登录 [cloud.tencent.com](https://cloud.tencent.com)
2. 控制台 → **轻量应用服务器** → 新建
3. 选择 **Ubuntu 22.04** 镜像
4. 购买后在"安全组"中开放以下端口：
   - **22**（SSH 连接）
   - **80**（HTTP）
   - **443**（HTTPS，可选但推荐）
   - 不需要单独开放 8080（Nginx 代理）

> 阿里云步骤类似：产品 → 云服务器 ECS 或轻量应用服务器

---

## 3. 连接服务器

### Windows 用户（推荐用 PowerShell 或 MobaXterm）

```powershell
# 在 PowerShell 中运行（将 IP 替换为你的服务器公网 IP）
ssh root@你的服务器IP
```

第一次连接会提示"是否信任此主机"，输入 `yes` 回车。

> 💡 **简便方式**：下载 [MobaXterm](https://mobaxterm.mobatek.net/) 免费版，图形界面连接，自带文件传输。

---

## 4. 安装服务器环境

连接成功后，在服务器终端依次运行以下命令：

```bash
# 更新系统包列表
apt update && apt upgrade -y

# 安装 Nginx（网页服务器）
apt install -y nginx

# 安装 Git
apt install -y git

# 安装 Go 1.22
wget -q https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证 Go 安装成功
go version
# 应该输出：go version go1.22.4 linux/amd64

# 安装 Node.js 20（用于构建前端）
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt install -y nodejs

# 验证 Node.js 安装成功
node --version   # 应输出 v20.x.x
npm --version    # 应输出 10.x.x
```

---

## 5. 部署 Go 后端

```bash
# 创建项目目录
mkdir -p /opt/damo
cd /opt/damo

# 从你的 Git 仓库克隆代码（替换为你的仓库地址）
git clone https://github.com/你的用户名/game_damo.git .
# 如果没有 Git 仓库，也可以用 scp 上传文件（见文末说明）

# 进入后端目录
cd /opt/damo/backend

# 下载 Go 依赖
go mod tidy

# 编译后端（生成可执行文件）
go build -o damo-backend main.go

# 创建配置文件
cp .env.example .env
nano .env  # 编辑配置文件
```

在 `.env` 文件中填写（nano 编辑器：Ctrl+O 保存，Ctrl+X 退出）：

```bash
# 服务端 LLM 配置（可选，配置后玩家无需自己填 Key 即可游玩）
# LLM_PROVIDER=qwen
# LLM_API_KEY=sk-你的APIKey
# LLM_MODEL=qwen-plus

# 后端端口（保持默认即可）
PORT=8080
```

```bash
# 测试后端是否能正常启动
./damo-backend
# 应输出：[Server] DAMO Backend starting on :8080
# 按 Ctrl+C 停止（接下来用 systemd 管理）
```

---

## 6. 构建并部署前端

**在你的本地电脑上**运行：

```powershell
# 进入前端目录（你的本地项目）
cd d:\IT\game_damo\legacy_vue

# 安装依赖
npm install

# 创建生产环境配置（生产环境不需要设置 VITE_BACKEND_URL，Nginx 会处理）
# 新建文件 .env.production，内容留空或不创建都行
# 前端代码已经处理了空值情况

# 构建生产版本
npm run build
# 构建完成后，dist/ 目录就是要上传的静态文件
```

**将 `dist/` 目录上传到服务器**：

```powershell
# 方法1：使用 scp 上传（在本地 PowerShell 运行）
scp -r d:\IT\game_damo\legacy_vue\dist\* root@你的服务器IP:/var/www/damo/

# 方法2：在服务器上直接拉取（如果代码在 Git 仓库）
# 在服务器上运行：
# cd /opt/damo/legacy_vue && npm install && npm run build
# mkdir -p /var/www/damo && cp -r dist/* /var/www/damo/
```

如果用 MobaXterm，可以直接拖拽上传文件夹。

在服务器上确认文件：
```bash
# 在服务器上运行
mkdir -p /var/www/damo
ls /var/www/damo  # 应该能看到 index.html 和 assets 等文件
```

---

## 7. 配置 Nginx

```bash
# 创建 Nginx 配置文件
nano /etc/nginx/sites-available/damo
```

粘贴以下内容（将 `你的服务器IP或域名` 替换为实际值）：

```nginx
server {
    listen 80;
    server_name 你的服务器IP或域名;

    # 前端静态文件
    root /var/www/damo;
    index index.html;

    # 将 /api/* 请求代理到 Go 后端
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 60s;
    }

    # SPA 路由支持（刷新页面不 404）
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

```bash
# 启用配置
ln -s /etc/nginx/sites-available/damo /etc/nginx/sites-enabled/
# 删除默认配置（避免冲突）
rm -f /etc/nginx/sites-enabled/default
# 检查配置是否有语法错误
nginx -t
# 重启 Nginx
systemctl restart nginx
```

---

## 8. 配置 HTTPS（推荐）

有域名才能用 HTTPS。如果你只有 IP，可以跳过此步。

```bash
# 安装 Certbot（免费 SSL 证书工具）
apt install -y certbot python3-certbot-nginx

# 申请证书（替换为你的域名和邮箱）
certbot --nginx -d 你的域名.com -m 你的邮箱@example.com --agree-tos

# Certbot 会自动修改 Nginx 配置，无需手动操作
# 证书每 90 天自动续期
```

---

## 9. 设置开机自启

用 **systemd** 管理后端进程，确保服务器重启后自动运行：

```bash
# 创建 systemd 服务文件
nano /etc/systemd/system/damo-backend.service
```

粘贴以下内容：

```ini
[Unit]
Description=DAMO Game Backend
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/damo/backend
ExecStart=/opt/damo/backend/damo-backend
Restart=on-failure
RestartSec=5s
# 如果需要服务器端 API Key，在这里配置环境变量（可选）
# Environment="LLM_PROVIDER=qwen"
# Environment="LLM_API_KEY=sk-你的APIKey"

[Install]
WantedBy=multi-user.target
```

```bash
# 重新加载 systemd 配置
systemctl daemon-reload

# 启动后端服务
systemctl start damo-backend

# 设置开机自启
systemctl enable damo-backend

# 查看运行状态
systemctl status damo-backend
# 应显示 Active: active (running)
```

---

## ✅ 验证一切正常

```bash
# 1. 检查后端是否在运行
systemctl status damo-backend

# 2. 测试后端 API 是否正常响应
curl http://localhost:8080/api/health
# 应返回：{"status":"ok"}

# 3. 检查 Nginx 是否运行
systemctl status nginx

# 4. 测试 Nginx 代理是否正常
curl http://你的服务器IP/api/health
# 应返回：{"status":"ok"}
```

用浏览器访问 `http://你的服务器IP`，应该能看到游戏主页！

---

## 10. 常见问题排查

### 问题：网页打不开
```bash
# 检查 Nginx 日志
tail -f /var/log/nginx/error.log
# 检查端口是否监听
ss -tlnp | grep -E '80|8080'
# 检查安全组，确认已开放 80 端口
```

### 问题：AI 对话没有响应
```bash
# 检查后端日志
journalctl -u damo-backend -f
# 测试后端直接响应
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"user_message":"test","history":[],"rounds_left":10,"affection":0}'
```

### 问题：刷新页面 404
检查 Nginx 配置里是否有 `try_files $uri $uri/ /index.html;`

### 问题：更新代码后如何重新部署
```bash
# 在本地电脑重新构建前端（如果前端有改动）
npm run build
scp -r dist/* root@你的服务器IP:/var/www/damo/

# 在服务器上重新编译后端（如果后端有改动）
cd /opt/damo/backend
go build -o damo-backend main.go
systemctl restart damo-backend
```

### 没有 Git 仓库，如何上传代码
```powershell
# 在本地 PowerShell 中，用 scp 上传整个 backend 目录
scp -r d:\IT\game_damo\backend root@你的服务器IP:/opt/damo/
```

---

## 费用参考

| 项目 | 费用 |
|------|------|
| 腾讯云轻量服务器（2核2G） | 约 50-100元/月 |
| 阿里云 ECS（1核2G） | 约 40-80元/月 |
| SSL 证书（certbot） | **免费** |
| 域名（可选） | 约 50-100元/年 |
| AI API Key 费用 | 按调用量计费，千问 qwen-plus 约 0.004元/千token |

> 💡 每个玩家**使用自己的 API Key** 时，服务器运营成本接近零。
