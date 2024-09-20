package client_commands

import (
	"fmt"
	client_utils "screenium/src/client/utils"

	"github.com/urfave/cli/v2"
)

func UnknownCommand(args cli.Args) {
	cmd := args.Get(0)

	if cmd == "help" || cmd == "h" {
		return
	}

	fmt.Println(client_utils.ColorRed("\nerror") +
		client_utils.ColorGrey(": ") +
		client_utils.ColorWhiteBold("unknown command ") +
		client_utils.ColorBoldMagenta("\""+cmd+"\"") +
		client_utils.ColorGrey("."),
	)
	return
}
