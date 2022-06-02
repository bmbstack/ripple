package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/cmd/ripple/templates/internal/controllers"
	"github.com/bmbstack/ripple/cmd/ripple/templates/internal/initial"
)

func RunServer() {
	initial.Logger.Info("Run server ....")
	controllers.RouteAPI()
	ripple.Default().Run()
}
