package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/internal/initial"
	"github.com/bmbstack/ripple/logger"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	ripple.InitConfigWithPath(c.String("env"), c.String("conf"))
	ripple.Default().AddLogType(ripple.LogTypeSLS)
	initial.InitLogger()
	logger.With(map[string]interface{}{
		"userId": 101,
		"traceId": "lskajdfouiaadgvv",
	}).Info("hello, tom")
	logger.Info("hello, jack")
}
