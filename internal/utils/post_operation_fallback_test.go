package utils

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type fallbackTestModel struct {
	Name      types.String `tfsdk:"name"`
	Label     types.String `tfsdk:"label"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	Count     types.Int64  `tfsdk:"count"`
	Weight    types.Number `tfsdk:"weight"`
	ModeField types.String `tfsdk:"bgp_as_number"`
}

func TestMergeMissingPlanScalars(t *testing.T) {
	plan := fallbackTestModel{
		Name: types.StringValue("test"), Enabled: types.BoolValue(true),
		Count: types.Int64Value(3), Weight: types.NumberValue(big.NewFloat(1.25)),
		Label:     types.StringValue("fallback"),
		ModeField: types.StringValue("must-not-be-restored"),
	}
	data := map[string]interface{}{"name": "test"}

	merged := MergeMissingPlanScalars(data, plan, "switchpoints", "campus")
	if merged["enabled"] != true {
		t.Fatalf("missing bool fallback = %#v, want true", merged["enabled"])
	}
	if merged["count"] != int64(3) {
		t.Fatalf("missing int fallback = %#v, want 3", merged["count"])
	}
	if merged["label"] != "fallback" {
		t.Fatalf("missing string fallback = %#v", merged["label"])
	}
	if merged["weight"] != "1.25" {
		t.Fatalf("missing number fallback = %#v, want 1.25", merged["weight"])
	}
	if _, ok := merged["bgp_as_number"]; ok {
		t.Fatal("mode-inapplicable field was restored")
	}
}

func TestMergeMissingPlanScalarsPresentNullAndUnknown(t *testing.T) {
	plan := fallbackTestModel{Enabled: types.BoolValue(true), Label: types.StringUnknown()}
	merged := MergeMissingPlanScalars(map[string]interface{}{"enabled": nil}, plan, "switchpoints", "datacenter")
	if merged["enabled"] != nil {
		t.Fatalf("explicit API null must win: %#v", merged["enabled"])
	}
	if _, ok := merged["label"]; ok {
		t.Fatal("unknown plan value was restored")
	}
}

func TestApplyPostOperationFallbackRequiresContext(t *testing.T) {
	data := map[string]interface{}{}
	if result := ApplyPostOperationFallback(context.Background(), data); len(result) != 0 {
		t.Fatalf("ordinary Read must not receive fallback: %#v", result)
	}
}

func TestApplyPostOperationFallbackWithContext(t *testing.T) {
	ctx := WithPostOperationFallback(context.Background(), fallbackTestModel{Enabled: types.BoolValue(true)}, "switchpoints", "datacenter")
	result := ApplyPostOperationFallback(ctx, map[string]interface{}{})
	if result["enabled"] != true {
		t.Fatalf("context fallback = %#v", result)
	}
}

func TestNullifyUnknownPlanScalars(t *testing.T) {
	plan := fallbackTestModel{Enabled: types.BoolUnknown(), Label: types.StringUnknown(), Count: types.Int64Unknown(), Weight: types.NumberUnknown()}
	result := nullifyUnknownPlanScalars(plan).(fallbackTestModel)
	if !result.Enabled.IsNull() || !result.Label.IsNull() || !result.Count.IsNull() || !result.Weight.IsNull() {
		t.Fatalf("unknown scalars must become null: %#v", result)
	}
}
