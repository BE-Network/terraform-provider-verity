package genericresource

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type versionZeroState struct {
	Name types.String `tfsdk:"name"`
}

type versionOneState struct {
	Name    types.String `tfsdk:"name"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func TestStateUpgraderRegistryPreservesVersionZeroSchemaAndUpgrades(t *testing.T) {
	ctx := context.Background()
	versionZeroSchema := schema.Schema{Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{Required: true},
	}}
	versionOneSchema := schema.Schema{Version: 1, Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{Required: true}, "enabled": schema.BoolAttribute{Optional: true},
	}}
	registry, err := NewStateUpgraderRegistry(1, []PriorStateDescriptor{{
		Version: 0, PriorSchema: versionZeroSchema,
		Upgrade: resource.StateUpgrader{StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
			var oldState versionZeroState
			resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.State = tfsdk.State{Schema: versionOneSchema}
			resp.Diagnostics.Append(resp.State.Set(ctx, versionOneState{Name: oldState.Name, Enabled: types.BoolValue(true)})...)
		}},
	}})
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}

	upgrader, exists := registry.UpgradeState(ctx)[0]
	if !exists || upgrader.PriorSchema == nil {
		t.Fatal("version-zero upgrader and prior schema are required")
	}
	if upgrader.PriorSchema.Type().TerraformType(ctx).Equal(versionOneSchema.Type().TerraformType(ctx)) {
		t.Fatal("prior schema must remain the version-zero schema")
	}
	raw := tftypes.NewValue(versionZeroSchema.Type().TerraformType(ctx), map[string]tftypes.Value{"name": tftypes.NewValue(tftypes.String, "edge")})
	request := resource.UpgradeStateRequest{State: &tfsdk.State{Raw: raw, Schema: versionZeroSchema}}
	response := &resource.UpgradeStateResponse{}
	upgrader.StateUpgrader(ctx, request, response)
	if response.Diagnostics.HasError() {
		t.Fatalf("upgrade diagnostics: %#v", response.Diagnostics)
	}
	var upgraded versionOneState
	if diags := response.State.Get(ctx, &upgraded); diags.HasError() {
		t.Fatalf("read upgraded state: %#v", diags)
	}
	if upgraded.Name.ValueString() != "edge" || !upgraded.Enabled.ValueBool() {
		t.Fatalf("upgraded state = %#v", upgraded)
	}
}

func TestStateUpgraderRegistryRejectsInvalidPriorVersions(t *testing.T) {
	priorSchema := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{Required: true}}}
	validUpgrade := resource.StateUpgrader{StateUpgrader: func(context.Context, resource.UpgradeStateRequest, *resource.UpgradeStateResponse) {}}
	for _, test := range []struct {
		name    string
		current int64
		prior   []PriorStateDescriptor
		want    string
	}{
		{"current version", 0, []PriorStateDescriptor{{Version: 0, PriorSchema: priorSchema, Upgrade: validUpgrade}}, "less than current"},
		{"missing callback", 1, []PriorStateDescriptor{{Version: 0, PriorSchema: priorSchema}}, "no upgrader"},
		{"duplicate version", 2, []PriorStateDescriptor{{Version: 0, PriorSchema: priorSchema, Upgrade: validUpgrade}, {Version: 0, PriorSchema: priorSchema, Upgrade: validUpgrade}}, "duplicate"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewStateUpgraderRegistry(test.current, test.prior)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}
