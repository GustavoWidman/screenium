package server_utils

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	server_config "screenium/src/server/config"
	"strings"
	"time"
)

func PingDaemon() bool {
	conn, err := net.Dial("unix", server_config.SOCKET_FILE_NAME)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	doneChan := make(chan bool, 1)
	go func() {
		buf := bufio.NewReader(conn)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			command := strings.Trim(line, " \r\n")
			if command == "pong" {
				doneChan <- true
				return
			}
		}
	}()
	go func() {
		time.Sleep(time.Second * 5)
		conn.Close()
		doneChan <- false
	}()
	fmt.Fprintf(conn, "ping\n")

	return <-doneChan
}

func CheckDaemon() bool {
	cmd := exec.Command("fuser", server_config.DAEMON_PID_FILE_NAME)

	out, err := cmd.Output()

	if err != nil {
		return false
	}

	locked := string(out) != ""

	if locked {
		available := PingDaemon()
		if available {
			return true
		}
	}

	return false
}
