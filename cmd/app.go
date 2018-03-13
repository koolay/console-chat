// Package main provides ...
package cmd

import (
	"fmt"
	"os"

	"github.com/koolay/console-chat/rethink"
	"github.com/urfave/cli"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "console-chat"
	app.Usage = "realtime chat"
	app.Commands = []cli.Command{
		NewFeedCmd(),
		NewCreateRoomCmd(),
		NewJoinCmd(),
	}
	app.Action = func(c *cli.Context) error {
		var username string
		var password string
		var errCount = 0
		for {
			promptUsername := &survey.Input{
				Message: "Please type your username",
			}
			survey.AskOne(promptUsername, &username, nil)
			promptPassword := &survey.Password{
				Message: "Please type your password",
			}
			survey.AskOne(promptPassword, &password, nil)
			err := rethink.RethinkActor.Login(username, password)
			if err == nil {
				rethink.LogInfo("Login with %s successfully!\n", username)
				break
			} else if errCount > 3 {
				rethink.LogInfo("Invalid username or passwn")
				os.Exit(0)
			} else {
				errCount++
				rethink.LogError(err.Error())
			}
		}

		// list rooms and choose room in
		selectedRoom := ""
		rooms, err := rethink.RethinkActor.GetAllRooms()
		if err != nil {
			rethink.LogError(err.Error())
			os.Exit(0)
		}
		prompt := &survey.Select{
			Message: "Choose a room:",
			Options: rooms,
		}
		survey.AskOne(prompt, &selectedRoom, nil)
		fmt.Printf("Join in room %s \n", selectedRoom)
		err = rethink.RethinkActor.SwitchRoom(selectedRoom)
		if err != nil {
			return err
		}

		ui := NewChatUI()
		return ui.Run()
	}
	return app
}
