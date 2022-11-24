package zlogger

import "fmt"

// Foreground colors.
const (
	fgBlack color = iota + 30
	fgRed
	fgGreen
	fgYellow
	fgBlue
	fgMagenta
	fgCyan
	fgWhite
)

// Background colors.
const (
	bgBlack color = iota + 40
	bgRed
	bgGreen
	bgYellow
	bgBlue
	bgMagenta
	bgCyan
	bgWhite
)

type color uint8


func (c color) coloredString(value string) string {
	return fmt.Sprintf("\x1b[%d;%dm %s \x1b[0m", c, 1, value)
}