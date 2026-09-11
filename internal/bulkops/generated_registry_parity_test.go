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

// assertFixedHeadersMatchLegacySplitKey ties a generated discriminator to the
// legacy bulk configuration. Resources that share a bulk key are separated by a
// fixed header, and the legacy manager splits batches on exactly that key, so
// the two representations must name the same parameter.
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

// TestBulkMetadataMatchesFrozenSnapshot pins the transport facts now filled from
// the spec registry to the values the bulk registry carried as literals. The
// resource type keys every operation's logging and status tracking, and the
// header split key decides how ACL batches are divided, so deriving them had to
// change nothing.
func TestBulkMetadataMatchesFrozenSnapshot(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "bulk_metadata_snapshot.json"))
	if err != nil {
		t.Fatal(err)
	}
	var snapshot map[string]map[string]string
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(resourceRegistry) != len(snapshot) {
		t.Fatalf("bulk registry has %d entries, snapshot has %d", len(resourceRegistry), len(snapshot))
	}
	for key, want := range snapshot {
		config, exists := resourceRegistry[key]
		if !exists {
			t.Errorf("bulk registry is missing %q", key)
			continue
		}
		if config.ResourceType != want["resource_type"] {
			t.Errorf("%s resource type = %q, snapshot = %q", key, config.ResourceType, want["resource_type"])
		}
		if config.HeaderSplitKey != want["header_split_key"] {
			t.Errorf("%s header split key = %q, snapshot = %q", key, config.HeaderSplitKey, want["header_split_key"])
		}
	}
}
