package genericresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// The post-operation fallback carries a just-applied plan through the read that
// follows a write, so a response which omits a field the configuration set does
// not turn that field into a null.
//
// The handwritten resources do the same thing by reflecting over their model
// struct for `tfsdk` tags. The generic engine has no struct to reflect over, and
// does not need one: it already holds the plan keyed by Terraform name and the
// spec that says what each field is.
type fallbackKey struct{}

type fallback struct {
	plan     map[string]attr.Value
	merge    func(map[string]interface{}, map[string]attr.Value) map[string]interface{}
	settle   func(map[string]attr.Value) map[string]attr.Value
	setState func(context.Context, *tfsdk.State, map[string]attr.Value) diag.Diagnostics
}

// withFallback records the plan an operation applied, for the read it triggers.
func (r *Resource) withFallback(ctx context.Context, plan map[string]attr.Value) context.Context {
	return context.WithValue(ctx, fallbackKey{}, fallback{
		plan:     plan,
		merge:    r.mergePlanScalars,
		settle:   func(values map[string]attr.Value) map[string]attr.Value { return nullifyUnknown(r.spec.Fields, values) },
		setState: r.setState,
	})
}

// applyFallback fills members the response left out from the plan that produced
// it. Outside a post-operation read there is no fallback and the data stands.
func applyFallback(ctx context.Context, data map[string]interface{}) map[string]interface{} {
	carried, ok := ctx.Value(fallbackKey{}).(fallback)
	if !ok {
		return data
	}
	return carried.merge(data, carried.plan)
}

// setFallbackState writes the applied plan as state when a read declines to
// verify, which happens while the provider's own writes are still in flight.
// Unknowns settle to null: they are Computed, so the next read fills them in.
func setFallbackState(ctx context.Context, state *tfsdk.State) (bool, diag.Diagnostics) {
	carried, ok := ctx.Value(fallbackKey{}).(fallback)
	if !ok {
		return false, nil
	}
	return true, carried.setState(ctx, state, carried.settle(carried.plan))
}
