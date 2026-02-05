# 开发环境搭建

> **学习目标**：完成本指南后，你将能够搭建完整的火宝短剧开发环境，具备本地开发和调试的能力。
> 
> **前置知识**：具备基本的编程经验和命令行使用能力。
> 
> **难度等级**：⭐（入门教程）
> 
> **预计学习时间**：1-2小时

---

## 第一章：环境准备

### 1.1 系统要求

在开始搭建开发环境之前，请确保你的系统满足以下要求。这些要求针对开发环境进行了优化，与生产环境有所不同。

| 组件 | 最低要求 | 推荐配置 | 说明 |
|------|----------|----------|------|
| 操作系统 | Windows 10/macOS 10.15/Ubuntu 20.04 | Windows 11/macOS 14/Ubuntu 22.04 | 建议使用64位操作系统 |
| 处理器 | x86-64或ARM64 | 多核处理器 | 多核可提升编译速度 |
| 内存 | 8GB | 16GB或更高 | 开发环境需要更多内存 |
| 存储空间 | 20GB可用空间 | 50GB可用空间 | 需要存放代码、依赖和项目 |
| 网络 | 宽带互联网 | 高速互联网 | 用于下载依赖和访问AI服务 |

### 1.2 开发工具链

开发火宝短剧需要安装以下基础工具。

**版本控制（必需）**：Git是版本控制工具，用于管理代码变更和协作开发。访问https://git-scm.com/downloads下载并安装。安装完成后，在终端执行`git --version`验证安装。

**代码编辑器（必需）**：推荐使用Visual Studio Code（免费）或GoLand（商业，需要许可证）。Visual Studio Code推荐安装以下扩展：Go扩展（代码提示和格式化）、ESLint（JavaScript代码检查）、Prettier（代码格式化）、Docker扩展（容器管理）。

**容器工具（推荐）**：Docker用于本地运行服务，简化开发和测试流程。请参考安装指南章节完成Docker安装。

---

## 第二章：Go开发环境

### 2.1 Go安装

火宝短剧后端使用Go语言开发，需要安装Go运行时和SDK。

**Windows和macOS安装**：

1. 访问Go官网下载页面：https://go.dev/dl/
2. 下载对应操作系统的安装包（Windows用`.msi`，macOS用`.pkg`）
3. 运行安装包，按照向导完成安装
4. 打开新的终端窗口，执行`go version`验证安装

**Linux安装**（以Ubuntu为例）：

```bash
# 下载Go安装包
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz

# 解压到/usr/local（需要sudo权限）
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

**国内镜像配置**：国内网络访问Go模块可能较慢，建议配置国内镜像：

```bash
# 配置GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct

# 查看当前配置
go env
```

### 2.2 Go模块管理

火宝短剧使用Go Modules管理依赖。

**Go.mod文件**：`go.mod`文件定义了项目的模块路径和依赖。首次克隆项目后，需要下载所有依赖：

```bash
# 下载所有依赖
go mod download

# 验证依赖下载
go list -m all | head -20

# 清理未使用的依赖
go mod tidy
```

**依赖版本管理**：当添加新依赖时，使用`go get`命令：

```bash
# 添加特定版本的依赖
go get github.com/gin-gonic/gin@v1.9.1

# 添加最新的依赖
go get github.com/some/package@latest

# 升级依赖
go get -u github.com/some/package
```

### 2.3 编译运行

**编译项目**：

```bash
# 开发模式编译
go build -o huobao-drama-dev .

# 生产模式编译
go build -o huobao-drama -ldflags="-s -w" .

# 只检查编译错误
go build -n
```

**运行项目**：

```bash
# 直接运行（开发模式）
go run main.go

# 使用编译后的二进制运行
./huobao-drama-dev
```

**热重载**：开发过程中，每次修改代码后手动停止并重启服务比较繁琐。使用Air工具实现热重载：

```bash
# 安装Air
go install github.com/air-verse/air@latest

# 创建Air配置文件
cat > .air.toml << 'EOF'
[root]
target = "main.go"

[build]
cmd = "go build -o ./tmp/main ."
bin = "./tmp/main"
include_ext = ["go", "tpl", "tmpl", "html"]
exclude_dir = ["tmp", "vendor", "web/node_modules"]
delay = 1000

[log]
time_format = "2006/01/02 - 15:04:05"
main_only = true

[screen]
clear_on_rebuild = true
EOF

# 启动开发服务器（自动热重载）
air
```

---

## 第三章：前端开发环境

### 3.1 Node.js安装

火宝短剧前端使用Vue 3框架，需要安装Node.js运行时。

**使用nvm安装（推荐）**：nvm（Node Version Manager）允许在同一系统上安装多个Node.js版本。

```bash
# macOS/Linux安装nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 重新加载shell配置
source ~/.bashrc

# 安装LTS版本的Node.js
nvm install --lts

# 验证安装
node --version
npm --version
```

**国内镜像配置**：

```bash
# 配置npm镜像
npm config set registry https://registry.npmmirror.com/

# 配置yarn镜像（如果使用yarn）
yarn config set registry https://registry.npmmirror.com/
```

### 3.2 前端依赖安装

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 验证安装
npm list --depth=0
```

**依赖说明**：前端项目的主要依赖包括Vue 3（前端框架）、Vite（构建工具）、Element Plus（UI组件库）、TailwindCSS（样式框架）、Pinia（状态管理）、Vue Router（路由管理）。

### 3.3 前端开发服务器

**开发模式运行**：

```bash
cd web
npm run dev
```

开发服务器运行后，会自动打开浏览器窗口访问开发地址。通常是http://localhost:5173或http://localhost:3012。

**构建生产版本**：

```bash
cd web
npm run build

# 查看构建产物
ls -la dist/
```

**预览生产构建**：

```bash
npm run preview
```

---

## 第四章：Docker开发环境

### 4.1 Docker Compose配置

对于不熟悉本地环境配置的开发者，可以使用Docker搭建开发环境。

**Docker Compose开发配置**：

创建`docker-compose.dev.yml`文件：

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: huobao-drama-dev
    ports:
      - "5678:5678"
      - "3012:3012"
    volumes:
      - .:/app
      - app_modules:/app/web/node_modules
    environment:
      - GO_ENV=development
      - AI_DEFAULT_PROVIDER=openai
    command: >
      bash -c "air -c .air.toml"
    networks:
      - huobao-network

volumes:
  app_modules:

networks:
  huobao-network:
    driver: bridge
```

**Dockerfile.dev**：

```dockerfile
# 开发环境Dockerfile
FROM golang:1.23-alpine

# 安装基础工具
RUN apk add --no-cache git curl unzip nodejs npm

# 安装Air实现热重载
RUN go install github.com/air-verse/air@latest

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 5678 3012

# 启动命令（在docker-compose中指定）
CMD ["tail", "-f", "/dev/null"]
```

**启动开发环境**：

```bash
# 构建并启动
docker compose -f docker-compose.dev.yml up -d --build

# 查看日志
docker compose -f docker-compose.dev.yml logs -f

# 进入容器
docker exec -it huobao-drama-dev sh
```

### 4.2 开发工具容器

可以创建专门的开发工具容器，避免在主机安装过多工具。

```yaml
services:
  dev-tools:
    image: golang:1.23
    container_name: huobao-dev-tools
    volumes:
      - .:/app
    working_dir: /app
    command: sleep infinity
    networks:
      - huobao-network

networks:
  huobao-network:
    driver: bridge
```

使用开发工具容器：

```bash
# 执行Go命令
docker exec huobao-dev-tools go build -o bin/app .

# 执行npm命令
docker exec huobao-dev-tools npm run build --prefix web
```

---

## 第五章：数据库环境

### 5.1 SQLite配置

火宝短剧使用SQLite作为默认数据库。SQLite是嵌入式数据库，不需要单独的服务器进程。

**数据库文件位置**：默认数据库文件位于`./data/drama_generator.db`。首次运行应用时，会自动创建数据库文件并执行迁移。

**数据库工具**：可以使用以下工具查看和编辑SQLite数据库：

| 工具 | 平台 | 说明 |
|------|------|------|
| SQLite Browser | 全平台 | 免费的可视化工具 |
| DBeaver | 全平台 | 通用数据库管理工具 |
| TablePlus | macOS/Windows | 现代化数据库工具 |

**命令行访问**：

```bash
# 打开数据库
sqlite3 data/drama_generator.db

# SQLite命令
.tables          # 查看所有表
.schema table_name  # 查看表结构
.exit            # 退出
```

### 5.2 数据库迁移

首次运行或代码更新后，需要执行数据库迁移：

```bash
# 手动执行迁移
go run main.go migrate

# 或者使用CLI工具
./huobao-drama migrate
```

**迁移文件位置**：迁移文件位于`migrations/`目录，每个迁移文件包含`up`和`down`两个迁移脚本。

---

## 第六章：外部服务配置

### 6.1 AI服务配置

开发环境需要配置AI服务才能使用AI功能。

**配置方法**：在项目根目录创建`.env.development`文件：

```bash
# AI服务配置
AI_DEFAULT_PROVIDER=openai

# OpenAI配置
OPENAI_API_KEY=your_api_key_here
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4

# DeepSeek配置（备选）
DEEPSEEK_API_KEY=your_deepseek_api_key_here
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-chat
```

**模拟模式**：如果暂时没有API密钥，可以启用模拟模式（功能受限）：

```bash
AI_MODE=simulation
```

### 6.2 本地服务配置

开发时可能需要配置本地服务的访问地址。

**配置文件位置**：`configs/config.yaml`

```yaml
app:
  name: "Huobao Drama Dev"
  version: "dev"
  debug: true

server:
  port: 5678
  host: "0.0.0.0"
  cors_origins:
    - "http://localhost:5173"
    - "http://localhost:3012"

storage:
  type: "local"
  local_path: "./data/storage"
  base_url: "http://localhost:5678/static"
```

---

## 第七章：开发工具配置

### 7.1 IDE配置

**Visual Studio Code推荐配置**：

创建`.vscode/settings.json`：

```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.formatOnSave": true
  },
  "[vue]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode",
    "editor.formatOnSave": true
  },
  "files.eol": "\n"
}
```

创建`.vscode/launch.json`用于调试：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}",
      "env": {},
      "args": []
    }
  ]
}
```

### 7.2 代码检查与格式化

**Go代码检查**：

```bash
# 安装gofmt（Go自带）
gofmt -w .

# 安装golint
go install golang.org/x/lint/golint@latest

# 运行检查
golint ./...

# 使用golangci-lint（推荐）
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

**前端代码检查**：

```bash
# 检查代码
npm run lint

# 修复可自动修复的问题
npm run lint -- --fix
```

---

## 第八章：故障排查

### 8.1 常见问题

**Go编译错误**：

```bash
# 清理缓存
go clean -cache

# 重新下载依赖
go mod download
go mod tidy

# 检查Go版本
go version
```

**前端依赖安装失败**：

```bash
# 清理node_modules
rm -rf node_modules package-lock.json

# 重新安装
npm install --registry=https://registry.npmmirror.com/
```

**端口被占用**：

```bash
# 查看占用端口的进程
lsof -i :5678

# 或使用
netstat -tulpn | grep 5678

# 结束进程
kill -9 <PID>
```

### 8.2 日志查看

**Go应用日志**：

开发模式下，日志直接输出到终端。使用`air`时，终端会显示所有日志。

**Docker日志**：

```bash
# 查看容器日志
docker compose logs -f app
```

---

## 练习任务

### 练习1：开发环境搭建（⭐）

**任务目标**：完成完整的开发环境搭建

**任务要求**：安装Go和Node.js运行环境，克隆项目代码并成功运行，使用热重载功能进行开发。

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| Go版本≥1.21 | ☐ |
| Node.js版本≥18 | ☐ |
| 项目成功克隆 | ☐ |
| `go run main.go`成功运行 | ☐ |
| `npm run dev`成功运行 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [代码结构详解](code_structure.md) | 项目架构和模块划分 | ⭐⭐ |
| [开发指南](development_guide.md) | 核心开发流程和规范 | ⭐⭐⭐ |
| [高级开发主题](advanced_dev.md) | 架构决策和最佳实践 | ⭐⭐⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本，包含完整开发环境搭建指南 |
