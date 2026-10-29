package bulkops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"terraform-provider-verity/openapi"
)

func TestBulkOrdersMatchGolden(t *testing.T) {
	raw, err := os.ReadFile("testdata/bulk_order_golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var want map[string]map[string][]string
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	got := map[string]map[string][]string{
		"datacenter": {"PUT": datacenterPutOrder, "PATCH": datacenterPatchOrder, "DELETE": datacenterDeleteOrder},
		"campus":     {"PUT": campusPutOrder, "PATCH": campusPatchOrder, "DELETE": campusDeleteOrder},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bulk ordering changed:\ngot: %v\nwant: %v", got, want)
	}
}

type bulkCaptureTransport struct {
	request      *http.Request
	responseBody string
	body         []byte
}

func (transport *bulkCaptureTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	transport.request = request.Clone(request.Context())
	if request.Body != nil {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		transport.body = body
	}
	body := transport.responseBody
	if body == "" {
		body = "{}"
	}
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: request}, nil
}

func TestBulkBindingsMatchWireGolden(t *testing.T) {
	raw, err := os.ReadFile("testdata/bulk_wire_golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var golden map[string]struct {
		Path            string   `json:"path"`
		Wrapper         string   `json:"wrapper"`
		DeleteParameter string   `json:"delete_parameter"`
		RequestType     string   `json:"request_type"`
		Operations      []string `json:"operations"`
	}
	if err := json.Unmarshal(raw, &golden); err != nil {
		t.Fatal(err)
	}
	if len(resourceRegistry) != len(golden) {
		t.Fatalf("binding count = %d, want %d", len(resourceRegistry), len(golden))
	}
	for key, want := range golden {
		t.Run(key, func(t *testing.T) {
			config, ok := resourceRegistry[key]
			if !ok {
				t.Fatal("missing binding")
			}
			for _, requestType := range []reflect.Type{config.PutRequestType, config.PatchRequestType} {
				if requestType == nil || requestType.Name() != want.RequestType {
					t.Fatalf("request type = %v, want %s", requestType, want.RequestType)
				}
			}
			if client, ok := config.APIClientGetter(nil).(*GenericAPIClient); !ok || client.resourceType != key {
				t.Fatal("wrong API client binding")
			}
			for _, operation := range []string{"Put", "Patch", "Delete", "Get"} {
				enabled := false
				for _, expected := range want.Operations {
					enabled = enabled || expected == operation
				}
				present := map[string]bool{
					"Put":    config.PutFunc != nil || config.HeaderPutFunc != nil,
					"Patch":  config.PatchFunc != nil || config.HeaderPatchFunc != nil,
					"Delete": config.DeleteFunc != nil || config.HeaderDeleteFunc != nil,
					"Get":    config.GetFunc != nil || config.HeaderGetFunc != nil,
				}[operation]
				if present != enabled {
					t.Fatalf("%s callback present = %t, want %t", operation, present, enabled)
				}
				if !enabled {
					continue
				}
				selectors := []string{""}
				if config.HeaderSplitKey != "" {
					selectors = []string{"4", "6"}
				}
				for _, selector := range selectors {
					t.Run(operation+selector, func(t *testing.T) {
						transport := &bulkCaptureTransport{}
						settings := openapi.NewConfiguration()
						settings.Servers = openapi.ServerConfigurations{{URL: "http://bulk.test"}}
						settings.HTTPClient = &http.Client{Transport: transport}
						client := openapi.NewAPIClient(settings)
						headers := map[string]string{config.HeaderSplitKey: selector}
						requestType := config.PutRequestType
						if operation == "Patch" {
							requestType = config.PatchRequestType
						}
						body := reflect.New(requestType).Interface()
						payload := `{"` + want.Wrapper + `":{"unit":{}}}`
						if err := json.Unmarshal([]byte(payload), body); err != nil {
							t.Fatal(err)
						}
						names := []string{"unit name", "unit+second/&"}
						ctx := context.Background()
						var response *http.Response
						var err error
						switch operation {
						case "Put":
							if config.HeaderPutFunc != nil {
								response, err = config.HeaderPutFunc(client, ctx, body, headers)
							} else {
								response, err = config.PutFunc(client, ctx, body)
							}
						case "Patch":
							if config.HeaderPatchFunc != nil {
								response, err = config.HeaderPatchFunc(client, ctx, body, headers)
							} else {
								response, err = config.PatchFunc(client, ctx, body)
							}
						case "Delete":
							if config.HeaderDeleteFunc != nil {
								response, err = config.HeaderDeleteFunc(client, ctx, names, headers)
							} else {
								response, err = config.DeleteFunc(client, ctx, names)
							}
						case "Get":
							if config.HeaderGetFunc != nil {
								response, err = config.HeaderGetFunc(client, ctx, headers)
							} else {
								response, err = config.GetFunc(client, ctx)
							}
						}
						if err != nil {
							t.Fatal(err)
						}
						if response == nil || transport.request == nil {
							t.Fatal("callback sent no request")
						}
						response.Body.Close()
						request := transport.request
						if request.Method != strings.ToUpper(operation) || request.URL.Path != want.Path {
							t.Fatalf("sent %s %s, want %s %s", request.Method, request.URL.Path, strings.ToUpper(operation), want.Path)
						}
						if operation == "Delete" && !reflect.DeepEqual(request.URL.Query()[want.DeleteParameter], names) {
							t.Fatalf("delete query = %v, want %v", request.URL.Query(), names)
						}
						if selector != "" && request.URL.Query().Get(config.HeaderSplitKey) != selector {
							t.Fatalf("missing split query in %s", request.URL)
						}
						if operation == "Put" || operation == "Patch" {
							var content map[string]map[string]json.RawMessage
							if err := json.Unmarshal(transport.body, &content); err != nil {
								t.Fatal(err)
							}
							if len(content) != 1 || len(content[want.Wrapper]) != 1 || content[want.Wrapper]["unit"] == nil {
								t.Fatalf("wrong request wrapper: %s", transport.body)
							}
						}
					})
				}
			}
		})
	}
}

func TestACLResponseExtraction(t *testing.T) {
	config := resourceRegistry["acl"]
	for _, test := range []struct{ selector, collection string }{{"4", "ipv4_filter"}, {"6", "ipv6_filter"}, {"", "ipv4_filter"}} {
		expected := map[string]interface{}{"unit": test.selector}
		actual, err := config.HeaderResponseExtractor(map[string]interface{}{test.collection: expected}, map[string]string{"ip_version": test.selector})
		if err != nil || !reflect.DeepEqual(actual, expected) {
			t.Fatalf("selector %q extracted %v, %v", test.selector, actual, err)
		}
		actual, err = config.HeaderResponseExtractor(map[string]interface{}{test.collection: "invalid"}, map[string]string{"ip_version": test.selector})
		if err != nil || actual == nil || len(actual) != 0 {
			t.Fatalf("invalid collection extracted %v, %v", actual, err)
		}
	}
}
