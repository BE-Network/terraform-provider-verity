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
		if err := supportedField(field, ""); err != nil {
			return err
		}
	}
	return nil
}

// container is the kind of block a member sits in: "" at the top level, or
// object or list for a member of a singleton or of a list entry.
func supportedField(field spec.FieldSpec, container spec.FieldKind) error {
	nested := container != ""
	if field.AutoAssignment != nil && nested {
		return fmt.Errorf("%s is an auto-assignment pair inside an object, which the engine does not implement yet", field.APIName)
	}
	switch field.Kind {
	case spec.FieldKindString, spec.FieldKindBool, spec.FieldKindInt64, spec.FieldKindNumber:
		// A nullable member of a list entry is served: the configuration scan
		// records each entry's attributes under its index. One inside a
		// singleton is not yet. The only resource with one is
		// verity_switchpoint, whose auto-assignment pairs do not all follow the
		// shared rule, so serving it waits on that decision rather than on this.
		if container == spec.FieldKindObject && field.Nullable {
			return fmt.Errorf("%s is a nullable member of a singleton block, which the engine does not serve yet", field.APIName)
		}
		return nil
	case spec.FieldKindList:
		if nested {
			return fmt.Errorf("%s is a list inside a block, which the engine does not serve yet", field.APIName)
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
		for _, member := range field.Fields {
			if member.Unmanaged {
				continue
			}
			if err := supportedField(member, spec.FieldKindList); err != nil {
				return err
			}
		}
		return nil
	case spec.FieldKindObject:
		if nested {
			return fmt.Errorf("%s is an object inside an object, which the engine does not serve yet", field.APIName)
		}
		if field.Collection == nil || field.Collection.Strategy != spec.CollectionSingleton {
			return fmt.Errorf("%s is an object without the singleton strategy, which the engine does not serve", field.APIName)
		}
		if len(field.Fields) == 0 {
			return fmt.Errorf("%s is an object with no members", field.APIName)
		}
		for _, member := range field.Fields {
			if member.Unmanaged {
				continue
			}
			if err := supportedField(member, spec.FieldKindObject); err != nil {
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
