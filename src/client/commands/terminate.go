package client_commands

import (
	"fmt"
	"io"
	"os"
	client_utils "screenium/src/client/utils"

	"github.com/urfave/cli/v2"
)

func MakeTerminateCommand() *cli.Command {
	return &cli.Command{
		Name:        "terminate",
		Aliases:     []string{"term", "t"},
		Usage:       client_utils.ColorGreenBold("screenium terminate ") + "<shell name>",
		UsageText:   client_utils.ColorGrey("shell-name"),
		Description: "Forcibly terminates a given shell session.\n",
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
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowSubcommandHelp(c)
				fmt.Println(client_utils.ColorRed("error") + client_utils.ColorGrey(": ") + client_utils.ColorWhiteBold("shell name not specified"))
				return nil
			}

			shell_name := c.Args().Get(0)

			terminate(shell_name)

			return nil
		},
	}
}

func terminate(shell_name string) {
	conn := client_utils.QuickCommand(fmt.Sprintf("terminate %s\n", shell_name), true)
	defer conn.Close()

	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
	return
}
