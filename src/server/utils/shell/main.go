package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type Shell struct {
	Name        string
	Cmd         *exec.Cmd
	Pty         *os.File
	HistoryFile *os.File
	PID         int
	Attached    bool
	SessionID   string
	DetachChan  chan bool
}

// type ReaderPair struct {
// 	Reader    *io.PipeReader
// 	TeeReader io.Reader
// }

func CreateShell(shell_name string, shell_path string) (Shell, error) {
	c := exec.Command(shell_path)

	ptty, err := pty.Start(c)
	if err != nil {
		return Shell{}, err
	}

	shell := Shell{
		Name: shell_name,
		Cmd:  c,
		Pty:  ptty,
	}

	fmt.Println("PID:", c.Process.Pid)

	shell.PID = c.Process.Pid

	c.Process.Release()

	if err := shell.MakeHistoryFile(); err != nil {
		return Shell{}, err
	}

	shell.DetachChan = make(chan bool)

	return shell, nil
}

func (shell *Shell) Close() error {
	return shell.Pty.Close()
}
