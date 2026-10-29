package bulkops

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"terraform-provider-verity/openapi"
)

func TestEveryBulkBindingPreparesTypedValues(t *testing.T) {
	for key, config := range resourceRegistry {
		for _, operation := range []string{"PUT", "PATCH"} {
			preparer, requestType := config.PreparePut, config.PutRequestType
			if operation == "PATCH" {
				preparer, requestType = config.PreparePatch, config.PatchRequestType
			}
			if preparer == nil {
				if (operation == "PUT" && (config.PutFunc != nil || config.HeaderPutFunc != nil)) || (operation == "PATCH" && (config.PatchFunc != nil || config.HeaderPatchFunc != nil)) {
					t.Fatalf("%s %s has no preparer", key, operation)
				}
				continue
			}
			t.Run(key+operation, func(t *testing.T) {
				var wrapper string
				var valueType reflect.Type
				for index := 0; index < requestType.NumField(); index++ {
					field := requestType.Field(index)
					if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Map {
						wrapper = strings.Split(field.Tag.Get("json"), ",")[0]
						valueType = field.Type.Elem().Elem()
					}
				}
				if valueType == nil {
					t.Fatal("missing name-keyed wrapper")
				}
				data := map[string]interface{}{"unit": reflect.Zero(valueType).Interface()}
				request, err := (&Manager{}).createRequestPreparerWithError(config, operation)(data)
				if err != nil {
					t.Fatal(err)
				}
				if reflect.TypeOf(request) != reflect.PointerTo(requestType) {
					t.Fatalf("wrong request type %T", request)
				}
				raw, err := json.Marshal(request)
				if err != nil {
					t.Fatal(err)
				}
				var body map[string]map[string]json.RawMessage
				if err := json.Unmarshal(raw, &body); err != nil {
					t.Fatal(err)
				}
				if len(body) != 1 || len(body[wrapper]) != 1 || body[wrapper]["unit"] == nil {
					t.Fatalf("wrong prepared body %s", raw)
				}
				if _, err := (&Manager{}).createRequestPreparerWithError(config, operation)(map[string]interface{}{"unit": "wrong type"}); err == nil {
					t.Fatal("invalid value type accepted")
				}
			})
		}
	}
	if _, err := (&Manager{}).createRequestPreparerWithError(ResourceConfig{ResourceType: "missing"}, "PUT")(nil); err == nil {
		t.Fatal("missing preparer accepted")
	}
}

func TestEveryCreateBindingChecksPreExistingResources(t *testing.T) {
	raw, err := os.ReadFile("testdata/bulk_wire_golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var golden map[string]struct {
		Path       string `json:"path"`
		Collection string `json:"response_collection"`
	}
	if err := json.Unmarshal(raw, &golden); err != nil {
		t.Fatal(err)
	}
	for key, config := range resourceRegistry {
		if config.PreparePut == nil || config.HeaderSplitKey != "" {
			continue
		}
		t.Run(key, func(t *testing.T) {
			want := golden[key]
			if want.Collection == "" {
				t.Fatal("missing collection fixture")
			}
			transport := &bulkCaptureTransport{responseBody: `{"ignored":true,"` + want.Collection + `":{"existing":{}}}`}
			settings := openapi.NewConfiguration()
			settings.Servers = openapi.ServerConfigurations{{URL: "http://bulk.test"}}
			settings.HTTPClient = &http.Client{Transport: transport}
			manager := NewManager(openapi.NewAPIClient(settings), nil, nil, "datacenter")
			checker := manager.createPreExistenceChecker(config, "PUT")
			names, values, err := checker(context.Background(), []string{"existing", "new"}, map[string]interface{}{"existing": true, "new": true})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(names, []string{"new"}) || !reflect.DeepEqual(values, map[string]interface{}{"new": true}) {
				t.Fatalf("existing resource was not filtered: %v %v", names, values)
			}
			if transport.request.Method != "GET" || transport.request.URL.Path != want.Path {
				t.Fatalf("wrong existence query %s %s", transport.request.Method, transport.request.URL.Path)
			}
		})
	}
	manager := NewManager(nil, nil, nil, "datacenter")
	if _, _, err := manager.createPreExistenceChecker(ResourceConfig{ResourceType: "missing"}, "PUT")(context.Background(), []string{"unit"}, nil); err == nil {
		t.Fatal("missing pre-existence binding accepted")
	}
}

func TestGeneratedCacheRefreshKeysMatchGolden(t *testing.T) {
	raw, err := os.ReadFile("testdata/bulk_cache_keys_golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var expected []string
	if err := json.Unmarshal(raw, &expected); err != nil {
		t.Fatal(err)
	}
	keys := append([]string(nil), finalCacheRefreshKeys...)
	sort.Strings(keys)
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("cache keys = %v, want %v", keys, expected)
	}
}

func TestGeneratedPreparationErrorBecomesDiagnostic(t *testing.T) {
	manager := NewManager(nil, nil, nil, "datacenter")
	operationID := manager.AddPut(context.Background(), "gateway", "unit", "wrong value type")
	config := resourceRegistry["gateway"]
	diagnostics := manager.executeBulkOperation(context.Background(), BulkOperationConfig{
		ResourceType:            "gateway",
		OperationType:           "PUT",
		ExtractOperations:       manager.createExtractor("gateway", "PUT"),
		PrepareRequest:          manager.createRequestPreparer(config, "PUT"),
		PrepareRequestWithError: manager.createRequestPreparerWithError(config, "PUT"),
		ExecuteRequest: func(context.Context, interface{}) (*http.Response, error) {
			t.Fatal("invalid preparation must stop before SDK execution")
			return nil, nil
		},
		UpdateRecentOps: func() {},
	})
	if !diagnostics.HasError() {
		t.Fatal("invalid typed batch produced no diagnostic")
	}
	if _, exists := manager.operationErrors[operationID]; !exists {
		t.Fatal("preparation error was not recorded for the operation")
	}
}
