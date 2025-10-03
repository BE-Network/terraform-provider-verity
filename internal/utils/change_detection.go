package utils

import (
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

// CompareAndSetNullableInt64Field compares plan vs state and sets API nullable field if changed
func CompareAndSetNullableInt64Field(plan, state types.Int64, setter func(*openapi.NullableInt32), hasChanges *bool) {
	if !plan.Equal(state) {
		if !plan.IsNull() {
			val := int32(plan.ValueInt64())
			nullableVal := *openapi.NewNullableInt32(&val)
			setter(&nullableVal)
		} else {
			nullableVal := *openapi.NewNullableInt32(nil)
			setter(&nullableVal)
		}
		*hasChanges = true
	}
}

// CompareAndSetFloat64Field compares plan vs state and sets API field if changed (converts to float32)
func CompareAndSetFloat64Field(plan, state types.Float64, setter func(*float32), hasChanges *bool) {
	if !plan.Equal(state) {
		val := float32(plan.ValueFloat64())
		setter(openapi.PtrFloat32(val))
		*hasChanges = true
	}
}

// CompareAndSetNullableFloat64Field compares plan vs state and sets API nullable field if changed
func CompareAndSetNullableFloat64Field(plan, state types.Float64, setter func(*openapi.NullableFloat32), hasChanges *bool) {
	if !plan.Equal(state) {
		if !plan.IsNull() {
			val := float32(plan.ValueFloat64())
			nullableVal := *openapi.NewNullableFloat32(&val)
			setter(&nullableVal)
		} else {
			nullableVal := *openapi.NewNullableFloat32(nil)
			setter(&nullableVal)
		}
		*hasChanges = true
	}
}

// CompareAndSetStringFieldWithEmpty compares and sets string field, using empty string for null values
func CompareAndSetStringFieldWithEmpty(plan, state types.String, setter func(*string), hasChanges *bool) {
	if !plan.Equal(state) {
		if !plan.IsNull() {
			setter(openapi.PtrString(plan.ValueString()))
		} else {
			setter(openapi.PtrString(""))
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
