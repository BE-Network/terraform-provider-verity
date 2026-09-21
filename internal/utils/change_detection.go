package utils

import (
	"fmt"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RefTypeSupportMode string

const (
	RefTypeSupportOne      RefTypeSupportMode = "one"
	RefTypeSupportMultiple RefTypeSupportMode = "multiple"
)

type RefTypeFieldMapping struct {
	FieldName        string
	RefTypeFieldName string
	APIField         **string
	RefTypeAPIField  **string
	TFValue          types.String
	RefTypeTFValue   types.String
}

type RefTypeFieldWithComparison struct {
	FieldName         string
	RefTypeFieldName  string
	APIField          func(*string)
	RefTypeAPIField   func(*string)
	PlanValue         types.String
	StateValue        types.String
	PlanRefTypeValue  types.String
	StateRefTypeValue types.String
	SupportMode       RefTypeSupportMode
}

func SetRefTypeFields(fields []RefTypeFieldMapping) {
	for _, field := range fields {
		SetStringFields([]StringFieldMapping{
			{FieldName: field.FieldName, APIField: field.APIField, TFValue: field.TFValue},
			{FieldName: field.RefTypeFieldName, APIField: field.RefTypeAPIField, TFValue: field.RefTypeTFValue},
		})
	}
}

func CompareAndSetRefTypeFields(fields []RefTypeFieldWithComparison, hasChanges *bool, diags *diag.Diagnostics) bool {
	for _, field := range fields {
		if !applyRefTypeFieldChange(
			field.PlanValue, field.StateValue,
			field.PlanRefTypeValue, field.StateRefTypeValue,
			field.APIField, field.RefTypeAPIField,
			field.FieldName, field.RefTypeFieldName,
			field.SupportMode,
			hasChanges, diags,
		) {
			return false
		}
	}
	return true
}

func applyRefTypeFieldChange(
	planBase, stateBase, planRefType, stateRefType types.String,
	baseSetter, refTypeSetter func(*string),
	baseFieldName, refTypeFieldName string,
	supportMode RefTypeSupportMode,
	hasChanges *bool,
	diags *diag.Diagnostics,
) bool {
	baseChanged := !planBase.Equal(stateBase)
	refTypeChanged := !planRefType.Equal(stateRefType)

	if !baseChanged && !refTypeChanged {
		return true
	}

	switch supportMode {
	case RefTypeSupportMultiple:
		if !ValidateMultipleRefTypesSupported(diags, planBase, planRefType, baseFieldName, refTypeFieldName) {
			return false
		}

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
		return true

	default:
		if !ValidateOneRefTypeSupported(diags, planBase, planRefType, baseFieldName, refTypeFieldName, baseChanged, refTypeChanged) {
			return false
		}

		if baseChanged && !refTypeChanged {
			if !planBase.IsNull() && planBase.ValueString() != "" {
				baseSetter(openapi.PtrString(planBase.ValueString()))
			} else {
				baseSetter(openapi.PtrString(""))
			}
			*hasChanges = true
			return true
		}

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
		return true
	}
}

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

func CompareAndSetBoolField(plan, state types.Bool, setter func(*bool), hasChanges *bool) {
	if !plan.Equal(state) {
		setter(openapi.PtrBool(plan.ValueBool()))
		*hasChanges = true
	}
}

func CompareAndSetInt64Field(plan, state types.Int64, setter func(*int64), hasChanges *bool) {
	if !plan.Equal(state) {
		val := plan.ValueInt64()
		setter(openapi.PtrInt64(val))
		*hasChanges = true
	}
}

func CompareAndSetNullableInt64Field(configVal, stateVal types.Int64, isConfigured bool, setter func(*openapi.NullableInt64), hasChanges *bool) {

	if !isConfigured {
		return
	}

	if !configVal.Equal(stateVal) {
		if !configVal.IsNull() {
			val := configVal.ValueInt64()
			nullableVal := *openapi.NewNullableInt64(&val)
			setter(&nullableVal)
		} else {
			nullableVal := *openapi.NewNullableInt64(nil)
			setter(&nullableVal)
		}
		*hasChanges = true
	}
}

func CompareAndSetNullableNumberField(configVal, stateVal types.Number, isConfigured bool, setter func(*openapi.NullableFloat64), hasChanges *bool) {

	if !isConfigured {
		return
	}

	if !configVal.Equal(stateVal) {
		if !configVal.IsNull() {
			bigVal := configVal.ValueBigFloat()
			float64Val, _ := bigVal.Float64()
			val := float64Val
			nullableVal := *openapi.NewNullableFloat64(&val)
			setter(&nullableVal)
		} else {
			nullableVal := *openapi.NewNullableFloat64(nil)
			setter(&nullableVal)
		}
		*hasChanges = true
	}
}

func HandleMultipleRefTypesSupported(
	planBase, stateBase, planRefType, stateRefType types.String,
	baseSetter, refTypeSetter func(*string),
	baseFieldName, refTypeFieldName string,
	hasChanges *bool,
	diags *diag.Diagnostics,
) bool {
	return applyRefTypeFieldChange(
		planBase, stateBase, planRefType, stateRefType,
		baseSetter, refTypeSetter,
		baseFieldName, refTypeFieldName,
		RefTypeSupportMultiple,
		hasChanges, diags,
	)
}

func HandleOneRefTypeSupported(
	planBase, stateBase, planRefType, stateRefType types.String,
	baseSetter, refTypeSetter func(*string),
	baseFieldName, refTypeFieldName string,
	hasChanges *bool,
	diags *diag.Diagnostics,
) bool {
	return applyRefTypeFieldChange(
		planBase, stateBase, planRefType, stateRefType,
		baseSetter, refTypeSetter,
		baseFieldName, refTypeFieldName,
		RefTypeSupportOne,
		hasChanges, diags,
	)
}

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
