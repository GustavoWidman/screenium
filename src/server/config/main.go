package server_config

import "os"

const (
	SOCKET_FILE_NAME     = "/tmp/screenium/manager.sock"
	DAEMON_PID_FILE_NAME = "/tmp/screenium/manager.pid"
	DAEMON_LOG_FILE_NAME = "/tmp/screenium/manager.log"
)

var ServerSignalChan = make(chan os.Signal, 1)
