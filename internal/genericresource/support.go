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
//   - an indexed collection needs the add/update/delete strategy Phase 4
//     implements, and treating it as a value would send the whole list;
//   - an auto-assignment pair has a flag that changes whether its value is sent
//     at all, and treating the flag as an ordinary bool would send both;
//   - a fixed header selects between resources sharing one endpoint, and the
//     write path does not yet pass it, so a write would reach the wrong objects;
//   - a nullable member inside a singleton needs the configuration scan at a
//     nested path, which the engine reads only at the top level.
func Supported(resource spec.ResourceSpec) error {
	if len(resource.API.FixedHeaders) != 0 {
		return fmt.Errorf("it is selected by a fixed header, which the write path does not pass yet")
	}
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
	if field.AutoAssignment != nil {
		return fmt.Errorf("%s is an auto-assignment pair, which the engine does not implement yet", field.APIName)
	}
	switch field.Kind {
	case spec.FieldKindString, spec.FieldKindBool, spec.FieldKindInt64, spec.FieldKindNumber:
		if nested && field.Nullable {
			return fmt.Errorf("%s is a nullable member of an object, which needs a nested configuration scan the engine does not do yet", field.APIName)
		}
		return nil
	case spec.FieldKindList:
		return fmt.Errorf("%s is an indexed collection, which the engine does not serve yet", field.APIName)
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
