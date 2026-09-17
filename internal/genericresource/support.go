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
//   - a nullable member of a block needs the configuration scan at a nested
//     path, which the engine reads only at the top level.
func Supported(resource spec.ResourceSpec) error {
	for _, field := range resource.Fields {
		if field.Unmanaged {
			continue
		}
		if err := supportedField(field, false); err != nil {
			return err
		}
	}
	return nil
}

func supportedField(field spec.FieldSpec, nested bool) error {
	if field.AutoAssignment != nil && nested {
		return fmt.Errorf("%s is an auto-assignment pair inside an object, which the engine does not implement yet", field.APIName)
	}
	switch field.Kind {
	case spec.FieldKindString, spec.FieldKindBool, spec.FieldKindInt64, spec.FieldKindNumber:
		if nested && field.Nullable {
			return fmt.Errorf("%s is a nullable member of a block, which needs a nested configuration scan the engine does not do yet", field.APIName)
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
			if err := supportedField(member, true); err != nil {
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
			if err := supportedField(member, true); err != nil {
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
