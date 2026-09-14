package provider

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TestGeneratedResourceKeysCoverEveryResource checks that the generated key table
// serves every resource the provider registers. The resource implementations now
// take their endpoint name and cache key from it rather than declaring their own
// constants, so a missing entry would leave a resource addressing an empty
// endpoint rather than failing to compile.
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

// TestRegistrationFollowsTheRegistry checks that registration is driven by the
// registry rather than by a handwritten list. Every registry resource needs a
// constructor, every constructor must belong to a registry resource, and each one
// must report the Terraform type it is registered under. A drift in any direction
// would either drop a resource from the provider or register one the registry does
// not describe.
func TestRegistrationFollowsTheRegistry(t *testing.T) {
	registry := generatedRegistry(t)
	if len(generatedResourceOrder) != len(registry) {
		t.Fatalf("registration order lists %d resources, the registry has %d", len(generatedResourceOrder), len(registry))
	}
	if len(resourceConstructors) != len(registry) {
		t.Fatalf("constructor map has %d entries, the registry has %d", len(resourceConstructors), len(registry))
	}

	inRegistry := make(map[string]bool, len(registry))
	for _, resourceSpec := range registry {
		inRegistry[resourceSpec.TerraformType] = true
	}
	ctx := context.Background()
	for terraformType, constructor := range resourceConstructors {
		if !inRegistry[terraformType] {
			t.Errorf("constructor registered for %q, which the registry does not describe", terraformType)
			continue
		}
		var metadata resource.MetadataResponse
		constructor().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		if metadata.TypeName != terraformType {
			t.Errorf("constructor keyed %q reports type %q", terraformType, metadata.TypeName)
		}
	}

	// The order the provider registers in must be canonical, so the list is
	// reproducible from the registry rather than from edit history.
	if !sort.StringsAreSorted(generatedResourceOrder) {
		t.Error("registration order is not in canonical Terraform-name order")
	}

	// Every registered resource must appear exactly once, including the non-API
	// ones the plan keeps bespoke.
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
