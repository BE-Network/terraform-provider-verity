package genericresource

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

func listEntries(field spec.FieldSpec, value attr.Value) ([]map[string]attr.Value, bool, error) {
	if value == nil || value.IsNull() || value.IsUnknown() {
		return nil, false, nil
	}
	list, ok := value.(types.List)
	if !ok {
		return nil, false, fmt.Errorf("%s is a list block but holds %T", field.TerraformName, value)
	}
	entries := make([]map[string]attr.Value, 0, len(list.Elements()))
	for position, element := range list.Elements() {
		entry, ok := element.(types.Object)
		if !ok {
			return nil, false, fmt.Errorf("%s[%d] is %T, not an object", field.TerraformName, position, element)
		}
		if entry.IsNull() || entry.IsUnknown() {
			continue
		}
		entries = append(entries, entry.Attributes())
	}
	return entries, len(entries) != 0, nil
}

func createEntry(field spec.FieldSpec, entry map[string]attr.Value, nullables nullableSource) (transport.WireValue, error) {
	object := make(transport.WireObject, len(field.Fields))
	for _, member := range field.Fields {
		if member.Unmanaged {
			continue
		}
		value, held := entry[member.TerraformName]
		if member.Nullable {

			value, held = nullables.entryMember(field, member, entry)
		}
		if !held {
			continue
		}
		wire, send, err := createWire(member, value)
		if err != nil {
			return transport.WireValue{}, fmt.Errorf("%s.%w", field.TerraformName, err)
		}
		if send {
			object[member.APIName] = wire
		}
	}
	return transport.Object(object), nil
}

func createList(field spec.FieldSpec, value attr.Value, nullables nullableSource) (transport.WireValue, bool, error) {
	entries, present, err := listEntries(field, value)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	if !present {
		return createWire(field, nullFor(field))
	}
	wires := make([]transport.WireValue, 0, len(entries))
	for _, entry := range entries {
		wire, err := createEntry(field, entry, nullables)
		if err != nil {
			return transport.WireValue{}, false, err
		}
		wires = append(wires, wire)
	}
	return transport.List(wires), true, nil
}

func entryIndex(field spec.FieldSpec, entry map[string]attr.Value) (int64, bool) {
	value, held := entry[field.Collection.IdentityField]
	index, ok := value.(types.Int64)
	if !held || !ok || index.IsNull() {
		return 0, false
	}
	return index.ValueInt64(), true
}

func updateList(field spec.FieldSpec, plan, state attr.Value, nullables nullableSource, diagnostics *diag.Diagnostics) (transport.WireValue, bool, error) {
	planned, _, err := listEntries(field, plan)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	previous, _, err := listEntries(field, state)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	identity, found := memberNamed(field.Fields, field.Collection.IdentityField)
	if !found {
		return transport.WireValue{}, false, fmt.Errorf("%s: identity %q names no member", field.TerraformName, field.Collection.IdentityField)
	}

	byIndex := make(map[int64]map[string]attr.Value, len(previous))
	for _, entry := range previous {
		if index, ok := entryIndex(field, entry); ok {
			byIndex[index] = entry
		}
	}

	companions, paired, err := referencePairs(field.Fields)
	if err != nil {
		return transport.WireValue{}, false, fmt.Errorf("%s: %w", field.TerraformName, err)
	}

	var sent []transport.WireValue
	kept := make(map[int64]bool, len(planned))
	for _, entry := range planned {
		index, ok := entryIndex(field, entry)
		if !ok {
			continue
		}
		kept[index] = true

		before, exists := byIndex[index]
		if !exists {
			wire, err := createEntry(field, entry, nullables)
			if err != nil {
				return transport.WireValue{}, false, err
			}
			sent = append(sent, wire)
			continue
		}

		object := transport.WireObject{identity.APIName: transport.Int64(index)}
		entryChanged := false
		for _, member := range field.Fields {
			if member.Unmanaged || member.Reference == nil {
				continue
			}
			pairChanged, err := applyReferencePair(member, companions[member.TerraformName], entry, before, object, diagnostics)
			if err != nil {
				return transport.WireValue{}, false, fmt.Errorf("%s: %w", field.TerraformName, err)
			}
			if diagnostics.HasError() {
				return transport.WireValue{}, false, nil
			}
			entryChanged = entryChanged || pairChanged
		}
		for _, member := range field.Fields {
			if member.Unmanaged || paired[member.TerraformName] || member.TerraformName == identity.TerraformName {
				continue
			}
			value, held := entry[member.TerraformName]
			if member.Nullable {
				value, held = nullables.entryMember(field, member, entry)
			}
			if !held {
				continue
			}
			if prior, had := before[member.TerraformName]; had && value.Equal(prior) {
				continue
			}
			wire, send, memberChanged, err := updateWire(member, value)
			if err != nil {
				return transport.WireValue{}, false, fmt.Errorf("%s.%w", field.TerraformName, err)
			}
			entryChanged = entryChanged || memberChanged
			if send {
				object[member.APIName] = wire
			}
		}
		if entryChanged {
			sent = append(sent, transport.Object(object))
		}
	}

	var removed []int64
	for index := range byIndex {
		if !kept[index] {
			removed = append(removed, index)
		}
	}
	sort.Slice(removed, func(i, j int) bool { return removed[i] < removed[j] })
	for _, index := range removed {
		sent = append(sent, transport.Object(transport.WireObject{identity.APIName: transport.Int64(index)}))
	}

	if len(sent) == 0 {
		return transport.WireValue{}, false, nil
	}
	return transport.List(sent), true, nil
}

func listFromAPI(field spec.FieldSpec, raw interface{}, mode string) (attr.Value, error) {
	items, ok := raw.([]interface{})
	if !ok || len(items) == 0 {
		return nullFor(field), nil
	}
	entries := make([]map[string]attr.Value, 0, len(items))
	for _, item := range items {
		object, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		members, err := stateFromAPI(field.Fields, object, mode, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", field.TerraformName, err)
		}
		entries = append(entries, members)
	}
	if len(entries) == 0 {
		return nullFor(field), nil
	}
	return listValue(field, entries)
}

func listValue(field spec.FieldSpec, entries []map[string]attr.Value) (attr.Value, error) {
	elements := make([]attr.Value, 0, len(entries))
	for _, members := range entries {
		entry, diags := types.ObjectValue(attributeTypes(field.Fields), members)
		if diags.HasError() {
			return nil, fmt.Errorf("%s: %v", field.TerraformName, diags)
		}
		elements = append(elements, entry)
	}
	list, diags := types.ListValue(elementType(field), elements)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %v", field.TerraformName, diags)
	}
	return list, nil
}

func settleList(field spec.FieldSpec, value attr.Value) attr.Value {
	entries, present, err := listEntries(field, value)
	if err != nil || !present {
		if value == nil || value.IsUnknown() {
			return nullFor(field)
		}
		return value
	}
	settled := make([]map[string]attr.Value, 0, len(entries))
	for _, entry := range entries {
		settled = append(settled, nullifyUnknown(field.Fields, entry))
	}
	result, err := listValue(field, settled)
	if err != nil {
		return nullFor(field)
	}
	return result
}
