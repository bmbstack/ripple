test:
	go test -v

define build_arch
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn
	go mod tidy
	go mod download
	go mod vendor
	CGO_ENABLED=0 GOOS=$(1) GOARCH=amd64 go build -mod=vendor -o rippleApp -v cmd/main.go
endef

build:
	$(call build_arch,linux)

build_mac:
	$(call build_arch,darwin)

build_win:
	$(call build_arch,windows)
