package provider

import (
	"context"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-verity/internal/genericresource"
)

func TestRegistrationFollowsTheRegistry(t *testing.T) {
	registry := generatedRegistry(t)
	order := registryResourceOrder()
	if len(order) != len(registry) {
		t.Fatalf("registration order lists %d resources, the registry has %d", len(order), len(registry))
	}
	inRegistry := make(map[string]bool, len(registry))
	for _, resourceSpec := range registry {
		inRegistry[resourceSpec.TerraformType] = true
	}
	ctx := context.Background()
	nonAPI := make(map[string]bool, len(nonAPIResources))
	for _, constructor := range nonAPIResources {
		var metadata resource.MetadataResponse
		constructor().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		nonAPI[metadata.TypeName] = true
	}
	for _, constructor := range getAllResources() {
		served := constructor()
		var metadata resource.MetadataResponse
		served.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		if nonAPI[metadata.TypeName] {
			continue
		}
		if !inRegistry[metadata.TypeName] {
			t.Errorf("%s is registered but the registry does not describe it", metadata.TypeName)
		}
		if _, generic := served.(*genericresource.Resource); !generic {
			t.Errorf("%s is served by %T; API-backed resources are served by the generic engine", metadata.TypeName, served)
		}
	}

	if !sort.StringsAreSorted(order) {
		t.Error("registration order is not in canonical Terraform-name order")
	}

	seen := map[string]int{}
	for _, constructor := range getAllResources() {
		var metadata resource.MetadataResponse
		constructor().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		seen[metadata.TypeName]++
	}
	for name, count := range seen {
		if count != 1 {
			t.Errorf("%s is registered %d times", name, count)
		}
	}
	if len(seen) != len(registry)+len(nonAPIResources) {
		t.Errorf("provider registers %d resources, want %d registry resources plus %d non-API", len(seen), len(registry), len(nonAPIResources))
	}
}
