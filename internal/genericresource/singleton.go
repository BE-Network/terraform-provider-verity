package genericresource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

// A singleton is an object the API sends whole, which every handwritten resource
// exposes as a list block holding at most one entry. The API shape and the state
// shape differ on purpose: changing the block to an attribute would break every
// existing state file, so the engine keeps the list and reads and writes its
// first entry.
//
// The rules below are the handwritten resources' own, measured against them by
// the differential tests rather than assumed:
//
//   - create sends the object when the block is written, carrying each member
//     its create policies allow, and nothing when the block is absent;
//   - update considers the object only when both the plan and the state hold an
//     entry, sends only the members that changed, and sends nothing when the
//     block is added or removed;
//   - a response object becomes a one-entry list, and a missing one an absent
//     block.

// elementType is the object type of a singleton block's one entry.
func elementType(field spec.FieldSpec) types.ObjectType {
	return types.ObjectType{AttrTypes: attributeTypes(field.Fields)}
}

// nullFor is the null value of a field's state type, which for a singleton is a
// null list of its entry type rather than a null scalar.
func nullFor(field spec.FieldSpec) attr.Value {
	if field.Kind == spec.FieldKindObject {
		return types.ListNull(elementType(field))
	}
	return nullOf(field.Kind)
}

// singletonMembers returns the first entry of a singleton block, and whether the
// block holds one at all. A null, unknown, or empty list is an absent block.
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

// createSingleton decides what a create sends for a singleton. A block that is
// not written goes through the object's own create policy; one that is written
// is sent as an object, each member decided by its own policies, and sent even
// when no member survives, because the handwritten resources send it that way.
func createSingleton(field spec.FieldSpec, value attr.Value) (transport.WireValue, bool, error) {
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
		if !held {
			continue
		}
		wire, send, err := createWire(member, memberValue)
		if err != nil {
			return transport.WireValue{}, false, fmt.Errorf("%s.%w", field.TerraformName, err)
		}
		if send {
			object[member.APIName] = wire
		}
	}
	return transport.Object(object), true, nil
}

// updateSingleton decides what an update sends for a singleton that differs from
// state. It reports whether anything changed; only then is the object sent, and
// it carries only the members that changed.
//
// Adding or removing the block sends nothing, as it does in the handwritten
// resources. Removal leaves the server's object alone, and an addition is picked
// up on a later update once the read has put the server's object into state.
func updateSingleton(field spec.FieldSpec, plan, state attr.Value, diagnostics *diag.Diagnostics) (transport.WireValue, bool, error) {
	planned, plannedPresent, err := singletonMembers(field, plan)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	previous, previousPresent, err := singletonMembers(field, state)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	if !plannedPresent || !previousPresent {
		return transport.WireValue{}, false, nil
	}

	object := make(transport.WireObject, len(field.Fields))
	changed := false

	// A reference pair inside the object is decided together, exactly as at the
	// top level, and both halves clear to an empty string rather than by the
	// omission its siblings use.
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
		if member.Unmanaged || paired[member.TerraformName] {
			continue
		}
		value, held := planned[member.TerraformName]
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

// singletonFromAPI decodes a response object into a one-entry block. Each member
// follows its own mode and absence policy, so a member the running mode does not
// expose reads as null even when the response carries it.
func singletonFromAPI(field spec.FieldSpec, raw interface{}, mode string) (attr.Value, error) {
	object, ok := raw.(map[string]interface{})
	if !ok {
		// Not an object: the handwritten resources record no block in that case.
		return nullFor(field), nil
	}
	members, err := stateFromAPI(field.Fields, object, mode, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", field.TerraformName, err)
	}
	return singletonValue(field, members)
}

// singletonValue builds the one-entry list a block's state holds.
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

// settleSingleton replaces unknown members with nulls so a planned block can be
// written to state. An unknown cannot be stored, and a Computed member is filled
// in by the next read.
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
