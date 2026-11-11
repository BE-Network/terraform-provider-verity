package genericresource

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
)

func settingsResource(t *testing.T, behavior spec.DeleteBehavior) (*Resource, tfsdk.State) {
	t.Helper()
	runtime := &readRuntime{mode: "datacenter"}
	r := secretResource(t, runtime)
	r.spec.TerraformType = "verity_sfp_breakout"
	r.spec.Operations = spec.OperationSpec{Read: true, Update: true}
	r.spec.DeleteBehavior = behavior
	state := tfsdk.State{Schema: r.schema}
	if diags := r.setState(context.Background(), &state, secretValues(map[string]attr.Value{"name": types.StringValue("SFP Breakouts")})); diags.HasError() {
		t.Fatal(diags)
	}
	return r, state
}

func TestDeleteWithForgetPolicyRemovesStateOnly(t *testing.T) {
	t.Parallel()
	r, state := settingsResource(t, spec.DeleteBehaviorForget)
	resp := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	if resp.Diagnostics.WarningsCount() != 1 || !strings.Contains(resp.Diagnostics[0].Detail(), "settings remain applied") {
		t.Fatalf("want one warning saying the settings remain applied, got %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("the resource must be removed from state")
	}
}

func TestDeleteWithoutForgetPolicyStillFails(t *testing.T) {
	t.Parallel()
	r, state := settingsResource(t, spec.DeleteBehaviorAPI)
	resp := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
	if !resp.Diagnostics.HasError() || resp.Diagnostics[0].Summary() != "Delete Not Supported" {
		t.Fatalf("a resource the API can't delete, without the forget policy, must still fail: %v", resp.Diagnostics)
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("state must be kept when the delete fails")
	}
}
