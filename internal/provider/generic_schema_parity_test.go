package provider

import (
	"testing"

	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

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

func TestCompileSchemaRefusesUnsupportedResources(t *testing.T) {
	t.Parallel()

	for _, resourceSpec := range generatedRegistry(t) {
		if err := genericresource.Supported(resourceSpec); err != nil {
			t.Errorf("%s is unsupported: %v", resourceSpec.TerraformType, err)
		}
	}

	unsupported := registrySpec(t, "verity_switchpoint")
	unsupported.Fields = append([]spec.FieldSpec(nil), unsupported.Fields...)
	for i, field := range unsupported.Fields {
		if field.TerraformName != "object_properties" {
			continue
		}
		block := field
		block.Fields = append(append([]spec.FieldSpec(nil), field.Fields...), spec.FieldSpec{
			TerraformName: "nested", APIName: "nested", Kind: spec.FieldKindObject,
			Access: spec.AccessOptional, Modes: field.Modes,
		})
		unsupported.Fields[i] = block
	}
	if genericresource.Supported(unsupported) == nil {
		t.Fatal("the altered switchpoint is supported, so this check proves nothing")
	}
	if _, err := genericresource.CompileSchema(unsupported); err == nil {
		t.Fatal("an unsupported resource compiled without error")
	}
}

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
