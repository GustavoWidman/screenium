package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	server_commands "screenium/src/server/commands"
	server_config "screenium/src/server/config"
	server_utils "screenium/src/server/utils"
	"strings"
	"time"

	"github.com/sevlyar/go-daemon"
)

func RunDaemon() {
	screenium_folder := path.Dir(server_config.SOCKET_FILE_NAME)
	if _, err := os.Stat(screenium_folder); os.IsNotExist(err) {
		err := os.MkdirAll(screenium_folder, 0755)
		if err != nil {
			panic(err)
		}
	}

	if server_utils.CheckDaemon() {
		// log.Println("Daemon is already running")
		return
	} else {
		// log.Println("Starting daemon")
	}

	context := daemon.Context{
		PidFileName: server_config.DAEMON_PID_FILE_NAME,
		PidFilePerm: 0644,
		LogFileName: server_config.DAEMON_LOG_FILE_NAME,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"screenium-daemon"},
	}

	child, _ := context.Reborn()

	if child != nil {
		timeout := time.After(time.Second * 5)
		for {
			select {
			case <-timeout:
				log.Println("Failed to bringup daemon")
				os.Exit(1)
				return
			default:
				if server_utils.CheckDaemon() {
					return
				}
				time.Sleep(time.Millisecond * 5)
			}
		}

	} else {
		defer context.Release()
		ListenToSocket()
	}
}

func ListenToSocket() {
	server_utils.OverrideSignals()

	listener, err := net.Listen("unix", server_config.SOCKET_FILE_NAME)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		log.Println("Accepted new connection.")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := bufio.NewReader(conn)
	line, err := buf.ReadString('\n')
	if err != nil {
		log.Println("Error reading command:", err)
		fmt.Fprintf(conn, "internal error: %s\n", err.Error())
		conn.Close()
		return
	}

	args := strings.Split(strings.Trim(line, " \r\n"), " ")

	switch args[0] {
	case "ping":
		server_commands.PingCommand(args, conn)
	case "create":
		server_commands.CreateCommand(args, conn)
	case "attach":
		server_commands.AttachCommand(args, conn)
	case "io":
		server_commands.IOControlCommand(args, conn)
	case "list":
		server_commands.ListCommand(args, conn)
	case "kill":
		server_commands.KillCommand(args, conn)
	case "detach":
		server_commands.DetachCommand(args, conn)
	case "terminate":
		server_commands.TerminateCommand(args, conn)
	default:
		server_commands.UnknownCommand(args, conn)
	}
}
