package main

import (
	"github.com/bmbstack/ripple/cmd/ripple/scripts"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ripple"
	app.Usage = "Command line tool to managing your Ripple application"
	app.Version = "0.0.1"
	app.Authors = []*cli.Author{{Name: "wangmingjob", Email: "wangmingjob@icloud.com"}}
	app.Commands = scripts.Commands()
	_ = app.Run(os.Args)
}
