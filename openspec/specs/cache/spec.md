# cache Specification

## Purpose
TBD - created by archiving change simplify-cache-protoc-upgrade-rpcx. Update Purpose after archive.
## Requirements
### Requirement: 薄 KV 缓存抽象
缓存层 SHALL 仅基于 Redis 实现，`Store` 接口 SHALL 只暴露通用键值操作（`Set`/`Get`/`Del`/`Delete`/`DeleteByPrefix`/`Expire`/`TTL`/`Exists`/`Exist`/`Incr`/`Decr`/`Close`/`Key`），不暴露 Redis 专属命令。

#### Scenario: 通用键值读写
- **WHEN** 用户通过 `Cache().Set(key, value, ttl)` 写入并 `Cache().Get(key)` 读取
- **THEN** 数据在过期时间前可正确读写，键名按 `{section}:{key}` 前缀隔离

#### Scenario: 批量键删除
- **WHEN** 用户通过 `Cache().DeleteByPrefix(prefix)` 删除一批键
- **THEN** 所有匹配前缀的键被删除

### Requirement: Client 逃生通道
缓存层 SHALL 通过 `Client()` 方法返回底层 `*redis.Client`，用于访问接口未暴露的 Redis 高级命令（hash/list/set/zset/bitmap/script）。

#### Scenario: 访问 Redis 高级命令
- **WHEN** 用户通过 `Cache().Client()` 获取客户端并调用 `HGetAll` 等高级命令
- **THEN** 返回类型为 `*redis.Client`，可直接调用 Redis 命令而无需类型断言

### Requirement: 仅支持 Redis 后端
缓存层 SHALL 移除 memory 后端，`NewCache` SHALL 仅创建 Redis 缓存实例，配置中 SHALL 不再包含 `adapter` 字段。

#### Scenario: 创建缓存实例
- **WHEN** 用户通过 `cache.NewCache(Options)` 创建缓存实例
- **THEN** 创建的是 Redis 缓存，且 `Options` 无 `Adapter` 字段

