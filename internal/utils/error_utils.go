package utils

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FormatOpenAPIError formats OpenAPI errors for better diagnostics
func FormatOpenAPIError(err error, message string) diag.Diagnostics {
	var diagnostics diag.Diagnostics

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

// ValidateMissingReferenceType checks if a base field is non-empty but the reference type is empty
// Returns true if validation passes, false if an error was added to diagnostics
func ValidateMissingReferenceType(diags *diag.Diagnostics, baseField types.String, refTypeField types.String, baseFieldName, refTypeFieldName string) bool {
	if !baseField.IsNull() && baseField.ValueString() != "" &&
		(refTypeField.IsNull() || refTypeField.ValueString() == "") {
		diags.AddError(
			"Missing reference type",
			fmt.Sprintf("When setting '%s' to a non-empty value, you must also specify '%s'. Please check the API documentation for valid values.", baseFieldName, refTypeFieldName),
		)
		return false
	}
	return true
}

// ValidateMissingBaseField checks if a reference type field is non-empty but the base field is empty
// Returns true if validation passes, false if an error was added to diagnostics
func ValidateMissingBaseField(diags *diag.Diagnostics, baseField types.String, refTypeField types.String, baseFieldName, refTypeFieldName string) bool {
	if !refTypeField.IsNull() && refTypeField.ValueString() != "" &&
		(baseField.IsNull() || baseField.ValueString() == "") {
		diags.AddError(
			"Missing base field",
			fmt.Sprintf("When setting '%s' to a non-empty value, you must also specify '%s'. The API requires both fields to be set together.", refTypeFieldName, baseFieldName),
		)
		return false
	}
	return true
}

// ValidateInconsistentFields checks if one field is changing to empty while the other remains non-empty
// Returns true if validation passes, false if an error was added to diagnostics
func ValidateInconsistentFields(diags *diag.Diagnostics, baseField types.String, refTypeField types.String, baseFieldName, refTypeFieldName string) bool {
	// Check if base field is null but ref type has value
	if baseField.IsNull() && !refTypeField.IsNull() && refTypeField.ValueString() != "" {
		diags.AddError(
			"Inconsistent fields",
			fmt.Sprintf("You cannot set '%s' to empty while '%s' has a value. Please set both fields together.", baseFieldName, refTypeFieldName),
		)
		return false
	}

	// Check if ref type is null but base field has value
	if refTypeField.IsNull() && !baseField.IsNull() && baseField.ValueString() != "" {
		diags.AddError(
			"Inconsistent fields",
			fmt.Sprintf("You cannot set '%s' to empty while '%s' has a value. Please set both fields together.", refTypeFieldName, baseFieldName),
		)
		return false
	}

	return true
}

// AddIneffectiveChangeWarning adds warnings for ineffective changes in reference fields
func AddIneffectiveChangeWarning(diags *diag.Diagnostics, baseField types.String, refTypeField types.String, baseFieldName, refTypeFieldName string) {
	// Case 1: Base field is empty but ref type has a value
	if (baseField.IsNull() || baseField.ValueString() == "") &&
		!refTypeField.IsNull() && refTypeField.ValueString() != "" {
		diags.AddWarning(
			"Ineffective change",
			fmt.Sprintf("Setting '%s' while '%s' is empty won't have any effect. Both fields need to be set together.", refTypeFieldName, baseFieldName),
		)
	}

	// Case 2: Base field has a value but ref type is empty
	if !baseField.IsNull() && baseField.ValueString() != "" &&
		(refTypeField.IsNull() || refTypeField.ValueString() == "") {
		diags.AddWarning(
			"Ineffective change",
			fmt.Sprintf("Setting '%s' while '%s' is empty won't have any effect. Both fields need to be set together.", baseFieldName, refTypeFieldName),
		)
	}
}

// ValidateReferenceFields performs all reference field validations in one call
// Returns true if all validations pass, false if any validation failed
func ValidateReferenceFields(diags *diag.Diagnostics, baseField types.String, refTypeField types.String, baseFieldName, refTypeFieldName string) bool {
	if !ValidateMissingReferenceType(diags, baseField, refTypeField, baseFieldName, refTypeFieldName) {
		return false
	}

	if !ValidateMissingBaseField(diags, baseField, refTypeField, baseFieldName, refTypeFieldName) {
		return false
	}

	if !ValidateInconsistentFields(diags, baseField, refTypeField, baseFieldName, refTypeFieldName) {
		return false
	}

	AddIneffectiveChangeWarning(diags, baseField, refTypeField, baseFieldName, refTypeFieldName)

	return true
}

// ValidateMultipleRefTypesSupported performs validation for fields that support multiple ref types
// Always sends both fields when either changes
// Returns true if validation passes, false if there are errors
func ValidateMultipleRefTypesSupported(diags *diag.Diagnostics, baseField, refTypeField types.String, baseFieldName, refTypeFieldName string) bool {
	// Run all validations
	return ValidateReferenceFields(diags, baseField, refTypeField, baseFieldName, refTypeFieldName)
}

// ValidateOneRefTypeSupported performs validation for fields that support only one ref type
// Behavior differs based on which field changed
// Returns true if validation passes, false if there are errors
func ValidateOneRefTypeSupported(diags *diag.Diagnostics, baseField, refTypeField types.String, baseFieldName, refTypeFieldName string, baseChanged, refTypeChanged bool) bool {
	if baseChanged && !refTypeChanged {
		// When base field changes but ref_type doesn't, validate that ref_type is set when base is non-empty
		if !ValidateMissingReferenceType(diags, baseField, refTypeField, baseFieldName, refTypeFieldName) {
			return false
		}

		// Validate consistency when base field is being set to empty
		if baseField.IsNull() && !refTypeField.IsNull() && refTypeField.ValueString() != "" {
			diags.AddError(
				"Inconsistent fields",
				fmt.Sprintf("You cannot set '%s' to empty while '%s' has a value. Please set both fields together.", baseFieldName, refTypeFieldName),
			)
			return false
		}
	} else if refTypeChanged {
		// When ref_type changes (or both change)

		// Validate base field is set when ref_type is non-empty
		if !ValidateMissingBaseField(diags, baseField, refTypeField, baseFieldName, refTypeFieldName) {
			return false
		}

		// Validate ref_type is not being set to empty while base has a value
		if refTypeField.IsNull() && !baseField.IsNull() && baseField.ValueString() != "" {
			diags.AddError(
				"Inconsistent fields",
				fmt.Sprintf("You cannot set '%s' to empty while '%s' has a value. Please set both fields together.", refTypeFieldName, baseFieldName),
			)
			return false
		}

		AddIneffectiveChangeWarning(diags, baseField, refTypeField, baseFieldName, refTypeFieldName)
	}

	return true
}
