package server_commands

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"screenium/src/server/utils/shell"
	"strconv"
	"strings"
)

func IOControlCommand(args []string, conn net.Conn) {
	if len(args) < 3 {
		conn.Write([]byte("Usage: io <shell-name> <shell-session-id>\n"))
		conn.Close()
		return
	}
	shell_name := args[1]
	session_id := args[2]
	shell, ok := shells[shell_name]
	if !ok {
		conn.Write([]byte(fmt.Sprintf("A shell with the name %s does not exist.\n", shell_name)))
		conn.Close()
		return
	}

	if shell.SessionID != session_id {
		conn.Write([]byte("This shell does not have the session id you are looking for!\n"))
		conn.Close()
		return
	}

	fmt.Fprint(conn, "ack\n")

	attachShellToConn(conn, shell)
	return
}

func attachShellToConn(conn net.Conn, shell *shell.Shell) {
	// wait for initial resize
	initialResize(conn, shell)

	// there has to be a better way to do this :P
	fmt.Fprintf(conn, "[attached to \"\033[1;97m%s\033[0m\" \033[38;2;127;127;127m(PID %d)\033[0m]\n\r", shell.Name, shell.PID)

	detachChan := shell.Attach(conn, conn)

	go func() {
		detached := <-detachChan

		if !detached {
			shell.Close()
			delete(shells, shell.Name)
			conn.Close()
		} else {
			shell.Attached = false
			shell.SessionID = ""
			log.Printf("Detached shell %s\n", shell.Name)
			conn.Close()
		}
		return
	}()
}

func initialResize(conn net.Conn, shell *shell.Shell) {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	pair := strings.Split(strings.Trim(line, " \r\n"), ",")
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
}
