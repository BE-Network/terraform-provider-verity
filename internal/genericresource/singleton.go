package genericresource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

func elementType(field spec.FieldSpec) types.ObjectType {
	return types.ObjectType{AttrTypes: attributeTypes(field.Fields)}
}

func nullFor(field spec.FieldSpec) attr.Value {
	if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
		return types.ListNull(elementType(field))
	}
	return nullOf(field.Kind)
}

func singletonMembers(field spec.FieldSpec, value attr.Value) (map[string]attr.Value, bool, error) {
	if value == nil || value.IsNull() || value.IsUnknown() {
		return nil, false, nil
	}
	list, ok := value.(types.List)
	if !ok {
		return nil, false, fmt.Errorf("%s is a singleton block but holds %T", field.TerraformName, value)
	}
	elements := list.Elements()
	if len(elements) == 0 {
		return nil, false, nil
	}
	entry, ok := elements[0].(types.Object)
	if !ok {
		return nil, false, fmt.Errorf("%s entry is %T, not an object", field.TerraformName, elements[0])
	}
	if entry.IsNull() || entry.IsUnknown() {
		return nil, false, nil
	}
	return entry.Attributes(), true, nil
}

func createSingleton(field spec.FieldSpec, value attr.Value, nullables nullableSource) (transport.WireValue, bool, error) {
	members, present, err := singletonMembers(field, value)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	if !present {
		return createWire(field, nullFor(field))
	}
	object := make(transport.WireObject, len(field.Fields))
	for _, member := range field.Fields {
		if member.Unmanaged {
			continue
		}
		memberValue, held := members[member.TerraformName]
		if member.Nullable {

			memberValue, held = nullables.singletonMember(field, member, members)
		}
		if !held {
			continue
		}
		var (
			wire transport.WireValue
			send bool
			err  error
		)
		if member.Kind == spec.FieldKindList {

			wire, send, err = createList(member, memberValue, nullableSource{})
			if err == nil && !send {
				wire, send = transport.List([]transport.WireValue{}), true
			}
		} else {
			wire, send, err = createWire(member, memberValue)
		}
		if err != nil {
			return transport.WireValue{}, false, fmt.Errorf("%s.%w", field.TerraformName, err)
		}
		if send {
			object[member.APIName] = wire
		}
	}
	return transport.Object(object), true, nil
}

func updateEmptySingleton(field spec.FieldSpec, plan, state attr.Value) (wire transport.WireValue, send, changed bool, err error) {
	_, plannedPresent, err := singletonMembers(field, plan)
	if err != nil {
		return transport.WireValue{}, false, false, err
	}
	_, previousPresent, err := singletonMembers(field, state)
	if err != nil {
		return transport.WireValue{}, false, false, err
	}
	switch {
	case plannedPresent == previousPresent:
		return transport.WireValue{}, false, false, nil
	case plannedPresent:
		return transport.Object(transport.WireObject{}), true, true, nil
	default:
		return transport.WireValue{}, false, true, nil
	}
}

func hasManagedMembers(field spec.FieldSpec) bool {
	for _, member := range field.Fields {
		if !member.Unmanaged {
			return true
		}
	}
	return false
}

func updateSingleton(field spec.FieldSpec, plan, state attr.Value, nullables nullableSource, diagnostics *diag.Diagnostics) (transport.WireValue, bool, error) {
	planned, plannedPresent, err := singletonMembers(field, plan)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	previous, previousPresent, err := singletonMembers(field, state)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	if !plannedPresent && !previousPresent {
		return transport.WireValue{}, false, nil
	}

	object := make(transport.WireObject, len(field.Fields))
	changed := false

	for _, member := range field.Fields {
		if member.Unmanaged || member.Kind != spec.FieldKindList {
			continue
		}
		wire, listChanged, err := updateList(member, planned[member.TerraformName], previous[member.TerraformName], nullableSource{}, diagnostics)
		if err != nil {
			return transport.WireValue{}, false, fmt.Errorf("%s: %w", field.TerraformName, err)
		}
		if diagnostics.HasError() {
			return transport.WireValue{}, false, nil
		}
		if listChanged {
			object[member.APIName] = wire
			changed = true
		}
	}

	if !plannedPresent || !previousPresent {

		if !changed {
			return transport.WireValue{}, false, nil
		}
		return transport.Object(object), true, nil
	}

	companions, paired, err := referencePairs(field.Fields)
	if err != nil {
		return transport.WireValue{}, false, fmt.Errorf("%s: %w", field.TerraformName, err)
	}
	for _, member := range field.Fields {
		if member.Unmanaged || member.Reference == nil {
			continue
		}
		pairChanged, err := applyReferencePair(member, companions[member.TerraformName], planned, previous, object, diagnostics)
		if err != nil {
			return transport.WireValue{}, false, fmt.Errorf("%s: %w", field.TerraformName, err)
		}
		if diagnostics.HasError() {
			return transport.WireValue{}, false, nil
		}
		changed = changed || pairChanged
	}

	for _, member := range field.Fields {
		if member.Unmanaged || paired[member.TerraformName] || member.Kind == spec.FieldKindList {
			continue
		}
		value, held := planned[member.TerraformName]
		if member.Nullable {
			value, held = nullables.singletonMember(field, member, planned)
		}
		if !held {
			continue
		}
		if before, had := previous[member.TerraformName]; had && value.Equal(before) {
			continue
		}
		wire, send, memberChanged, err := updateWire(member, value)
		if err != nil {
			return transport.WireValue{}, false, fmt.Errorf("%s.%w", field.TerraformName, err)
		}
		changed = changed || memberChanged
		if send {
			object[member.APIName] = wire
		}
	}
	if !changed {
		return transport.WireValue{}, false, nil
	}
	return transport.Object(object), true, nil
}

func singletonFromAPI(field spec.FieldSpec, raw interface{}, mode string) (attr.Value, error) {
	object, ok := raw.(map[string]interface{})
	if !ok {

		return nullFor(field), nil
	}
	members, err := stateFromAPI(field.Fields, object, mode, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", field.TerraformName, err)
	}

	return singletonValue(field, members)
}

func singletonValue(field spec.FieldSpec, members map[string]attr.Value) (attr.Value, error) {
	entry, diags := types.ObjectValue(attributeTypes(field.Fields), members)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %v", field.TerraformName, diags)
	}
	list, diags := types.ListValue(elementType(field), []attr.Value{entry})
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %v", field.TerraformName, diags)
	}
	return list, nil
}

func settleSingleton(field spec.FieldSpec, value attr.Value) attr.Value {
	members, present, err := singletonMembers(field, value)
	if err != nil || !present {
		if value == nil || value.IsUnknown() {
			return nullFor(field)
		}
		return value
	}
	settled, err := singletonValue(field, nullifyUnknown(field.Fields, members))
	if err != nil {
		return nullFor(field)
	}
	return settled
}
