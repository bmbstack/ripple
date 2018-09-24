package scripts

import (
	"github.com/bmbstack/ripple/cmd/ripple/logger"
	"github.com/urfave/cli"
	"os/exec"
)

// Commands
func Commands() []cli.Command {
	return []cli.Command{
		//New application
		{
			Name:  "new",
			Usage: "Create a Ripple application",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					logger.Logger.Error("Please input the application name[ripple new appName]")
					return
				}
				applicationName := c.Args()[0]
				NewApplication(applicationName)
			},
		},
		//Run application
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run the Ripple applicastion",
			Action: func(c *cli.Context) {
				if len(c.Args()) > 0 {
					logger.Logger.Error("Don't need input args")
					return
				}
				RunApplication()
			},
		},
	}
}

// runCommand runs a command with exec.Command
func RunCommand(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}
