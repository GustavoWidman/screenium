package shell

import (
	"errors"
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

	shell.PID = c.Process.Pid

	if err := shell.MakeHistoryFile(); err != nil {
		return Shell{}, err
	}

	shell.DetachChan = make(chan bool)

	return shell, nil
}

func (shell *Shell) Close() error {
	_, err := shell.Pty.Write([]byte("exit\n"))

	if err != nil {
		return err
	}

	state, err := shell.Cmd.Process.Wait()

	if err != nil {
		return err
	}

	if !state.Exited() {
		return errors.New("shell did not exit cleanly")
	}

	err = shell.Pty.Close()

	if err != nil {
		return err
	}

	return nil
}
