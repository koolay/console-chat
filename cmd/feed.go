package cmd

import (
	"errors"
	"fmt"

	"github.com/koolay/console-chat/rethink"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	database = "chat"
	address  = "172.105.233.187:28015"
)

func NewFeedCmd() cli.Command {

	feedCmd := cli.Command{
		Name:  "feed",
		Usage: "feed a room",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "room, r"},
		},
		Action: func(c *cli.Context) error {
			options := &rethink.Options{
				Database: database,
				Address:  address,
			}
			r := rethink.NewRethink(options)
			return r.FeedsPublic()
		},
	}
	return feedCmd
}

func NewCreateRoomCmd() cli.Command {

	return cli.Command{
		Name:  "create",
		Usage: "create a room",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n"},
		},
		Action: func(c *cli.Context) error {
			options := &rethink.Options{
				Database: database,
				Address:  address,
			}
			r := rethink.NewRethink(options)
			return r.CreateRoom(c.String("name"))
		},
	}
}

func NewSendCmd() cli.Command {

	return cli.Command{
		Name:  "send",
		Usage: "send message",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "room, r"},
			cli.StringFlag{Name: "to, t"},
		},
		Action: func(c *cli.Context) error {
			options := &rethink.Options{
				Database: database,
				Address:  address,
			}
			r := rethink.NewRethink(options)
			return r.CreateRoom(c.String("name"))
		},
	}
}

func NewJoinCmd() cli.Command {

	return cli.Command{
		Name:  "join",
		Usage: "join app, and register a new user.",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "username, u"},
		},
		Action: func(c *cli.Context) error {
			options := &rethink.Options{
				Database: database,
				Address:  address,
			}
			r := rethink.NewRethink(options)
			username := c.String("username")
			password := ""
			confirmPassword := ""
			prompt := &survey.Password{
				Message: "Please type your password",
			}
			survey.AskOne(prompt, &password, nil)
			confirmPrompt := &survey.Password{
				Message: "Please repeat your password",
			}
			survey.AskOne(confirmPrompt, &confirmPassword, nil)
			if password != confirmPassword {
				return errors.New("passwords not match")
			}
			if username == "" || password == "" {
				return errors.New("username and password should not be empty")
			}
			if err := r.Join(username, password); err != nil {
				return err
			} else {
				fmt.Println("Created successfully!")
			}
			return nil
		},
	}
}

func NewLoginCmd() cli.Command {

	return cli.Command{
		Name:  "login",
		Usage: "login with username and password.",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "username, u"},
		},
		Action: func(c *cli.Context) error {
			options := &rethink.Options{
				Database: database,
				Address:  address,
			}
			r := rethink.NewRethink(options)
			username := c.String("username")
			password := ""
			prompt := &survey.Password{
				Message: "Please type your password",
			}
			survey.AskOne(prompt, &password, nil)
			if err := r.Login(username, password); err != nil {
				return err
			}
			fmt.Printf("Login with %s successfully!\n", username)
			return nil
		},
	}
}
