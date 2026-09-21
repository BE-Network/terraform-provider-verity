package bulkops

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"terraform-provider-verity/internal/spec"
)

func TestGeneratedSpecsMatchLegacyBulkRegistry(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "specs", "generated_registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	var artifact struct {
		Resources spec.Registry `json:"resources"`
	}
	if err := json.Unmarshal(raw, &artifact); err != nil {
		t.Fatal(err)
	}
	if len(artifact.Resources) == 0 {
		t.Fatal("generated registry is empty")
	}
	for _, resourceSpec := range artifact.Resources {
		config, exists := resourceRegistry[resourceSpec.API.BulkKey]
		if !exists || config.ResourceType != resourceSpec.API.BulkKey {
			t.Fatalf("legacy bulk registry does not match generated bulk key %q for %s", resourceSpec.API.BulkKey, resourceSpec.TerraformType)
		}
		assertFixedHeadersMatchLegacySplitKey(t, resourceSpec, config)
	}
}

func assertFixedHeadersMatchLegacySplitKey(t *testing.T, resourceSpec spec.ResourceSpec, config ResourceConfig) {
	t.Helper()
	if config.HeaderSplitKey == "" {
		if len(resourceSpec.API.FixedHeaders) != 0 {
			t.Fatalf("%s declares fixed headers %v but legacy bulk key %q does not split on a header", resourceSpec.TerraformType, resourceSpec.API.FixedHeaders, resourceSpec.API.BulkKey)
		}
		return
	}
	value, exists := resourceSpec.API.FixedHeaders[config.HeaderSplitKey]
	if !exists {
		t.Fatalf("%s must pin legacy split key %q, got fixed headers %v", resourceSpec.TerraformType, config.HeaderSplitKey, resourceSpec.API.FixedHeaders)
	}
	if len(resourceSpec.API.FixedHeaders) != 1 {
		t.Fatalf("%s pins %d fixed headers but legacy bulk key %q splits only on %q", resourceSpec.TerraformType, len(resourceSpec.API.FixedHeaders), resourceSpec.API.BulkKey, config.HeaderSplitKey)
	}
	if value == "" {
		t.Fatalf("%s pins an empty value for legacy split key %q", resourceSpec.TerraformType, config.HeaderSplitKey)
	}
}

func TestBulkMetadataMatchesGolden(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "bulk_metadata_golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden map[string]map[string]string
	if err := json.Unmarshal(raw, &golden); err != nil {
		t.Fatal(err)
	}
	if len(resourceRegistry) != len(golden) {
		t.Fatalf("bulk registry has %d entries, golden has %d", len(resourceRegistry), len(golden))
	}
	for key, want := range golden {
		config, exists := resourceRegistry[key]
		if !exists {
			t.Errorf("bulk registry is missing %q", key)
			continue
		}
		if config.ResourceType != want["resource_type"] {
			t.Errorf("%s resource type = %q, golden = %q", key, config.ResourceType, want["resource_type"])
		}
		if config.HeaderSplitKey != want["header_split_key"] {
			t.Errorf("%s header split key = %q, golden = %q", key, config.HeaderSplitKey, want["header_split_key"])
		}
	}
}
