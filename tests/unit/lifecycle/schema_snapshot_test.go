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

const schemaSnapshotPath = "../testdata/schema_snapshot_v6_6.json"

type schemaSnapshot struct {
	FormatVersion int                      `json:"format_version"`
	Resources     []resourceSchemaSnapshot `json:"resources"`
}

type resourceSchemaSnapshot struct {
	TerraformType string                `json:"terraform_type"`
	SchemaVersion int64                 `json:"schema_version"`
	StateType     string                `json:"state_type"`
	Attributes    []attributeSnapshot   `json:"attributes"`
	Blocks        []blockSchemaSnapshot `json:"blocks"`
}

type attributeSnapshot struct {
	Name          string              `json:"name"`
	Type          string              `json:"type"`
	ElementType   string              `json:"element_type,omitempty"`
	Required      bool                `json:"required"`
	Optional      bool                `json:"optional"`
	Computed      bool                `json:"computed"`
	Sensitive     bool                `json:"sensitive"`
	PlanModifiers []extensionSnapshot `json:"plan_modifiers,omitempty"`
	Validators    []extensionSnapshot `json:"validators,omitempty"`
}

type blockSchemaSnapshot struct {
	Name          string                `json:"name"`
	Nesting       string                `json:"nesting"`
	Attributes    []attributeSnapshot   `json:"attributes"`
	Blocks        []blockSchemaSnapshot `json:"blocks,omitempty"`
	PlanModifiers []extensionSnapshot   `json:"plan_modifiers,omitempty"`
	Validators    []extensionSnapshot   `json:"validators,omitempty"`
}

type extensionSnapshot struct {
	Type       string `json:"type"`
	Parameters string `json:"parameters"`
}

func TestSchemaSnapshot(t *testing.T) {
	got, err := marshalSchemaSnapshot(buildSchemaSnapshot(t))
	if err != nil {
		t.Fatalf("marshal schema snapshot: %v", err)
	}

	if os.Getenv("UPDATE_SCHEMA_SNAPSHOT") == "1" {
		if err := os.WriteFile(schemaSnapshotPath, got, 0o644); err != nil {
			t.Fatalf("write schema snapshot: %v", err)
		}
		return
	}

	want, err := os.ReadFile(schemaSnapshotPath)
	if err != nil {
		t.Fatalf("read schema snapshot: %v\nRegenerate intentionally with UPDATE_SCHEMA_SNAPSHOT=1.", err)
	}
	if !bytes.Equal(want, got) {
		t.Fatalf("version-zero resource schema contract changed; inspect the diff and regenerate intentionally with UPDATE_SCHEMA_SNAPSHOT=1\nwant: %s\ngot: %s", schemaSnapshotPath, string(got))
	}
}

func buildSchemaSnapshot(t *testing.T) schemaSnapshot {
	t.Helper()
	ctx := context.Background()
	provider := verityprovider.New("schema-snapshot")()
	factories := provider.Resources(ctx)
	resources := make([]resourceSchemaSnapshot, 0, len(factories))
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
		resources = append(resources, resourceSchemaSnapshot{
			TerraformType: metadata.TypeName,
			SchemaVersion: schemaResponse.Schema.Version,
			StateType:     fmt.Sprint(schemaResponse.Schema.Type().TerraformType(ctx)),
			Attributes:    snapshotAttributes(t, schemaResponse.Schema.Attributes),
			Blocks:        snapshotBlocks(t, schemaResponse.Schema.Blocks),
		})
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].TerraformType < resources[j].TerraformType
	})
	return schemaSnapshot{FormatVersion: 1, Resources: resources}
}

func snapshotAttributes(t *testing.T, attributes map[string]fwschema.Attribute) []attributeSnapshot {
	t.Helper()
	result := make([]attributeSnapshot, 0, len(attributes))
	for name, attribute := range attributes {
		result = append(result, snapshotAttribute(t, name, attribute))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func snapshotAttribute(t *testing.T, name string, attribute fwschema.Attribute) attributeSnapshot {
	t.Helper()
	base := attributeSnapshot{Name: name}
	switch value := attribute.(type) {
	case fwschema.StringAttribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "string", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = snapshotExtensions(value.PlanModifiers), snapshotExtensions(value.Validators)
	case fwschema.BoolAttribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "bool", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = snapshotExtensions(value.PlanModifiers), snapshotExtensions(value.Validators)
	case fwschema.Int64Attribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "int64", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = snapshotExtensions(value.PlanModifiers), snapshotExtensions(value.Validators)
	case fwschema.NumberAttribute:
		base.Type, base.Required, base.Optional, base.Computed, base.Sensitive = "number", value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = snapshotExtensions(value.PlanModifiers), snapshotExtensions(value.Validators)
	case fwschema.ListAttribute:
		base.Type, base.ElementType, base.Required, base.Optional, base.Computed, base.Sensitive = "list", fmt.Sprint(value.ElementType), value.Required, value.Optional, value.Computed, value.Sensitive
		base.PlanModifiers, base.Validators = snapshotExtensions(value.PlanModifiers), snapshotExtensions(value.Validators)
	default:
		t.Fatalf("unsupported attribute %q of type %T; extend the snapshot contract", name, attribute)
	}
	return base
}

func snapshotBlocks(t *testing.T, blocks map[string]fwschema.Block) []blockSchemaSnapshot {
	t.Helper()
	result := make([]blockSchemaSnapshot, 0, len(blocks))
	for name, block := range blocks {
		result = append(result, snapshotBlock(t, name, block))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func snapshotBlock(t *testing.T, name string, block fwschema.Block) blockSchemaSnapshot {
	t.Helper()
	snapshot := blockSchemaSnapshot{Name: name}
	switch value := block.(type) {
	case fwschema.ListNestedBlock:
		snapshot.Nesting = "list"
		snapshot.Attributes = snapshotAttributes(t, value.NestedObject.Attributes)
		snapshot.Blocks = snapshotBlocks(t, value.NestedObject.Blocks)
		snapshot.PlanModifiers, snapshot.Validators = snapshotExtensions(value.PlanModifiers), snapshotExtensions(value.Validators)
	case fwschema.SingleNestedBlock:
		snapshot.Nesting = "single"
		snapshot.Attributes = snapshotAttributes(t, value.Attributes)
		snapshot.Blocks = snapshotBlocks(t, value.Blocks)
		snapshot.PlanModifiers, snapshot.Validators = snapshotExtensions(value.PlanModifiers), snapshotExtensions(value.Validators)
	default:
		t.Fatalf("unsupported block %q of type %T; extend the snapshot contract", name, block)
	}
	return snapshot
}

func snapshotExtensions[T any](values []T) []extensionSnapshot {
	if len(values) == 0 {
		return nil
	}
	result := make([]extensionSnapshot, 0, len(values))
	for _, value := range values {
		result = append(result, extensionSnapshot{Type: fmt.Sprintf("%T", value), Parameters: snapshotParameterValue(reflect.ValueOf(value))})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type == result[j].Type {
			return result[i].Parameters < result[j].Parameters
		}
		return result[i].Type < result[j].Type
	})
	return result
}

func snapshotParameterValue(value reflect.Value) string {
	if !value.IsValid() {
		return "null"
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return "null"
		}
		return snapshotParameterValue(value.Elem())
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
		return "&" + snapshotParameterValue(value.Elem())
	case reflect.Func:
		return "<func>"
	case reflect.Slice, reflect.Array:
		parts := make([]string, value.Len())
		for index := range parts {
			parts[index] = snapshotParameterValue(value.Index(index))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case reflect.Map:
		keys := value.MapKeys()
		sort.Slice(keys, func(i, j int) bool { return snapshotParameterValue(keys[i]) < snapshotParameterValue(keys[j]) })
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, snapshotParameterValue(key)+":"+snapshotParameterValue(value.MapIndex(key)))
		}
		return "{" + strings.Join(parts, ",") + "}"
	case reflect.Struct:
		parts := make([]string, 0, value.NumField())
		for index := 0; index < value.NumField(); index++ {
			parts = append(parts, value.Type().Field(index).Name+":"+snapshotParameterValue(value.Field(index)))
		}
		return value.Type().String() + "{" + strings.Join(parts, ",") + "}"
	default:
		return value.Type().String() + "<" + value.Kind().String() + ">"
	}
}

func marshalSchemaSnapshot(snapshot schemaSnapshot) ([]byte, error) {
	encoded, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}

func TestSnapshotParameterValueIsDeterministicForFunctions(t *testing.T) {
	type parameterFixture struct {
		Callback func()
		Label    string
		Values   []int
	}
	first := parameterFixture{Callback: func() {}, Label: "stable", Values: []int{2, 1}}
	second := parameterFixture{Callback: func() {}, Label: "stable", Values: []int{2, 1}}
	firstSnapshot := snapshotParameterValue(reflect.ValueOf(first))
	secondSnapshot := snapshotParameterValue(reflect.ValueOf(second))
	if firstSnapshot != secondSnapshot || strings.Contains(firstSnapshot, "0x") {
		t.Fatalf("non-deterministic function parameter snapshot: %q vs %q", firstSnapshot, secondSnapshot)
	}
}
