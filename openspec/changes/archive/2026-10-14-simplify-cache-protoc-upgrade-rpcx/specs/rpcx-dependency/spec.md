## ADDED Requirements

### Requirement: rpcx v1.9.3 依赖
框架 SHALL 依赖 `github.com/smallnest/rpcx` v1.9.3，替换原有的 v1.7.3。

#### Scenario: 依赖版本
- **WHEN** 查看 `go.mod` 依赖声明
- **THEN** `github.com/smallnest/rpcx` 版本为 v1.9.3

### Requirement: Go 1.24 构建要求
框架 SHALL 要求 Go 1.24 及以上版本构建，`go.mod` 的 `go` 指令 SHALL 为 1.24.0。

#### Scenario: 构建工具链
- **WHEN** 使用 Go 1.24+ 工具链执行 `go build ./...`
- **THEN** 构建成功，rpcx 相关 API（`server`/`client`/`util`/`log`）编译通过
