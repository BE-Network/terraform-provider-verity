package genericresource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
)

// A reference is two fields the API insists on seeing together: a value naming
// an object, and a companion naming that object's type.
//
// They cannot be decided one at a time. Changing the type without the value, or
// clearing one and leaving the other, describes an object the server cannot
// resolve, so the update rule is about the pair rather than about either half:
// when the type changes, both halves are sent, and both clear to an empty
// string rather than being omitted.
//
// The validation is the handwritten helpers' own, called rather than
// reimplemented. They decide which combinations are refused and word the
// diagnostics, and a migrated resource should refuse exactly what it refused
// before, in the same words.
func applyReferencePair(base, companion spec.FieldSpec, plan, state map[string]attr.Value,
	object transport.WireObject, diagnostics *diag.Diagnostics) (changed bool, err error) {

	planBase, err := stringOf(base, plan)
	if err != nil {
		return false, err
	}
	stateBase, err := stringOf(base, state)
	if err != nil {
		return false, err
	}
	planType, err := stringOf(companion, plan)
	if err != nil {
		return false, err
	}
	stateType, err := stringOf(companion, state)
	if err != nil {
		return false, err
	}

	baseChanged := !planBase.Equal(stateBase)
	typeChanged := !planType.Equal(stateType)
	if !baseChanged && !typeChanged {
		return false, nil
	}

	// More than one permitted target type means the companion is never implied by
	// the value, so both halves are always sent and the unchanged half is taken
	// from state rather than left out.
	if len(base.Reference.AllowedTypes) > 1 {
		if !utils.ValidateMultipleRefTypesSupported(diagnostics, planBase, planType, base.TerraformName, companion.TerraformName) {
			return false, nil
		}
		effectiveBase, effectiveType := stateBase, stateType
		if baseChanged {
			effectiveBase = planBase
		}
		if typeChanged {
			effectiveType = planType
		}
		object[base.APIName] = referenceWire(effectiveBase)
		object[companion.APIName] = referenceWire(effectiveType)
		return true, nil
	}

	if !utils.ValidateOneRefTypeSupported(diagnostics, planBase, planType, base.TerraformName, companion.TerraformName, baseChanged, typeChanged) {
		return false, nil
	}

	// With one permitted type the companion is implied, so a value-only change
	// does not need to restate it.
	if baseChanged && !typeChanged {
		object[base.APIName] = referenceWire(planBase)
		return true, nil
	}

	object[base.APIName] = referenceWire(planBase)
	object[companion.APIName] = referenceWire(planType)
	return true, nil
}

// referenceWire is how a reference half reaches the wire: its value, or an empty
// string when it has none. A reference is never cleared by omission and never
// by a null, whichever half it is.
func referenceWire(value types.String) transport.WireValue {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return transport.String("")
	}
	return transport.String(value.ValueString())
}

func stringOf(field spec.FieldSpec, values map[string]attr.Value) (types.String, error) {
	value, held := values[field.TerraformName]
	if !held || value == nil {
		return types.StringNull(), nil
	}
	typed, ok := value.(types.String)
	if !ok {
		return types.StringNull(), fmt.Errorf("%s is half of a reference but holds %T; both halves must be strings", field.TerraformName, value)
	}
	return typed, nil
}

// referencePairs indexes the fields that are half of a reference, so the update
// loop can hand a pair to applyReferencePair once and skip its companion.
//
// A companion the spec names but the resource does not carry is a registry
// error rather than something to work around: it would mean sending one half of
// a pair the API requires whole.
func referencePairs(fields []spec.FieldSpec) (map[string]spec.FieldSpec, map[string]bool, error) {
	byName := make(map[string]spec.FieldSpec, len(fields))
	for _, field := range fields {
		byName[field.TerraformName] = field
	}
	companions := make(map[string]spec.FieldSpec)
	paired := make(map[string]bool)
	for _, field := range fields {
		if field.Reference == nil {
			continue
		}
		companion, found := byName[field.Reference.TypeField]
		if !found {
			return nil, nil, fmt.Errorf("%s references type field %q, which the resource does not carry",
				field.TerraformName, field.Reference.TypeField)
		}
		if companion.Kind != spec.FieldKindString || field.Kind != spec.FieldKindString {
			return nil, nil, fmt.Errorf("%s and %s are a reference pair, so both must be strings",
				field.TerraformName, companion.TerraformName)
		}
		companions[field.TerraformName] = companion
		paired[field.TerraformName] = true
		paired[companion.TerraformName] = true
	}
	return companions, paired, nil
}
