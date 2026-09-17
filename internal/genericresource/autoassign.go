package genericresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

// An auto-assignment pair is a value the server can choose and a boolean flag
// that asks it to. While the flag is on, the value belongs to the server: the
// configuration may not set it, the plan cannot know it in advance, and the
// request carries the flag instead of the value.
//
// The rules below are the ones verity_service, verity_tenant, and verity_fabric
// implement identically, surveyed from their handwritten code rather than
// designed here. verity_switchpoint implements the same plan rules but a
// narrower update rule for seven of its ten pairs; it is not servable yet, and
// that difference has to be decided before it is.

type autoAssignedPair struct {
	value spec.FieldSpec
	flag  spec.FieldSpec
}

// autoAssignmentPairs indexes each auto-assigned value with its flag, and every
// field that belongs to a pair so the ordinary loops can skip both halves.
func autoAssignmentPairs(fields []spec.FieldSpec) ([]autoAssignedPair, map[string]bool, error) {
	byName := make(map[string]spec.FieldSpec, len(fields))
	for _, field := range fields {
		byName[field.TerraformName] = field
	}
	var pairs []autoAssignedPair
	paired := make(map[string]bool)
	for _, field := range fields {
		if field.Unmanaged || field.AutoAssignment == nil {
			continue
		}
		flag, found := byName[field.AutoAssignment.FlagField]
		if !found || flag.Kind != spec.FieldKindBool {
			return nil, nil, fmt.Errorf("%s is auto-assigned by %q, which is not a bool field of the resource",
				field.TerraformName, field.AutoAssignment.FlagField)
		}
		pairs = append(pairs, autoAssignedPair{value: field, flag: flag})
		paired[field.TerraformName] = true
		paired[flag.TerraformName] = true
	}
	return pairs, paired, nil
}

func isTrue(value attr.Value) bool {
	typed, ok := value.(types.Bool)
	return ok && !typed.IsNull() && !typed.IsUnknown() && typed.ValueBool()
}

func isKnown(value attr.Value) bool {
	return value != nil && !value.IsNull() && !value.IsUnknown()
}

// isSet reports whether a configuration names a value. An empty string counts as
// not set, as the handwritten validation treats it.
func isSet(value attr.Value) bool {
	if !isKnown(value) {
		return false
	}
	if text, ok := value.(types.String); ok {
		return text.ValueString() != ""
	}
	return true
}

// createAutoAssigned decides what a create sends for one pair. With the flag on,
// the flag is sent and the value left to the server. Otherwise the value follows
// its own policies — including the configuration scan for a nullable value — and
// the flag is sent when the plan knows it.
func createAutoAssigned(pair autoAssignedPair, plan map[string]attr.Value, nullables nullableSource, object transport.WireObject) error {
	flag := plan[pair.flag.TerraformName]
	if isTrue(flag) {
		object[pair.flag.APIName] = transport.Bool(true)
		return nil
	}

	value := plan[pair.value.TerraformName]
	if pair.value.Nullable {
		written, isWritten := nullables.known(pair.value)
		if !isWritten {
			value = nil
		} else {
			value = written
		}
	}
	if value != nil {
		wire, send, err := createWire(pair.value, value)
		if err != nil {
			return err
		}
		if send {
			object[pair.value.APIName] = wire
		}
	}
	if isKnown(flag) {
		object[pair.flag.APIName] = transport.Bool(flag.(types.Bool).ValueBool())
	}
	return nil
}

// updateAutoAssigned decides what an update sends for one pair, and reports
// whether anything changed.
//
// Two API behaviors shape it. A flag is sent only when the configuration states
// it, so a flag the server set is not echoed back as a user decision. And turning
// assignment off is ignored by the API unless the value travels with it, so the
// value is resent — the planned one, or failing that the one in state.
func updateAutoAssigned(pair autoAssignedPair, plan, state, config map[string]attr.Value, object transport.WireObject) (bool, error) {
	planValue, stateValue := plan[pair.value.TerraformName], state[pair.value.TerraformName]
	planFlag, stateFlag := plan[pair.flag.TerraformName], state[pair.flag.TerraformName]

	valueChanged := planValue != nil && !planValue.IsUnknown() && !planValue.Equal(stateValue)
	flagChanged := planFlag != nil && !planFlag.Equal(stateFlag)
	if !valueChanged && !flagChanged {
		return false, nil
	}

	if valueChanged {
		wire, send, _, err := updateWire(pair.value, planValue)
		if err != nil {
			return false, err
		}
		if send {
			object[pair.value.APIName] = wire
		}
	}

	switch {
	case flagChanged:
		if configured := config[pair.flag.TerraformName]; configured != nil && !configured.IsNull() {
			if isKnown(planFlag) {
				object[pair.flag.APIName] = transport.Bool(planFlag.(types.Bool).ValueBool())
			}
			if isTrue(stateFlag) && isKnown(planFlag) && !isTrue(planFlag) {
				resend := planValue
				if !isKnown(resend) {
					resend = stateValue
				}
				if isKnown(resend) {
					wire, err := toWire(pair.value, resend)
					if err != nil {
						return false, err
					}
					object[pair.value.APIName] = wire
				}
			}
		}
	case valueChanged:
		// The flag travels with a changed value so the server does not keep a
		// stale assignment decision alongside it.
		switch {
		case isKnown(planFlag):
			object[pair.flag.APIName] = transport.Bool(planFlag.(types.Bool).ValueBool())
		case isKnown(stateFlag):
			object[pair.flag.APIName] = transport.Bool(stateFlag.(types.Bool).ValueBool())
		default:
			object[pair.flag.APIName] = transport.Bool(false)
		}
	}
	return true, nil
}

// planAutoAssignment applies the pair rules to the plan: refuse a value written
// while the flag is on, mark the value unknown whenever the server is about to
// choose it, and keep the state value when a change to it would be ignored.
//
// It reads the plan as it arrived rather than as ModifyPlan has already changed
// it, which is what the handwritten resources compare against.
func (r *Resource) planAutoAssignment(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	pairs, _, err := autoAssignmentPairs(r.spec.Fields)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource Specification", err.Error())
		return
	}
	if len(pairs) == 0 {
		return
	}

	plan, diags := readScalars(ctx, req.Plan, r.spec.Fields)
	resp.Diagnostics.Append(diags...)
	config, diags := readScalars(ctx, req.Config, r.spec.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, pair := range pairs {
		if isTrue(config[pair.flag.TerraformName]) && isSet(config[pair.value.TerraformName]) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("'%s' cannot be specified when auto-assigned", pair.value.TerraformName),
				fmt.Sprintf("The '%s' field cannot be specified in the configuration when '%s' is set to true. The API will assign this value automatically.",
					pair.value.TerraformName, pair.flag.TerraformName),
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	valuePath := func(pair autoAssignedPair) path.Path { return path.Root(pair.value.TerraformName) }

	if req.State.Raw.IsNull() {
		for _, pair := range pairs {
			if isTrue(plan[pair.flag.TerraformName]) {
				resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, valuePath(pair), unknownOf(pair.value.Kind))...)
			}
		}
		return
	}

	state, diags := readScalars(ctx, req.State, r.spec.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, pair := range pairs {
		planFlag := plan[pair.flag.TerraformName]
		planValue, stateValue := plan[pair.value.TerraformName], state[pair.value.TerraformName]
		trigger := changedTrigger(pair, plan, state)

		if isTrue(planFlag) {
			switch {
			case !planFlag.Equal(state[pair.flag.TerraformName]):
				resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, valuePath(pair), unknownOf(pair.value.Kind))...)
				resp.Diagnostics.AddWarning(
					fmt.Sprintf("'%s' will be assigned by the API", pair.value.TerraformName),
					fmt.Sprintf("The '%s' field will be automatically assigned by the API because '%s' is being set to true.",
						pair.value.TerraformName, pair.flag.TerraformName),
				)
			case trigger != "":
				resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, valuePath(pair), unknownOf(pair.value.Kind))...)
				resp.Diagnostics.AddWarning(
					fmt.Sprintf("'%s' will be updated by the API", pair.value.TerraformName),
					fmt.Sprintf("The '%s' field will be automatically updated by the API because '%s' is set to true and '%s' is changing.",
						pair.value.TerraformName, pair.flag.TerraformName, trigger),
				)
			case !planValue.Equal(stateValue):
				resp.Diagnostics.AddWarning(
					fmt.Sprintf("Ignoring %s changes with auto-assignment enabled", pair.value.TerraformName),
					fmt.Sprintf("The '%s' field changes will be ignored because '%s' is set to true.",
						pair.value.TerraformName, pair.flag.TerraformName),
				)
				if isKnown(stateValue) {
					resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, valuePath(pair), stateValue)...)
				}
			}
			continue
		}

		// With assignment off, a change to a field the value is derived from still
		// makes the server recompute it — unless the configuration pins the value.
		if trigger != "" && planValue.Equal(stateValue) && !isSet(config[pair.value.TerraformName]) {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, valuePath(pair), unknownOf(pair.value.Kind))...)
		}
	}
}

// changedTrigger returns the first field the value is recomputed from whose plan
// differs from state, or "" when none does.
func changedTrigger(pair autoAssignedPair, plan, state map[string]attr.Value) string {
	for _, name := range pair.value.AutoAssignment.RecomputedWhen {
		planned, before := plan[name], state[name]
		if planned != nil && !planned.Equal(before) {
			return name
		}
	}
	return ""
}

func unknownOf(kind spec.FieldKind) attr.Value {
	switch kind {
	case spec.FieldKindString:
		return types.StringUnknown()
	case spec.FieldKindBool:
		return types.BoolUnknown()
	case spec.FieldKindInt64:
		return types.Int64Unknown()
	case spec.FieldKindNumber:
		return types.NumberUnknown()
	default:
		return nil
	}
}
