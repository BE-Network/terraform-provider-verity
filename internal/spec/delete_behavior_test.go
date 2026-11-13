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

func TestCreateBehaviorValidation(t *testing.T) {
	base, err := registry.Lookup("verity_sfp_breakout")
	if err != nil {
		t.Fatal(err)
	}
	if base.CreateBehavior != spec.CreateBehaviorAdoptExisting {
		t.Fatalf("verity_sfp_breakout create_behavior = %q, want adopt_existing", base.CreateBehavior)
	}
	if err := base.Validate(); err != nil {
		t.Fatalf("adopt_existing on an update-only resource must be valid: %v", err)
	}

	creatable := base
	creatable.Operations.Create = true
	if err := creatable.Validate(); err == nil || !strings.Contains(err.Error(), "create_behavior adopt_existing") {
		t.Fatalf("adopt_existing on a creatable resource must be rejected, got %v", err)
	}

	unknown := base
	unknown.CreateBehavior = "clone"
	if err := unknown.Validate(); err == nil || !strings.Contains(err.Error(), "unknown create_behavior") {
		t.Fatalf("an unknown create_behavior must be rejected, got %v", err)
	}

	all, err := registry.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, resource := range all {
		if resource.CreateBehavior != spec.CreateBehaviorAPI && resource.TerraformType != "verity_sfp_breakout" {
			t.Fatalf("%s has create_behavior %q; only verity_sfp_breakout opts in", resource.TerraformType, resource.CreateBehavior)
		}
	}
}
