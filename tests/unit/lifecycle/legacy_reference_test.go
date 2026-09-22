package lifecycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func legacyReference(t *testing.T) map[string]map[string]interface{} {
	t.Helper()
	path := legacyReferencePath(t)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no recorded legacy reference for %s: %v", t.Name(), err)
	}
	var captured map[string]map[string]interface{}
	if err := json.Unmarshal(raw, &captured); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return captured
}

func legacyReferencePath(t *testing.T) string {
	name := strings.NewReplacer("/", "__", " ", "_", ":", "_", "'", "", "\"", "").Replace(t.Name())
	return filepath.Join("testdata", "legacy_reference", name+"-1.json")
}
