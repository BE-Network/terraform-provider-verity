package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"terraform-provider-verity/internal/spec"
)

func TestBulkBindingsSupportABodyWithNoFields(t *testing.T) {
	dir := t.TempDir()
	const source = `package openapi
import ("context"; "net/http")
type APIClient struct { GroupsAPI *GroupServiceAPIService }
type GroupServiceAPIService struct {}
type GroupBody struct { Items map[string]map[string]interface{} ` + "`json:\"items,omitempty\"`" + ` }
type ApiGroupDeleteRequest struct {}
func (r ApiGroupDeleteRequest) GroupName(names []string) ApiGroupDeleteRequest { return r }
func (r ApiGroupDeleteRequest) Execute() (*http.Response, error) { return nil, nil }
func parameterAddToHeaderOrQuery(values interface{}, name string, value interface{}) {}
func (a *GroupServiceAPIService) DeleteGroups(ctx context.Context) ApiGroupDeleteRequest { return ApiGroupDeleteRequest{} }
func (a *GroupServiceAPIService) DeleteGroupsExecute(r ApiGroupDeleteRequest) (*http.Response, error) {
 var localVarHTTPMethod = http.MethodDelete
 localVarPath := "" + "/groups"
 var localVarQueryParams interface{}
 parameterAddToHeaderOrQuery(localVarQueryParams, "group_name", nil)
 _ = localVarHTTPMethod
 _ = localVarPath
 return nil, nil
}
type ApiGroupRequest struct {}
func (r ApiGroupRequest) Payload(body GroupBody) ApiGroupRequest { return r }
func (r ApiGroupRequest) Execute() (*http.Response, error) { return nil, nil }
func (a *GroupServiceAPIService) MakeGroups(ctx context.Context) ApiGroupRequest { return ApiGroupRequest{} }
func (a *GroupServiceAPIService) MakeGroupsExecute(r ApiGroupRequest) (*http.Response, error) {
 var localVarHTTPMethod = http.MethodPut
 localVarPath := "" + "/groups"
 _ = localVarHTTPMethod
 _ = localVarPath
 return nil, nil
}
func (a *GroupServiceAPIService) ReadGroups(ctx context.Context) ApiGroupRequest { return ApiGroupRequest{} }
func (a *GroupServiceAPIService) ReadGroupsExecute(r ApiGroupRequest) (*http.Response, error) {
 var localVarHTTPMethod = http.MethodGet
 localVarPath := "" + "/groups"
 _ = localVarHTTPMethod
 _ = localVarPath
 return nil, nil
}
`
	if err := os.WriteFile(filepath.Join(dir, "groups.go"), []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	sdk, err := discoverBulkSDK(dir)
	if err != nil {
		t.Fatal(err)
	}
	registry := spec.Registry{{TerraformType: "verity_group", API: spec.APIResourceSpec{EndpointPath: "/groups", BulkKey: "group", RequestWrapperKey: "items", DeleteParameter: "group_name", CacheKey: "groups", ResponseCollectionKey: "group"}, Modes: []spec.Mode{spec.ModeCampus}, Operations: spec.OperationSpec{Create: true, Read: true, Delete: true}, BulkOrder: map[spec.Mode]spec.BulkOrderSpec{spec.ModeCampus: {Put: 10, Delete: 10}}}}
	bindings, err := planBulkBindings(registry, sdk)
	if err != nil {
		t.Fatal(err)
	}
	put := bindings[0].Operations["PUT"]
	if !put.Untyped || put.ValueType != "map[string]interface{}" || put.BodyType != "GroupBody" || put.BodyField != "Items" {
		t.Fatalf("binding for a body with no fields = %+v", put)
	}
	output, err := renderBulkBindings(registry, bindings)
	if err != nil {
		t.Fatal(err)
	}
	if want := "values map[string]map[string]interface{}) { request.Items = values }"; !strings.Contains(strings.Join(strings.Fields(string(output)), " "), want) {
		t.Fatalf("rendered bindings do not assign the untyped map:\n%s", output)
	}
}
