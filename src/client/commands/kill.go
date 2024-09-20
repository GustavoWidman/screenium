package client_commands

import (
	"fmt"
	client_utils "screenium/src/client/utils"
	server_utils "screenium/src/server/utils"

	"github.com/urfave/cli/v2"
)

// kills the daemon
func KillCommand(args []string) {
	conn := client_utils.QuickCommand("kill\n", true)
	defer conn.Close()

	return
}

func MakeKillCommand() *cli.Command {
	return &cli.Command{
		Name:        "kill",
		Aliases:     []string{"k"},
		Usage:       client_utils.ColorGreenBold("screenium kill"),
		UsageText:   "",
		Description: "Forcibly kills the screenium daemon if it is running",
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
			conn := client_utils.QuickCommand("kill\n", true)
			defer conn.Close()

			if server_utils.CheckDaemon() {
				fmt.Println(client_utils.ColorRed("error") +
					client_utils.ColorGrey(": ") +
					client_utils.ColorWhiteBold("daemon is still running! wtf?"),
				)
			} else {
				fmt.Println(client_utils.ColorGreen("success") +
					client_utils.ColorGrey(": ") +
					client_utils.ColorWhiteBold("daemon has been killed"),
				)
			}

			return nil
		},
	}
}
