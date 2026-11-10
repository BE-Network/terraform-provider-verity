package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const maxAPIErrorLines = 20

var htmlLineBreak = regexp.MustCompile(`(?i)<br\s*/?>|</br>`)

func APIErrorMessage(err error) string {
	withBody, ok := err.(interface{ Body() []byte })
	if !ok {
		return ""
	}
	var body map[string]interface{}
	if json.Unmarshal(withBody.Body(), &body) != nil {
		return ""
	}
	for _, key := range []string{"message", "Error", "error", "payload"} {
		if text, ok := body[key].(string); ok && strings.TrimSpace(text) != "" {
			return strings.TrimSpace(htmlLineBreak.ReplaceAllString(text, "\n"))
		}
	}
	return ""
}

func truncateLines(text string, limit int) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= limit {
		return text
	}
	return strings.Join(lines[:limit], "\n") + fmt.Sprintf("\n… and %d more lines", len(lines)-limit)
}

func FormatOpenAPIError(err error, message string) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	if err != nil && strings.Contains(err.Error(), "404") {
		diagnostics.AddError(
			"Unsupported Resource Type",
			fmt.Sprintf("%s - This resource type is not supported", message),
		)
		return diagnostics
	}

	if reason := APIErrorMessage(err); reason != "" {
		diagnostics.AddError(
			message,
			fmt.Sprintf("%v\nDetails: %s", err, truncateLines(reason, maxAPIErrorLines)),
		)
		return diagnostics
	}

	diagnostics.AddError(
		message,
		err.Error(),
	)

	return diagnostics
}
