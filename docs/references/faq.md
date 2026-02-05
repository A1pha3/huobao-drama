# 常见问题解答

> 本文档收录了火宝短剧使用者最常遇到的问题及其解决方案。

---

## 第一部分：安装部署问题

### Q1：Docker安装失败怎么办？

**问题描述**：在安装Docker时遇到各种错误，导致安装失败。

**可能原因**：

- 系统版本不兼容
- 网络连接问题
- 权限不足
- 磁盘空间不足

**解决方案**：

1. **检查系统要求**：
   - Windows 10/11专业版/企业版（家庭版不支持Hyper-V）
   - macOS 10.15或更高版本
   - Linux内核版本3.10或更高

2. **Windows用户**：
   ```powershell
   # 检查是否启用了Hyper-V
   Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V
   
   # 如果未启用，以管理员身份运行PowerShell
   Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
   ```

3. **macOS用户**：
   ```bash
   # 检查Apple Silicon还是Intel芯片
   uname -m
   # arm64 = Apple Silicon
   # x86_64 = Intel
   ```

4. **Linux用户**：
   ```bash
   # 检查系统架构
   uname -m
   
   # 对于Ubuntu 20.04+
   sudo apt update
   sudo apt install docker.io docker-compose
   ```

### Q2：端口被占用怎么办？

**问题描述**：启动服务时提示端口已被占用。

**解决方案**：

```bash
# Windows
netstat -ano | findstr :5678
taskkill /PID <PID> /F

# macOS/Linux
lsof -i :5678
kill -9 <PID>

# 或者修改配置文件使用其他端口
# 在.env文件中修改
SERVER_PORT=5679
```

### Q3：AI服务配置错误

**问题描述**：配置AI服务后无法正常使用AI功能。

**检查步骤**：

1. **验证API Key是否正确**：
   ```bash
   # 测试OpenAI API
   curl -X POST "https://api.openai.com/v1/chat/completions" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_API_KEY" \
     -d '{"model":"gpt-4","messages":[{"role":"user","content":"test"}]}'
   ```

2. **检查网络连接**：
   ```bash
   # 测试API端点连通性
   ping api.openai.com
   telnet api.openai.com 443
   ```

3. **验证配置文件格式**：
   ```yaml
   # 确保格式正确（注意缩进）
   ai:
     default_provider: "openai"
     providers:
       openai:
         api_key: "sk-xxx"
         base_url: "https://api.openai.com/v1"
   ```

---

## 第二部分：功能使用问题

### Q4：AI生成的剧本不满意怎么办？

**问题描述**：使用AI生成剧本时，结果不符合预期。

**优化建议**：

1. **改进提示词**：
   ```
   ❌ 模糊的提示词
   "生成一个爱情故事"
   
   ✅ 具体的提示词
   "生成一个现代都市爱情故事，
    男女主角在咖啡店相遇，
    女生是咖啡师，男生是作家，
    故事风格轻松幽默，
    需要包含3个场景：相遇、相知、表白"
   ```

2. **使用否定约束**：
   ```
   "生成一个爱情故事，
    不要包含出轨情节，
    不要有悲剧结局，
    不要出现暴力元素"
   ```

3. **多次尝试和迭代**：
   - 先生成故事大纲
   - 审核大纲后要求AI细化每个场景
   - 最后要求AI润色对白

4. **使用模板**：
   选择合适的模板作为起点，然后进行修改

### Q5：角色形象不一致怎么办？

**问题描述**：同一个角色在不同分镜中的形象差异很大。

**解决方案**：

1. **详细描述角色特征**：
   ```
   角色名称：小红
   核心特征：
   - 黑色长发，及腰
   - 圆脸，眼睛大而有神
   - 身高165cm，身材匀称
   - 常穿淡蓝色连衣裙
   
   每次生成时强调：
   "保持角色小红的特征：
    黑色长发、圆脸、蓝色连衣裙"
   ```

2. **使用角色参考**：
   - 上传一张参考图片作为基础
   - 在提示词中引用参考图的特征

3. **批量生成**：
   - 先生成角色的多张图片
   - 选择最满意的一张作为标准
   - 后续生成时引用该图片的ID

### Q6：视频生成失败

**问题描述**：视频生成任务失败或生成时间过长。

**排查步骤**：

1. **检查分镜图像**：
   - 确认所有分镜都有对应的图像
   - 检查图像文件是否完整

2. **检查资源**：
   ```bash
   # 检查磁盘空间
   df -h
   
   # 检查内存
   free -h
   
   # 检查CPU使用率
   top
   ```

3. **查看错误日志**：
   ```bash
   # Docker环境
   docker logs huobao-drama-app 2>&1 | grep -i error
   
   # 本地环境
   tail -f logs/app.log
   ```

4. **降低生成参数**：
   - 降低分辨率（720p代替1080p）
   - 缩短视频时长
   - 减少并发任务数

---

## 第三部分：性能问题

### Q7：系统运行缓慢

**问题描述**：平台响应变慢，操作有延迟。

**优化建议**：

1. **清理缓存**：
   ```bash
   # 清理Docker缓存
   docker system prune -af
   
   # 清理项目缓存
   rm -rf tmp/*
   ```

2. **检查资源使用**：
   ```bash
   # Docker环境
   docker stats
   
   # 本地环境
   htop  # Linux/macOS
   Task Manager  # Windows
   ```

3. **优化数据库**：
   ```bash
   # SQLite优化
   sqlite3 data/drama_generator.db "PRAGMA integrity_check;"
   sqlite3 data/drama_generator.db "VACUUM;"
   ```

4. **重启服务**：
   ```bash
   # Docker环境
   docker compose restart
   
   # 本地环境
   # 重启应用进程
   ```

### Q8：内存不足

**问题描述**：系统提示内存不足或频繁崩溃。

**解决方案**：

1. **增加虚拟内存（Windows）**：
   ```
   控制面板 → 系统 → 高级系统设置 → 性能 → 设置 → 高级 → 虚拟内存 → 更改
   ```

2. **增加交换空间（Linux）**：
   ```bash
   sudo fallocate -l 4G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

3. **限制Docker内存**：
   ```yaml
   # docker-compose.yml
   services:
     app:
       mem_limit: 4g
   ```

---

## 第四部分：数据问题

### Q9：如何备份和恢复数据？

**备份步骤**：

```bash
# 1. 停止服务
docker compose down

# 2. 备份数据库
cp data/drama_generator.db data/drama_generator.db.backup
tar -czf data_backup_$(date +%Y%m%d).tar.gz data/

# 3. 备份配置文件
cp configs/config.yaml configs/config.yaml.backup

# 4. 重启服务
docker compose up -d
```

**恢复步骤**：

```bash
# 1. 停止服务
docker compose down

# 2. 恢复数据库
cp data/drama_generator.db.backup data/drama_generator.db

# 3. 验证完整性
sqlite3 data/drama_generator.db "PRAGMA integrity_check;"

# 4. 重启服务
docker compose up -d
```

### Q10：数据库损坏怎么办？

**问题描述**：SQLite数据库文件损坏，无法正常访问。

**修复步骤**：

```bash
# 1. 备份损坏的数据库
cp data/drama_generator.db data/drama_generator.db.corrupt

# 2. 尝试恢复
sqlite3 data/drama_generator.db ".recover" > recovered.db

# 3. 替换数据库
mv recovered.db data/drama_generator.db

# 4. 验证
sqlite3 data/drama_generator.db "PRAGMA integrity_check;"

# 5. 如果恢复失败，可能需要从备份恢复
```

### Q11：如何导出项目？

**导出步骤**：

1. 在项目列表页面，选择要导出的项目
2. 点击「更多」→「导出」
3. 选择导出格式（ZIP压缩包或单独文件）
4. 选择包含的内容：
   - 剧本数据
   - 角色信息
   - 分镜图像
   - 视频文件
5. 点击「确认导出」
6. 下载导出的文件

---

## 第五部分：网络问题

### Q12：无法访问外部AI服务

**问题描述**：无法连接到OpenAI等服务。

**排查步骤**：

1. **检查网络连接**：
   ```bash
   # 测试网络连通性
   ping api.openai.com
   
   # 测试HTTPS连接
   curl -v https://api.openai.com
   ```

2. **检查代理设置**：
   ```bash
   # 查看环境变量
   env | grep -i proxy
   
   # 如果需要代理，在.env中配置
   HTTP_PROXY=http://proxy.example.com:7890
   HTTPS_PROXY=http://proxy.example.com:7890
   ```

3. **检查防火墙**：
   ```bash
   # Windows
   Windows Defender Firewall
   
   # macOS
   sudo /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate
   
   # Linux
   sudo ufw status
   ```

### Q13：国内访问AI服务慢

**问题描述**：国内网络访问OpenAI等服务速度慢。

**解决方案**：

1. **配置代理**：
   ```yaml
   # .env
   HTTP_PROXY=http://127.0.0.1:7890
   HTTPS_PROXY=http://127.0.0.1:7890
   ```

2. **使用国内镜像服务**：
   ```yaml
   # 使用OpenAI的代理镜像
   OPENAI_BASE_URL=https://api.openai-proxy.com/v1
   ```

3. **更换AI服务商**：
   - 使用DeepSeek
   - 使用豆包
   - 使用本地部署的模型

---

## 第六部分：开发相关

### Q14：如何进行本地开发？

**开发环境搭建**：

1. **安装依赖**：
   ```bash
   # 安装Go
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/air-verse/air@latest
   
   # 安装Node.js
   nvm install --lts
   nvm use --lts
   
   # 安装前端依赖
   cd web
   npm install
   ```

2. **启动开发服务器**：
   ```bash
   # 终端1：后端
   air
   
   # 终端2：前端
   cd web
   npm run dev
   ```

3. **访问开发环境**：
   - 前端：http://localhost:5173
   - 后端API：http://localhost:5678

### Q15：如何运行测试？

**测试命令**：

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./domain/...
go test ./application/...
go test ./api/...

# 运行测试并显示覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 运行基准测试
go test -bench=. -benchmem
```

---

## 第七部分：其他问题

### Q16：如何获取帮助？

**获取帮助的方式**：

1. **查阅文档**：
   - 快速入门指南
   - 用户手册
   - API参考
   - 故障排查文档

2. **GitHub Issues**：
   - https://github.com/chatfire-AI/huobao-drama/issues
   - 搜索已有问题
   - 创建新问题

3. **社区讨论**：
   - GitHub Discussions
   - 项目Wiki

4. **联系维护者**：
   - 邮箱：见项目首页

### Q17：如何报告Bug？

**Bug报告模板**：

```markdown
## Bug描述
[简要描述问题]

## 复现步骤
1. [步骤1]
2. [步骤2]
3. [步骤3]

## 预期行为
[期望的结果]

## 实际行为
[实际的结果]

## 环境信息
- 操作系统：
- Go版本：
- Docker版本（如果使用）：
- 浏览器（如果是前端问题）：

## 错误日志
[粘贴相关错误日志]

## 截图（如适用）
[添加截图]
```

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本 |
