package spec

import (
	"reflect"
	"strings"
	"testing"
)

func bulkTestRegistry() Registry {
	return Registry{
		{TerraformType: "verity_one", API: APIResourceSpec{BulkKey: "one"}, Modes: []Mode{ModeDatacenter}, Operations: OperationSpec{Create: true, Read: true, Update: true, Delete: true}, BulkOrder: map[Mode]BulkOrderSpec{ModeDatacenter: {Put: 10, Patch: 20, Delete: 10}}},
		{TerraformType: "verity_two", API: APIResourceSpec{BulkKey: "two"}, Modes: []Mode{ModeDatacenter}, Operations: OperationSpec{Read: true, Update: true}, BulkOrder: map[Mode]BulkOrderSpec{ModeDatacenter: {Patch: 10}}},
	}
}

func TestBulkOrderValidation(t *testing.T) {
	for _, test := range []struct {
		name string
		edit func(Registry)
		want string
	}{
		{"valid", func(Registry) {}, ""},
		{"missing mode", func(r Registry) { delete(r[0].BulkOrder, ModeDatacenter) }, "has no datacenter bulk order"},
		{"unsupported mode", func(r Registry) { r[0].BulkOrder[ModeCampus] = BulkOrderSpec{Put: 10} }, "unsupported mode"},
		{"negative rank", func(r Registry) { r[0].BulkOrder[ModeDatacenter] = BulkOrderSpec{Put: -10, Patch: 20, Delete: 10} }, "disagrees with operation support"},
		{"missing operation", func(r Registry) { r[0].BulkOrder[ModeDatacenter] = BulkOrderSpec{Patch: 20, Delete: 10} }, "disagrees with operation support"},
		{"patch only put", func(r Registry) { r[1].BulkOrder[ModeDatacenter] = BulkOrderSpec{Put: 30, Patch: 10} }, "disagrees with operation support"},
		{"patch only delete", func(r Registry) { r[1].BulkOrder[ModeDatacenter] = BulkOrderSpec{Patch: 10, Delete: 30} }, "disagrees with operation support"},
		{"rank collision", func(r Registry) { r[1].BulkOrder[ModeDatacenter] = BulkOrderSpec{Patch: 20} }, "is shared by"},
		{"variant disagreement", func(r Registry) { r[1].API.BulkKey = "one" }, "variants for one disagree"},
		{"delete order mismatch", func(r Registry) {
			r[1].Operations = r[0].Operations
			r[1].BulkOrder[ModeDatacenter] = BulkOrderSpec{Put: 20, Patch: 30, Delete: 20}
		}, "DELETE order must reverse PUT"},
		{"shared variant", func(r Registry) { r[1] = r[0]; r[1].TerraformType = "verity_one_variant" }, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			registry := bulkTestRegistry()
			test.edit(registry)
			err := registry.ValidateBulkOrders()
			if test.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestBulkOrdersDeduplicateVariantsAndKeepIndependentPolicies(t *testing.T) {
	registry := bulkTestRegistry()
	registry = append(registry, registry[0])
	for operation, want := range map[string][]string{"PUT": {"one"}, "PATCH": {"two", "one"}, "DELETE": {"one"}} {
		if got := registry.BulkOperationOrder(ModeDatacenter, operation); !reflect.DeepEqual(got, want) {
			t.Fatalf("%s order = %v, want %v", operation, got, want)
		}
	}
	if got := registry.BulkOperationOrder(ModeCampus, "PATCH"); len(got) != 0 {
		t.Fatalf("unsupported mode order = %v", got)
	}
}
