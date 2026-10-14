# Changelog

本项目的所有重要变更都会记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [1.4.7] - 2026-10-14

本次版本聚焦框架瘦身与依赖升级：缓存层简化为纯 Redis 薄 KV 抽象、protoc 模块裁剪、rpcx 升级到 v1.9.3，并引入 OpenSpec 规范驱动开发。

### ⚠️ 破坏性变更（Breaking Changes）

- **构建要求提升**：`go.mod` 的 `go` 指令从 `1.21` 提升到 `1.24.0`，编译需 Go ≥ 1.24（因 rpcx v1.9.3 要求）。
- **cache API 变更**：
  - `Store`/`Cache` 接口移除 hash / list / set / zset / bitmap / script 等 Redis 专属方法（60+ 方法瘦身为 14 个通用 KV 方法）。
  - `Client()` 返回类型由 `interface{}` 改为 `*redis.Client`，调用高级命令需改为 `Cache().Client().HGetAll(...)` 等。
- **配置变更**：`CacheConfig.Adapter` 字段移除；`bindAllTag` 配置项移除，统一改用 `validateTag`。

### 新增（Added）

- 引入 **OpenSpec** 规范驱动开发（proposal → specs → design → tasks → archive 工作流），建立系统行为真相源 `openspec/specs/`。
- 新增内置 `Validator`（基于 `go-playground/validator`），支持自定义校验 tag（默认 `validate`）。
- 参数绑定完成后自动校验：`middleware/bind` 的 `DefaultBinder.Bind` 在绑定结束后自动调用 `Echo.Validator.Validate`。
- 新增单元测试：`protoc/ripple`、`nacos/rpcxnacos/client`、`middleware/bind`、`validator`。

### 变更（Changed）

- **rpcx 升级**：`github.com/smallnest/rpcx` `v1.7.3 → v1.9.3`，同步升级 quic-go、thrift、redis v9 等间接依赖。
- **缓存层简化**：`cache` 采用「薄 KV 抽象 + `Client()` 逃生通道」设计（对齐 Spring Boot 的 `Cache` + `StringRedisTemplate` 分工），仅保留 Redis 后端。
- **nacos-sdk-go 去 vendor**：移除 `nacos/nacos-sdk-go/v2` vendored 源码，改为直接依赖官方 `github.com/nacos-group/nacos-sdk-go/v2 v2.3.5`，import 路径同步调整。
- **参数绑定重构**：`middleware/binding` 目录删除，由 `middleware/bind`（fork echo `DefaultBinder`）取代。
- **文档同步**：更新 `README.md` 与 `doc/architecture.md`（cache 适配器矩阵、Store 接口、protoc 目录树、配置样本）。

### 移除（Removed）

- **cache memory 后端**：`cache/memory.go`（从未实现，原为 `panic("implement me")` 占位）。
- **protoc 3 个未使用插件**：`protoc-gen-gogofast`、`protoc-gen-gogofaster`、`protoc-gen-gogoslick`，仅保留实际使用的 `protoc-gen-gofast`。
- **middleware/binding 目录**：旧的 binding 实现，被 `middleware/bind` 取代。
- **vendored nacos-sdk-go 源码**：改为官方 module 依赖。

---

## [1.4.6]

- 修复 mysql dsn 的 location 处理
- 文档更新

> 更早版本（1.4.5 及之前）未维护 CHANGELOG，详见 git 提交历史。
