package main

import (
	"github.com/bmbstack/ripple"
	"github.com/bmbstack/ripple/cmd/ripple/scripts"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ripple"
	app.Usage = "Command line tool to managing your Ripple application"
	app.Version = ripple.Version()
	app.Authors = []*cli.Author{{Name: "wangmingjob", Email: "wangmingjob@icloud.com"}}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "env",
			Value:       "prod",
			DefaultText: "prod",
			Usage:       "执行环境 (开发环境dev(只有ripple作者会使用)、线上环境prod)",
		},
	}
	app.Commands = scripts.Commands()
	_ = app.Run(os.Args)
}
