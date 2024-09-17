package client_commands

import (
	"fmt"
	"io"
	"os"
	client_utils "screenium/src/client/utils"

	"github.com/urfave/cli/v2"
)

func DetachCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: detach <shell-name>")
		return
	}
	shell_name := args[2]

	detach(shell_name)
}

func detach(shell_name string) {
	conn := client_utils.QuickCommand(fmt.Sprintf("detach %s\n", shell_name), true)
	defer conn.Close()

	io.Copy(os.Stdout, conn)
	return
}

func MakeDetachCommand() *cli.Command {
	return &cli.Command{
		Name:        "detach",
		Usage:       client_utils.ColorGreenBold("screenium detach ") + "<shell name>",
		UsageText:   client_utils.ColorGrey("shell-name"),
		Description: "Forcibly detaches anyone that is currently the given shell\n",
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

			detach(shell_name)

			return nil
		},
	}
}
