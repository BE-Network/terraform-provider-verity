package genericresource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"terraform-provider-verity/internal/spec"
)

// CompileSchema turns a reviewed ResourceSpec into the Terraform schema the
// provider serves. It is the first half of the generic engine: the half that
// decides what a configuration may say, before anything decides what to send.
//
// Phase 2 is the scalar pilot, so this compiles scalar fields only and refuses
// anything else rather than emitting an approximation. An object or list needs
// the collection strategy from CollectionSpec, which Phase 4 implements; a spec
// that reaches here carrying one is a spec this engine cannot serve yet, and
// saying so is the whole point of refusing.
func CompileSchema(resource spec.ResourceSpec) (schema.Schema, error) {
	attributes := make(map[string]schema.Attribute, len(resource.Fields))
	for _, field := range resource.Fields {
		attribute, err := compileAttribute(field)
		if err != nil {
			return schema.Schema{}, fmt.Errorf("%s.%s: %w", resource.TerraformType, field.TerraformName, err)
		}
		if _, duplicate := attributes[field.TerraformName]; duplicate {
			return schema.Schema{}, fmt.Errorf("%s: duplicate Terraform attribute %q", resource.TerraformType, field.TerraformName)
		}
		attributes[field.TerraformName] = attribute
	}
	if len(attributes) == 0 {
		return schema.Schema{}, fmt.Errorf("%s: no fields to compile", resource.TerraformType)
	}
	return schema.Schema{
		Description: resource.Description,
		Version:     resource.SchemaVersion,
		Attributes:  attributes,
	}, nil
}

// compileAttribute maps one field's declared access, kind, and replacement rule
// onto a Framework attribute. Every branch is driven by the spec; nothing here
// consults the field's name.
func compileAttribute(field spec.FieldSpec) (schema.Attribute, error) {
	if field.Unmanaged {
		return nil, fmt.Errorf("unmanaged fields have no schema attribute")
	}
	required, optional, computed, err := accessFlags(field.Access)
	if err != nil {
		return nil, err
	}
	// A Required attribute is always supplied by the configuration, so a plan
	// modifier that fills it in would never fire and a default would contradict
	// it. Catching it here keeps an override that sets both from producing a
	// schema Terraform would reject at a less obvious moment.
	if required && field.Default != nil {
		return nil, fmt.Errorf("a required attribute cannot carry a default")
	}

	switch field.Kind {
	case spec.FieldKindString:
		attribute := schema.StringAttribute{
			Description: field.Description,
			Required:    required,
			Optional:    optional,
			Computed:    computed,
			Sensitive:   field.Sensitive,
		}
		if field.Replace {
			attribute.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		}
		return attribute, nil

	case spec.FieldKindBool:
		attribute := schema.BoolAttribute{
			Description: field.Description,
			Required:    required,
			Optional:    optional,
			Computed:    computed,
			Sensitive:   field.Sensitive,
		}
		if field.Replace {
			attribute.PlanModifiers = []planmodifier.Bool{boolplanmodifier.RequiresReplace()}
		}
		if field.Default != nil {
			value, err := boolLiteral(*field.Default)
			if err != nil {
				return nil, err
			}
			attribute.Default = booldefault.StaticBool(value)
		}
		return attribute, nil

	case spec.FieldKindInt64:
		attribute := schema.Int64Attribute{
			Description: field.Description,
			Required:    required,
			Optional:    optional,
			Computed:    computed,
			Sensitive:   field.Sensitive,
		}
		if field.Replace {
			attribute.PlanModifiers = []planmodifier.Int64{int64planmodifier.RequiresReplace()}
		}
		return attribute, nil

	case spec.FieldKindNumber:
		attribute := schema.NumberAttribute{
			Description: field.Description,
			Required:    required,
			Optional:    optional,
			Computed:    computed,
			Sensitive:   field.Sensitive,
		}
		if field.Replace {
			attribute.PlanModifiers = []planmodifier.Number{numberplanmodifier.RequiresReplace()}
		}
		return attribute, nil

	case spec.FieldKindObject, spec.FieldKindList:
		return nil, fmt.Errorf("kind %q is a collection, which the scalar engine does not compile; "+
			"it needs the strategy in CollectionSpec", field.Kind)

	default:
		return nil, fmt.Errorf("unknown field kind %q", field.Kind)
	}
}

func accessFlags(access spec.Access) (required, optional, computed bool, err error) {
	switch access {
	case spec.AccessRequired:
		return true, false, false, nil
	case spec.AccessOptional:
		return false, true, false, nil
	case spec.AccessComputed:
		return false, false, true, nil
	case spec.AccessOptionalComputed:
		return false, true, true, nil
	default:
		return false, false, false, fmt.Errorf("unknown access %q", access)
	}
}

func boolLiteral(literal spec.LiteralSpec) (bool, error) {
	if literal.Kind != spec.LiteralBool {
		return false, fmt.Errorf("a bool attribute cannot default to a %s literal", literal.Kind)
	}
	switch literal.Value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("bool default %q is neither true nor false", literal.Value)
	}
}
