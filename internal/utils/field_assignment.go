package utils

import (
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringFieldMapping represents a mapping between a Terraform string field and an API field
type StringFieldMapping struct {
	FieldName string
	APIField  **string
	TFValue   types.String
}

// BoolFieldMapping represents a mapping between a Terraform bool field and an API field
type BoolFieldMapping struct {
	FieldName string
	APIField  **bool
	TFValue   types.Bool
}

// Int64FieldMapping represents a mapping between a Terraform int64 field and an API int32 field
type Int64FieldMapping struct {
	FieldName string
	APIField  **int32
	TFValue   types.Int64
}

// NullableInt64FieldMapping represents a mapping for nullable int64 fields
type NullableInt64FieldMapping struct {
	FieldName string
	APIField  *openapi.NullableInt32
	TFValue   types.Int64
}

// Float64FieldMapping represents a mapping between a Terraform float64 field and an API float32 field
type Float64FieldMapping struct {
	FieldName string
	APIField  **float32
	TFValue   types.Float64
}

// NullableFloat64FieldMapping represents a mapping for nullable float64 fields
type NullableFloat64FieldMapping struct {
	FieldName string
	APIField  *openapi.NullableFloat32
	TFValue   types.Float64
}

// SetStringFields processes a slice of string field mappings
func SetStringFields(fields []StringFieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() {
			*field.APIField = openapi.PtrString(field.TFValue.ValueString())
		}
	}
}

// SetBoolFields processes a slice of boolean field mappings
func SetBoolFields(fields []BoolFieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() {
			*field.APIField = openapi.PtrBool(field.TFValue.ValueBool())
		}
	}
}

// SetInt64Fields processes a slice of int64 field mappings (API uses int32)
func SetInt64Fields(fields []Int64FieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() {
			val := int32(field.TFValue.ValueInt64())
			*field.APIField = openapi.PtrInt32(val)
		}
	}
}

// SetNullableInt64Fields processes a slice of nullable int64 field mappings
func SetNullableInt64Fields(fields []NullableInt64FieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() {
			val := int32(field.TFValue.ValueInt64())
			*field.APIField = *openapi.NewNullableInt32(&val)
		} else {
			*field.APIField = *openapi.NewNullableInt32(nil)
		}
	}
}

// SetFloat64Fields processes a slice of float64 field mappings (API uses float32)
func SetFloat64Fields(fields []Float64FieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() {
			val := float32(field.TFValue.ValueFloat64())
			*field.APIField = openapi.PtrFloat32(val)
		}
	}
}

// SetNullableFloat64Fields processes a slice of nullable float64 field mappings
func SetNullableFloat64Fields(fields []NullableFloat64FieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() {
			val := float32(field.TFValue.ValueFloat64())
			*field.APIField = *openapi.NewNullableFloat32(&val)
		} else {
			*field.APIField = *openapi.NewNullableFloat32(nil)
		}
	}
}
