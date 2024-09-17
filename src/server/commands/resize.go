package server_commands

import (
	"net"
	"screenium/src/server/utils/shell"
	"strconv"
	"strings"
)

func ResizeCommand(args []string, conn net.Conn, shell *shell.Shell) {
	if len(args) < 2 {
		conn.Write([]byte("Usage: resize <rows>,<cols>\n"))
		return
	}

	pair := strings.Split(args[1], ",")
	if len(pair) != 2 {
		conn.Write([]byte("Usage: resize <rows>,<cols>\n"))
		return
	}

	rows, err := strconv.Atoi(pair[0])
	if err != nil {
		conn.Write([]byte("Usage: resize <rows>,<cols>\n"))
		return
	}

	cols, err := strconv.Atoi(pair[1])
	if err != nil {
		conn.Write([]byte("Usage: resize <rows>,<cols>\n"))
		return
	}

	shell.Resize(rows, cols)
	return
}
