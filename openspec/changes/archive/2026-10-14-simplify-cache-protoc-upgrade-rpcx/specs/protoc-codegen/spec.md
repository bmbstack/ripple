## ADDED Requirements

### Requirement: 仅保留 protoc-gen-gofast 插件
protobuf 代码生成 SHALL 仅保留 `protoc-gen-gofast` 插件，SHALL 移除 `protoc-gen-gogofast`、`protoc-gen-gogofaster`、`protoc-gen-gogoslick` 三个未使用的插件。

#### Scenario: 生成 pb.go 代码
- **WHEN** 用户执行 `ripple gen proto` 或 `protoc --gofast_out=plugins=ripple:.`
- **THEN** 使用 `protoc-gen-gofast` 插件生成 `*.pb.go` 与 `*.rpc.go` 代码

#### Scenario: 插件目录结构
- **WHEN** 查看 `protoc/` 目录
- **THEN** 仅存在 `command/`、`ripple/`、`protoc-gen-gofast/`、`testdata/` 四个子目录
