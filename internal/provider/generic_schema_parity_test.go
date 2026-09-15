package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/spec"
)

// The generic engine must serve the schema the provider already ships.
//
// Phase 2 migrates verity_ipv4_list to a compiled schema with no handwritten
// attribute map. The registry is verified against the legacy schema elsewhere,
// but that proves the *description* is faithful, not that compiling it back
// produces the same thing. This closes the loop: compile the reviewed spec and
// compare it to what the handwritten resource returns, attribute by attribute.
//
// A schema difference is a state-compatibility break, so this is the check that
// has to pass before the generic resource is allowed to serve real state.
func TestCompiledSchemaMatchesLegacyIPv4List(t *testing.T) {
	t.Parallel()

	resourceSpec := registrySpec(t, "verity_ipv4_list")
	compiled, err := genericresource.CompileSchema(resourceSpec)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	legacy := legacySchema(t, NewVerityIpv4ListResource)

	if compiled.Description != legacy.Description {
		t.Errorf("schema description = %q, legacy is %q", compiled.Description, legacy.Description)
	}
	if compiled.Version != legacy.Version {
		t.Errorf("schema version = %d, legacy is %d", compiled.Version, legacy.Version)
	}
	if len(compiled.Attributes) != len(legacy.Attributes) {
		t.Fatalf("compiled %d attributes, legacy has %d: %v vs %v",
			len(compiled.Attributes), len(legacy.Attributes),
			attributeNames(compiled.Attributes), attributeNames(legacy.Attributes))
	}

	for name, want := range legacy.Attributes {
		got, present := compiled.Attributes[name]
		if !present {
			t.Errorf("compiled schema has no attribute %q", name)
			continue
		}
		compareAttributes(t, name, got, want)
	}
}

// compareAttributes checks the facts Terraform acts on: the value type, the three
// access flags, sensitivity, the description shown to users, and whether a change
// forces replacement. Plan modifiers are compared by identifying RequiresReplace
// rather than by counting, since an unrelated modifier is not replacement.
func compareAttributes(t *testing.T, name string, got, want schema.Attribute) {
	t.Helper()

	if gotType, wantType := got.GetType(), want.GetType(); !gotType.Equal(wantType) {
		t.Errorf("%s: type = %s, legacy is %s", name, gotType, wantType)
	}
	if got.IsRequired() != want.IsRequired() || got.IsOptional() != want.IsOptional() || got.IsComputed() != want.IsComputed() {
		t.Errorf("%s: required/optional/computed = %v/%v/%v, legacy is %v/%v/%v",
			name, got.IsRequired(), got.IsOptional(), got.IsComputed(),
			want.IsRequired(), want.IsOptional(), want.IsComputed())
	}
	if got.IsSensitive() != want.IsSensitive() {
		t.Errorf("%s: sensitive = %v, legacy is %v", name, got.IsSensitive(), want.IsSensitive())
	}
	if got.GetDescription() != want.GetDescription() {
		t.Errorf("%s: description = %q, legacy is %q", name, got.GetDescription(), want.GetDescription())
	}
	if gotReplace, wantReplace := attributeRequiresReplace(got), attributeRequiresReplace(want); gotReplace != wantReplace {
		t.Errorf("%s: requires replace = %v, legacy is %v", name, gotReplace, wantReplace)
	}
}

// attributeRequiresReplace reports whether an attribute carries the
// RequiresReplace plan modifier for its kind. The modifier list is typed per
// kind, so this switches on the concrete attribute type to reach it.
func attributeRequiresReplace(attribute schema.Attribute) bool {
	switch typed := attribute.(type) {
	case schema.StringAttribute:
		return countRequiresReplace(typed.PlanModifiers) > 0
	case schema.BoolAttribute:
		return countRequiresReplace(typed.PlanModifiers) > 0
	case schema.Int64Attribute:
		return countRequiresReplace(typed.PlanModifiers) > 0
	case schema.NumberAttribute:
		return countRequiresReplace(typed.PlanModifiers) > 0
	default:
		return false
	}
}

func attributeNames(attributes map[string]schema.Attribute) []string {
	names := make([]string, 0, len(attributes))
	for name := range attributes {
		names = append(names, name)
	}
	return names
}

func legacySchema(t *testing.T, factory func() resource.Resource) schema.Schema {
	t.Helper()
	var response resource.SchemaResponse
	factory().Schema(context.Background(), resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("legacy schema reported diagnostics: %v", response.Diagnostics)
	}
	return response.Schema
}

func registrySpec(t *testing.T, terraformType string) spec.ResourceSpec {
	t.Helper()
	for _, candidate := range generatedRegistry(t) {
		if candidate.TerraformType == terraformType {
			return candidate
		}
	}
	t.Fatalf("no registry entry for %s", terraformType)
	return spec.ResourceSpec{}
}

// The compiler must refuse what it cannot serve rather than approximate it. A
// collection needs the strategy Phase 4 implements, and silently emitting a bare
// attribute for one would produce a schema that accepts configuration the engine
// then ignores.
func TestCompileSchemaRefusesCollections(t *testing.T) {
	t.Parallel()

	for _, resourceSpec := range generatedRegistry(t) {
		hasCollection := false
		for _, field := range resourceSpec.Fields {
			if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
				hasCollection = true
				break
			}
		}
		if !hasCollection {
			continue
		}
		if _, err := genericresource.CompileSchema(resourceSpec); err == nil {
			t.Errorf("%s carries a collection but compiled without error", resourceSpec.TerraformType)
		}
		return
	}
	t.Fatal("no registry resource carries a collection, so this check proved nothing")
}

// Every scalar-only resource in the registry must compile, not just the pilot.
// If one does not, the engine's reach is narrower than the registry claims and
// the next migration would discover it the hard way.
func TestEveryScalarOnlyResourceCompiles(t *testing.T) {
	t.Parallel()

	compiled := 0
	for _, resourceSpec := range generatedRegistry(t) {
		scalarOnly := true
		for _, field := range resourceSpec.Fields {
			if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
				scalarOnly = false
				break
			}
		}
		if !scalarOnly {
			continue
		}
		if _, err := genericresource.CompileSchema(resourceSpec); err != nil {
			t.Errorf("%s is scalar-only but did not compile: %v", resourceSpec.TerraformType, err)
			continue
		}
		compiled++
	}
	if compiled == 0 {
		t.Fatal("no scalar-only resource was compiled, so this check proved nothing")
	}
	t.Logf("%d scalar-only resources compile from the registry", compiled)
}
