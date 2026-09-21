package utils

import (
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StringFieldMapping struct {
	FieldName string
	APIField  **string
	TFValue   types.String
}

type BoolFieldMapping struct {
	FieldName string
	APIField  **bool
	TFValue   types.Bool
}

type Int64FieldMapping struct {
	FieldName string
	APIField  **int64
	TFValue   types.Int64
}

type NullableInt64FieldMapping struct {
	FieldName    string
	APIField     *openapi.NullableInt64
	TFValue      types.Int64
	IsConfigured bool
}

type NullableNumberFieldMapping struct {
	FieldName    string
	APIField     *openapi.NullableFloat64
	TFValue      types.Number
	IsConfigured bool
}

func SetStringFields(fields []StringFieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() && !field.TFValue.IsUnknown() {
			*field.APIField = openapi.PtrString(field.TFValue.ValueString())
		}
	}
}

func SetBoolFields(fields []BoolFieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() && !field.TFValue.IsUnknown() {
			*field.APIField = openapi.PtrBool(field.TFValue.ValueBool())
		}
	}
}

func SetInt64Fields(fields []Int64FieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() && !field.TFValue.IsUnknown() {
			val := field.TFValue.ValueInt64()
			*field.APIField = openapi.PtrInt64(val)
		}
	}
}

func SetNullableInt64Fields(fields []NullableInt64FieldMapping) {
	for _, field := range fields {

		if !field.IsConfigured {
			continue
		}

		if field.TFValue.IsUnknown() {
			continue
		}

		if !field.TFValue.IsNull() {

			val := field.TFValue.ValueInt64()
			*field.APIField = *openapi.NewNullableInt64(&val)
		} else {

			*field.APIField = *openapi.NewNullableInt64(nil)
		}
	}
}

func SetNullableNumberFields(fields []NullableNumberFieldMapping) {
	for _, field := range fields {

		if !field.IsConfigured {
			continue
		}

		if field.TFValue.IsUnknown() {
			continue
		}

		if !field.TFValue.IsNull() {

			bigVal := field.TFValue.ValueBigFloat()
			float64Val, _ := bigVal.Float64()
			val := float64Val
			*field.APIField = *openapi.NewNullableFloat64(&val)
		} else {

			*field.APIField = *openapi.NewNullableFloat64(nil)
		}
	}
}
