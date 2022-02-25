package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/cmd/ripple/templates/controllers"
	_ "github.com/bmbstack/ripple/cmd/ripple/templates/controllers"
	"github.com/bmbstack/ripple/cmd/ripple/templates/initial"
)

func RunServer() {
	initial.Logger.Info("Run server ....")
	controllers.RouteAPI()
	ripple.Default().Run()
}
