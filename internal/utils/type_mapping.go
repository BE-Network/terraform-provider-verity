package utils

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapStringFromAPI converts an API interface{} value to types.String
func MapStringFromAPI(apiValue interface{}) types.String {
	if apiValue == nil {
		return types.StringNull()
	}
	if strVal, ok := apiValue.(string); ok {
		return types.StringValue(strVal)
	}
	return types.StringNull()
}

// MapBoolFromAPI converts an API interface{} value to types.Bool
func MapBoolFromAPI(apiValue interface{}) types.Bool {
	if apiValue == nil {
		return types.BoolNull()
	}
	if boolVal, ok := apiValue.(bool); ok {
		return types.BoolValue(boolVal)
	}
	return types.BoolNull()
}

// MapInt64FromAPI converts an API interface{} value to types.Int64
// Handles multiple numeric types that might come from JSON
func MapInt64FromAPI(apiValue interface{}) types.Int64 {
	if apiValue == nil {
		return types.Int64Null()
	}

	switch v := apiValue.(type) {
	case int:
		return types.Int64Value(int64(v))
	case int32:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case float32:
		return types.Int64Value(int64(v))
	case float64:
		return types.Int64Value(int64(v))
	case string:
		if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
			return types.Int64Value(intVal)
		}
	}

	return types.Int64Null()
}

// MapNullableInt64FromAPI is specifically for nullable fields that might be null in API
func MapNullableInt64FromAPI(apiValue interface{}) types.Int64 {
	return MapInt64FromAPI(apiValue)
}

// MapNumberFromAPI converts an API interface{} value to types.Number
// Uses string parsing to avoid float32/float64 precision issues
// This is the recommended way to handle decimal numbers in Terraform
//
// NOTE: We must parse floats as strings to match how Terraform parses HCL.
// When Terraform reads "0.99" from HCL, it uses string parsing to create an exact big.Float.
// If we use big.NewFloat(float64), we get precision artifacts that cause spurious diffs.
func MapNumberFromAPI(apiValue interface{}) types.Number {
	if apiValue == nil {
		return types.NumberNull()
	}

	switch v := apiValue.(type) {
	case float32:
		// Convert to string first to preserve decimal representation
		// Use %g to get the shortest representation that round-trips
		str := fmt.Sprintf("%g", v)
		if bf, _, err := big.ParseFloat(str, 10, 256, big.ToNearestEven); err == nil {
			return types.NumberValue(bf)
		}
	case float64:
		// Convert to string first to preserve decimal representation
		// Use %g to get the shortest representation that round-trips
		str := fmt.Sprintf("%g", v)
		if bf, _, err := big.ParseFloat(str, 10, 256, big.ToNearestEven); err == nil {
			return types.NumberValue(bf)
		}
	case int:
		bf := big.NewFloat(float64(v))
		return types.NumberValue(bf)
	case int32:
		bf := big.NewFloat(float64(v))
		return types.NumberValue(bf)
	case int64:
		bf := big.NewFloat(float64(v))
		return types.NumberValue(bf)
	case string:
		if bf, _, err := big.ParseFloat(v, 10, 256, big.ToNearestEven); err == nil {
			return types.NumberValue(bf)
		}
	}

	return types.NumberNull()
}

// Mode-aware mapping functions
// These functions check if a field applies to the current mode and return null if not

// MapStringWithMode maps a string field from API data, returning null if field doesn't apply to mode
func MapStringWithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.String {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.StringNull()
	}
	return MapStringFromAPI(data[fieldName])
}

// MapBoolWithMode maps a bool field from API data, returning null if field doesn't apply to mode
func MapBoolWithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.Bool {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.BoolNull()
	}
	return MapBoolFromAPI(data[fieldName])
}

// MapInt64WithMode maps an int64 field from API data, returning null if field doesn't apply to mode
func MapInt64WithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.Int64 {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.Int64Null()
	}
	return MapInt64FromAPI(data[fieldName])
}

// MapNumberWithMode maps a number field from API data, returning null if field doesn't apply to mode
// Uses big.Float to avoid float precision issues
func MapNumberWithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.Number {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.NumberNull()
	}
	return MapNumberFromAPI(data[fieldName])
}

// Nested field mapping functions - for fields inside nested blocks like object_properties
// These take a separate dataKey (for API data lookup) and fieldPath (for mode checking)

// MapStringWithModeNested maps a string field from nested API data
// dataKey is the key in the data map, fieldPath is the full path for mode checking (e.g., "object_properties.group")
func MapStringWithModeNested(data map[string]interface{}, dataKey, resourceType, fieldPath, mode string) types.String {
	if !FieldAppliesToMode(resourceType, fieldPath, mode) {
		return types.StringNull()
	}
	return MapStringFromAPI(data[dataKey])
}

// MapBoolWithModeNested maps a bool field from nested API data
func MapBoolWithModeNested(data map[string]interface{}, dataKey, resourceType, fieldPath, mode string) types.Bool {
	if !FieldAppliesToMode(resourceType, fieldPath, mode) {
		return types.BoolNull()
	}
	return MapBoolFromAPI(data[dataKey])
}

// MapInt64WithModeNested maps an int64 field from nested API data
func MapInt64WithModeNested(data map[string]interface{}, dataKey, resourceType, fieldPath, mode string) types.Int64 {
	if !FieldAppliesToMode(resourceType, fieldPath, mode) {
		return types.Int64Null()
	}
	return MapInt64FromAPI(data[dataKey])
}
