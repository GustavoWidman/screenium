package server_commands

import (
	"net"
	server_config "screenium/src/server/config"
	"syscall"
)

// kills the daemon
func KillCommand(args []string, conn net.Conn) {
	// that should do it
	server_config.ServerSignalChan <- syscall.SIGINT
	return
}
