package client_utils

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	server_config "screenium/src/server/config"
	"strings"
	"time"
)

// returns a connection to the socket
func QuickConnect(panicOnError bool) net.Conn {
	conn, err := net.Dial("unix", server_config.SOCKET_FILE_NAME)
	if err != nil {
		if panicOnError {
			panic(err)
		}
		return nil
	}
	return conn
}

func QuickCommand(command string, panicOnError bool) net.Conn {
	conn := QuickConnect(panicOnError)
	if conn == nil {
		return nil
	}

	fmt.Fprintf(conn, command)

	return conn
}

type CommandResult string

const (
	CommandResultSuccess CommandResult = "success"
	CommandResultError   CommandResult = "error"
	CommandResultTimeout CommandResult = "timeout"
	CommandResultUnknown CommandResult = "unknown"
)

/*
awaits for a specific string to be received before returning

if timeout given is time.Duration(0), waits forever

returns CommandResultSuccess if the command was received

returns CommandResultError if the command was not received

returns CommandResultTimeout if the command was not received before the timeout

returns CommandResultUnknown if the command was not received and the response was unknown

in case of a CommandResultUnknown, the second return value is the unknown command response
*/
func WaitForCommand(conn net.Conn, command string, timeout time.Duration) (CommandResult, string) {
	ch := make(chan string, 1)
	go func() {
		buf := bufio.NewReader(conn)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				ch <- err.Error()
				return
			}
			cmd := strings.Trim(line, " \r\n")
			if cmd == command {
				ch <- "success"
				return
			} else {
				ch <- fmt.Sprintf("unknown:%s", cmd)
			}
		}
	}()

	if timeout != time.Duration(0) {
		go func() {
			time.Sleep(timeout)
			ch <- "timeout"
		}()
	}

	result := <-ch

	if strings.HasPrefix(result, "unknown:") {
		return CommandResultUnknown, strings.TrimPrefix(result, "unknown:")
	}

	switch result {
	case "success":
		return CommandResultSuccess, ""
	case "timeout":
		return CommandResultTimeout, ""
	default:
		return CommandResultError, result
	}
}

func WaitForSomething(conn net.Conn, timeout time.Duration) (string, error) {
	ch := make(chan string, 1)
	go func() {
		buf := bufio.NewReader(conn)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				ch <- fmt.Sprintf("err:%s", err.Error())
				return
			}
			cmd := strings.Trim(line, " \r\n")
			ch <- cmd
		}
	}()
	if timeout != time.Duration(0) {
		go func() {
			time.Sleep(timeout)
			ch <- "err:timeout"
		}()
	}
	result := <-ch

	if strings.HasPrefix(result, "err:") {
		return "", errors.New(strings.TrimPrefix(result, "err:"))
	}

	return result, nil
}
