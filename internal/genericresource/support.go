package genericresource

import (
	"fmt"

	"terraform-provider-verity/internal/spec"
)

// Supported reports whether the generic engine can serve a resource, and if not,
// the first reason it cannot.
//
// It is the single definition of the engine's reach. The engine refuses to
// construct anything it rejects, the adapter generator emits adapters only for
// what it accepts, and the tests enumerate the servable set through it, so the
// three cannot disagree about which resources can be migrated.
//
// Each refusal names something the engine would otherwise get silently wrong
// rather than something it merely has not been asked to do:
//
//   - a list strategy other than indexed_patch has no registry user, so there
//     is nothing to check an implementation of it against;
//   - an auto-assignment pair inside an object has a flag that changes whether
//     its value is sent at all, and the pair rules run only at the top level;
//   - a nullable member of a singleton block waits on verity_switchpoint, the
//     only resource with one, whose auto-assignment decision is still open.
func Supported(resource spec.ResourceSpec) error {
	for _, field := range resource.Fields {
		if field.Unmanaged {
			continue
		}
		if err := supportedField(field, nil); err != nil {
			return err
		}
	}
	return nil
}

// containers lists the blocks a field sits in, outermost first: empty at the top
// level, [object] for a member of a singleton, [list] for a member of a list
// entry, and [object list] for a member of a list inside a singleton.
func supportedField(field spec.FieldSpec, containers []spec.FieldKind) error {
	nested := len(containers) != 0
	insideSingleton := nested && containers[0] == spec.FieldKindObject
	if field.AutoAssignment != nil && nested {
		return fmt.Errorf("%s is an auto-assignment pair inside an object, which the engine does not implement yet", field.APIName)
	}
	switch field.Kind {
	case spec.FieldKindString, spec.FieldKindBool, spec.FieldKindInt64, spec.FieldKindNumber:
		// A nullable member of a list entry is served: the scan records each
		// entry's attributes under its index. So is one of a singleton, which the
		// scan records under "block.member". A nullable member of a list inside a
		// singleton would need the scan's key for a doubly nested entry, and no
		// resource has one.
		if len(containers) > 1 && field.Nullable {
			return fmt.Errorf("%s is a nullable member of a list inside a singleton block, which the engine does not serve", field.APIName)
		}
		return nil
	case spec.FieldKindList:
		// A list may sit at the top level or directly inside a singleton, which
		// is verity_fabric's object_properties.system_graphs, the provider's
		// only two-level nesting. Nothing nests deeper.
		if nested && !(len(containers) == 1 && insideSingleton) {
			return fmt.Errorf("%s is a list nested more deeply than inside a singleton, which the engine does not serve", field.APIName)
		}
		// indexed_patch is the only list strategy any registry resource uses, so
		// it is the only one implemented; a strategy with no user has nothing to
		// be checked against.
		if field.Collection == nil || field.Collection.Strategy != spec.CollectionIndexedPatch {
			return fmt.Errorf("%s is a list without the indexed_patch strategy, which the engine does not serve", field.APIName)
		}
		if field.ElementKind != spec.FieldKindObject || len(field.Fields) == 0 {
			return fmt.Errorf("%s is a list of %s, and the engine serves lists of objects", field.APIName, field.ElementKind)
		}
		identity, found := memberNamed(field.Fields, field.Collection.IdentityField)
		if !found || identity.Kind != spec.FieldKindInt64 {
			return fmt.Errorf("%s identifies its entries by %q, which is not an int64 member", field.APIName, field.Collection.IdentityField)
		}
		inner := append(append([]spec.FieldKind(nil), containers...), spec.FieldKindList)
		for _, member := range field.Fields {
			if member.Unmanaged {
				continue
			}
			if err := supportedField(member, inner); err != nil {
				return err
			}
		}
		return nil
	case spec.FieldKindObject:
		if nested {
			return fmt.Errorf("%s is an object inside a block, which the engine does not serve yet", field.APIName)
		}
		if field.Collection == nil || field.Collection.Strategy != spec.CollectionSingleton {
			return fmt.Errorf("%s is an object without the singleton strategy, which the engine does not serve", field.APIName)
		}
		// A singleton with no members is the API declaring an object with no
		// properties. Two resources expose it as an empty block; it carries no
		// values, so only its presence has any meaning, and updateEmptySingleton
		// decides what a change of presence sends.
		for _, member := range field.Fields {
			if member.Unmanaged {
				continue
			}
			if err := supportedField(member, []spec.FieldKind{spec.FieldKindObject}); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("%s has unknown kind %q", field.APIName, field.Kind)
	}
}

func memberNamed(fields []spec.FieldSpec, name string) (spec.FieldSpec, bool) {
	for _, field := range fields {
		if field.TerraformName == name {
			return field, true
		}
	}
	return spec.FieldSpec{}, false
}
