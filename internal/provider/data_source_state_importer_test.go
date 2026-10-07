package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/zclconf/go-cty/cty"

	"terraform-provider-verity/internal/importer"
	"terraform-provider-verity/openapi"
)

func TestImportBlocksPreserveLiteralIDs(t *testing.T) {
	for _, id := range []string{"plain", "quoted\"name", "back\\slash", "${reference}", "%{directive}", "braces{name}", "line\nwith\ttabs"} {
		t.Run(id, func(t *testing.T) {
			path, err := createImportBlocks(t.TempDir(), []importer.ImportedResource{{
				TerraformType: "verity_service", TerraformName: "edge", ID: id,
			}})
			if err != nil {
				t.Fatal(err)
			}
			got := readImportBlocks(t, path)
			if want := map[string]string{"verity_service.edge": id}; !reflect.DeepEqual(got, want) {
				t.Fatalf("imports = %#v, want %#v", got, want)
			}
		})
	}
}

func TestStateImporterReportsOnlyCurrentRun(t *testing.T) {
	for _, mode := range []string{"datacenter", "campus"} {
		t.Run(mode, func(t *testing.T) {
			services := map[string]interface{}{
				"edge\"node": map[string]interface{}{"name": "different-body-name", "enable": true},
				"node2":      map[string]interface{}{"enable": true},
				"node10":     map[string]interface{}{"enable": true},
			}
			var empty atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				body := map[string]interface{}{}
				if r.URL.Path == "/services" {
					body["service"] = services
					if empty.Load() {
						body["service"] = map[string]interface{}{}
					}
				}
				if err := json.NewEncoder(w).Encode(body); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			config := openapi.NewConfiguration()
			config.Servers = openapi.ServerConfigurations{{URL: server.URL}}
			config.HTTPClient = server.Client()
			source := &stateImporterDataSource{client: &providerContext{client: openapi.NewAPIClient(config), mode: mode}}
			ctx := context.Background()
			var schema datasource.SchemaResponse
			source.Schema(ctx, datasource.SchemaRequest{}, &schema)
			outputDir := t.TempDir()
			unrelatedPath := filepath.Join(outputDir, "user.tf")
			unrelated := []byte("this is deliberately not valid HCL")
			if err := os.WriteFile(unrelatedPath, unrelated, 0644); err != nil {
				t.Fatal(err)
			}
			manualPath := filepath.Join(outputDir, "manual.tf")
			manual := []byte("resource \"verity_tenant\" \"manual\" { name = \"manual\" }")
			if err := os.WriteFile(manualPath, manual, 0644); err != nil {
				t.Fatal(err)
			}
			raw := tftypes.NewValue(schema.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"id":             tftypes.NewValue(tftypes.String, nil),
				"output_dir":     tftypes.NewValue(tftypes.String, outputDir),
				"imported_files": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
			})
			request := datasource.ReadRequest{Config: tfsdk.Config{Raw: raw, Schema: schema.Schema}}
			wantImports := map[string]string{"verity_service.edge_node": "edge\"node", "verity_service.node2": "node2", "verity_service.node10": "node10"}
			stageImports := map[string]string{}
			for run := 0; run < 4; run++ {
				if run == 2 {
					empty.Store(true)
				}
				if run == 3 {
					if err := os.Remove(filepath.Join(outputDir, "services.tf")); err != nil {
						t.Fatal(err)
					}
				}
				response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
				source.Read(ctx, request, &response)
				if run == 2 {
					checkStaleFileDiagnostic(t, response, filepath.Join(outputDir, "services.tf"))
					if _, err := os.Stat(filepath.Join(outputDir, "services.tf")); err != nil {
						t.Fatalf("stale resource file was removed: %v", err)
					}
					if got := readImportBlocks(t, filepath.Join(outputDir, "import_blocks.tf")); !reflect.DeepEqual(got, wantImports) {
						t.Fatalf("failed generation replaced import blocks: %#v", got)
					}
					continue
				}
				if response.Diagnostics.HasError() {
					t.Fatalf("run %d: %v", run, response.Diagnostics)
				}
				if run == 0 {
					stageImports = readStageImports(t, filepath.Join(outputDir, "stages.tf"))
					for address, id := range stageImports {
						wantImports[address] = id
					}
				}
				var result stateImporterDataSourceModel
				if diags := response.State.Get(ctx, &result); diags.HasError() {
					t.Fatal(diags)
				}
				wantFiles := []string{filepath.Join(outputDir, "services.tf"), filepath.Join(outputDir, "stages.tf"), filepath.Join(outputDir, "import_blocks.tf")}
				if run == 3 {
					wantFiles = wantFiles[1:]
					wantImports = stageImports
				}
				var gotFiles []string
				for _, file := range result.ImportedFiles {
					gotFiles = append(gotFiles, file.ValueString())
				}
				if !reflect.DeepEqual(gotFiles, wantFiles) {
					t.Fatalf("run %d files = %v, want %v", run, gotFiles, wantFiles)
				}
				if got := readImportBlocks(t, filepath.Join(outputDir, "import_blocks.tf")); !reflect.DeepEqual(got, wantImports) {
					t.Fatalf("run %d imports = %#v, want %#v", run, got, wantImports)
				}
			}
			if got, err := os.ReadFile(unrelatedPath); err != nil || string(got) != string(unrelated) {
				t.Fatalf("unrelated file changed: %q, %v", got, err)
			}
			if got, err := os.ReadFile(manualPath); err != nil || string(got) != string(manual) {
				t.Fatalf("manual resource file changed: %q, %v", got, err)
			}
			stalePath := filepath.Join(outputDir, "racks.tf")
			if err := os.WriteFile(stalePath, []byte("resource \"verity_rack\" \"stale\" { name = \"stale\" }"), 0644); err != nil {
				t.Fatal(err)
			}
			response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
			source.Read(ctx, request, &response)
			checkStaleFileDiagnostic(t, response, stalePath)
		})
	}
}

func checkStaleFileDiagnostic(t *testing.T, response datasource.ReadResponse, path string) {
	t.Helper()
	if !response.Diagnostics.HasError() {
		t.Fatal("stale resource file did not stop generation")
	}
	for _, diagnostic := range response.Diagnostics.Errors() {
		if strings.Contains(diagnostic.Detail(), path) && strings.Contains(diagnostic.Detail(), "was not regenerated") && strings.Contains(diagnostic.Detail(), "move or remove") {
			return
		}
	}
	t.Fatalf("missing actionable stale-file diagnostic for %q: %v", path, response.Diagnostics)
}

func readImportBlocks(t *testing.T, path string) map[string]string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	file, diags := hclsyntax.ParseConfig(content, path, hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	imports := map[string]string{}
	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		if block.Type != "import" {
			t.Fatalf("unexpected block %q", block.Type)
		}
		traversal, diags := hcl.AbsTraversalForExpr(block.Body.Attributes["to"].Expr)
		if diags.HasErrors() || len(traversal) != 2 {
			t.Fatalf("invalid resource address: %v, %v", traversal, diags)
		}
		address := traversal.RootName() + "." + traversal[1].(hcl.TraverseAttr).Name
		value, diags := block.Body.Attributes["id"].Expr.Value(nil)
		if diags.HasErrors() || value.Type() != cty.String {
			t.Fatalf("invalid literal import ID: %v, %v", value, diags)
		}
		if _, duplicate := imports[address]; duplicate {
			t.Fatalf("duplicate import %q", address)
		}
		imports[address] = value.AsString()
	}
	return imports
}

func readStageImports(t *testing.T, path string) map[string]string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	file, diags := hclsyntax.ParseConfig(content, path, hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	imports := map[string]string{}
	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		if block.Type != "resource" || len(block.Labels) != 2 || block.Labels[0] != "verity_operation_stage" {
			t.Fatalf("unexpected stage block: %#v", block)
		}
		imports[block.Labels[0]+"."+block.Labels[1]] = "stage"
	}
	if len(imports) == 0 {
		t.Fatal("no generated stages")
	}
	return imports
}
