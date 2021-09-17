package parser

import (
	"strings"
	"unicode"
)

var sceneheadings = [2]string{
	"INT.",
	"EXT.",
}

func validSceneHeading(contents string) bool {
	if !isUpper(contents) {
		return false
	}

	for _, heading := range sceneheadings {
		if x := strings.Index(contents, heading); x == 0 {
			return true
		}
	}

	// Length check to avoid "." false positives
	if x := strings.Index(contents, "."); x == 0 && len(contents) > 1 {
		return true
	}

	return false
}

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func isAlphaNumeric(s string) bool {
	for _, r := range s {
		alphanumeric := unicode.IsLetter(r) || unicode.IsDigit(r)
		if !alphanumeric && r != ' ' {
			return false
		}
	}

	return true
}
