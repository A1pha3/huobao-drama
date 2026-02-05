# 安装指南

> **学习目标**：完成本指南后，你将能够独立完成火宝短剧平台的完整安装部署，包括开发环境和生产环境的配置。
> 
> **前置知识**：无需任何前置知识，本指南为零基础设计。
> 
> **难度等级**：⭐（入门教程）
> 
> **预计学习时间**：30-60分钟

---

## 概述

### 关于火宝短剧

火宝短剧是一款基于人工智能的短剧生成平台，能够自动化完成从剧本创作、角色设计、分镜制作到视频生成的完整工作流程。平台利用大语言模型理解剧本内容，提取角色形象和场景描述，再通过扩散模型生成高质量的视觉素材，最终使用视频合成技术将静态图片转化为动态视频。

本指南将帮助你完成火宝短剧平台的安装部署。我们提供两种部署方式：**Docker容器化部署**（推荐）和**本地直接部署**。无论你选择哪种方式，只要按照本指南的步骤操作，都能够顺利完成安装。

### 系统要求

在开始安装之前，请确保你的系统满足以下要求。这些要求同时适用于**开发环境**和**生产环境**，不过开发环境的配置可以适当放宽。

| 软件组件 | 最低版本 | 推荐版本 | 说明 |
|----------|----------|----------|------|
| 操作系统 | Windows 10/macOS 10.15/Ubuntu 20.04 | Windows 11/macOS 14/Ubuntu 22.04 | 建议使用64位操作系统 |
| 处理器 | x86-64或ARM64 | 多核处理器 | 多核可提升编译和运行性能 |
| 内存 | 4GB | 8GB或更高 | 视频生成需要较大内存 |
| 存储空间 | 10GB可用空间 | 50GB可用空间 | 需要存放项目文件和依赖 |
| 网络 | 宽带互联网 | 高速互联网 | 用于下载依赖和访问AI服务 |

---

## 第一部分：Docker容器化部署（推荐）

### 1.1 为什么选择Docker

Docker容器化技术是当前软件部署的主流方案，具有以下优势：**环境一致性**——开发、测试、生产环境使用相同的容器镜像，彻底消除「在我机器上能运行」的问题；**隔离性**——容器与应用主机环境隔离，不会互相影响；**可移植性**——容器可以在任何支持Docker的系统上运行；**快速启动**——容器可以在几秒钟内启动和停止。

对于大多数用户，我们推荐使用Docker进行部署。即使你之前没有使用过Docker的经验，本指南也将带你一步步完成安装过程。

### 1.2 安装Docker

**Windows系统安装**：

1. 访问Docker官方网站下载Docker Desktop：https://www.docker.com/products/docker-desktop/
2. 下载适用于Windows的Docker Desktop安装包
3. 运行安装程序，按照向导完成安装
4. 安装完成后，Docker Desktop会自动启动
5. 在系统托盘中看到Docker图标表示安装成功

安装过程中，如果系统提示需要启用WSL 2（Windows Subsystem for Linux 2），请按照提示完成WSL 2的安装和配置。WSL 2是Docker在Windows上运行的必要组件。

**macOS系统安装**：

1. 访问Docker官方网站下载Docker Desktop：https://www.docker.com/products/docker-desktop/
2. 下载适用于Apple Silicon（M系列芯片）或Intel芯片的Docker Desktop安装包
3. 将Docker.dmg拖拽到应用程序文件夹完成安装
4. 从应用程序文件夹启动Docker Desktop
5. 在菜单栏中看到Docker图标表示安装成功

**Linux系统安装**（以Ubuntu为例）：

```bash
# 更新软件包列表
sudo apt update

# 安装Docker依赖包
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# 添加Docker官方GPG密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 添加Docker软件仓库
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 将当前用户添加到docker组（免sudo运行Docker）
sudo usermod -aG docker $USER

# 重新登录使组权限生效
newgrp docker
```

### 1.3 验证Docker安装

无论使用哪种操作系统，安装完成后都应该验证Docker是否正常工作。打开终端（Windows用户使用PowerShell或命令提示符），执行以下命令：

```bash
# 查看Docker版本
docker --version

# 预期输出：
# Docker version 24.0.7 或更高版本

# 运行Docker hello-world测试
docker run --rm hello-world

# 预期输出：
# Hello from Docker!
# This message shows that your installation appears to be working correctly.
```

如果命令执行后显示了版本信息和欢迎信息，说明Docker安装成功。如果出现错误信息，请根据错误提示进行排查。

### 1.4 获取项目代码

项目代码托管在GitHub上，你可以使用Git克隆或直接下载压缩包两种方式获取代码。

**方式一：使用Git克隆**（需要先安装Git）：

```bash
# 克隆项目代码
git clone https://github.com/chatfire-AI/huobao-drama.git

# 进入项目目录
cd huobao-drama

# 切换到稳定版本（生产环境建议使用tag）
git checkout v1.0.4
```

**方式二：直接下载**：

1. 访问项目 releases 页面：https://github.com/chatfire-AI/huobao-drama/releases
2. 下载最新版本的源代码压缩包
3. 解压到本地目录
4. 进入解压后的目录

### 1.5 配置环境变量

项目使用`.env`文件管理环境变量，这种方式使得配置与代码分离，便于管理不同环境的配置。

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置文件
vim .env
```

配置文件中的关键配置项说明如下：

```bash
# ===========================================
# 应用配置
# ===========================================
APP_NAME=Huobao Drama
APP_VERSION=v1.0.4
APP_DEBUG=false

# ===========================================
# 服务端口配置
# ===========================================
SERVER_PORT=5678
SERVER_HOST=0.0.0.0
```

对于国内用户，建议配置Docker镜像加速以提高镜像拉取速度：

```bash
# 创建Docker配置目录
sudo mkdir -p /etc/docker

# 创建daemon配置文件
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://dockerproxy.com",
    "https://docker.mirrors.ustc.edu.cn"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  }
}
EOF

# 重启Docker服务
sudo systemctl restart docker
```

### 1.6 启动服务

完成所有配置后，可以启动服务了。首次启动会自动构建Docker镜像，可能需要较长时间（取决于网络速度）。

```bash
# 构建并启动所有服务（后台运行）
docker compose up -d --build

# 查看服务状态
docker compose ps

# 预期输出：
# NAME                 STATUS
# huobao-drama-app-1   Up (healthy)

# 查看应用日志
docker compose logs -f app
```

服务启动后，可以通过浏览器访问`http://localhost:5678`进入平台。如果一切正常，你应该能看到火宝短剧的用户界面。

### 1.7 验证安装成功

打开浏览器，依次访问以下地址确认各组件正常工作：

| 地址 | 预期结果 |
|------|----------|
| http://localhost:5678 | 显示火宝短剧用户界面 |
| http://localhost:5678/health | 返回`{"status":"ok"}` |
| http://localhost:5678/api/v1/projects | 返回授权错误（表示API正常） |

如果所有地址都能正常响应，说明安装成功！

---

## 第二部分：本地直接部署

### 2.1 环境准备

如果无法使用Docker，或者希望进行本地开发调试，可以选择直接部署。这种方式需要手动安装所有依赖。

**Go语言安装**：

```bash
# 下载Go安装包（Linux/macOS）
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz

# 解压到/usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

**Node.js安装**（仅用于前端开发）：

```bash
# 使用nvm安装Node.js（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install --lts
nvm use --lts

# 验证安装
node --version
npm --version
```

**FFmpeg安装**（必需组件）：

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install -y ffmpeg

# 验证安装
ffmpeg -version
```

### 2.2 项目配置

```bash
# 进入项目目录
cd huobao-drama

# 安装Go依赖
go mod download

# 构建后端
go build -o huobao-drama .
```

### 2.3 前端开发环境

```bash
# 进入前端目录
cd web

# 安装依赖
npm install --registry=https://registry.npmmirror.com/

# 启动开发服务器
npm run dev
```

### 2.4 启动服务

**开发模式**（前后端分离）：

```bash
# 终端1：启动后端
cd huobao-drama
go run main.go

# 终端2：启动前端
cd huobao-drama/web
npm run dev
```

**生产模式**（后端同时提供前端服务）：

```bash
# 构建前端
cd web
npm run build
cd ..

# 构建后端
go build -o huobao-drama .

# 启动服务
./huobao-drama
```

---

## 第三部分：配置AI服务

### 3.1 AI服务概述

火宝短剧的AI功能需要连接外部AI服务才能使用。平台支持多种AI服务提供商，你可以根据需求选择合适的服务。

| 服务提供商 | 主要功能 | 获取地址 |
|------------|----------|----------|
| OpenAI | 文本生成、图像生成 | https://platform.openai.com |
| DeepSeek | 文本生成 | https://platform.deepseek.com |
| 豆包 | 文本生成、图像生成 | https://www.volcengine.com |

### 3.2 获取API密钥

1. 访问所选AI服务提供商的官网，注册账号并完成实名认证
2. 进入API管理页面，创建新的API密钥
3. 复制并妥善保管API密钥（API密钥具有与密码同等的重要性）

**安全提示**：API密钥是访问AI服务的凭证，请勿将API密钥泄露给他人或上传到公开的代码仓库。建议定期更换API密钥。

### 3.3 配置AI服务

编辑项目根目录下的`.env`文件：

```bash
# OpenAI配置
OPENAI_API_KEY=sk-your-api-key-here
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4

# DeepSeek配置
DEEPSEEK_API_KEY=sk-your-api-key-here
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-chat

# 豆包配置
DOUBAO_API_KEY=your-api-key-here
DOUBAO_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
DOUBAO_MODEL=doubao-pro-32k

# 默认AI服务提供商
AI_DEFAULT_PROVIDER=openai
```

配置完成后，重启服务使配置生效。

---

## 第四部分：常见问题

### Q1：Docker启动失败怎么办？

首先检查端口是否被占用：

```bash
# 检查5678端口占用情况
netstat -tulpn | grep 5678
# 或使用
ss -tulpn | grep 5678
```

如果端口被占用，可以修改`.env`文件中的`SERVER_PORT`配置项，使用其他端口。

### Q2：国内访问AI服务不稳定怎么办？

建议使用国内AI服务提供商（如DeepSeek或豆包），或者配置代理服务器。在`.env`文件中配置代理：

```bash
HTTP_PROXY=http://proxy.example.com:7890
HTTPS_PROXY=http://proxy.example.com:7890
```

### Q3：Docker镜像下载缓慢怎么办？

配置Docker镜像加速器，参考1.5节的配置步骤。配置完成后需要重启Docker服务。

### Q4：启动后无法访问Web界面？

1. 检查防火墙是否允许对应端口
2. 检查服务是否正常启动：`docker compose ps`
3. 查看应用日志排查问题：`docker compose logs -f app`

---

## 练习任务

### 练习1：Docker环境搭建（⭐）

**任务目标**：完成Docker环境搭建和火宝短剧服务启动

**任务要求**：

1. 安装Docker Desktop/Docker Engine
2. 验证Docker安装成功
3. 获取火宝短剧项目代码
4. 启动服务并验证可以访问Web界面

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| `docker --version`返回正确版本 | ☐ |
| `docker run hello-world`执行成功 | ☐ |
| `docker compose ps`显示服务状态为Up | ☐ |
| 浏览器可访问http://localhost:5678 | ☐ |

---

## 自检清单

### 环境检查

| 检查项 | 完成情况 |
|--------|----------|
| 操作系统满足要求 | ☐ |
| 内存≥4GB | ☐ |
| 存储空间≥10GB | ☐ |
| 网络连接正常 | ☐ |

### Docker检查

| 检查项 | 完成情况 |
|--------|----------|
| Docker安装成功 | ☐ |
| Docker版本≥20.10 | ☐ |
| Docker Compose安装成功 | ☐ |
| Docker镜像加速配置完成（国内用户） | ☐ |

### 服务检查

| 检查项 | 完成情况 |
|--------|----------|
| 服务启动成功 | ☐ |
| 健康检查通过 | ☐ |
| Web界面可访问 | ☐ |
| AI服务配置完成 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [快速入门](quick-start.md) | 15分钟创建第一个项目 | ⭐ |
| [用户手册](user-manual.md) | 完整功能操作指南 | ⭐⭐ |
| [高级用法](advanced_usage.md) | 复杂场景处理 | ⭐⭐⭐ |
| [部署运维指南](../deployment.md) | 生产环境部署详细指南 | ⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本，包含完整安装指南 |
