# 系统架构设计文档

> **学习目标**：完成本手册学习后，你将能够深入理解火宝短剧平台的技术架构设计，掌握领域驱动设计（DDD）在项目中的实践应用，理解关键架构决策的背景和权衡，具备参与平台核心开发的能力。
> 
> **前置知识**：建议具备Go语言开发经验，熟悉RESTful API设计模式，了解SQL数据库基本操作。建议先完成《API接口参考手册》的学习，理解平台的功能边界和数据模型。
> 
> **难度等级**：⭐⭐⭐⭐（专家设计）
> 
> **预计学习时间**：8-16小时

---

## 第一章：架构设计哲学

### 1.1 设计原则与核心理念

火宝短剧平台的架构设计遵循**领域驱动设计（Domain-Driven Design，简称DDD）**方法论，这是一种以业务领域为核心的软件开发方法。DDD的核心思想是将软件的复杂性映射到业务领域的复杂性上，通过建立**通用语言（Ubiquitous Language）**来消除技术和业务之间的鸿沟。理解这一设计理念对于理解整个系统的架构至关重要。

**为什么选择DDD**：在火宝短剧这样的AI创作平台中，核心业务逻辑涉及剧本创作、角色设计、分镜制作、视频生成等多个领域，每个领域都有其独特的概念、规则和变化模式。传统的三层架构（表现层-业务逻辑层-数据访问层）往往难以清晰地表达这些领域概念，导致业务逻辑分散在各处，难以维护和演进。DDD通过将系统划分为**有界上下文（Bounded Context）**，每个上下文内部拥有自己的领域模型，使得复杂业务能够被清晰地组织和管理。

**架构设计的三大原则**指导着整个系统的构建：

**原则一：领域优先（Domain-First）**。所有设计决策都应该服务于业务需求，而非技术便利。当技术方案与业务需求发生冲突时，优先考虑业务需求。这体现在代码层面，就是领域层（Domain Layer）不依赖任何外部框架和基础设施，所有的领域逻辑都是纯代码实现，可以在任何环境中运行。

**原则二：清晰边界（Clear Boundaries）**。系统的每个部分都有明确的职责边界，层与层之间通过明确定义的接口交互。这种清晰边界带来的好处是**可测试性**——每个组件都可以独立测试；**可替换性**——可以替换某个组件而不影响其他部分；**可理解性**——新成员可以通过边界快速理解系统结构。

**原则三：演进式设计（Evolutionary Design）**。架构不是一次性设计完成的，而是随着对业务理解的深入而不断演进的。系统设计应该预留扩展点，支持未来的变化。这意味着在设计时要识别出**变化点**，将这些变化点设计为可配置的机制，而不是硬编码在代码中。

### 1.2 架构分层设计

火宝短剧采用**四层架构**设计，从上到下依次为：**API层（API Layer）**、**应用服务层（Application Service Layer）**、**领域层（Domain Layer）**和**基础设施层（Infrastructure Layer）**。这种分层遵循经典的DDD分层模式，并针对项目特点进行了适当调整。

```
┌─────────────────────────────────────────────────────────────┐
│                      API 层                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐   │
│  │ HTTP Handler │ │  Request   │ │    Response        │   │
│  │  (Gin)      │ │ Validator  │ │    Formatter       │   │
│  └─────────────┘ └─────────────┘ └─────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    应用服务层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐   │
│  │ Application │ │    Task    │ │    Event            │   │
│  │   Service   │ │  Scheduler │ │    Handler          │   │
│  └─────────────┘ └─────────────┘ └─────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      领域层                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐   │
│  │  Aggregate  │ │   Value    │ │      Event         │   │
│  │   Root     │ │   Object   │ │    (Domain)        │   │
│  └─────────────┘ └─────────────┘ └─────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    基础设施层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐   │
│  │ Repository  │ │ External   │ │    Utility         │   │
│  │ Implementation│ │ Service   │ │                    │   │
│  └─────────────┘ │ Adapter    │ └─────────────────────┘   │
│                  └─────────────┘                           │
└─────────────────────────────────────────────────────────────┘
```

**各层职责定义**：

**API层**负责处理所有外部请求，包括HTTP请求解析、参数验证、认证授权、响应格式化等。API层是系统的「门面」，所有外部调用都通过这一层进入系统。API层不包含任何业务逻辑，它的作用是「接收请求、调用应用服务、返回响应」。

**应用服务层**是系统的「编排层」，负责协调领域对象完成业务用例。应用服务不包含业务规则，它只知道「做什么」（What），不知道「怎么做」（How）。例如，一个「生成视频」的应用服务会调用领域服务创建视频生成任务，协调各个聚合根更新状态，发布领域事件通知相关方。

**领域层**是系统的「核心」，包含所有的业务实体、值对象、领域服务、领域事件等。这一层包含了系统最本质的业务逻辑——不管用什么技术实现，这些逻辑都是不变的。领域层完全独立于任何框架和基础设施，可以独立测试和验证。

**基础设施层**提供技术能力的实现，包括数据库访问、外部服务调用、文件存储、日志记录等。基础设施层的代码实现应用服务和仓库接口，但不包含业务决策。这一层是「可替换的」——可以更换数据库、替换AI服务提供商，而不影响上层代码。

### 1.3 目录结构设计

```
huobao-drama/
├── api/                    # API层实现
│   ├── handlers/          # HTTP处理器
│   ├── middleware/        # 中间件（认证、日志、限流）
│   ├── request/          # 请求参数定义和验证
│   └── response/         # 响应格式化
├── application/           # 应用服务层
│   ├── services/        # 应用服务实现
│   ├── tasks/           # 任务调度
│   └── events/          # 事件处理器
├── domain/               # 领域层（核心）
│   ├── aggregates/       # 聚合根定义
│   ├── entities/         # 实体定义
│   ├── valueobjects/    # 值对象定义
│   ├── services/        # 领域服务
│   ├── events/          # 领域事件
│   └── repositories/    # 仓库接口
├── infrastructure/       # 基础设施层
│   ├── persistence/     # 数据库访问实现
│   ├── external/        # 外部服务适配器
│   └── storage/         # 文件存储实现
├── pkg/                  # 公共包
│   ├── config/          # 配置管理
│   ├── logger/          # 日志封装
│   └── utils/           # 工具函数
├── configs/             # 配置文件
├── web/                  # 前端应用
├── main.go              # 程序入口
└── go.mod               # Go模块定义
```

**目录设计的设计考量**：

这种目录结构的设计遵循了「按职责分」和「按层次分」两个维度。顶层按照职责划分为API、应用服务、领域、基础设施四个包；每个包内部再按照具体的类型进行细分。这种设计的优势是**高内聚、低耦合**——相关代码放在一起，不相关代码自然分离；修改一处代码影响范围有限。

**值得注意的设计点**：

1. **仓库接口定义在领域层**：这是DDD的一个重要实践。领域层定义了仓库应该提供什么能力（接口），基础设施层负责实现这些能力（实现）。这样领域层完全不需要知道数据是如何持久化的。
2. **领域事件是领域层的一部分**：领域事件是领域模型的一部分，它们的定义和触发都在领域层。事件的发布和订阅是基础设施层的职责。
3. **不直接依赖外部框架**：领域层和应用服务层都不直接依赖Gin、SQLite等外部框架，而是通过依赖注入的方式接入。这种设计使得核心业务逻辑可以在不同环境中运行。

---

## 第二章：领域模型设计

### 2.1 聚合根设计

**聚合根（Aggregate Root）**是DDD中的核心概念，它是聚合的入口点，负责维护聚合的不变性。聚合是一组相关对象的集合，这些对象被视为一个数据变更的单元。理解聚合根的设计对于维护数据一致性和系统性能至关重要。

**聚合根设计原则**：

| 原则 | 说明 | 示例 |
|------|------|------|
| 边界清晰 | 每个聚合有明确的边界，内部对象外部不可直接访问 | Project聚合包含Script、Character等 |
| 引用唯一 | 聚合之间通过ID引用，而非直接对象引用 | Storyboard引用CharacterID而非Character对象 |
| 一致性边界 | 聚合内部的数据变更必须保持一致 | Project的更新必须通过Project聚合根 |

**火宝短剧的核心聚合**：

```
聚合边界图：

┌─────────────────────────────────────────────────────────┐
│                    Project Aggregate                     │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │  Project        │──▶│    Script       │              │
│  │  (聚合根)        │  │                 │              │
│  └────────┬────────┘  └─────────────────┘              │
│           │                                                │
│           │  ┌─────────────────┐  ┌─────────────────┐   │
│           ├──▶│   Character    │──▶│ CharacterImage │   │
│           │  └─────────────────┘  └─────────────────┘   │
│           │                                                │
│           │  ┌─────────────────┐  ┌─────────────────┐   │
│           ├──▶│  Storyboard    │──▶│ StoryboardImage│   │
│           │  └─────────────────┘  └─────────────────┘   │
│           │                                                │
│           │  ┌─────────────────┐                          │
│           └──▶│     Video      │                          │
│              └─────────────────┘                          │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    Asset Aggregate                       │
│  ┌─────────────────┐                                   │
│  │     Asset       │                                   │
│  │  (聚合根)        │                                   │
│  └─────────────────┘                                   │
└─────────────────────────────────────────────────────────┘
```

**Project聚合根的详细设计**：

Project是火宝短剧中最顶层的聚合，所有其他聚合（Script、Character、Storyboard、Video）都归属于某个Project。这种设计的考量是：

1. **所有权清晰**：每个资源都属于特定的项目，不存在跨项目的资源引用
2. **权限控制简单**：通过项目的访问控制即可控制所有子资源的访问
3. **数据删除简单**：删除项目时级联删除所有子资源

```go
// Project 聚合根定义
type Project struct {
    ID          string        `json:"id" gorm:"primaryKey;type:varchar(36)"`
    Name        string        `json:"name" gorm:"type:varchar(100);not null"`
    Description string        `json:"description" gorm:"type:text"`
    OwnerID     string        `json:"owner_id" gorm:"type:varchar(36);not null;index"`
    Settings    Settings      `json:"settings" gorm:"type:text"` // JSON存储
    
    // 聚合内部状态
    scripts     []*Script
    characters  []*Character
    storyboards []*Storyboard
    videos     []*Video
    
    // 元数据
    CreatedAt   time.Time     `json:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
}

// Settings 项目配置
type Settings struct {
    Resolution   string `json:"resolution"`
    FrameRate   int    `json:"frame_rate"`
    AspectRatio string `json:"aspect_ratio"`
}
```

**聚合内部的不变性约束**：

聚合根负责维护聚合内部的不变性约束，这些约束确保了业务规则的正确执行：

```go
// AddScript 添加剧本
func (p *Project) AddScript(script *Script) error {
    // 不变性检查1：剧本名称在项目内唯一
    for _, s := range p.scripts {
        if s.Name == script.Name {
            return ErrScriptNameDuplicate
        }
    }
    
    // 不变性检查2：项目状态允许添加剧本
    if p.Status == ProjectStatusClosed {
        return ErrProjectClosed
    }
    
    // 添加剧本
    script.ProjectID = p.ID
    p.scripts = append(p.scripts, script)
    
    // 发布领域事件
    p.RegisterEvent(ScriptCreatedEvent{
        ProjectID: p.ID,
        ScriptID:  script.ID,
    })
    
    return nil
}
```

### 2.2 实体与值对象

**实体（Entity）**是具有唯一标识的对象，其相等性通过标识判断，而非属性值。即使两个实体的所有属性都相同，如果ID不同，它们也是不同的对象。

**值对象（Value Object）**是没有标识的对象，其相等性通过属性值判断。值对象是不可变的，创建后不能修改，如果需要修改则创建新的值对象替代。

**实体示例：Script（剧本）**：

```go
// Script 实体
type Script struct {
    ID          string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    ProjectID   string    `json:"project_id" gorm:"type:varchar(36);not null;index"`
    Title       string    `json:"title" gorm:"type:varchar(200);not null"`
    Content     string    `json:"content" gorm:"type:text"`
    Status      string    `json:"status" gorm:"type:varchar(20);not null;default:draft"`
    Metadata    Metadata  `json:"metadata" gorm:"type:text"`
    
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Script之间的相等性通过ID判断
func (s *Script) Equals(other *Script) bool {
    return s.ID == other.ID
}
```

**值对象示例：ScriptMetadata（剧本元数据）**：

```go
// ScriptMetadata 值对象
type ScriptMetadata struct {
    WordCount       int       `json:"word_count"`
    SceneCount      int       `json:"scene_count"`
    CharacterCount  int       `json:"character_count"`
    EstimatedDuration float64 `json:"estimated_duration"` // 预计时长（秒）
    Style           string    `json:"style"`
    Language        string    `json:"language"`
}

// 值对象是不可变的
func NewScriptMetadata(wordCount, sceneCount int) *ScriptMetadata {
    return &ScriptMetadata{
        WordCount:   wordCount,
        SceneCount:  sceneCount,
        Language:    "zh-CN",
    }
}

// 修改值对象需要创建新的实例
func (m *ScriptMetadata) WithUpdatedDuration(duration float64) *ScriptMetadata {
    return &ScriptMetadata{
        WordCount:        m.WordCount,
        SceneCount:       m.SceneCount,
        EstimatedDuration: duration,
        Style:            m.Style,
        Language:         m.Language,
    }
}

// 值对象的相等性通过属性值判断
func (m *ScriptMetadata) Equals(other *ScriptMetadata) bool {
    return m.WordCount == other.WordCount &&
           m.SceneCount == other.SceneCount &&
           m.EstimatedDuration == other.EstimatedDuration &&
           m.Style == other.Style &&
           m.Language == other.Language
}
```

**值对象的设计考量**：

使用值对象代替简单类型可以带来以下好处：

1. **类型安全**：将`int`和`float64`封装为有意义类型，避免传入错误的值
2. **业务规则封装**：值对象的创建逻辑包含验证规则
3. **无标识比较**：值对象的比较基于属性值，直观清晰
4. **性能优化**：值对象可以自由复制，不需要复杂的引用管理

### 2.3 领域服务

**领域服务（Domain Service）**用于表达那些不属于任何实体或值对象的业务逻辑。当某个操作需要协调多个聚合根，或者需要与外部系统交互时，应该使用领域服务。

**领域服务的设计原则**：

1. **无状态**：领域服务应该是无状态的，不保存任何状态信息
2. **幂等性**：相同的输入应该产生相同的输出
3. **纯粹性**：只包含领域逻辑，不涉及技术细节

**示例：ScriptGenerationService（剧本生成服务）**：

```go
// IScriptGenerationService 剧本生成领域服务接口
type IScriptGenerationService interface {
    GenerateFromPrompt(ctx context.Context, prompt string, options GenerationOptions) (*Script, error)
    EnhanceScript(ctx context.Context, script *Script, feedback string) (*Script, error)
}

// GenerationOptions 生成选项
type GenerationOptions struct {
    Style       string
    Length      LengthOption
    Language    string
    IncludeStageDirections bool
}

// LengthOption 生成长度选项
type LengthOption string

const (
    LengthShort   LengthOption = "short"   // 短（约5场景）
    LengthMedium  LengthOption = "medium"  // 中（约10场景）
    LengthLong    LengthOption = "long"    // 长（约20场景）
)
```

### 2.4 领域事件

**领域事件（Domain Event）**是领域中发生的事情的记录，用于解耦系统各部分之间的依赖。当某个业务操作发生时，可以发布领域事件，其他模块可以订阅这些事件并做出响应。

**领域事件的设计考量**：

| 设计点 | 考量 |
|--------|------|
| 事件命名 | 使用过去式命名，如`ScriptCreatedEvent` |
| 事件数据 | 包含足够的信息供订阅者处理，必要时包含完整快照 |
| 事件顺序 | 确保同一聚合根的事件按时间顺序处理 |
| 幂等性 | 事件处理器应该能够处理重复事件 |

**示例：领域事件定义与发布**：

```go
// ScriptCreatedEvent 剧本创建事件
type ScriptCreatedEvent struct {
    EventBase
    ProjectID   string `json:"project_id"`
    ScriptID    string `json:"script_id"`
    ScriptTitle string `json:"script_title"`
    GeneratedBy string `json:"generated_by"` // AI或用户
}

// StoryboardGeneratedEvent 分镜生成事件
type StoryboardGeneratedEvent struct {
    EventBase
    ProjectID    string `json:"project_id"`
    ScriptID     string `json:"script_id"`
    StoryboardID string `json:"storyboard_id"`
    SceneNumber  int    `json:"scene_number"`
}

// 发布领域事件
func (p *Project) RegisterEvent(event domain.Event) {
    p.events = append(p.events, event)
}

// 发布所有已注册的事件
func (p *Project) PublishEvents(ctx context.Context, eventBus IEventBus) error {
    for _, event := range p.events {
        if err := eventBus.Publish(ctx, event); err != nil {
            return err
        }
    }
    p.events = nil // 清空已发布的事件
    return nil
}
```

---

## 第三章：应用服务层设计

### 3.1 应用服务与领域服务的区别

应用服务（Application Service）和领域服务（Domain Service）虽然名称相似，但职责截然不同。理解它们的区别是掌握DDD的关键。

| 对比维度 | 应用服务 | 领域服务 |
|----------|----------|----------|
| 职责 | 编排业务流程 | 执行业务规则 |
| 所在层 | 应用服务层 | 领域层 |
| 依赖 | 依赖领域对象和领域服务 | 只依赖领域概念 |
| 事务 | 管理事务边界 | 不涉及事务 |
| 返回值 | DTO或领域对象 | 领域对象或值对象 |

**简单类比**：如果把系统比作一家餐厅，领域服务就是「厨师」，负责「做饭」这个核心业务；应用服务就是「服务员」，负责「接待顾客、点单、上菜」这些流程编排工作。

**应用服务示例：ProjectAppService（项目应用服务）**：

```go
// ProjectAppService 项目应用服务
type ProjectAppService struct {
    projectRepo   domain.IProjectRepository
    scriptService domain.IScriptGenerationService
    eventBus      domain.IEventBus
    logger        *zap.Logger
}

// CreateProject 创建项目
func (s *ProjectAppService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectResponse, error) {
    // 1. 创建领域对象
    project, err := domain.NewProject(req.Name, req.Description, req.OwnerID)
    if err != nil {
        return nil, err
    }
    
    // 2. 设置配置
    project.UpdateSettings(req.Settings)
    
    // 3. 持久化
    if err := s.projectRepo.Save(ctx, project); err != nil {
        return nil, err
    }
    
    // 4. 发布事件
    project.RegisterEvent(domain.ProjectCreatedEvent{
        ProjectID: project.ID,
        OwnerID:   req.OwnerID,
    })
    if err := project.PublishEvents(ctx, s.eventBus); err != nil {
        s.logger.Error("failed to publish events", zap.Error(err))
    }
    
    return ToProjectResponse(project), nil
}

// GenerateScript 生成剧本
func (s *ProjectAppService) GenerateScript(ctx context.Context, req *GenerateScriptRequest) (*TaskResponse, error) {
    // 1. 获取项目
    project, err := s.projectRepo.FindByID(ctx, req.ProjectID)
    if err != nil {
        return nil, err
    }
    
    // 2. 调用领域服务生成剧本
    script, err := s.scriptService.GenerateFromPrompt(
        ctx,
        req.Prompt,
        domain.GenerationOptions{
            Style:  domain.LengthOption(req.Length),
            Style:  req.Style,
            Length: req.Length,
        },
    )
    if err != nil {
        return nil, err
    }
    
    // 3. 将剧本添加到项目
    if err := project.AddScript(script); err != nil {
        return nil, err
    }
    
    // 4. 持久化
    if err := s.projectRepo.Save(ctx, project); err != nil {
        return nil, err
    }
    
    // 5. 发布事件
    if err := project.PublishEvents(ctx, s.eventBus); err != nil {
        s.logger.Error("failed to publish events", zap.Error(err))
    }
    
    return ToTaskResponse(script.GetGenerationTask()), nil
}
```

### 3.2 任务调度设计

视频生成等AI任务通常是耗时操作，不能在HTTP请求中同步处理。系统采用**异步任务**模式，将耗时任务放入队列，由后台工作进程处理。

**任务状态流转**：

```
任务状态机：

                    ┌──────────────┐
                    │   PENDING    │  任务创建，等待处理
                    └──────┬───────┘
                           │
                           │ 任务调度器拾取任务
                           ▼
                    ┌──────────────┐
              ┌──────│  PROCESSING │──────┐
              │      └──────┬───────┘      │
              │             │              │
              │             ▼              │
              │      ┌──────────────┐       │
              │      │   SUCCESS   │◀──────┘  任务成功完成
              │      └──────┬───────┘
              │             │
              │             │ 任务失败
              │             ▼
              │      ┌──────────────┐
              └─────▶│   FAILED    │  任务执行失败
                     └──────────────┘
                              │
                              │ 重试策略判断
                              ▼
                     ┌──────────────┐
                     │   RETRYING   │  等待重试
                     └──────┬───────┘
                              │
                              │ 重试次数耗尽
                              ▼
                     ┌──────────────┐
                     │   DEAD       │  进入死信队列
                     └──────────────┘
```

**任务数据结构**：

```go
// Task 任务聚合根
type Task struct {
    ID          string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
    Type        TaskType    `json:"type" gorm:"type:varchar(30);not null;index"`
    Status      TaskStatus  `json:"status" gorm:"type:varchar(20);not null;index"`
    Payload     string      `json:"payload" gorm:"type:text"` // JSON存储
    Result      string      `json:"result" gorm:"type:text"`  // JSON存储
    Progress    int         `json:"progress" gorm:"default:0"`
    RetryCount  int         `json:"retry_count" gorm:"default:0"`
    MaxRetries  int         `json:"max_retries" gorm:"default:3"`
    Error       string      `json:"error" gorm:"type:text"`
    StartedAt   *time.Time  `json:"started_at"`
    CompletedAt *time.Time  `json:"completed_at"`
    
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

// TaskType 任务类型
type TaskType string

const (
    TaskTypeScriptGeneration    TaskType = "script_generation"
    TaskTypeCharacterImageGen   TaskType = "character_image_generation"
    TaskTypeStoryboardGen       TaskType = "storyboard_generation"
    TaskTypeVideoGeneration     TaskType = "video_generation"
)

// TaskStatus 任务状态
type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "pending"
    TaskStatusProcessing TaskStatus = "processing"
    TaskStatusSuccess     TaskStatus = "success"
    TaskStatusFailed      TaskStatus = "failed"
    TaskStatusRetrying    TaskStatus = "retrying"
    TaskStatusDead        TaskStatus = "dead"
)
```

**任务调度器设计**：

```go
// TaskScheduler 任务调度器
type TaskScheduler struct {
    taskRepo    repository.ITaskRepository
    workerPool  *WorkerPool
    eventBus    domain.IEventBus
    logger      *zap.Logger
}

// WorkerPool 工作线程池
type WorkerPool struct {
    workers     int
    taskQueue   chan *Task
    wg          sync.WaitGroup
    stopChan    chan struct{}
}

// Start 启动调度器
func (s *TaskScheduler) Start(ctx context.Context) error {
    // 启动多个工作协程
    for i := 0; i < s.workerPool.workers; i++ {
        s.workerPool.wg.Add(1)
        go s.worker(ctx, i)
    }
    
    // 定时扫描待处理任务
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            s.workerPool.stopChan <- struct{}{}
            s.workerPool.wg.Wait()
            return ctx.Err()
        case <-ticker.C:
            if err := s.dispatchTasks(ctx); err != nil {
                s.logger.Error("failed to dispatch tasks", zap.Error(err))
            }
        }
    }
}

// dispatchTasks 派发待处理任务
func (s *TaskScheduler) dispatchTasks(ctx context.Context) error {
    tasks, err := s.taskRepo.FindPendingTasks(ctx, 100) // 每次最多派发100个
    if err != nil {
        return err
    }
    
    for _, task := range tasks {
        select {
        case s.workerPool.taskQueue <- task:
            s.logger.Info("task dispatched", zap.String("task_id", task.ID))
        default:
            // 队列满了，停止派发
            return nil
        }
    }
    return nil
}
```

### 3.3 依赖注入设计

系统采用**依赖注入（Dependency Injection）**模式管理组件之间的依赖关系。这种设计使得组件之间的耦合度降低，便于测试和替换实现。

**依赖注入容器**：

```go
// Container 依赖注入容器
type Container struct {
    services map[string]interface{}
}

// Register 注册服务
func (c *Container) Register(name string, service interface{}) {
    c.services[name] = service
}

// Resolve 获取服务
func (c *Container) Resolve(name string) interface{} {
    return c.services[name]
}

// NewContainer 创建新容器
func NewContainer() *Container {
    return &Container{
        services: make(map[string]interface{}),
    }
}
```

**服务注册示例**：

```go
func InitializeContainer(cfg *Config) *Container {
    container := NewContainer()
    
    // 配置
    container.Register("config", cfg)
    
    // 日志
    logger := NewLogger(cfg.Log)
    container.Register("logger", logger)
    
    // 数据库
    db := NewDatabase(cfg.Database)
    container.Register("db", db)
    
    // 领域仓储实现
    projectRepo := persistence.NewProjectRepository(db)
    scriptRepo := persistence.NewScriptRepository(db)
    taskRepo := persistence.NewTaskRepository(db)
    container.Register("project_repo", projectRepo)
    container.Register("script_repo", scriptRepo)
    container.Register("task_repo", taskRepo)
    
    // 事件总线
    eventBus := NewInMemoryEventBus()
    container.Register("event_bus", eventBus)
    
    // 领域服务
    scriptService := application.NewScriptGenerationService(
        projectRepo,
        GetAIService(cfg),
        eventBus,
        logger,
    )
    container.Register("script_service", scriptService)
    
    // 应用服务
    projectAppService := application.NewProjectAppService(
        projectRepo,
        scriptService,
        eventBus,
        logger,
    )
    container.Register("project_app_service", projectAppService)
    
    return container
}
```

---

## 第四章：基础设施层设计

### 4.1 仓储模式实现

**仓储（Repository）**是DDD中用于访问聚合根的抽象接口。仓储隐藏了数据访问的细节，领域层只需要调用仓储的方法，而不需要知道数据是如何存储的。

**仓储接口定义（在领域层）**：

```go
// IProjectRepository 项目仓储接口
type IProjectRepository interface {
    // 查询操作
    FindByID(ctx context.Context, id string) (*Project, error)
    FindByOwner(ctx context.Context, ownerID string, pagination Pagination) ([]*Project, error)
    FindByName(ctx context.Context, name string) ([]*Project, error)
    
    // 持久化操作
    Save(ctx context.Context, project *Project) error
    Update(ctx context.Context, project *Project) error
    Delete(ctx context.Context, id string) error
    
    // 聚合查询
    FindWithScripts(ctx context.Context, projectID string) (*Project, error)
    FindWithCharacters(ctx context.Context, projectID string) (*Project, error)
}
```

**仓储实现（在基础设施层）**：

```go
// GormProjectRepository 基于GORM的项目仓储实现
type GormProjectRepository struct {
    db *gorm.DB
}

// NewProjectRepository 创建仓储实例
func NewProjectRepository(db *gorm.DB) *GormProjectRepository {
    return &GormProjectRepository{db: db}
}

// FindByID 实现
func (r *GormProjectRepository) FindByID(ctx context.Context, id string) (*domain.Project, error) {
    var projectEntity ProjectEntity
    
    // 查询项目
    err := r.db.WithContext(ctx).First(&projectEntity, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrProjectNotFound
        }
        return nil, err
    }
    
    // 转换为领域对象
    return r.toDomain(&projectEntity), nil
}

// Save 实现
func (r *GormProjectRepository) Save(ctx context.Context, project *domain.Project) error {
    entity := r.toEntity(project)
    
    // 判断是新增还是更新
    exists, err := r.exists(ctx, project.ID)
    if err != nil {
        return err
    }
    
    if exists {
        return r.db.WithContext(ctx).Save(entity).Error
    }
    return r.db.WithContext(ctx).Create(entity).Error
}

// 实体与领域对象转换
func (r *GormProjectRepository) toDomain(entity *ProjectEntity) *domain.Project {
    // ... 转换逻辑
}

func (r *GormProjectRepository) toEntity(project *domain.Project) *ProjectEntity {
    // ... 转换逻辑
}
```

**数据模型定义**：

```go
// ProjectEntity 项目数据库实体
type ProjectEntity struct {
    ID          string    `gorm:"primaryKey;type:varchar(36)"`
    Name        string    `gorm:"type:varchar(100);not null"`
    Description string    `gorm:"type:text"`
    OwnerID     string    `gorm:"type:varchar(36);not null;index"`
    Settings    string    `gorm:"type:text"`
    Status      string    `gorm:"type:varchar(20);not null;default:active"`
    
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
    
    // 关联
    Scripts     []ScriptEntity  `gorm:"foreignKey:ProjectID"`
    Characters  []CharacterEntity `gorm:"foreignKey:ProjectID"`
}

// TableName 指定表名
func (ProjectEntity) TableName() string {
    return "projects"
}
```

### 4.2 外部服务适配器

系统需要与多个外部AI服务交互，包括文本生成、图像生成和视频生成服务。通过**适配器模式**将这些外部服务封装为统一的接口，使得上层代码不直接依赖特定的服务实现。

**外部服务接口定义**：

```go
// ITextGenerationService 文本生成服务接口
type ITextGenerationService interface {
    // Generate 生成文本
    Generate(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error)
    
    // StreamGenerate 流式生成
    StreamGenerate(ctx context.Context, req *TextGenerationRequest) (<-chan TextGenerationResponse, error)
}

// TextGenerationRequest 文本生成请求
type TextGenerationRequest struct {
    Model     string            `json:"model"`
    Prompt    string            `json:"prompt"`
    MaxTokens int               `json:"max_tokens"`
    Temperature float64         `json:"temperature"`
    Messages  []ChatMessage     `json:"messages"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
    Role    string `json:"role"`    // system, user, assistant
    Content string `json:"content"`
}

// TextGenerationResponse 文本生成响应
type TextGenerationResponse struct {
    Content    string    `json:"content"`
    Usage      TokenUsage `json:"usage"`
    Model      string    `json:"model"`
    FinishReason string  `json:"finish_reason"`
}
```

**OpenAI适配器实现**：

```go
// OpenAITextGenerationService OpenAI文本生成适配器
type OpenAITextGenerationService struct {
    client  *openai.Client
    model   string
    logger  *zap.Logger
}

// NewOpenAITextGenerationService 创建OpenAI适配器
func NewOpenAITextGenerationService(apiKey, model string, logger *zap.Logger) *OpenAITextGenerationService {
    client := openai.NewClient(apiKey)
    return &OpenAITextGenerationService{
        client: client,
        model:  model,
        logger: logger,
    }
}

// Generate 实现
func (s *OpenAITextGenerationService) Generate(ctx context.Context, req *ITextGenerationRequest) (*ITextGenerationResponse, error) {
    s.logger.Debug("calling OpenAI API",
        zap.String("model", s.model),
        zap.Int("max_tokens", req.MaxTokens),
    )
    
    resp, err := s.client.ChatCompletions(ctx, openai.ChatCompletionRequest{
        Model:       s.model,
        Messages:    s.toOpenAIMessages(req.Messages),
        MaxTokens:   req.MaxTokens,
        Temperature: req.Temperature,
    })
    if err != nil {
        s.logger.Error("OpenAI API error", zap.Error(err))
        return nil, ErrAIServiceUnavailable
    }
    
    return &ITextGenerationResponse{
        Content:      resp.Choices[0].Message.Content,
        Usage:        TokenUsage(resp.Usage),
        Model:        resp.Model,
        FinishReason: resp.Choices[0].FinishReason,
    }, nil
}

// StreamGenerate 实现
func (s *OpenAITextGenerationService) StreamGenerate(ctx context.Context, req *ITextGenerationRequest) (<-chan ITextGenerationResponse, error) {
    stream, err := s.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
        Model:       s.model,
        Messages:    s.toOpenAIMessages(req.Messages),
        MaxTokens:   req.MaxTokens,
        Temperature: req.Temperature,
        Stream:      true,
    })
    if err != nil {
        return nil, err
    }
    
    ch := make(chan ITextGenerationResponse)
    
    go func() {
        defer close(ch)
        defer stream.Close()
        
        reader := bufio.NewReader(stream)
        for {
            response, err := reader.ReadJSON()
            if err != nil {
                if err != io.EOF {
                    s.logger.Error("stream read error", zap.Error(err))
                }
                break
            }
            
            ch <- &ITextGenerationResponse{
                Content: response.Choices[0].Delta.Content,
                Model:   response.Model,
            }
        }
    }()
    
    return ch, nil
}
```

**服务选择器**：

```go
// AIServiceProvider AI服务提供商
type AIServiceProvider struct {
    textServices    map[string]ITextGenerationService
    imageServices   map[string]IImageGenerationService
    videoServices   map[string]IVideoGenerationService
    defaultProvider string
}

// GetTextService 获取文本生成服务
func (p *AIServiceProvider) GetTextService(provider string) (ITextGenerationService, error) {
    if provider == "" {
        provider = p.defaultProvider
    }
    
    service, ok := p.textServices[provider]
    if !ok {
        return nil, ErrUnsupportedProvider
    }
    return service, nil
}

// ProviderRegistry 服务提供商注册表
func NewAIServiceProvider(cfg *Config) *AIServiceProvider {
    p := &AIServiceProvider{
        textServices:  make(map[string]ITextGenerationService),
        imageServices: make(map[string]IImageGenerationService),
        videoServices: make(map[string]IVideoGenerationService),
    }
    
    // 注册各服务提供商
    if cfg.OpenAI.APIKey != "" {
        p.textServices["openai"] = NewOpenAITextGenerationService(
            cfg.OpenAI.APIKey,
            cfg.OpenAI.TextModel,
            logger,
        )
    }
    
    if cfg.DeepSeek.APIKey != "" {
        p.textServices["deepseek"] = NewDeepSeekTextService(
            cfg.DeepSeek.APIKey,
            cfg.DeepSeek.TextModel,
            logger,
        )
    }
    
    // 设置默认提供商
    p.defaultProvider = cfg.AI.DefaultProvider
    
    return p
}
```

### 4.3 文件存储设计

平台生成的所有文件（图片、视频等）都存储在文件系统中。系统设计了统一的**存储接口**，支持本地存储和云存储的切换。

**存储接口定义**：

```go
// IStorage 存储接口
type IStorage interface {
    // 上传文件
    Put(ctx context.Context, key string, reader io.Reader, contentType string) error
    
    // 下载文件
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    
    // 删除文件
    Delete(ctx context.Context, key string) error
    
    // 获取文件URL
    GetURL(ctx context.Context, key string) (string, error)
    
    // 检查文件是否存在
    Exists(ctx context.Context, key string) (bool, error)
}

// StorageType 存储类型
type StorageType string

const (
    StorageTypeLocal StorageType = "local"
    StorageTypeS3    StorageType = "s3"
    StorageTypeOSS   StorageType = "oss"
)
```

**本地存储实现**：

```go
// LocalStorage 本地存储实现
type LocalStorage struct {
    rootDir    string
    baseURL    string
    logger     *zap.Logger
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(rootDir, baseURL string, logger *zap.Logger) *LocalStorage {
    // 确保目录存在
    if err := os.MkdirAll(rootDir, 0755); err != nil {
        logger.Fatal("failed to create storage directory", zap.Error(err))
    }
    
    return &LocalStorage{
        rootDir: rootDir,
        baseURL: baseURL,
        logger:  logger,
    }
}

// Put 上传文件
func (s *LocalStorage) Put(ctx context.Context, key string, reader io.Reader, contentType string) error {
    // 构建完整路径
    fullPath := filepath.Join(s.rootDir, key)
    
    // 确保目录存在
    dir := filepath.Dir(fullPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    // 创建文件
    file, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // 写入内容
    _, err = io.Copy(file, reader)
    return err
}

// GetURL 获取文件URL
func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
    return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}
```

---

## 第五章：架构决策记录

### 5.1 为什么选择SQLite而非MySQL/PostgreSQL

**决策**：使用SQLite作为主数据库，而非MySQL或PostgreSQL。

**背景**：火宝短剧定位为个人创作者或小团队的创作工具，初期用户规模有限。但系统需要支持快速的开发迭代和易于部署的特性。

**考量因素**：

| 因素 | SQLite | MySQL | PostgreSQL |
|------|--------|-------|------------|
| 部署复杂度 | 无需单独服务 | 需要部署MySQL服务 | 需要部署PostgreSQL服务 |
| 运维成本 | 几乎为零 | 需要DBA维护 | 需要DBA维护 |
| 并发能力 | 适合单用户场景 | 高并发支持 | 高并发支持 |
| 数据安全 | 单文件，可靠 | 主从复制 | 主从复制 |
| 开发效率 | 开箱即用 | 需要配置连接 | 需要配置连接 |

**最终选择SQLite的理由**：

1. **部署简化**：对于目标用户（个人创作者）而言，不需要额外部署数据库服务大大降低了使用门槛
2. **数据可移植性**：SQLite数据库就是一个文件，可以轻松备份、迁移
3. **性能足够**：在预期的用户规模下（单用户或少并发），SQLite性能完全足够
4. **成本优势**：无需额外的数据库服务器成本

**风险与缓解**：

| 风险 | 缓解措施 |
|------|----------|
| 并发写入瓶颈 | 启用WAL模式，提升并发性能 |
| 数据损坏 | 定期备份，使用PRAGMA integrity_check检查 |
| 单文件大小限制 | 监控文件大小，支持数据归档 |

### 5.2 为什么选择Go而非Python/Node.js

**决策**：后端使用Go语言开发，而非Python或Node.js。

**背景**：平台需要高性能的并发处理能力（视频生成是CPU密集型任务），同时保持开发效率。

**考量因素**：

| 因素 | Go | Python | Node.js |
|------|-----|--------|---------|
| 并发模型 | goroutine天然支持 | GIL限制真正并发 | 单线程事件循环 |
| 执行效率 | 编译型，高性能 | 解释型，较低 | JIT，中等 |
| 开发效率 | 中等 | 高 | 高 |
| 生态丰富度 | 中等 | 丰富 | 丰富 |
| 部署便利性 | 单二进制 | 需要运行时 | 需要运行时 |

**最终选择Go的理由**：

1. **性能优势**：Go是编译型语言，执行效率高，适合CPU密集型的视频生成任务
2. **并发简单**：goroutine使得并发编程变得简单有效
3. **部署简单**：单二进制文件，部署极其便利
4. **静态类型**：编译时检查类型错误，提高代码质量

### 5.3 为什么选择Gin而非Echo/Fiber

**决策**：使用Gin作为Web框架，而非Echo或Fiber。

**背景**：需要一个高性能、易用、社区活跃的Web框架。

**考量因素**：

| 因素 | Gin | Echo | Fiber |
|------|-----|------|-------|
| 性能 | 高 | 高 | 最高 |
| API风格 | RESTful | 灵活 | 灵活 |
| 文档质量 | 优秀 | 优秀 | 一般 |
| 社区规模 | 最大 | 中等 | 增长中 |
| 学习曲线 | 低 | 中等 | 中等 |

**最终选择Gin的理由**：

1. **平衡性能与易用**：Gin在高性能和易用性之间取得了很好的平衡
2. **社区成熟**：最大的Go Web框架社区，问题容易找到答案
3. **中间件丰富**：大量现成的中间件可供使用
4. **学习资源多**：教程、示例代码丰富

---

## 第六章：性能优化策略

### 6.1 数据库优化

**连接池配置**：

```go
// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    MaxIdleConns    int `yaml:"max_idle"`    // 最大空闲连接
    MaxOpenConns    int `yaml:"max_open"`    // 最大打开连接
    ConnMaxLifetime int `yaml:"conn_max_lifetime"` // 连接最大生命周期（秒）
}

// 创建数据库连接
func NewDatabase(cfg *DatabaseConfig) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("./data/drama_generator.db"), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
    })
    
    // 配置连接池
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
    
    return db
}
```

**索引优化**：

```go
// 项目表的索引优化
func (ProjectEntity) TableName() string {
    return "projects"
}

// 创建索引
func (ProjectEntity) BeforeCreate(tx *gorm.DB) error {
    tx.AutoMigrate(&ProjectEntity{})
    return nil
}

// 手动创建复合索引
func MigrateIndexes(db *gorm.DB) error {
    // 按所有者查询的索引
    if !db.Migrator().HasIndex(&ProjectEntity{}, "idx_project_owner") {
        db.Exec("CREATE INDEX idx_project_owner ON projects(owner_id)")
    }
    
    // 按状态和创建时间查询的索引
    if !db.Migrator().HasIndex(&ProjectEntity{}, "idx_project_status_created") {
        db.Exec("CREATE INDEX idx_project_status_created ON projects(status, created_at DESC)")
    }
    
    return nil
}
```

**WAL模式优化**：

```go
// 启用WAL模式
func EnableWALMode(db *gorm.DB) error {
    return db.Exec("PRAGMA journal_mode=WAL").Error
}
```

### 6.2 缓存策略

**内存缓存**：

```go
// MemoryCache 内存缓存实现
type MemoryCache struct {
    cache *ccache.Cache
    ttl   time.Duration
}

// NewMemoryCache 创建缓存实例
func NewMemoryCache(ttl time.Duration) *MemoryCache {
    return &MemoryCache{
        cache: ccache.New(ccache.Configure()),
        ttl:   ttl,
    }
}

// Get 获取缓存
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
    return c.cache.Get(key)
}

// Set 设置缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}) {
    c.cache.Set(key, value, c.ttl)
}

// Delete 删除缓存
func (c *MemoryCache) Delete(ctx context.Context, key string) {
    c.cache.Delete(key)
}
```

### 6.3 并发优化

**工作池模式**：

```go
// WorkerPool 工作池
type WorkerPool struct {
    tasks     chan func()
    workers   int
    wg        sync.WaitGroup
    stopChan  chan struct{}
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
    return &WorkerPool{
        tasks:    make(chan func(), queueSize),
        workers: workers,
        stopChan: make(chan struct{}),
    }
}

// Start 启动工作池
func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go func() {
            defer p.wg.Done()
            for {
                select {
                case task := <-p.tasks:
                    task()
                case <-p.stopChan:
                    return
                }
            }
        }()
    }
}

// Submit 提交任务
func (p *WorkerPool) Submit(task func()) {
    select {
    case p.tasks <- task:
    default:
        // 队列满，阻塞等待
        p.tasks <- task
    }
}

// Stop 停止工作池
func (p *WorkerPool) Stop() {
    close(p.stopChan)
    p.wg.Wait()
}
```

---

## 第七章：可扩展性设计

### 7.1 AI服务扩展

系统设计了统一的AI服务接口，支持无缝切换和添加新的AI服务提供商。

**添加新服务提供商的步骤**：

```go
// 步骤1：实现ITextGenerationService接口
type MyAIService struct {
    apiKey string
    model  string
}

func (s *MyAIService) Generate(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
    // 实现生成逻辑
}

// 步骤2：在配置中添加配置项
type AIConfig struct {
    MyAI struct {
        APIKey string `yaml:"api_key"`
        Model  string `yaml:"model"`
    } `yaml:"myai"`
}

// 步骤3：注册到服务提供商
func (p *AIServiceProvider) RegisterMyAI(cfg *AIConfig) {
    if cfg.MyAI.APIKey != "" {
        p.textServices["myai"] = &MyAIService{
            apiKey: cfg.MyAI.APIKey,
            model:  cfg.MyAI.Model,
        }
    }
}
```

### 7.2 存储扩展

系统支持扩展新的存储后端，只需实现IStorage接口即可。

```go
// S3Storage S3存储实现
type S3Storage struct {
    client     *s3.Client
    bucket     string
    region     string
}

func (s *S3Storage) Put(ctx context.Context, key string, reader io.Reader, contentType string) error {
    // 实现S3上传逻辑
}

// StorageFactory 存储工厂
type StorageFactory struct {
    storages map[StorageType]func(*Config) IStorage
}

// NewStorageFactory 创建存储工厂
func NewStorageFactory() *StorageFactory {
    f := &StorageFactory{
        storages: make(map[StorageType]func(*Config) IStorage),
    }
    
    // 注册本地存储
    f.storages[StorageTypeLocal] = func(cfg *Config) IStorage {
        return NewLocalStorage(cfg.Local.Path, cfg.Local.BaseURL, logger)
    }
    
    // 注册S3存储
    f.storages[StorageTypeS3] = func(cfg *Config) IStorage {
        return NewS3Storage(cfg.S3.Bucket, cfg.S3.Region, logger)
    }
    
    return f
}

// Create 创建指定类型的存储
func (f *StorageFactory) Create(cfg *Config, storageType StorageType) IStorage {
    constructor, ok := f.storages[storageType]
    if !ok {
        panic("unsupported storage type")
    }
    return constructor(cfg)
}
```

---

## 练习任务

### 练习1：领域模型设计实践（⭐⭐⭐⭐）

**任务目标**：设计一个新的领域聚合，并实现完整的CRUD操作

**任务要求**：

1. 设计一个「素材库（AssetLibrary）」聚合，包含以下概念：
   - 素材库（聚合根）
   - 素材项（实体）
   - 素材标签（值对象）

2. 实现以下功能：
   - 创建素材库
   - 向素材库添加素材
   - 从素材库删除素材
   - 根据标签查询素材

**验收标准**：

| 维度 | 标准 | 权重 |
|------|------|------|
| 聚合设计 | 聚合边界清晰，引用关系正确 | 30% |
| 不变性约束 | 所有业务规则正确实现 | 30% |
| 仓储接口 | 接口定义符合DDD规范 | 20% |
| 测试覆盖 | 核心逻辑有单元测试 | 20% |

---

## 自检清单

### 架构理解自检

| 检查项 | 完成情况 |
|--------|----------|
| 能够解释四层架构的职责划分 | ☐ |
| 能够识别系统中的聚合根及其边界 | ☐ |
| 能够区分实体和值对象的使用场景 | ☐ |
| 能够说明领域事件的作用和实现方式 | ☐ |
| 能够解释应用服务和领域服务的区别 | ☐ |
| 能够设计新的聚合和仓储实现 | ☐ |

### 设计能力自检

| 检查项 | 完成情况 |
|--------|----------|
| 能够根据业务需求设计聚合 | ☐ |
| 能够识别系统中的不变性约束 | ☐ |
| 能够实现依赖注入 | ☐ |
| 能够设计可扩展的外部服务适配器 | ☐ |
| 能够制定性能优化策略 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [快速入门指南](QUICK_START.md) | 15分钟快速上手 | ⭐ |
| [用户手册](USER_MANUAL.md) | 平台功能的完整操作指南 | ⭐⭐ |
| [API接口文档](API_REFERENCE.md) | 开发者API参考手册 | ⭐⭐⭐ |
| [部署运维指南](DEPLOYMENT.md) | 生产环境部署和运维指南 | ⭐⭐ |
| [故障排查文档](TROUBLESHOOTING.md) | 常见问题和解决方案 | ⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-05 | 初始版本，包含完整架构设计文档 |
