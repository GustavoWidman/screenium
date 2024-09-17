package client

import (
	"fmt"
	client_commands "screenium/src/client/commands"
	client_utils "screenium/src/client/utils"

	"github.com/urfave/cli/v2"
)

func EvaluateArgs(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: program [daemon|create|list|attach]")
		return
	}

	switch args[1] {
	case "create":
		client_commands.CreateCommand(args)
	case "attach":
		client_commands.AttachCommand(args)
	case "list":
		client_commands.ListCommand(args)
	case "kill":
		client_commands.KillCommand(args)
	case "detach":
		client_commands.DetachCommand(args)
	// case "runlocal":
	// 	shell_name := "local"
	// 	shell_path := "/bin/bash"
	// 	shell, err := shell.Create(shell_name, shell_path)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	shell.Attach()
	default:
		client_commands.UnknownCommand(args)
	}
}

func MakeApp() *cli.App {
	return &cli.App{
		Name:        client_utils.ColorBoldMagenta("screenium"),
		Description: "is a modern remake of GNU \"screen\" in golang",
		Usage:       client_utils.ColorWhiteBold("screenium <command> ") + client_utils.ColorBoldCyan("[...flags]") + client_utils.ColorWhiteBold(" [...args]"),
		Authors: []*cli.Author{{
			Name:  client_utils.ColorBoldMagenta("r3dlust"),
			Email: "admin@r3dlust.com",
		}},
		Version:         client_utils.ColorGrey("1.0.0-beta"),
		HelpName:        "screenium",
		HideHelpCommand: true,
		DefaultCommand:  "help",

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
			client_commands.MakeKillCommand(),
			client_commands.MakeListCommand(),
			client_commands.MakeHelpCommand(),
		},
	}
}
