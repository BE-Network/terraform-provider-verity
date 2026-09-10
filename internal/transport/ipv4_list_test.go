package transport

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIPv4ListAdapterBuildPutPreservesKnownZeroValues(t *testing.T) {
	request, err := (IPv4ListAdapter{}).BuildPut(map[string]WireObject{
		"edge-list": {
			"name":      String("edge-list"),
			"enable":    Bool(false),
			"ipv4_list": String(""),
		},
	})
	if err != nil {
		t.Fatalf("BuildPut() error = %v", err)
	}
	got, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	want := `{"ipv4_list_filter":{"edge-list":{"enable":false,"ipv4_list":"","name":"edge-list"}}}`
	if string(got) != want {
		t.Fatalf("PUT wire JSON = %s, want %s", got, want)
	}
}

func TestIPv4ListAdapterBuildPatchOmitsAbsentFields(t *testing.T) {
	request, err := (IPv4ListAdapter{}).BuildPatch(map[string]WireObject{
		"edge-list": {"enable": Bool(false)},
	})
	if err != nil {
		t.Fatalf("BuildPatch() error = %v", err)
	}
	got, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	want := `{"ipv4_list_filter":{"edge-list":{"enable":false}}}`
	if string(got) != want {
		t.Fatalf("PATCH wire JSON = %s, want %s", got, want)
	}
}

func TestWireObjectPreservesExplicitNull(t *testing.T) {
	encoded, err := json.Marshal(WireObject{"nullable_field": Null()})
	if err != nil {
		t.Fatalf("marshal wire object: %v", err)
	}
	if got, want := string(encoded), `{"nullable_field":null}`; got != want {
		t.Fatalf("wire JSON = %s, want %s", got, want)
	}

	_, err = (IPv4ListAdapter{}).BuildPatch(map[string]WireObject{"edge-list": {"ipv4_list": Null()}})
	if err == nil || !strings.Contains(err.Error(), "explicit null") {
		t.Fatalf("BuildPatch() error = %v, want explicit-null diagnostic", err)
	}
}
