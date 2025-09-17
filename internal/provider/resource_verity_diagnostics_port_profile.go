package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	if !plan.Enable.IsNull() {
		diagnosticsPortProfileProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.EnableSflow.IsNull() {
		diagnosticsPortProfileProps.EnableSflow = openapi.PtrBool(plan.EnableSflow.ValueBool())
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "diagnostics_port_profile", name, *diagnosticsPortProfileProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for diagnostics port profile creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Diagnostics Port Profile %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Diagnostics Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_port_profiles")
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("diagnosticsportprofile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping diagnostics port profile %s verification – trusting recent successful API operation", diagnosticsPortProfileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching diagnostics port profiles for verification of %s", diagnosticsPortProfileName))

	type DiagnosticsPortProfilesResponse struct {
		DiagnosticsPortProfile map[string]interface{} `json:"diagnostics_port_profile"`
	}

	var result DiagnosticsPortProfilesResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		diagnosticsPortProfilesData, fetchErr := getCachedResponse(ctx, r.provCtx, "diagnosticsportprofiles", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch diagnostics port profiles")
			respAPI, err := r.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading diagnostics port profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res DiagnosticsPortProfilesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode diagnostics port profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d diagnostics port profiles", len(res.DiagnosticsPortProfile)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch diagnostics port profiles on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = diagnosticsPortProfilesData.(DiagnosticsPortProfilesResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Diagnostics Port Profile %s", diagnosticsPortProfileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for diagnostics port profile with ID: %s", diagnosticsPortProfileName))
	var diagnosticsPortProfileData map[string]interface{}
	exists := false

	if data, ok := result.DiagnosticsPortProfile[diagnosticsPortProfileName].(map[string]interface{}); ok {
		diagnosticsPortProfileData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found diagnostics port profile directly by ID: %s", diagnosticsPortProfileName))
	} else {
		for apiName, d := range result.DiagnosticsPortProfile {
			diagnosticsPortProfile, ok := d.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := diagnosticsPortProfile["name"].(string); ok && name == diagnosticsPortProfileName {
				diagnosticsPortProfileData = diagnosticsPortProfile
				diagnosticsPortProfileName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found diagnostics port profile with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Diagnostics Port Profile with ID '%s' not found in API response", diagnosticsPortProfileName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", diagnosticsPortProfileData["name"]))

	if enable, ok := diagnosticsPortProfileData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if enableSflow, ok := diagnosticsPortProfileData["enable_sflow"].(bool); ok {
		state.EnableSflow = types.BoolValue(enableSflow)
	} else {
		state.EnableSflow = types.BoolNull()
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

	if !plan.Name.Equal(state.Name) {
		diagnosticsPortProfileProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		diagnosticsPortProfileProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.EnableSflow.Equal(state.EnableSflow) {
		diagnosticsPortProfileProps.EnableSflow = openapi.PtrBool(plan.EnableSflow.ValueBool())
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "diagnostics_port_profile", name, diagnosticsPortProfileProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for diagnostics port profile update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Diagnostics Port Profile %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "diagnostics_port_profile", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for diagnostics port profile deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Diagnostics Port Profile %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Diagnostics Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnostics_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityDiagnosticsPortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
