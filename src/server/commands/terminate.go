package server_commands

import (
	"fmt"
	"net"
	server_utils "screenium/src/server/utils"
	"time"
)

// forcibly detaches the currently attached user from a shell
func TerminateCommand(args []string, conn net.Conn) {
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

	if shell.Attached {
		msg := server_utils.ColorYellow("warning") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("someone is currently attached to this shell! Are you sure you want to terminate it? (y/N):")
		conn.Write([]byte(msg))

		// read the input
		input := make([]byte, 1)
		_, err := conn.Read(input)
		if (string(input) != "y" && string(input) != "Y") || err != nil {
			msg := server_utils.ColorYellow("warning") +
				server_utils.ColorGrey(": ") +
				server_utils.ColorWhiteBold("shell ") +
				server_utils.ColorBoldMagenta("\""+args[1]+"\"") +
				server_utils.ColorGrey(" will not be terminated.\n")
			conn.Write([]byte(msg))
			conn.Close()
			return
		}

		// detach gracefully
		shell.DetachChan <- true

		for shell.Attached {
			time.Sleep(time.Millisecond * 10)
		}
	}

	err := shell.Close()

	if err != nil {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("failed to terminate shell ") +
			server_utils.ColorBoldMagenta("\""+shell_name+"\"") +
			server_utils.ColorGrey(fmt.Sprintf(" (PID %d)", shell.PID)) + "\n" +
			server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold(err.Error()) + "\n"
		conn.Write([]byte(msg))
		conn.Close()
		return
	}

	delete(shells, shell.Name)

	msg := server_utils.ColorGreen("success") +
		server_utils.ColorGrey(": ") +
		server_utils.ColorWhiteBold("terminated shell ") +
		server_utils.ColorBoldMagenta("\""+shell_name+"\"") +
		server_utils.ColorGrey(fmt.Sprintf(" (PID %d)", shell.PID)) + "\n"
	conn.Write([]byte(msg))
	conn.Close()
	return
}
