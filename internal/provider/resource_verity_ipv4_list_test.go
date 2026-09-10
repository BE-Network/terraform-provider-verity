package provider

import (
	"testing"

	"terraform-provider-verity/internal/spec"
)

func TestIPv4ListSeedMatchesLegacyCacheKey(t *testing.T) {
	if err := spec.SeedRegistry.Validate(); err != nil {
		t.Fatalf("seed registry invalid: %v", err)
	}
	for _, resourceSpec := range spec.SeedRegistry {
		if resourceSpec.TerraformType == "verity_ipv4_list" {
			if resourceSpec.API.CacheKey != ipv4ListCacheKey {
				t.Fatalf("spec cache key = %q, legacy cache key = %q", resourceSpec.API.CacheKey, ipv4ListCacheKey)
			}
			return
		}
	}
	t.Fatal("IPv4 List spec seed is missing")
}
