package shell

import (
	"io"
)

func (shell *Shell) Attach(in io.Reader, out io.Writer) chan bool {
	shell.WriteOutHistory(out)

	detachChan := make(chan bool)

	wrappedWriter := &CustomWriter{
		wrappedWriter:   shell.Pty,
		excludeByte:     byte(24),
		ExcludeByteChan: detachChan,
	}

	go func() {
		io.Copy(wrappedWriter, in)
		detachChan <- false
	}()

	go func() {
		writer := shell.MakeInterceptedWriter(out)
		io.Copy(writer, shell.Pty)
		detachChan <- false
	}()

	go func() {
		select {
		case <-shell.DetachChan:
			detachChan <- true
		case <-detachChan:
			// If detachChan receives a signal, this goroutine should exit
			return
		}
	}()

	return detachChan
}

type CustomWriter struct {
	wrappedWriter   io.Writer
	excludeByte     byte
	ExcludeByteChan chan bool // notifies when the byte has been captured and excluded
}

// Write method for CustomWriter that filters out the specified byte before writing.
func (cw *CustomWriter) Write(p []byte) (n int, err error) {
	var filteredData []byte
	for _, b := range p {
		if b == cw.excludeByte {
			cw.ExcludeByteChan <- true
		} else {
			filteredData = append(filteredData, b)
		}
	}
	n, err = cw.wrappedWriter.Write(filteredData)
	return n, err
}
