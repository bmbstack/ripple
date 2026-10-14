## Why

ripple 框架存在累积的技术债务与冗余：缓存层抽象过厚且 memory 后端从未实现、protoc 模块冗余保留了 3 个从未使用的代码生成插件、rpcx 依赖停留在 2021 年的 v1.7.3。这些冗余增加了维护成本并误导使用者，需要一次集中清理与升级。

## What Changes

- **缓存层简化（BREAKING）**：删除从未实现的 memory 缓存后端；`Store` 接口从 60+ 个 Redis 命令直通方法瘦身为 14 个通用 KV 方法；Redis 高级命令改为通过 `Client() *redis.Client` 逃生通道访问。
- **protoc 模块简化**：删除 3 个从未被使用的插件目录 `protoc-gen-gogofast`、`protoc-gen-gogofaster`、`protoc-gen-gogoslick`，仅保留实际使用的 `protoc-gen-gofast`。
- **rpcx 依赖升级（BREAKING）**：`github.com/smallnest/rpcx` 从 v1.7.3 升级到 v1.9.3，`go` 指令从 1.21 提升到 1.24.0（构建需 Go ≥ 1.24）。
- **文档与配置清理**：移除已废弃的 `bindAllTag` 配置项，统一使用 `validateTag`；同步更新 `README.md` 与 `doc/architecture.md`。

## Capabilities

### New Capabilities

- `cache`: 缓存层能力——基于 Redis 的薄 KV 抽象，提供 `Client()` 逃生通道访问高级命令。
- `protoc-codegen`: protobuf 代码生成能力——仅保留 `protoc-gen-gofast` 插件及其 ripple 扩展。
- `rpcx-dependency`: rpcx RPC 依赖与构建要求——依赖 rpcx v1.9.3，构建需 Go ≥ 1.24。

### Modified Capabilities

<!-- 无已有 spec，本次全部为新建能力 -->

## Impact

- **代码**：`cache/`（删除 `memory.go`，重写 `cache.go`、`redis.go`）、`protoc/`（删除 3 个插件目录）、`config.go`、`ripple.go`。
- **API 破坏**：`cache.Client()` 返回类型由 `interface{}` 变为 `*redis.Client`；`Store`/`Cache` 接口移除 hash/list/set/zset/bitmap/script 等高级方法。
- **依赖**：`github.com/smallnest/rpcx` v1.7.3 → v1.9.3，`go 1.21` → `go 1.24.0`。
- **配置**：`CacheConfig.Adapter` 字段移除；配置样本 `bindAllTag` 替换为 `validateTag`。
- **文档**：`README.md`、`doc/architecture.md` 同步。
