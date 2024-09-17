package server_commands

import (
	"fmt"
	"log"
	"net"
	server_utils "screenium/src/server/utils"
	"screenium/src/server/utils/shell"
)

func CreateCommand(args []string, conn net.Conn) {
	if len(args) < 2 {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("improper client-side usage!\n")

		conn.Write([]byte(msg))
		conn.Close()
		return
	}
	shell_name := args[1]
	shell_path := args[2]

	// check if shell with that name already exists
	_, exists := shells[shell_name]
	if exists {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("shell ") +
			server_utils.ColorBoldMagenta("\""+shell_name+"\"") +
			server_utils.ColorWhiteBold(" already exists!\n")

		conn.Write([]byte(msg))
		conn.Close()
		return
	}

	shell, err := shell.CreateShell(shell_name, shell_path)
	if err != nil {
		log.Println(err)
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("internal error! check logs.\n")

		conn.Write([]byte(msg))
		conn.Close()
		return
	}
	shells[shell_name] = &shell
	conn.Write([]byte(fmt.Sprintf("%d\n", shell.PID)))
	conn.Close()
	return
}
