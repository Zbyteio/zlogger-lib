// Package log provides context-aware and structured logging capabilities.
package zlogger

import (
	"unicode"
)

// ---------------------- INTERNAL --------------------------------------




func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
