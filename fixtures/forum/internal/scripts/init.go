package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/internal/initial"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	ripple.InitConfigWithPath(c.String("env"), c.String("conf"))
	ripple.Default().AddLogType(ripple.LogTypeSLS)
	initial.InitLogger()
}
