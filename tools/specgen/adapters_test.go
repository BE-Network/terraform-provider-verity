package main

import (
	"path/filepath"
	"strings"
	"testing"

	"terraform-provider-verity/internal/spec"
)

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

func TestAdapterForAResourceWithNoBodyFields(t *testing.T) {
	root := filepath.Join("..", "..")
	registry, err := readRegistry(filepath.Join(root, "specs", "generated_registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	structs, err := discoverStructs(filepath.Join(root, "openapi"))
	if err != nil {
		t.Fatal(err)
	}
	var collection spec.ResourceSpec
	for _, resource := range registry {
		if resource.TerraformType == "verity_fabric_collection" {
			collection = resource
		}
	}
	adapter, reason := planOne(collection, structs)
	if reason != "" || !adapter.EmptyBody {
		t.Fatalf("fabric collection adapter = %+v, reason %q", adapter, reason)
	}

	withField := collection
	withField.Fields = append(append([]spec.FieldSpec(nil), collection.Fields...), spec.FieldSpec{
		TerraformName: "enable", APIName: "enable", Kind: spec.FieldKindBool, Access: spec.AccessOptionalComputed,
	})
	if _, reason := planOne(withField, structs); !strings.Contains(reason, `no "enable" field`) {
		t.Fatalf("a field that the SDK body can't carry must be rejected, got reason %q", reason)
	}
}
