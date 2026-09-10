package bulkops

import (
	"context"
	"net/http"
	"testing"

	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/openapi"
)

func TestIPv4ListTransportPreparerUsesWireAdapter(t *testing.T) {
	config := resourceRegistry["ipv4_list"]
	preparer := (&Manager{}).createRequestPreparerWithError(config, "PATCH")

	request, err := preparer(map[string]interface{}{
		"edge-list": transport.WireObject{
			"enable": transport.Bool(false),
		},
	})
	if err != nil {
		t.Fatalf("prepare request: %v", err)
	}

	ipv4Request, ok := request.(*openapi.Ipv4listsPutRequest)
	if !ok {
		t.Fatalf("request type = %T, want *openapi.Ipv4listsPutRequest", request)
	}
	value := ipv4Request.GetIpv4ListFilter()["edge-list"]
	if !value.HasEnable() || value.GetEnable() {
		t.Fatalf("enable = %v (set=%v), want false and set", value.GetEnable(), value.HasEnable())
	}
	if value.HasName() || value.HasIpv4List() {
		t.Fatalf("patch request included omitted fields: %#v", value)
	}
}

func TestIPv4ListTransportPreparerRejectsMixedBatch(t *testing.T) {
	config := resourceRegistry["ipv4_list"]
	preparer := (&Manager{}).createRequestPreparerWithError(config, "PUT")

	_, err := preparer(map[string]interface{}{
		"wire": transport.WireObject{"name": transport.String("wire")},
		"legacy": openapi.Ipv4listsPutRequestIpv4ListFilterValue{
			Name: openapi.PtrString("legacy"),
		},
	})
	if err == nil {
		t.Fatal("expected mixed batch to fail")
	}
}

func TestIPv4ListTransportPreparationErrorBecomesDiagnostic(t *testing.T) {
	manager := NewManager(nil, nil, nil, "datacenter")
	operationID := manager.AddPut(context.Background(), "ipv4_list", "edge-list", transport.WireObject{
		"name": transport.Null(),
	})
	config := resourceRegistry["ipv4_list"]

	diagnostics := manager.executeBulkOperation(context.Background(), BulkOperationConfig{
		ResourceType:            "ipv4_list",
		OperationType:           "PUT",
		ExtractOperations:       manager.createExtractor("ipv4_list", "PUT"),
		PrepareRequest:          manager.createRequestPreparer(config, "PUT"),
		PrepareRequestWithError: manager.createRequestPreparerWithError(config, "PUT"),
		ExecuteRequest: func(context.Context, interface{}) (*http.Response, error) {
			t.Fatal("request execution must not occur when preparation fails")
			return nil, nil
		},
		UpdateRecentOps: func() {},
	})
	if !diagnostics.HasError() {
		t.Fatal("expected request preparation diagnostic")
	}
	if _, exists := manager.operationErrors[operationID]; !exists {
		t.Fatal("expected failed operation to retain the preparation error")
	}
}
