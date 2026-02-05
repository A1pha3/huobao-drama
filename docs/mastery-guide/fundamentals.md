# 基础知识回顾

> **学习目标**：通过本指南回顾项目涉及的核心技术概念，为深入学习系统架构打下坚实基础。
> 
> **前置知识**：具备基本的编程概念和软件开发经验。
> 
> **难度等级**：⭐（入门教程）
> 
> **预计学习时间**：1-2小时

---

## 第一章：编程基础

### 1.1 面向对象编程

**面向对象编程（OOP）**是一种以对象为中心的编程范式。对象是类的实例，类定义了对象的属性和行为。OOP的四大支柱是封装、继承、多态和抽象。

**封装**是将数据和操作数据的方法绑定在一起，隐藏内部实现细节。在Go中，通过结构体和方法实现封装：

```go
// Person 结构体（类）
type Person struct {
    name string  // 私有字段（首字母小写）
    age  int
}

// 方法（行为）
func (p *Person) GetName() string {
    return p.name  // 可以访问私有字段
}

func (p *Person) SetName(name string) {
    p.name = name
}
```

**继承**允许一个类继承另一个类的属性和方法。在Go中，使用组合模拟继承：

```go
// Animal 基础结构体
type Animal struct {
    Name string
}

// Eat 方法
func (a *Animal) Eat() {
    fmt.Println("Eating...")
}

// Dog 组合Animal
type Dog struct {
    Animal  // 组合
    Breed  string
}
```

**多态**是同一操作作用于不同对象产生不同行为。在Go中，通过接口实现多态：

```go
// Writer 接口
type Writer interface {
    Write(p []byte) (n int, err error)
}

// File 实现了Writer接口
type File struct{}

func (f *File) Write(p []byte) (n int, err error) {
    // 实现写入逻辑
    return len(p), nil
}
```

### 1.2 函数式编程

**函数式编程**是一种将计算视为函数求值的编程范式。核心概念包括纯函数、高阶函数和不可变性。

**纯函数**是相同输入总是产生相同输出，且没有副作用的函数：

```go
// 纯函数
func add(a, b int) int {
    return a + b
}

// 非纯函数（有副作用）
var counter int

func increment() int {
    counter++
    return counter
}
```

**高阶函数**是接受函数作为参数或返回函数的函数：

```go
// 接受函数作为参数
func filter(numbers []int, predicate func(int) bool) []int {
    result := []int{}
    for _, n := range numbers {
        if predicate(n) {
            result = append(result, n)
        }
    }
    return result
}

// 使用
evens := filter([]int{1, 2, 3, 4}, func(n int) bool {
    return n%2 == 0
})
```

### 1.3 并发编程

**并发**是同时处理多个任务的能力。Go语言通过goroutine和channel实现轻量级并发。

**Goroutine**是轻量级线程，由Go运行时管理：

```go
// 启动goroutine
go func() {
    fmt.Println("Running in background")
}()

// 并发执行多个任务
go task1()
go task2()
go task3()
```

**Channel**是goroutine之间通信的管道：

```go
// 创建channel
ch := make(chan int, 10)  // 带缓冲的channel

// 发送数据
ch <- 42

// 接收数据
value := <-ch
```

**同步机制**：

```go
// WaitGroup 等待一组goroutine完成
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("Task %d completed\n", id)
    }(i)
}
wg.Wait()  // 等待所有任务完成
```

---

## 第二章：Web开发基础

### 2.1 HTTP协议

**HTTP（超文本传输协议）**是Web应用通信的基础协议。HTTP请求包含请求行、请求头和请求体；HTTP响应包含状态行、响应头和响应体。

**HTTP方法**：

| 方法 | 语义 | 幂等 | 安全性 |
|------|------|------|--------|
| GET | 获取资源 | 是 | 是 |
| POST | 创建资源 | 否 | 否 |
| PUT | 更新资源 | 是 | 否 |
| DELETE | 删除资源 | 是 | 否 |

**HTTP状态码**：

| 状态码 | 含义 | 示例 |
|--------|------|------|
| 200 | 成功 | OK |
| 201 | 创建成功 | Created |
| 400 | 客户端错误 | Bad Request |
| 401 | 未授权 | Unauthorized |
| 404 | 资源不存在 | Not Found |
| 500 | 服务器错误 | Internal Server Error |

### 2.2 RESTful API

**REST（表述性状态转移）**是一种软件架构风格。RESTful API遵循REST原则，使用HTTP方法表示操作，使用URL表示资源。

**URL设计规范**：

| 规范 | 错误示例 | 正确示例 |
|------|----------|----------|
| 使用名词表示资源 | /getUsers | /users |
| 使用HTTP方法 | /users/delete/1 | DELETE /users/1 |
| 使用复数 | /user/1 | /users/1 |
| 层次结构 | /users/1/orders/2 | /users/1/orders/2 |
| 使用连字符 | /userorders | /user-orders |

**API版本控制**：

```go
// URL路径版本控制
/api/v1/users
/api/v2/users

// Header版本控制
Accept: application/vnd.hubobao.v1+json
```

### 2.3 数据格式

**JSON（JavaScript对象表示法）**是API常用的数据交换格式：

```json
{
    "id": "123",
    "name": "测试项目",
    "description": "这是一个测试",
    "created_at": "2026-02-06T00:00:00Z",
    "tags": ["AI", "短剧"],
    "config": {
        "resolution": "1080p",
        "frame_rate": 24
    }
}
```

**JSON与Go结构体映射**：

```go
type Project struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    Tags        []string  `json:"tags"`
    Config      Config    `json:"config"`
}

type Config struct {
    Resolution string `json:"resolution"`
    FrameRate  int    `json:"frame_rate"`
}
```

---

## 第三章：数据库基础

### 3.1 SQL基础

**SQL（结构化查询语言）**是操作关系型数据库的标准语言。

**基本操作**：

```sql
-- 查询数据
SELECT id, name FROM projects WHERE status = 'active';

-- 插入数据
INSERT INTO projects (name, description) VALUES ('测试项目', '描述');

-- 更新数据
UPDATE projects SET status = 'completed' WHERE id = '123';

-- 删除数据
DELETE FROM projects WHERE id = '123';
```

**连接查询**：

```sql
-- 内连接
SELECT p.*, u.username 
FROM projects p
INNER JOIN users u ON p.owner_id = u.id;

-- 左连接
SELECT p.*, c.name as category_name
FROM projects p
LEFT JOIN categories c ON p.category_id = c.id;
```

### 3.2 ORM

**ORM（对象关系映射）**是将面向对象概念映射到关系数据库的技术。

**GORM基本操作**：

```go
// 创建
project := Project{Name: "测试"}
db.Create(&project)

// 查询
var project Project
db.First(&project, "id = ?", "123")

// 更新
db.Model(&project).Update("name", "新名称")

// 删除
db.Delete(&project)
```

### 3.3 数据库事务

**事务**是数据库操作的原子单元，确保数据一致性：

```go
// 开启事务
tx := db.Begin()

defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 操作
if err := tx.Create(&project).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&script).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
tx.Commit()
```

---

## 第四章：软件架构

### 4.1 分层架构

**分层架构**是最常见的软件架构模式，将系统划分为多个层次，每层有特定职责。

**三层架构**：

```
┌──────────────────┐
│   表示层（UI）   │  展示数据，接收用户输入
├──────────────────┤
│   业务逻辑层     │  处理业务规则和逻辑
├──────────────────┤
│   数据访问层     │  与数据库交互
└──────────────────┘
```

**四层架构（在三层基础上演进）**：

```
┌──────────────────┐
│     API层        │  处理外部请求
├──────────────────┤
│   应用服务层     │  用例编排
├──────────────────┤
│     领域层       │  业务规则和领域模型
├──────────────────┤
│   基础设施层     │  技术实现（数据库、外部服务）
└──────────────────┘
```

### 4.2 设计模式

**设计模式**是针对特定问题的通用解决方案。

**工厂模式**封装对象创建逻辑：

```go
type Storage interface {
    Put(key string, data []byte) error
    Get(key string) ([]byte, error)
}

type StorageType string

const (
    LocalStorage StorageType = "local"
    S3Storage    StorageType = "s3"
)

func NewStorage(t StorageType) Storage {
    switch t {
    case LocalStorage:
        return &Local{}
    case S3Storage:
        return &S3{}
    default:
        panic("unknown storage type")
    }
}
```

**策略模式**封装可互换的算法：

```go
type PricingStrategy interface {
    Calculate(price float64) float64
}

type RegularPricing struct{}

func (r *RegularPricing) Calculate(price float64) float64 {
    return price
}

type DiscountPricing struct {
    Discount float64
}

func (d *DiscountPricing) Calculate(price float64) float64 {
    return price * (1 - d.Discount)
}
```

### 4.3 依赖注入

**依赖注入（DI）**是一种解耦组件的技术，通过外部注入依赖而非内部创建。

```go
// 不使用DI
type UserService struct {
    repo *UserRepository  // 直接创建依赖
}

func NewUserService() *UserService {
    return &UserService{
        repo: NewUserRepository(),  // 紧耦合
    }
}

// 使用DI
type UserService struct {
    repo UserRepositoryInterface  // 接口依赖
}

func NewUserService(repo UserRepositoryInterface) *UserService {
    return &UserService{
        repo: repo,  // 依赖注入
    }
}
```

---

## 第五章：版本控制

### 5.1 Git基础

**Git**是最流行的分布式版本控制系统。

**基本命令**：

```bash
# 初始化仓库
git init

# 克隆仓库
git clone https://github.com/user/repo.git

# 添加文件
git add .

# 提交
git commit -m "message"

# 推送
git push origin main

# 拉取
git pull origin main

# 查看状态
git status
```

### 5.2 分支管理

**Git分支**支持并行开发：

```bash
# 创建分支
git checkout -b feature/new-feature

# 切换分支
git checkout main

# 合并分支
git merge feature/new-feature

# 删除分支
git branch -d feature/new-feature
```

### 5.3 协作工作流

**功能分支工作流**：

```
main分支 ────────────────────────────
           \            /
            \          /
             feature  ─ 开发和测试
              \
               ── 合并到main
```

---

## 练习任务

### 练习1：概念回顾（⭐）

**任务目标**：巩固核心技术概念

**任务要求**：回顾OOP的四大支柱，解释HTTP方法和状态码的含义，说明三层架构各层的职责。

**验收标准**：

| 检查项 | 完成情况 |
|--------|----------|
| 理解封装、继承、多态、抽象 | ☐ |
| 能说出至少5种HTTP方法 | ☐ |
| 能说明三层架构各层职责 | ☐ |

---

## 相关文档链接

| 文档名称 | 说明 | 难度 |
|----------|------|------|
| [架构深度解析](architecture_analysis.md) | DDD架构和设计原理 | ⭐⭐ |
| [性能优化](performance.md) | 性能分析和优化方法 | ⭐⭐⭐ |
| [设计模式](design_patterns.md) | 架构模式和最佳实践 | ⭐⭐⭐⭐ |

---

## 版本信息

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0.0 | 2026-02-06 | 初始版本，包含基础知识回顾 |
