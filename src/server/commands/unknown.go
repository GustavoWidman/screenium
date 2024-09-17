package server_commands

import "net"

func UnknownCommand(args []string, conn net.Conn) {
	conn.Write([]byte("error: unknown command \"" + args[0] + "\"\n"))
	conn.Close()
	return
}
