package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectedVersionRejectsInvalidOrAmbiguousSelection(t *testing.T) {
	for _, tt := range []struct {
		name string
		yaml string
		want string
	}{
		{"multi_digit_minor", "format_version: 1\napi_version: '6.10'\n", "6.10"},
		{"missing", "format_version: 1\n", ""},
		{"patch", "format_version: 1\napi_version: '6.6.1'\n", ""},
		{"path", "format_version: 1\napi_version: '../6.6'\n", ""},
		{"noncanonical", "format_version: 1\napi_version: '06.6'\n", ""},
		{"unknown_format", "format_version: 2\napi_version: '6.6'\n", ""},
		{"multiple_documents", "format_version: 1\napi_version: '6.6'\n---\napi_version: '6.7'\n", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "overrides.yaml")
			if err := os.WriteFile(path, []byte(tt.yaml), 0644); err != nil {
				t.Fatal(err)
			}
			got, err := selectedVersion(path)
			if tt.want == "" {
				if err == nil {
					t.Fatal("invalid version selection accepted")
				}
			} else if err != nil || got != tt.want {
				t.Fatalf("selected version = %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}

func TestVerifyVersionRefusesInputsFromAnotherRelease(t *testing.T) {
	dir := t.TempDir()
	raw := filepath.Join(dir, "raw.json")
	if err := os.WriteFile(raw, []byte(`{"openapi":"3.0.0","paths":{}}`), 0644); err != nil {
		t.Fatal(err)
	}
	inputs := filepath.Join(dir, "inputs")
	if err := normalize(normalizeOptions{
		Version: "6.10", Datacenter: raw, Campus: raw, OutputDir: inputs,
		SourceExportDate: "2026-10-06", Provenance: "test fixture",
	}); err != nil {
		t.Fatal(err)
	}
	if err := verifyVersion(inputs, "6.10"); err != nil {
		t.Fatal(err)
	}
	if err := verifyVersion(inputs, "6.11"); err == nil || !strings.Contains(err.Error(), "does not match selected version") {
		t.Fatalf("mismatched input version error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(inputs, "campus.json"), []byte("{}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := verifyVersion(inputs, "6.10"); err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("corrupted inputs error = %v", err)
	}
}
