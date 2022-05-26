package scripts

import (
	"fmt"
	"github.com/bmbstack/ripple/cmd/ripple/logger"
	"os"
	"strings"
)

func Generate(protoPath string) {
	logger.Logger.Info("Generate code (*.pb.go) ...")
	currentPath, _ := os.Getwd()
	logger.Logger.Info(fmt.Sprintf("auto generate file: %s/%s/*.proto", currentPath, protoPath))

	goPathArray := strings.Split(os.Getenv("GOPATH"), ":")
	goPath := goPathArray[0]
	out, err := RunCommand("bash", "-c", fmt.Sprintf("protoc -I.:%s/src --gofast_out=plugins=ripple:. ./%s/*.proto", goPath, protoPath))
	logger.Logger.Info(fmt.Sprintf("protoc gen result: %s", string(out)))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("protoc gen err: %s", err.Error()))
	}
	logger.Logger.Info("auto generate file finish")
}
