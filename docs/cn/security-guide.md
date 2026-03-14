# 安全配置指南

> **学习目标**：完成本指南后，你将能够配置火宝短剧的安全设置，保护系统免受常见安全威胁。
>
> **前置知识**：建议具备基本的 Linux 系统管理和网络安全知识。
>
> **难度等级**：⭐⭐（进阶）
>
> **预计学习时间**：1-2 小时

---

## 概述

### 安全重要性

在生产环境中，系统的安全性至关重要。本指南涵盖以下安全领域：

| 领域 | 内容 | 重要性 |
|------|------|--------|
| 认证授权 | JWT 配置、密码策略 | ⭐⭐⭐ |
| 网络安全 | HTTPS 配置、CORS、Rate Limiting | ⭐⭐⭐ |
| 数据安全 | 加密存储、敏感信息保护 | ⭐⭐⭐ |
| 运行安全 | 沙箱隔离、资源限制 | ⭐⭐ |

---

## 认证与授权

### JWT 配置

火宝短剧使用 JWT（JSON Web Token）进行用户认证。

**配置项**：

```yaml
# configs/config.yaml
auth:
  jwt_secret: "${YOUR_STRONG_SECRET_KEY}"  # 必须使用强密钥
  token_expiry: 86400  # 24 小时（秒）
  refresh_token_expiry: 604800  # 7 天（秒）
```

**强密钥生成**：

```bash
# 使用 OpenSSL 生成
openssl rand -base64 32

# 或使用随机字符串（至少 32 字符）
# 建议：包含大小写字母、数字、特殊字符
```

**JWT 安全最佳实践**：

| 实践 | 说明 |
|------|------|
| 使用强密钥 | 长度至少 256 位 |
| 设置合理过期时间 | 不建议超过 24 小时 |
| 启用令牌刷新 | 平衡安全性与用户体验 |
| 记录令牌使用 | 便于审计和异常检测 |

### 密码策略

**推荐配置**：

```yaml
auth:
  password:
    min_length: 12  # 最小长度
    require_uppercase: true  # 必须包含大写字母
    require_lowercase: true  # 必须包含小写字母
    require_numbers: true    # 必须包含数字
    require_special: true    # 必须包含特殊字符
    max_age: 90             # 最大使用天数（90 天后强制修改）
```

---

## 网络安全

### HTTPS 配置

**生产环境必须使用 HTTPS**。以下是常见配置方式：

**方式一：Nginx 反向代理（推荐）**：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书配置
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # SSL 安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location / {
        proxy_pass http://localhost:5678;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**方式二：Let's Encrypt 免费证书**：

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期测试
sudo certbot renew --dry-run
```

### CORS 配置

配置跨域资源共享策略：

```yaml
# configs/config.yaml
server:
  cors_origins:
    - "https://your-domain.com"  # 生产环境只允许自有域名
    # - "http://localhost:3012"   # 开发环境可添加
  cors_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
  cors_headers:
    - "Authorization"
    - "Content-Type"
  cors_max_age: 86400  # 24 小时
```

**安全建议**：

- ❌ 不要使用通配符 `*`
- ✅ 生产环境只允许自有域名
- ✅ 限制允许的 HTTP 方法

### Rate Limiting（限流）

防止 API 被滥用和 DDoS 攻击：

```go
// middleware/rate_limit.go 示例
package middleware

import (
    "golang.org/x/time/rate"
    "github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
    // 限制：每秒 10 个请求，突发 20
    limiter := rate.NewLimiter(10, 20)
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "Too Many Requests",
                "message": "请求频率超限，请稍后重试",
            })
            return
        }
        c.Next()
    }
}
```

**推荐的限流策略**：

| 端点类型 | 限制 |
|----------|------|
| 登录/注册 | 5 次/分钟/IP |
| 普通 API | 100 次/分钟/用户 |
| 视频生成 | 10 次/小时/用户 |
| 资源上传 | 50 次/小时/用户 |

---

## 数据安全

### 敏感信息保护

**永不暴露的信息**：

| 类型 | 示例 | 处理方式 |
|------|------|----------|
| 密码 | admin123 | BCrypt 哈希存储 |
| API 密钥 | sk-xxxxx | 环境变量加密存储 |
| JWT 密钥 | xxxxxx | 强密钥，定期轮换 |
| 数据库密码 | xxxxxx | 独立配置，仅限必要人员访问 |

**环境变量配置**：

```bash
# 创建 .env 文件（永不提交到版本控制）
cp .env.example .env

# 编辑配置
vim .env
```

```env
# .env 示例（敏感信息使用环境变量）
DATABASE_PASSWORD=your_strong_db_password
JWT_SECRET=$(openssl rand -base64 32)
AI_API_KEY=sk-your-api-key-here
ENCRYPTION_KEY=$(openssl rand -base64 32)
```

### 数据库安全

**SQLite 安全建议**：

1. **文件权限**：限制数据库文件访问
   ```bash
   chmod 600 data/drama_generator.db
   ```

2. **备份加密**：备份文件加密存储
   ```bash
   # 加密备份
   tar czf - data/ | openssl enc -aes-256-cbc -salt -out backup_$(date +%Y%m%d).tar.gz.enc
   
   # 解密恢复
   openssl enc -d -aes-256-cbc -in backup_xxx.tar.gz.enc | tar xzf -
   ```

3. **WAL 模式安全**：
   ```go
   // 确保 WAL 模式正确配置
   db, err := gorm.Open(sqlite.Open("file:security.db?cache=shared"), &gorm.Config{
       ConnPool: db,
   })
   ```

### 文件上传安全

**限制上传文件类型**：

```go
// 上传安全检查
func validateUpload(filename string, allowedExts map[string]bool) error {
    ext := filepath.Ext(filename)
    if !allowedExts[ext] {
        return fmt.Errorf("不允许的文件类型: %s", ext)
    }
    
    // 检查 MIME 类型
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    buffer := make(512)
    _, err = file.Read(buffer)
    if err != nil {
        return err
    }
    
    contentType := http.DetectContentType(buffer)
    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png": true,
        "image/gif": true,
        "video/mp4": true,
    }
    
    if !allowedTypes[contentType] {
        return fmt.Errorf("不允许的文件类型: %s", contentType)
    }
    
    return nil
}
```

**文件大小限制**：

```yaml
# configs/config.yaml
upload:
  max_file_size: 104857600  # 100MB
  max_image_size: 10485760  # 10MB
  max_video_size: 524288000 # 500MB
  allowed_image_types:
    - ".jpg"
    - ".jpeg"
    - ".png"
    - ".gif"
  allowed_video_types:
    - ".mp4"
    - ".mov"
    - ".avi"
```

---

## 运行安全

### 容器安全

**Docker 安全最佳实践**：

```yaml
# docker-compose.yml 安全配置
services:
  app:
    image: huobao/huobao-drama:latest
    security_opt:
      - no-new-privileges:true
    read_only: true  # 只读文件系统
    cap_drop:
      - ALL
    tmpfs:
      - /tmp:size=10M,mode=1777
    
    # 非 root 用户运行
    user: "1000:1000"
```

**安全加固清单**：

| 配置项 | 推荐值 | 说明 |
|--------|--------|------|
| `no-new-privileges` | true | 防止进程获取额外权限 |
| `read_only` | true | 限制文件系统写入 |
| `cap_drop` | ALL | 移除所有 Linux 能力 |
| `user` | 非 root | 以非 root 用户运行 |
| `tmpfs` | /tmp | 临时文件隔离 |

### 资源限制

```yaml
services:
  app:
    image: huobao/huobao-drama:latest
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M
```

---

## 安全监控

### 日志审计

**配置安全日志**：

```go
// middleware/security_logger.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "log"
    "time"
)

func SecurityLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // 记录安全相关事件
        log.Printf("[SECURITY] %s | %s | %s | %s",
            c.ClientIP(),
            c.Request.Method,
            c.Request.URL.Path,
            c.Request.UserAgent(),
        )
        
        c.Next()
        
        // 记录异常访问
        if c.Writer.Status() >= 400 {
            log.Printf("[SECURITY ALERT] Status: %d | IP: %s | Path: %s",
                c.Writer.Status(),
                c.ClientIP(),
                c.Request.URL.Path,
            )
        }
    }
}
```

**监控告警规则**：

| 事件 | 阈值 | 告警方式 |
|------|------|----------|
| 登录失败 | 10 次/分钟 | 邮件/钉钉 |
| 401/403 错误 | 50 次/分钟 | 邮件 |
| API 限流触发 | 持续触发 | 告警 |
| 异常 IP 访问 | 新 IP 大量请求 | 告警 |

### 入侵检测

**常见攻击模式检测**：

| 攻击类型 | 检测方法 | 响应 |
|----------|----------|------|
| SQL 注入 | 参数化查询 + 日志监控 | 限流 + 告警 |
| XSS | 输入过滤 + CSP 头 | 阻断 + 告警 |
| CSRF | Token 验证 | 阻断 |
|暴力破解|登录限流|账号锁定|

---

## 合规检查清单

### 部署前检查

| 检查项 | 状态 | 说明 |
|--------|------|------|
| HTTPS 启用 | ☐ | 生产环境强制 |
| JWT 强密钥 | ☐ | 32+ 字符随机密钥 |
| API 限流配置 | ☐ | 防止滥用 |
| CORS 严格配置 | ☐ | 不使用通配符 |
| 敏感信息加密 | ☐ | API 密钥、密码 |
| 日志脱敏 | ☐ | 不记录敏感数据 |
| 备份加密 | ☐ | 备份文件加密 |
| 容器安全加固 | ☐ | 非 root 运行 |

### 定期维护任务

| 周期 | 任务 |
|------|------|
| 每日 | 检查异常登录日志 |
| 每周 | 轮换 API 密钥 |
| 每月 | 安全补丁更新 |
| 每季度 | 安全审计 |

---

## 练习任务

### 练习 1：配置 HTTPS ⭐⭐

**任务目标**：为火宝短剧配置 HTTPS

**任务要求**：

1. 获取 SSL 证书（Let's Encrypt）
2. 配置 Nginx 反向代理
3. 启用安全响应头

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| HTTPS 访问正常 | ☐ |
| 证书自动续期配置 | ☐ |
| HSTS 头启用 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [部署指南](deployment.md) | 生产环境部署详细指南 | ⭐⭐ |
| [故障排查](troubleshooting.md) | 安全相关问题诊断 | ⭐⭐ |
| [API 参考](api-reference.md) | 认证相关 API | ⭐⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本 |
