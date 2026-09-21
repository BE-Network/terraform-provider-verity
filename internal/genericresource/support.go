package genericresource

import (
	"fmt"

	"terraform-provider-verity/internal/spec"
)

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

func supportedField(field spec.FieldSpec, containers []spec.FieldKind) error {
	nested := len(containers) != 0
	insideSingleton := nested && containers[0] == spec.FieldKindObject
	if field.AutoAssignment != nil && nested {
		return fmt.Errorf("%s is an auto-assignment pair inside an object, which the engine does not implement yet", field.APIName)
	}
	switch field.Kind {
	case spec.FieldKindString, spec.FieldKindBool, spec.FieldKindInt64, spec.FieldKindNumber:

		if len(containers) > 1 && field.Nullable {
			return fmt.Errorf("%s is a nullable member of a list inside a singleton block, which the engine does not serve", field.APIName)
		}
		return nil
	case spec.FieldKindList:

		if nested && !(len(containers) == 1 && insideSingleton) {
			return fmt.Errorf("%s is a list nested more deeply than inside a singleton, which the engine does not serve", field.APIName)
		}

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
