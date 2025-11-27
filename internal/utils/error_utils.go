package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// FormatOpenAPIError formats OpenAPI errors for better diagnostics
func FormatOpenAPIError(err error, message string) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	if err != nil && strings.Contains(err.Error(), "404") {
		diagnostics.AddError(
			"Unsupported Resource Type",
			fmt.Sprintf("%s - This resource type is not supported", message),
		)
		return diagnostics
	}

	if openAPIErr, ok := err.(interface{ Body() []byte }); ok {
		body := openAPIErr.Body()
		var respBody map[string]interface{}
		if jsonErr := json.Unmarshal(body, &respBody); jsonErr == nil {
			if payload, ok := respBody["payload"].(string); ok {
				diagnostics.AddError(
					message,
					fmt.Sprintf("%v\nDetails: %s", err, payload),
				)
				return diagnostics
			}
		}
	}

	diagnostics.AddError(
		message,
		err.Error(),
	)

	return diagnostics
}
