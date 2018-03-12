// Package main provides ...
package cmd

import (
	"github.com/urfave/cli"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "console-chat"
	app.Usage = "realtime chat"
	app.Commands = []cli.Command{
		NewFeedCmd(),
		NewCreateRoomCmd(),
		NewJoinCmd(),
		NewLoginCmd(),
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}
	return app
}
