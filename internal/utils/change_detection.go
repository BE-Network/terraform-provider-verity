package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func ValidateInconsistentFields(diags *diag.Diagnostics, baseField types.String, refTypeField types.String, baseFieldName, refTypeFieldName string) bool {

	if baseField.IsNull() && !refTypeField.IsNull() && refTypeField.ValueString() != "" {
		diags.AddError(
			"Inconsistent fields",
			fmt.Sprintf("You cannot set '%s' to empty while '%s' has a value. Please set both fields together.", baseFieldName, refTypeFieldName),
		)
		return false
	}

	if refTypeField.IsNull() && !baseField.IsNull() && baseField.ValueString() != "" {
		diags.AddError(
			"Inconsistent fields",
			fmt.Sprintf("You cannot set '%s' to empty while '%s' has a value. Please set both fields together.", refTypeFieldName, baseFieldName),
		)
		return false
	}

	return true
}

func AddIneffectiveChangeWarning(diags *diag.Diagnostics, baseField types.String, refTypeField types.String, baseFieldName, refTypeFieldName string) {

	if (baseField.IsNull() || baseField.ValueString() == "") &&
		!refTypeField.IsNull() && refTypeField.ValueString() != "" {
		diags.AddWarning(
			"Ineffective change",
			fmt.Sprintf("Setting '%s' while '%s' is empty won't have any effect. Both fields need to be set together.", refTypeFieldName, baseFieldName),
		)
	}

	if !baseField.IsNull() && baseField.ValueString() != "" &&
		(refTypeField.IsNull() || refTypeField.ValueString() == "") {
		diags.AddWarning(
			"Ineffective change",
			fmt.Sprintf("Setting '%s' while '%s' is empty won't have any effect. Both fields need to be set together.", baseFieldName, refTypeFieldName),
		)
	}
}

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

func ValidateMultipleRefTypesSupported(diags *diag.Diagnostics, baseField, refTypeField types.String, baseFieldName, refTypeFieldName string) bool {

	return ValidateReferenceFields(diags, baseField, refTypeField, baseFieldName, refTypeFieldName)
}

func ValidateOneRefTypeSupported(diags *diag.Diagnostics, baseField, refTypeField types.String, baseFieldName, refTypeFieldName string, baseChanged, refTypeChanged bool) bool {
	if baseChanged && !refTypeChanged {

		if !ValidateMissingReferenceType(diags, baseField, refTypeField, baseFieldName, refTypeFieldName) {
			return false
		}

		if baseField.IsNull() && !refTypeField.IsNull() && refTypeField.ValueString() != "" {
			diags.AddError(
				"Inconsistent fields",
				fmt.Sprintf("You cannot set '%s' to empty while '%s' has a value. Please set both fields together.", baseFieldName, refTypeFieldName),
			)
			return false
		}
	} else if refTypeChanged {

		if !ValidateMissingBaseField(diags, baseField, refTypeField, baseFieldName, refTypeFieldName) {
			return false
		}

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
