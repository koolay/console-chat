package main

import (
	"log"
	"os"

	"github.com/koolay/console-chat/cmd"
)

func main() {

	app := cmd.NewApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
