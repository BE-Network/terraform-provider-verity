package spec_test

import (
	"strings"
	"testing"

	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/spec"
)

func TestDeleteBehaviorValidation(t *testing.T) {
	base, err := registry.Lookup("verity_sfp_breakout")
	if err != nil {
		t.Fatal(err)
	}
	if base.DeleteBehavior != spec.DeleteBehaviorForget {
		t.Fatalf("verity_sfp_breakout delete_behavior = %q, want forget", base.DeleteBehavior)
	}
	if err := base.Validate(); err != nil {
		t.Fatalf("forget on a resource the API can't delete must be valid: %v", err)
	}

	deletable := base
	deletable.Operations.Delete = true
	deletable.API.DeleteParameter = "sfp_breakout_name"
	if err := deletable.Validate(); err == nil || !strings.Contains(err.Error(), "delete_behavior forget") {
		t.Fatalf("forget on a deletable resource must be rejected, got %v", err)
	}

	unknown := base
	unknown.DeleteBehavior = "archive"
	if err := unknown.Validate(); err == nil || !strings.Contains(err.Error(), "unknown delete_behavior") {
		t.Fatalf("an unknown delete_behavior must be rejected, got %v", err)
	}
}
