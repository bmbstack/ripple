package scripts

import (
	"fmt"
	"github.com/bmbstack/ripple/cmd/ripple/logger"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// collect files with suffix in current path
func collect(path string, suffix string) []string {
	var list []string
	err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return err
		}
		if strings.HasSuffix(file, suffix) {
			list = append(list, file)
		}
		return nil
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("filepath.Walk err: %v", err))
	}
	return list
}

func Generate(path string) {
	logger.Logger.Info("generate code (*.pb.go, *.controller.go, *.service.go) ...")
	generatePb(path)
	generateCtlService(path)
	logger.Logger.Notice("auto generate file finish")
}

func generatePb(path string) {
	goPathArray := strings.Split(os.Getenv("GOPATH"), ":")
	goPath := goPathArray[0]

	list := collect(path, ".proto")
	for _, item := range list {
		logger.Logger.Info(fmt.Sprintf("auto generate *.pb.go, ref file: %s", item))
		cmd := fmt.Sprintf("protoc -I.:%s/src --gofast_out=plugins=ripple:. %s", goPath, item)
		logger.Logger.Debug(fmt.Sprintf("Run command: %s", cmd))
		out, err := RunCommand("bash", "-c", cmd)
		logger.Logger.Info(fmt.Sprintf("protoc gen out: %s", string(out)))
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("protoc gen err: %s", err.Error()))
		}
	}
}

func generateCtlService(path string) {
	list := collect(path, ".dto.go")
	for _, item := range list {
		logger.Logger.Info(fmt.Sprintf("auto generate *.controller.go, ref file: %s", item))

		logger.Logger.Info(fmt.Sprintf("auto generate *.service.go, ref file: %s", item))

	}
}
