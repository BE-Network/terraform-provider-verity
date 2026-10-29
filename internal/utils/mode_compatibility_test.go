package utils

import (
	"fmt"
	"strings"
	"testing"

	"terraform-provider-verity/internal/registry"
)

func TestAPIVersionCompatibilityUsesEmbeddedSelection(t *testing.T) {
	version, err := registry.Version()
	if err != nil {
		t.Fatal(err)
	}
	if GetSupportedAPIVersionString() != version.String() {
		t.Fatal("supported API version differs from embedded selection")
	}
	for _, server := range []string{version.String(), version.String() + ".0.269"} {
		if err := ValidateAPIVersion(server); err != nil {
			t.Errorf("matching API version %s refused: %v", server, err)
		}
	}
	for _, server := range []string{
		fmt.Sprintf("%d.%d", version.Major+1, version.Minor),
		fmt.Sprintf("%d.%d", version.Major, version.Minor+1),
		"invalid",
	} {
		if err := ValidateAPIVersion(server); err == nil || !strings.Contains(err.Error(), version.String()) {
			t.Errorf("mismatch diagnostic for %q = %v", server, err)
		}
	}
}

func TestResourceCompatibilityForSwaggerModeChanges(t *testing.T) {
	tests := []struct {
		resource string
		mode     string
		want     bool
	}{
		{"verity_sfp_breakout", "datacenter", true},
		{"verity_sfp_breakout", "campus", true},
		{"verity_plane", "datacenter", true},
		{"verity_plane", "campus", false},
		{"verity_rack", "datacenter", true},
		{"verity_rack", "campus", false},
		{"verity_operation_stage", "datacenter", true},
		{"verity_operation_stage", "campus", true},
		{"verity_ipv4_list", "datacenter", true},
		{"verity_ipv4_list", "campus", false},
		{"verity_ipv6_list", "datacenter", true},
		{"verity_ipv6_list", "campus", false},
	}

	for _, tt := range tests {
		t.Run(tt.resource+"/"+tt.mode, func(t *testing.T) {
			if got := IsResourceCompatibleWithMode(tt.resource, tt.mode); got != tt.want {
				t.Errorf("IsResourceCompatibleWithMode(%q, %q) = %t, want %t", tt.resource, tt.mode, got, tt.want)
			}
		})
	}
}
