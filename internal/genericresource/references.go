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

	if planBase.IsUnknown() {
		planBase = stateBase
	}
	if planType.IsUnknown() {
		planType = stateType
	}

	baseChanged := !planBase.Equal(stateBase)
	typeChanged := !planType.Equal(stateType)
	if !baseChanged && !typeChanged {
		return false, nil
	}

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

	if baseChanged && !typeChanged {
		object[base.APIName] = referenceWire(planBase)
		return true, nil
	}

	object[base.APIName] = referenceWire(planBase)
	object[companion.APIName] = referenceWire(planType)
	return true, nil
}

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
