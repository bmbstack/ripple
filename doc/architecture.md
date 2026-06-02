# Ripple 框架知识库

> 版本: v1.4.6 | 基于 [Echo](https://github.com/labstack/echo) 的轻量级 Go Web 框架

---

## 目录

1. [概述](#1-概述)
2. [技术架构](#2-技术架构)
3. [项目结构](#3-项目结构)
4. [核心组件](#4-核心组件)
   - [4.1 Ripple 核心引擎](#41-ripple-核心引擎)
   - [4.2 配置系统](#42-配置系统)
   - [4.3 Controller 路由系统](#43-controller-路由系统)
   - [4.4 ORM 数据库层](#44-orm-数据库层)
   - [4.5 Cache 缓存层](#45-cache-缓存层)
   - [4.6 Logger 日志系统](#46-logger-日志系统)
   - [4.7 Renderer 模板渲染](#47-renderer-模板渲染)
   - [4.8 Middleware 中间件](#48-middleware-中间件)
   - [4.9 RPC 远程调用](#49-rpc-远程调用)
   - [4.10 CLI 代码生成工具](#410-cli-代码生成工具)
5. [技术特点](#5-技术特点)
6. [应用项目结构](#6-应用项目结构)
7. [快速开始](#7-快速开始)
8. [配置详解](#8-配置详解)
9. [依赖拓扑](#9-依赖拓扑)

---

## 1. 概述

Ripple 是一个基于 [Echo v4](https://github.com/labstack/echo) 构建的轻量级 Go Web 框架。它在 Echo 基础上集成了 ORM（GORM）、缓存（Redis/Memory）、日志（Logrus + 阿里云 SLS / 腾讯云 CLS）、RPC 通信（Rpcx + Nacos）和代码生成 CLI 工具，提供"开箱即用"的全栈 Web 开发体验。

**核心设计理念：**

- **Simplicity（简洁）** — 不过度封装，用户可在一天内理解整个框架
- **Relevance（适度）** — 只做必要的事，不引入过多抽象层
- **Elegance（优雅）** — 遵循 Go 最佳实践，源码注释完备

**框架三要素：**

```
请求 → Router（路由） → Middleware（中间件） → Handler（处理器） → 响应
```

---

## 2. 技术架构

Ripple 采用分层架构设计，各层职责清晰，通过全局单例 `Ripple` 实例协调工作：

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI 工具层                            │
│   ripple new  │  ripple run  │  ripple gen (代码生成)        │
├─────────────────────────────────────────────────────────────┤
│                     Controller 路由层                        │
│   controller tag 反射解析 → Echo Group → RESTful 路由注册     │
├─────────────────────────────────────────────────────────────┤
│                      Middleware 中间件层                      │
│   Recover  │  Logger  │  Bind/Binding  │  Pprof              │
├─────────────────────────────────────────────────────────────┤
│                      Service 业务逻辑层                       │
│   自动生成的 Service 接口 + 手动实现                           │
├──────────┬──────────┬──────────┬──────────┬─────────────────┤
│   ORM    │  Cache   │  Logger  │ Renderer │      RPC        │
│  (GORM)  │ (Redis/  │(Logrus + │(Pongo2)  │ (Rpcx + Nacos)  │
│          │ Memory)  │ SLS/CLS) │          │                 │
├──────────┴──────────┴──────────┴──────────┴─────────────────┤
│                      Infrastructure 基础设施                  │
│   MySQL/PostgreSQL  │  Redis  │  Nacos  │  SLS/CLS 日志云    │
└─────────────────────────────────────────────────────────────┘
```

**Ripple 核心结构体：**

```go
type Ripple struct {
    Logger              *logger.Logger        // 内置彩色的控制台 Logger
    Logrus              *logrus.Logger        // 结构化 Logger (JSON/云端)
    Echo                *echo.Echo            // Echo HTTP 引擎
    Orms                map[string]*Orm       // 多数据库连接池 (alias → Orm)
    Caches              map[string]*Cache     // 多缓存实例 (alias → Cache)
    Logs                map[string]*logrus.Logger // 多日志实例 (alias → Logger)
    RpcServer           *server.Server        // RPC 服务端
    NacosRegisterPlugin *serverplugin.NacosRegisterPlugin // Nacos 注册插件
}
```

`Ripple` 实例通过 `sync.Once` 实现全局单例模式，通过 `ripple.Default()` 获取。

---

## 3. 项目结构

Ripple 框架本身的源码目录：

```
ripple/
├── cache/                         # 缓存抽象层
│   ├── cache.go                   #   Store 接口 + Cache 门面
│   ├── redis.go                   #   Redis 适配器 (go-redis/v8)
│   └── memory.go                  #   内存适配器 (未完整实现)
├── cmd/ripple/                    # CLI 脚手架工具
│   ├── main.go                    #   入口 (urfave/cli)
│   ├── logger/                    #   CLI 内部 logger
│   ├── scripts/
│   │   ├── commands.go            #   命令定义 (new/run/gen)
│   │   ├── gen.go                 #   代码生成核心逻辑 (956行)
│   │   ├── new.go                 #   创建项目脚手架
│   │   └── run.go                 #   运行项目
│   ├── templates/                 #   脚手架模板
│   └── util/                      #   CLI 工具函数
├── config.go                      # 配置结构体 + Viper 加载
├── controller.go                  # Controller 接口 + 反射路由注册
├── helper/helper.go               # 通用工具函数
├── logger/                        # 日志系统
│   ├── logger.go                  #   基于 Logrus 的全局 Logger
│   ├── cls/hook.go                #   腾讯云 CLS Hook
│   └── sls/hook.go                #   阿里云 SLS Hook
├── middleware/                    # 中间件
│   ├── bind/bind.go               #   全 Tag 绑定器 (form/json/query/param/header)
│   ├── binding/                   #   Echo 默认绑定器包装
│   ├── logger/logger.go           #   Worker 级别彩色 Logger
│   ├── pprof/pprof.go             #   pprof 性能分析路由
│   └── recover/recover.go         #   Panic 恢复中间件
├── nacos/                         # Nacos 服务发现集成
│   ├── nacos-sdk-go/              #   内嵌 Nacos Go SDK
│   └── rpcxnacos/
│       ├── client/                #   RPC Nacos 客户端发现
│       └── serverplugin/          #   RPC Nacos 服务注册插件
├── orm.go                         # GORM 封装
├── protoc/                        # Protobuf 代码生成插件
│   ├── command/                   #   protoc 命令封装
│   ├── ripple/ripple.go           #   protoc-gen-gofast 的 ripple 插件 (生成 *.rpc.go)
│   └── protoc-gen-gofast/         #   基于 gogo/protobuf 的 protoc 插件
├── render.go                      # 模板渲染 (Pongo2)
├── ripple.go                      # Ripple 核心引擎 (单例 + 生命周期管理)
├── rpc.go                         # RPC 服务端创建 (Rpcx + Nacos)
├── rst/                           # Go AST 工具库 (解析/修改源码)
│   ├── rst.go, file.go, scope.go  #   AST 解析核心
│   ├── func_decl.go, func_body.go #   函数声明/体操作
│   ├── func_call.go, func_lit.go  #   函数调用/字面量操作
│   └── comment.go, util.go        #   注释 + 工具
├── sync/once.go                   # 同步工具
├── util/ip.go                     # 内网 IP 获取
├── fixtures/forum/                # 示例项目 (论坛)
└── vendor/                        # Go vendor 依赖
```

---

## 4. 核心组件

### 4.1 Ripple 核心引擎

**文件：** `ripple.go`

Ripple 核心负责整个框架生命周期的管理：

| 职责 | 方法 | 说明 |
|------|------|------|
| 实例化 | `NewRipple()` / `Default()` | 创建全局单例，初始化 Echo、ORM、Cache、Logger |
| 控制器注册 | `RegisterController(c)` | 通过反射解析 struct tag 注册路由 |
| 模型注册 | `RegisterModels(orm, models...)` | 注册 GORM 模型，支持 AutoMigrate |
| RPC 注册 | `RegisterRpc(name, rpc, metadata)` | 注册 RPC 服务并上报 Nacos |
| RPC 运行 | `RunRpc()` | 启动 RPC 服务端监听 |
| HTTP 运行 | `Run()` / `RunWith(domain)` | 启动 HTTP 服务 |
| 日志类型 | `AddLogType(LogTypeSLS/CLS)` | 添加云端日志 Hook |
| 资源释放 | `CloseOrm()` / `CloseCache()` / `StopRpc()` | 优雅关闭 |

**生命周期流程：**

```
InitConfig(env)
  └→ viper 加载 config.{env}.yaml
      └→ Default() / NewRipple()
          ├→ NewLogger()              创建内置 Logger
          ├→ echo.New()               创建 Echo 实例
          ├→ 注册 Recover + Logger 中间件
          ├→ 设置 Binder (bind/binding)
          ├→ 设置 Renderer (Pongo2)
          ├→ 注册 Static 静态文件目录
          ├→ 遍历 Databases → NewOrm()  创建数据库连接
          ├→ 遍历 Caches → NewCache()   创建缓存连接
          └→ 遍历 Logs → NewLog()       创建日志实例
      └→ RegisterController()         注册业务路由
      └→ RegisterModels()             注册数据模型
      └→ RegisterRpc() + RunRpc()     启动 RPC 服务
      └→ Run()                        启动 HTTP 监听
      └→ OnStop(clean)                监听退出信号，优雅关闭
```

### 4.2 配置系统

**文件：** `config.go`

基于 [Viper](https://github.com/spf13/viper) 实现多环境 YAML 配置加载。

**配置结构：**

```go
type BaseConfig struct {
    AutoMigrate bool             // 启动时自动迁移表结构
    BindAllTag  bool             // true=bind中间件, false=echo默认binding
    Debug       bool             // 调试模式
    Domain      string           // 监听地址 (如 127.0.0.1:8090)
    Static      string           // 静态文件目录
    Templates   string           // 模板文件目录
    Databases   []DatabaseConfig // 数据库配置列表
    Caches      []CacheConfig    // 缓存配置列表
    Logs        []LogConfig      // 日志配置列表
    Nacos       NacosConfig      // Nacos 服务发现配置
    SLS         SlsConfig        // 阿里云日志服务配置
    CLS         ClsConfig        // 腾讯云日志服务配置
}
```

**按环境加载：**

- `config.dev.yaml` — 开发环境
- `config.test.yaml` — 测试环境
- `config.prod.yaml` — 生产环境

运行时通过 `--env` 参数指定：`go run cmd/main.go --env dev`

### 4.3 Controller 路由系统

**文件：** `controller.go`

Ripple 的 Controller 系统通过 **struct tag + 反射** 实现声明式路由注册，避免手写路由配置。

**核心接口：**

```go
type Controller interface {
    Path() string  // 路由组前缀
}
```

**Tag 格式：**

在 struct 字段上使用 `controller` tag，格式为 `"HTTP方法 路由路径"`：

```go
type HomeController struct {
    Index  echo.HandlerFunc `controller:"GET /"`        // → GET /
    Html   echo.HandlerFunc `controller:"GET html"`     // → GET /html
    Create echo.HandlerFunc `controller:"POST /create"` // → POST /create
    // 中间件字段
    Auth   echo.MiddlewareFunc `controller:",middleware"` // → Group.Use(Auth)
}
```

**工作流程：**

```
RegisterController(c)
  └→ reflectCtrl(c) → 获取 struct 类型信息
  └→ echoMux.Group(c.Path()) → 创建路由组
  └→ 遍历 struct 字段:
      └→ parseTag(tag) → 解析 "GET /path"
      └→ getResourceFunc() → 查找 Action{Name} 方法
      └→ grp.METHOD(path, handler) → 注册到 Echo Group
```

**方法命名约定：** 字段 `Index` 对应方法 `ActionIndex`。

**Namespace 类型：**

```go
type Namespace string  // 实现 Controller 接口
func (n Namespace) Path() string { return string(n) }
```

支持 HTTP 方法：GET / POST / PUT / PATCH / DELETE / HEAD / OPTIONS / CONNECT / TRACE。

### 4.4 ORM 数据库层

**文件：** `orm.go`

基于 [GORM v2](https://gorm.io) 封装的数据库操作层，支持 MySQL 和 PostgreSQL。

**核心结构：**

```go
type Orm struct {
    models map[string]reflect.Value  // 已注册模型
    *gorm.DB                          // 嵌入 GORM DB
}
```

**数据库配置：**

```yaml
databases:
  - alias: "forum"          # 别名，通过 GetOrm("forum") 获取
    dialect: "mysql"        # mysql / postgresql
    host: "127.0.0.1"
    port: 3306
    name: "forum"
    username: "root"
    password: "123456"
    maxIdleConns: 200       # 最大空闲连接
    maxOpenConns: 1000      # 最大打开连接
```

**关键方法：**

| 方法 | 说明 |
|------|------|
| `NewOrm(config, debug)` | 创建数据库连接，支持 mysql/postgresql |
| `AddModels(values...)` | 注册数据模型 |
| `AutoMigrateAll()` | 自动迁移所有注册的模型 |
| `Close()` | 关闭数据库连接 |
| `ripple.GetOrm("alias")` | 获取指定别名的 ORM 实例 |

**多数据库支持：** Ripple 支持同时连接多个数据库（不同 alias），通过 `ripple.GetOrm("alias")` 获取指定实例。

**模型注册示例：**

```go
type User struct { gorm.Model /* ... */ }

func init() {
    orm := ripple.Default().GetOrm("forum")
    ripple.Default().RegisterModels(orm, &User{})
}
```

### 4.5 Cache 缓存层

**文件：** `cache/cache.go`, `cache/redis.go`, `cache/memory.go`

提供统一的缓存抽象，通过 **Store 接口** 解耦后端实现。

**Store 接口（精简）：**

```go
type Store interface {
    connect(opt Options) error
    Client() interface{}
    Close() error
    Key(key string) string
    // String
    Set(key, val string, expiration time.Duration) error
    Get(key string) (string, error)
    // Hash
    HGet/HGetAll/HSet/HExists/HMGet/HMSet/HSetNX
    // Set
    SCard/SAdd/SRem/SScan
    // List
    LRange/LPush/LLen/RPop/RPopCount/RPopLPush/RPush/RPushX
    // Sorted Set
    ZAdd/ZRange/ZRank/ZScore/ZRem
    // Advanced
    SetEX/SetNX/RunScript/Incr/Decr/Expire/Del/Flush
    // Bitmap
    GetBit/SetBit/BitCount/BitOpAnd/BitOpOr/BitOpXor/BitOpNot/BitPos/BitField
}
```

**适配器矩阵：**

| 适配器 | 后端 | 实现状态 | Key 前缀 |
|--------|------|----------|----------|
| `redis` | go-redis/v8 | ✅ 完整 | `{section}:{key}` |
| `memory` | sync.Map | ⚠️ 骨架 | 无 |

**缓存配置：**

```yaml
caches:
  - alias: "forum"
    section: "forum"        # Key 前缀
    adapter: "redis"        # redis / memory
    host: "127.0.0.1"
    port: 6379
    password: "123456"
    db: 0
    poolFIFO: true          # 连接池 FIFO
    poolSize: 10
    minIdleConns: 5
```

**使用方式：**

```go
cache := ripple.Default().GetCache("forum")
// 方式1: 使用 Store 接口
cache.Set("key", "value", 0)
cache.HGetAll("hash_key")
// 方式2: 获取底层 Redis Client
client := cache.Client().(*redis.Client)
client.HGetAll(ctx, "hash_key").Result()
```

### 4.6 Logger 日志系统

Ripple 有两套日志系统，各有分工：

#### 4.6.1 内置 Logger（控制台彩色输出）

**文件：** `middleware/logger/logger.go`

基于 `log` 标准库的彩色终端日志，支持级别：CRITICAL / ERROR / WARNING / NOTICE / DEBUG / INFO。

> 这是 Ripple 核心引擎 `Ripple.Logger` 的类型，在 `NewRipple()` 时自动创建。

#### 4.6.2 Structured Logger（Logrus + 云端）

**文件：** `logger/logger.go`

基于 [Logrus](https://github.com/sirupsen/logrus) 的结构化日志，支持 JSON 格式输出，可通过 Hook 上报云端日志服务。

**日志级别：** Trace / Debug / Info / Warn / Error / Fatal / Panic

**云端日志 Hook：**

| Hook | 云服务 | 文件 |
|------|--------|------|
| SLS Hook | 阿里云日志服务 (Simple Log Service) | `logger/sls/hook.go` |
| CLS Hook | 腾讯云日志服务 (Cloud Log Service) | `logger/cls/hook.go` |

**日志配置（YAML）：**

```yaml
logs:
  - alias: "box"
    type: "sls"                     # sls / cls
    lumberjackConfig:               # 本地文件滚动 (可选)
      filename: "./log/box.log"
      maxSize: 10                   # MB
      maxBackups: 3
      maxAge: 28                    # days
      compress: true
    accessKeyId: "xxx"
    accessKeySecret: "xxx"
    endpoint: "cn-beijing.log.aliyuncs.com"
    allowLogLevel: "info"
    closeStdout: false
    project: "my-project"
    logstore: "my-logstore"
    topic: "topic"
    source: "source"
```

**使用方式：**

```go
// 全局 Logger
ripple.Default().AddLogType(ripple.LogTypeSLS)  // 添加云端 Hook

// 结构化日志
logger.With(map[string]interface{}{
    "userId":  101,
    "traceId": "xxx",
}).Info("hello, tom")

logger.Info("hello, jack")

// 多实例日志
boxLog := ripple.Default().GetLog("box")
logger.WithTo(boxLog, nil).Info("box log message")
```

**输出通道策略：**

```
日志条目
├→ stdout (closeStdout=false)
├→ Lumberjack 本地滚动文件 (lumberjackConfig 配置时)
└→ 云端 Hook (SLS / CLS)
```

### 4.7 Renderer 模板渲染

**文件：** `render.go`

基于 [Pongo2](https://github.com/flosch/pongo2)（Django 模板语法的 Go 实现）的模板渲染引擎。

**特性：**

- Django 风格模板语法（`{{ variable }}`, `{% if %}`, `{% for %}` 等）
- 支持模板缓存（Reload 模式可选）
- 自动从配置的 `templates` 目录加载模板文件

**Pongo2 模板示例：**

```html
<!DOCTYPE html>
<html>
<body>
    <h3>{{ title }}</h3>
    {% if users %}
        {% for user in users %}
            <li>{{ user.Name }}</li>
        {% endfor %}
    {% endif %}
</body>
</html>
```

**控制器中使用：**

```go
func (this HomeController) ActionIndex(ctx echo.Context) error {
    return ctx.Render(http.StatusOK, "home/index.html", map[string]interface{}{
        "title": "Hello, Ripple",
    })
}
```

### 4.8 Middleware 中间件

#### 4.8.1 Recover — Panic 恢复

**文件：** `middleware/recover/recover.go`

自动捕获 Handler 链中的 panic，记录错误日志和堆栈，防止服务崩溃。

**默认配置：** 堆栈大小 4KB，打印完整堆栈。

#### 4.8.2 Bind — 全 Tag 绑定器

**文件：** `middleware/bind/bind.go`

对标 Echo 默认 Binder 的增强版，支持同时从 query/param/form/header/json/xml 绑定参数。

**绑定顺序：** Path Params → Query Params → Request Body → Validation

**支持的 Content-Type：**

- `application/json`
- `application/xml` / `text/xml`
- `application/x-www-form-urlencoded`
- `multipart/form-data`

**配置开关：** `bindAllTag: true` 使用 Bind，`false` 使用 Echo 默认 Binding。

#### 4.8.3 Binding — Echo 默认绑定器包装

**文件：** `middleware/binding/`

对 Echo 默认 binder 的轻量包装，支持 `validator/v10` 数据校验，通过 `binding:"required"` tag 声明。

#### 4.8.4 Pprof — 性能分析

**文件：** `middleware/pprof/pprof.go`

一键注册 Go 标准 pprof 路由：

```
GET /debug/pprof/              → Index
GET /debug/pprof/heap          → 堆内存分析
GET /debug/pprof/goroutine     → Goroutine 分析
GET /debug/pprof/block         → 阻塞分析
GET /debug/pprof/threadcreate  → 线程创建分析
GET /debug/pprof/cmdline       → 命令行参数
GET /debug/pprof/profile       → CPU Profile
GET /debug/pprof/symbol        → 符号表
```

使用：`ripplepprof.Wrap(echoInstance)`

#### 4.8.5 Middleware Logger

**文件：** `middleware/logger/logger.go`

Echo 中间件级别的请求日志记录（使用 Echo 内置的 `middleware.Logger()`）。

### 4.9 RPC 远程调用

**文件：** `rpc.go`, `nacos/rpcxnacos/`

Ripple 集成了 [Rpcx](https://github.com/smallnest/rpcx) 作为 RPC 框架，[Nacos](https://nacos.io) 作为服务注册与发现中心。

**RPC 架构：**

```
┌─────────────────┐         ┌─────────────┐         ┌─────────────────┐
│  Service A       │ register│   Nacos     │discover │  Service B       │
│  (RPC Server)    │────────▶│  Registry   │◀────────│  (RPC Client)    │
│  Rpcx + Nacos    │         │             │         │  XClientPool     │
└─────────────────┘         └─────────────┘         └─────────────────┘
```

**服务端（Provider）：**

```go
// 1. 注册 RPC 服务
ripple.Default().RegisterRpc("StudentRpc", &StudentRpcImpl{}, "")
// 2. 启动 RPC 监听
ripple.Default().RunRpc()

// 实现
type StudentRpcImpl struct{}
func (s *StudentRpcImpl) Learn(ctx context.Context, req *proto.LearnReq, reply *proto.LearnReply) error {
    reply.Name = "learn result"
    return nil
}
```

**客户端（Consumer）：**

```go
sc := services.GetStudentClient()  // 自动生成的 Client 工厂
reply, _ := sc.Learn(context.Background(), &proto.LearnReq{Id: 1})
```

**Nacos 配置：**

```yaml
nacos:
  host: "my.nacos.com"     # Nacos 服务地址
  port: 8848
  namespaceId: "public"
  cluster: "ripple_user"    # 集群名
  group: "DEFAULT_GROUP"
  server: "127.0.0.1:18090" # 本机 RPC 监听地址
  failMode: "failover"     # failover / failfast / failbackup
  selectMode: "roundRobin" # randomSelect / roundRobin / weightedRoundRobin
  clientPoolSize: 10
```

**RPC 生命周期：**

```
NewRpcServerNacos(config)
  └→ 创建 Rpcx Server
  └→ 创建 NacosRegisterPlugin
      └→ plugin.Start() → 连接 Nacos
      └→ s.Plugins.Add(plugin) → 注册插件
RegisterRpc(name, rpc, metadata)
  └→ plugin.Register() → 向 Nacos 注册服务实例
RunRpc()
  └→ s.Serve("tcp", address) → goroutine 启动监听
StopRpc()
  └→ plugin.Stop() → 注销 Nacos 所有服务
  └→ s.Close() → 关闭所有连接
```

### 4.10 CLI 代码生成工具

**文件：** `cmd/ripple/`

Ripple CLI 是一个命令行工具，用于项目创建、运行和代码生成。

**三大命令：**

| 命令 | 别名 | 功能 |
|------|------|------|
| `ripple new` | — | 创建新项目脚手架 |
| `ripple run` | `r` | 运行项目 |
| `ripple gen` | `g` | 自动生成代码 |

#### 4.10.1 代码生成 (`ripple gen`)

这是 Ripple 最核心的 CLI 功能，通过解析 `.proto` 和 `.dto.go` 文件自动生成大量代码。

**生成映射表：**

| 输入文件 | 输出文件 | 输出目录 | 生成方式 |
|----------|----------|----------|----------|
| `*.proto` | `*.pb.go` | proto/ | protoc + protoc-gen-gofast |
| `*.proto` (含 `@RippleRpc`) | `*.rpc.go` | internal/rpc/ | protoc ripple 插件 (AST 生成) |
| `*.dto.go` (含 `@RippleApi`) | `*.controller.go` | internal/controllers/v1/ | AST 解析 + 模板生成 |
| `*.dto.go` (含 `@RippleApi`) | `*.service.go` | internal/services/ | AST 解析 + 模板生成 |
| `*.pb.go` (外部引用) | `rpc.client.go` | internal/rpcclient/ | AST 解析 + 模板生成 |
| — | `ecode.go` | internal/ecode/ | 固定模板生成 |

**生成方式详解：**

1. **Proto → pb.go：** 调用 `protoc` 命令，使用 `protoc-gen-gofast` 插件（基于 gogo/protobuf）
2. **Proto → rpc.go：** protoc 的 ripple 插件解析 proto 文件中的 `@RippleRpc` / `@NacosGroup` / `@NacosCluster` 注释，通过 Go AST 工具（`rst/` 包）生成 RPC 服务端/客户端代码
3. **DTO → controller/service：** 使用 `rst/` 包解析 Go 源码 AST，提取 `@RippleApi` / `@Uri` / `@Method` 注释，自动生成 Controller 和 Service 代码

**示例注释标记：**

```protobuf
// @RippleRpc
// @NacosGroup DEFAULT_GROUP
// @NacosCluster ripple
service Student { ... }
```

```go
// @RippleApi
// @Uri /student/learn
// @Method POST
type ReqStudentLearn struct { ... }
```

---

## 5. 技术特点

### 5.1 技术栈一览

| 层面 | 技术选型 | 说明 |
|------|----------|------|
| HTTP 框架 | Echo v4.6.1 | 高性能、极简的 Go HTTP 框架 |
| ORM | GORM v1.22 | 全功能 ORM，支持 MySQL / PostgreSQL |
| 缓存 | go-redis/v8 | 生产级 Redis 客户端 |
| 结构化日志 | Logrus v1.6 | 结构化日志库 |
| 模板引擎 | Pongo2 v4 | Django 模板语法的 Go 实现 |
| RPC 框架 | Rpcx v1.7.3 | 高性能 Go RPC 框架 |
| 服务注册 | Nacos | 阿里巴巴开源的服务发现与配置中心 |
| 云端日志 | 阿里云 SLS / 腾讯云 CLS | 生产级日志采集 |
| 日志滚动 | lumberjack v2 | 本地日志文件滚动管理 |
| 配置管理 | Viper v1.11 | 多格式配置文件加载 |
| 数据校验 | validator/v10 | struct tag 声明的数据校验 |
| 定时任务 | gocron v1.16 | 定时任务调度 |
| Protobuf | gogo/protobuf | 高性能 protobuf 实现 |
| 代码生成 | Go AST (dst 库) | 声明式代码生成 |
| CLI | urfave/cli v2 | 命令行框架 |

### 5.2 框架特性

- **声明式路由** — struct tag 驱动，告别手写路由
- **多数据库支持** — 同一应用连接多个数据库实例
- **多缓存实例** — 支持多个 Redis 实例，Key prefix 自动隔离
- **双日志系统** — 开发用彩色控制台，生产用 JSON + 云端 Hook
- **全 Tag 参数绑定** — 同时从 query/param/form/header/json/xml 绑定
- **代码自动生成** — 从 proto 和 dto 自动生成 Controller、Service、RPC 桩代码
- **服务注册发现** — 内置 Nacos 集成，RPC 调用开箱即用
- **优雅关闭** — 信号拦截，自动清理 ORM/Cache/RPC 资源

---

## 6. 应用项目结构

Ripple 脚手架生成的项目结构：

```
rippleApp/
├── Makefile                        # 构建脚本
├── cmd/
│   └── main.go                     # 程序入口
├── config/
│   ├── config.dev.yaml             # 开发环境配置
│   ├── config.test.yaml            # 测试环境配置
│   └── config.prod.yaml            # 生产环境配置
├── frontend/
│   ├── static/                     # 静态资源 (CSS/JS/图片)
│   └── templates/                  # Pongo2 模板文件
├── go.mod
├── go.sum
├── internal/
│   ├── controllers/                # 控制器层
│   │   └── v1/                     #   API v1 (ripple gen 固定输出)
│   ├── dto/                        # 数据传输对象 (ripple gen 输入源)
│   ├── ecode/                      # 错误码定义 (ripple gen 固定输出)
│   ├── helper/                     # 业务辅助函数
│   ├── initial/                    # 初始化逻辑
│   ├── models/                     # GORM 数据模型
│   ├── rpc/                        # RPC 服务实现 (ripple gen 固定输出)
│   ├── rpcclient/                  # RPC 客户端 (ripple gen 固定输出)
│   ├── scripts/                    # 业务脚本
│   └── services/                   # 服务层 (ripple gen 固定输出)
└── proto/
    ├── user.proto                  # Protobuf 定义 (ripple gen 输入源)
    └── user.pb.go                  # 生成的 protobuf 代码
```

**分层职责：**

| 层 | 目录 | 职责 |
|----|------|------|
| 入口 | `cmd/main.go` | 初始化配置、注册组件、启动服务 |
| 控制器 | `internal/controllers/` | 处理 HTTP 请求、参数校验、调用 Service |
| 服务 | `internal/services/` | 业务逻辑实现 |
| 数据模型 | `internal/models/` | GORM 模型定义 |
| 数据传输 | `internal/dto/` | 请求/响应结构体定义 |
| RPC 服务 | `internal/rpc/` | RPC 接口实现 |
| RPC 客户端 | `internal/rpcclient/` | 远程服务调用客户端 |
| Proto | `proto/` | Protobuf 接口定义 |

---

## 7. 快速开始

### 7.1 安装

```shell
go get github.com/bmbstack/ripple@latest
go install github.com/bmbstack/ripple/cmd/ripple@latest
go install github.com/bmbstack/ripple/protoc/protoc-gen-gofast@latest
```

### 7.2 创建项目

```shell
ripple new rippleApp
cd $GOPATH/src/rippleApp
go mod init
go mod tidy
go mod vendor
go run cmd/main.go --env dev
```

### 7.3 创建模型

在 `internal/models/user.go`：

```go
package models

import "github.com/bmbstack/ripple"

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"size:255;not null"`
    Email    string `gorm:"size:255"`
}

func init() {
    ripple.Default().RegisterModels(
        ripple.Default().GetOrm("forum"),
        &User{},
    )
}
```

### 7.4 创建控制器

在 `internal/controllers/v1/user.controller.go`：

```go
package v1

import (
    "net/http"
    "github.com/bmbstack/ripple"
    "github.com/labstack/echo/v4"
)

type UserController struct {
    List echo.HandlerFunc `controller:"GET /users"`
}

func (c UserController) Path() string { return "/api/v1" }

func (c UserController) ActionList(ctx echo.Context) error {
    return ctx.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}

func init() {
    ripple.Default().RegisterController(&UserController{})
}
```

### 7.5 使用代码生成

定义 DTO 文件 `internal/dto/user.dto.go`：

```go
// @RippleApi
// @Uri /user/create
// @Method POST
type ReqUserCreate struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}
```

执行生成：

```shell
ripple g
```

自动生成 `internal/controllers/v1/user.controller.go` 和 `internal/services/user.service.go`。

### 7.6 优雅关闭

```go
func main() {
    // ... 初始化 ...
    ripple.helper.OnStop(func() {
        ripple.Default().StopRpc()
        ripple.Default().CloseOrm()
        ripple.Default().CloseCache()
    })
}
```

---

## 8. 配置详解

### 8.1 完整 YAML 配置参考

```yaml
# === 基础配置 ===
domain: "127.0.0.1:8090"          # HTTP 监听地址
static: "frontend/static"          # 静态文件目录
templates: "frontend/templates"    # 模板文件目录
autoMigrate: false                 # 启动时自动迁移表结构
bindAllTag: false                  # false=Echo默认绑定, true=全tag绑定
debug: false                       # 调试模式(GORM日志/echo debug)

# === 数据库配置 ===
databases:
  - alias: "default"               # 唯一别名
    dialect: "mysql"               # mysql | postgresql
    host: "127.0.0.1"
    port: 3306
    name: "database_name"
    username: "root"
    password: "password"
    maxIdleConns: 200
    maxOpenConns: 1000
    charset: "utf8mb4"            # 默认 utf8mb4
    parseTime: "True"             # 默认 True
    loc: "Local"                  # 默认 Local

# === 缓存配置 ===
caches:
  - alias: "default"
    section: "app"                 # Key 前缀隔离
    adapter: "redis"               # redis | memory
    host: "127.0.0.1"
    port: 6379
    password: ""
    db: 0
    poolFIFO: true
    poolSize: 10
    minIdleConns: 5
    # url: "redis://user:pass@host:6379/0"  # URL 模式(优先级高)

# === Nacos 服务发现 ===
nacos:
  host: "nacos-server.com"
  port: 8848
  namespaceId: "public"
  cluster: "my-cluster"
  group: "DEFAULT_GROUP"
  server: "127.0.0.1:18090"       # 本机 RPC 监听地址
  failMode: "failover"            # failover | failfast | failbackup
  selectMode: "roundRobin"        # randomSelect | roundRobin | weightedRoundRobin
  clientPoolSize: 10
  cacheDir: "./cache"
  logDir: "./log"

# === 日志配置 (多实例) ===
logs:
  - alias: "biz"
    type: "sls"                   # sls | cls
    accessKeyId: "xxx"
    accessKeySecret: "xxx"
    endpoint: "cn-beijing.log.aliyuncs.com"
    allowLogLevel: "info"
    closeStdout: false
    project: "my-project"
    logstore: "my-logstore"
    topic: "topic"
    source: "source"
    lumberjackConfig:             # 可选：本地文件滚动
      filename: "./log/app.log"
      maxSize: 10                 # MB
      maxBackups: 3
      maxAge: 28                  # days
      compress: true

# === 全局 SLS 配置 (直接附加到主 Logrus) ===
sls:
  accessKeyId: "xxx"
  accessKeySecret: "xxx"
  endpoint: "cn-beijing-intranet.log.aliyuncs.com"
  allowLogLevel: "info"
  closeStdout: false
  project: "my-project"
  logstore: "my-logstore"
  topic: "topic"
  source: "source"

# === 全局 CLS 配置 (直接附加到主 Logrus) ===
cls:
  accessKeyId: "xxx"
  accessKeySecret: "xxx"
  endpoint: "ap-beijing.cls.tencentcs.com"
  allowLogLevel: "info"
  closeStdout: false
  topic: "topic-id"
```

### 8.2 环境切换

通过 `--env` 参数选择配置文件：

```shell
go run cmd/main.go --env dev    # 加载 config/config.dev.yaml
go run cmd/main.go --env prod   # 加载 config/config.prod.yaml
go run cmd/main.go --env test   # 加载 config/config.test.yaml
```

---

## 9. 依赖拓扑

```
ripple
├── Echo v4              (HTTP 框架)
│   └── gommon            (color/bytes 工具)
├── GORM v2              (ORM)
│   ├── mysql driver
│   └── postgres driver
├── go-redis/v8          (Redis 客户端)
├── Logrus v1.6          (结构化日志)
├── Pongo2 v4            (模板引擎)
├── Rpcx v1.7.3          (RPC 框架)
│   └── Nacos SDK v2      (服务发现)
├── Viper v1.11          (配置管理)
├── go-co-op/gocron      (定时任务)
├── gogo/protobuf        (protobuf 实现)
├── dave/dst             (Go AST 装饰器)
├── urfave/cli v2        (CLI 框架)
├── validator/v10        (参数校验)
├── lumberjack v2        (日志滚动)
├── AliCloud SDK
│   ├── sls (log)
│   └── kms (key management)
└── TencentCloud CLS SDK (日志)
```

---

> **文档版本：** v1.0 | **最后更新：** 2026-06-02 | **框架版本：** Ripple v1.4.6
