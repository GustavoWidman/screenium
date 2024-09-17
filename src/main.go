package main

import (
	"log"
	"os"
	"screenium/src/client"
	"screenium/src/server"
)

func main() {
	is_daemon := os.Getenv("_GO_DAEMON")

	if is_daemon != "" {
		server.ListenToSocket()
		return
	} else {
		server.RunDaemon()
	}

	// client.EvaluateArgs(os.Args)

	if err := client.MakeApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
