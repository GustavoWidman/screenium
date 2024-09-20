package main

import (
	"fmt"
	"log"
	"os"
	"screenium/src/client"
	client_utils "screenium/src/client/utils"
	"screenium/src/server"
)

func main() {
	is_daemon := os.Getenv("_GO_DAEMON")
	is_screenium := os.Getenv("SCREENIUM_SHELL")

	if is_screenium != "" {
		fmt.Print(client_utils.ColorRed("error") +
			client_utils.ColorGrey(": ") +
			client_utils.ColorWhiteBold("screenium is not meant to be run inside another screenium session!\n") +
			client_utils.ColorGreen("tip") +
			client_utils.ColorGrey(": ") +
			client_utils.ColorWhiteBold("If you wish to detach you can do so by pressing Ctrl+X.\n"))
		return
	}

	if is_daemon != "" {
		server.ListenToSocket()
		return
	}

	server.RunDaemon()

	if err := client.MakeApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
