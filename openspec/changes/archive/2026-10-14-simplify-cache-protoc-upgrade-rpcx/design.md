## Context

ripple 是轻量级 Go web + RPC 框架。本次集中清理三类技术债务：

1. **缓存层**：原 `Store` 接口是 go-redis 客户端 API 的 1:1 直通（60+ 方法），且 `MemoryCache` 从未实现（所有方法 `panic("implement me")`）。框架实际只用 `Get`/`Set`/`Delete`/`DeleteByPrefix` 四个方法。
2. **protoc 模块**：复制了 gogo/protobuf 官方 4 个 vanity 变体（gofast/gogofast/gogofaster/gogoslick），实际代码生成只调用 `--gofast_out`。
3. **rpcx 依赖**：停留在 v1.7.3（2021 年），新版本 v1.9.3 要求 Go ≥ 1.24。

## Goals / Non-Goals

**Goals:**
- 缓存层简化为"薄 KV 抽象 + Client() 逃生通道"，对齐 Spring Boot 的 `Cache` 抽象 + `StringRedisTemplate` 设计。
- protoc 模块仅保留实际使用的 gofast 插件。
- rpcx 升级到 v1.9.3，`go` 指令提升到 1.24.0。
- 清理废弃配置字段与过时文档。

**Non-Goals:**
- 不引入多后端缓存 provider 机制（memory/caffeine 等），保持 Redis 单一后端。
- 不替换 gogo/protobuf 为官方 protobuf-go（超出本次范围）。
- 不重命名 `AdapterConfig` 结构体（虽语义已变，但为降低破坏面保留）。

## Decisions

1. **缓存抽象采用薄 KV 接口**：`Store` 接口只保留 14 个通用 KV 方法，Redis 高级命令（hash/list/set/zset/bitmap/script）通过 `Client() *redis.Client` 访问。理由：参照 Spring Boot 的 `Cache`/`StringRedisTemplate` 分工，抽象只暴露通用语义，专属能力走逃生通道。备选方案（保留完整命令接口）被否——那会导致抽象过厚、未来若有非 Redis 后端无法实现。

2. **`Client()` 返回类型改为 `*redis.Client`**：由 `interface{}` 改为强类型，消除下游类型断言。代价是接口暴露了 Redis 具体类型，但既然锁定单一 Redis 后端，这是可接受的。

3. **删除 `Register`/`stores` 全局 map 与 `Adapter` 字段**：多适配器注册机制在仅有单一 Redis 实现时是虚设，`NewCache` 直接返回 `NewRedisCache()`。

4. **protoc 只保留 gofast**：删除 3 个未被引用的插件目录。三个插件是独立 `main` 包，删除无循环依赖风险。

5. **rpcx 升级与 Go 版本**：升级 rpcx 到 v1.9.3 强制 `go 1.24.0`，本地用 Go 1.26.0 工具链构建验证。

## Risks / Trade-offs

- [破坏性 API 变更：`Client()` 返回类型、`Store` 接口方法移除] → 已排查框架内与 fixtures 调用点，仅 4 个方法被使用；下游消费者若用过高级方法需迁移到 `Client()`。
- [Go 1.24 构建要求提升] → 需在 CI/文档中明确 Go 1.24+ 工具链要求。
- [gogo/protobuf 依赖保留] → `command` 与 gofast 插件仍依赖 gogo/protobuf，删除 3 个插件不影响该依赖。
