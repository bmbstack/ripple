package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/cmd/ripple/templates/internal/initial"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	ripple.InitConfig(c.String("env"))
	initial.InitLogger()
}
