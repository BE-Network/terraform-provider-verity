package utils

import (
	"fmt"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CompareAndSetStringField compares plan vs state and sets API field if changed
func CompareAndSetStringField(plan, state types.String, setter func(*string), hasChanges *bool) {
	if !plan.Equal(state) {
		if !plan.IsNull() {
			setter(openapi.PtrString(plan.ValueString()))
		} else {
			setter(openapi.PtrString(""))
		}
		*hasChanges = true
	}
}

// CompareAndSetBoolField compares plan vs state and sets API field if changed
func CompareAndSetBoolField(plan, state types.Bool, setter func(*bool), hasChanges *bool) {
	if !plan.Equal(state) {
		setter(openapi.PtrBool(plan.ValueBool()))
		*hasChanges = true
	}
}

// CompareAndSetInt64Field compares plan vs state and sets API field if changed (converts to int32)
func CompareAndSetInt64Field(plan, state types.Int64, setter func(*int32), hasChanges *bool) {
	if !plan.Equal(state) {
		val := int32(plan.ValueInt64())
		setter(openapi.PtrInt32(val))
		*hasChanges = true
	}
}

// CompareAndSetNullableInt64Field compares config vs state and sets API nullable field if changed.
// Only processes the field if isConfigured is true (field is explicitly written in HCL).
// If isConfigured is false, the field is omitted from HCL and should not be sent in the PATCH request.
func CompareAndSetNullableInt64Field(configVal, stateVal types.Int64, isConfigured bool, setter func(*openapi.NullableInt32), hasChanges *bool) {
	// Skip if field is not configured in HCL
	if !isConfigured {
		return
	}

	// Only send if value differs from state
	if !configVal.Equal(stateVal) {
		if !configVal.IsNull() {
			val := int32(configVal.ValueInt64())
			nullableVal := *openapi.NewNullableInt32(&val)
			setter(&nullableVal)
		} else {
			nullableVal := *openapi.NewNullableInt32(nil)
			setter(&nullableVal)
		}
		*hasChanges = true
	}
}

// CompareAndSetNullableNumberField compares config vs state and sets API nullable field if changed.
// Uses types.Number (backed by big.Float) to avoid float precision issues.
// Only processes the field if isConfigured is true (field is explicitly written in HCL).
// If isConfigured is false, the field is omitted from HCL and should not be sent in the PATCH request.
func CompareAndSetNullableNumberField(configVal, stateVal types.Number, isConfigured bool, setter func(*openapi.NullableFloat32), hasChanges *bool) {
	// Skip if field is not configured in HCL
	if !isConfigured {
		return
	}

	// Only send if value differs from state
	if !configVal.Equal(stateVal) {
		if !configVal.IsNull() {
			bigVal := configVal.ValueBigFloat()
			float64Val, _ := bigVal.Float64()
			val := float32(float64Val)
			nullableVal := *openapi.NewNullableFloat32(&val)
			setter(&nullableVal)
		} else {
			nullableVal := *openapi.NewNullableFloat32(nil)
			setter(&nullableVal)
		}
		*hasChanges = true
	}
}

// HandleMultipleRefTypesSupported handles ref type logic for "many ref types supported" pattern
// Always sends both fields when either changes
func HandleMultipleRefTypesSupported(
	planBase, stateBase, planRefType, stateRefType types.String,
	baseSetter, refTypeSetter func(*string),
	baseFieldName, refTypeFieldName string,
	hasChanges *bool,
	diags *diag.Diagnostics,
) bool {
	baseChanged := !planBase.Equal(stateBase)
	refTypeChanged := !planRefType.Equal(stateRefType)

	if baseChanged || refTypeChanged {
		// Validate using "many ref types supported" rules
		if !ValidateMultipleRefTypesSupported(diags, planBase, planRefType, baseFieldName, refTypeFieldName) {
			return false
		}

		// Always send both fields when either changes
		baseValue := stateBase
		if baseChanged {
			baseValue = planBase
		}
		if !baseValue.IsNull() && baseValue.ValueString() != "" {
			baseSetter(openapi.PtrString(baseValue.ValueString()))
		} else {
			baseSetter(openapi.PtrString(""))
		}

		refTypeValue := stateRefType
		if refTypeChanged {
			refTypeValue = planRefType
		}
		if !refTypeValue.IsNull() && refTypeValue.ValueString() != "" {
			refTypeSetter(openapi.PtrString(refTypeValue.ValueString()))
		} else {
			refTypeSetter(openapi.PtrString(""))
		}

		*hasChanges = true
	}
	return true
}

// HandleOneRefTypeSupported handles ref type logic for "one ref type supported" pattern
// Behavior differs based on which field changed
func HandleOneRefTypeSupported(
	planBase, stateBase, planRefType, stateRefType types.String,
	baseSetter, refTypeSetter func(*string),
	baseFieldName, refTypeFieldName string,
	hasChanges *bool,
	diags *diag.Diagnostics,
) bool {
	baseChanged := !planBase.Equal(stateBase)
	refTypeChanged := !planRefType.Equal(stateRefType)

	if baseChanged || refTypeChanged {
		// Validate using "one ref type supported" rules
		if !ValidateOneRefTypeSupported(diags, planBase, planRefType, baseFieldName, refTypeFieldName, baseChanged, refTypeChanged) {
			return false
		}

		// Only send the base field if only it changed
		if baseChanged && !refTypeChanged {
			// Just send the base field
			if !planBase.IsNull() && planBase.ValueString() != "" {
				baseSetter(openapi.PtrString(planBase.ValueString()))
			} else {
				baseSetter(openapi.PtrString(""))
			}
			*hasChanges = true
		} else if refTypeChanged {
			// Send both fields
			if !planBase.IsNull() && planBase.ValueString() != "" {
				baseSetter(openapi.PtrString(planBase.ValueString()))
			} else {
				baseSetter(openapi.PtrString(""))
			}

			if !planRefType.IsNull() && planRefType.ValueString() != "" {
				refTypeSetter(openapi.PtrString(planRefType.ValueString()))
			} else {
				refTypeSetter(openapi.PtrString(""))
			}
			*hasChanges = true
		}
	}
	return true
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
