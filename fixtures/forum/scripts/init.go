package scripts

import (
	"github.com/bmbstack/ripple/fixtures/forum/logger"
	"os"
)

// init commands
func GetInitCommands() []string {
	commands := make([]string, 8)
	dir, err := os.Getwd()
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	commands = append(commands, "SCRIPTPATH="+dir)
	return commands
}
