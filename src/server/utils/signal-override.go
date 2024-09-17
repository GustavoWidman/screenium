package server_utils

import (
	"log"
	"os"
	"os/signal"
	server_config "screenium/src/server/config"
	"syscall"
)

// override all the signals to properly shutdown the daemon
func OverrideSignals() {
	signal.Notify(server_config.ServerSignalChan,
		os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, // keyboard
		syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, // os termination
	)
	go func() {
		for range server_config.ServerSignalChan {
			log.Println("Received signal, shutting down...")
			os.Remove(server_config.DAEMON_PID_FILE_NAME)
			os.Remove(server_config.SOCKET_FILE_NAME)
			os.Exit(0)
		}
	}()
}
