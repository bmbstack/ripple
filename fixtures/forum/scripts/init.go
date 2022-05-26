package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/initial"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	ripple.InitConfig(c.String("env"), nil)
	initial.InitLogger()
}
