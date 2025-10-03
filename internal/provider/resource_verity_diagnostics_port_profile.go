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

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityDiagnosticsPortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityDiagnosticsPortProfileResource{}
	_ resource.ResourceWithImportState = &verityDiagnosticsPortProfileResource{}
)

func NewVerityDiagnosticsPortProfileResource() resource.Resource {
	return &verityDiagnosticsPortProfileResource{}
}

type verityDiagnosticsPortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
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
			},
			"enable_sflow": schema.BoolAttribute{
				Description: "Enable sFlow for this Diagnostics Profile",
				Optional:    true,
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "diagnostics_port_profile", name, diagnosticsPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_port_profiles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	state.Name = utils.MapStringFromAPI(diagnosticsPortProfileMap["name"])

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":       &state.Enable,
		"enable_sflow": &state.EnableSflow,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(diagnosticsPortProfileMap[apiKey])
	}

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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "diagnostics_port_profile", name, diagnosticsPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Port Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_port_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "diagnostics_port_profile", name, nil, &resp.Diagnostics)
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
