# 故障排查与常见问题

> **学习目标**：完成本手册学习后，你将能够独立诊断和解决火宝短剧平台的常见故障，掌握系统化的故障排查方法，具备配置监控和预防故障的能力。
> 
> **前置知识**：建议具备基本的Linux系统操作能力，熟悉Docker基本命令，了解日志查看方法。建议先完成《部署运维指南》的学习，理解系统的部署架构。
> 
> **难度等级**：⭐⭐（进阶）
> 
> **预计学习时间**：2-4小时

---

## 第一章：故障排查方法论

### 1.1 科学排查框架

故障排查是运维工作的核心技能之一。遵循系统化的排查方法可以大大提高问题定位的效率，避免盲目试错。本节介绍的排查框架适用于各种类型的系统故障，核心思想是**科学假设验证**：首先基于观察现象形成假设，然后设计验证方案，最后根据验证结果得出结论或形成新假设。

**故障排查的四个阶段**：

| 阶段 | 核心问题 | 主要动作 |
|------|----------|----------|
| 观察 | 发生了什么？何时开始？ | 收集信息、记录时间线 |
| 分析 | 可能的原因是什么？ | 分类讨论、形成假设 |
| 验证 | 如何确认假设？ | 设计验证方案、执行测试 |
| 解决 | 如何修复问题？ | 实施修复、验证效果 |

**阶段一：观察与信息收集**

故障排查的第一步是全面收集信息。很多时候，故障定位失败是因为信息收集不充分。收集的信息应该包括：

1. **时间信息**：故障开始的确切时间、持续时长、是否间歇性出现
2. **影响范围**：影响的用户数量、功能模块、地理位置
3. **错误表现**：错误消息、返回状态码、异常行为描述
4. **环境信息**：系统版本、配置变更历史、最近的操作记录

**阶段二：分析与假设形成**

基于收集到的信息，分析可能的原因。分析时应该考虑：
- 最近是否有配置变更？
- 是否有资源瓶颈（CPU、内存、磁盘、网络）？
- 依赖服务是否正常？
- 是否有类似的故障历史？

**阶段三：验证与排除**

设计验证方案时，应该遵循「**最小化验证**」原则：用最小的代价确认或排除某个假设。验证方案应该是具体的、可执行的、有明确结果的。

**阶段四：解决与预防**

故障修复后，需要验证修复效果。同时，应该分析根本原因，制定预防措施，避免类似故障再次发生。

### 1.2 常用诊断工具

有效的故障排查需要借助各种诊断工具。本节介绍常用的诊断工具及其使用方法。

**系统资源诊断工具**：

| 工具 | 功能 | 使用场景 |
|------|------|----------|
| top/htop | 实时查看CPU和内存使用 | 性能问题排查 |
| df -h | 查看磁盘空间使用 | 存储问题排查 |
| iostat | 查看I/O性能 | 磁盘性能问题 |
| netstat/ss | 查看网络连接状态 | 网络问题排查 |
| dmesg | 查看系统内核日志 | 系统级错误排查 |

**Docker诊断工具**：

| 工具 | 功能 | 使用场景 |
|------|------|----------|
| docker ps | 查看运行中的容器 | 服务状态检查 |
| docker logs | 查看容器日志 | 应用错误排查 |
| docker stats | 查看容器资源使用 | 性能问题排查 |
| docker inspect | 查看容器详细信息 | 配置问题排查 |
| docker exec | 在容器内执行命令 | 深入诊断 |

**应用诊断工具**：

| 工具 | 功能 | 使用场景 |
|------|------|----------|
| curl | HTTP请求测试 | API接口测试 |
| sqlite3 | SQLite数据库交互 | 数据库问题排查 |
| journalctl | systemd日志查看 | 系统服务日志 |

---

## 第二章：常见故障与解决方案

### 2.1 服务启动失败

**故障现象**：应用服务无法启动，Docker容器退出或systemd服务启动失败。

**故障排查流程**：

```bash
# 第一步：检查Docker容器状态
docker ps -a

# 查看停止的容器
docker ps -a --filter "status=exited"

# 查看具体容器的退出原因
docker logs huobao-drama-app --tail 100
```

**常见原因一：端口被占用**

```bash
# 检查5678端口的使用情况
netstat -tulpn | grep 5678
# 或使用ss命令
ss -tulpn | grep 5678

# 如果端口被占用，找到占用进程
lsof -i :5678

# 解决方案：停止占用端口的进程，或修改应用端口配置
```

**解决方案**：修改配置文件中的端口设置：

```yaml
# docker-compose.yml 或 configs/config.yaml
server:
  port: 5679  # 改为其他可用端口
```

**常见原因二：数据目录权限不足**

```bash
# 检查目录权限
ls -la /opt/huobao-drama/data/

# 检查目录所有者
ls -la /opt/huobao-drama/ | grep data

# 解决方案：设置正确的目录权限
sudo chown -R huobao:huobao /opt/huobao-drama/data
sudo chmod -R 755 /opt/huobao-drama/data
```

**常见原因三：配置文件错误**

```bash
# 检查配置文件语法
cat /opt/huobao-drama/configs/config.yaml

# 验证YAML语法
python3 -c "import yaml; yaml.safe_load(open('/opt/huobao-drama/configs/config.yaml'))"

# 如果有语法错误，Python会报错并指示错误位置
```

**常见原因四：数据库文件损坏**

```bash
# 检查SQLite数据库完整性
sqlite3 /opt/huobao-drama/data/drama_generator.db "PRAGMA integrity_check;"

# 如果发现损坏，尝试修复
sqlite3 /opt/huobao-drama/data/drama_generator.db ".recover" > /tmp/recovered.db
cp /tmp/recovered.db /opt/huobao-drama/data/drama_generator.db
```

**故障排查检查清单**：

| 步骤 | 检查项 | 操作命令 |
|------|--------|----------|
| 1 | 检查容器日志 | docker logs huobao-drama-app |
| 2 | 检查端口占用 | netstat -tulpn \| grep 5678 |
| 3 | 检查权限 | ls -la /opt/huobao-drama/data/ |
| 4 | 检查配置 | cat configs/config.yaml |
| 5 | 检查数据库 | sqlite3 db_file "PRAGMA integrity_check" |

### 2.2 API请求错误

**故障现象**：API请求返回错误状态码，如400、401、403、404、500等。

**按状态码分类的排查方法**：

**400 Bad Request（请求参数错误）**：

```bash
# 构造请求并查看详细错误
curl -v http://localhost:5678/api/v1/projects

# 预期输出应包含详细错误信息
```

**排查步骤**：

1. 检查请求URL是否正确
2. 检查请求参数格式（JSON、Form等）
3. 检查必填字段是否完整
4. 查看应用日志获取详细错误信息

```bash
# 查看相关错误日志
docker logs huobao-drama-app 2>&1 | grep -i "validation\|parameter"
```

**401 Unauthorized（未授权）**：

```bash
# 检查认证Token
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:5678/api/v1/projects

# Token是否有效？
# 1. 检查Token是否过期
# 2. 检查Token格式是否正确
# 3. 检查Token是否被撤销
```

**解决方案**：

```bash
# 重新获取Token（如果使用用户名密码登录）
curl -X POST http://localhost:5678/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your_user","password":"your_password"}'
```

**403 Forbidden（禁止访问）**：

此错误表示请求被服务器拒绝，通常是因为权限不足。

```bash
# 检查当前用户权限
curl -H "Authorization: Bearer TOKEN" http://localhost:5678/api/v1/users/me

# 查看返回的权限信息
```

**常见原因**：
- 用户被禁用
- 角色权限被修改
- 访问超出权限范围的资源

**404 Not Found（资源不存在）**：

```bash
# 检查资源ID是否正确
curl -H "Authorization: Bearer TOKEN" http://localhost:5678/api/v1/projects/non_existent_id

# 检查项目是否存在
```

**500 Internal Server Error（服务器内部错误）**：

这是最严重的错误类型，需要详细分析。

```bash
# 查看应用错误日志
docker logs huobao-drama-app 2>&1 | grep -A 10 "ERROR\|PANIC\|FATAL"

# 如果是Go应用，检查是否有panic
docker logs huobao-drama-app 2>&1 | grep -i "panic"
```

**常见原因与解决方案**：

| 原因 | 排查方法 | 解决方案 |
|------|----------|----------|
| 数据库连接失败 | 检查数据库日志 | 重启数据库服务 |
| 外部服务超时 | 检查AI服务配置 | 增加超时时间 |
| 内存不足 | docker stats检查 | 增加服务器内存 |
| 代码异常 | 查看应用日志 | 修复代码bug |

### 2.3 视频生成失败

**故障现象**：视频生成任务失败，任务状态显示为「failed」。

**故障排查流程**：

```bash
# 第一步：查询任务详情
curl http://localhost:5678/api/v1/tasks/{task_id}

# 查看任务状态和错误信息
# 特别关注 error 字段和 result 中的失败原因
```

**常见原因一：AI服务调用失败**：

```bash
# 检查AI服务配置
cat configs/config.yaml | grep -A 5 "ai:"

# 测试AI服务连通性
curl -X POST https://api.openai.com/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"test"}]}'
```

**解决方案**：检查API Key是否正确配置，网络连接是否正常。

**常见原因二：分镜图像不存在**：

```bash
# 检查分镜图像文件
ls -la /opt/huobao-drama/data/storage/storyboards/

# 检查文件路径配置
curl http://localhost:5678/api/v1/projects/{project_id}/storyboards/{storyboard_id}
```

**解决方案**：重新生成分镜图像，或检查存储路径配置。

**常见原因三：FFmpeg不可用**：

```bash
# 检查FFmpeg是否安装
ffmpeg -version

# 检查FFmpeg是否在PATH中
which ffmpeg

# 如果使用Docker，检查容器内FFmpeg
docker exec huobao-drama-app which ffmpeg
```

**解决方案**：确保FFmpeg正确安装并添加到系统PATH。

**故障日志分析**：

```bash
# 查看视频生成相关日志
docker logs huobao-drama-app 2>&1 | grep -i "video\|ffmpeg"

# 查看最近1小时的日志
docker logs --since 1h huobao-drama-app 2>&1 | grep -i error
```

### 2.4 数据库问题

**故障现象**：数据库相关错误，如连接失败、数据查询异常、数据损坏等。

**数据库连接问题**：

```bash
# 测试数据库连接
sqlite3 /opt/huobao-drama/data/drama_generator.db ".tables"

# 查看数据库文件大小
ls -lh /opt/huobao-drama/data/drama_generator.db

# 检查数据库文件完整性
sqlite3 /opt/huobao-drama/data/drama_generator.db "PRAGMA integrity_check;"
```

**数据库性能问题**：

```bash
# 开启慢查询日志（在配置文件中）
# log:
#   slow_query: true
#   slow_query_threshold: 1000  # 毫秒

# 查看当前查询
sqlite3 /opt/huobao-drama/data/drama_generator.db ".schema"
```

**常见SQLite问题与解决方案**：

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| database locked | 并发写入冲突 | 启用WAL模式 |
| database is locked | 文件权限问题 | 检查目录权限 |
| database disk image is malformed | 文件损坏 | 使用.recover恢复 |
| no such table | 表不存在 | 检查数据库迁移 |

**启用WAL模式解决并发问题**：

```bash
# 备份当前数据库
cp /opt/huobao-drama/data/drama_generator.db /opt/huobao-drama/data/drama_generator.db.backup

# 启用WAL模式
sqlite3 /opt/huobao-drama/data/drama_generator.db "PRAGMA journal_mode=WAL;"

# 验证WAL模式已启用
sqlite3 /opt/huobao-drama/data/drama_generator.db "PRAGMA journal_mode;"
# 应返回 wal
```

---

## 第三章：性能问题排查

### 3.1 系统资源瓶颈

**CPU使用率过高**：

```bash
# 查看实时CPU使用
top

# 查看CPU使用率最高的进程
ps aux --sort=-%cpu | head -10

# Docker容器CPU使用
docker stats --no-stream huobao-drama-app
```

**常见原因与处理**：

| 原因 | 排查方法 | 解决方案 |
|------|----------|----------|
| 视频生成任务 | docker stats检查容器CPU | 限制并发任务数 |
| 数据库查询 | 检查慢查询日志 | 优化查询或添加索引 |
| 恶意请求 | 查看访问日志 | 配置速率限制 |

**内存使用过高**：

```bash
# 查看内存使用
free -h

# 查看内存使用详情
cat /proc/meminfo

# 查看进程内存使用
ps aux --sort=-%mem | head -10
```

**处理方法**：

```bash
# Docker环境，查看容器内存使用
docker stats --no-stream huobao-drama-app

# 如果内存持续过高，考虑：
# 1. 增加服务器内存
# 2. 优化应用内存使用
# 3. 配置SWAP分区
```

**磁盘空间不足**：

```bash
# 查看磁盘使用情况
df -h

# 查看目录大小
du -sh /opt/huobao-drama/data/* | sort -hr | head -10

# 找出大文件
find /opt/huobao-drama -type f -size +100M -exec ls -lh {} \;
```

**磁盘清理操作**：

```bash
# 清理Docker未使用的资源
docker system prune -af --volumes

# 清理日志文件
find /opt/huobao-drama/logs -name "*.log" -mtime +7 -delete

# 清理临时文件
rm -rf /tmp/huobao-*

# 清理旧的视频文件
find /opt/huobao-drama/data/storage/videos -type f -mtime +30 -delete
```

### 3.2 API响应慢

**问题定位步骤**：

```bash
# 第一步：测试API响应时间
time curl -o /dev/null -s http://localhost:5678/api/v1/projects

# 第二步：检查数据库查询时间
# 在配置中启用慢查询日志

# 第三步：检查外部服务响应时间
# 如果使用AI服务，测试AI服务API响应

# 第四步：检查网络延迟
ping -c 5 api.openai.com  # 如果使用OpenAI
```

**常见原因与解决方案**：

| 原因 | 排查方法 | 解决方案 |
|------|----------|----------|
| 数据库查询慢 | 开启慢查询日志 | 添加索引、优化查询 |
| 外部服务慢 | 测试AI服务API | 更换服务商、增加超时 |
| 内存不足导致Swap | 检查内存和Swap使用 | 增加内存、优化代码 |
| 并发请求过多 | 查看并发连接数 | 配置限流 |

---

## 第四章：日志分析与调试

### 4.1 日志查看技巧

**Docker日志查看**：

```bash
# 实时查看日志
docker logs -f huobao-drama-app

# 查看最近100行
docker logs --tail 100 huobao-drama-app

# 按时间过滤
docker logs --since "2026-02-05T10:00:00" huobao-drama-app

# 查看包含特定内容的日志
docker logs huobao-drama-app 2>&1 | grep "ERROR"

# 使用JSON格式查看（如果配置了JSON日志）
docker logs huobao-drama-app 2>&1 | jq '.'
```

**日志级别筛选**：

```bash
# 查看ERROR级别日志
docker logs huobao-drama-app 2>&1 | jq 'select(.level == "error")'

# 查看WARN级别及以上日志
docker logs huobao-drama-app 2>&1 | jq 'select(.level | . == "error" or . == "warn")'
```

### 4.2 日志分析实践

**典型错误模式识别**：

| 错误模式 | 可能原因 | 排查方向 |
|----------|----------|----------|
|大量timeout错误| 服务响应慢/网络问题 | 检查服务状态和网络 |
|大量connection错误| 数据库连接/网络中断 | 检查数据库和网络 |
|panic错误| 代码异常 | 查看堆栈信息 |
|validation错误| 用户输入不合法 | 检查请求参数 |

**请求链路追踪**：

```bash
# 查看单个请求的完整日志
# 1. 在日志中搜索请求ID
docker logs huobao-drama-app 2>&1 | grep "req_id"

# 2. 或者使用trace ID（如果已配置）
curl -H "X-Trace-ID: test-123" http://localhost:5678/api/v1/projects
```

---

## 第五章：预防措施与监控

### 5.1 监控配置建议

**关键监控指标**：

| 指标 | 监控工具 | 告警阈值 | 采集频率 |
|------|----------|----------|----------|
| API可用性 | curl健康检查 | < 99.9% | 每分钟 |
| API响应时间P99 | 日志分析 | > 3秒 | 每5分钟 |
| 错误率 | 日志统计 | > 1% | 每5分钟 |
| CPU使用率 | 系统监控 | > 80% | 每分钟 |
| 内存使用率 | 系统监控 | > 85% | 每分钟 |
| 磁盘使用率 | 系统监控 | > 80% | 每小时 |
| 任务失败率 | API统计 | > 5% | 每10分钟 |

**健康检查接口**：

平台提供了健康检查接口，可以用于监控服务状态：

```bash
# 简单健康检查
curl http://localhost:5678/health

# 预期输出
{"status":"ok","timestamp":"2026-02-05T12:00:00Z"}

# 详细健康检查（如果配置）
curl http://localhost:5678/health/details
```

### 5.2 告警配置

**告警通知配置**（使用Prometheus Alertmanager为例）：

```yaml
# alertmanager.yml

global:
  resolve_timeout: 5m

route:
  group_by: ['alertname']
  receiver: 'default'
  routes:
    - match:
        severity: critical
      receiver: 'critical-alerts'
    - match:
        severity: warning
      receiver: 'warning-alerts'

receivers:
- name: 'default'
  email_configs:
  - to: 'ops@example.com'
    send_resolved: true
- name: 'critical-alerts'
  email_configs:
  - to: 'critical@example.com'
    send_resolved: true
  webhook_configs:
  - url: 'https://oapi.dingtalk.com/robot/send?access_token=XXX'
```

### 5.3 定期巡检

**每日巡检任务**：

| 任务 | 命令 | 检查标准 |
|------|------|----------|
| 服务状态检查 | systemctl status huobao-drama | active (running) |
| 磁盘空间检查 | df -h /opt | < 80% |
| 错误日志检查 | journalctl -u huobao-drama -p err --since yesterday | 无新ERROR |
| 健康检查 | curl http://localhost:5678/health | {"status":"ok"} |

**每周巡检任务**：

| 任务 | 命令 | 检查标准 |
|------|------|----------|
| 备份验证 | ls -la /opt/huobao-drama/backups/ | 备份存在 |
| SSL证书检查 | certbot certificates | 有效期>30天 |
| 安全更新 | apt list --upgradable | 无高危更新 |
| 性能趋势 | 查看监控图表 | 无明显恶化 |

---

## 第六章：故障恢复预案

### 6.1 故障级别定义

| 级别 | 定义 | 响应时间 | 通知范围 |
|------|------|----------|----------|
| P0（紧急） | 服务完全不可用 | 立即响应 | 全员通知 |
| P1（严重） | 核心功能不可用 | 15分钟内 | 技术团队 |
| P2（一般） | 非核心功能异常 | 1小时内 | 相关人员 |
| P3（轻微） | 轻微问题 | 24小时内 | 可延迟 |

### 6.2 故障恢复流程

**P0级别故障恢复流程**：

```
故障发现（P0）
    │
    ▼
立即通知（1分钟内）
    │
    ▼
├─→ 服务停机？
│   │
│   ├─→ 是：执行回滚操作 → 验证恢复
│   │
│   └─→ 否：定位问题根因 → 修复 → 验证
│
▼
恢复验证（30分钟内）
    │
    ▼
├─→ 服务恢复正常：发送恢复通知 → 进入故障复盘
│
└─→ 无法恢复：启动备用方案 → 人工介入
```

### 6.3 回滚操作

**Docker回滚**：

```bash
# 查看历史镜像
docker images huobao-drama

# 回滚到上一版本
docker compose down
docker compose pull
docker compose up -d

# 或者回滚到指定版本
docker compose down
docker tag huobao-drama:v1.0.3 huobao-drama:latest
docker compose up -d
```

**数据库回滚**：

```bash
# 停止服务
docker compose stop app

# 恢复备份
cp /opt/huobao-drama/backups/db_20260204.sqlite3.gz /opt/huobao-drama/data/drama_generator.db.gz
gunzip /opt/huobao-drama/data/drama_generator.db.gz

# 验证完整性
sqlite3 /opt/huobao-drama/data/drama_generator.db "PRAGMA integrity_check;"

# 启动服务
docker compose start app
```

---

## 练习任务

### 练习1：故障诊断实践（⭐⭐）

**任务目标**：模拟常见故障并进行诊断

**任务要求**：

1. 模拟「服务启动失败」故障
2. 按照排查流程逐步诊断
3. 记录排查过程和解决方案
4. 总结故障预防措施

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| 成功模拟故障场景 | ☐ |
| 遵循排查流程 | ☐ |
| 准确定位故障原因 | ☐ |
| 成功恢复服务 | ☐ |
| 总结预防措施 | ☐ |

---

## 常见问题速查表

| 问题 | 可能原因 | 快速解决方案 |
|------|----------|--------------|
| 服务无法启动 | 端口被占用 | 停止占用进程或改端口 |
| API返回401 | Token过期 | 重新登录获取Token |
| 视频生成失败 | AI服务配置错误 | 检查API Key配置 |
| 数据库连接失败 | 数据库文件损坏 | 修复或恢复数据库 |
| 页面加载慢 | 磁盘空间不足 | 清理磁盘空间 |
| 登录失败 | 用户被禁用 | 联系管理员 |
| 上传失败 | 文件大小超限 | 压缩或分割文件 |
| Docker容器退出 | OOM（内存不足） | 增加服务器内存 |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [快速入门指南](quick-start.md) | 15分钟快速上手 | ⭐ |
| [用户手册](user-manual.md) | 平台功能的完整操作指南 | ⭐⭐ |
| [API接口文档](api-reference.md) | 开发者API参考手册 | ⭐⭐⭐ |
| [架构设计文档](architecture.md) | 系统架构深度解析 | ⭐⭐⭐⭐ |
| [部署运维指南](deployment.md) | 生产环境部署和运维指南 | ⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-05 | 初始版本，包含完整故障排查指南 |
