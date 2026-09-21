package transport

import (
	"fmt"
	"strconv"

	"terraform-provider-verity/openapi"
)

type ResourceValueAdapter interface {
	ResourceValue(object WireObject) (interface{}, error)
}

func wireStringPtr(value WireValue, target **string) error {
	if value.Kind() == ValueKindNull {
		return fmt.Errorf("explicit null, but this field cannot represent null")
	}
	if value.Kind() != ValueKindString {
		return fmt.Errorf("expected string, got %s", value.Kind())
	}
	result := value.StringValue()
	*target = &result
	return nil
}

func wireBoolPtr(value WireValue, target **bool) error {
	if value.Kind() == ValueKindNull {
		return fmt.Errorf("explicit null, but this field cannot represent null")
	}
	if value.Kind() != ValueKindBool {
		return fmt.Errorf("expected bool, got %s", value.Kind())
	}
	result := value.BoolValue()
	*target = &result
	return nil
}

func wireInt64Ptr(value WireValue, target **int64) error {
	if value.Kind() == ValueKindNull {
		return fmt.Errorf("explicit null, but this field cannot represent null")
	}
	result, err := int64Of(value)
	if err != nil {
		return err
	}
	*target = &result
	return nil
}

func wireInt32Ptr(value WireValue, target **int32) error {
	if value.Kind() == ValueKindNull {
		return fmt.Errorf("explicit null, but this field cannot represent null")
	}
	result, err := int32Of(value)
	if err != nil {
		return err
	}
	*target = &result
	return nil
}

func wireFloat64Ptr(value WireValue, target **float64) error {
	if value.Kind() == ValueKindNull {
		return fmt.Errorf("explicit null, but this field cannot represent null")
	}
	result, err := float64Of(value)
	if err != nil {
		return err
	}
	*target = &result
	return nil
}

func wireFloat32Ptr(value WireValue, target **float32) error {
	if value.Kind() == ValueKindNull {
		return fmt.Errorf("explicit null, but this field cannot represent null")
	}
	result, err := float32Of(value)
	if err != nil {
		return err
	}
	*target = &result
	return nil
}

func wireNullableInt64(value WireValue, target *openapi.NullableInt64) error {
	if value.Kind() == ValueKindNull {
		*target = *openapi.NewNullableInt64(nil)
		return nil
	}
	result, err := int64Of(value)
	if err != nil {
		return err
	}
	*target = *openapi.NewNullableInt64(&result)
	return nil
}

func wireNullableInt32(value WireValue, target *openapi.NullableInt32) error {
	if value.Kind() == ValueKindNull {
		*target = *openapi.NewNullableInt32(nil)
		return nil
	}
	result, err := int32Of(value)
	if err != nil {
		return err
	}
	*target = *openapi.NewNullableInt32(&result)
	return nil
}

func wireNullableFloat64(value WireValue, target *openapi.NullableFloat64) error {
	if value.Kind() == ValueKindNull {
		*target = *openapi.NewNullableFloat64(nil)
		return nil
	}
	result, err := float64Of(value)
	if err != nil {
		return err
	}
	*target = *openapi.NewNullableFloat64(&result)
	return nil
}

func wireNullableFloat32(value WireValue, target *openapi.NullableFloat32) error {
	if value.Kind() == ValueKindNull {
		*target = *openapi.NewNullableFloat32(nil)
		return nil
	}
	result, err := float32Of(value)
	if err != nil {
		return err
	}
	*target = *openapi.NewNullableFloat32(&result)
	return nil
}

func wireNullableString(value WireValue, target *openapi.NullableString) error {
	if value.Kind() == ValueKindNull {
		*target = *openapi.NewNullableString(nil)
		return nil
	}
	if value.Kind() != ValueKindString {
		return fmt.Errorf("expected string, got %s", value.Kind())
	}
	result := value.StringValue()
	*target = *openapi.NewNullableString(&result)
	return nil
}

func int64Of(value WireValue) (int64, error) {
	switch value.Kind() {
	case ValueKindInt64:
		return value.Int64Value(), nil
	case ValueKindDecimal:
		parsed, err := strconv.ParseInt(value.DecimalValue(), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%q is not an integer", value.DecimalValue())
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("expected integer, got %s", value.Kind())
	}
}

func int32Of(value WireValue) (int32, error) {
	result, err := int64Of(value)
	if err != nil {
		return 0, err
	}
	if result < -1<<31 || result > 1<<31-1 {
		return 0, fmt.Errorf("%d is outside the int32 range", result)
	}
	return int32(result), nil
}

func float64Of(value WireValue) (float64, error) {
	switch value.Kind() {
	case ValueKindInt64:
		return float64(value.Int64Value()), nil
	case ValueKindDecimal:
		parsed, err := strconv.ParseFloat(value.DecimalValue(), 64)
		if err != nil {
			return 0, fmt.Errorf("%q is not a number", value.DecimalValue())
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("expected number, got %s", value.Kind())
	}
}

func float32Of(value WireValue) (float32, error) {
	switch value.Kind() {
	case ValueKindInt64:
		return float32(value.Int64Value()), nil
	case ValueKindDecimal:
		parsed, err := strconv.ParseFloat(value.DecimalValue(), 32)
		if err != nil {
			return 0, fmt.Errorf("%q is not a float32", value.DecimalValue())
		}
		return float32(parsed), nil
	default:
		return 0, fmt.Errorf("expected number, got %s", value.Kind())
	}
}

func wireObject(value WireValue) (WireObject, error) {
	if value.Kind() != ValueKindObject {
		return nil, fmt.Errorf("expected object, got %s", value.Kind())
	}
	return value.ObjectValue(), nil
}

func wireList(value WireValue) ([]WireValue, error) {
	if value.Kind() != ValueKindList {
		return nil, fmt.Errorf("expected list, got %s", value.Kind())
	}
	return value.ListValue(), nil
}

func wireEmptyObject(value WireValue, target *map[string]interface{}) error {
	members, err := wireObject(value)
	if err != nil {
		return err
	}
	if len(members) != 0 {
		return fmt.Errorf("an object with no properties was given %d members", len(members))
	}
	*target = map[string]interface{}{}
	return nil
}
