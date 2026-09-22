package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const docsTestRegistry = "../../specs/generated_registry.json"

func docsTestOptions(dir string, check bool) docsOptions {
	return docsOptions{
		Registry:  docsTestRegistry,
		OutputDir: filepath.Join(dir, "resources"),
		Keep:      "verity_operation_stage.md",
		Check:     check,
	}
}

func TestDocsCheckReportsStalePagesAndKeepsHandwrittenOnes(t *testing.T) {
	dir := t.TempDir()
	opts := docsTestOptions(dir, false)
	if err := generateDocs(opts); err != nil {
		t.Fatalf("generate: %v", err)
	}

	stale := filepath.Join(opts.OutputDir, "verity_retired.md")
	handwritten := filepath.Join(opts.OutputDir, "verity_operation_stage.md")
	if err := os.WriteFile(stale, []byte("# retired\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(handwritten, []byte("# operation stage\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := generateDocs(docsTestOptions(dir, true))
	if err == nil {
		t.Fatal("check passed with a page the registry does not generate")
	}
	if !strings.Contains(err.Error(), "verity_retired.md is not generated from the registry") {
		t.Fatalf("check error does not name the stale page: %v", err)
	}
	if strings.Contains(err.Error(), "verity_operation_stage.md") {
		t.Fatalf("check reported the handwritten page: %v", err)
	}

	if err := generateDocs(opts); err != nil {
		t.Fatalf("regenerate: %v", err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("regeneration left the stale page: %v", err)
	}
	kept, err := os.ReadFile(handwritten)
	if err != nil || string(kept) != "# operation stage\n" {
		t.Fatalf("regeneration changed the handwritten page: %q, %v", kept, err)
	}
	if err := generateDocs(docsTestOptions(dir, true)); err != nil {
		t.Fatalf("check failed after regeneration: %v", err)
	}
}

func TestDocsCheckReportsAnEditedPage(t *testing.T) {
	dir := t.TempDir()
	opts := docsTestOptions(dir, false)
	if err := generateDocs(opts); err != nil {
		t.Fatalf("generate: %v", err)
	}
	page := filepath.Join(opts.OutputDir, "verity_service.md")
	raw, err := os.ReadFile(page)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(page, append(raw, []byte("edited by hand\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	err = generateDocs(docsTestOptions(dir, true))
	if err == nil || !strings.Contains(err.Error(), "verity_service.md differs") {
		t.Fatalf("check did not report the edited page: %v", err)
	}
}

func TestDocsGenerateOnePagePerRegistryResource(t *testing.T) {
	dir := t.TempDir()
	opts := docsTestOptions(dir, false)
	if err := generateDocs(opts); err != nil {
		t.Fatalf("generate: %v", err)
	}
	registry, err := readRegistry(docsTestRegistry)
	if err != nil {
		t.Fatal(err)
	}
	pages, err := filepath.Glob(filepath.Join(opts.OutputDir, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != len(registry) {
		t.Fatalf("generated %d pages for %d registry resources", len(pages), len(registry))
	}
	for _, resource := range registry {
		raw, err := os.ReadFile(filepath.Join(opts.OutputDir, resource.TerraformType+".md"))
		if err != nil {
			t.Fatalf("%s: %v", resource.TerraformType, err)
		}
		if !strings.Contains(string(raw), "terraform import "+resource.TerraformType+".") {
			t.Errorf("%s page has no import command", resource.TerraformType)
		}
	}
}
