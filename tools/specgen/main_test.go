package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCanonicalJSON(t *testing.T) {
	got, err := canonicalJSON([]byte(`{"z":1.20,"a":"<value>","nested":{"b":false,"a":null}}`))
	if err != nil {
		t.Fatalf("canonicalJSON() error = %v", err)
	}
	want := "{\n  \"a\": \"<value>\",\n  \"nested\": {\n    \"a\": null,\n    \"b\": false\n  },\n  \"z\": 1.20\n}\n"
	if string(got) != want {
		t.Fatalf("canonicalJSON() = %q, want %q", got, want)
	}
}

func TestNormalizeAndVerify(t *testing.T) {
	tempDir := t.TempDir()
	datacenter := filepath.Join(tempDir, "raw-datacenter.json")
	campus := filepath.Join(tempDir, "raw-campus.json")
	if err := os.WriteFile(datacenter, []byte(`{"z":true,"a":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(campus, []byte(`{"b":null,"a":"campus"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := filepath.Join(tempDir, "canonical")
	opts := normalizeOptions{Version: "6.6", Datacenter: datacenter, Campus: campus, OutputDir: outputDir, SourceExportDate: "2026-08-11", Provenance: "test fixture"}
	if err := normalize(opts); err != nil {
		t.Fatalf("normalize() error = %v", err)
	}
	if err := verify(outputDir); err != nil {
		t.Fatalf("verify() error = %v", err)
	}

	first, err := os.ReadFile(filepath.Join(outputDir, "datacenter.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := normalize(opts); err != nil {
		t.Fatalf("second normalize() error = %v", err)
	}
	second, err := os.ReadFile(filepath.Join(outputDir, "datacenter.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatal("normalization output changed for identical input")
	}

	if err := os.WriteFile(filepath.Join(outputDir, "campus.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verify(outputDir); err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("verify() after mutation error = %v, want checksum mismatch", err)
	}
}

func TestExtractCoverageMergesModeSpecificOpenAPIWithoutGuessingPolicies(t *testing.T) {
	tempDir := t.TempDir()
	datacenter := filepath.Join(tempDir, "datacenter-raw.json")
	campus := filepath.Join(tempDir, "campus-raw.json")
	datacenterDocument := `{"openapi":"3.0.0","paths":{"/widgets":{"put":{"requestBody":{"content":{"application/json":{"schema":{"type":"object","properties":{"widget":{"type":"object","properties":{"widget_name":{"type":"object","properties":{"name":{"type":"string"},"enabled":{"type":"boolean"}}}}}}}}}}}}}}`
	campusDocument := `{"openapi":"3.0.0","paths":{"/widgets":{"patch":{"requestBody":{"content":{"application/json":{"schema":{"type":"object","properties":{"widget":{"type":"object","properties":{"widget_name":{"type":"object","properties":{"name":{"type":"string"},"campus_only":{"type":"integer","nullable":true}}}}}}}}}}}}}}`
	if err := os.WriteFile(datacenter, []byte(datacenterDocument), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(campus, []byte(campusDocument), 0o644); err != nil {
		t.Fatal(err)
	}
	inputs := filepath.Join(tempDir, "inputs")
	if err := normalize(normalizeOptions{Version: "6.6", Datacenter: datacenter, Campus: campus, OutputDir: inputs, SourceExportDate: "2026-09-03", Provenance: "test"}); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(tempDir, "coverage.json")
	if err := extract(extractOptions{InputDir: inputs, Output: output}); err != nil {
		t.Fatal(err)
	}
	if err := extract(extractOptions{InputDir: inputs, Output: output, Check: true}); err != nil {
		t.Fatal(err)
	}
	var report coverageManifest
	raw, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Resources) != 1 {
		t.Fatalf("resources = %#v", report.Resources)
	}
	resource := report.Resources[0]
	if resource.RequestWrapperKey != "widget" || strings.Join(resource.Modes, ",") != "campus,datacenter" {
		t.Fatalf("resource = %#v", resource)
	}
	if len(resource.Fields) != 3 || resource.Fields[0].APIName != "campus_only" || !resource.Fields[0].Nullable {
		t.Fatalf("fields = %#v", resource.Fields)
	}
	if !strings.Contains(strings.Join(resource.Fields[0].Review, " "), "lifecycle") {
		t.Fatalf("field policy review was not reported: %#v", resource.Fields[0])
	}
}

func TestExtractionCapturesGetResponseAndDeleteParameterShapes(t *testing.T) {
	resource := &coverageResource{}
	pathItem := map[string]any{
		"get": map[string]any{
			"responses": map[string]any{
				"200": map[string]any{
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"properties": map[string]any{"widgets": map[string]any{"type": "object"}},
							},
						},
					},
				},
			},
		},
		"delete": map[string]any{
			"parameters": []any{map[string]any{
				"in": "query", "name": "widget_name", "required": true,
				"schema": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			}},
		},
	}
	operations := resourceOperations(pathItem)
	if strings.Join(operations, ",") != "get,delete" {
		t.Fatalf("operations = %#v", operations)
	}
	mergeResponseShape(resource, "datacenter", pathItem["get"])
	mergeDeleteParameters(resource, "datacenter", pathItem["delete"])
	if resource.ResponseCollectionKey != "widgets" {
		t.Fatalf("response key = %q", resource.ResponseCollectionKey)
	}
	if len(resource.DeleteParameters) != 1 {
		t.Fatalf("delete parameters = %#v", resource.DeleteParameters)
	}
	parameter := resource.DeleteParameters[0]
	if parameter.Name != "widget_name" || !parameter.Required || parameter.Kind != "array" || parameter.ItemKind != "string" {
		t.Fatalf("delete parameter = %#v", parameter)
	}
}
