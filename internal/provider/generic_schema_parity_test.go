package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"sort"

	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
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
func TestCompiledSchemaMatchesLegacy(t *testing.T) {
	t.Parallel()

	// Every resource the engine can serve, not just the pilot. Checking one of
	// them left the others' schemas unverified, and a schema difference is what
	// decides whether Terraform plans a change at all.
	servable := make([]string, 0, len(transport.GeneratedAdapters))
	for terraformType := range transport.GeneratedAdapters {
		servable = append(servable, terraformType)
	}
	sort.Strings(servable)
	if len(servable) == 0 {
		t.Fatal("no generated adapters, so this test proved nothing")
	}

	for _, terraformType := range servable {
		t.Run(terraformType, func(t *testing.T) {
			t.Parallel()
			assertSchemaParity(t, terraformType)
		})
	}
}

func assertSchemaParity(t *testing.T, terraformType string) {
	t.Helper()

	resourceSpec := registrySpec(t, terraformType)
	compiled, err := genericresource.CompileSchema(resourceSpec)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	constructor, registered := resourceConstructors[terraformType]
	if !registered {
		t.Fatalf("%s has no handwritten constructor to compare against", terraformType)
	}
	legacy := legacySchema(t, constructor)

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

	// Blocks are compared as strictly as attributes. A block's type decides the
	// state shape, so a singleton compiled to anything but the list block the
	// handwritten resource declares would break existing state.
	if len(compiled.Blocks) != len(legacy.Blocks) {
		t.Fatalf("compiled %d blocks, legacy has %d", len(compiled.Blocks), len(legacy.Blocks))
	}
	for name, want := range legacy.Blocks {
		got, present := compiled.Blocks[name]
		if !present {
			t.Errorf("compiled schema has no block %q", name)
			continue
		}
		compareBlocks(t, name, got, want)
	}
}

func compareBlocks(t *testing.T, name string, got, want schema.Block) {
	t.Helper()

	gotList, gotOK := got.(schema.ListNestedBlock)
	wantList, wantOK := want.(schema.ListNestedBlock)
	if !gotOK || !wantOK {
		t.Errorf("%s: block is %T, legacy is %T; only list blocks are compared", name, got, want)
		return
	}
	if gotList.Description != wantList.Description {
		t.Errorf("%s: description = %q, legacy is %q", name, gotList.Description, wantList.Description)
	}
	if len(gotList.Validators) != len(wantList.Validators) || len(gotList.PlanModifiers) != len(wantList.PlanModifiers) {
		t.Errorf("%s: %d validators and %d plan modifiers, legacy has %d and %d", name,
			len(gotList.Validators), len(gotList.PlanModifiers), len(wantList.Validators), len(wantList.PlanModifiers))
	}
	if gotType, wantType := got.Type(), want.Type(); !gotType.Equal(wantType) {
		t.Errorf("%s: type = %s, legacy is %s", name, gotType, wantType)
	}
	gotMembers, wantMembers := gotList.NestedObject.Attributes, wantList.NestedObject.Attributes
	if len(gotMembers) != len(wantMembers) || len(gotList.NestedObject.Blocks) != len(wantList.NestedObject.Blocks) {
		t.Errorf("%s: %d members and %d nested blocks, legacy has %d and %d", name,
			len(gotMembers), len(gotList.NestedObject.Blocks), len(wantMembers), len(wantList.NestedObject.Blocks))
	}
	for member, wantMember := range wantMembers {
		gotMember, present := gotMembers[member]
		if !present {
			t.Errorf("%s: no member %q", name, member)
			continue
		}
		compareAttributes(t, name+"."+member, gotMember, wantMember)
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

// The compiler must refuse what the engine cannot serve rather than approximate
// it. An indexed collection needs the strategy Phase 4 implements, and silently
// emitting a bare attribute for one would produce a schema that accepts
// configuration the engine then ignores.
func TestCompileSchemaRefusesUnsupportedResources(t *testing.T) {
	t.Parallel()

	refused := 0
	for _, resourceSpec := range generatedRegistry(t) {
		if genericresource.Supported(resourceSpec) == nil {
			continue
		}
		if _, err := genericresource.CompileSchema(resourceSpec); err == nil {
			t.Errorf("%s is unsupported but compiled without error", resourceSpec.TerraformType)
		}
		refused++
	}
	if refused == 0 {
		t.Fatal("no registry resource is unsupported, so this check proved nothing")
	}
}

// Every resource the support rule accepts must compile and have a generated
// adapter. The support rule, the compiler, and the generator are separate code,
// and a resource one of them accepts and another does not is not migratable.
func TestEverySupportedResourceCompilesAndHasAnAdapter(t *testing.T) {
	t.Parallel()

	supported := 0
	for _, resourceSpec := range generatedRegistry(t) {
		if genericresource.Supported(resourceSpec) != nil {
			continue
		}
		supported++
		if _, err := genericresource.CompileSchema(resourceSpec); err != nil {
			t.Errorf("%s is supported but did not compile: %v", resourceSpec.TerraformType, err)
		}
		if _, present := transport.GeneratedAdapters[resourceSpec.TerraformType]; !present {
			t.Errorf("%s is supported but has no generated transport adapter", resourceSpec.TerraformType)
		}
	}
	if supported == 0 {
		t.Fatal("no registry resource is supported, so this check proved nothing")
	}
	if supported != len(transport.GeneratedAdapters) {
		t.Errorf("%d resources are supported but %d adapters were generated", supported, len(transport.GeneratedAdapters))
	}
	t.Logf("%d resources are supported, compile, and have generated adapters", supported)
}
