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
	_ resource.Resource                = &veritySspGroupResource{}
	_ resource.ResourceWithConfigure   = &veritySspGroupResource{}
	_ resource.ResourceWithImportState = &veritySspGroupResource{}
	_ resource.ResourceWithModifyPlan  = &veritySspGroupResource{}
)

const sspGroupResourceType = "sspgroups"
const sspGroupTerraformType = "verity_ssp_group"

func NewVeritySspGroupResource() resource.Resource {
	return &veritySspGroupResource{}
}

type veritySspGroupResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type veritySspGroupResourceModel struct {
	Name             types.String                          `tfsdk:"name"`
	Enable           types.Bool                            `tfsdk:"enable"`
	Site             types.String                          `tfsdk:"site"`
	SiteRefType      types.String                          `tfsdk:"site_ref_type_"`
	Position         types.Number                          `tfsdk:"position"`
	ObjectProperties []veritySspGroupObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySspGroupObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *veritySspGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssp_group"
}

func (r *veritySspGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySspGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity SuperSpine Group resource",
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
				Description: "Fabric this SuperSpine Group is assigned to",
				Optional:    true,
				Computed:    true,
			},
			"site_ref_type_": schema.StringAttribute{
				Description: "Object type for site field",
				Optional:    true,
				Computed:    true,
			},
			"position": schema.NumberAttribute{
				Description: "Position of the Switch",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the superspine group",
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

func (r *veritySspGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan veritySspGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySspGroupResourceModel
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
	sspGroupReq := &openapi.SspgroupsPutRequestSuperspineGroupValue{
		Name: openapi.PtrString(name),
	}

	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Site", APIField: &sspGroupReq.Site, TFValue: plan.Site},
		{FieldName: "SiteRefType", APIField: &sspGroupReq.SiteRefType, TFValue: plan.SiteRefType},
	})

	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &sspGroupReq.Enable, TFValue: plan.Enable},
	})

	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, sspGroupTerraformType, name)

	utils.SetNullableNumberFields([]utils.NullableNumberFieldMapping{
		{FieldName: "Position", APIField: &sspGroupReq.Position, TFValue: config.Position, IsConfigured: configuredAttrs.IsConfigured("position")},
	})

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		sspGroupReq.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "ssp_group", name, *sspGroupReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SSP Group %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ssp_groups")

	var minState veritySspGroupResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if sspGroupData, exists := bulkMgr.GetResourceResponse("ssp_group", name); exists {
			newState := populateSspGroupState(ctx, minState, utils.MergeMissingPlanScalars(sspGroupData, plan, sspGroupResourceType, r.provCtx.mode), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
			return
		}
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, sspGroupResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *veritySspGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySspGroupResourceModel
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

	sspGroupName := state.Name.ValueString()

	if r.bulkOpsMgr != nil {
		if sspGroupData, exists := r.bulkOpsMgr.GetResourceResponse("ssp_group", sspGroupName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached ssp_group data for %s from recent operation", sspGroupName))
			state = populateSspGroupState(ctx, state, utils.ApplyPostOperationFallback(ctx, sspGroupData), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ssp_group") {
		tflog.Info(ctx, fmt.Sprintf("Skipping SSP Group %s verification - trusting recent successful API operation", sspGroupName))
		if handled, diags := utils.SetPostOperationFallbackState(ctx, &resp.State); handled {
			resp.Diagnostics.Append(diags...)
		}
		return
	}

	type SspGroupsResponse struct {
		SuperspineGroup map[string]map[string]interface{} `json:"superspine_group"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "ssp_groups", sspGroupName,
		func() (SspGroupsResponse, error) {
			respAPI, err := r.client.SuperSpineGroupsAPI.SspgroupsGet(ctx).Execute()
			if err != nil {
				return SspGroupsResponse{}, fmt.Errorf("error reading SSP Group: %v", err)
			}
			defer respAPI.Body.Close()

			var res SspGroupsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SspGroupsResponse{}, fmt.Errorf("failed to decode SSP Groups response: %v", err)
			}

			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read SSP Group %s", sspGroupName))...,
		)
		return
	}

	sspGroupData, _, exists := utils.FindResourceByAPIName(
		result.SuperspineGroup,
		sspGroupName,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	sspGroupMap, ok := (interface{}(sspGroupData)).(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid SSP Group Data",
			fmt.Sprintf("SSP Group data is not in expected format for %s", sspGroupName),
		)
		return
	}

	state = populateSspGroupState(ctx, state, utils.ApplyPostOperationFallback(ctx, sspGroupMap), r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySspGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySspGroupResourceModel

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

	var config veritySspGroupResourceModel
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
	sspGroupReq := openapi.SspgroupsPutRequestSuperspineGroupValue{}
	hasChanges := false

	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, sspGroupTerraformType, name)

	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { sspGroupReq.Name = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { sspGroupReq.Enable = v }, &hasChanges)
	utils.CompareAndSetNullableNumberField(config.Position, state.Position, configuredAttrs.IsConfigured("position"), func(v *openapi.NullableFloat64) { sspGroupReq.Position = *v }, &hasChanges)

	if !utils.HandleOneRefTypeSupported(
		plan.Site, state.Site, plan.SiteRefType, state.SiteRefType,
		func(v *string) { sspGroupReq.Site = v },
		func(v *string) { sspGroupReq.SiteRefType = v },
		"site", "site_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objProps.Notes},
		}, &objPropsChanged)

		if objPropsChanged {
			sspGroupReq.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "ssp_group", name, sspGroupReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SSP Group %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ssp_groups")

	var minState veritySspGroupResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if sspGroupData, exists := bulkMgr.GetResourceResponse("ssp_group", name); exists {
			newState := populateSspGroupState(ctx, minState, utils.MergeMissingPlanScalars(sspGroupData, plan, sspGroupResourceType, r.provCtx.mode), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
			return
		}
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, sspGroupResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *veritySspGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state veritySspGroupResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "ssp_group", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SSP Group %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ssp_groups")
	resp.State.RemoveResource(ctx)
}

func (r *veritySspGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateSspGroupState(ctx context.Context, state veritySspGroupResourceModel, data map[string]interface{}, mode string) veritySspGroupResourceModel {
	const resourceType = sspGroupResourceType

	state.Name = utils.MapStringFromAPI(data["name"])
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.Site = utils.MapStringWithMode(data, "site", resourceType, mode)
	state.SiteRefType = utils.MapStringWithMode(data, "site_ref_type_", resourceType, mode)
	state.Position = utils.MapNumberWithMode(data, "position", resourceType, mode)

	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			state.ObjectProperties = []veritySspGroupObjectPropertiesModel{
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

func (r *veritySspGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySspGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const resourceType = sspGroupResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"site", "site_ref_type_",
	)

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyNumbers(
		"position",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"notes"},
	})

	if req.State.Raw.IsNull() {
		return
	}

	var state veritySspGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySspGroupResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, sspGroupTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		NumberFields: []utils.NullableNumberField{
			{AttrName: "position", ConfigVal: config.Position, StateVal: state.Position},
		},
	})
}
