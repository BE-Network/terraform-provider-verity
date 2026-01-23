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
	_ resource.Resource                = &verityDiagnosticsPortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityDiagnosticsPortProfileResource{}
	_ resource.ResourceWithImportState = &verityDiagnosticsPortProfileResource{}
	_ resource.ResourceWithModifyPlan  = &verityDiagnosticsPortProfileResource{}
)

const diagnosticsPortProfileResourceType = "diagnosticsportprofiles"

func NewVerityDiagnosticsPortProfileResource() resource.Resource {
	return &verityDiagnosticsPortProfileResource{}
}

type verityDiagnosticsPortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityDiagnosticsPortProfileResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Enable      types.Bool   `tfsdk:"enable"`
	EnableSflow types.Bool   `tfsdk:"enable_sflow"`
}

func (r *verityDiagnosticsPortProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_diagnostics_port_profile"
}

func (r *verityDiagnosticsPortProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityDiagnosticsPortProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Diagnostics Port Profile",
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
			"enable_sflow": schema.BoolAttribute{
				Description: "Enable sFlow for this Diagnostics Profile",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *verityDiagnosticsPortProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityDiagnosticsPortProfileResourceModel
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
	diagnosticsPortProfileProps := &openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", TFValue: plan.Enable, APIField: &diagnosticsPortProfileProps.Enable},
		{FieldName: "EnableSflow", TFValue: plan.EnableSflow, APIField: &diagnosticsPortProfileProps.EnableSflow},
	})

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "diagnostics_port_profile", name, *diagnosticsPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_port_profiles")

	var minState verityDiagnosticsPortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if diagnosticsPortProfileData, exists := bulkMgr.GetResourceResponse("diagnostics_port_profile", name); exists {
			state := populateDiagnosticsPortProfileState(ctx, minState, diagnosticsPortProfileData, r.provCtx.mode)
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

func (r *verityDiagnosticsPortProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityDiagnosticsPortProfileResourceModel
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

	diagnosticsPortProfileName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if diagnosticsPortProfileData, exists := r.bulkOpsMgr.GetResourceResponse("diagnostics_port_profile", diagnosticsPortProfileName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached diagnostics port profile data for %s from recent operation", diagnosticsPortProfileName))
			state = populateDiagnosticsPortProfileState(ctx, state, diagnosticsPortProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("diagnostics_port_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping diagnostics port profile %s verification – trusting recent successful API operation", diagnosticsPortProfileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching diagnostics port profiles for verification of %s", diagnosticsPortProfileName))

	type DiagnosticsPortProfilesResponse struct {
		DiagnosticsPortProfile map[string]interface{} `json:"diagnostics_port_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "diagnostics_port_profiles", diagnosticsPortProfileName,
		func() (DiagnosticsPortProfilesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch diagnostics port profiles")
			respAPI, err := r.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(ctx).Execute()
			if err != nil {
				return DiagnosticsPortProfilesResponse{}, fmt.Errorf("error reading diagnostics port profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res DiagnosticsPortProfilesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return DiagnosticsPortProfilesResponse{}, fmt.Errorf("failed to decode diagnostics port profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d diagnostics port profiles", len(res.DiagnosticsPortProfile)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Diagnostics Port Profile %s", diagnosticsPortProfileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for diagnostics port profile with name: %s", diagnosticsPortProfileName))

	diagnosticsPortProfileData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.DiagnosticsPortProfile,
		diagnosticsPortProfileName,
		func(data interface{}) (string, bool) {
			if diagnosticsPortProfile, ok := data.(map[string]interface{}); ok {
				if name, ok := diagnosticsPortProfile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Diagnostics Port Profile with name '%s' not found in API response", diagnosticsPortProfileName))
		resp.State.RemoveResource(ctx)
		return
	}

	diagnosticsPortProfileMap, ok := diagnosticsPortProfileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Diagnostics Port Profile Data",
			fmt.Sprintf("Diagnostics Port Profile data is not in expected format for %s", diagnosticsPortProfileName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found diagnostics port profile '%s' under API key '%s'", diagnosticsPortProfileName, actualAPIName))

	state = populateDiagnosticsPortProfileState(ctx, state, diagnosticsPortProfileMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityDiagnosticsPortProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityDiagnosticsPortProfileResourceModel

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
	diagnosticsPortProfileProps := openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { diagnosticsPortProfileProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { diagnosticsPortProfileProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableSflow, state.EnableSflow, func(v *bool) { diagnosticsPortProfileProps.EnableSflow = v }, &hasChanges)

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "diagnostics_port_profile", name, diagnosticsPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Port Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_port_profiles")

	var minState verityDiagnosticsPortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if diagnosticsPortProfileData, exists := bulkMgr.GetResourceResponse("diagnostics_port_profile", name); exists {
			newState := populateDiagnosticsPortProfileState(ctx, minState, diagnosticsPortProfileData, r.provCtx.mode)
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

func (r *verityDiagnosticsPortProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityDiagnosticsPortProfileResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "diagnostics_port_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityDiagnosticsPortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateDiagnosticsPortProfileState(ctx context.Context, state verityDiagnosticsPortProfileResourceModel, data map[string]interface{}, mode string) verityDiagnosticsPortProfileResourceModel {
	const resourceType = diagnosticsPortProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.EnableSflow = utils.MapBoolWithMode(data, "enable_sflow", resourceType, mode)

	return state
}

func (r *verityDiagnosticsPortProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityDiagnosticsPortProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = diagnosticsPortProfileResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable", "enable_sflow",
	)
}
