## 1. rpcx 依赖升级

- [x] 1.1 将 `github.com/smallnest/rpcx` 从 v1.7.3 升级到 v1.9.3
- [x] 1.2 将 `go.mod` 的 `go` 指令从 1.21 提升到 1.24.0
- [x] 1.3 执行 `go mod tidy` 解析连锁升级的间接依赖（quic-go/thrift/redis 等）
- [x] 1.4 更新 `README.md` 中 nacos 安装命令的 rpcx 版本号

## 2. cache 模块简化

- [x] 2.1 重写 `cache/cache.go`：`Store` 接口瘦身为 14 个薄 KV 方法，删除 `Register`/`stores` 全局 map 与 `Adapter` 字段
- [x] 2.2 重写 `cache/redis.go`：只实现薄接口方法，`Client()` 返回 `*redis.Client`
- [x] 2.3 删除 `cache/memory.go`
- [x] 2.4 清理 `config.go` 的 `CacheConfig.Adapter` 字段与 `ripple.go` 的 `Adapter` 传参
- [x] 2.5 清理配置 yaml 中的 `adapter` 字段（fixtures/forum/config/*.yaml）

## 3. protoc 模块简化

- [x] 3.1 删除 `protoc/protoc-gen-gogofast/` 目录
- [x] 3.2 删除 `protoc/protoc-gen-gogofaster/` 目录
- [x] 3.3 删除 `protoc/protoc-gen-gogoslick/` 目录

## 4. 文档清理

- [x] 4.1 更新 `README.md` 配置样本：`bindAllTag` 替换为 `validateTag`
- [x] 4.2 同步更新 `doc/architecture.md`（cache 适配器矩阵、Store 接口、protoc 目录树）

## 5. 验证

- [x] 5.1 `go build ./...` 通过
- [x] 5.2 `go test -vet=off ./...` 全部通过
- [x] 5.3 `gofmt -l` 无格式问题
