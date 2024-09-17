package server_commands

import (
	"fmt"
	"net"
	server_utils "screenium/src/server/utils"
	"text/tabwriter"
)

func ListCommand(_ []string, conn net.Conn) {
	if len(shells) == 0 {
		msg := server_utils.ColorRed("error") +
			server_utils.ColorGrey(": ") +
			server_utils.ColorWhiteBold("no shells are currently running!\n")

		conn.Write([]byte(msg))
		conn.Close()
		return
	}

	prettyString := server_utils.ColorWhiteBold("Shells:\n")
	w := tabwriter.NewWriter(conn, 20, 30, 1, '\t', tabwriter.AlignRight)
	fmt.Fprint(w, prettyString)
	for shell_name, shell := range shells {
		attachedStr := ""
		switch shell.Attached {
		case true:
			attachedStr = server_utils.ColorRed("attached")
		case false:
			attachedStr = server_utils.ColorGreen("detached")
		}

		fmt.Fprintf(
			w,
			"   %s\t%s\t%s\t%s\n",
			server_utils.ColorBoldMagenta(shell_name),
			server_utils.ColorGrey(fmt.Sprintf("(PID %d)", shell.PID)),
			server_utils.ColorCyan(fmt.Sprintf("(%s)", shell.Cmd.Args[0])),
			attachedStr,
		)
	}
	w.Flush()

	conn.Close()
	return
}
