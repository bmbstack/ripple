package main

import (
	"fmt"
	"github.com/bmbstack/ripple/cmd/ripple/templates/scripts"
	"github.com/urfave/cli"
	"github.com/labstack/gommon/color"
	"os"
)

func main() {
	line := "==============================="
	app := cli.NewApp()
	app.Name = "{{rippleApp}}"
	app.Usage = "A {{rippleApp}} application powered by Ripple framework"
	app.Author = "wangmingjob"
	app.Email = "wangmingjob@icloud.com"
	app.Version = "0.0.1"
	app.Commands = scripts.Commands()
	fmt.Println(fmt.Sprintf("%s%s%s", color.White(line), color.Bold(color.Green("Application command")), color.White(line)))
	app.Run(os.Args)
}
