package registry

import (
	"encoding/json"
	"testing"
)

func TestVersionComesFromEmbeddedRegistry(t *testing.T) {
	var artifact struct {
		APIVersion string `json:"api_version"`
	}
	if err := json.Unmarshal(embedded, &artifact); err != nil {
		t.Fatal(err)
	}
	version, err := Version()
	if err != nil {
		t.Fatal(err)
	}
	if version.String() != artifact.APIVersion {
		t.Fatalf("runtime version = %s, embedded API version = %s", version, artifact.APIVersion)
	}
}
