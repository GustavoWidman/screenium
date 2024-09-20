package client_commands

import (
	"fmt"
	"io"
	"os"
	client_utils "screenium/src/client/utils"

	"github.com/urfave/cli/v2"
)

func ListCommand(_ []string) {
	list()
}

func list() {
	conn := client_utils.QuickConnect(true)
	defer conn.Close()

	fmt.Fprintf(conn, "list\n")
	io.Copy(os.Stdout, conn)
}

func MakeListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Aliases:     []string{"ls"},
		Usage:       client_utils.ColorGreenBold("screenium list"),
		UsageText:   "               ",
		Description: "List all running shells and their statuses",
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
			list()

			return nil
		},
	}
}
