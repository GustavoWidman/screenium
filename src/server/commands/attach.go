package server_commands

import (
	"bufio"
	"fmt"
	"net"
	server_utils "screenium/src/server/utils"
	"strings"

	"github.com/google/uuid"
)

func AttachCommand(args []string, conn net.Conn) {
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
			server_utils.ColorWhiteBold(" does not exist!\n")

		conn.Write([]byte(msg))
		conn.Close()
		return
	}

	if shell.Attached {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("shell ") +
			server_utils.ColorBoldMagenta("\""+shell_name+"\"") +
			server_utils.ColorGrey(fmt.Sprintf(" (PID %d)", shell.PID)) +
			server_utils.ColorWhiteBold(" is currently attached to another terminal!\n")

		conn.Write([]byte(msg))
		conn.Close()
		return
	}
	shell.Attached = true
	sessionID := uuid.New().String()
	shell.SessionID = sessionID
	fmt.Fprintf(conn, "%s\n", sessionID)

	// this will be turned into the connection for resizing and sending output that isnt shell i/o
	// shell i/o has been moved to the "io" command, which requires the shell session id
	reader := bufio.NewReader(conn)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			cmd := strings.Split(strings.Trim(line, " \r\n"), " ")[0]
			switch cmd {
			case "resize":
				ResizeCommand(args, conn, shell)
			}
		}
	}()

	return
}
