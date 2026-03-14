# Docker 主机访问配置指南

> **学习目标**：完成本指南后，你将能够配置 Docker 容器访问宿主机服务（如 Ollama 本地模型）。
>
> **前置知识**：建议具备基本的 Docker 和 Linux 命令行知识。
>
> **难度等级**：⭐（入门）
>
> **预计学习时间**：15-30 分钟

---

## 概述

### 为什么需要主机访问

在 Docker 容器化部署场景中，容器默认运行在隔离的网络环境中。如果需要在容器内访问宿主机上运行的服务（如本地 Ollama 推理引擎），需要进行特殊配置。

**典型场景**：

- 使用本地运行的 Ollama 作为 AI 服务提供商
- 访问宿主机上的数据库或其他服务
- 调用本地文件系统的资源

### 技术原理

Docker 容器访问宿主机有两种主要方式：

| 方式 | 原理 | 适用平台 | 难度 |
|------|------|----------|------|
| `host.docker.internal` | Docker Desktop 内置 DNS 解析 | macOS/Windows | ⭐ |
| `--add-host` 参数 | 直接映射主机 IP | Linux | ⭐⭐ |
| Docker 网络桥接 | 自定义网络配置 | 所有平台 | ⭐⭐⭐ |

---

## 各平台配置

### macOS（Docker Desktop）

macOS 上的 Docker Desktop 自动提供 `host.docker.internal` 解析，无需额外配置。

**验证方法**：

```bash
# 进入容器
docker exec -it huobao-drama-app-1 /bin/sh

# 测试主机访问
ping -c 3 host.docker.internal
nslookup host.docker.internal
```

**预期输出**：

```
PING host.docker.internal (192.168.65.2): 56 data bytes
64 bytes from 192.168.65.2: icmp_seq=0 ttl=254 time=2.385 ms
```

### Windows（Docker Desktop）

与 macOS 类似，Docker Desktop 自动支持 `host.docker.internal`。

**注意事项**：

- 确保 Docker Desktop 已启用 "Expose daemon on tcp://localhost:2375" 选项
- 防火墙可能需要放行 2375 端口

### Linux

Linux 系统需要使用 `--add-host` 参数手动添加主机映射。

**方法一：运行时指定**：

```bash
docker run -d \
  --name huobao-drama \
  -p 5678:5678 \
  --add-host=host.docker.internal:host-gateway \
  -v $(pwd)/data:/app/data \
  huobao/huobao-drama:latest
```

**方法二：docker-compose 配置**（推荐）：

```yaml
# docker-compose.yml
services:
  app:
    image: huobao/huobao-drama:latest
    ports:
      - "5678:5678"
    extra_hosts:
      - "host.docker.internal:host-gateway"
    volumes:
      - ./data:/app/data
```

**验证配置**：

```bash
# 重新构建并启动
docker compose down
docker compose up -d

# 进入容器验证
docker exec -it huobao-drama-app-1 /bin/sh

# 测试解析
cat /etc/hosts | grep host.docker.internal
# 应该显示：host.docker.internal x.x.x.x
```

---

## Ollama 访问配置

### 场景描述

在容器内访问宿主机上运行的 Ollama 服务，用于本地大语言模型推理。

### 宿主机 Ollama 配置

**步骤 1**：确保 Ollama 监听所有接口

```bash
# 宿主机上执行
export OLLAMA_HOST=0.0.0.0:11434
ollama serve
```

**说明**：默认情况下 Ollama 只监听 `127.0.0.1`，需要显式配置为监听所有网络接口。

**步骤 2**：验证 Ollama 服务

```bash
# 宿主机测试
curl http://localhost:11434/api/version
```

**预期输出**：

```json
{
  "version": "0.5.41",
  "date": "2024-08-04",
  "OLLAMA_HOST": "0.0.0.0:11434"
}
```

### 容器内 Ollama 配置

**步骤 1**：配置 AI 服务 URL

在火宝短剧的配置文件中，设置 Ollama 的 Base URL：

```yaml
# configs/config.yaml
ai:
  default_text_provider: "openai"  # Ollama 兼容 OpenAI API 格式
  providers:
    openai:
      base_url: "http://host.docker.internal:11434/v1"
      api_key: "ollama"  # Ollama 不需要真实密钥
      model: "qwen2.5:latest"  # 使用本地模型
```

**步骤 2**：验证容器内访问

```bash
# 进入容器
docker exec -it huobao-drama-app-1 /bin/sh

# 测试 Ollama 连接
curl http://host.docker.internal:11434/api/tags

# 预期输出：列出可用的模型
```

### 常见问题排查

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| 连接超时 | Ollama 未监听所有接口 | 设置 `OLLAMA_HOST=0.0.0.0:11434` |
| 连接拒绝 | 防火墙阻止端口 | 开放 11434 端口 |
| DNS 解析失败 | Linux 未配置 host-gateway | 使用 `--add-host` 参数 |
| 模型不存在 | 本地未下载模型 | 宿主机运行 `ollama pull qwen2.5` |

---

## 完整配置示例

### docker-compose.yml 完整配置

```yaml
version: '3.8'

services:
  app:
    image: huobao/huobao-drama:latest
    container_name: huobao-drama-app
    restart: unless-stopped
    ports:
      - "5678:5678"
    environment:
      - TZ=Asia/Shanghai
    volumes:
      - ./data:/app/data
    extra_hosts:
      - "host.docker.internal:host-gateway"
    networks:
      - huobao-network

networks:
  huobao-network:
    driver: bridge
```

### Ollama 完整启动流程

```bash
# 1. 宿主机下载并启动 Ollama
export OLLAMA_HOST=0.0.0.0:11434
nohup ollama serve > ollama.log 2>&1 &

# 2. 下载模型（以 Qwen2.5 为例）
ollama pull qwen2.5:latest

# 3. 验证 Ollama 运行
curl http://localhost:11434/api/tags

# 4. 启动火宝短剧
docker compose up -d

# 5. 验证容器内访问
docker exec -it huobao-drama-app-1 \
  curl http://host.docker.internal:11434/api/tags
```

---

## 性能注意事项

### 网络延迟

通过 `host.docker.internal` 访问会有轻微网络延迟，对于大多数场景影响可忽略。

**优化建议**：

- 如果延迟明显，考虑使用 Docker 网络桥接模式
- 高并发场景下，建议将 Ollama 也容器化部署在同一网络中

### 安全考虑

**风险**：暴露 Ollama 服务端口可能导致未授权访问

**防护措施**：

1. **网络隔离**：使用 Docker 网络限制访问范围
2. **认证**：配置 Ollama 认证（如果版本支持）
3. **防火墙**：限制可访问 Ollama 的 IP 范围

```bash
# 示例：仅允许 Docker 网络访问
iptables -A INPUT -s 172.0.0.0/8 -p tcp --dport 11434 -j ACCEPT
iptables -A INPUT -p tcp --dport 11434 -j DROP
```

---

## 练习任务

### 练习 1：配置 Linux 主机访问 ⭐

**任务目标**：在 Linux 服务器上配置 Docker 容器访问宿主机服务

**任务要求**：

1. 编辑 `docker-compose.yml`，添加 `extra_hosts` 配置
2. 重启服务并验证
3. 在容器内成功执行 `curl http://host.docker.internal:11434/api/tags`

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| `extra_hosts` 配置正确 | ☐ |
| 容器能解析 `host.docker.internal` | ☐ |
| Ollama API 返回有效响应 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [安装指南](getting-started/installation.md) | Docker 完整安装配置 | ⭐ |
| [快速入门](getting-started/quick-start.md) | 15 分钟创建第一个项目 | ⭐ |
| [部署指南](deployment.md) | 生产环境部署详细指南 | ⭐⭐ |
| [故障排查](troubleshooting.md) | 常见问题和解决方案 | ⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本 |
