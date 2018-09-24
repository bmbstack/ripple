package logger

import (
	"os"
	"github.com/bmbstack/ripple/middleware/logger"
)

var Logger *logger.Logger

func init() {
	Logger = NewLogger()
}

func NewLogger() *logger.Logger {
	log, err := logger.NewLogger("ripple", 1, os.Stdout)
	if err != nil {
		log.Error(err.Error())
		panic(err) // Check for error
	}
	return log
}
