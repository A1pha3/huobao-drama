# 火宝短剧分镜模块改进建议

## 文档概述

本文档汇总了火宝短剧分镜模块的所有改进建议，按照优先级排序。每条建议包含问题描述、影响分析、改进方案、工作量估算和风险评估，方便技术负责人做出实施决策。

**适用读者**：技术负责人、项目经理、架构师

**决策类型**：是否批准实施改进方案

**预计阅读时间**：20-30分钟

---

## 改进建议速览

| 优先级 | 问题 | 影响 | 估算工作量 | 决策状态 |
|--------|------|------|------------|----------|
| **P0** | 缺少单元测试 | 质量风险 | 3-5人天 | ☐ 待决定 |
| **P1** | AI输出容错不足 | 稳定性 | 2-3人天 | ☐ 待决定 |
| **P2** | 单文件过大 | 可维护性 | 1-2人天 | ☐ 待决定 |
| **P3** | 超长剧本处理 | 功能完整性 | 5-8人天 | ☐ 待决定 |
| **P4** | 审核机制缺失 | 生产可用性 | 3-5人天 | ☐ 待决定 |
| **P5** | 硬编码值散落 | 可配置性 | 0.5-1人天 | ☐ 待决定 |

---

## P0 级别：缺少单元测试

### 问题描述

当前`application/services/`目录下没有任何测试文件。作为核心业务模块，分镜服务（997行代码）完全没有单元测试覆盖，存在以下风险：

1. **重构风险**：任何代码重构都无法验证正确性
2. **回归风险**：修复bug可能引入新的bug
3. **知识流失**：开发人员离职后，核心逻辑无法被正确理解
4. **持续集成缺失**：无法接入CI/CD自动化测试

### 影响分析

| 影响维度 | 具体表现 | 严重程度 |
|----------|----------|----------|
| 质量风险 | 核心逻辑无验证 | 🔴 高 |
| 维护成本 | 人工测试成本高 | 🟡 中 |
| 上线风险 | Bug逃逸到生产环境 | 🔴 高 |
| 团队效率 | 回归测试耗时长 | 🟡 中 |

### 改进方案

#### 方案一：完整测试覆盖（推荐）

为所有核心逻辑编写单元测试：

```go
// application/services/storyboard_service_test.go

package services

import (
    "testing"
    "github.com/drama-generator/backend/domain/models"
)

// 测试时长计算逻辑
func TestDurationCalculation(t *testing.T) {
    tests := []struct {
        name           string
        script         string
        expectedMin    int
        expectedMax    int
    }{
        {
            name:        "纯对话场景",
            script:      "A：你好。B：再见。",
            expectedMin: 4,
            expectedMax: 6,
        },
        {
            name:        "纯动作场景",
            script:      "他站起来，走过去，打开门。",
            expectedMin: 5,
            expectedMax: 8,
        },
        {
            name:        "混合场景",
            script:      "他站起来，A：你好。B：再见。",
            expectedMin: 6,
            expectedMax: 10,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            duration := calculateDuration(tt.script)
            if duration < tt.expectedMin || duration > tt.expectedMax {
                t.Errorf("expected duration between %d and %d, got %d",
                    tt.expectedMin, tt.expectedMax, duration)
            }
        })
    }
}

// 测试Prompt构建
func TestBuildStoryboardPrompt(t *testing.T) {
    // 测试各种输入组合
}

// 测试JSON解析容错
func TestParseStoryboardJSON(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        expectErr bool
    }{
        {
            name:      "有效JSON数组",
            input:     `[{"shot_number": 1}]`,
            expectErr: false,
        },
        {
            name:      "有效JSON对象",
            input:     `{"storyboards": [{"shot_number": 1}]}`,
            expectErr: false,
        },
        {
            name:      "Markdown代码块包裹",
            input:     "```json\n[{\"shot_number\": 1}]\n```",
            expectErr: false,
        },
    }
}

// 测试边界条件
func TestExtractInitialPose(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {
            input:    "他站起来，然后走过去",
            expected: "他站起来",
        },
        {
            input:    "陈峥弯腰双手握住撬棍用力撬动保险箱门",
            expected: "陈峥弯腰双手握住撬棍用力撬动保险箱门",
        },
    }
}
```

#### 方案二：关键路径测试

仅测试最核心的时长计算和JSON解析逻辑，工作量减半。

### 工作量估算

| 任务 | 工作量 |
|------|--------|
| 时长计算测试 | 0.5人天 |
| Prompt构建测试 | 0.5人天 |
| JSON解析测试 | 0.5人天 |
| 边界条件测试 | 0.5人天 |
| 数据保存测试 | 1人天 |
| **合计** | **2-3人天** |

### 风险评估

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|----------|
| 测试设计不当 | 低 | 中 | 参照项目现有测试风格 |
| 改动现有逻辑 | 低 | 高 | 避免改动生产代码 |
| 覆盖率不足 | 中 | 中 | 逐步补充 |

### 实施建议

**推荐方案**：方案一（完整测试覆盖）

**最佳时机**：
- 新功能开发间隙
- 大规模重构前
- 版本发布前的质量加固期

**后续跟进**：
- 将测试覆盖率纳入CI（目标>80%）
- 每次PR要求包含测试

---

## P1 级别：AI输出容错不足

### 问题描述

当前代码假设AI会返回合法JSON，但实际场景中AI可能：

1. 返回Markdown代码块（```json ... ```）
2. 返回被截断的不完整JSON
3. 返回缺少字段的JSON
4. 返回不符合Schema的数据

### 影响分析

| 影响维度 | 具体表现 | 严重程度 |
|----------|----------|----------|
| 用户体验 | AI抽风时流程中断 | 🔴 高 |
| 排查难度 | 错误信息不清晰 | 🟡 中 |
| 容错能力 | 无法处理异常输出 | 🔴 高 |

### 改进方案

```go
// pkg/utils/ai_json.go

package utils

import (
    "bytes"
    "encoding/json"
    "strings"
)

// SafeParseAIJSON 增强版AI输出解析
func SafeParseAIJSON(raw string, v interface{}) error {
    // 1. 清理Markdown代码块
    cleaned := cleanMarkdownCodeBlock(raw)
    
    // 2. 修复不完整JSON
    cleaned = fixIncompleteJSON(cleaned)
    
    // 3. 宽松解析（允许未知字段）
    decoder := json.NewDecoder(bytes.NewReader([]byte(cleaned)))
    decoder.DisallowUnknownFields()
    return decoder.Decode(v)
}

// cleanMarkdownCodeBlock 清理Markdown代码块包装
func cleanMarkdownCodeBlock(raw string) string {
    // 移除 ```json 和 ```
    cleaned := strings.TrimPrefix(raw, "```json")
    cleaned = strings.TrimPrefix(cleaned, "```")
    cleaned = strings.TrimSuffix(cleaned, "```")
    cleaned = strings.TrimSpace(cleaned)
    return cleaned
}

// fixIncompleteJSON 修复被截断的JSON
func fixIncompleteJSON(raw string) string {
    // 平衡括号数量
    openCount := strings.Count(raw, "{")
    closeCount := strings.Count(raw, "}")
    
    // 如果开括号多于闭括号，需要补全
    for openCount > closeCount {
        // 在最后一个逗号或开括号后截断
        lastComma := strings.LastIndex(raw, ",")
        lastOpen := strings.LastIndex(raw, "{")
        
        if lastComma > lastOpen {
            raw = raw[:lastComma]
        }
        openCount--
    }
    
    return raw
}

// ParseResult 解析结果包装
type ParseResult struct {
    Success    bool
    Data       interface{}
    Raw        string
    WasFixed   bool
    ErrorMsg   string
}

// SafeParseWithResult 带详细结果的解析
func SafeParseWithResult(raw string, v interface{}) ParseResult {
    result := ParseResult{
        Success: false,
        Raw:     raw,
    }
    
    // 检查是否包含Markdown
    if strings.Contains(raw, "```") {
        result.WasFixed = true
        raw = cleanMarkdownCodeBlock(raw)
    }
    
    // 尝试解析
    err := json.Unmarshal([]byte(raw), v)
    if err != nil {
        // 尝试修复
        fixed := fixIncompleteJSON(raw)
        if fixed != raw {
            result.WasFixed = true
            err = json.Unmarshal([]byte(fixed), v)
        }
    }
    
    if err != nil {
        result.ErrorMsg = err.Error()
        return result
    }
    
    result.Success = true
    return result
}
```

### 工作量估算

| 任务 | 工作量 |
|------|--------|
| 实现SafeParseAIJSON增强版 | 1人天 |
| 集成到storyboard_service | 0.5人天 |
| 集成到frame_prompt_service | 0.5人天 |
| 测试用例编写 | 0.5人天 |
| **合计** | **2.5人天** |

### 风险评估

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|----------|
| 修复逻辑有bug | 中 | 中 | 全面的测试用例 |
| 性能下降 | 低 | 低 | 仅在解析失败时执行修复 |
| 兼容性问题 | 低 | 中 | 保留原有解析逻辑作为fallback |

### 实施建议

**推荐方案**：实施方案一

**测试场景**：
- 各种Markdown包装格式
- 各种截断位置
- 各种缺失字段情况
- 性能基准测试

---

## P2 级别：单文件过大

### 问题描述

`storyboard_service.go`文件达到997行，违反单一职责原则。包含：
- 主服务逻辑（300行）
- Prompt构建逻辑（200行）
- DTO定义（150行）
- 辅助函数（150行）
- 数据访问逻辑（200行）

### 影响分析

| 影响维度 | 具体表现 | 严重程度 |
|----------|----------|----------|
| 可维护性 | 难以定位问题 | 🟡 中 |
| 并行开发 | 冲突概率高 | 🟡 中 |
| 代码复用 | 难以复用 | 🟢 低 |
| 学习曲线 | 新人难以理解 | 🟡 中 |

### 改进方案

```
application/services/
├── storyboard_service.go              # 主服务（~300行）
│   ├── 结构体定义
│   ├── GenerateStoryboard方法
│   └── 核心业务流程
│
├── storyboard_prompt_builder.go        # Prompt构建（~200行）
│   ├── GetStoryboardSystemPrompt()
│   ├── BuildStoryboardPrompt()
│   ├── BuildCharacterContext()
│   └── BuildSceneContext()
│
├── storyboard_dto.go                   # DTO定义（~150行）
│   ├── Storyboard结构体
│   ├── GenerateStoryboardResult
│   ├── CreateStoryboardRequest
│   └── FramePrompt相关
│
├── storyboard_repository.go           # 数据访问（~200行）
│   ├── SaveStoryboards()
│   ├── LoadStoryboardsByEpisode()
│   ├── DeleteStoryboard()
│   └── UpdateStoryboard()
│
├── storyboard_extractor.go            # 文本提取（~150行）
│   ├── ExtractInitialPose()
│   ├── ExtractSimpleLocation()
│   └── ExtractSimplePose()
│
└── storyboard_test.go                 # 单元测试
```

**迁移策略**：使用IDE重构功能，确保不改变代码逻辑

```bash
# 步骤1：创建新文件
touch application/services/storyboard_prompt_builder.go

# 步骤2：使用IDE重构移动代码

# 步骤3：验证编译通过
go build ./...

# 步骤4：运行现有测试（如有）
go test ./application/services/...
```

### 工作量估算

| 任务 | 工作量 |
|------|--------|
| 创建新文件结构 | 0.5人天 |
| 迁移Prompt构建逻辑 | 0.5人天 |
| 迁移DTO定义 | 0.25人天 |
| 迁移数据访问逻辑 | 0.5人天 |
| 迁移辅助函数 | 0.25人天 |
| 编译验证 | 0.25人天 |
| **合计** | **2.25人天** |

### 风险评估

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|----------|
| 代码迁移出错 | 中 | 高 | 充分的编译验证 |
| import循环依赖 | 中 | 高 | 遵循依赖方向 |
| 功能回归 | 低 | 高 | 回归测试 |

### 实施建议

**推荐方案**：实施方案一

**实施步骤**：
1. 先建立新文件结构
2. 使用IDE重构移动代码（手动复制粘贴容易出错）
3. 每次移动后编译验证
4. 最后删除原文件中的重复代码

---

## P3 级别：超长剧本处理

### 问题描述

当剧本非常长（如10万字）时：
1. Prompt可能超过模型上下文限制（通常8K-32K tokens）
2. AI输出可能不完整
3. 生成质量可能下降

### 影响分析

| 影响维度 | 具体表现 | 严重程度 |
|----------|----------|----------|
| 功能完整性 | 超长剧本无法处理 | 🔴 高 |
| 用户体验 | 流程中断 | 🔴 高 |
| 生成质量 | 长剧本质量下降 | 🟡 中 |

### 改进方案

```go
// application/services/storyboard_chunk_service.go

package services

import (
    "math"
    "strings"
    "github.com/sashabaranov/go-openai"
)

const (
    // Token估算（中文约1.5token/字，英文约4token/词）
    avgTokensPerChineseChar = 1.5
    avgTokensPerEnglishWord = 1.3
    
    // 保留空间
    reservedTokens = 2000
)

// ChunkConfig 分块配置
type ChunkConfig struct {
    MaxTokens      int     // 模型最大token数
    OverlapRatio   float64 // 重叠比例（0.1 = 10%）
    MinChunkSize   int     // 最小chunk大小（tokens）
}

// ChunkResult 分块结果
type ChunkResult struct {
    Chunks        []string
    TotalTokens   int
    ChunkCount    int
}

// ChunkScript 脚本分块
func ChunkScript(script string, config ChunkConfig) ChunkResult {
    // 1. 估算总tokens
    totalTokens := estimateTokens(script)
    
    // 2. 如果不需要分块
    maxPayload := config.MaxTokens - reservedTokens
    if totalTokens <= maxPayload {
        return ChunkResult{
            Chunks:      []string{script},
            TotalTokens: totalTokens,
            ChunkCount:  1,
        }
    }
    
    // 3. 计算chunk大小和数量
    chunkSize := int(float64(maxPayload) * (1 - config.OverlapRatio))
    chunkCount := int(math.Ceil(float64(totalTokens) / float64(chunkSize)))
    
    // 4. 按段落分割（保持语义完整性）
    paragraphs := splitByParagraphs(script)
    
    // 5. 构建chunks
    chunks := buildChunks(paragraphs, chunkSize, config.OverlapRatio)
    
    return ChunkResult{
        Chunks:      chunks,
        TotalTokens: totalTokens,
        ChunkCount:  len(chunks),
    }
}

// splitByParagraphs 按段落分割
func splitByParagraphs(script string) []string {
    // 按空行分割
    paragraphs := strings.Split(script, "\n\n")
    
    // 过滤空段落
    var result []string
    for _, p := range paragraphs {
        if strings.TrimSpace(p) != "" {
            result = append(result, p)
        }
    }
    return result
}

// buildChunks 构建语义完整的chunks
func buildChunks(paragraphs []string, chunkSize int, overlapRatio float64) []string {
    var chunks []string
    var currentChunk strings.Builder
    currentSize := 0
    
    overlapSize := int(float64(chunkSize) * overlapRatio)
    
    for i, para := range paragraphs {
        paraTokens := estimateTokens(para)
        
        // 如果加上这个段落会超过限制
        if currentSize+paraTokens > chunkSize {
            // 保存当前chunk
            chunks = append(chunks, currentChunk.String())
            
            // 计算重叠部分
            overlap := ""
            if len(chunks) > 0 && overlapSize > 0 {
                prevChunk := chunks[len(chunks)-1]
                overlap = extractOverlap(prevChunk, overlapSize)
            }
            
            // 开始新chunk（包含重叠部分）
            currentChunk.Reset()
            if overlap != "" {
                currentChunk.WriteString(overlap)
                currentChunk.WriteString("\n\n")
            }
            currentChunk.WriteString(para)
            currentSize = estimateTokens(currentChunk.String())
        } else {
            // 添加到当前chunk
            if currentChunk.Len() > 0 {
                currentChunk.WriteString("\n\n")
            }
            currentChunk.WriteString(para)
            currentSize = estimateTokens(currentChunk.String())
        }
    }
    
    // 保存最后一个chunk
    if currentChunk.Len() > 0 {
        chunks = append(chunks, currentChunk.String())
    }
    
    return chunks
}

// estimateTokens 估算token数量
func estimateTokens(text string) int {
    // 简化估算：中文按字符数*1.5，英文按词数*1.3
    chineseCount := countChineseChars(text)
    englishCount := countEnglishWords(text)
    
    return int(float64(chineseCount)*avgTokensPerChineseChar) + 
           int(float64(englishCount)*avgTokensPerEnglishWord)
}

// extractOverlap 提取重叠部分
func extractOverlap(chunk string, maxTokens int) string {
    sentences := strings.Split(chunk, "。")
    var result strings.Builder
    currentTokens := 0
    
    // 从后往前取句子，保持在maxTokens内
    for i := len(sentences) - 1; i >= 0; i-- {
        sentence := sentences[i] + "。"
        sentenceTokens := estimateTokens(sentence)
        
        if currentTokens+sentenceTokens > maxTokens {
            break
        }
        
        result.WriteString(sentence)
        currentTokens += sentenceTokens
    }
    
    return result.String()
}
```

### 工作量估算

| 任务 | 工作量 |
|------|--------|
| Chunk服务核心逻辑 | 3人天 |
| 集成到StoryboardService | 1人天 |
| 超长剧本测试 | 2人天 |
| 性能优化 | 1人天 |
| **合计** | **7人天** |

### 风险评估

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|----------|
| 分块语义断裂 | 中 | 高 | 按段落分块，保持语义 |
| 成本增加 | 高 | 中 | 明确告知用户 |
| 实现复杂度高 | 中 | 中 | 渐进式实现 |

### 实施建议

**推荐方案**：可延后实施

**前置条件**：
- P0（测试覆盖）先完成
- 有明确的用户反馈需求

**简化方案**：
- 先实现8K token限制+错误提示
- 超长剧本提示用户分拆

---

## P4 级别：缺少分镜审核机制

### 问题描述

当前流程是「AI生成 → 直接保存」，对于生产级应用，缺少：
- 人工审核环节
- 分镜版本管理
- 审核状态流转

### 影响分析

| 影响维度 | 具体表现 | 严重程度 |
|----------|----------|----------|
| 内容质量 | AI错误无法拦截 | 🔴 高 |
| 合规风险 | 不当内容流出 | 🔴 高 |
| 用户体验 | 一次性生成，无法预览 | 🟡 中 |

### 改进方案

```go
// domain/models/storyboard_status.go

package models

// StoryboardStatus 分镜状态
type StoryboardStatus string

const (
    StatusPending   StoryboardStatus = "pending"     // 待审核（AI生成后）
    StatusApproved  StoryboardStatus = "approved"    // 已审核
    StatusRejected  StoryboardStatus = "rejected"   // 已驳回
    StatusGenerating StoryboardStatus = "generating" // 生成中
    StatusCompleted StoryboardStatus = "completed"   // 已完成（图片/视频都生成完）
)

// StoryboardAuditLog 审核日志
type StoryboardAuditLog struct {
    ID            uint              `gorm:"primaryKey" json:"id"`
    StoryboardID  uint              `gorm:"index" json:"storyboard_id"`
    Action        string            `gorm:"size:50" json:"action"`  // approve/reject/edit
    FromStatus    StoryboardStatus  `json:"from_status"`
    ToStatus      StoryboardStatus  `json:"to_status"`
    OperatorID    uint              `json:"operator_id"`
    OperatorType  string            `gorm:"size:20" json:"operator_type"` // user/admin
    Comment       string            `gorm:"type:text" json:"comment"`
    CreatedAt     time.Time         `gorm:"autoCreateTime" json:"created_at"`
}

// UpdateStoryboardStatus 更新分镜状态
func (sb *Storyboard) UpdateStoryboardStatus(db *gorm.DB, newStatus StoryboardStatus, operatorID uint, comment string) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. 记录状态变更
        auditLog := StoryboardAuditLog{
            StoryboardID:  sb.ID,
            Action:        "status_change",
            FromStatus:    StoryboardStatus(sb.Status),
            ToStatus:      newStatus,
            OperatorID:    operatorID,
            OperatorType:  "user",
            Comment:       comment,
        }
        if err := tx.Create(&auditLog).Error; err != nil {
            return err
        }
        
        // 2. 更新状态
        return tx.Model(sb).Update("status", newStatus).Error
    })
}

// 前端API设计

// GET /api/v1/storyboards/:id/pending
// 返回待审核的分镜列表

// POST /api/v1/storyboards/:id/approve
// Body: { "comment": "审核通过" }

// POST /api/v1/storyboards/:id/reject
// Body: { "comment": "驳回原因" }
```

### 工作量估算

| 任务 | 工作量 |
|------|--------|
| 数据库迁移（新增状态字段+审核日志表） | 0.5人天 |
| 后端审核逻辑 | 1人天 |
| 前端审核界面 | 2人天 |
| 审核工作流设计 | 0.5人天 |
| 测试 | 1人天 |
| **合计** | **5人天** |

### 风险评估

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|----------|
| 工作流设计不合理 | 中 | 中 | 参考成熟产品设计 |
| 前端改动量大 | 高 | 中 | MVP优先 |
| 用户不接受 | 中 | 中 | 可配置是否开启 |

### 实施建议

**推荐方案**：可延后实施（V2功能）

**MVP范围**：
- 仅增加pending状态
- 后端审核API
- 简单的前端审核按钮
- 暂不实现审核日志

**前置条件**：
- P0（测试覆盖）先完成

---

## P5 级别：硬编码值散落

### 问题描述

多处出现Magic Number，缺乏可配置性：

```go
ai.WithMaxTokens(16000)                          // 为什么是16000？
if totalDuration > 300 { ... }                   // 300代表什么？
if count == 3 { frames[0] = ... }               // 3是magic number
const MaxShotDuration = 12
const MinShotDuration = 4
```

### 改进方案

```go
// config/config.go 或 pkg/constants/constants.go

package constants

// AI 配置常量
const (
    // Token限制
    DefaultMaxTokens       = 16000
    MaxStoryboardTokens   = 16000
    MaxEpisodeTokens      = 32000
    ReservedPromptTokens  = 2000
    
    // 模型配置
    DefaultModel          = "gpt-4"
    FallbackModel         = "gpt-3.5-turbo"
)

// 时长配置常量
const (
    // 单镜头时长限制（秒）
    MaxSingleShotDuration = 12
    MinSingleShotDuration = 4
    DefaultShotDuration   = 5
    
    // 剧集时长限制（秒）
    MaxEpisodeDuration    = 300  // 5分钟
    MinEpisodeDuration    = 60   // 1分钟
    
    // 调整参数
    DialogueShortLimit    = 20   // 短对话字数上限
    DialogueMediumLimit   = 50   // 中等对话字数上限
)

// 分镜配置常量
const (
    DefaultPanelCount     = 3
    MaxPanelCount         = 4
    ActionSequenceCount   = 5
    
    // 详细度要求（字数下限）
    MinTimeDescLength     = 15
    MinLocationDescLength = 20
    MinActionDescLength   = 25
    MinResultDescLength   = 25
    MinAtmosphereLength   = 20
)

// 重构代码示例

// 重构前
func calculateDuration(script string) int {
    if totalDuration > 300 {  // 300是什么？
        return 300
    }
    return totalDuration
}

// 重构后
func calculateDuration(script string) int {
    if totalDuration > constants.MaxEpisodeDuration {
        return constants.MaxEpisodeDuration
    }
    return totalDuration
}
```

### 工作量估算

| 任务 | 工作量 |
|------|--------|
| 创建constants文件 | 0.25人天 |
| 迁移AI相关常量 | 0.25人天 |
| 迁移时长相关常量 | 0.25人天 |
| 迁移其他常量 | 0.25人天 |
| **合计** | **1人天** |

### 风险评估

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|----------|
| 遗漏迁移 | 低 | 低 | 全局搜索Magic Number |
| 配置冲突 | 低 | 低 | 保留原有默认值 |

### 实施建议

**推荐方案**：立即实施

**最佳时机**：任何代码改动间隙

**实施策略**：
1. 新常量文件
2. 逐个迁移
3. 保留原有默认值（向后兼容）

---

## 决策矩阵

### 成本效益分析

| 优先级 | 改进项 | 成本 | 收益 | ROI |
|--------|--------|------|------|-----|
| P0 | 单元测试 | 3人天 | 高（质量保障） | ⭐⭐⭐⭐⭐ |
| P1 | AI容错 | 2.5人天 | 高（稳定性） | ⭐⭐⭐⭐ |
| P2 | 代码拆分 | 2人天 | 中（可维护） | ⭐⭐⭐ |
| P3 | 超长剧本 | 7人天 | 中（功能完整） | ⭐⭐ |
| P4 | 审核机制 | 5人天 | 高（生产可用） | ⭐⭐⭐⭐ |
| P5 | 常量抽取 | 1人天 | 低（编码规范） | ⭐⭐ |

### 推荐实施顺序

```
第一阶段（V1.1，本周完成）
├── P5 常量抽取（1人天）
└── P1 AI容错（2.5人天）

第二阶段（V1.2，下周完成）
├── P0 单元测试（3人天）
└── P2 代码拆分（2人天）

第三阶段（V2.0规划）
├── P4 审核机制（5人天）
└── P3 超长剧本（7人天）
```

### 决策确认表

请根据实际情况勾选：

| 改进项 | 批准实施 | 暂缓 | 不实施 | 备注 |
|--------|----------|------|--------|------|
| P0 单元测试 | ☐ | ☐ | ☐ | |
| P1 AI容错 | ☐ | ☐ | ☐ | |
| P2 代码拆分 | ☐ | ☐ | ☐ | |
| P3 超长剧本 | ☐ | ☐ | ☐ | |
| P4 审核机制 | ☐ | ☐ | ☐ | |
| P5 常量抽取 | ☐ | ☐ | ☐ | |

---

## 附录A：实施检查清单

### 实施前

- [ ] 代码冻结（避免冲突）
- [ ] 备份当前代码
- [ ] 确认测试环境
- [ ] 通知相关开发人员

### 实施中

- [ ] 每次提交包含测试
- [ ] 保持编译通过
- [ ] 记录变更日志
- [ ] 及时合并到主分支

### 实施后

- [ ] 全量回归测试
- [ ] 更新文档
- [ ] 代码审查
- [ ] 监控上线后指标

---

## 附录B：相关文档链接

- [分镜模块深度分析报告](../analysis/storyboard-module-analysis.md)
- [代码结构说明](../developer-guide/code-structure.md)
- [API参考文档](../api-reference.md)

---

**文档版本**：v1.0  
**编写日期**：2026年2月  
**作者**：Sisyphus AI Assistant
