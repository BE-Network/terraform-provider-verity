package lifecycle

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

const importerGoldenEnv = "UPDATE_IMPORTER_GOLDEN"

func TestImporterOutputMatchesGolden(t *testing.T) {
	for _, mode := range []string{"datacenter", "campus"} {
		t.Run(mode, func(t *testing.T) {
			generated := runImporter(t, mode)
			goldenDir := filepath.Join("testdata", "importer_golden", mode)

			if os.Getenv(importerGoldenEnv) == "1" {
				if err := os.RemoveAll(goldenDir); err != nil {
					t.Fatal(err)
				}
				if err := os.MkdirAll(goldenDir, 0o755); err != nil {
					t.Fatal(err)
				}
				for name, content := range generated {
					if err := os.WriteFile(filepath.Join(goldenDir, name), content, 0o644); err != nil {
						t.Fatal(err)
					}
				}
				t.Logf("recorded %d files in %s", len(generated), goldenDir)
				return
			}

			golden := readDir(t, goldenDir)
			if len(golden) == 0 {
				t.Fatalf("no golden files in %s; record them with %s=1", goldenDir, importerGoldenEnv)
			}
			if got, want := fileNames(generated), fileNames(golden); got != want {
				t.Fatalf("generated files differ from the golden set\n  generated: %s\n  golden:    %s", got, want)
			}
			for name, want := range golden {
				if string(generated[name]) != string(want) {
					t.Errorf("%s differs from %s; review the diff and re-record with %s=1",
						name, filepath.Join(goldenDir, name), importerGoldenEnv)
				}
			}
		})
	}
}

func TestImporterOutputIsDeterministic(t *testing.T) {
	for _, mode := range []string{"datacenter", "campus"} {
		t.Run(mode, func(t *testing.T) {
			first := runImporter(t, mode)
			second := runImporter(t, mode)
			if got, want := fileNames(first), fileNames(second); got != want {
				t.Fatalf("two runs wrote different files\n  %s\n  %s", got, want)
			}
			for name, content := range first {
				if string(second[name]) != string(content) {
					t.Errorf("%s differs between two runs of the same import", name)
				}
			}
		})
	}
}

func TestImporterCreatesTheOutputDirectory(t *testing.T) {
	ms := newImporterServer(t, "datacenter")
	outputDir := filepath.Join(t.TempDir(), "does", "not", "exist")
	applyImporter(t, ms, "datacenter", outputDir)

	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("the importer did not create %s: %v", outputDir, err)
	}
	if len(entries) == 0 {
		t.Fatalf("%s is empty", outputDir)
	}
}

func runImporter(t *testing.T, mode string) map[string][]byte {
	t.Helper()
	ms := newImporterServer(t, mode)
	outputDir := t.TempDir()
	applyImporter(t, ms, mode, outputDir)
	return readDir(t, outputDir)
}

func newImporterServer(t *testing.T, mode string) *mock.MockServer {
	t.Helper()
	ms := mock.NewMockServer(mode)
	t.Cleanup(ms.Close)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(mode)); err != nil {
		t.Fatal(err)
	}
	return ms
}

func applyImporter(t *testing.T, ms *mock.MockServer, mode, outputDir string) {
	t.Helper()
	config := mock.ProviderConfig(ms.URL(), mode) + `
data "verity_state_importer" "test" {
  output_dir = "` + filepath.ToSlash(outputDir) + `"
}
`
	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps:                    []fwresource.TestStep{{Config: config}},
	})
}

func readDir(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	files := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		files[entry.Name()] = content
	}
	return files
}

func fileNames(files map[string][]byte) string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
