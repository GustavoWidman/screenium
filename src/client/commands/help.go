package client_commands

import (
	"fmt"
	client_utils "screenium/src/client/utils"

	"github.com/urfave/cli/v2"
)

func MakeHelpCommand() *cli.Command {
	return &cli.Command{
		Name:        "help",
		Aliases:     []string{"h"},
		Usage:       client_utils.ColorGreenBold("screenium help"),
		UsageText:   "",
		Description: "Shows the help message.",
		Hidden:      true,
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
			cli.ShowAppHelp(c)
			return nil
		},
	}
}
