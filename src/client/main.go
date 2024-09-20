package client

import (
	"fmt"
	client_commands "screenium/src/client/commands"
	client_utils "screenium/src/client/utils"

	"github.com/urfave/cli/v2"
)

func MakeApp() *cli.App {
	return &cli.App{
		Name:        client_utils.ColorBoldMagenta("screenium"),
		Description: "is a modern remake of GNU \"screen\" in golang",
		Usage:       client_utils.ColorWhiteBold("screenium <command> ") + client_utils.ColorBoldCyan("[...flags]") + client_utils.ColorWhiteBold(" [...args]"),
		Authors: []*cli.Author{{
			Name:  client_utils.ColorBoldMagenta("r3dlust"),
			Email: "admin@r3dlust.com",
		}},
		Version:         client_utils.ColorGrey("1.0.1-beta"),
		HelpName:        "screenium",
		HideHelpCommand: true,
		CustomAppHelpTemplate: fmt.Sprintf(`{{.Name}} {{.Description}} %s

%s {{.Usage}}
{{if .Commands}}
%s
{{range .Commands}}{{if not .Hidden}}   %s{{ "\t"}}{{.UsageText}}{{ "\t"}}{{.Description}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
%s
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if len .Authors}}
%s
   {{range .Authors}}{{ . }}{{end}}{{end}}
`,
			client_utils.ColorGrey("({{.Version}}")+client_utils.ColorGrey(")"),
			client_utils.ColorWhiteBold("Usage:"),
			client_utils.ColorWhiteBold("Commands:"),
			client_utils.ColorBlueBold("{{.Name}}"),
			client_utils.ColorWhiteBold("Global Options:"),
			client_utils.ColorWhiteBold("Author:"),
		),
		Commands: []*cli.Command{
			client_commands.MakeCreateCommand(),
			client_commands.MakeAttachCommand(),
			client_commands.MakeDetachCommand(),
			client_commands.MakeTerminateCommand(),
			client_commands.MakeKillCommand(),
			client_commands.MakeListCommand(),
			// client_commands.MakeHelpCommand(),
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)

			if c.Args().Len() > 0 {
				client_commands.UnknownCommand(c.Args())
			}

			return nil
		},
	}
}
