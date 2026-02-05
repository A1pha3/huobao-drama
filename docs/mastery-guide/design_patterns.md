# 设计模式

> **学习目标**：掌握火宝短剧项目中使用的核心设计模式，能够在日常开发中灵活运用这些模式解决实际问题。
> 
> **前置知识**：建议先完成《架构深度解析》的学习，理解系统架构后再学习设计模式。
> 
> **难度等级**：⭐⭐⭐⭐（专家设计）
> 
> **预计学习时间**：8-16小时

---

## 第一章：创建型模式

### 1.1 工厂模式

**意图**：封装对象的创建逻辑，客户端无需关心对象是如何创建的。

**问题**：当对象创建涉及复杂逻辑（如条件判断、依赖注入）时，将创建逻辑分散在代码中会导致难以维护。

**解决方案**：将创建逻辑集中到工厂类中，通过统一接口创建对象。

**火宝短剧中的应用**：

```go
// storage_factory.go

// IStorage 存储接口
type IStorage interface {
    Put(ctx context.Context, key string, data []byte) error
    Get(ctx context.Context, key string) ([]byte, error)
    Delete(ctx context.Context, key string) error
}

// StorageType 存储类型
type StorageType string

const (
    LocalStorage StorageType = "local"
    S3Storage    StorageType = "s3"
    OSSStorage   StorageType = "oss"
)

// StorageFactory 存储工厂
type StorageFactory struct {
    factories map[StorageType]func(*Config) IStorage
}

// NewStorageFactory 创建工厂
func NewStorageFactory() *StorageFactory {
    return &StorageFactory{
        factories: map[StorageType]func(*Config) IStorage{
            LocalStorage: NewLocalStorage,
            S3Storage:    NewS3Storage,
            OSSStorage:   NewOSSStorage,
        },
    }
}

// Create 创建存储实例
func (f *StorageFactory) Create(cfg *Config, storageType StorageType) (IStorage, error) {
    factory, ok := f.factories[storageType]
    if !ok {
        return nil, fmt.Errorf("unsupported storage type: %s", storageType)
    }
    return factory(cfg), nil
}
```

**使用示例**：

```go
factory := NewStorageFactory()

// 创建本地存储
localStorage, err := factory.Create(config, LocalStorage)

// 创建S3存储
s3Storage, err := factory.Create(config, S3Storage)
```

**使用场景**：

| 场景 | 示例 |
|------|------|
| 多存储后端支持 | Local、S3、OSS |
| 多AI服务支持 | OpenAI、DeepSeek、豆包 |
| 多数据库支持 | SQLite、MySQL、PostgreSQL |

### 1.2 单例模式

**意图**：确保一个类只有一个实例，并提供全局访问点。

**问题**：某些资源（如配置、连接池）只需要一个实例，重复创建会浪费资源或导致问题。

**解决方案**：控制实例的创建过程，确保只有一个实例存在。

**火宝短剧中的应用**：

```go
// config.go

type Config struct {
    App      AppConfig
    Server   ServerConfig
    Database DatabaseConfig
    AI       AIConfig
    Storage  StorageConfig
    
    instance *Config
    once    sync.Once
}

// GetInstance 获取配置单例
func (c *Config) GetInstance(path string) (*Config, error) {
    c.once.Do(func() {
        c.instance = &Config{}
        // 初始化配置...
    })
    return c.instance, nil
}
```

**使用场景**：

| 场景 | 示例 |
|------|------|
| 配置管理 | 全局配置单例 |
| 连接池 | 数据库连接池 |
| 日志器 | 全局日志实例 |

### 1.3 构建者模式

**意图**：将复杂对象的构建过程与表示分离，使同样的构建过程可以创建不同的表示。

**问题**：当对象有多个可选参数时，构造函数参数列表会变得很长且难以使用。

**解决方案**：使用链式调用的Builder来逐步设置参数，最后构建对象。

**火宝短剧中的应用**：

```go
// project_builder.go

// ProjectBuilder 项目构建者
type ProjectBuilder struct {
    name        string
    description string
    ownerID     string
    settings    *Settings
    tags        []string
}

// NewProjectBuilder 创建构建者
func NewProjectBuilder() *ProjectBuilder {
    return &ProjectBuilder{
        settings: &Settings{},
        tags:     []string{},
    }
}

// Name 设置名称
func (b *ProjectBuilder) Name(name string) *ProjectBuilder {
    b.name = name
    return b
}

// Description 设置描述
func (b *ProjectBuilder) Description(desc string) *ProjectBuilder {
    b.description = desc
    return b
}

// OwnerID 设置所有者ID
func (b *ProjectBuilder) OwnerID(id string) *ProjectBuilder {
    b.ownerID = id
    return b
}

// Settings 设置配置
func (b *ProjectBuilder) Settings(settings *Settings) *ProjectBuilder {
    b.settings = settings
    return b
}

// Tags 设置标签
func (b *ProjectBuilder) Tags(tags ...string) *ProjectBuilder {
    b.tags = append(b.tags, tags...)
    return b
}

// Build 构建项目
func (b *ProjectBuilder) Build() (*Project, error) {
    if b.name == "" {
        return nil, ErrNameRequired
    }
    if b.ownerID == "" {
        return nil, ErrOwnerIDRequired
    }
    
    return &Project{
        ID:          uuid.New().String(),
        Name:        b.name,
        Description: b.description,
        OwnerID:     b.ownerID,
        Settings:    b.settings,
        Tags:        b.tags,
        Status:      ProjectStatusDraft,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}
```

**使用示例**：

```go
project, err := NewProjectBuilder().
    Name("我的短剧").
    Description("AI生成的短剧项目").
    OwnerID("user-123").
    Settings(&Settings{Resolution: "1080p"}).
    Tags("AI", "短剧").
    Build()
```

---

## 第二章：结构型模式

### 2.1 适配器模式

**意图**：将一个类的接口转换成客户期望的另一个接口，使原本不兼容的类可以合作。

**问题**：外部系统的接口与项目内部接口不一致，需要适配才能使用。

**解决方案**：创建一个适配器类，包装外部系统，暴露项目期望的接口。

**火宝短剧中的应用**：

```go
// openai_adapter.go

// ITextGenerationService 文本生成服务接口（项目内部接口）
type ITextGenerationService interface {
    Generate(ctx context.Context, req *GenerationRequest) (*GenerationResponse, error)
}

// OpenAIClient OpenAI官方客户端（外部接口）
type OpenAIClient struct {
    APIKey string
}

// GenerationRequest 生成请求（外部接口）
type OpenAIGenerationRequest struct {
    Model       string `json:"model"`
    Messages    []struct {
        Role    string `json:"role"`
        Content string `json:"content"`
    } `json:"messages"`
    MaxTokens   int     `json:"max_tokens"`
    Temperature float64 `json:"temperature"`
}

// OpenAIAdapter 适配器
type OpenAIAdapter struct {
    client *OpenAIClient
}

// Generate 实现ITextGenerationService接口
func (a *OpenAIAdapter) Generate(ctx context.Context, req *GenerationRequest) (*GenerationResponse, error) {
    // 转换请求格式
    openaiReq := OpenAIGenerationRequest{
        Model:       req.Model,
        MaxTokens:   req.MaxTokens,
        Temperature: req.Temperature,
    }
    for _, msg := range req.Messages {
        openaiReq.Messages = append(openaiReq.Messages, struct {
            Role    string `json:"role"`
            Content string `json:"content"`
        }{
            Role:    msg.Role,
            Content: msg.Content,
        })
    }
    
    // 调用外部接口
    resp, err := a.client.CreateChatCompletion(ctx, &openaiReq)
    if err != nil {
        return nil, err
    }
    
    // 转换响应格式
    return &GenerationResponse{
        Content:      resp.Choices[0].Message.Content,
        Usage:        TokenUsage(resp.Usage),
        Model:        resp.Model,
    }, nil
}
```

**使用场景**：

| 场景 | 示例 |
|------|------|
| 多AI服务适配 | OpenAI、DeepSeek、豆包 |
| 多存储后端适配 | Local、S3、OSS |
| 遗留系统集成 | 适配旧API |

### 2.2 装饰器模式

**意图**：动态地给对象添加额外的职责，比继承更灵活。

**问题**：需要在运行时给对象添加功能，但继承会导致类爆炸。

**解决方案**：使用装饰器包装原对象，在保持接口不变的情况下添加功能。

**火宝短剧中的应用**：

```go
// logging_decorator.go

// IScriptService 脚本服务接口
type IScriptService interface {
    Generate(ctx context.Context, prompt string) (*Script, error)
}

// ScriptService 基础服务
type ScriptService struct {
    generator AIGenerator
}

// LoggingDecorator 日志装饰器
type LoggingDecorator struct {
    service IScriptService
    logger  *zap.Logger
}

// Generate 添加日志功能
func (d *LoggingDecorator) Generate(ctx context.Context, prompt string) (*Script, error) {
    start := time.Now()
    
    d.logger.Info("开始生成剧本",
        zap.String("prompt", prompt),
    )
    
    script, err := d.service.Generate(ctx, prompt)
    
    duration := time.Since(start)
    if err != nil {
        d.logger.Error("生成剧本失败",
            zap.String("prompt", prompt),
            zap.Error(err),
            zap.Duration("duration", duration),
        )
        return nil, err
    }
    
    d.logger.Info("生成剧本成功",
        zap.String("script_id", script.ID),
        zap.Duration("duration", duration),
    )
    
    return script, nil
}

// CachingDecorator 缓存装饰器
type CachingDecorator struct {
    service IScriptService
    cache   Cache
}

// Generate 添加缓存功能
func (d *CachingDecorator) Generate(ctx context.Context, prompt string) (*Script, error) {
    cacheKey := fmt.Sprintf("script:%s", hashPrompt(prompt))
    
    if cached, err := d.cache.Get(ctx, cacheKey); err == nil {
        return cached.(*Script), nil
    }
    
    script, err := d.service.Generate(ctx, prompt)
    if err != nil {
        return nil, err
    }
    
    d.cache.Set(ctx, cacheKey, script, time.Hour)
    
    return script, nil
}
```

**使用示例**：

```go
// 基础服务
baseService := &ScriptService{generator: openai}

// 添加日志装饰器
loggedService := &LoggingDecorator{
    service: baseService,
    logger:  logger,
}

// 添加缓存装饰器
cachedService := &CachingDecorator{
    service: loggedService,
    cache:   cache,
}

// 使用（调用链：缓存 → 日志 → 基础服务）
script, err := cachedService.Generate(ctx, "生成一个短剧")
```

### 2.3 代理模式

**意图**：为其他对象提供一种代理以控制对这个对象的访问。

**问题**：需要在访问对象前后添加额外的逻辑（如缓存、权限、日志）。

**解决方案**：创建一个代理对象，包装真实对象，在适当的时候委托给真实对象。

**火宝短剧中的应用**：

```go
// repository_proxy.go

// IProjectRepository 项目仓储接口
type IProjectRepository interface {
    FindByID(ctx context.Context, id string) (*Project, error)
    Save(ctx context.Context, project *Project) error
}

// GormProjectRepository 真实仓储
type GormProjectRepository struct {
    db *gorm.DB
}

// ProjectRepositoryProxy 代理仓储
type ProjectRepositoryProxy struct {
    real    *GormProjectRepository
    cache   Cache
    logger  *zap.Logger
}

// FindByID 代理方法
func (p *ProjectRepositoryProxy) FindByID(ctx context.Context, id string) (*Project, error) {
    // 1. 检查缓存
    if cached, err := p.cache.Get(ctx, "project:"+id); err == nil {
        p.logger.Debug("缓存命中",
            zap.String("project_id", id),
        )
        return cached.(*Project), nil
    }
    
    // 2. 调用真实方法
    project, err := p.real.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    p.cache.Set(ctx, "project:"+id, project, time.Hour)
    
    return project, nil
}
```

---

## 第三章：行为型模式

### 3.1 观察者模式

**意图**：定义对象间的一对多依赖关系，当一个对象状态改变时，所有依赖它的对象都得到通知。

**问题**：一个对象的状态变化需要通知多个其他对象。

**解决方案**：建立发布-订阅机制，让对象可以订阅感兴趣的事件。

**火宝短剧中的应用**：

```go
// event_bus.go

// IEventBus 事件总线接口
type IEventBus interface {
    Publish(ctx context.Context, event DomainEvent) error
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handlerID string) error
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event DomainEvent) error

// InMemoryEventBus 内存事件总线
type InMemoryEventBus struct {
    mu       sync.RWMutex
    handlers map[string][]*HandlerRegistration
}

// HandlerRegistration 处理器注册
type HandlerRegistration struct {
    id      string
    handler EventHandler
}

// Publish 发布事件
func (b *InMemoryEventBus) Publish(ctx context.Context, event DomainEvent) error {
    b.mu.RLock()
    handlers, ok := b.handlers[event.EventType()]
    b.mu.RUnlock()
    
    if !ok {
        return nil  // 没有订阅者
    }
    
    for _, reg := range handlers {
        if err := reg.handler(ctx, event); err != nil {
            // 记录错误但继续处理其他处理器
            log.Printf("事件处理失败: %v", err)
        }
    }
    
    return nil
}

// Subscribe 订阅事件
func (b *InMemoryEventBus) Subscribe(eventType string, handlerID string, handler EventHandler) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    b.handlers[eventType] = append(b.handlers[eventType], &HandlerRegistration{
        id:      handlerID,
        handler: handler,
    })
    
    return nil
}
```

**使用示例**：

```go
// 创建事件总线
eventBus := NewInMemoryEventBus()

// 订阅项目创建事件
eventBus.Subscribe("project.created", "project-notification", 
    func(ctx context.Context, event DomainEvent) error {
        e := event.(*ProjectCreatedEvent)
        // 发送通知
        notificationService.Send(e.OwnerID, "项目已创建: "+e.ProjectName)
        return nil
    },
)

// 订阅剧本创建事件
eventBus.Subscribe("script.created", "script-stats",
    func(ctx context.Context, event DomainEvent) error {
        // 更新统计
        statsService.IncrementScriptCount()
        return nil
    },
)

// 发布事件（在应用服务中）
func (s *ProjectAppService) Create(ctx context.Context, req *CreateRequest) (*Project, error) {
    project, err := domain.NewProject(req.Name, req.Description, req.OwnerID)
    if err != nil {
        return nil, err
    }
    
    // 保存项目
    if err := s.repo.Save(ctx, project); err != nil {
        return nil, err
    }
    
    // 发布事件
    project.RegisterEvent(&ProjectCreatedEvent{
        ProjectID:   project.ID,
        ProjectName: project.Name,
        OwnerID:     req.OwnerID,
        Timestamp:   time.Now(),
    })
    project.PublishEvents(ctx, s.eventBus)
    
    return project, nil
}
```

### 3.2 策略模式

**意图**：定义一系列算法，将每个算法封装起来，使它们可以互相替换。

**问题**：同一个问题有多种算法，需要在运行时选择算法。

**解决方案**：将算法封装成独立的策略类，通过统一接口使用。

**火宝短剧中的应用**：

```go
// pricing_strategy.go

// PricingStrategy 定价策略接口
type PricingStrategy interface {
    Calculate(usage *Usage) (*Pricing, error)
}

// Usage 使用量
type Usage struct {
    APIRequests   int
    StorageGB     float64
    VideoMinutes int
}

// Pricing 价格信息
type Pricing struct {
    TotalCost   float64
    Breakdown   map[string]float64
}

// TokenUsagePricing 按Token使用量计费
type TokenUsagePricing struct {
    PerTokenCost float64
}

func (p *TokenUsagePricing) Calculate(usage *Usage) (*Pricing, error) {
    tokenCost := float64(usage.APIRequests) * p.PerTokenCost
    
    return &Pricing{
        TotalCost: tokenCost,
        Breakdown: map[string]float64{
            "token_usage": tokenCost,
        },
    }, nil
}

// StoragePricing 按存储量计费
type StoragePricing struct {
    PerGBPerMonth float64
}

func (p *StoragePricing) Calculate(usage *Usage) (*Pricing, error) {
    storageCost := usage.StorageGB * p.PerGBPerMonth
    
    return &Pricing{
        TotalCost: storageCost,
        Breakdown: map[string]float64{
            "storage": storageCost,
        },
    }, nil
}

// PricingContext 定价上下文
type PricingContext struct {
    strategy PricingStrategy
}

// Calculate 计算价格
func (c *PricingContext) Calculate(usage *Usage) (*Pricing, error) {
    return c.strategy.Calculate(usage)
}

// SetStrategy 设置策略
func (c *PricingContext) SetStrategy(strategy PricingStrategy) {
    c.strategy = strategy
}
```

**使用示例**：

```go
// 根据用户类型选择策略
var strategy PricingStrategy
switch user.Type {
case "token":
    strategy = &TokenUsagePricing{PerTokenCost: 0.001}
case "storage":
    strategy = &StoragePricing{PerGBPerMonth: 0.1}
default:
    strategy = &TokenUsagePricing{PerTokenCost: 0.001}
}

pricing := &PricingContext{strategy: strategy}
result, err := pricing.Calculate(usage)
```

### 3.3 状态模式

**意图**：允许一个对象在其内部状态改变时改变其行为。

**问题**：对象的行为随着内部状态的改变而改变，使用条件语句会导致复杂难维护的代码。

**解决方案**：将每种状态封装成独立的状态类，状态转换由状态类处理。

**火宝短剧中的应用**：

```go
// video_state.go

// VideoState 视频状态接口
type VideoState interface {
    Start(generation *VideoGeneration) error
    Process(generation *VideoGeneration) error
    Complete(generation *VideoGeneration) error
    Fail(generation *VideoGeneration, err error) error
    Cancel(generation *VideoGeneration) error
}

// PendingState 待处理状态
type PendingState struct{}

func (s *PendingState) Start(generation *VideoGeneration) error {
    generation.Status = VideoStatusProcessing
    generation.StartedAt = time.Now()
    return nil
}

func (s *PendingState) Process(generation *VideoGeneration) error {
    return ErrInvalidStateTransition
}

// ProcessingState 处理中状态
type ProcessingState struct{}

func (s *ProcessingState) Start(generation *VideoGeneration) error {
    return ErrInvalidStateTransition
}

func (s *ProcessingState) Process(generation *VideoGeneration) error {
    // 执行视频生成...
    return nil
}

func (s *ProcessingState) Complete(generation *VideoGeneration) error {
    generation.Status = VideoStatusCompleted
    generation.CompletedAt = time.Now()
    return nil
}

// VideoGeneration 视频生成聚合根
type VideoGeneration struct {
    ID        string
    Status    VideoStatus
    state     VideoState
    Progress  float64
    OutputURL string
    Error    string
}

// SetState 设置状态
func (v *VideoGeneration) SetState(state VideoState) {
    v.state = state
}

// Process 处理视频生成
func (v *VideoGeneration) Process() error {
    return v.state.Process(v)
}

// Complete 完成
func (v *VideoGeneration) Complete() error {
    return v.state.Complete(v)
}
```

---

## 第四章：架构模式

### 4.1 依赖注入

**意图**：通过外部注入而非内部创建的方式提供对象的依赖。

**问题**：对象内部创建依赖会导致紧耦合，难以测试和替换实现。

**解决方案**：通过构造函数或容器注入依赖。

**火宝短剧中的应用**：

```go
// service.go

// ProjectService 项目应用服务
type ProjectService struct {
    // 依赖通过构造函数注入
    projectRepo    domain.IProjectRepository
    scriptService domain.IScriptGenerationService
    eventBus      domain.IEventBus
    logger        *zap.Logger
    cache         Cache
}

// NewProjectService 创建服务（依赖注入点）
func NewProjectService(
    projectRepo domain.IProjectRepository,
    scriptService domain.IScriptGenerationService,
    eventBus domain.IEventBus,
    logger *zap.Logger,
) *ProjectService {
    return &ProjectService{
        projectRepo:    projectRepo,
        scriptService: scriptService,
        eventBus:      eventBus,
        logger:        logger,
    }
}
```

**依赖注入容器**：

```go
// container.go

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

// Initialize 初始化容器
func InitializeContainer(cfg *Config) *Container {
    container := NewContainer()
    
    // 初始化依赖
    db := NewDatabase(cfg.Database)
    
    projectRepo := persistence.NewProjectRepository(db)
    scriptRepo := persistence.NewScriptRepository(db)
    
    // 注册服务
    container.Register("project_repo", projectRepo)
    container.Register("script_repo", scriptRepo)
    container.Register("project_service", NewProjectService(
        projectRepo,
        scriptRepo,
        nil,  // eventBus
        nil,  // logger
    ))
    
    return container
}
```

### 4.2 事件溯源

**意图**：通过记录对象状态变化的事件来代替存储当前状态。

**问题**：难以追溯状态变化历史，复杂的状态转换难以实现。

**解决方案**：存储所有状态变化事件，通过重放事件来重建任意时刻的状态。

**火宝短剧中的应用**：

```go
// event_sourced_project.go

// EventSourcedProject 事件溯源项目
type EventSourcedProject struct {
    ID          string
    Name        string
    Description string
    Status      ProjectStatus
    Version     int
    events     []DomainEvent
}

// ProjectCreatedEvent 项目创建事件
type ProjectCreatedEvent struct {
    ProjectID   string
    Name        string
    Description string
    OwnerID     string
    Timestamp   time.Time
}

// Apply 应用事件
func (p *EventSourcedProject) Apply(event DomainEvent) {
    p.Version++
    
    switch e := event.(type) {
    case *ProjectCreatedEvent:
        p.ID = e.ProjectID
        p.Name = e.Name
        p.Description = e.Description
        p.Status = ProjectStatusDraft
    }
    
    p.events = append(p.events, event)
}

// LoadFromEvents 从事件重建状态
func LoadFromEvents(events []DomainEvent) *EventSourcedProject {
    project := &EventSourcedProject{}
    for _, event := range events {
        project.Apply(event)
    }
    return project
}

// GetEventStream 获取事件流
func (p *EventSourcedProject) GetEventStream() []DomainEvent {
    return p.events
}
```

---

## 第五章：模式选择指南

### 5.1 模式选择决策树

```
需要创建对象？
├─→ 简单创建 → 直接new
├─→ 创建逻辑复杂 → 工厂模式
├─→ 需要单例 → 单例模式
└─→ 参数很多 → 构建者模式

需要适配接口？
├─→ 外部系统 → 适配器模式
└─→ 运行时增强 → 装饰器模式

需要控制访问？
├─→ 添加额外逻辑 → 代理模式
├─→ 一对多通知 → 观察者模式
├─→ 状态相关行为 → 状态模式
└─→ 算法可变 → 策略模式
```

### 5.2 模式使用注意事项

| 注意事项 | 说明 |
|----------|------|
| 避免过度设计 | 不要为了用模式而用模式 |
| 考虑复杂度 | 模式会增加代码复杂度 |
| 团队共识 | 确保团队理解使用的模式 |
| 文档记录 | 记录为什么选择这个模式 |

---

## 练习任务

### 练习1：设计模式应用实践（⭐⭐⭐⭐）

**任务目标**：在一个新功能中应用合适的设计模式

**任务要求**：设计一个「导出服务」，支持导出为PDF、Word、Markdown等格式，应用合适的创建型、结构型和行为型模式。

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| 选择合适的创建型模式 | ☐ |
| 选择合适的行为型模式 | ☐ |
| 代码清晰易维护 | ☐ |
| 模式使用合理 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [基础知识回顾](fundamentals.md) | 核心技术概念回顾 | ⭐ |
| [架构深度解析](architecture_analysis.md) | DDD架构和设计原理 | ⭐⭐ |
| [性能优化](performance.md) | 性能分析和优化方法 | ⭐⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本，包含设计模式完整内容 |
