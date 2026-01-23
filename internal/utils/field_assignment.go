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

// NullableInt64FieldMapping represents a mapping for nullable int64 fields.
// IsConfigured should be set using ConfiguredAttributes.IsConfigured() from HCL parsing.
type NullableInt64FieldMapping struct {
	FieldName    string
	APIField     *openapi.NullableInt32
	TFValue      types.Int64 // From plan or config - used to get the actual value
	IsConfigured bool        // From HCL parsing - true if key exists in .tf file
}

// NullableNumberFieldMapping represents a mapping for nullable number fields (using big.Float).
// This avoids float32/float64 precision issues by using types.Number backed by big.Float.
// IsConfigured should be set using ConfiguredAttributes.IsConfigured() from HCL parsing.
type NullableNumberFieldMapping struct {
	FieldName    string
	APIField     *openapi.NullableFloat32
	TFValue      types.Number // From plan or config - uses big.Float for precision
	IsConfigured bool         // From HCL parsing - true if key exists in .tf file
}

// SetStringFields processes a slice of string field mappings.
// Skips fields that are null or unknown (known after apply) - these should not be sent to the API.
func SetStringFields(fields []StringFieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() && !field.TFValue.IsUnknown() {
			*field.APIField = openapi.PtrString(field.TFValue.ValueString())
		}
	}
}

// SetBoolFields processes a slice of boolean field mappings.
// Skips fields that are null or unknown (known after apply) - these should not be sent to the API.
func SetBoolFields(fields []BoolFieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() && !field.TFValue.IsUnknown() {
			*field.APIField = openapi.PtrBool(field.TFValue.ValueBool())
		}
	}
}

// SetInt64Fields processes a slice of int64 field mappings (API uses int32).
// Skips fields that are null or unknown (known after apply) - these should not be sent to the API.
func SetInt64Fields(fields []Int64FieldMapping) {
	for _, field := range fields {
		if !field.TFValue.IsNull() && !field.TFValue.IsUnknown() {
			val := int32(field.TFValue.ValueInt64())
			*field.APIField = openapi.PtrInt32(val)
		}
	}
}

// SetNullableInt64Fields processes a slice of nullable int64 field mappings.
// Fields where IsConfigured=false are skipped entirely (not included in API request).
// Fields where IsConfigured=true AND TFValue.IsNull()=true → send explicit null to API.
// Fields where IsConfigured=true AND TFValue.IsNull()=false → send the value to API.
func SetNullableInt64Fields(fields []NullableInt64FieldMapping) {
	for _, field := range fields {
		// Skip fields not explicitly written in the .tf file
		if !field.IsConfigured {
			continue
		}

		if !field.TFValue.IsNull() {
			// Explicit value
			val := int32(field.TFValue.ValueInt64())
			*field.APIField = *openapi.NewNullableInt32(&val)
		} else {
			// Explicit null
			*field.APIField = *openapi.NewNullableInt32(nil)
		}
	}
}

// SetNullableNumberFields processes a slice of nullable number field mappings.
// Uses types.Number (backed by big.Float) to avoid float precision issues.
// Fields where IsConfigured=false are skipped entirely (not included in API request).
// Fields where IsConfigured=true AND TFValue.IsNull()=true → send explicit null to API.
// Fields where IsConfigured=true AND TFValue.IsNull()=false → send the value to API.
func SetNullableNumberFields(fields []NullableNumberFieldMapping) {
	for _, field := range fields {
		// Skip fields not explicitly written in the .tf file
		if !field.IsConfigured {
			continue
		}

		if !field.TFValue.IsNull() {
			// Explicit value - convert big.Float to float32 for API
			bigVal := field.TFValue.ValueBigFloat()
			float64Val, _ := bigVal.Float64()
			val := float32(float64Val)
			*field.APIField = *openapi.NewNullableFloat32(&val)
		} else {
			// Explicit null
			*field.APIField = *openapi.NewNullableFloat32(nil)
		}
	}
}
