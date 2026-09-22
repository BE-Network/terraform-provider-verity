package provider

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-verity/internal/genericresource"
)

func TestGeneratedResourceKeysCoverEveryResource(t *testing.T) {
	registry := generatedRegistry(t)
	if len(generatedResourceKeys) != len(registry) {
		t.Fatalf("key table has %d entries, the registry has %d", len(generatedResourceKeys), len(registry))
	}
	for _, resourceSpec := range registry {
		keys, exists := generatedResourceKeys[resourceSpec.TerraformType]
		if !exists {
			t.Errorf("no generated keys for %s", resourceSpec.TerraformType)
			continue
		}
		if keys.Endpoint == "" || keys.CacheKey == "" || keys.ResponseCollectionKey == "" {
			t.Errorf("%s has an incomplete key entry: %+v", resourceSpec.TerraformType, keys)
		}
		if want := strings.TrimPrefix(resourceSpec.API.EndpointPath, "/"); keys.Endpoint != want {
			t.Errorf("%s endpoint = %q, registry says %q", resourceSpec.TerraformType, keys.Endpoint, want)
		}
		if keys.CacheKey != resourceSpec.API.CacheKey {
			t.Errorf("%s cache key = %q, registry says %q", resourceSpec.TerraformType, keys.CacheKey, resourceSpec.API.CacheKey)
		}
	}
}

func TestRegistrationFollowsTheRegistry(t *testing.T) {
	registry := generatedRegistry(t)
	if len(generatedResourceOrder) != len(registry) {
		t.Fatalf("registration order lists %d resources, the registry has %d", len(generatedResourceOrder), len(registry))
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

	if !sort.StringsAreSorted(generatedResourceOrder) {
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
