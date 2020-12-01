package scripts

import (
	"github.com/bmbstack/ripple/cmd/ripple/templates/logger"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Commands
func Commands() []*cli.Command {
	return []*cli.Command{
		//Init application
		{
			Name:  "init",
			Usage: "Init application(go packages / DB migration /Frontend compiler)",
			Action: func(c *cli.Context) error {
				commands := GetInitCommands()
				RunScript(commands)
				logger.Logger.Info("Init application done")
				return nil
			},
		},
		//Run server
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Run server",
			Action: func(c *cli.Context) error {
				db := ""
				if c.Args().Len() > 0 {
					db = c.Args().First()
				}
				commands := GetServerCommands(db)
				RunScript(commands)
				RunServer()
				return nil
			},
		},
	}
}

// Run Script
func RunScript(commands []string) {
	entireScript := strings.NewReader(strings.Join(commands, "\n"))
	bash := exec.Command(Bash)
	stdin, _ := bash.StdinPipe()
	stdout, _ := bash.StdoutPipe()
	stderr, _ := bash.StderrPipe()

	wait := sync.WaitGroup{}
	wait.Add(3)
	go func() {
		_, _ = io.Copy(stdin, entireScript)
		_ = stdin.Close()
		wait.Done()

	}()
	go func() {
		_, _ = io.Copy(os.Stdout, stdout)
		wait.Done()

	}()
	go func() {
		_, _ = io.Copy(os.Stderr, stderr)
		wait.Done()

	}()

	_ = bash.Start()
	wait.Wait()
	_ = bash.Wait()
}
