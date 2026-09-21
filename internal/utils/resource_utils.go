package utils

import (
	"regexp"
	"unicode"
)

func SanitizeResourceName(name string) string {
	if len(name) > 0 {
		firstChar := rune(name[0])
		if unicode.IsDigit(firstChar) {

			name = "_" + name
		} else if firstChar == '-' {

			name = "_" + name[1:]
		}
	}

	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return reg.ReplaceAllString(name, "_")
}
