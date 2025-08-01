package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityLagResource{}
	_ resource.ResourceWithConfigure   = &verityLagResource{}
	_ resource.ResourceWithImportState = &verityLagResource{}
)

func NewVerityLagResource() resource.Resource {
	return &verityLagResource{}
}

type verityLagResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

func (r *verityLagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lag"
}

func (r *verityLagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Link Aggregation Group (LAG)",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable object.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_peer_link": schema.BoolAttribute{
				Description: "Is this a peer link LAG.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"color": schema.StringAttribute{
				Description: "UI display color.",
				Optional:    true,
			},
			"lacp": schema.BoolAttribute{
				Description: "Enable LACP.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"eth_port_profile": schema.StringAttribute{
				Description: "Ethernet port profile name.",
				Optional:    true,
			},
			"peer_link_vlan": schema.Int64Attribute{
				Description: "VLAN ID for peer link.",
				Optional:    true,
			},
			"fallback": schema.BoolAttribute{
				Description: "Enable fallback mode.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"fast_rate": schema.BoolAttribute{
				Description: "Enable fast rate.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"eth_port_profile_ref_type_": schema.StringAttribute{
				Description: "Reference type for the Ethernet port profile.",
				Optional:    true,
				Computed:    true,
			},
			"uplink": schema.BoolAttribute{
				Description: "Indicates this LAG is designated as an uplink in the case of a spineless pod. Link State Tracking will be applied to BGP Egress VLANs/Interfaces and the MCLAG Peer Link VLAN",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						// No attributes defined - object_properties is an empty object in the schema
					},
				},
			},
		},
	}
}

func (r *verityLagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provCtx, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.provCtx = provCtx
	r.client = provCtx.client
	r.bulkOpsMgr = provCtx.bulkOpsMgr
	r.notifyOperationAdded = provCtx.NotifyOperationAdded
}

type verityLagResourceModel struct {
	Name                  types.String                     `tfsdk:"name"`
	Enable                types.Bool                       `tfsdk:"enable"`
	ObjectProperties      []verityLagObjectPropertiesModel `tfsdk:"object_properties"`
	IsPeerLink            types.Bool                       `tfsdk:"is_peer_link"`
	Color                 types.String                     `tfsdk:"color"`
	Lacp                  types.Bool                       `tfsdk:"lacp"`
	EthPortProfile        types.String                     `tfsdk:"eth_port_profile"`
	PeerLinkVlan          types.Int64                      `tfsdk:"peer_link_vlan"`
	Fallback              types.Bool                       `tfsdk:"fallback"`
	FastRate              types.Bool                       `tfsdk:"fast_rate"`
	EthPortProfileRefType types.String                     `tfsdk:"eth_port_profile_ref_type_"`
	Uplink                types.Bool                       `tfsdk:"uplink"`
}

type verityLagObjectPropertiesModel struct {
	// No attributes defined - object_properties is an empty object in the schema
}

func (r *verityLagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityLagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	lagReq := openapi.NewLagsPutRequestLagValue()
	lagReq.Name = openapi.PtrString(name)

	if !plan.Enable.IsNull() {
		lagReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.IsPeerLink.IsNull() {
		lagReq.IsPeerLink = openapi.PtrBool(plan.IsPeerLink.ValueBool())
	}
	if !plan.Color.IsNull() {
		color := plan.Color.ValueString()
		if color != "" {
			lagReq.Color = openapi.PtrString(color)
		}
	}
	if !plan.Lacp.IsNull() {
		lagReq.Lacp = openapi.PtrBool(plan.Lacp.ValueBool())
	}
	if !plan.EthPortProfile.IsNull() {
		ethPortProfile := plan.EthPortProfile.ValueString()
		if ethPortProfile != "" {
			lagReq.EthPortProfile = openapi.PtrString(ethPortProfile)
		}
	}
	if !plan.PeerLinkVlan.IsNull() {
		peerLinkVlan := int32(plan.PeerLinkVlan.ValueInt64())
		lagReq.PeerLinkVlan = *openapi.NewNullableInt32(&peerLinkVlan)
	} else {
		lagReq.PeerLinkVlan = *openapi.NewNullableInt32(nil)
	}
	if !plan.Fallback.IsNull() {
		lagReq.Fallback = openapi.PtrBool(plan.Fallback.ValueBool())
	}
	if !plan.FastRate.IsNull() {
		lagReq.FastRate = openapi.PtrBool(plan.FastRate.ValueBool())
	}
	if !plan.EthPortProfileRefType.IsNull() {
		lagReq.EthPortProfileRefType = openapi.PtrString(plan.EthPortProfileRefType.ValueString())
	}
	if !plan.Uplink.IsNull() {
		lagReq.Uplink = openapi.PtrBool(plan.Uplink.ValueBool())
	}

	if len(plan.ObjectProperties) > 0 {
		lagReq.ObjectProperties = make(map[string]interface{})
	} else {
		lagReq.ObjectProperties = nil
	}

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddPut(ctx, "lag", name, *lagReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for LAG creation operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create LAG %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LAG %s creation operation completed successfully", name))
	clearCache(ctx, provCtx, "lags")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityLagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityLagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	lagName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("lag") {
		tflog.Info(ctx, fmt.Sprintf("Skipping LAG %s verification - trusting recent successful API operation", lagName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent LAG operations found, performing normal verification for %s", lagName))

	type LagsResponse struct {
		Lag map[string]map[string]interface{} `json:"lag"`
	}

	var result LagsResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		lagsData, fetchErr := getCachedResponse(ctx, r.provCtx, "lags", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch LAGs")
			respAPI, err := r.client.LAGsAPI.LagsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading LAG: %v", err)
			}
			defer respAPI.Body.Close()

			var res LagsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode LAGs response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d LAGs from API", len(res.Lag)))
			return res, nil
		})

		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch LAGs on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = lagsData.(LagsResponse)
		break
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read LAG %s", lagName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for LAG with ID: %s", lagName))
	var lagData map[string]interface{}
	exists := false

	if data, ok := result.Lag[lagName]; ok {
		lagData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found LAG directly by ID: %s", lagName))
	} else {
		for apiName, l := range result.Lag {
			if name, ok := l["name"].(string); ok && name == lagName {
				lagData = l
				lagName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found LAG with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("LAG with ID '%s' not found in API response", lagName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", lagData["name"]))

	if val, ok := lagData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else {
		state.Enable = types.BoolNull()
	}
	if val, ok := lagData["is_peer_link"].(bool); ok {
		state.IsPeerLink = types.BoolValue(val)
	} else {
		state.IsPeerLink = types.BoolNull()
	}
	if val, ok := lagData["color"].(string); ok {
		state.Color = types.StringValue(val)
	} else {
		state.Color = types.StringNull()
	}
	if val, ok := lagData["lacp"].(bool); ok {
		state.Lacp = types.BoolValue(val)
	} else {
		state.Lacp = types.BoolNull()
	}
	if val, ok := lagData["eth_port_profile"].(string); ok {
		state.EthPortProfile = types.StringValue(val)
	} else {
		state.EthPortProfile = types.StringNull()
	}
	if val, ok := lagData["peer_link_vlan"]; ok {
		switch v := val.(type) {
		case float64:
			state.PeerLinkVlan = types.Int64Value(int64(v))
		case int:
			state.PeerLinkVlan = types.Int64Value(int64(v))
		default:
			state.PeerLinkVlan = types.Int64Null()
		}
	} else {
		state.PeerLinkVlan = types.Int64Null()
	}
	if val, ok := lagData["fallback"].(bool); ok {
		state.Fallback = types.BoolValue(val)
	} else {
		state.Fallback = types.BoolNull()
	}
	if val, ok := lagData["fast_rate"].(bool); ok {
		state.FastRate = types.BoolValue(val)
	} else {
		state.FastRate = types.BoolNull()
	}
	if val, ok := lagData["eth_port_profile_ref_type_"].(string); ok {
		state.EthPortProfileRefType = types.StringValue(val)
	} else {
		state.EthPortProfileRefType = types.StringNull()
	}
	if val, ok := lagData["uplink"].(bool); ok {
		state.Uplink = types.BoolValue(val)
	} else {
		state.Uplink = types.BoolNull()
	}

	// Only set object_properties if it exists in the API response
	if _, ok := lagData["object_properties"]; ok {
		state.ObjectProperties = []verityLagObjectPropertiesModel{{}}
	} else {
		state.ObjectProperties = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityLagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityLagResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	lagReq := &openapi.LagsPutRequestLagValue{}
	hasChanges := false

	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) == 0 {
		lagReq.ObjectProperties = make(map[string]interface{})
		hasChanges = true
	}

	boolFields := []string{"enable", "is_peer_link", "lacp", "fallback", "fast_rate", "uplink"}
	for _, field := range boolFields {
		if planField, stateField := getBoolField(field, plan, state); !planField.Equal(stateField) {
			boolVal := planField.ValueBool()
			switch field {
			case "enable":
				lagReq.Enable = openapi.PtrBool(boolVal)
			case "is_peer_link":
				lagReq.IsPeerLink = openapi.PtrBool(boolVal)
			case "lacp":
				lagReq.Lacp = openapi.PtrBool(boolVal)
			case "fallback":
				lagReq.Fallback = openapi.PtrBool(boolVal)
			case "fast_rate":
				lagReq.FastRate = openapi.PtrBool(boolVal)
			case "uplink":
				lagReq.Uplink = openapi.PtrBool(boolVal)
			}
			hasChanges = true
		}
	}

	planColor, stateColor := getStringField("color", plan, state)
	if !planColor.Equal(stateColor) {
		lagReq.Color = openapi.PtrString(planColor.ValueString())
		hasChanges = true
	}

	ethPortProfileChanged := !plan.EthPortProfile.Equal(state.EthPortProfile)
	ethPortProfileRefTypeChanged := !plan.EthPortProfileRefType.Equal(state.EthPortProfileRefType)

	if ethPortProfileChanged || ethPortProfileRefTypeChanged {
		// Validate using multiple ref types supported rules
		if !utils.ValidateReferenceFields(&resp.Diagnostics,
			plan.EthPortProfile, plan.EthPortProfileRefType,
			"eth_port_profile", "eth_port_profile_ref_type_") {
			return
		}

		// For multiple ref types supported: When either field changes, always send both fields
		if !plan.EthPortProfile.IsNull() && plan.EthPortProfile.ValueString() != "" {
			lagReq.EthPortProfile = openapi.PtrString(plan.EthPortProfile.ValueString())
		} else {
			lagReq.EthPortProfile = openapi.PtrString("")
		}

		if !plan.EthPortProfileRefType.IsNull() && plan.EthPortProfileRefType.ValueString() != "" {
			lagReq.EthPortProfileRefType = openapi.PtrString(plan.EthPortProfileRefType.ValueString())
		} else {
			lagReq.EthPortProfileRefType = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.PeerLinkVlan.Equal(state.PeerLinkVlan) {
		if !plan.PeerLinkVlan.IsNull() {
			peerLinkVlan := int32(plan.PeerLinkVlan.ValueInt64())
			lagReq.PeerLinkVlan = *openapi.NewNullableInt32(&peerLinkVlan)
		} else {
			lagReq.PeerLinkVlan = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddPatch(ctx, "lag", name, *lagReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for LAG update operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update LAG %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LAG %s update operation completed successfully", name))
	clearCache(ctx, provCtx, "lags")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityLagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityLagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := state.Name.ValueString()
	operationID := r.bulkOpsMgr.AddDelete(ctx, "lag", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for LAG deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete LAG %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LAG %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "lags")
	resp.State.RemoveResource(ctx)
}

func (r *verityLagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func getBoolField(field string, plan verityLagResourceModel, state verityLagResourceModel) (types.Bool, types.Bool) {
	switch field {
	case "enable":
		return plan.Enable, state.Enable
	case "is_peer_link":
		return plan.IsPeerLink, state.IsPeerLink
	case "lacp":
		return plan.Lacp, state.Lacp
	case "fallback":
		return plan.Fallback, state.Fallback
	case "fast_rate":
		return plan.FastRate, state.FastRate
	case "uplink":
		return plan.Uplink, state.Uplink
	default:
		return types.BoolNull(), types.BoolNull()
	}
}

func getStringField(field string, plan verityLagResourceModel, state verityLagResourceModel) (types.String, types.String) {
	switch field {
	case "color":
		return plan.Color, state.Color
	case "eth_port_profile":
		return plan.EthPortProfile, state.EthPortProfile
	case "eth_port_profile_ref_type_":
		return plan.EthPortProfileRefType, state.EthPortProfileRefType
	default:
		return types.StringNull(), types.StringNull()
	}
}
