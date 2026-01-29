package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityLagResource{}
	_ resource.ResourceWithConfigure   = &verityLagResource{}
	_ resource.ResourceWithImportState = &verityLagResource{}
	_ resource.ResourceWithModifyPlan  = &verityLagResource{}
)

const lagResourceType = "lags"

func NewVerityLagResource() resource.Resource {
	return &verityLagResource{}
}

type verityLagResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
			},
			"is_peer_link": schema.BoolAttribute{
				Description: "Is this a peer link LAG.",
				Optional:    true,
				Computed:    true,
			},
			"color": schema.StringAttribute{
				Description: "UI display color.",
				Optional:    true,
				Computed:    true,
			},
			"lacp": schema.BoolAttribute{
				Description: "Enable LACP.",
				Optional:    true,
				Computed:    true,
			},
			"eth_port_profile": schema.StringAttribute{
				Description: "Ethernet port profile name.",
				Optional:    true,
				Computed:    true,
			},
			"peer_link_vlan": schema.Int64Attribute{
				Description: "VLAN ID for peer link.",
				Optional:    true,
				Computed:    true,
			},
			"fallback": schema.BoolAttribute{
				Description: "Enable fallback mode.",
				Optional:    true,
				Computed:    true,
			},
			"fast_rate": schema.BoolAttribute{
				Description: "Enable fast rate.",
				Optional:    true,
				Computed:    true,
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

	var config verityLagResourceModel
	diags = req.Config.Get(ctx, &config)
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
	lagReq := &openapi.LagsPutRequestLagValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Color", APIField: &lagReq.Color, TFValue: plan.Color},
		{FieldName: "EthPortProfile", APIField: &lagReq.EthPortProfile, TFValue: plan.EthPortProfile},
		{FieldName: "EthPortProfileRefType", APIField: &lagReq.EthPortProfileRefType, TFValue: plan.EthPortProfileRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &lagReq.Enable, TFValue: plan.Enable},
		{FieldName: "IsPeerLink", APIField: &lagReq.IsPeerLink, TFValue: plan.IsPeerLink},
		{FieldName: "Lacp", APIField: &lagReq.Lacp, TFValue: plan.Lacp},
		{FieldName: "Fallback", APIField: &lagReq.Fallback, TFValue: plan.Fallback},
		{FieldName: "FastRate", APIField: &lagReq.FastRate, TFValue: plan.FastRate},
		{FieldName: "Uplink", APIField: &lagReq.Uplink, TFValue: plan.Uplink},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_lag", name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "PeerLinkVlan", APIField: &lagReq.PeerLinkVlan, TFValue: config.PeerLinkVlan, IsConfigured: configuredAttrs.IsConfigured("peer_link_vlan")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		lagReq.ObjectProperties = make(map[string]interface{})
	} else {
		lagReq.ObjectProperties = nil
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "lag", name, *lagReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LAG %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "lags")

	var minState verityLagResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if lagData, exists := bulkMgr.GetResourceResponse("lag", name); exists {
			state := populateLagState(ctx, minState, lagData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if lagData, exists := r.bulkOpsMgr.GetResourceResponse("lag", lagName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached lag data for %s from recent operation", lagName))
			state = populateLagState(ctx, state, lagData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("lag") {
		tflog.Info(ctx, fmt.Sprintf("Skipping LAG %s verification – trusting recent successful API operation", lagName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching LAG for verification of %s", lagName))

	type LagsResponse struct {
		Lag map[string]map[string]interface{} `json:"lag"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "lags", lagName,
		func() (LagsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch LAGs")
			respAPI, err := r.client.LAGsAPI.LagsGet(ctx).Execute()
			if err != nil {
				return LagsResponse{}, fmt.Errorf("error reading LAG: %v", err)
			}
			defer respAPI.Body.Close()

			var res LagsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return LagsResponse{}, fmt.Errorf("failed to decode LAGs response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d LAGs from API", len(res.Lag)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read LAG %s", lagName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for LAG with name: %s", lagName))

	lagData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Lag,
		lagName,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("LAG with name '%s' not found in API response", lagName))
		resp.State.RemoveResource(ctx)
		return
	}

	lagMap, ok := lagData, true
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid LAG Data",
			fmt.Sprintf("LAG data is not in expected format for %s", lagName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found LAG '%s' under API key '%s'", lagName, actualAPIName))

	state = populateLagState(ctx, state, lagMap, r.provCtx.mode)
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

	// Get config for nullable field handling
	var config verityLagResourceModel
	diags = req.Config.Get(ctx, &config)
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
	lagReq := openapi.LagsPutRequestLagValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_lag", name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { lagReq.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Color, state.Color, func(v *string) { lagReq.Color = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { lagReq.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IsPeerLink, state.IsPeerLink, func(v *bool) { lagReq.IsPeerLink = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Lacp, state.Lacp, func(v *bool) { lagReq.Lacp = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Fallback, state.Fallback, func(v *bool) { lagReq.Fallback = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.FastRate, state.FastRate, func(v *bool) { lagReq.FastRate = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Uplink, state.Uplink, func(v *bool) { lagReq.Uplink = v }, &hasChanges)

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.PeerLinkVlan, state.PeerLinkVlan, configuredAttrs.IsConfigured("peer_link_vlan"), func(v *openapi.NullableInt32) { lagReq.PeerLinkVlan = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) == 0 {
		lagReq.ObjectProperties = make(map[string]interface{})
		hasChanges = true
	}

	// Handle EthPortProfile and EthPortProfileRefType using "Many ref types supported" pattern
	if !utils.HandleMultipleRefTypesSupported(
		plan.EthPortProfile, state.EthPortProfile, plan.EthPortProfileRefType, state.EthPortProfileRefType,
		func(v *string) { lagReq.EthPortProfile = v },
		func(v *string) { lagReq.EthPortProfileRefType = v },
		"eth_port_profile", "eth_port_profile_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "lag", name, lagReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LAG %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "lags")

	var minState verityLagResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if lagData, exists := bulkMgr.GetResourceResponse("lag", name); exists {
			newState := populateLagState(ctx, minState, lagData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "lag", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LAG %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "lags")
	resp.State.RemoveResource(ctx)
}

func (r *verityLagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateLagState(ctx context.Context, state verityLagResourceModel, data map[string]interface{}, mode string) verityLagResourceModel {
	const resourceType = lagResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Int fields
	state.PeerLinkVlan = utils.MapInt64WithMode(data, "peer_link_vlan", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.IsPeerLink = utils.MapBoolWithMode(data, "is_peer_link", resourceType, mode)
	state.Lacp = utils.MapBoolWithMode(data, "lacp", resourceType, mode)
	state.Fallback = utils.MapBoolWithMode(data, "fallback", resourceType, mode)
	state.FastRate = utils.MapBoolWithMode(data, "fast_rate", resourceType, mode)
	state.Uplink = utils.MapBoolWithMode(data, "uplink", resourceType, mode)

	// String fields
	state.Color = utils.MapStringWithMode(data, "color", resourceType, mode)
	state.EthPortProfile = utils.MapStringWithMode(data, "eth_port_profile", resourceType, mode)
	state.EthPortProfileRefType = utils.MapStringWithMode(data, "eth_port_profile_ref_type_", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if _, ok := data["object_properties"]; ok {
			state.ObjectProperties = []verityLagObjectPropertiesModel{{}}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityLagResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityLagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = lagResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"color", "eth_port_profile", "eth_port_profile_ref_type_",
	)

	nullifier.NullifyBools(
		"enable", "is_peer_link", "lacp", "fallback", "fast_rate", "uplink",
	)

	nullifier.NullifyInt64s(
		"peer_link_vlan",
	)

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityLagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityLagResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable Int64 fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_lag", name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "peer_link_vlan", ConfigVal: config.PeerLinkVlan, StateVal: state.PeerLinkVlan},
		},
	})
}
