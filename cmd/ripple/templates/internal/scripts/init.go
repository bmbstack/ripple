package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/urfave/cli/v2"
	
	"github.com/bmbstack/ripple/cmd/ripple/templates/internal/initial"
)

func Init(c *cli.Context) {
	ripple.InitConfigWithPath(c.String("env"), c.String("conf"))
	initial.InitLogger()
}
