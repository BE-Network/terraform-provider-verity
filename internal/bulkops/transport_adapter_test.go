package bulkops

import (
	"encoding/json"
	"testing"

	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/openapi"
)

func TestTransportPreparerUsesGeneratedAdapter(t *testing.T) {
	tests := []struct {
		operation string
		object    transport.WireObject
		want      string
	}{
		{
			operation: "PUT",
			object: transport.WireObject{
				"name":      transport.String("edge-list"),
				"enable":    transport.Bool(false),
				"ipv4_list": transport.String(""),
			},
			want: `{"ipv4_list_filter":{"edge-list":{"enable":false,"ipv4_list":"","name":"edge-list"}}}`,
		},
		{
			operation: "PATCH",
			object:    transport.WireObject{"enable": transport.Bool(false)},
			want:      `{"ipv4_list_filter":{"edge-list":{"enable":false}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.operation, func(t *testing.T) {
			value, err := transport.GeneratedAdapters["verity_ipv4_list"].ResourceValue(tt.object)
			if err != nil {
				t.Fatalf("convert resource: %v", err)
			}
			config := resourceRegistry["ipv4_list"]
			request, err := (&Manager{}).createRequestPreparerWithError(config, tt.operation)(map[string]interface{}{"edge-list": value})
			if err != nil {
				t.Fatalf("prepare request: %v", err)
			}
			if _, ok := request.(*openapi.Ipv4listsPutRequest); !ok {
				t.Fatalf("request type = %T, want *openapi.Ipv4listsPutRequest", request)
			}
			got, err := json.Marshal(request)
			if err != nil {
				t.Fatalf("marshal request: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("%s wire JSON = %s, want %s", tt.operation, got, tt.want)
			}
		})
	}
}

func TestTransportPreparerRejectsUnconvertedWireObjects(t *testing.T) {
	config := resourceRegistry["ipv4_list"]
	preparer := (&Manager{}).createRequestPreparerWithError(config, "PUT")

	_, err := preparer(map[string]interface{}{
		"wire": transport.WireObject{"name": transport.String("wire")},
		"typed": openapi.Ipv4listsPutRequestIpv4ListFilterValue{
			Name: openapi.PtrString("typed"),
		},
	})
	if err == nil {
		t.Fatal("expected unconverted wire object to fail")
	}
}
