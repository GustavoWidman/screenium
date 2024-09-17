package shell

import (
	"io"
	"log"
	"os"
)

func (shell *Shell) MakeHistoryFile() error {
	// make a tmp file in /tmp/screenium/history/<shell_name>-*.hist
	// if the dir does not exist, recursively create it
	if _, err := os.Stat("/tmp/screenium/history/"); os.IsNotExist(err) {
		os.MkdirAll("/tmp/screenium/history/", 0755)
	}

	tmpFile, err := os.CreateTemp("/tmp/screenium/history/", shell.Name+"-*.hist")
	if err != nil {
		return err
	}

	shell.HistoryFile = tmpFile

	return nil
}

// write out the contents of the history file to the writer
func (shell *Shell) WriteOutHistory(writer io.Writer) {
	_, err := shell.HistoryFile.Seek(0, 0)
	if err != nil {
		return
	}

	_, err = io.Copy(writer, shell.HistoryFile)
	if err != nil {
		return
	}
}

// write the contents of the reader to the history file
func (shell *Shell) WriteInHistory(reader io.Reader) {
	_, err := shell.HistoryFile.Seek(0, 2)
	if err != nil {
		return
	}

	writer := shell.HistoryFile

	log.Println("Started writing to history")
	_, err = io.Copy(writer, reader)
	log.Println("Finished writing to history")
	if err != nil {
		return
	}
}

func (shell *Shell) MakeInterceptedWriter(writers ...io.Writer) io.Writer {
	_, err := shell.HistoryFile.Seek(0, 2)
	if err != nil {
		return nil
	}

	writers = append([]io.Writer{shell.HistoryFile}, writers...)
	writer := io.MultiWriter(writers...)

	return writer
}
