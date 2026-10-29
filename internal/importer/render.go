package importer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/utils"
)

func fieldsByAPIName(fields []spec.FieldSpec) map[string]spec.FieldSpec {
	byAPIName := make(map[string]spec.FieldSpec, len(fields))
	for _, field := range fields {
		if !field.Unmanaged {
			byAPIName[field.APIName] = field
		}
	}
	return byAPIName
}

func naturalNameLess(left, right string) bool {
	leftParts, rightParts := getNaturalSortParts(left), getNaturalSortParts(right)
	for n := 0; n < len(leftParts) && n < len(rightParts); n++ {
		leftNumber, leftIsNumber := leftParts[n].(int)
		rightNumber, rightIsNumber := rightParts[n].(int)
		if leftIsNumber && rightIsNumber {
			if leftNumber != rightNumber {
				return leftNumber < rightNumber
			}
		} else if leftIsNumber != rightIsNumber {
			return leftIsNumber
		} else if leftParts[n].(string) != rightParts[n].(string) {
			return leftParts[n].(string) < rightParts[n].(string)
		}
	}
	if len(leftParts) != len(rightParts) {
		return len(leftParts) < len(rightParts)
	}
	return left < right
}

func (i *Importer) generateResourceTF(objects map[string]map[string]interface{}, config ResourceConfig) (string, []ImportedResource, error) {
	names := make([]string, 0, len(objects))
	for name := range objects {
		names = append(names, name)
	}
	sort.Slice(names, func(a, b int) bool { return naturalNameLess(names[a], names[b]) })
	var output strings.Builder
	resources := make([]ImportedResource, 0, len(names))
	for _, name := range names {
		resource := ImportedResource{
			TerraformType: "verity_" + config.ResourceType,
			TerraformName: utils.SanitizeResourceName(name),
			ID:            name,
		}
		fmt.Fprintf(&output, "\nresource \"%s\" \"%s\" {\n", resource.TerraformType, resource.TerraformName)
		nameValue, _ := scalarHCL(name)
		fmt.Fprintf(&output, "    name = %s\n", nameValue)
		fmt.Fprintf(&output, "    depends_on = [verity_operation_stage.%s]\n", config.StageName)
		if err := renderFields(&output, objects[name], config.Fields, config.SkipTopLevelKeys, "\t", "", true); err != nil {
			return "", nil, fmt.Errorf("verity_%s %q: %w", config.ResourceType, name, err)
		}
		output.WriteString("}\n\n")
		resources = append(resources, resource)
	}
	return output.String(), resources, nil
}

func renderFields(output *strings.Builder, object map[string]interface{}, fields []spec.FieldSpec, skip map[string]bool, indent, identity string, objectsFirst bool) error {
	byAPIName := fieldsByAPIName(fields)
	keys := make([]string, 0, len(object))
	for key := range object {
		field, found := byAPIName[key]
		if !found || skip[key] || skip[field.TerraformName] {
			continue
		}
		if field.AutoAssignment != nil {
			for flagKey, flagField := range byAPIName {
				if flagField.TerraformName == field.AutoAssignment.FlagField && object[flagKey] == true {
					found = false
					break
				}
			}
			if !found {
				continue
			}
		}
		keys = append(keys, key)
	}
	sort.Slice(keys, func(a, b int) bool {
		left, right := byAPIName[keys[a]], byAPIName[keys[b]]
		if (left.TerraformName == identity) != (right.TerraformName == identity) {
			return left.TerraformName == identity
		}
		if objectsFirst && (left.Kind == spec.FieldKindObject) != (right.Kind == spec.FieldKindObject) {
			return left.Kind == spec.FieldKindObject
		}
		return keys[a] < keys[b]
	})
	for _, key := range keys {
		if err := renderField(output, byAPIName[key], object[key], indent); err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}
	return nil
}

func renderField(output *strings.Builder, field spec.FieldSpec, value interface{}, indent string) error {
	switch {
	case field.Kind == spec.FieldKindObject:
		object, ok := value.(map[string]interface{})
		if value != nil && !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
		return renderBlock(output, field, object, indent)
	case field.Kind == spec.FieldKindList && (field.ElementKind == spec.FieldKindObject || len(field.Fields) > 0):
		if value == nil {
			return nil
		}
		entries, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("expected list of objects, got %T", value)
		}
		for index, entry := range entries {
			object, ok := entry.(map[string]interface{})
			if !ok {
				return fmt.Errorf("entry %d: expected object, got %T", index, entry)
			}
			if err := renderBlock(output, field, object, indent); err != nil {
				return fmt.Errorf("entry %d: %w", index, err)
			}
		}
		return nil
	case field.Kind == spec.FieldKindList && value != nil:
		entries, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("expected scalar list, got %T", value)
		}
		if len(entries) == 0 {
			fmt.Fprintf(output, "%s%s = []\n", indent, field.TerraformName)
			return nil
		}
		fmt.Fprintf(output, "%s%s = [\n", indent, field.TerraformName)
		for index, entry := range entries {
			formatted, err := scalarHCL(entry)
			if err != nil {
				return fmt.Errorf("entry %d: %w", index, err)
			}
			fmt.Fprintf(output, "%s\t%s,\n", indent, formatted)
		}
		fmt.Fprintf(output, "%s]\n", indent)
		return nil
	default:
		formatted, err := scalarHCL(value)
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%s%s = %s\n", indent, field.TerraformName, formatted)
		return nil
	}
}

func renderBlock(output *strings.Builder, field spec.FieldSpec, object map[string]interface{}, indent string) error {
	var contents strings.Builder
	identity := ""
	if field.Collection != nil {
		identity = field.Collection.IdentityField
	}
	if err := renderFields(&contents, object, field.Fields, nil, indent+"\t", identity, false); err != nil {
		return err
	}
	if contents.Len() == 0 && field.Kind == spec.FieldKindObject {
		fmt.Fprintf(output, "%s%s {}\n", indent, field.TerraformName)
		return nil
	}
	fmt.Fprintf(output, "%s%s {\n", indent, field.TerraformName)
	output.WriteString(contents.String())
	fmt.Fprintf(output, "%s}\n", indent)
	return nil
}

func scalarHCL(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return string(hclwrite.TokensForValue(cty.StringVal(v)).Bytes()), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v)), nil
		}
		return fmt.Sprintf("%g", v), nil
	case nil:
		return "null", nil
	default:
		return "", fmt.Errorf("expected JSON scalar, got %T", value)
	}
}
