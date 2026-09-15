package utils

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/openapi"
)

// The nullable setters distinguish three states, and an unknown is the third.
//
// A nullable field carries an explicit null on this API, so its setter cannot use
// the "is it null?" test the non-nullable setters use to mean "was it left out?" —
// null is a value it has to send. It gates on IsConfigured instead, which reads
// the .tf file. That leaves unknown unaccounted for: it is neither null nor
// absent from configuration, so it fell through to the value branch, where
// ValueInt64 reports zero for an unknown and ValueBigFloat reports zero likewise.
// The request then stored a zero the configuration never asked for.
//
// SetStringFields, SetBoolFields and SetInt64Fields all skip an unknown. These
// tests hold the nullable pair to the same contract: an unknown is omitted and
// left to the read that follows the operation, which is what the registry records
// for these fields as unknown_plan: omit_and_read.
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

// The three states the unknown guard must not disturb: a known value still
// serializes, and an explicit null still sends null rather than being omitted.
// Without these, skipping unknown could be widened into skipping null and the
// omission tests above would still pass.
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
