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
