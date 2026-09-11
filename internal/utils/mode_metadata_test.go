package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestModeMetadataMatchesFrozenSnapshot pins the mode tables to the behavior the
// provider shipped before they were derived from the spec registry. Mode data
// decides which fields are sent for a datacenter or campus system, so moving its
// source must not change a single entry. The snapshot is the frozen "before";
// regenerate it deliberately only when a mode change is itself the intent.
func TestModeMetadataMatchesFrozenSnapshot(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "mode_metadata_snapshot.json"))
	if err != nil {
		t.Fatal(err)
	}
	var snapshot struct {
		ResourceCompatibility map[string]ResourceMode         `json:"resource_compatibility"`
		ModeFields            map[string]map[string]FieldMode `json:"mode_fields"`
	}
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		t.Fatal(err)
	}

	if len(ResourceCompatibility) != len(snapshot.ResourceCompatibility) {
		t.Fatalf("resource compatibility has %d entries, snapshot has %d", len(ResourceCompatibility), len(snapshot.ResourceCompatibility))
	}
	for name, want := range snapshot.ResourceCompatibility {
		if got, exists := ResourceCompatibility[name]; !exists || got != want {
			t.Errorf("resource %q compatibility = %q, snapshot = %q", name, got, want)
		}
	}

	if len(ModeFields) != len(snapshot.ModeFields) {
		t.Fatalf("mode fields cover %d resources, snapshot covers %d", len(ModeFields), len(snapshot.ModeFields))
	}
	for resource, wantFields := range snapshot.ModeFields {
		gotFields, exists := ModeFields[resource]
		if !exists {
			t.Errorf("mode fields are missing resource %q", resource)
			continue
		}
		if len(gotFields) != len(wantFields) {
			t.Errorf("resource %q has %d mode fields, snapshot has %d", resource, len(gotFields), len(wantFields))
		}
		for field, want := range wantFields {
			if got, exists := gotFields[field]; !exists || got != want {
				t.Errorf("%s.%s mode = %q, snapshot = %q", resource, field, got, want)
			}
		}
	}
}

// The generated and pending tables must stay disjoint, otherwise a resource could
// be represented in the registry while a stale hand-maintained entry silently
// overrode it.
func TestPendingModeMetadataDoesNotShadowGenerated(t *testing.T) {
	for name := range pendingResourceCompatibility {
		if _, exists := generatedResourceCompatibility[name]; exists {
			t.Errorf("resource %q is both generated and pending; remove the pending entry", name)
		}
	}
	for resource := range pendingModeFields {
		if _, exists := generatedModeFields[resource]; exists {
			t.Errorf("endpoint %q is both generated and pending; remove the pending entry", resource)
		}
	}
}
