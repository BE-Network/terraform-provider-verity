package main

import "testing"

// Every type setterFor accepts must name a transport helper with the matching
// generated-SDK target type. This table is deliberately exhaustive: adding a
// supported type without its width-specific helper otherwise first fails when a
// future OpenAPI document happens to use it.
func TestSetterForSupportsEverySDKScalarType(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"*string":         "wireStringPtr",
		"*bool":           "wireBoolPtr",
		"*int32":          "wireInt32Ptr",
		"*int64":          "wireInt64Ptr",
		"*float32":        "wireFloat32Ptr",
		"*float64":        "wireFloat64Ptr",
		"NullableInt32":   "wireNullableInt32",
		"NullableInt64":   "wireNullableInt64",
		"NullableFloat32": "wireNullableFloat32",
		"NullableFloat64": "wireNullableFloat64",
		"NullableString":  "wireNullableString",
	}
	for goType, want := range cases {
		got, supported := setterFor(goType)
		if !supported || got != want {
			t.Errorf("setterFor(%q) = (%q, %v), want (%q, true)", goType, got, supported, want)
		}
	}
}

func TestSetterForRejectsUnsupportedSDKType(t *testing.T) {
	t.Parallel()

	if setter, supported := setterFor("*time.Duration"); supported || setter != "" {
		t.Errorf("setterFor unsupported type = (%q, %v), want (\"\", false)", setter, supported)
	}
}
