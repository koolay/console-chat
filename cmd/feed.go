package cmd

import (
	"log"

	"github.com/koolay/console-chat/rethink"
	"github.com/urfave/cli"
)

func NewFeedCmd() cli.Command {

	feedCmd := cli.Command{
		Name:  "feed",
		Usage: "feed a room",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "room, r"},
		},
		Action: func(c *cli.Context) error {
			options := &rethink.RethinkOptions{
				Database: "dmp",
				Address:  "172.105.233.187:28015",
			}
			rth := rethink.NewRethink(options)
			if err := rth.Feeds(c.String("room")); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
	return feedCmd
}
