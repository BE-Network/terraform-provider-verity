package utils

import (
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

// MapFloat64FromAPI converts an API interface{} value to types.Float64
// Handles multiple numeric types that might come from JSON
func MapFloat64FromAPI(apiValue interface{}) types.Float64 {
	if apiValue == nil {
		return types.Float64Null()
	}

	switch v := apiValue.(type) {
	case float32:
		return types.Float64Value(float64(v))
	case float64:
		return types.Float64Value(v)
	case int:
		return types.Float64Value(float64(v))
	case int32:
		return types.Float64Value(float64(v))
	case int64:
		return types.Float64Value(float64(v))
	case string:
		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			return types.Float64Value(floatVal)
		}
	}

	return types.Float64Null()
}

// MapNullableFloat64FromAPI is specifically for nullable fields that might be null in API
func MapNullableFloat64FromAPI(apiValue interface{}) types.Float64 {
	return MapFloat64FromAPI(apiValue)
}
