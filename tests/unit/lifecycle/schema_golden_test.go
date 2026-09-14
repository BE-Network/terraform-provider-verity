package lifecycle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	verityprovider "terraform-provider-verity/internal/provider"
)

const schemaGoldenPath = "../testdata/schema_golden_v6_6.json"

type schemaGolden struct {
	FormatVersion int                    `json:"format_version"`
	Resources     []resourceSchemaGolden `json:"resources"`
}

type resourceSchemaGolden struct {
	TerraformType string              `json:"terraform_type"`
	SchemaVersion int64               `json:"schema_version"`
	StateType     string              `json:"state_type"`
	Attributes    []attributeGolden   `json:"attributes"`
	Blocks        []blockSchemaGolden `json:"blocks"`
}

type attributeGolden struct {
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	ElementType   string            `json:"element_type,omitempty"`
	Required      bool              `json:"required"`
	Optional      bool              `json:"optional"`
	Computed      bool              `json:"computed"`
	Sensitive     bool              `json:"sensitive"`
	PlanModifiers []extensionGolden `json:"plan_modifiers,omitempty"`
	Validators    []extensionGolden `json:"validators,omitempty"`
}

type blockSchemaGolden struct {
	Name          string              `json:"name"`
	Nesting       string              `json:"nesting"`
	Attributes    []attributeGolden   `json:"attributes"`
	Blocks        []blockSchemaGolden `json:"blocks,omitempty"`
	PlanModifiers []extensionGolden   `json:"plan_modifiers,omitempty"`
	Validators    []extensionGolden   `json:"validators,omitempty"`
}

type extensionGolden struct {
	Type       string `json:"type"`
	Parameters string `json:"parameters"`
}

func TestSchemaGolden(t *testing.T) {
	got, err := marshalSchemaGolden(buildSchemaGolden(t))
	if err != nil {
		t.Fatalf("marshal schema golden: %v", err)
	}

	if os.Getenv("UPDATE_SCHEMA_SNAPSHOT") == "1" {
		if err := os.WriteFile(schemaGoldenPath, got, 0o644); err != nil {
			t.Fatalf("write schema golden: %v", err)
		}
		return
	}

	want, err := os.ReadFile(schemaGoldenPath)
	if err != nil {
		t.Fatalf("read schema golden: %v\nRegenerate intentionally with UPDATE_SCHEMA_SNAPSHOT=1.", err)
	}
	if !bytes.Equal(want, got) {
		t.Fatalf("version-zero resource schema contract changed; inspect the diff and regenerate intentionally with UPDATE_SCHEMA_SNAPSHOT=1\nwant: %s\ngot: %s", schemaGoldenPath, string(got))
	}
}

func buildSchemaGolden(t *testing.T) schemaGolden {
	t.Helper()
	ctx := context.Background()
	provider := verityprovider.New("schema-golden")()
	factories := provider.Resources(ctx)
	resources := make([]resourceSchemaGolden, 0, len(factories))
	seen := make(map[string]struct{}, len(factories))

	for _, factory := range factories {
		res := factory()
		var metadata resource.MetadataResponse
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		if metadata.TypeName == "" {
			t.Fatalf("resource %T did not return a Terraform type name", res)
		}
		if _, exists := seen[metadata.TypeName]; exists {
			t.Fatalf("duplicate Terraform resource type %q", metadata.TypeName)
		}
		seen[metadata.TypeName] = struct{}{}

		var schemaResponse resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
		if schemaResponse.Diagnostics.HasError() {
			t.Fatalf("read schema for %s: %s", metadata.TypeName, schemaResponse.Diagnostics)
		}
		resources = append(resources, resourceSchemaGolden{
			TerraformType: metadata.TypeName,
			SchemaVersion: schemaResponse.Schema.Version,
			StateType:     fmt.Sprint(schemaResponse.Schema.Type().TerraformType(ctx)),
			Attributes:    goldenAttributes(t, schemaResponse.Schema.Attributes),
			Blocks:        goldenBlocks(t, schemaResponse.Schema.Blocks),
		})
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].TerraformType < resources[j].TerraformType
	})
	return schemaGolden{FormatVersion: 1, Resources: resources}
}

func goldenAttributes(t *testing.T, attributes map[string]fwschema.Attribute) []attributeGolden {
	t.Helper()
	result := make([]attributeGolden, 0, len(attributes))
	for name, attribute := range attributes {
		result = append(result, goldenAttribute(t, name, attribute))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func goldenAttribute(t *testing.T, name string, attribute fwschema.Attribute) attributeGolden {
	t.Helper()
	base := attributeGolden{Name: name}
	switch value := attribute.(type) {
	case fwschema.StringAttribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "string", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = goldenExtensions(value.PlanModifiers), goldenExtensions(value.Validators)
	case fwschema.BoolAttribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "bool", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = goldenExtensions(value.PlanModifiers), goldenExtensions(value.Validators)
	case fwschema.Int64Attribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "int64", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = goldenExtensions(value.PlanModifiers), goldenExtensions(value.Validators)
	case fwschema.NumberAttribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "number", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = goldenExtensions(value.PlanModifiers), goldenExtensions(value.Validators)
	case fwschema.ListAttribute:
		base.Type, base.ElementType, base.Required, base.Optional, base.Computed, base.Sensitive = "list", fmt.Sprint(value.ElementType), value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = goldenExtensions(value.PlanModifiers), goldenExtensions(value.Validators)
	default:
		t.Fatalf("unsupported attribute %q of type %T; extend the golden contract", name, attribute)
	}
	return base
}

func goldenBlocks(t *testing.T, blocks map[string]fwschema.Block) []blockSchemaGolden {
	t.Helper()
	result := make([]blockSchemaGolden, 0, len(blocks))
	for name, block := range blocks {
		result = append(result, goldenBlock(t, name, block))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func goldenBlock(t *testing.T, name string, block fwschema.Block) blockSchemaGolden {
	t.Helper()
	golden := blockSchemaGolden{Name: name}
	switch value := block.(type) {
	case fwschema.ListNestedBlock:
		golden.Nesting = "list"
		golden.Attributes = goldenAttributes(t, value.NestedObject.Attributes)
		golden.Blocks = goldenBlocks(t, value.NestedObject.Blocks)
		golden.PlanModifiers, golden.Validators = goldenExtensions(value.PlanModifiers), goldenExtensions(value.Validators)
	case fwschema.SingleNestedBlock:
		golden.Nesting = "single"
		golden.Attributes = goldenAttributes(t, value.Attributes)
		golden.Blocks = goldenBlocks(t, value.Blocks)
		golden.PlanModifiers, golden.Validators = goldenExtensions(value.PlanModifiers), goldenExtensions(value.Validators)
	default:
		t.Fatalf("unsupported block %q of type %T; extend the golden contract", name, block)
	}
	return golden
}

func goldenExtensions[T any](values []T) []extensionGolden {
	if len(values) == 0 {
		return nil
	}
	result := make([]extensionGolden, 0, len(values))
	for _, value := range values {
		result = append(result, extensionGolden{Type: fmt.Sprintf("%T", value), Parameters: goldenParameterValue(reflect.ValueOf(value))})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type == result[j].Type {
			return result[i].Parameters < result[j].Parameters
		}
		return result[i].Type < result[j].Type
	})
	return result
}

func goldenParameterValue(value reflect.Value) string {
	if !value.IsValid() {
		return "null"
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return "null"
		}
		return goldenParameterValue(value.Elem())
	}
	switch value.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, value.Type().Bits())
	case reflect.String:
		return strconv.Quote(value.String())
	case reflect.Ptr:
		if value.IsNil() {
			return "null"
		}
		return "&" + goldenParameterValue(value.Elem())
	case reflect.Func:
		return "<func>"
	case reflect.Slice, reflect.Array:
		parts := make([]string, value.Len())
		for index := range parts {
			parts[index] = goldenParameterValue(value.Index(index))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case reflect.Map:
		keys := value.MapKeys()
		sort.Slice(keys, func(i, j int) bool { return goldenParameterValue(keys[i]) < goldenParameterValue(keys[j]) })
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, goldenParameterValue(key)+":"+goldenParameterValue(value.MapIndex(key)))
		}
		return "{" + strings.Join(parts, ",") + "}"
	case reflect.Struct:
		parts := make([]string, 0, value.NumField())
		for index := 0; index < value.NumField(); index++ {
			parts = append(parts, value.Type().Field(index).Name+":"+goldenParameterValue(value.Field(index)))
		}
		return value.Type().String() + "{" + strings.Join(parts, ",") + "}"
	default:
		return value.Type().String() + "<" + value.Kind().String() + ">"
	}
}

func marshalSchemaGolden(golden schemaGolden) ([]byte, error) {
	encoded, err := json.MarshalIndent(golden, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}

func TestGoldenParameterValueIsDeterministicForFunctions(t *testing.T) {
	type parameterFixture struct {
		Callback func()
		Label    string
		Values   []int
	}
	first := parameterFixture{Callback: func() {}, Label: "stable", Values: []int{2, 1}}
	second := parameterFixture{Callback: func() {}, Label: "stable", Values: []int{2, 1}}
	firstGolden := goldenParameterValue(reflect.ValueOf(first))
	secondGolden := goldenParameterValue(reflect.ValueOf(second))
	if firstGolden != secondGolden || strings.Contains(firstGolden, "0x") {
		t.Fatalf("non-deterministic function parameter golden: %q vs %q", firstGolden, secondGolden)
	}
}
