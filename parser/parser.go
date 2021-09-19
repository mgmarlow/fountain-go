package parser

import (
	"unicode"
)

func isAlphaNumeric(s string) bool {
	for _, r := range s {
		alphanumeric := unicode.IsLetter(r) || unicode.IsDigit(r)
		if !alphanumeric && r != ' ' {
			return false
		}
	}

	return true
}
