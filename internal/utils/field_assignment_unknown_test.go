package utils

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/openapi"
)

func TestSetNullableInt64FieldsOmitsUnknown(t *testing.T) {
	t.Parallel()

	var field openapi.NullableInt64
	SetNullableInt64Fields([]NullableInt64FieldMapping{{
		FieldName:    "mtu",
		APIField:     &field,
		TFValue:      types.Int64Unknown(),
		IsConfigured: true,
	}})

	if field.IsSet() {
		t.Fatalf("an unknown int64 was serialized as %v; it must be omitted so the read after the operation supplies it", *field.Get())
	}
}

func TestSetNullableNumberFieldsOmitsUnknown(t *testing.T) {
	t.Parallel()

	var field openapi.NullableFloat64
	SetNullableNumberFields([]NullableNumberFieldMapping{{
		FieldName:    "weight",
		APIField:     &field,
		TFValue:      types.NumberUnknown(),
		IsConfigured: true,
	}})

	if field.IsSet() {
		t.Fatalf("an unknown number was serialized as %v; it must be omitted so the read after the operation supplies it", *field.Get())
	}
}

func TestSetNullableFieldsStillSendKnownValuesAndExplicitNull(t *testing.T) {
	t.Parallel()

	var known openapi.NullableInt64
	SetNullableInt64Fields([]NullableInt64FieldMapping{{
		FieldName:    "mtu",
		APIField:     &known,
		TFValue:      types.Int64Value(1500),
		IsConfigured: true,
	}})
	if !known.IsSet() || known.Get() == nil || *known.Get() != 1500 {
		t.Fatalf("a known int64 did not reach the request: set=%v value=%v", known.IsSet(), known.Get())
	}

	var cleared openapi.NullableInt64
	SetNullableInt64Fields([]NullableInt64FieldMapping{{
		FieldName:    "mtu",
		APIField:     &cleared,
		TFValue:      types.Int64Null(),
		IsConfigured: true,
	}})
	if !cleared.IsSet() || cleared.Get() != nil {
		t.Fatalf("an explicit null did not reach the request as null: set=%v value=%v", cleared.IsSet(), cleared.Get())
	}

	var number openapi.NullableFloat64
	SetNullableNumberFields([]NullableNumberFieldMapping{{
		FieldName:    "weight",
		APIField:     &number,
		TFValue:      types.NumberValue(big.NewFloat(2.5)),
		IsConfigured: true,
	}})
	if !number.IsSet() || number.Get() == nil || *number.Get() != 2.5 {
		t.Fatalf("a known number did not reach the request: set=%v value=%v", number.IsSet(), number.Get())
	}

	var clearedNumber openapi.NullableFloat64
	SetNullableNumberFields([]NullableNumberFieldMapping{{
		FieldName:    "weight",
		APIField:     &clearedNumber,
		TFValue:      types.NumberNull(),
		IsConfigured: true,
	}})
	if !clearedNumber.IsSet() || clearedNumber.Get() != nil {
		t.Fatalf("an explicit null number did not reach the request as null: set=%v value=%v", clearedNumber.IsSet(), clearedNumber.Get())
	}
}
