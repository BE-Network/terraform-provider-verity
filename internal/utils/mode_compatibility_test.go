package utils

import "testing"

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
