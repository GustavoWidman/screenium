package shell

import (
	"github.com/creack/pty"
)

func (shell *Shell) Resize(rows, cols int) error {
	winSize := pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	}

	if err := pty.Setsize(shell.Pty, &winSize); err != nil {
		return err
	}

	return nil
}
