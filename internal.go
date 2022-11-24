// Package log provides context-aware and structured logging capabilities.
package zlogger

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
)

// ---------------------- INTERNAL --------------------------------------


func (c color) coloredString(value string) string {
	return fmt.Sprintf("\x1b[%d;%dm %s \x1b[0m", c, 1, value)
}

func colorifySatusCode(statusCode int) string {
	var statusCodeString string = strconv.Itoa(statusCode)
	if statusCode >= 500 {
		return fgRed.coloredString(statusCodeString)
	} else if statusCode >= 400 {
		return fgYellow.coloredString(statusCodeString)
	} else if statusCode >= 300 {
		return fgCyan.coloredString(statusCodeString)
	} else if statusCode >= 200 {
		return fgGreen.coloredString(statusCodeString)
	}
	// default value
	return fgWhite.coloredString(statusCodeString)
}

func colorifyRequestMethod(methodName string) string {
	switch methodName {
	case "GET":
		return bgGreen.coloredString(methodName)
	case "POST":
		return bgYellow.coloredString(methodName)
	case "PUT":
		return bgBlue.coloredString(methodName)
	case "DELETE":
		return bgRed.coloredString(methodName)
	case "OPTION":
		return bgCyan.coloredString(methodName)
	case "PATCH":
		return bgMagenta.coloredString(methodName)
	default:
		return bgWhite.coloredString(methodName)
	}
}

func colorifyRequestLatency(latency time.Duration) string {
	if latency < time.Second {
		return fgGreen.coloredString(latency.String())
	} else if latency < time.Second*2 {
		return fgYellow.coloredString(latency.String())
	}
	return fgRed.coloredString(latency.String())
}



func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
