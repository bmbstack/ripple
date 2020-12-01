package scripts

import (
	"errors"
	"github.com/bmbstack/ripple/cmd/ripple/logger"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Commands
func Commands() []*cli.Command {
	return []*cli.Command{
		//New application
		{
			Name:  "new",
			Usage: "Create a Ripple application",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					msg := "Please input the application name[ripple new appName]"
					logger.Logger.Error(msg)
					return errors.New(msg)
				}
				applicationName := c.Args().First()
				NewApplication(applicationName)
				return nil
			},
		},
		//Run application
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run the Ripple application",
			Action: func(c *cli.Context) error {
				if c.Args().Len() > 0 {
					msg := "Don't need input args"
					logger.Logger.Error(msg)
					return errors.New(msg)
				}
				RunApplication()
				return nil
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
