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
	_ resource.Resource                = &veritySpinePlaneResource{}
	_ resource.ResourceWithConfigure   = &veritySpinePlaneResource{}
	_ resource.ResourceWithImportState = &veritySpinePlaneResource{}
	_ resource.ResourceWithModifyPlan  = &veritySpinePlaneResource{}
)

const spinePlaneResourceType = "spineplanes"

func NewVeritySpinePlaneResource() resource.Resource {
	return &veritySpinePlaneResource{}
}

type veritySpinePlaneResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type veritySpinePlaneResourceModel struct {
	Name             types.String                            `tfsdk:"name"`
	Enable           types.Bool                              `tfsdk:"enable"`
	ObjectProperties []veritySpinePlaneObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySpinePlaneObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *veritySpinePlaneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spine_plane"
}

func (r *veritySpinePlaneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySpinePlaneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Spine Plane resource",
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
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the spine plane",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"notes": schema.StringAttribute{
							Description: "User Notes.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *veritySpinePlaneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan veritySpinePlaneResourceModel
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
	spinePlaneReq := &openapi.SpineplanesPutRequestSpinePlaneValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &spinePlaneReq.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		spinePlaneReq.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "spine_plane", name, *spinePlaneReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Spine Plane %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "spine_planes")

	var minState veritySpinePlaneResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if spinePlaneData, exists := bulkMgr.GetResourceResponse("spine_plane", name); exists {
			newState := populateSpinePlaneState(ctx, minState, spinePlaneData, r.provCtx.mode)
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

func (r *veritySpinePlaneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySpinePlaneResourceModel
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

	spinePlaneName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if spinePlaneData, exists := r.bulkOpsMgr.GetResourceResponse("spine_plane", spinePlaneName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached spine_plane data for %s from recent operation", spinePlaneName))
			state = populateSpinePlaneState(ctx, state, spinePlaneData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("spine_plane") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Spine Plane %s verification - trusting recent successful API operation", spinePlaneName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent Spine Plane operations found, performing normal verification for %s", spinePlaneName))

	type SpinePlanesResponse struct {
		SpinePlane map[string]map[string]interface{} `json:"spine_plane"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "spine_planes", spinePlaneName,
		func() (SpinePlanesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Spine Planes")
			respAPI, err := r.client.SpinePlanesAPI.SpineplanesGet(ctx).Execute()
			if err != nil {
				return SpinePlanesResponse{}, fmt.Errorf("error reading Spine Plane: %v", err)
			}
			defer respAPI.Body.Close()

			var res SpinePlanesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SpinePlanesResponse{}, fmt.Errorf("failed to decode Spine Planes response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Spine Planes from API", len(res.SpinePlane)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Spine Plane %s", spinePlaneName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Spine Plane with name: %s", spinePlaneName))

	spinePlaneData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.SpinePlane,
		spinePlaneName,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Spine Plane with name '%s' not found in API response", spinePlaneName))
		resp.State.RemoveResource(ctx)
		return
	}

	spinePlaneMap, ok := (interface{}(spinePlaneData)).(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Spine Plane Data",
			fmt.Sprintf("Spine Plane data is not in expected format for %s", spinePlaneName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Spine Plane '%s' under API key '%s'", spinePlaneName, actualAPIName))

	state = populateSpinePlaneState(ctx, state, spinePlaneMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySpinePlaneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySpinePlaneResourceModel

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
	spinePlaneReq := openapi.SpineplanesPutRequestSpinePlaneValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { spinePlaneReq.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { spinePlaneReq.Enable = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objProps.Notes},
		}, &objPropsChanged)

		if objPropsChanged {
			spinePlaneReq.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "spine_plane", name, spinePlaneReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Spine Plane %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "spine_planes")

	var minState veritySpinePlaneResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if spinePlaneData, exists := bulkMgr.GetResourceResponse("spine_plane", name); exists {
			newState := populateSpinePlaneState(ctx, minState, spinePlaneData, r.provCtx.mode)
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

func (r *veritySpinePlaneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state veritySpinePlaneResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "spine_plane", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Spine Plane %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "spine_planes")
	resp.State.RemoveResource(ctx)
}

func (r *veritySpinePlaneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateSpinePlaneState(ctx context.Context, state veritySpinePlaneResourceModel, data map[string]interface{}, mode string) veritySpinePlaneResourceModel {
	const resourceType = spinePlaneResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			state.ObjectProperties = []veritySpinePlaneObjectPropertiesModel{
				{
					Notes: utils.MapStringWithModeNested(objProps, "notes", resourceType, "object_properties.notes", mode),
				},
			}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *veritySpinePlaneResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySpinePlaneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = spinePlaneResourceType
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
}
