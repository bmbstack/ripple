package scripts

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/cmd/ripple/templates/controllers"
	_ "github.com/bmbstack/ripple/cmd/ripple/templates/controllers"
	"github.com/bmbstack/ripple/cmd/ripple/templates/logger"
)

// Server commands
func GetServerCommands(db string) []string {
	commands := make([]string, 2)
	switch db {
	case "mysql":
		commands = append(commands, "/usr/local/bin/mysqld_safe &")
		commands = append(commands, "sleep 10s")
	}
	return commands
}

// Run server
func RunServer() {
	logger.Logger.Info("Run server ....")
	controllers.RouteAPI()
	ripple.Run()
}
