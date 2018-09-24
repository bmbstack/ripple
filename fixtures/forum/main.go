package main

import (
	"fmt"
	"github.com/bmbstack/ripple/fixtures/forum/scripts"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli"
	"os"
)

func main() {
	line := "==============================="
	app := cli.NewApp()
	app.Name = "forum"
	app.Usage = "A forum application powered by Ripple framework"
	app.Author = "wangmingjob"
	app.Email = "wangmingjob@icloud.com"
	app.Version = "0.0.1"
	app.Commands = scripts.Commands()
	fmt.Println(fmt.Sprintf("%s%s%s", color.White(line), color.Bold(color.Green("Application command")), color.White(line)))
	app.Run(os.Args)
}
