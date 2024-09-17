package client_commands

import (
	"fmt"
	client_utils "screenium/src/client/utils"
)

func UnknownCommand(args []string) {
	fmt.Println(client_utils.ColorRed("error") +
		client_utils.ColorGrey(": ") +
		client_utils.ColorWhiteBold("unknown command ") +
		client_utils.ColorBoldMagenta("\""+args[0]+"\"") +
		client_utils.ColorGrey(".\n"),
	)
	return
}
