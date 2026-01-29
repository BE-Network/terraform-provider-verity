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
	_ resource.Resource                = &veritySflowCollectorResource{}
	_ resource.ResourceWithConfigure   = &veritySflowCollectorResource{}
	_ resource.ResourceWithImportState = &veritySflowCollectorResource{}
	_ resource.ResourceWithModifyPlan  = &veritySflowCollectorResource{}
)

const sflowCollectorResourceType = "sflowcollectors"

func NewVeritySflowCollectorResource() resource.Resource {
	return &veritySflowCollectorResource{}
}

type veritySflowCollectorResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type veritySflowCollectorResourceModel struct {
	Name   types.String `tfsdk:"name"`
	Enable types.Bool   `tfsdk:"enable"`
	Ip     types.String `tfsdk:"ip"`
	Port   types.Int64  `tfsdk:"port"`
}

func (r *veritySflowCollectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sflow_collector"
}

func (r *veritySflowCollectorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provCtx, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T", req.ProviderData),
		)
		return
	}

	r.provCtx = provCtx
	r.client = provCtx.client
	r.bulkOpsMgr = provCtx.bulkOpsMgr
	r.notifyOperationAdded = provCtx.NotifyOperationAdded
}

func (r *veritySflowCollectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity SFlow Collector",
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
			"ip": schema.StringAttribute{
				Description: "IP address of the sFlow Collector",
				Optional:    true,
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Port (maximum 65535)",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *veritySflowCollectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan veritySflowCollectorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySflowCollectorResourceModel
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
	sflowCollectorReq := &openapi.SflowcollectorsPutRequestSflowCollectorValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Ip", APIField: &sflowCollectorReq.Ip, TFValue: plan.Ip},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &sflowCollectorReq.Enable, TFValue: plan.Enable},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_sflow_collector", name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "Port", APIField: &sflowCollectorReq.Port, TFValue: config.Port, IsConfigured: configuredAttrs.IsConfigured("port")},
	})

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "sflow_collector", name, *sflowCollectorReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SFlow Collector %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sflow_collectors")

	var minState veritySflowCollectorResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if sflowCollectorData, exists := bulkMgr.GetResourceResponse("sflow_collector", name); exists {
			state := populateSflowCollectorState(ctx, minState, sflowCollectorData, r.provCtx.mode)
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

func (r *veritySflowCollectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySflowCollectorResourceModel
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

	sflowCollectorName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if sflowCollectorData, exists := r.bulkOpsMgr.GetResourceResponse("sflow_collector", sflowCollectorName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached sflow_collector data for %s from recent operation", sflowCollectorName))
			state = populateSflowCollectorState(ctx, state, sflowCollectorData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("sflow_collector") {
		tflog.Info(ctx, fmt.Sprintf("Skipping sflow collector %s verification – trusting recent successful API operation", sflowCollectorName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching sflow collectors for verification of %s", sflowCollectorName))

	type SflowCollectorsResponse struct {
		SflowCollector map[string]interface{} `json:"sflow_collector"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "sflow_collectors", sflowCollectorName,
		func() (SflowCollectorsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch sflow collectors")
			respAPI, err := r.client.SFlowCollectorsAPI.SflowcollectorsGet(ctx).Execute()
			if err != nil {
				return SflowCollectorsResponse{}, fmt.Errorf("error reading sflow collectors: %v", err)
			}
			defer respAPI.Body.Close()

			var res SflowCollectorsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SflowCollectorsResponse{}, fmt.Errorf("failed to decode sflow collectors response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d sflow collectors", len(res.SflowCollector)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read SFlow Collector %s", sflowCollectorName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for sflow collector with name: %s", sflowCollectorName))

	sflowCollectorData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.SflowCollector,
		sflowCollectorName,
		func(data interface{}) (string, bool) {
			if sflowCollector, ok := data.(map[string]interface{}); ok {
				if name, ok := sflowCollector["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("SFlow Collector with name '%s' not found in API response", sflowCollectorName))
		resp.State.RemoveResource(ctx)
		return
	}

	sflowCollectorMap, ok := sflowCollectorData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid SFlow Collector Data",
			fmt.Sprintf("SFlow Collector data is not in expected format for %s", sflowCollectorName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found sflow collector '%s' under API key '%s'", sflowCollectorName, actualAPIName))

	state = populateSflowCollectorState(ctx, state, sflowCollectorMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySflowCollectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySflowCollectorResourceModel

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
	var config veritySflowCollectorResourceModel
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
	sflowCollectorProps := openapi.SflowcollectorsPutRequestSflowCollectorValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_sflow_collector", name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { sflowCollectorProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Ip, state.Ip, func(v *string) { sflowCollectorProps.Ip = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { sflowCollectorProps.Enable = v }, &hasChanges)

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.Port, state.Port, configuredAttrs.IsConfigured("port"), func(v *openapi.NullableInt32) { sflowCollectorProps.Port = *v }, &hasChanges)

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "sflow_collector", name, sflowCollectorProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SFlow Collector %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sflow_collectors")

	var minState veritySflowCollectorResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if sflowCollectorData, exists := bulkMgr.GetResourceResponse("sflow_collector", name); exists {
			newState := populateSflowCollectorState(ctx, minState, sflowCollectorData, r.provCtx.mode)
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

func (r *veritySflowCollectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state veritySflowCollectorResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "sflow_collector", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SFlow Collector %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sflow_collectors")
	resp.State.RemoveResource(ctx)
}

func (r *veritySflowCollectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateSflowCollectorState(ctx context.Context, state veritySflowCollectorResourceModel, data map[string]interface{}, mode string) veritySflowCollectorResourceModel {
	const resourceType = sflowCollectorResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// String fields
	state.Ip = utils.MapStringWithMode(data, "ip", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Int64 fields
	state.Port = utils.MapInt64WithMode(data, "port", resourceType, mode)

	return state
}

func (r *veritySflowCollectorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySflowCollectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = sflowCollectorResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyStrings(
		"ip",
	)

	nullifier.NullifyInt64s(
		"port",
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
	var state veritySflowCollectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySflowCollectorResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_sflow_collector", name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "port", ConfigVal: config.Port, StateVal: state.Port},
		},
	})
}
