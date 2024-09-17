package client_utils

import "github.com/mgutz/ansi"

var (
	ColorRed         = ansi.ColorFunc("red")
	ColorGreen       = ansi.ColorFunc("green")
	ColorGreenBold   = ansi.ColorFunc("green+b")
	ColorYellow      = ansi.ColorFunc("yellow")
	ColorBlue        = ansi.ColorFunc("blue")
	ColorBlueBold    = ansi.ColorFunc("blue+b")
	ColorWhiteBold   = ansi.ColorFunc("white+b")
	ColorGrey        = ansi.ColorFunc("244")
	ColorBoldMagenta = ansi.ColorFunc("magenta+b")
	ColorBoldCyan    = ansi.ColorFunc("cyan+b")
	ColorCyan        = ansi.ColorFunc("cyan")
)
