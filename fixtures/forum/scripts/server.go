package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/fixtures/forum/controllers"
	_ "github.com/bmbstack/ripple/fixtures/forum/controllers"
	. "github.com/bmbstack/ripple/fixtures/forum/initial"
	_ "github.com/bmbstack/ripple/fixtures/forum/models/one"
	_ "github.com/bmbstack/ripple/fixtures/forum/models/two"
)

func RunServer() {
	Logger.Info("Run server ....")
	controllers.RouteAPI()
	ripple.Default().Run()
}
