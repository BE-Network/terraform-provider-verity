package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"terraform-provider-verity/internal/spec"
)

func TestBulkGenerationIsDeterministicAndChecksDrift(t *testing.T) {
	root := filepath.Join("..", "..")
	options := bulkOptions{Registry: filepath.Join(root, "specs", "generated_registry.json"), OpenAPIDir: filepath.Join(root, "openapi"), Output: filepath.Join(t.TempDir(), "bulk.go")}
	if err := generateBulk(options); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(options.Output)
	if err != nil {
		t.Fatal(err)
	}
	if err := generateBulk(options); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(options.Output)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("generation is not deterministic")
	}
	options.Check = true
	if err := generateBulk(options); err != nil {
		t.Fatal(err)
	}
	registry, err := readRegistry(options.Registry)
	if err != nil {
		t.Fatal(err)
	}
	for i, j := 0, len(registry)-1; i < j; i, j = i+1, j-1 {
		registry[i], registry[j] = registry[j], registry[i]
	}
	sdk, err := discoverBulkSDK(options.OpenAPIDir)
	if err != nil {
		t.Fatal(err)
	}
	bindings, err := planBulkBindings(registry, sdk)
	if err != nil {
		t.Fatal(err)
	}
	reversed, err := renderBulkBindings(registry, bindings)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, reversed) {
		t.Fatal("registry input order changes output")
	}
	if err := os.WriteFile(options.Output, append(first, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	if err := generateBulk(options); err == nil || !strings.Contains(err.Error(), "tools/generate_provider.sh --write") {
		t.Fatalf("drift error = %v", err)
	}
}

func TestBulkBindingsRejectMissingOrConflictingContracts(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, test := range []struct {
		name string
		edit func(spec.Registry, bulkSDK)
		want string
	}{
		{"missing operation", func(_ spec.Registry, sdk bulkSDK) { delete(sdk.methods, "GatewaysAPIService.GatewaysPutExecute") }, "needs one SDK binding"},
		{"missing constructor", func(_ spec.Registry, sdk bulkSDK) { delete(sdk.methods, "GatewaysAPIService.GatewaysPut") }, "no compatible SDK request constructor"},
		{"changed Execute signature", func(_ spec.Registry, sdk bulkSDK) {
			sdk.methods["ApiGatewaysGetRequest.Execute"].Type.Params.List = []*ast.Field{{Type: ast.NewIdent("string")}}
		}, "no compatible SDK Execute method"},
		{"missing delete setter", func(_ spec.Registry, sdk bulkSDK) { delete(sdk.methods, "ApiGatewaysDeleteRequest.GatewayName") }, "needs one setter"},
		{"missing split setter", func(_ spec.Registry, sdk bulkSDK) { delete(sdk.methods, "ApiAclsPutRequest.IpVersion") }, "needs one setter"},
		{"wrong wrapper", func(r spec.Registry, _ bulkSDK) {
			for i := range r {
				if r[i].API.BulkKey == "gateway" {
					r[i].API.RequestWrapperKey = "missing"
				}
			}
		}, "needs one SDK body"},
		{"wrong delete parameter", func(r spec.Registry, _ bulkSDK) {
			for i := range r {
				if r[i].API.BulkKey == "gateway" {
					r[i].API.DeleteParameter = "missing"
				}
			}
		}, "has no query parameter"},
		{"conflicting variants", func(r spec.Registry, _ bulkSDK) {
			for i := range r {
				if r[i].TerraformType == "verity_acl_v6" {
					r[i].API.EndpointPath = "/different"
				}
			}
		}, "disagree on API route"},
		{"conflicting splits", func(r spec.Registry, _ bulkSDK) {
			for i := range r {
				if r[i].TerraformType == "verity_acl_v6" {
					r[i].API.FixedHeaders = map[string]string{"other": "6"}
				}
			}
		}, "inconsistent split parameters"},
		{"unknown split", func(r spec.Registry, _ bulkSDK) {
			for i := range r {
				if r[i].API.BulkKey == "gateway" {
					r[i].API.FixedHeaders = map[string]string{"missing": "4"}
				}
			}
		}, "has no split query parameter"},
		{"missing order", func(r spec.Registry, _ bulkSDK) { r[0].BulkOrder = nil }, "has no"},
	} {
		t.Run(test.name, func(t *testing.T) {
			registry, err := readRegistry(filepath.Join(root, "specs", "generated_registry.json"))
			if err != nil {
				t.Fatal(err)
			}
			sdk, err := discoverBulkSDK(filepath.Join(root, "openapi"))
			if err != nil {
				t.Fatal(err)
			}
			test.edit(registry, sdk)
			if _, err := planBulkBindings(registry, sdk); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestBulkBindingsDiscoverAConventionalNewEndpoint(t *testing.T) {
	dir := t.TempDir()
	const source = `package openapi
import ("context"; "net/http")
type APIClient struct { WidgetsAPI *WidgetServiceAPIService }
type WidgetServiceAPIService struct {}
type Widget struct {}
type WidgetBody struct { Items *map[string]Widget ` + "`json:\"items,omitempty\"`" + ` }
type ApiWidgetDeleteRequest struct {}
func (r ApiWidgetDeleteRequest) WidgetName(names []string) ApiWidgetDeleteRequest { return r }
func (r ApiWidgetDeleteRequest) Execute() (*http.Response, error) { return nil, nil }
func parameterAddToHeaderOrQuery(values interface{}, name string, value interface{}) {}
func (a *WidgetServiceAPIService) DeleteWidgets(ctx context.Context) ApiWidgetDeleteRequest { return ApiWidgetDeleteRequest{} }
func (a *WidgetServiceAPIService) DeleteWidgetsExecute(r ApiWidgetDeleteRequest) (*http.Response, error) {
 var localVarHTTPMethod = http.MethodDelete
 localVarPath := "" + "/widgets"
 var localVarQueryParams interface{}
 parameterAddToHeaderOrQuery(localVarQueryParams, "widget_name", nil)
 _ = localVarHTTPMethod
 _ = localVarPath
 return nil, nil
}
type ApiWidgetRequest struct {}
func (r ApiWidgetRequest) Payload(body WidgetBody) ApiWidgetRequest { return r }
func (r ApiWidgetRequest) Execute() (*http.Response, error) { return nil, nil }
func (a *WidgetServiceAPIService) MakeWidgets(ctx context.Context) ApiWidgetRequest { return ApiWidgetRequest{} }
func (a *WidgetServiceAPIService) MakeWidgetsExecute(r ApiWidgetRequest) (*http.Response, error) {
 var localVarHTTPMethod = http.MethodPut
 localVarPath := "" + "/widgets"
 _ = localVarHTTPMethod
 _ = localVarPath
 return nil, nil
}
func (a *WidgetServiceAPIService) ReadWidgets(ctx context.Context) ApiWidgetRequest { return ApiWidgetRequest{} }
func (a *WidgetServiceAPIService) ReadWidgetsExecute(r ApiWidgetRequest) (*http.Response, error) {
 var localVarHTTPMethod = http.MethodGet
 localVarPath := "" + "/widgets"
 _ = localVarHTTPMethod
 _ = localVarPath
 return nil, nil
}
`
	if err := os.WriteFile(filepath.Join(dir, "widgets.go"), []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	sdk, err := discoverBulkSDK(dir)
	if err != nil {
		t.Fatal(err)
	}
	registry := spec.Registry{{TerraformType: "verity_widget", API: spec.APIResourceSpec{EndpointPath: "/widgets", BulkKey: "widget", RequestWrapperKey: "items", DeleteParameter: "widget_name", CacheKey: "widgets", ResponseCollectionKey: "widget"}, Modes: []spec.Mode{spec.ModeCampus}, Operations: spec.OperationSpec{Create: true, Read: true, Delete: true}, BulkOrder: map[spec.Mode]spec.BulkOrderSpec{spec.ModeCampus: {Put: 10, Delete: 10}}}}
	bindings, err := planBulkBindings(registry, sdk)
	if err != nil {
		t.Fatal(err)
	}
	put := bindings[0].Operations["PUT"]
	if put.Service != "WidgetsAPI" || put.Method != "MakeWidgets" || put.Setter != "Payload" || put.BodyType != "WidgetBody" {
		t.Fatalf("new endpoint binding = %+v", put)
	}
	if got := registry.BulkOperationOrder(spec.ModeCampus, "PUT"); !reflect.DeepEqual(got, []string{"widget"}) {
		t.Fatalf("new endpoint order = %v", got)
	}
	output, err := renderBulkBindings(registry, bindings)
	if err != nil {
		t.Fatal(err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "generated.go", output, 0)
	if err != nil {
		t.Fatal(err)
	}
	var calls []string
	ast.Inspect(file, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			calls = append(calls, typeString(call.Fun))
		}
		return true
	})
	found := false
	for _, call := range calls {
		found = found || call == "c.WidgetsAPI.MakeWidgets(ctx).Payload"
	}
	if !found {
		t.Fatalf("generated callback did not use the discovered SDK setter: %v", calls)
	}
}
