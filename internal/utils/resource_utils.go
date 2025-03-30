package utils

import (
	"regexp"
	"unicode"
)

// SanitizeResourceName converts a name to a valid Terraform resource ID
// It handles:
// 1. Names starting with numbers (adds underscore prefix)
// 2. Names starting with hyphens (replaces hyphen with underscore)
// 3. Special characters (converts to underscores)
// Example: "123resource" -> "_123resource"
// Example: "My Resource" -> "My_Resource"
// Example: "-g" -> "_g"
func SanitizeResourceName(name string) string {
	if len(name) > 0 {
		firstChar := rune(name[0])
		if unicode.IsDigit(firstChar) {
			// Add underscore prefix if name starts with a number
			name = "_" + name
		} else if firstChar == '-' {
			// Replace hyphen with underscore
			name = "_" + name[1:]
		}
	}

	// Replace any character that's not alphanumeric, underscore, or hyphen with underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return reg.ReplaceAllString(name, "_")
}
