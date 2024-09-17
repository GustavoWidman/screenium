package server_commands

import "net"

func PingCommand(_ []string, conn net.Conn) {
	conn.Write([]byte("pong\n"))
	conn.Close()
	return
}
