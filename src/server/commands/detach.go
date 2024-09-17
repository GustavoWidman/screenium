package server_commands

import (
	"fmt"
	"net"
	server_utils "screenium/src/server/utils"
	"time"
)

// forcibly detaches the currently attached user from a shell
func DetachCommand(args []string, conn net.Conn) {
	if len(args) < 2 {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("improper client-side usage!\n")
		conn.Write([]byte(msg))
		conn.Close()
		return
	}
	shell_name := args[1]
	shell, ok := shells[shell_name]
	if !ok {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("shell ") +
			server_utils.ColorBoldMagenta("\""+shell_name+"\"") +
			server_utils.ColorGrey(" does not exist!\n")
		conn.Write([]byte(msg))
		conn.Close()
		return
	}

	if !shell.Attached {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("shell ") +
			server_utils.ColorBoldMagenta("\""+shell_name+"\"") +
			server_utils.ColorGrey(fmt.Sprintf(" (PID %d)", shell.PID)) +
			server_utils.ColorWhiteBold(" is not attached to anyone right now!\n")
		conn.Write([]byte(msg))
		conn.Close()
		return
	}

	shell.DetachChan <- true

	for shell.Attached {
		time.Sleep(time.Millisecond * 10)
	}

	msg := server_utils.ColorGreen("success") +
		server_utils.ColorGrey(": ") +
		server_utils.ColorWhiteBold("detached shell ") +
		server_utils.ColorBoldMagenta("\""+shell_name+"\"") +
		server_utils.ColorGrey(fmt.Sprintf(" (PID %d)", shell.PID)) + "\n"
	conn.Write([]byte(msg))
	conn.Close()
	return
}
