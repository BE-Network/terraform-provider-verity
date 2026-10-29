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

func CompileSchema(resource spec.ResourceSpec) (schema.Schema, error) {
	if err := Supported(resource); err != nil {
		return schema.Schema{}, fmt.Errorf("%s: %w", resource.TerraformType, err)
	}
	attributes := make(map[string]schema.Attribute, len(resource.Fields))
	blocks := make(map[string]schema.Block)
	for _, field := range resource.Fields {
		if field.Unmanaged {

			continue
		}
		if _, duplicate := attributes[field.TerraformName]; duplicate {
			return schema.Schema{}, fmt.Errorf("%s: duplicate Terraform attribute %q", resource.TerraformType, field.TerraformName)
		}
		if _, duplicate := blocks[field.TerraformName]; duplicate {
			return schema.Schema{}, fmt.Errorf("%s: duplicate Terraform block %q", resource.TerraformType, field.TerraformName)
		}
		if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
			block, err := compileSingletonBlock(field)
			if err != nil {
				return schema.Schema{}, fmt.Errorf("%s.%s: %w", resource.TerraformType, field.TerraformName, err)
			}
			blocks[field.TerraformName] = block
			continue
		}
		attribute, err := compileAttribute(field)
		if err != nil {
			return schema.Schema{}, fmt.Errorf("%s.%s: %w", resource.TerraformType, field.TerraformName, err)
		}
		attributes[field.TerraformName] = attribute
	}
	if len(attributes) == 0 && len(blocks) == 0 {
		return schema.Schema{}, fmt.Errorf("%s: no fields to compile", resource.TerraformType)
	}
	compiled := schema.Schema{
		Description: resource.Description,
		Version:     resource.SchemaVersion,
		Attributes:  attributes,
	}
	if len(blocks) != 0 {
		compiled.Blocks = blocks
	}
	return compiled, nil
}

func compileSingletonBlock(field spec.FieldSpec) (schema.Block, error) {
	members := make(map[string]schema.Attribute, len(field.Fields))
	var blocks map[string]schema.Block
	for _, member := range field.Fields {
		if member.Unmanaged {
			continue
		}
		if member.Kind == spec.FieldKindList {

			block, err := compileSingletonBlock(member)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", member.TerraformName, err)
			}
			if blocks == nil {
				blocks = make(map[string]schema.Block)
			}
			blocks[member.TerraformName] = block
			continue
		}
		attribute, err := compileAttribute(member)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", member.TerraformName, err)
		}
		members[member.TerraformName] = attribute
	}
	return schema.ListNestedBlock{
		Description:  field.Description,
		NestedObject: schema.NestedBlockObject{Attributes: members, Blocks: blocks},
	}, nil
}

func compileAttribute(field spec.FieldSpec) (schema.Attribute, error) {
	if field.Unmanaged {
		return nil, fmt.Errorf("unmanaged fields have no schema attribute")
	}
	required, optional, computed, err := accessFlags(field.Access)
	if err != nil {
		return nil, err
	}

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

		return nil, fmt.Errorf("kind %q is a collection, which is compiled as a block rather than an attribute", field.Kind)

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
