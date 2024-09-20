package client_commands

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	client_utils "screenium/src/client/utils"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func AttachCommand(args []string) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program attach <shell-name>")
		return
	}
	shell_name := os.Args[2]
	attach(shell_name)
	return
}

func attach(shell_name string) {
	conn := client_utils.QuickCommand(fmt.Sprintf("attach %s\n", shell_name), true)
	defer conn.Close()

	// wait until the server acknowledges the attachment
	sessionID, err := client_utils.WaitForSomething(conn, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	// if sessionID is not a valid UUID, its an error and we should print it and exit
	if _, err := uuid.Parse(sessionID); err != nil {
		fmt.Println(sessionID)
		return
	}

	ioConn := client_utils.QuickConnect(true)
	ioExitChan := make(chan bool)

	// hook up sigwinch to resize the session
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			width, height, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(conn, "resize %d,%d\n", height, width)
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	go func() {
		// if we get any output on the normal conn, immediately stop the ioConn by sending true to ioExitChan and print the output
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			ioExitChan <- true
			fmt.Println(line)
		}
	}()

	AttachToConn(ioConn, shell_name, sessionID, ioExitChan)
}

func AttachToConn(conn net.Conn, shell_name, sessionID string, ch chan bool) {
	defer conn.Close()
	fmt.Fprintf(conn, "io %s %s\n", shell_name, sessionID)
	result, unknown := client_utils.WaitForCommand(conn, "ack", time.Second*5)
	if result != client_utils.CommandResultSuccess {
		fmt.Println("Failed to attach to the shell")
		if unknown != "" {
			fmt.Println(unknown)
		}
		return
	}

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	go func() {
		io.Copy(conn, os.Stdin)
		ch <- true
	}()
	go func() {
		io.Copy(os.Stdout, conn)
		ch <- true
	}()

	// send back our ack since now we are ready
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(conn, "%d,%d\n", height, width)

	<-ch

	_ = term.Restore(int(os.Stdin.Fd()), oldState)

	cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	return
}

func MakeAttachCommand() *cli.Command {
	return &cli.Command{
		Name:        "attach",
		Aliases:     []string{"a"},
		Usage:       client_utils.ColorGreenBold("screenium attach ") + client_utils.ColorCyan("[flags]") + " " + "<shell name>",
		UsageText:   client_utils.ColorGrey("shell-name"),
		Description: "Attaches to a pre-existing screenium shell session",
		CustomHelpTemplate: fmt.Sprintf(`%s: {{.Description}}

%s
    {{.Usage}}
{{if .VisibleFlags}}
%s
    {{range .VisibleFlags}}{{.}}
    {{end}}{{end}}
`,
			client_utils.ColorBoldCyan("{{.Name}}"),
			client_utils.ColorWhiteBold("Usage:"),
			client_utils.ColorWhiteBold("Flags:"),
		),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "detach",
				Required: false,
				Aliases:  []string{"d"},
				Usage:    "Forcibly detaches anyone that is currently attached to the shell",
			},
			&cli.BoolFlag{
				Name:     "debug",
				Usage:    "Enables debug mode",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowSubcommandHelp(c)
				fmt.Println(client_utils.ColorRed("error") + client_utils.ColorGrey(": ") + client_utils.ColorWhiteBold("shell name not specified"))
				return nil
			}

			shell_name := c.Args().Get(0)

			if c.Bool("detach") {
				detach(shell_name)
			}

			attach(shell_name)
			return nil
		},
	}
}
