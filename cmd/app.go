// Package main provides ...
package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "console-chat"
	app.Usage = "realtime chat"
	app.Commands = []cli.Command{
		NewFeedCmd(),
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello friend!")
		return nil
	}
	return app
}
