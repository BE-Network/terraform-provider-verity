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

// An indexed collection is a list of objects the API identifies by an index it
// assigns. Terraform holds it as a list block; the API receives, on update, only
// the entries that changed, each named by its index.
//
// Every handwritten resource with a list reconciles it through
// ProcessIndexedArrayUpdates, with per-resource closures deciding what a new,
// changed, or removed entry sends. The rules below are that function's, and the
// closures become the members' declared policies:
//
//   - create sends every entry, each member by its create policies;
//   - update matches plan entries to state entries by index. An index state does
//     not hold is a new entry, sent whole. An index it does hold sends the index
//     and the members that changed, and only when one did. A state index the
//     plan no longer holds is sent as the index alone, which deletes it;
//   - a response array becomes the list, in the order the API returns it, and an
//     empty or missing one an absent block.
//
// Two handwritten details are kept deliberately. An entry whose index is null is
// ignored, and one whose index is unknown is treated as index zero: that is what
// makes an entry written without an index reach the wire as a create with no
// index at all, the behavior tests/unit/lifecycle/index_zero_test.go pins.

// listEntries returns a list block's entries as member maps, and whether the
// block holds any.
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

// createEntry builds one entry as a create sends it.
func createEntry(field spec.FieldSpec, entry map[string]attr.Value, nullables nullableSource) (transport.WireValue, error) {
	object := make(transport.WireObject, len(field.Fields))
	for _, member := range field.Fields {
		if member.Unmanaged {
			continue
		}
		value, held := entry[member.TerraformName]
		if member.Nullable {
			// As at the top level, only the configuration shows whether a null
			// was written, and a member that is not written is left to the
			// server.
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

// createList decides what a create sends for a list block.
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

// entryIndex reads an entry's identity: whether it has one at all, and its value.
// An unknown index reads as zero, as the handwritten reconciliation reads it.
func entryIndex(field spec.FieldSpec, entry map[string]attr.Value) (int64, bool) {
	value, held := entry[field.Collection.IdentityField]
	index, ok := value.(types.Int64)
	if !held || !ok || index.IsNull() {
		return 0, false
	}
	return index.ValueInt64(), true
}

// updateList decides what an update sends for a list block that differs from
// state, and whether anything changed.
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

	// Removed entries are sent as their index alone. The handwritten
	// reconciliation emits them in map iteration order, which differs from run to
	// run; sorting them changes nothing the API can observe and makes the request
	// reproducible.
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

// listFromAPI decodes a response array into a list block, entries in the order
// the API returns them. An entry that is not an object is skipped, and an empty
// or missing array is an absent block, as the handwritten resources record it.
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

// settleList replaces unknown members with nulls in every entry, so a planned
// list can be written to state.
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
