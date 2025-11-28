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
	_ resource.Resource                = &verityDiagnosticsProfileResource{}
	_ resource.ResourceWithConfigure   = &verityDiagnosticsProfileResource{}
	_ resource.ResourceWithImportState = &verityDiagnosticsProfileResource{}
)

func NewVerityDiagnosticsProfileResource() resource.Resource {
	return &verityDiagnosticsProfileResource{}
}

type verityDiagnosticsProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityDiagnosticsProfileResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	Enable               types.Bool   `tfsdk:"enable"`
	EnableSflow          types.Bool   `tfsdk:"enable_sflow"`
	FlowCollector        types.String `tfsdk:"flow_collector"`
	FlowCollectorRefType types.String `tfsdk:"flow_collector_ref_type_"`
	PollInterval         types.Int64  `tfsdk:"poll_interval"`
	VrfType              types.String `tfsdk:"vrf_type"`
}

func (r *verityDiagnosticsProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_diagnostics_profile"
}

func (r *verityDiagnosticsProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityDiagnosticsProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Diagnostics Profile",
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
			},
			"enable_sflow": schema.BoolAttribute{
				Description: "Enable sFlow for this Diagnostics Profile",
				Optional:    true,
			},
			"flow_collector": schema.StringAttribute{
				Description: "Flow Collector for this Diagnostics Profile",
				Optional:    true,
			},
			"flow_collector_ref_type_": schema.StringAttribute{
				Description: "Object type for flow_collector field",
				Optional:    true,
			},
			"poll_interval": schema.Int64Attribute{
				Description: "The sampling rate for sFlow polling (seconds)",
				Optional:    true,
			},
			"vrf_type": schema.StringAttribute{
				Description: "Management or Underlay",
				Optional:    true,
			},
		},
	}
}

func (r *verityDiagnosticsProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityDiagnosticsProfileResourceModel
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
	diagnosticsProfileProps := &openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "FlowCollector", APIField: &diagnosticsProfileProps.FlowCollector, TFValue: plan.FlowCollector},
		{FieldName: "FlowCollectorRefType", APIField: &diagnosticsProfileProps.FlowCollectorRefType, TFValue: plan.FlowCollectorRefType},
		{FieldName: "VrfType", APIField: &diagnosticsProfileProps.VrfType, TFValue: plan.VrfType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &diagnosticsProfileProps.Enable, TFValue: plan.Enable},
		{FieldName: "EnableSflow", APIField: &diagnosticsProfileProps.EnableSflow, TFValue: plan.EnableSflow},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "PollInterval", APIField: &diagnosticsProfileProps.PollInterval, TFValue: plan.PollInterval},
	})

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "diagnostics_profile", name, *diagnosticsProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_profiles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityDiagnosticsProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityDiagnosticsProfileResourceModel
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

	diagnosticsProfileName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("diagnostics_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping diagnostics profile %s verification – trusting recent successful API operation", diagnosticsProfileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching diagnostics profiles for verification of %s", diagnosticsProfileName))

	type DiagnosticsProfilesResponse struct {
		DiagnosticsProfile map[string]interface{} `json:"diagnostics_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "diagnostics_profiles", diagnosticsProfileName,
		func() (DiagnosticsProfilesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch diagnostics profiles")
			respAPI, err := r.client.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(ctx).Execute()
			if err != nil {
				return DiagnosticsProfilesResponse{}, fmt.Errorf("error reading diagnostics profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res DiagnosticsProfilesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return DiagnosticsProfilesResponse{}, fmt.Errorf("failed to decode diagnostics profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d diagnostics profiles", len(res.DiagnosticsProfile)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Diagnostics Profile %s", diagnosticsProfileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for diagnostics profile with name: %s", diagnosticsProfileName))

	diagnosticsProfileData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.DiagnosticsProfile,
		diagnosticsProfileName,
		func(data interface{}) (string, bool) {
			if diagnosticsProfile, ok := data.(map[string]interface{}); ok {
				if name, ok := diagnosticsProfile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Diagnostics Profile with name '%s' not found in API response", diagnosticsProfileName))
		resp.State.RemoveResource(ctx)
		return
	}

	diagnosticsProfileMap, ok := diagnosticsProfileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Diagnostics Profile Data",
			fmt.Sprintf("Diagnostics Profile data is not in expected format for %s", diagnosticsProfileName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found diagnostics profile '%s' under API key '%s'", diagnosticsProfileName, actualAPIName))

	state.Name = utils.MapStringFromAPI(diagnosticsProfileMap["name"])

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"flow_collector":           &state.FlowCollector,
		"flow_collector_ref_type_": &state.FlowCollectorRefType,
		"vrf_type":                 &state.VrfType,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(diagnosticsProfileMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":       &state.Enable,
		"enable_sflow": &state.EnableSflow,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(diagnosticsProfileMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"poll_interval": &state.PollInterval,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(diagnosticsProfileMap[apiKey])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityDiagnosticsProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityDiagnosticsProfileResourceModel

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
	diagnosticsProfileProps := openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { diagnosticsProfileProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.VrfType, state.VrfType, func(v *string) { diagnosticsProfileProps.VrfType = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { diagnosticsProfileProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableSflow, state.EnableSflow, func(v *bool) { diagnosticsProfileProps.EnableSflow = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.PollInterval, state.PollInterval, func(v *openapi.NullableInt32) { diagnosticsProfileProps.PollInterval = *v }, &hasChanges)

	// Handle FlowCollector and FlowCollectorRefType fields using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.FlowCollector, state.FlowCollector, plan.FlowCollectorRefType, state.FlowCollectorRefType,
		func(v *string) { diagnosticsProfileProps.FlowCollector = v },
		func(v *string) { diagnosticsProfileProps.FlowCollectorRefType = v },
		"flow_collector", "flow_collector_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "diagnostics_profile", name, diagnosticsProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityDiagnosticsProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityDiagnosticsProfileResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "diagnostics_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityDiagnosticsProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
