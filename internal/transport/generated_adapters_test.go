package transport

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGeneratedAdapterPreservesKnownZeroValues(t *testing.T) {
	value, err := GeneratedAdapters["verity_ipv4_list"].ResourceValue(WireObject{
		"name":      String("edge-list"),
		"enable":    Bool(false),
		"ipv4_list": String(""),
	})
	if err != nil {
		t.Fatalf("ResourceValue() error = %v", err)
	}
	got, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	want := `{"enable":false,"ipv4_list":"","name":"edge-list"}`
	if string(got) != want {
		t.Fatalf("PUT wire JSON = %s, want %s", got, want)
	}
}

func TestGeneratedAdapterOmitsAbsentFields(t *testing.T) {
	value, err := GeneratedAdapters["verity_ipv4_list"].ResourceValue(WireObject{
		"enable": Bool(false),
	})
	if err != nil {
		t.Fatalf("ResourceValue() error = %v", err)
	}
	got, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	want := `{"enable":false}`
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

	_, err = GeneratedAdapters["verity_ipv4_list"].ResourceValue(WireObject{"ipv4_list": Null()})
	if err == nil || !strings.Contains(err.Error(), "explicit null") {
		t.Fatalf("ResourceValue() error = %v, want explicit-null diagnostic", err)
	}
}
