package utils

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MapStringFromAPI(apiValue interface{}) types.String {
	if apiValue == nil {
		return types.StringNull()
	}
	if strVal, ok := apiValue.(string); ok {
		return types.StringValue(strVal)
	}
	return types.StringNull()
}

func MapBoolFromAPI(apiValue interface{}) types.Bool {
	if apiValue == nil {
		return types.BoolNull()
	}
	if boolVal, ok := apiValue.(bool); ok {
		return types.BoolValue(boolVal)
	}
	return types.BoolNull()
}

func MapInt64FromAPI(apiValue interface{}) types.Int64 {
	if apiValue == nil {
		return types.Int64Null()
	}

	switch v := apiValue.(type) {
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case float64:
		return types.Int64Value(int64(v))
	case string:
		if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
			return types.Int64Value(intVal)
		}
	}

	return types.Int64Null()
}

func MapNullableInt64FromAPI(apiValue interface{}) types.Int64 {
	return MapInt64FromAPI(apiValue)
}

func MapNumberFromAPI(apiValue interface{}) types.Number {
	if apiValue == nil {
		return types.NumberNull()
	}

	switch v := apiValue.(type) {
	case float64:

		str := fmt.Sprintf("%g", v)
		if bf, _, err := big.ParseFloat(str, 10, 256, big.ToNearestEven); err == nil {
			return types.NumberValue(bf)
		}
	case int:
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

func MapStringWithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.String {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.StringNull()
	}
	return MapStringFromAPI(data[fieldName])
}

func MapBoolWithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.Bool {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.BoolNull()
	}
	return MapBoolFromAPI(data[fieldName])
}

func MapInt64WithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.Int64 {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.Int64Null()
	}
	return MapInt64FromAPI(data[fieldName])
}

func MapNumberWithMode(data map[string]interface{}, fieldName, resourceType, mode string) types.Number {
	if !FieldAppliesToMode(resourceType, fieldName, mode) {
		return types.NumberNull()
	}
	return MapNumberFromAPI(data[fieldName])
}

func MapStringWithModeNested(data map[string]interface{}, dataKey, resourceType, fieldPath, mode string) types.String {
	if !FieldAppliesToMode(resourceType, fieldPath, mode) {
		return types.StringNull()
	}
	return MapStringFromAPI(data[dataKey])
}

func MapBoolWithModeNested(data map[string]interface{}, dataKey, resourceType, fieldPath, mode string) types.Bool {
	if !FieldAppliesToMode(resourceType, fieldPath, mode) {
		return types.BoolNull()
	}
	return MapBoolFromAPI(data[dataKey])
}

func MapInt64WithModeNested(data map[string]interface{}, dataKey, resourceType, fieldPath, mode string) types.Int64 {
	if !FieldAppliesToMode(resourceType, fieldPath, mode) {
		return types.Int64Null()
	}
	return MapInt64FromAPI(data[dataKey])
}
