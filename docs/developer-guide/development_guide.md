# 开发指南

> **学习目标**：完成本指南后，你将能够按照项目规范进行代码开发，掌握代码规范、测试方法和提交流程。
> 
> **前置知识**：建议先完成《代码结构详解》的学习，理解项目架构后再阅读本指南。
> 
> **难度等级**：⭐⭐⭐（进阶分析）
> 
> **预计学习时间**：4-8小时

---

## 第一章：开发规范

### 1.1 Go代码规范

火宝短剧项目遵循Go社区的通用代码规范，并在此基础上增加了项目特定的要求。

**命名规范**：

| 元素 | 规则 | 示例 |
|------|------|------|
| 包名 | 简短、全小写、使用单数 | domain、repository |
| 变量名 | 驼峰命名、避免缩写 | projectID而非pid |
| 常量名 | 全大写、下划线分隔 | ErrProjectNotFound |
| 接口名 | 以er结尾 | IRepository、IService |
| 结构体 | 驼峰命名 | ProjectService |

**函数和方法的顺序**：

```go
type ProjectService struct {
    // 依赖项放在前面
    projectRepo domain.IProjectRepository
    logger      *zap.Logger
}

// NewXxx 构造函数放在最前面
func NewProjectService(repo domain.IProjectRepository, logger *zap.Logger) *ProjectService {
    return &ProjectService{
        projectRepo: repo,
        logger:      logger,
    }
}

// 方法按功能分组，按调用顺序排列
// CRUD方法
func (s *ProjectService) Create(ctx context.Context, req *CreateRequest) (*Project, error) {}
func (s *ProjectService) GetByID(ctx context.Context, id string) (*Project, error) {}
func (s *ProjectService) Update(ctx context.Context, id string, req *UpdateRequest) error {}
func (s *ProjectService) Delete(ctx context.Context, id string) error {}

// 业务方法放在后面
func (s *ProjectService) GenerateScript(ctx context.Context, req *GenerateScriptRequest) (*Task, error) {}
```

**错误处理规范**：

```go
// 推荐：直接返回错误
func (s *ProjectService) Create(ctx context.Context, req *CreateRequest) (*Project, error) {
    if req.Name == "" {
        return nil, domain.ErrInvalidArgument
    }

    project, err := domain.NewProject(req.Name, req.Description, req.OwnerID)
    if err != nil {
        return nil, err
    }

    if err := s.projectRepo.Save(ctx, project); err != nil {
        return nil, err
    }

    return project, nil
}

// 不推荐：嵌套过深
func (s *ProjectService) Create(ctx context.Context, req *CreateRequest) (*Project, error) {
    if req.Name != "" {
        if project, err := domain.NewProject(req.Name, req.Description, req.OwnerID); err == nil {
            if err := s.projectRepo.Save(ctx, project); err == nil {
                return project, nil
            }
        }
    }
    return nil, errors.New("failed")
}
```

**注释规范**：

```go
// ProjectService 处理项目相关的业务逻辑
//
// 提供项目的创建、查询、更新、删除等操作，
// 以及与AI服务集成生成剧本等功能。
type ProjectService struct {
    // projectRepo 提供项目数据的持久化操作
    projectRepo domain.IProjectRepository
    // logger 用于记录服务运行日志
    logger *zap.Logger
}

// Create 创建新项目
//
// 该方法会验证项目名称的唯一性，
// 并将创建的项目保存到数据库中。
//
// 参数：
//   ctx - 上下文对象，用于控制请求超时和取消
//   req - 创建项目的请求参数
//
// 返回：
//   创建成功的项目对象
//   创建过程中的错误信息
func (s *ProjectService) Create(ctx context.Context, req *CreateRequest) (*Project, error) {
    // 实现代码
}
```

### 1.2 Git提交规范

**提交信息格式**：

```
<类型>(<范围>): <描述>

[可选的正文]

[可选的脚注]
```

**类型标识**：

| 类型 | 说明 | 示例 |
|------|------|------|
| feat | 新功能 | feat(script): 添加剧本模板功能 |
| fix | 修复bug | fix(api): 修复项目创建接口的验证bug |
| docs | 文档改进 | docs: 更新API接口文档 |
| style | 代码格式 | style: 运行goimports格式化 |
| refactor | 重构 | refactor: 重构聚合根创建逻辑 |
| perf | 性能优化 | perf: 优化数据库查询性能 |
| test | 添加测试 | test: 增加项目服务单元测试 |
| chore | 构建工具 | chore: 更新CI配置文件 |

**提交信息示例**：

```
feat(script): 支持从模板创建剧本

- 添加剧本模板定义和存储
- 实现从模板生成剧本的功能
- 添加相关的单元测试

Closes #123
```

### 1.3 代码审查标准

**提交PR前的自检清单**：

| 检查项 | 说明 | 状态 |
|--------|------|------|
| 代码编译通过 | go build ./... 无错误 | ☐ |
| 测试通过 | go test ./... 全部通过 | ☐ |
| 代码格式化 | gofmt -w . 无变化 | ☐ |
| 静态检查通过 | golangci-lint run 无警告 | ☐ |
| 注释完整 | 公开API都有文档注释 | ☐ |
| 变更日志更新 | CHANGELOG.md已更新 | ☐ |
| PR描述完整 | 包含目的、变更、测试说明 | ☐ |

---

## 第二章：功能开发流程

### 2.1 开发任务分解

进行功能开发时，需要先将大任务分解为可执行的小任务。

**任务分解模板**：

```markdown
## 功能名称：XXX

### 任务分解

- [ ] 任务1：设计领域模型
  - 涉及文件：domain/aggregates/
  - 依赖：无
  - 验收标准：聚合根定义完整

- [ ] 任务2：实现仓储接口
  - 涉及文件：domain/repositories/
  - 依赖：任务1
  - 验收标准：接口定义符合DDD规范

- [ ] 任务3：实现仓储
  - 涉及文件：infrastructure/persistence/
  - 依赖：任务2
  - 验收标准：能正确持久化数据

- [ ] 任务4：实现应用服务
  - 涉及文件：application/services/
  - 依赖：任务1、任务3
  - 验收标准：业务逻辑正确

- [ ] 任务5：实现API接口
  - 涉及文件：api/handlers/
  - 依赖：任务4
  - 验收标准：API响应符合规范
```

### 2.2 代码实现步骤

**步骤一：定义领域模型**

```go
// domain/aggregates/my_feature.go

// MyFeature 我的功能聚合根
type MyFeature struct {
    ID          string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    Name        string    `json:"name" gorm:"type:varchar(100);not null"`
    Description string    `json:"description" gorm:"type:text"`
    Status      string    `json:"status" gorm:"type:varchar(20);not null;default:draft"`
    
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// NewMyFeature 创建聚合根工厂函数
func NewMyFeature(name, description string) (*MyFeature, error) {
    if strings.TrimSpace(name) == "" {
        return nil, ErrNameRequired
    }
    
    return &MyFeature{
        ID:          uuid.New().String(),
        Name:        name,
        Description: description,
        Status:      StatusDraft,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}
```

**步骤二：定义仓储接口**

```go
// domain/repositories/my_feature_repository.go

// IMyFeatureRepository 我的功能仓储接口
type IMyFeatureRepository interface {
    // 查询操作
    FindByID(ctx context.Context, id string) (*MyFeature, error)
    FindByName(ctx context.Context, name string) ([]*MyFeature, error)
    FindByStatus(ctx context.Context, status string) ([]*MyFeature, error)
    
    // 持久化操作
    Save(ctx context.Context, feature *MyFeature) error
    Update(ctx context.Context, feature *MyFeature) error
    Delete(ctx context.Context, id string) error
}
```

**步骤三：实现仓储**

```go
// infrastructure/persistence/my_feature_repository.go

// GormMyFeatureRepository 基于GORM的仓储实现
type GormMyFeatureRepository struct {
    db *gorm.DB
}

// NewMyFeatureRepository 创建仓储实例
func NewMyFeatureRepository(db *gorm.DB) *GormMyFeatureRepository {
    return &GormMyFeatureRepository{db: db}
}

// FindByID 根据ID查询
func (r *GormMyFeatureRepository) FindByID(ctx context.Context, id string) (*MyFeature, error) {
    var entity MyFeatureEntity
    
    err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, err
    }
    
    return r.toDomain(&entity), nil
}

// Save 保存
func (r *GormMyFeatureRepository) Save(ctx context.Context, feature *MyFeature) error {
    entity := r.toEntity(feature)
    return r.db.WithContext(ctx).Create(entity).Error
}
```

**步骤四：实现应用服务**

```go
// application/services/my_feature_service.go

// MyFeatureAppService 我的功能应用服务
type MyFeatureAppService struct {
    repo     domain.IMyFeatureRepository
    eventBus domain.IEventBus
    logger   *zap.Logger
}

// Create 创建
func (s *MyFeatureAppService) Create(
    ctx context.Context,
    req *CreateRequest,
) (*MyFeature, error) {
    feature, err := domain.NewMyFeature(req.Name, req.Description)
    if err != nil {
        return nil, err
    }
    
    if err := s.repo.Save(ctx, feature); err != nil {
        return nil, err
    }
    
    feature.RegisterEvent(domain.MyFeatureCreatedEvent{
        FeatureID: feature.ID,
    })
    feature.PublishEvents(ctx, s.eventBus)
    
    return feature, nil
}
```

**步骤五：实现API接口**

```go
// api/handlers/my_feature_handler.go

// MyFeatureHandler 我的功能HTTP处理器
type MyFeatureHandler struct {
    service *application.MyFeatureAppService
}

// RegisterRoutes 注册路由
func (h *MyFeatureHandler) RegisterRoutes(r *gin.Engine) {
    features := r.Group("/api/v1/my-features")
    {
        features.POST("", h.Create)
        features.GET("/:id", h.GetByID)
    }
}

// Create 创建
func (h *MyFeatureHandler) Create(c *gin.Context) {
    var req request.CreateMyFeatureRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Fail(response.ErrInvalidParam))
        return
    }
    
    feature, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Fail(response.ErrInternal))
        return
    }
    
    c.JSON(http.StatusCreated, response.Success(feature))
}
```

---

## 第三章：测试规范

### 3.1 单元测试

**测试覆盖率要求**：

| 代码类型 | 最低覆盖率 |
|----------|-----------|
| 核心业务逻辑 | 80% |
| API接口 | 90% |
| 工具函数 | 100% |

**单元测试示例**：

```go
// project_service_test.go

func TestProjectService_Create(t *testing.T) {
    // Arrange - 准备测试数据
    mockRepo := mocks.NewMockProjectRepository(gomock.NewController(t))
    logger, _ := zap.NewDevelopment()
    service := NewProjectService(mockRepo, logger)
    
    req := &CreateRequest{
        Name:        "测试项目",
        Description: "测试描述",
    }
    
    // Expect - 设置预期行为
    mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
    
    // Act - 执行测试
    project, err := service.Create(context.Background(), req)
    
    // Assert - 验证结果
    assert.NoError(t, err)
    assert.Equal(t, "测试项目", project.Name)
    assert.Equal(t, "测试描述", project.Description)
}
```

### 3.2 集成测试

```go
// api_test.go

func TestCreateProject(t *testing.T) {
    // 启动测试服务器
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 模拟处理逻辑
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"code":0,"data":{"id":"test-id"}}`))
    }))
    defer ts.Close()
    
    // 发送请求
    resp, err := http.Post(
        ts.URL+"/api/v1/projects",
        "application/json",
        strings.NewReader(`{"name":"测试项目"}`),
    )
    
    // 验证响应
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

### 3.3 测试工具

项目使用以下测试工具：

| 工具 | 用途 |
|------|------|
| testing | Go内置测试框架 |
| gomock | 生成Mock对象 |
| testify | 断言和Mock工具 |
| httptest | HTTP测试工具 |

---

## 第四章：调试技巧

### 4.1 本地调试

**使用Delve调试器**：

```bash
# 安装Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试运行
dlv debug main.go

# 设置断点
(dlv) break main.go:45

# 继续执行
(dlv) continue

# 单步执行
(dlv) next

# 查看变量
(dlv) locals

# 退出
(dlv) quit
```

### 4.2 日志调试

**开发模式日志**：

```go
// 在代码中添加日志
logger.Debug("开始处理请求",
    zap.String("project_id", projectID),
    zap.Any("request", req),
)

logger.Info("项目创建成功",
    zap.String("project_id", project.ID),
    zap.String("name", project.Name),
)

logger.Error("处理失败",
    zap.Error(err),
    zap.String("project_id", projectID),
)
```

### 4.3 API调试

**使用curl测试API**：

```bash
# 创建项目
curl -X POST http://localhost:5678/api/v1/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"name":"测试项目","description":"测试描述"}'

# 获取项目列表
curl http://localhost:5678/api/v1/projects \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取单个项目
curl http://localhost:5678/api/v1/projects/{id} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 第五章：问题排查

### 5.1 常见错误

**编译错误**：

```bash
# 检查依赖
go mod tidy

# 清理缓存
go clean -cache

# 检查语法
go build -n
```

**运行时错误**：

```bash
# 启用详细日志
DEBUG=true ./huobao-drama

# 查看错误堆栈
# 日志中会显示详细的错误信息和堆栈跟踪
```

### 5.2 调试方法

**分步排查**：

1. 确认请求格式正确
2. 确认认证Token有效
3. 检查请求参数验证
4. 检查应用服务逻辑
5. 检查仓储操作
6. 检查数据库状态

---

## 练习任务

### 练习1：功能开发实践（⭐⭐⭐）

**任务目标**：按照规范实现一个新功能

**任务要求**：选择一个小功能（如「收藏项目」），按照开发流程完成从领域模型到API接口的完整实现，编写单元测试，提交PR并通过审查。

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| 领域模型设计合理 | ☐ |
| 仓储实现正确 | ☐ |
| 应用服务逻辑正确 | ☐ |
| API接口符合规范 | ☐ |
| 测试覆盖率达到要求 | ☐ |
| PR描述完整清晰 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [开发环境搭建](environment.md) | 本地开发环境配置 | ⭐ |
| [代码结构详解](code_structure.md) | 项目架构和模块划分 | ⭐⭐ |
| [高级开发主题](advanced_dev.md) | 架构决策和最佳实践 | ⭐⭐⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本，包含完整开发指南 |
