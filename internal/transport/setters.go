package transport

import (
	"fmt"
	"strconv"

	"terraform-provider-verity/openapi"
)

// ResourceValueAdapter converts the codec's canonical object into the typed
// value the bulk manager asserts for one resource.
//
// The generated adapters implement it. The interface lives here, beside the
// setters they call, so the generated file carries nothing but the per-resource
// mapping.
type ResourceValueAdapter interface {
	ResourceValue(object WireObject) (interface{}, error)
}

// The setters below are the whole vocabulary the generated adapters use. Each
// takes one canonical value and fills one field of a generated SDK struct,
// which is where the two type systems meet: a wire null becomes an explicit API
// null where the field can hold one, and is refused where it cannot.

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

// A nullable field is the one place an explicit null is a value rather than an
// error: it is how this API clears a number.
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

// int64Of and float64Of accept the two numeric kinds the codec produces. A
// decimal keeps its lexical form until here, so the conversion happens once, at
// the boundary, rather than being carried as a float through the codec.
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

// wireObject unwraps a singleton object for a generated nested adapter. An
// object is never sent as null — an absent block is omitted instead — so a null
// here means the codec and the adapter disagree.
func wireObject(value WireValue) (WireObject, error) {
	if value.Kind() != ValueKindObject {
		return nil, fmt.Errorf("expected object, got %s", value.Kind())
	}
	return value.ObjectValue(), nil
}

// wireList unwraps an indexed collection for a generated list adapter. Like an
// object, a list is omitted when absent rather than sent as null.
func wireList(value WireValue) ([]WireValue, error) {
	if value.Kind() != ValueKindList {
		return nil, fmt.Errorf("expected list, got %s", value.Kind())
	}
	return value.ListValue(), nil
}
