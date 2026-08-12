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
	_ resource.Resource                = &verityPlaneResource{}
	_ resource.ResourceWithConfigure   = &verityPlaneResource{}
	_ resource.ResourceWithImportState = &verityPlaneResource{}
	_ resource.ResourceWithModifyPlan  = &verityPlaneResource{}
)

const planeResourceType = "planes"
const planeTerraformType = "verity_plane"

func NewVerityPlaneResource() resource.Resource {
	return &verityPlaneResource{}
}

type verityPlaneResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityPlaneResourceModel struct {
	Name             types.String                       `tfsdk:"name"`
	Enable           types.Bool                         `tfsdk:"enable"`
	Site             types.String                       `tfsdk:"site"`
	SiteRefType      types.String                       `tfsdk:"site_ref_type_"`
	Position         types.Number                       `tfsdk:"position"`
	ObjectProperties []verityPlaneObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityPlaneObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityPlaneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plane"
}

func (r *verityPlaneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPlaneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Plane resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Template Name. Must be unique within type.",
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
			"site": schema.StringAttribute{
				Description: "Fabric this Plane is assigned to",
				Optional:    true,
				Computed:    true,
			},
			"site_ref_type_": schema.StringAttribute{
				Description: "Object type for site field",
				Optional:    true,
				Computed:    true,
			},
			"position": schema.NumberAttribute{
				Description: "Position of the Plane",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the Plane",
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

func (r *verityPlaneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPlaneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityPlaneResourceModel
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
	planeReq := &openapi.PlanesPutRequestPlaneValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Site", APIField: &planeReq.Site, TFValue: plan.Site},
		{FieldName: "SiteRefType", APIField: &planeReq.SiteRefType, TFValue: plan.SiteRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &planeReq.Enable, TFValue: plan.Enable},
	})

	// Handle nullable number fields - parse HCL to detect explicit config
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, planeTerraformType, name)

	utils.SetNullableNumberFields([]utils.NullableNumberFieldMapping{
		{FieldName: "Position", APIField: &planeReq.Position, TFValue: config.Position, IsConfigured: configuredAttrs.IsConfigured("position")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: plan.ObjectProperties[0].Notes, APIValue: &objProps.Notes},
		})
		planeReq.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "plane", name, *planeReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Plane %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "planes")

	var minState verityPlaneResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if planeData, exists := bulkMgr.GetResourceResponse("plane", name); exists {
			state := populatePlaneState(ctx, minState, utils.MergeMissingPlanScalars(planeData, plan, planeResourceType, r.provCtx.mode), r.provCtx.mode)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, planeResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityPlaneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPlaneResourceModel
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

	planeName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if planeData, exists := r.bulkOpsMgr.GetResourceResponse("plane", planeName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached plane data for %s from recent operation", planeName))
			state = populatePlaneState(ctx, state, utils.ApplyPostOperationFallback(ctx, planeData), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("plane") {
		tflog.Info(ctx, fmt.Sprintf("Skipping plane %s verification – trusting recent successful API operation", planeName))
		if handled, diags := utils.SetPostOperationFallbackState(ctx, &resp.State); handled {
			resp.Diagnostics.Append(diags...)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching planes for verification of %s", planeName))

	type PlanesResponse struct {
		Plane map[string]interface{} `json:"plane"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "planes", planeName,
		func() (PlanesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch planes")
			respAPI, err := r.client.PlanesAPI.PlanesGet(ctx).Execute()
			if err != nil {
				return PlanesResponse{}, fmt.Errorf("error reading planes: %v", err)
			}
			defer respAPI.Body.Close()

			var res PlanesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PlanesResponse{}, fmt.Errorf("failed to decode planes response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d planes", len(res.Plane)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Plane %s", planeName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for plane with name: %s", planeName))

	planeData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Plane,
		planeName,
		func(data interface{}) (string, bool) {
			if plane, ok := data.(map[string]interface{}); ok {
				if name, ok := plane["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Plane with name '%s' not found in API response", planeName))
		resp.State.RemoveResource(ctx)
		return
	}

	planeMap, ok := planeData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Plane Data",
			fmt.Sprintf("Plane data is not in expected format for %s", planeName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found plane '%s' under API key '%s'", planeName, actualAPIName))

	state = populatePlaneState(ctx, state, utils.ApplyPostOperationFallback(ctx, planeMap), r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPlaneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPlaneResourceModel

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
	var config verityPlaneResourceModel
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
	planeReq := openapi.PlanesPutRequestPlaneValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, planeTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { planeReq.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { planeReq.Enable = v }, &hasChanges)

	// Handle nullable number field changes
	utils.CompareAndSetNullableNumberField(config.Position, state.Position, configuredAttrs.IsConfigured("position"), func(v *openapi.NullableFloat64) { planeReq.Position = *v }, &hasChanges)

	// Handle Site and SiteRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Site, state.Site, plan.SiteRefType, state.SiteRefType,
		func(v *string) { planeReq.Site = v },
		func(v *string) { planeReq.SiteRefType = v },
		"site", "site_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		props := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		propsChanged := false
		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: plan.ObjectProperties[0].Notes, StateValue: state.ObjectProperties[0].Notes, APIValue: &props.Notes},
		}, &propsChanged)
		if propsChanged {
			planeReq.ObjectProperties = &props
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "plane", name, planeReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Plane %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "planes")

	var minState verityPlaneResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if planeData, exists := bulkMgr.GetResourceResponse("plane", name); exists {
			newState := populatePlaneState(ctx, minState, utils.MergeMissingPlanScalars(planeData, plan, planeResourceType, r.provCtx.mode), r.provCtx.mode)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, planeResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityPlaneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPlaneResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "plane", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Plane %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "planes")
	resp.State.RemoveResource(ctx)
}

func (r *verityPlaneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populatePlaneState(ctx context.Context, state verityPlaneResourceModel, data map[string]interface{}, mode string) verityPlaneResourceModel {
	const resourceType = planeResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Number fields
	state.Position = utils.MapNumberWithMode(data, "position", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// String fields
	state.Site = utils.MapStringWithMode(data, "site", resourceType, mode)
	state.SiteRefType = utils.MapStringWithMode(data, "site_ref_type_", resourceType, mode)

	// Handle object_properties list block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if props, ok := data["object_properties"].(map[string]interface{}); ok {
			state.ObjectProperties = []verityPlaneObjectPropertiesModel{
				{
					Notes: utils.MapStringWithModeNested(
						props,
						"notes",
						resourceType,
						"object_properties.notes",
						mode,
					),
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

func (r *verityPlaneResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityPlaneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = planeResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings("site", "site_ref_type_")
	nullifier.NullifyBools("enable")
	nullifier.NullifyNumbers("position")
	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"notes"},
	})

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityPlaneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityPlaneResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable number fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, planeTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		NumberFields: []utils.NullableNumberField{
			{AttrName: "position", ConfigVal: config.Position, StateVal: state.Position},
		},
	})
}
