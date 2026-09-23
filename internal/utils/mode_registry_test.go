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

func TestResponseCollectionKeyLookups(t *testing.T) {
	cases := []struct {
		endpoint, terraformType, bulkKey, want string
	}{
		{endpoint: "tenants", want: "tenant"},
		{endpoint: "/bundles", want: "endpoint_bundle"},
		{endpoint: "no_such_endpoint", want: ""},
		{terraformType: "verity_acl_v4", want: "ipv4_filter"},
		{terraformType: "verity_acl_v6", want: "ipv6_filter"},
		{terraformType: "verity_no_such_resource", want: ""},
		{bulkKey: "tenant", want: "tenant"},
		{bulkKey: "device_voice_settings", want: "device_voice_settings"},
		{bulkKey: "acl", want: ""},
		{bulkKey: "no_such_bulk_key", want: ""},
	}
	for _, tc := range cases {
		switch {
		case tc.endpoint != "":
			if got := ResponseCollectionKeyForEndpoint(tc.endpoint); got != tc.want {
				t.Errorf("ResponseCollectionKeyForEndpoint(%q) = %q, want %q", tc.endpoint, got, tc.want)
			}
		case tc.terraformType != "":
			if got := ResponseCollectionKeyForType(tc.terraformType); got != tc.want {
				t.Errorf("ResponseCollectionKeyForType(%q) = %q, want %q", tc.terraformType, got, tc.want)
			}
		default:
			if got := ResponseCollectionKeyForBulkKey(tc.bulkKey); got != tc.want {
				t.Errorf("ResponseCollectionKeyForBulkKey(%q) = %q, want %q", tc.bulkKey, got, tc.want)
			}
		}
	}
}

func TestFieldAppliesToMode(t *testing.T) {
	cases := []struct {
		endpoint, field, mode string
		want                  bool
	}{
		{"services", "name", "datacenter", true},
		{"services", "name", "campus", true},
		{"services", "anycast_ipv4_mask", "datacenter", true},
		{"services", "anycast_ipv4_mask", "campus", false},
		{"services", "packet_priority", "campus", true},
		{"services", "packet_priority", "datacenter", false},
		{"services", "object_properties.warn_on_no_external_source", "campus", true},
		{"services", "object_properties.warn_on_no_external_source", "datacenter", false},
		{"devicevoicesettings", "bit_rate", "campus", true},
		{"devicevoicesettings", "bit_rate", "datacenter", false},
		{"services", "no_such_field", "datacenter", true},
		{"no_such_endpoint", "anycast_ipv4_mask", "campus", true},
	}
	for _, tc := range cases {
		if got := FieldAppliesToMode(tc.endpoint, tc.field, tc.mode); got != tc.want {
			t.Errorf("FieldAppliesToMode(%q, %q, %q) = %v, want %v", tc.endpoint, tc.field, tc.mode, got, tc.want)
		}
	}
}
