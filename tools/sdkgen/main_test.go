package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareRetainsOnlySDKSourceAndIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	client := "package openapi\nfunc (c *APIClient) callAPI(request *http.Request) (*http.Response, error) { return nil, nil }\n"
	for name, content := range map[string]string{
		"client.go": client, "model.go": "package openapi\ntype Model struct{}\n",
		"go.mod": "module generated", "go.sum": "generated", "README.md": "generated",
		"model_test.go": "package openapi", "git_push.sh": "generated",
		"tests/api_test.go": "generated", "api/openapi.yaml": "generated",
		"docs/model.md": "generated", ".openapi-generator/FILES": "generated",
	} {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := prepare(dir); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || entries[0].Name() != "client.go" || entries[1].Name() != "model.go" {
		t.Fatalf("retained entries: %v", entries)
	}
	first, err := os.ReadFile(filepath.Join(dir, "client.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), "[REDACTED AUTH REQUEST BODY]") || !strings.Contains(string(first), "[REDACTED AUTH RESPONSE BODY]") {
		t.Fatal("credential redaction missing from prepared client")
	}
	if err := prepare(dir); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(filepath.Join(dir, "client.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatal("preparation is not idempotent")
	}
}

func TestCustomizationRefusesChangedGeneratorContract(t *testing.T) {
	for _, source := range []string{
		"package openapi",
		"package openapi\nfunc (c *APIClient) callAPI(request *http.Request, changed bool) (*http.Response, error) { return nil, nil }",
		"package openapi\nfunc (c *OtherClient) callAPI(request *http.Request) (*http.Response, error) { return nil, nil }",
	} {
		if _, err := customizeClient([]byte(source)); err == nil {
			t.Fatal("changed generator contract was accepted")
		}
	}
}
