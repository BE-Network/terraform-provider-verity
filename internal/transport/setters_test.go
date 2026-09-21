package transport

import (
	"testing"

	"terraform-provider-verity/openapi"
)

func TestWidthSpecificSettersMatchSDKTargets(t *testing.T) {
	t.Parallel()

	var int32Target *int32
	if err := wireInt32Ptr(Int64(42), &int32Target); err != nil || int32Target == nil || *int32Target != 42 {
		t.Fatalf("wireInt32Ptr = %v, %v", int32Target, err)
	}
	var float32Target *float32
	if err := wireFloat32Ptr(Decimal("2.5"), &float32Target); err != nil || float32Target == nil || *float32Target != 2.5 {
		t.Fatalf("wireFloat32Ptr = %v, %v", float32Target, err)
	}

	var nullableInt32 openapi.NullableInt32
	if err := wireNullableInt32(Int64(42), &nullableInt32); err != nil || nullableInt32.Get() == nil || *nullableInt32.Get() != 42 {
		t.Fatalf("wireNullableInt32 = %v, %v", nullableInt32.Get(), err)
	}
	var nullableFloat32 openapi.NullableFloat32
	if err := wireNullableFloat32(Decimal("2.5"), &nullableFloat32); err != nil || nullableFloat32.Get() == nil || *nullableFloat32.Get() != 2.5 {
		t.Fatalf("wireNullableFloat32 = %v, %v", nullableFloat32.Get(), err)
	}
}

func TestInt32SettersRejectOverflow(t *testing.T) {
	t.Parallel()

	var target *int32
	if err := wireInt32Ptr(Int64(1<<31), &target); err == nil {
		t.Fatal("wireInt32Ptr accepted an integer outside int32 range")
	}
}
