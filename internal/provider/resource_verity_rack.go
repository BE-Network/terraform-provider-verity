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
	_ resource.Resource                = &verityRackResource{}
	_ resource.ResourceWithConfigure   = &verityRackResource{}
	_ resource.ResourceWithImportState = &verityRackResource{}
	_ resource.ResourceWithModifyPlan  = &verityRackResource{}
)

const rackResourceType = "racks"
const rackTerraformType = "verity_rack"

func NewVerityRackResource() resource.Resource {
	return &verityRackResource{}
}

type verityRackResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityRackResourceModel struct {
	Name             types.String                      `tfsdk:"name"`
	Enable           types.Bool                        `tfsdk:"enable"`
	Position         types.Number                      `tfsdk:"position"`
	Su               types.String                      `tfsdk:"su"`
	SuRefType        types.String                      `tfsdk:"su_ref_type_"`
	ObjectProperties []verityRackObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityRackObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityRackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack"
}

func (r *verityRackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityRackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Rack resource",
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
			"position": schema.NumberAttribute{
				Description: "Position of the Rack",
				Optional:    true,
				Computed:    true,
			},
			"su": schema.StringAttribute{
				Description: "SU this Rack is assigned to",
				Optional:    true,
				Computed:    true,
			},
			"su_ref_type_": schema.StringAttribute{
				Description: "Object type for su field",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the Rack",
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

func (r *verityRackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityRackResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityRackResourceModel
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
	rackReq := &openapi.RacksPutRequestRackValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Su", APIField: &rackReq.Su, TFValue: plan.Su},
		{FieldName: "SuRefType", APIField: &rackReq.SuRefType, TFValue: plan.SuRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &rackReq.Enable, TFValue: plan.Enable},
	})

	// Handle nullable number fields - parse HCL to detect explicit config
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, rackTerraformType, name)

	utils.SetNullableNumberFields([]utils.NullableNumberFieldMapping{
		{FieldName: "Position", APIField: &rackReq.Position, TFValue: config.Position, IsConfigured: configuredAttrs.IsConfigured("position")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: plan.ObjectProperties[0].Notes, APIValue: &objProps.Notes},
		})
		rackReq.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "rack", name, *rackReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Rack %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "racks")

	var minState verityRackResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if rackData, exists := bulkMgr.GetResourceResponse("rack", name); exists {
			state := populateRackState(ctx, minState, utils.MergeMissingPlanScalars(rackData, plan, rackResourceType, r.provCtx.mode), r.provCtx.mode)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, rackResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityRackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityRackResourceModel
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

	rackName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if rackData, exists := r.bulkOpsMgr.GetResourceResponse("rack", rackName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached rack data for %s from recent operation", rackName))
			state = populateRackState(ctx, state, utils.ApplyPostOperationFallback(ctx, rackData), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("rack") {
		tflog.Info(ctx, fmt.Sprintf("Skipping rack %s verification – trusting recent successful API operation", rackName))
		if handled, diags := utils.SetPostOperationFallbackState(ctx, &resp.State); handled {
			resp.Diagnostics.Append(diags...)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching racks for verification of %s", rackName))

	type RacksResponse struct {
		Rack map[string]interface{} `json:"rack"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "racks", rackName,
		func() (RacksResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch racks")
			respAPI, err := r.client.RacksAPI.RacksGet(ctx).Execute()
			if err != nil {
				return RacksResponse{}, fmt.Errorf("error reading racks: %v", err)
			}
			defer respAPI.Body.Close()

			var res RacksResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return RacksResponse{}, fmt.Errorf("failed to decode racks response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d racks", len(res.Rack)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Rack %s", rackName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for rack with name: %s", rackName))

	rackData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Rack,
		rackName,
		func(data interface{}) (string, bool) {
			if rack, ok := data.(map[string]interface{}); ok {
				if name, ok := rack["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Rack with name '%s' not found in API response", rackName))
		resp.State.RemoveResource(ctx)
		return
	}

	rackMap, ok := rackData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Rack Data",
			fmt.Sprintf("Rack data is not in expected format for %s", rackName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found rack '%s' under API key '%s'", rackName, actualAPIName))

	state = populateRackState(ctx, state, utils.ApplyPostOperationFallback(ctx, rackMap), r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityRackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityRackResourceModel

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
	var config verityRackResourceModel
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
	rackReq := openapi.RacksPutRequestRackValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, rackTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { rackReq.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { rackReq.Enable = v }, &hasChanges)

	// Handle nullable number field changes
	utils.CompareAndSetNullableNumberField(config.Position, state.Position, configuredAttrs.IsConfigured("position"), func(v *openapi.NullableFloat64) { rackReq.Position = *v }, &hasChanges)

	// Handle Su and SuRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Su, state.Su, plan.SuRefType, state.SuRefType,
		func(v *string) { rackReq.Su = v },
		func(v *string) { rackReq.SuRefType = v },
		"su", "su_ref_type_",
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
			rackReq.ObjectProperties = &props
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "rack", name, rackReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Rack %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "racks")

	var minState verityRackResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if rackData, exists := bulkMgr.GetResourceResponse("rack", name); exists {
			newState := populateRackState(ctx, minState, utils.MergeMissingPlanScalars(rackData, plan, rackResourceType, r.provCtx.mode), r.provCtx.mode)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, rackResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityRackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityRackResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "rack", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Rack %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "racks")
	resp.State.RemoveResource(ctx)
}

func (r *verityRackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateRackState(ctx context.Context, state verityRackResourceModel, data map[string]interface{}, mode string) verityRackResourceModel {
	const resourceType = rackResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Number fields
	state.Position = utils.MapNumberWithMode(data, "position", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// String fields
	state.Su = utils.MapStringWithMode(data, "su", resourceType, mode)
	state.SuRefType = utils.MapStringWithMode(data, "su_ref_type_", resourceType, mode)

	// Handle object_properties list block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if props, ok := data["object_properties"].(map[string]interface{}); ok {
			state.ObjectProperties = []verityRackObjectPropertiesModel{
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

func (r *verityRackResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityRackResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = rackResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings("su", "su_ref_type_")
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
	var state verityRackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityRackResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, rackTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		NumberFields: []utils.NullableNumberField{
			{AttrName: "position", ConfigVal: config.Position, StateVal: state.Position},
		},
	})
}
