// Package log provides context-aware and structured logging capabilities.
package main

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
		return FgRed.coloredString(statusCodeString)
	} else if statusCode >= 400 {
		return FgYellow.coloredString(statusCodeString)
	} else if statusCode >= 300 {
		return FgCyan.coloredString(statusCodeString)
	} else if statusCode >= 200 {
		return FgGreen.coloredString(statusCodeString)
	}
	// default value
	return FgWhite.coloredString(statusCodeString)
}

func colorifyRequestMethod(methodName string) string {
	switch methodName {
	case "GET":
		return BgGreen.coloredString(methodName)
	case "POST":
		return BgYellow.coloredString(methodName)
	case "PUT":
		return BgBlue.coloredString(methodName)
	case "DELETE":
		return BgRed.coloredString(methodName)
	case "OPTION":
		return BgCyan.coloredString(methodName)
	case "PATCH":
		return BgMagenta.coloredString(methodName)
	default:
		return BgWhite.coloredString(methodName)
	}
}

func colorifyRequestLatency(latency time.Duration) string {
	if latency < time.Second {
		return FgGreen.coloredString(latency.String())
	} else if latency < time.Second*2 {
		return FgYellow.coloredString(latency.String())
	}
	return FgRed.coloredString(latency.String())
}



func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
