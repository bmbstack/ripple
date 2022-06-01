package initial

import (
	"github.com/bmbstack/ripple/middleware/logger"
	"os"
)

var Logger *logger.Logger

func InitLogger() {
	Logger = NewLogger()
}

func NewLogger() *logger.Logger {
	log, err := logger.NewLogger("forum", 1, os.Stdout)
	if err != nil {
		panic(err) // Check for error
	}
	return log
}
