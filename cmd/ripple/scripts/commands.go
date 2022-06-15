package scripts

import (
	"errors"
	"github.com/bmbstack/ripple/cmd/ripple/logger"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		//New application
		{
			Name: "new",
			Usage: "Create a Ripple application" +
				"\n\tdesc: ripple new appName, however this appName can be empty, will be generated in the current directory" +
				"\n\tripple new" +
				"\n\tripple new app",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					msg := "please input the application name[ripple new appName]"
					logger.Logger.Error(msg)
					return errors.New(msg)
				}
				applicationName := c.Args().First()
				NewApplication(c.String("env"), applicationName)
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
					msg := "don't need input args"
					logger.Logger.Error(msg)
					return errors.New(msg)
				}
				RunApplication()
				return nil
			},
		},
		//Generate file
		{
			Name:    "gen",
			Aliases: []string{"g"},
			Usage: "Auto generate code, *.proto => *.pb.go *.rpc.go rpc.client.go; *.dto.go => *.controller.go && *.service.go" +
				"\n\tdesc: ripple g path component name/pbPath (path: dir/file; component: ''/proto/controller/service, name: component name, pbPath: *.pb.go path)" +
				"\n\tripple g" +
				"\n\tripple g packages/app" +
				"\n\tripple g packages/app proto" +
				"\n\tripple g packages/app controller" +
				"\n\tripple g packages/app service" +
				"\n\tripple g packages/app service product" +
				"\n\tripple g packages/app ecode" +
				"\n\tripple g packages/app/proto/user.proto" +
				"\n\tripple g packages/app/internal/dto/user.dto.go" +
				"\n\tripple g packages/app2 rpc.client packages/app1/proto/user.pb.go" +
				"",
			Action: func(c *cli.Context) error {
				args := c.Args()
				path := "."
				component := "all"
				name := ""
				if !strings.EqualFold(args.Get(0), "") {
					path = args.Get(0)
				}
				if !strings.EqualFold(args.Get(1), "") {
					component = args.Get(1)
				}
				if !strings.EqualFold(args.Get(2), "") {
					name = args.Get(2)
				}
				Generate(path, component, name)
				return nil
			},
		},
	}
}

// RunCommand runs a command with exec.Command
func RunCommand(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}
