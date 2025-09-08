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
	bulkOpsMgr           *utils.BulkOperationManager
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

	if !plan.Enable.IsNull() {
		diagnosticsProfileProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.EnableSflow.IsNull() {
		diagnosticsProfileProps.EnableSflow = openapi.PtrBool(plan.EnableSflow.ValueBool())
	}

	if !plan.FlowCollector.IsNull() {
		diagnosticsProfileProps.FlowCollector = openapi.PtrString(plan.FlowCollector.ValueString())
	}

	if !plan.FlowCollectorRefType.IsNull() {
		diagnosticsProfileProps.FlowCollectorRefType = openapi.PtrString(plan.FlowCollectorRefType.ValueString())
	}

	if !plan.PollInterval.IsNull() {
		diagnosticsProfileProps.PollInterval = openapi.PtrInt32(int32(plan.PollInterval.ValueInt64()))
	}

	if !plan.VrfType.IsNull() {
		diagnosticsProfileProps.VrfType = openapi.PtrString(plan.VrfType.ValueString())
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "diagnostics_profile", name, *diagnosticsProfileProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for diagnostics profile creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Diagnostics Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Diagnostics Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnosticsprofiles")

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

	var result DiagnosticsProfilesResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		diagnosticsProfilesData, fetchErr := getCachedResponse(ctx, r.provCtx, "diagnosticsprofiles", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch diagnostics profiles")
			respAPI, err := r.client.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading diagnostics profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res DiagnosticsProfilesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode diagnostics profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d diagnostics profiles", len(res.DiagnosticsProfile)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch diagnostics profiles on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = diagnosticsProfilesData.(DiagnosticsProfilesResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Diagnostics Profile %s", diagnosticsProfileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for diagnostics profile with ID: %s", diagnosticsProfileName))
	var diagnosticsProfileData map[string]interface{}
	exists := false

	if data, ok := result.DiagnosticsProfile[diagnosticsProfileName].(map[string]interface{}); ok {
		diagnosticsProfileData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found diagnostics profile directly by ID: %s", diagnosticsProfileName))
	} else {
		for apiName, d := range result.DiagnosticsProfile {
			diagnosticsProfile, ok := d.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := diagnosticsProfile["name"].(string); ok && name == diagnosticsProfileName {
				diagnosticsProfileData = diagnosticsProfile
				diagnosticsProfileName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found diagnostics profile with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Diagnostics Profile with ID '%s' not found in API response", diagnosticsProfileName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", diagnosticsProfileData["name"]))

	if enable, ok := diagnosticsProfileData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if enableSflow, ok := diagnosticsProfileData["enable_sflow"].(bool); ok {
		state.EnableSflow = types.BoolValue(enableSflow)
	} else {
		state.EnableSflow = types.BoolNull()
	}

	stringAttrs := map[string]*types.String{
		"flow_collector":           &state.FlowCollector,
		"flow_collector_ref_type_": &state.FlowCollectorRefType,
		"vrf_type":                 &state.VrfType,
	}

	for apiKey, stateField := range stringAttrs {
		if value, ok := diagnosticsProfileData[apiKey].(string); ok {
			*stateField = types.StringValue(value)
		} else {
			*stateField = types.StringNull()
		}
	}

	if value, ok := diagnosticsProfileData["poll_interval"]; ok && value != nil {
		switch v := value.(type) {
		case int:
			state.PollInterval = types.Int64Value(int64(v))
		case int32:
			state.PollInterval = types.Int64Value(int64(v))
		case int64:
			state.PollInterval = types.Int64Value(v)
		case float64:
			state.PollInterval = types.Int64Value(int64(v))
		default:
			state.PollInterval = types.Int64Null()
		}
	} else {
		state.PollInterval = types.Int64Null()
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

	if !plan.Name.Equal(state.Name) {
		diagnosticsProfileProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		diagnosticsProfileProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.EnableSflow.Equal(state.EnableSflow) {
		diagnosticsProfileProps.EnableSflow = openapi.PtrBool(plan.EnableSflow.ValueBool())
		hasChanges = true
	}

	if !plan.VrfType.Equal(state.VrfType) {
		diagnosticsProfileProps.VrfType = openapi.PtrString(plan.VrfType.ValueString())
		hasChanges = true
	}

	if !plan.PollInterval.Equal(state.PollInterval) {
		diagnosticsProfileProps.PollInterval = openapi.PtrInt32(int32(plan.PollInterval.ValueInt64()))
		hasChanges = true
	}

	// Handle FlowCollector and FlowCollectorRefType according to "One ref type supported" rules
	flowCollectorChanged := !plan.FlowCollector.Equal(state.FlowCollector)
	flowCollectorRefTypeChanged := !plan.FlowCollectorRefType.Equal(state.FlowCollectorRefType)

	if flowCollectorChanged || flowCollectorRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.FlowCollector, plan.FlowCollectorRefType,
			"flow_collector", "flow_collector_ref_type_",
			flowCollectorChanged, flowCollectorRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if flowCollectorChanged && !flowCollectorRefTypeChanged {
			// Just send the base field
			if !plan.FlowCollector.IsNull() && plan.FlowCollector.ValueString() != "" {
				diagnosticsProfileProps.FlowCollector = openapi.PtrString(plan.FlowCollector.ValueString())
			} else {
				diagnosticsProfileProps.FlowCollector = openapi.PtrString("")
			}
			hasChanges = true
		} else if flowCollectorRefTypeChanged {
			// Send both fields
			if !plan.FlowCollector.IsNull() && plan.FlowCollector.ValueString() != "" {
				diagnosticsProfileProps.FlowCollector = openapi.PtrString(plan.FlowCollector.ValueString())
			} else {
				diagnosticsProfileProps.FlowCollector = openapi.PtrString("")
			}

			if !plan.FlowCollectorRefType.IsNull() && plan.FlowCollectorRefType.ValueString() != "" {
				diagnosticsProfileProps.FlowCollectorRefType = openapi.PtrString(plan.FlowCollectorRefType.ValueString())
			} else {
				diagnosticsProfileProps.FlowCollectorRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "diagnostics_profile", name, diagnosticsProfileProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for diagnostics profile update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Diagnostics Profile %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Diagnostics Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnosticsprofiles")
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "diagnostics_profile", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for diagnostics profile deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Diagnostics Profile %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Diagnostics Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "diagnosticsprofiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityDiagnosticsProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
