package utils

import (
	"strings"
	"testing"

	"terraform-provider-verity/internal/spec"
)

func TestHeaderSplitKeyRejectsMultipleFixedHeaders(t *testing.T) {
	resources := []spec.ResourceSpec{{
		TerraformType: "verity_example",
		API: spec.APIResourceSpec{
			EndpointPath: "/examples", BulkKey: "example",
			FixedHeaders: map[string]string{"ip_version": "4", "region": "east"},
		},
	}}
	_, err := HeaderSplitKey(resources, "example")
	if err == nil || !strings.Contains(err.Error(), "supports one split key") {
		t.Fatalf("HeaderSplitKey() error = %v, want rejection of multiple fixed headers", err)
	}
	if !strings.Contains(err.Error(), "ip_version, region") {
		t.Fatalf("error should name the headers in a stable order, got %v", err)
	}
}

func TestHeaderSplitKeyReadsTheRegistry(t *testing.T) {
	splitKey, err := HeaderSplitKeyForBulkKey("acl")
	if err != nil || splitKey != "ip_version" {
		t.Fatalf("HeaderSplitKeyForBulkKey(\"acl\") = %q, %v; want ip_version", splitKey, err)
	}
	splitKey, err = HeaderSplitKeyForBulkKey("tenant")
	if err != nil || splitKey != "" {
		t.Fatalf("HeaderSplitKeyForBulkKey(\"tenant\") = %q, %v; want no split key", splitKey, err)
	}
	if _, err = HeaderSplitKeyForBulkKey("no_such_bulk_key"); err == nil {
		t.Fatal("an unknown bulk key was accepted")
	}
}

func TestResponseCollectionKeyForBulkKey(t *testing.T) {
	cases := map[string]string{
		"tenant":                "tenant",
		"device_voice_settings": "device_voice_settings",
		"acl":                   "",
		"no_such_bulk_key":      "",
	}
	for bulkKey, want := range cases {
		if got := ResponseCollectionKeyForBulkKey(bulkKey); got != want {
			t.Errorf("ResponseCollectionKeyForBulkKey(%q) = %q, want %q", bulkKey, got, want)
		}
	}
}
