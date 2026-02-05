# 部署运维指南

> **学习目标**：完成本指南学习后，你将能够独立完成火宝短剧平台的生产环境部署，掌握Docker和传统部署两种方式，理解配置管理、监控告警、日志收集等运维最佳实践，具备日常运维和故障处理能力。
> 
> **前置知识**：建议具备Linux系统操作基础，熟悉Docker基本概念和命令，了解Nginx反向代理配置。建议先完成《快速入门指南》的学习，理解平台的整体架构。
> 
> **难度等级**：⭐⭐（进阶）
> 
> **预计学习时间**：2-4小时

---

## 第一章：部署策略概述

### 1.1 部署方式对比

火宝短剧平台支持多种部署方式，包括**Docker容器化部署**和**传统二进制部署**。选择合适的部署方式需要综合考虑团队技术栈、运维能力、扩展需求和成本预算等因素。本节详细分析各种部署方式的特点，帮助你做出正确的选择。

**Docker容器化部署**是推荐的部署方式，特别适合以下场景：团队熟悉容器化技术、需要快速弹性扩展、希望标准化部署流程、使用云原生基础设施。Docker部署的核心优势在于**环境一致性**——开发、测试、生产环境使用相同的容器镜像，彻底消除「在我机器上能运行」的问题。此外，Docker的隔离性使得多个服务可以互不干扰地运行在同一台服务器上，资源利用率更高。

**传统二进制部署**适合以下场景：服务器资源有限（无法运行Docker）、团队不熟悉容器技术、需要最大化性能、简单的单服务部署。传统部署的优势在于资源占用少、运维简单，但环境配置的一致性难以保证，大规模部署时配置管理复杂。

| 对比维度 | Docker部署 | 传统部署 |
|----------|-----------|----------|
| 环境一致性 | ✅ 完全一致 | ⚠️ 依赖配置管理 |
| 资源占用 | 中等（约100-200MB） | 低（仅应用本身） |
| 扩展能力 | ✅ 原生支持 | 需额外工具 |
| 学习成本 | 需要Docker知识 | 较低 |
| 运维复杂度 | 中等（容器编排） | 低（直接运行） |
| 健康检查 | ✅ 原生支持 | 需额外实现 |
| 版本回滚 | ✅ 快速回滚 | 需手动处理 |

**选择建议**：对于大多数生产环境，推荐使用Docker部署。只有在服务器资源极其有限或团队完全没有容器经验的情况下，才考虑传统部署。无论选择哪种方式，都建议先在测试环境充分验证后再部署到生产环境。

### 1.2 生产环境要求

生产环境的服务器配置直接影响平台的性能和稳定性。以下是生产环境部署的硬件和软件要求，按照小型、中型、大型三种规模分别说明。

**硬件配置要求**：

| 规模 | CPU | 内存 | 存储 | 带宽 | 适用场景 |
|------|-----|------|------|------|----------|
| 小型 | 2核 | 4GB | 50GB SSD | 5Mbps | 个人用户、小团队 |
| 中型 | 4核 | 8GB | 100GB SSD | 10Mbps | 中等规模团队 |
| 大型 | 8核+ | 16GB+ | 500GB SSD+ | 20Mbps+ | 企业级部署 |

**存储说明**：平台生成的内容（图片、视频）会占用大量存储空间。建议预留足够的存储容量，并实施定期清理策略。视频文件是存储增长的主要来源，1080p视频每分钟约占用100-200MB空间。

**软件环境要求**：

| 组件 | 最低版本 | 推荐版本 | 说明 |
|------|----------|----------|------|
| Docker | 20.10 | 24.0+ | 容器运行时 |
| Docker Compose | 2.0 | 2.21+ | 容器编排工具 |
| FFmpeg | 4.0 | 6.0+ | 视频处理引擎 |
| 系统 | Ubuntu 20.04 | Ubuntu 22.04 LTS | 操作系统 |

**网络要求**：平台需要访问外部AI服务API，确保服务器能够访问互联网。如果部署环境存在网络隔离，需要配置相应的代理或白名单。

---

## 第二章：Docker容器化部署

### 2.1 环境准备

在开始部署之前，需要完成服务器环境的准备工作。以下步骤适用于Ubuntu 22.04 LTS系统，其他Linux发行版可参考对应命令。

**步骤一：更新系统软件包**：

```bash
# 更新软件包列表
sudo apt update

# 升级已安装的软件包
sudo apt upgrade -y

# 安装必要的基础工具
sudo apt install -y curl wget git vim htop net-tools
```

**步骤二：安装Docker**：

Docker官方提供了自动安装脚本，可以快速完成Docker的安装。这种方式适用于初次接触Docker的用户，安装过程会自动检测系统环境并安装所需的依赖。

```bash
# 使用Docker官方安装脚本
curl -fsSL https://get.docker.com | sh -s -- --version 24.0

# 或者使用国内镜像安装
curl -fsSL https://get.docker.com | sh -s -- --mirror Aliyun

# 将当前用户添加到docker组（免sudo执行docker命令）
sudo usermod -aG docker $USER

# 重新登录以使组权限生效
newgrp docker
```

安装完成后，验证Docker是否正确安装：

```bash
# 查看Docker版本
docker --version
# 预期输出：Docker version 24.0.7 或更高版本

# 运行Docker hello-world测试
docker run --rm hello-world
# 预期输出：显示Hello from Docker!欢迎信息
```

**步骤三：安装Docker Compose**：

Docker Compose用于定义和管理多容器应用，是Docker部署的核心工具。

```bash
# 安装Docker Compose（官方方式）
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
# 预期输出：Docker Compose version v2.21.0 或更高版本
```

对于国内用户，建议配置Docker镜像加速器以提高镜像拉取速度：

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

# 验证配置生效
docker info | grep -A 5 "Registry Mirrors"
```

**步骤四：安装FFmpeg**：

FFmpeg是视频处理的核心组件，平台依赖它进行视频的编解码、格式转换等操作。正确安装FFmpeg是视频功能正常工作的前提。

```bash
# 安装FFmpeg
sudo apt install -y ffmpeg

# 验证安装
ffmpeg -version
# 预期输出：ffmpeg version 6.x 或更高版本

# 检查支持的编码器（确认h264编码器可用）
ffmpeg -encoders | grep h264
# 预期输出应包含h264_videotoolbox（macOS）或libx264（Linux）
```

### 2.2 项目配置

完成环境准备后，需要配置火宝短剧项目。项目提供了预设的配置文件模板，可以快速完成配置。

**步骤一：获取项目代码**：

```bash
# 克隆项目代码
git clone https://github.com/chatfire-AI/huobao-drama.git

# 进入项目目录
cd huobao-drama

# 切换到稳定版本分支（生产环境建议使用tag）
git checkout v1.0.4
```

**步骤二：配置环境变量**：

项目使用`.env`文件管理环境变量，这种方式使得配置与代码分离，便于管理不同环境的配置。

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置文件
vim .env
```

**配置文件详解**：

```bash
# ===========================================
# Docker镜像加速配置（国内用户推荐启用）
# ===========================================
DOCKER_REGISTRY=docker.1ms.run/

# ===========================================
# npm镜像加速（加速前端依赖下载）
# ===========================================
NPM_REGISTRY=https://registry.npmmirror.com/

# ===========================================
# Go模块代理（加速Go依赖下载）
# ===========================================
GO_PROXY=https://goproxy.cn,direct

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

# ===========================================
# 数据库配置（SQLite）
# ===========================================
DATABASE_PATH=./data/drama_generator.db

# ===========================================
# 存储配置
# ===========================================
STORAGE_TYPE=local
STORAGE_PATH=./data/storage
STORAGE_BASE_URL=http://localhost:5678/static

# ===========================================
# AI服务配置（根据实际需要填写）
# ===========================================
# OpenAI配置
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4

# DeepSeek配置
DEEPSEEK_API_KEY=your_deepseek_api_key_here
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-chat

# 豆包配置
DOUBAO_API_KEY=your_doubao_api_key_here
DOUBAO_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
DOUBAO_MODEL=doubao-pro-32k

# 默认AI服务提供商
AI_DEFAULT_PROVIDER=openai

# ===========================================
# 日志配置
# ===========================================
LOG_LEVEL=info
LOG_FORMAT=json
```

**配置说明**：配置文件中的AI服务API Key是敏感信息，务必妥善保管，不要将包含真实API Key的配置文件提交到代码仓库。建议在生产环境中使用密钥管理服务来存储这些敏感信息。

**步骤三：配置目录权限**：

Docker容器运行时需要对某些目录有读写权限。确保数据目录正确创建并设置权限：

```bash
# 创建数据目录
mkdir -p ./data/storage/{images,videos,cache,exports}

# 设置目录权限
chmod -R 755 ./data

# 确保目录所有者正确
ls -la ./data/
# 预期输出：显示data目录及其子目录的详细信息
```

### 2.3 Docker Compose配置

项目提供了完整的`docker-compose.yml`配置文件，定义了应用运行所需的所有服务。

**配置文件分析**：

```yaml
# docker-compose.yml 完整配置

services:
  # 主应用服务
  app:
    image: ${DOCKER_REGISTRY}huobao-drama:${APP_VERSION:-latest}
    build:
      context: .
      dockerfile: Dockerfile
    container_name: huobao-drama-app
    restart: unless-stopped
    ports:
      - "${SERVER_PORT}:5678"
    volumes:
      # 数据卷挂载
      - ./data:/app/data
      # 配置文件挂载
      - ./configs:/app/configs:ro
    environment:
      # 应用环境变量
      - APP_ENV=production
      - APP_DEBUG=${APP_DEBUG:-false}
      - DATABASE_PATH=/app/data/drama_generator.db
      - STORAGE_PATH=/app/data/storage
      - STORAGE_BASE_URL=http://localhost:${SERVER_PORT:-5678}/static
    networks:
      - huobao-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:5678/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "3"

networks:
  huobao-network:
    driver: bridge

volumes:
  app_data:
    driver: local
```

**关键配置说明**：

| 配置项 | 说明 | 推荐值 |
|--------|------|--------|
| restart | 容器重启策略 | unless-stopped（异常退出自动重启） |
| ports | 端口映射 | "${SERVER_PORT}:5678" |
| volumes | 数据卷挂载 | 分离数据和配置 |
| healthcheck | 健康检查 | 定期检测服务状态 |
| logging | 日志配置 | 限制日志文件大小和数量 |

### 2.4 启动与验证

完成所有配置后，可以启动服务。首次启动会自动构建Docker镜像，可能需要较长时间。

**步骤一：构建并启动服务**：

```bash
# 构建并启动所有服务（后台运行）
docker compose up -d --build

# 查看构建和启动过程（可选，实时查看日志）
docker compose logs -f

# 仅查看应用日志
docker compose logs -f app
```

**启动过程中的常见问题**：

如果在构建过程中遇到网络超时问题，可以检查Docker镜像加速配置是否生效：

```bash
# 检查Docker配置
docker info | grep "Registry Mirrors"

# 如果配置未生效，尝试重新加载
sudo systemctl daemon-reload
sudo systemctl restart docker
```

**步骤二：验证服务状态**：

服务启动后，需要验证各组件是否正常工作。

```bash
# 查看所有容器状态
docker compose ps

# 预期输出
# NAME                 STATUS
# huobao-drama-app-1   Up (healthy)

# 健康检查
curl http://localhost:5678/health

# 预期输出（JSON格式）
# {"status":"ok","timestamp":"2026-02-05T12:00:00Z"}

# 检查API可用性
curl http://localhost:5678/api/v1/projects

# 预期输出（未授权但表示服务正常）
# {"code":401,"message":"unauthorized"}
```

**步骤五：访问Web界面**：

打开浏览器访问`http://服务器IP:5678`，应该能看到火宝短剧的用户界面。如果页面加载缓慢或出现错误，请检查服务器防火墙是否允许5678端口的访问。

```bash
# 检查防火墙状态（Ubuntu系统）
sudo ufw status

# 如果防火墙已启用，允许5678端口
sudo ufw allow 5678
sudo ufw reload
```

---

## 第三章：传统二进制部署

### 3.1 编译构建

对于无法使用Docker的环境，可以直接编译Go应用程序进行部署。这种方式资源占用更少，但环境配置的一致性需要自行保证。

**步骤一：安装Go运行时**：

```bash
# 下载Go安装包
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

**步骤二：构建前端**：

```bash
# 进入前端目录
cd web

# 安装依赖
npm install --registry=https://registry.npmmirror.com/

# 构建生产版本
npm run build

# 验证构建结果
ls -la dist/
# 应该看到index.html和static目录
```

**步骤三：构建后端**：

```bash
# 返回项目根目录
cd ..

# 构建Go应用
go build -o huobao-drama .

# 验证构建结果
ls -la huobao-drama
# 应该看到可执行文件
```

### 3.2 系统服务配置

将应用配置为系统服务，可以实现开机自启动和进程管理。

**步骤一：创建服务用户**：

出于安全考虑，不建议使用root用户运行应用。

```bash
# 创建专门的应用用户
sudo useradd -r -s /sbin/nologin huobao

# 创建应用目录
sudo mkdir -p /opt/huobao-drama

# 复制应用文件
sudo cp huobao-drama /opt/huobao-drama/
sudo cp -r configs /opt/huobao-drama/
sudo cp -r web/dist /opt/huobao-drama/web/

# 创建数据目录
sudo mkdir -p /opt/huobao-drama/data/{storage,cache,exports}
sudo mkdir -p /opt/huobao-drama/logs

# 设置目录权限
sudo chown -R huobao:huobao /opt/huobao-drama
sudo chmod -R 755 /opt/huobao-drama
```

**步骤二：创建systemd服务**：

```bash
# 创建服务配置文件
sudo tee /etc/systemd/system/huobao-drama.service <<-'EOF'
[Unit]
Description=Huobao Drama AI Short Video Platform
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=huobao
Group=huobao
WorkingDirectory=/opt/huobao-drama
ExecStart=/opt/huobao-drama/huobao-drama
Restart=on-failure
RestartSec=10
Environment=APP_ENV=production
EnvironmentFile=/opt/huobao-drama/.env

# 日志配置
StandardOutput=journal
StandardError=journal
SyslogIdentifier=huobao-drama

# 资源限制
LimitNOFILE=65536
LimitNPROC=4096

# 安全加固
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/huobao-drama/data

[Install]
WantedBy=multi-user.target
EOF

# 重新加载systemd
sudo systemctl daemon-reload

# 启用服务（开机自启动）
sudo systemctl enable huobao-drama

# 启动服务
sudo systemctl start huobao-drama

# 查看服务状态
sudo systemctl status huobao-drama
```

**步骤三：验证服务运行**：

```bash
# 查看服务状态
sudo systemctl status huobao-drama

# 查看实时日志
sudo journalctl -u huobao-drama -f

# 检查进程是否运行
ps aux | grep huobao-drama
```

---

## 第四章：反向代理配置

### 4.1 Nginx反向代理

生产环境中，通常使用Nginx作为反向代理，将外部请求转发到应用服务。Nginx可以提供SSL终止、负载均衡、静态文件服务等功能。

**Nginx配置文件示例**：

```nginx
# /etc/nginx/sites-available/huobao-drama

server {
    listen 80;
    server_name your-domain.com;
    
    # 重定向到HTTPS（生产环境建议启用）
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    # SSL证书配置
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    
    # API请求代理
    location /api/ {
        proxy_pass http://127.0.0.1:5678;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 600s;
        proxy_read_timeout 600s;
        
        # WebSocket支持（如需要）
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    
    # 静态文件服务（如果前端已构建）
    location / {
        root /opt/huobao-drama/web/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    
    # 静态资源访问
    location /static/ {
        alias /opt/huobao-drama/data/storage/;
        expires 7d;
        add_header Cache-Control "public, immutable";
    }
    
    # 健康检查
    location /health {
        proxy_pass http://127.0.0.1:5678/health;
    }
}
```

**启用配置并测试**：

```bash
# 创建软链接启用配置
sudo ln -s /etc/nginx/sites-available/huobao-drama /etc/nginx/sites-enabled/

# 测试配置语法
sudo nginx -t

# 重载Nginx
sudo systemctl reload nginx

# 检查Nginx状态
sudo systemctl status nginx
```

### 4.2 SSL证书配置

生产环境必须启用HTTPS来保护数据传输安全。推荐使用Let's Encrypt提供的免费SSL证书。

```bash
# 安装Certbot
sudo apt install -y certbot python3-certbot-nginx

# 获取SSL证书
sudo certbot --nginx -d your-domain.com

# 自动续期测试
sudo certbot renew --dry-run

# 设置自动续期
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

---

## 第五章：监控与日志

### 5.1 日志配置

完善的日志记录是问题排查和系统审计的基础。火宝短剧平台使用Zap日志库，提供高性能的结构化日志输出。

**日志配置示例**（configs/config.yaml）：

```yaml
log:
  level: info           # 日志级别：debug, info, warn, error
  format: json         # 输出格式：json, console
  output: both         # 输出位置：file, stdout, both
  file:
    path: ./logs/app.log
    max_size: 100      # 单个日志文件大小（MB）
    max_backups: 10    # 保留的日志文件数量
    max_age: 30        # 日志文件保留天数（天）
```

**日志级别说明**：

| 级别 | 使用场景 |
|------|----------|
| debug | 开发调试，记录详细执行过程 |
| info | 正常运行，记录关键操作事件 |
| warn | 警告信息，不影响功能但需要注意 |
| error | 错误信息，需要关注和处理 |

**日志查看命令**（Docker环境）：

```bash
# 实时查看应用日志
docker compose logs -f app

# 查看最近100行日志
docker compose logs --tail 100 app

# 按时间过滤日志
docker compose logs --since "2026-02-05" app

# 查看错误日志
docker compose logs app 2>&1 | grep -i error
```

### 5.2 监控告警

建议配置基本的监控告警，以便及时发现和处理系统异常。

**系统级监控指标**：

| 指标 | 阈值建议 | 告警方式 |
|------|----------|----------|
| CPU使用率 | > 80% 持续5分钟 | 邮件/短信 |
| 内存使用率 | > 85% | 邮件/短信 |
| 磁盘使用率 | > 90% | 邮件/短信 |
| API响应时间 | P99 > 3秒 | 邮件 |
| 错误率 | > 1% | 邮件/短信 |

**使用Prometheus和Grafana监控**（可选）：

对于需要更精细监控的生产环境，可以集成Prometheus监控系统：

```yaml
# docker-compose.monitoring.yml

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: huobao-prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.enable-lifecycle'
    networks:
      - huobao-network

  grafana:
    image: grafana/grafana:latest
    container_name: huobao-grafana
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    networks:
      - huobao-network

volumes:
  prometheus_data:
  grafana_data:
```

### 5.3 性能基准

以下是生产环境的性能基准参考，用于评估系统负载能力和规划资源。

**单节点性能指标**：

| 指标 | 数值 | 测试条件 |
|------|------|----------|
| API响应时间（P50） | < 100ms | 单用户，无AI调用 |
| API响应时间（P99） | < 500ms | 单用户，无AI调用 |
| 并发请求上限 | ~500 req/s | CPU 4核 / 内存 8GB |
| 视频生成速度 | ~30秒/分钟 | 1080p分辨率 |

---

## 第六章：备份与恢复

### 6.1 备份策略

生产环境的数据是核心资产，必须制定并执行可靠的备份策略。建议采用**3-2-1备份原则**：保留3份数据副本，存储在2种不同的介质上，其中1份存放在异地。

**备份内容清单**：

| 备份项 | 备份频率 | 保留时间 | 存储位置 |
|--------|----------|----------|----------|
| 数据库 | 每日+变更时 | 30天 | 本地+云存储 |
| 上传文件 | 每日增量 | 7天 | 本地+云存储 |
| 配置文件 | 每次变更 | 永久 | Git仓库 |
| 日志文件 | 每日归档 | 90天 | 本地 |

**自动备份脚本**：

```bash
#!/bin/bash
# /opt/huobao-drama/scripts/backup.sh

# 配置变量
BACKUP_DIR=/opt/huobao-drama/backups
DATE=$(date +%Y%m%d_%H%M%S)
DB_FILE=/opt/huobao-drama/data/drama_generator.db
STORAGE_DIR=/opt/huobao-drama/data/storage

# 创建备份目录
mkdir -p ${BACKUP_DIR}

# 备份数据库
echo "备份数据库..."
cp ${DB_FILE} ${BACKUP_DIR}/db_${DATE}.sqlite3
gzip ${BACKUP_DIR}/db_${DATE}.sqlite3

# 备份配置文件
echo "备份配置文件..."
tar -czf ${BACKUP_DIR}/config_${DATE}.tar.gz /opt/huobao-drama/configs/

# 清理旧备份（保留最近7天）
echo "清理旧备份..."
find ${BACKUP_DIR} -name "db_*.sqlite3.gz" -mtime +7 -delete
find ${BACKUP_DIR} -name "config_*.tar.gz" -mtime +7 -delete

# 计算备份校验和
echo "计算校验和..."
md5sum ${BACKUP_DIR}/db_${DATE}.sqlite3.gz >> ${BACKUP_DIR}/checksums.txt

echo "备份完成：${BACKUP_DIR}/db_${DATE}.sqlite3.gz"
```

配置定时任务自动执行备份：

```bash
# 编辑crontab
crontab -e

# 添加每日凌晨2点执行备份
0 2 * * * /opt/huobao-drama/scripts/backup.sh >> /var/log/huobao-backup.log 2>&1
```

### 6.2 恢复流程

当发生数据丢失或系统故障时，需要按照以下步骤恢复数据。

**数据库恢复**：

```bash
# 停止应用服务
sudo systemctl stop huobao-drama

# 备份当前数据库（以防恢复失败）
sudo cp /opt/huobao-drama/data/drama_generator.db /opt/huobao-drama/data/drama_generator.db.backup

# 从备份恢复
gunzip /opt/huobao-drama/backups/db_20260205_020000.sqlite3.gz
cp /opt/huobao-drama/backups/db_20260205_020000.sqlite3 /opt/huobao-drama/data/drama_generator.db

# 修复数据库完整性
sqlite3 /opt/huobao-drama/data/drama_generator.db "PRAGMA integrity_check;"

# 重启应用服务
sudo systemctl start huobao-drama
```

**完整恢复检查清单**：

| 步骤 | 操作 | 验证方法 |
|------|------|----------|
| 1 | 停止服务 | systemctl status确认已停止 |
| 2 | 备份当前数据 | 确认备份文件存在 |
| 3 | 恢复数据库 | PRAGMA integrity_check通过 |
| 4 | 恢复文件 | 检查storage目录完整性 |
| 5 | 重启服务 | curl健康检查接口 |
| 6 | 验证功能 | 测试核心功能 |

---

## 第七章：安全加固

### 7.1 服务器安全

服务器安全是系统安全的第一道防线。以下是必要的安全加固措施。

**用户与权限安全**：

```bash
# 创建应用专用用户
sudo useradd -r -s /sbin/nologin huobao

# 禁用root远程登录
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# 配置防火墙
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

**定期安全更新**：

```bash
# 配置自动安全更新
sudo apt install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

### 7.2 应用安全

**API安全建议**：

1. **启用HTTPS**：所有生产环境必须启用HTTPS
2. **速率限制**：配置API请求频率限制，防止暴力破解
3. **输入验证**：严格验证所有用户输入
4. **敏感信息保护**：API Key等敏感信息使用环境变量或密钥管理服务
5. **CORS配置**：合理配置跨域资源共享策略

**Nginx速率限制配置**：

```nginx
http {
    # 限制请求频率
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    
    server {
        location /api/ {
            limit_req zone=api_limit burst=20;
            
            # 其他代理配置...
        }
    }
}
```

---

## 第八章：运维检查清单

### 8.1 日常检查项

每日运维检查是确保系统稳定运行的重要措施。以下是建议的检查项目和频率。

**每日检查（建议执行时间：上午9点前）**：

| 检查项 | 操作命令 | 正常标准 |
|--------|----------|----------|
| 服务状态 | systemctl status huobao-drama | active (running) |
| 磁盘空间 | df -h /opt | < 80% |
| 内存使用 | free -h | < 80% |
| 错误日志 | journalctl -u huobao-drama --since yesterday | 无ERROR级别日志 |
| 健康检查 | curl http://localhost:5678/health | {"status":"ok"} |

**每周检查**：

| 检查项 | 操作命令 | 正常标准 |
|--------|----------|----------|
| 备份完整性 | ls -la /opt/huobao-drama/backups/ | 备份文件存在 |
| SSL证书有效期 | certbot certificates | > 30天 |
| 慢查询日志 | 检查数据库慢查询 | < 10条/天 |
| 依赖安全漏洞 | npm audit / go list -m | 无高危漏洞 |

### 8.2 故障恢复流程

当系统发生故障时，按照以下流程进行恢复。

**故障响应流程**：

```
故障发现
    │
    ▼
初步评估（影响范围判断）
    │
    ▼
├─→ 服务可恢复？
│   │
│   ├─→ 是：执行恢复操作 → 验证恢复 → 记录故障
│   │
│   └─→ 否：回滚到上一版本 → 验证回滚 → 紧急修复
│
▼
故障恢复后分析
    │
    ▼
制定预防措施
```

---

## 练习任务

### 练习1：Docker部署实践（⭐⭐）

**任务目标**：完成生产环境的Docker部署

**任务要求**：

1. 在Linux服务器上安装Docker和Docker Compose
2. 配置项目环境变量
3. 使用Docker Compose启动服务
4. 验证所有功能正常工作

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| Docker安装验证通过 | ☐ |
| Docker Compose安装验证通过 | ☐ |
| 服务启动成功 | ☐ |
| 健康检查通过 | ☐ |
| Web界面可访问 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [快速入门指南](quick-start.md) | 15分钟快速上手 | ⭐ |
| [用户手册](user-manual.md) | 平台功能的完整操作指南 | ⭐⭐ |
| [API接口文档](api-reference.md) | 开发者API参考手册 | ⭐⭐⭐ |
| [架构设计文档](architecture.md) | 系统架构深度解析 | ⭐⭐⭐⭐ |
| [故障排查文档](troubleshooting.md) | 常见问题和解决方案 | ⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-05 | 初始版本，包含完整部署运维指南 |
