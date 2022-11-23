package zlogger

// Foreground colors.
const (
	FgBlack color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Background colors.
const (
	BgBlack color = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

type color uint8

const (
	Production  LogEnvironment = "production"
	Local       LogEnvironment = "local"
)

type LogEnvironment string