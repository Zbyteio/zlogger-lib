package zlogger

/* DOCS -
for dev mode only
not accesable outside the package
*/

import (
	"fmt"
	"time"
)

type color struct {
	fgBlack   uint8
	fgRed     uint8
	fgGreen   uint8
	fgYellow  uint8
	fgBlue    uint8
	fgMagenta uint8
	fgCyan    uint8
	fgWhite   uint8

	bgBlack   uint8
	bgRed     uint8
	bgGreen   uint8
	bgYellow  uint8
	bgBlue    uint8
	bgMagenta uint8
	bgCyan    uint8
	bgWhite   uint8
}

type icolor interface {
	colorfgBlack(value interface{}) string
	colorfgRed(value interface{}) string
	colorfgGreen(value interface{}) string
	colorfgYellow(value interface{}) string
	colorfgBlue(value interface{}) string
	colorfgMagenta(value interface{}) string
	colorfgCyan(value interface{}) string
	colorfgWhite(value interface{}) string

	colorbgBlack(value interface{}) string
	colorbgRed(value interface{}) string
	colorbgGreen(value interface{}) string
	colorbgYellow(value interface{}) string
	colorbgBlue(value interface{}) string
	colorbgMagenta(value interface{}) string
	colorbgCyan(value interface{}) string
	colorbgWhite(value interface{}) string
}

var colorPallet icolor

func init() {
	colorPallet = color{
		fgBlack:   30,
		fgRed:     31,
		fgGreen:   32,
		fgYellow:  33,
		fgBlue:    34,
		fgMagenta: 35,
		fgCyan:    36,
		fgWhite:   37,

		bgBlack:   40,
		bgRed:     41,
		bgGreen:   42,
		bgYellow:  43,
		bgBlue:    44,
		bgMagenta: 45,
		bgCyan:    46,
		bgWhite:   47,
	}

}

func (c color) colorfgBlack(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgBlack, 1, value)
}
func (c color) colorfgRed(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgRed, 1, value)
}
func (c color) colorfgGreen(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgGreen, 1, value)
}
func (c color) colorfgYellow(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgYellow, 1, value)
}
func (c color) colorfgBlue(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgBlue, 1, value)
}
func (c color) colorfgMagenta(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgMagenta, 1, value)
}
func (c color) colorfgCyan(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgCyan, 1, value)
}
func (c color) colorfgWhite(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm%v\x1b[0m", c.fgWhite, 1, value)
}

func (c color) colorbgBlack(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgBlack, 1, value)
}
func (c color) colorbgRed(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgRed, 1, value)
}
func (c color) colorbgGreen(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgGreen, 1, value)
}
func (c color) colorbgYellow(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgYellow, 1, value)
}
func (c color) colorbgBlue(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgBlue, 1, value)
}
func (c color) colorbgMagenta(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgMagenta, 1, value)
}
func (c color) colorbgCyan(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgCyan, 1, value)
}
func (c color) colorbgWhite(value interface{}) string {
	return fmt.Sprintf("\x1b[%d;%dm %v \x1b[0m", c.bgWhite, 1, value)
}

func colorifySatusCode(statusCode int) string {
	if statusCode >= 500 {
		return colorPallet.colorfgRed(statusCode)
	} else if statusCode >= 400 {
		return colorPallet.colorfgYellow(statusCode)
	} else if statusCode >= 300 {
		return colorPallet.colorfgCyan(statusCode)
	} else if statusCode >= 200 {
		return colorPallet.colorfgGreen(statusCode)
	}
	// default value
	return colorPallet.colorfgWhite(statusCode)
}

func colorifyRequestMethod(methodName string) string {
	switch methodName {
	case "GET":
		return colorPallet.colorbgGreen(methodName)
	case "POST":
		return colorPallet.colorbgYellow(methodName)
	case "PUT":
		return colorPallet.colorbgBlue(methodName)
	case "DELETE":
		return colorPallet.colorbgRed(methodName)
	case "OPTION":
		return colorPallet.colorbgCyan(methodName)
	case "PATCH":
		return colorPallet.colorbgMagenta(methodName)
	default:
		return colorPallet.colorbgWhite(methodName)
	}
}

func colorifyRequestLatency(latency time.Duration) string {
	if latency < time.Second {
		return colorPallet.colorfgGreen(latency.String())
	} else if latency < time.Second*2 {
		return colorPallet.colorfgYellow(latency.String())
	}
	return colorPallet.colorfgRed(latency.String())
}

func colorifySqlLatency(latency time.Duration, threshold time.Duration) string {
	if latency < threshold {
		return colorPallet.colorfgGreen(latency.String())
	} else if latency < threshold*2 {
		return colorPallet.colorfgYellow(latency.String())
	}
	return colorPallet.colorfgRed(latency.String())
}

func colorifyRequestError(errorMsg string) string {
	return colorPallet.colorfgRed(errorMsg)
}