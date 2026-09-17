package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/openapi"
)

func TestGenericRuntimeGetDecodesResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/api/ipv4lists" {
			t.Fatalf("request = %s %s, want GET /api/ipv4lists", request.Method, request.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ipv4_list_filter":{"example":{"name":"example","enable":true}}}`))
	}))
	defer server.Close()

	config := openapi.NewConfiguration()
	config.Debug = true
	config.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
	client := openapi.NewAPIClient(config)
	runtime := genericRuntime{provCtx: &providerContext{client: client}}
	specification := spec.ResourceSpec{
		TerraformType: "verity_ipv4_list",
		API: spec.APIResourceSpec{
			EndpointPath:          "/ipv4lists",
			ResponseCollectionKey: "ipv4_list_filter",
		},
	}

	collection, err := runtime.get(context.Background(), specification)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got, want := collection["example"].(map[string]interface{})["name"], "example"; got != want {
		t.Errorf("decoded name = %v, want %q", got, want)
	}
}

// A resource selected by a fixed parameter must read only its own objects.
//
// The ACLs share /acls and are told apart by ip_version, which the OpenAPI
// documents declare `in: query` for every operation. The registry records it as
// fixed_headers, a name taken from the bulk manager's HeaderParams, and sending
// it as an HTTP header would reach a server that ignores it and returns an
// unfiltered or default collection. The server here answers by the query
// parameter alone, the way the API and the mock do, so a header-borne value
// reads the wrong version's objects.
func TestGenericRuntimeSendsFixedParametersInTheQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/api/acls" {
			t.Fatalf("request = %s %s, want GET /api/acls", request.Method, request.URL.Path)
		}
		if header := request.Header.Get("ip_version"); header != "" {
			t.Errorf("ip_version was sent as a header (%q); the API reads it from the query", header)
		}
		w.Header().Set("Content-Type", "application/json")
		switch request.URL.Query().Get("ip_version") {
		case "6":
			_, _ = w.Write([]byte(`{"ipv6_filter":{"v6-only":{"name":"v6-only"}}}`))
		case "4":
			_, _ = w.Write([]byte(`{"ipv4_filter":{"v4-only":{"name":"v4-only"}}}`))
		default:
			http.Error(w, "ip_version is required", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	config := openapi.NewConfiguration()
	config.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
	runtime := genericRuntime{provCtx: &providerContext{client: openapi.NewAPIClient(config)}}

	for version, want := range map[string]string{"4": "v4-only", "6": "v6-only"} {
		specification := spec.ResourceSpec{
			TerraformType: "verity_acl_v" + version,
			API: spec.APIResourceSpec{
				EndpointPath:          "/acls",
				ResponseCollectionKey: "ipv" + version + "_filter",
				FixedHeaders:          map[string]string{"ip_version": version},
			},
		}
		collection, err := runtime.get(context.Background(), specification)
		if err != nil {
			t.Fatalf("ip_version %s: get: %v", version, err)
		}
		if _, found := collection[want]; !found || len(collection) != 1 {
			t.Errorf("ip_version %s read %v, want only %q", version, collection, want)
		}
	}
}
