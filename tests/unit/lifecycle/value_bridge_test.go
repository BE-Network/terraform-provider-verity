package lifecycle

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	_ provider.Provider = &valueBridgeProvider{}
	_ resource.Resource = &valueBridgeResource{}
)

type valueBridgeProvider struct{}

func (p *valueBridgeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bridge"
}

func (p *valueBridgeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *valueBridgeProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *valueBridgeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{func() resource.Resource { return &valueBridgeResource{} }}
}

func (p *valueBridgeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

type valueBridgeResource struct{}

func (r *valueBridgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_value_bridge"
}

func (r *valueBridgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		Attributes: map[string]resourceschema.Attribute{
			"name":                     resourceschema.StringAttribute{Required: true},
			"optional_computed_string": resourceschema.StringAttribute{Optional: true, Computed: true},
			"optional_computed_bool":   resourceschema.BoolAttribute{Optional: true, Computed: true},
			"optional_computed_int64":  resourceschema.Int64Attribute{Optional: true, Computed: true},
			"optional_computed_number": resourceschema.NumberAttribute{Optional: true, Computed: true},
			"nullable_number":          resourceschema.NumberAttribute{Optional: true},
			"computed_string":          resourceschema.StringAttribute{Computed: true},
			"computed_number":          resourceschema.NumberAttribute{Computed: true},
		},
		Blocks: map[string]resourceschema.Block{
			"singleton": resourceschema.ListNestedBlock{
				NestedObject: resourceschema.NestedBlockObject{
					Attributes: map[string]resourceschema.Attribute{
						"label": resourceschema.StringAttribute{Required: true},
					},
					Blocks: map[string]resourceschema.Block{
						"second_level": resourceschema.ListNestedBlock{
							NestedObject: resourceschema.NestedBlockObject{Attributes: map[string]resourceschema.Attribute{
								"enabled": resourceschema.BoolAttribute{Required: true},
							}},
						},
					},
				},
			},
			"indexed_items": resourceschema.ListNestedBlock{
				NestedObject: resourceschema.NestedBlockObject{Attributes: map[string]resourceschema.Attribute{
					"index": resourceschema.Int64Attribute{Required: true},
					"label": resourceschema.StringAttribute{Required: true},
				}},
			},
		},
	}
}

func (r *valueBridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config, plan types.Object
	var name types.String
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := validateBridgeConfigAndPlan(config, plan, name); err != nil {
		resp.Diagnostics.AddError("Generic value bridge failed", err.Error())
		return
	}

	state, diags := bridgeState(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("computed_string"), types.StringValue("server-computed"))...)
}

func (r *valueBridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("computed_number"), types.NumberValue(big.NewFloat(42.5)))...)
}

func (r *valueBridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state, diags := bridgeState(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *valueBridgeResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func validateBridgeConfigAndPlan(config, plan types.Object, name types.String) error {
	if config.IsNull() || config.IsUnknown() || plan.IsNull() || plan.IsUnknown() {
		return fmt.Errorf("root config and plan must be known objects")
	}

	if name.ValueString() != "bridge-test" {
		return fmt.Errorf("unexpected generic name value %q", name.ValueString())
	}

	computed, ok := plan.Attributes()["computed_string"].(types.String)
	if !ok || !computed.IsUnknown() {
		return fmt.Errorf("computed_string must remain unknown in the generic plan tree")
	}
	nullable, ok := config.Attributes()["nullable_number"].(types.Number)
	if !ok || !nullable.IsNull() {
		return fmt.Errorf("nullable_number must remain null in the generic config tree")
	}

	singleton, ok := plan.Attributes()["singleton"].(types.List)
	if !ok || singleton.IsNull() || singleton.IsUnknown() || len(singleton.Elements()) != 1 {
		return fmt.Errorf("singleton block was not decoded as a known one-element generic list")
	}
	singletonObject, ok := singleton.Elements()[0].(types.Object)
	if !ok {
		return fmt.Errorf("singleton list element was not decoded as a generic object")
	}
	secondLevel, ok := singletonObject.Attributes()["second_level"].(types.List)
	if !ok || len(secondLevel.Elements()) != 1 {
		return fmt.Errorf("second-level block was not decoded as a generic list")
	}
	if _, ok := secondLevel.Elements()[0].(types.Object); !ok {
		return fmt.Errorf("second-level list element was not decoded as a generic object")
	}

	indexed, ok := plan.Attributes()["indexed_items"].(types.List)
	if !ok || len(indexed.Elements()) != 2 {
		return fmt.Errorf("indexed list was not decoded as two generic objects")
	}
	return nil
}

func bridgeState(ctx context.Context, plan types.Object) (types.Object, diag.Diagnostics) {
	attributes := plan.Attributes()
	attributes["computed_string"] = types.StringValue("server-computed")
	attributes["computed_number"] = types.NumberValue(big.NewFloat(42.5))
	return types.ObjectValue(plan.AttributeTypes(ctx), attributes)
}

func valueBridgeProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"bridge": providerserver.NewProtocol6WithError(&valueBridgeProvider{}),
	}
}

func TestFrameworkValueBridge(t *testing.T) {
	config := `
provider "bridge" {}

resource "bridge_value_bridge" "test" {
  name                     = "bridge-test"
  optional_computed_string = "configured"
  optional_computed_bool   = false
  optional_computed_int64  = 0
  optional_computed_number = 1.25
  nullable_number          = null

  singleton {
    label = "singleton"

    second_level {
      enabled = false
    }
  }

  indexed_items {
    index = 1
    label = "one"
  }

  indexed_items {
    index = 2
    label = "two"
  }
}
`

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: valueBridgeProviderFactories(),
		Steps: []fwresource.TestStep{{
			Config: config,
			Check: fwresource.ComposeAggregateTestCheckFunc(
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "optional_computed_string", "configured"),
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "optional_computed_bool", "false"),
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "optional_computed_int64", "0"),
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "computed_string", "server-computed"),
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "computed_number", "42.5"),
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "singleton.#", "1"),
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "singleton.0.second_level.#", "1"),
				fwresource.TestCheckResourceAttr("bridge_value_bridge.test", "indexed_items.#", "2"),
			),
		}},
	})
}
