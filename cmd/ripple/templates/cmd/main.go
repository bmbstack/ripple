package main

import (
	"fmt"
	. "github.com/bmbstack/ripple/helper"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli/v2"
	"os"
	"time"

	"github.com/bmbstack/ripple/cmd/ripple/templates/internal/scripts"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "server",
			Aliases:     []string{"s"},
			Usage:       "Run server",
			Description: "Run web server",
			Action: func(c *cli.Context) error {
				OnStop(func() {
					fmt.Println("OnStop, clean server...")
				})
				scripts.Init(c)
				scripts.RunServer()
				return nil
			},
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "{{rippleApp}}"
	app.Usage = "A {{rippleApp}} application powered by Ripple framework"
	app.Authors = []*cli.Author{{
		Name:  "wangmingjob",
		Email: "wangmingjob@icloud.com",
	}}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "conf",
			Value:       "./config",
			DefaultText: "./config",
			Usage:       "配置文件目录",
		},
		&cli.StringFlag{
			Name:        "env",
			Value:       "dev",
			DefaultText: "dev",
			Usage:       "执行环境 (开发环境dev、测试环境test、线上环境prod)",
		},
	}
	app.Version = "1.0.0"
	app.Commands = Commands()
	_ = app.Run(os.Args)

	line := "==============================="
	fmt.Println(fmt.Sprintf("%s%s%s%s",
		color.White(line),
		color.Bold(color.Green("任务列表")),
		color.Bold(color.Yellow("["+time.Now().Format(DateFullLayout))+"]"),
		color.White(line)))
	fmt.Println(color.Bold(color.White("包含以下任务:")))
	for key, command := range app.Commands {
		fmt.Println(fmt.Sprintf("任务%d：%s %s %s", key+1, command.Name, command.Usage, command.Description))
	}
}
