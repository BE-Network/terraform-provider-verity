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
