package genericresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type fallbackKey struct{}

type fallback struct {
	plan     map[string]attr.Value
	merge    func(map[string]interface{}, map[string]attr.Value) map[string]interface{}
	settle   func(map[string]attr.Value) map[string]attr.Value
	setState func(context.Context, *tfsdk.State, map[string]attr.Value) diag.Diagnostics
}

func (r *Resource) withFallback(ctx context.Context, plan map[string]attr.Value) context.Context {
	return context.WithValue(ctx, fallbackKey{}, fallback{
		plan:     plan,
		merge:    r.mergePlanScalars,
		settle:   func(values map[string]attr.Value) map[string]attr.Value { return nullifyUnknown(r.spec.Fields, values) },
		setState: r.setState,
	})
}

func applyFallback(ctx context.Context, data map[string]interface{}) map[string]interface{} {
	carried, ok := ctx.Value(fallbackKey{}).(fallback)
	if !ok {
		return data
	}
	return carried.merge(data, carried.plan)
}

func setFallbackState(ctx context.Context, state *tfsdk.State) (bool, diag.Diagnostics) {
	carried, ok := ctx.Value(fallbackKey{}).(fallback)
	if !ok {
		return false, nil
	}
	return true, carried.setState(ctx, state, carried.settle(carried.plan))
}
