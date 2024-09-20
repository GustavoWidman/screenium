package client_commands

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	client_utils "screenium/src/client/utils"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// invokes the create command on the server
func CreateCommand(args []string) {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program create <shell_name>")
		return
	}
	shell_name := os.Args[2]
	shell_path := os.Getenv("SHELL")

	if shell_path == "" {
		fmt.Println("SHELL environment variable not set, defaulting to /bin/sh")
		shell_path = "/bin/sh"
		return
	}

	create(shell_name, shell_path)
	return
}

func MakeCreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "create",
		Usage:       client_utils.ColorGreenBold("screenium create ") + client_utils.ColorCyan("[flags]") + " " + "<shell name>",
		UsageText:   client_utils.ColorGrey("shell-name"),
		Description: "Creates a new screenium shell session",
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
			&cli.StringFlag{
				Name:     "shell",
				Required: true,
				EnvVars:  []string{"SHELL"},
				Aliases:  []string{"s"},
				Usage:    "Overrides the SHELL environment variable",
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
			shell_path := c.String("shell")

			create(shell_name, shell_path)
			return nil
		},
	}
}

func create(shell_name, shell_path string) {
	conn, err := net.Dial("unix", "/tmp/screenium/manager.sock")
	if err != nil {
		log.Fatalf("Failed to connect to the socket: %v\n", err)
		return
	}

	defer conn.Close()

	pidChan := make(chan string, 1)
	go func() {
		buf := bufio.NewReader(conn)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			trimmed := strings.Trim(line, " \r\n")

			pidChan <- trimmed
		}
	}()
	go func() {
		time.Sleep(time.Second * 5)
		conn.Close()
		pidChan <- ""
	}()
	fmt.Fprintf(conn, "create %s %s\n", shell_name, shell_path)

	// print the pid
	pidStr := <-pidChan

	// try to parse pid as int
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Println(pidStr)
		return
	}

	// if pid is -1, shell was not created
	if pid == -1 {
		fmt.Println("Failed to create shell")
		return
	}

	attach(shell_name)
}
