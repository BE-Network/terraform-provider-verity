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
	_ resource.Resource                = &veritySuResource{}
	_ resource.ResourceWithConfigure   = &veritySuResource{}
	_ resource.ResourceWithImportState = &veritySuResource{}
	_ resource.ResourceWithModifyPlan  = &veritySuResource{}
)

const suResourceType = "sus"
const suTerraformType = "verity_su"

func NewVeritySuResource() resource.Resource {
	return &veritySuResource{}
}

type veritySuResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type veritySuResourceModel struct {
	Name             types.String                    `tfsdk:"name"`
	Enable           types.Bool                      `tfsdk:"enable"`
	Pod              types.String                    `tfsdk:"pod"`
	PodRefType       types.String                    `tfsdk:"pod_ref_type_"`
	Position         types.Number                    `tfsdk:"position"`
	ObjectProperties []veritySuObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySuObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *veritySuResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_su"
}

func (r *veritySuResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySuResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity SU resource",
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
			"pod": schema.StringAttribute{
				Description: "Pod this SU is assigned to",
				Optional:    true,
				Computed:    true,
			},
			"pod_ref_type_": schema.StringAttribute{
				Description: "Object type for pod field",
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
				Description: "Object properties for the SU",
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

func (r *veritySuResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan veritySuResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySuResourceModel
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
	suReq := &openapi.SusPutRequestSuValue{
		Name: openapi.PtrString(name),
	}

	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Pod", APIField: &suReq.Pod, TFValue: plan.Pod},
		{FieldName: "PodRefType", APIField: &suReq.PodRefType, TFValue: plan.PodRefType},
	})

	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &suReq.Enable, TFValue: plan.Enable},
	})

	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, suTerraformType, name)

	utils.SetNullableNumberFields([]utils.NullableNumberFieldMapping{
		{FieldName: "Position", APIField: &suReq.Position, TFValue: config.Position, IsConfigured: configuredAttrs.IsConfigured("position")},
	})

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		suReq.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "su", name, *suReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SU %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sus")

	var minState veritySuResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if suData, exists := bulkMgr.GetResourceResponse("su", name); exists {
			newState := populateSuState(ctx, minState, suData, r.provCtx.mode)
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

	r.Read(ctx, readReq, &readResp)
}

func (r *veritySuResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySuResourceModel
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

	suName := state.Name.ValueString()

	if r.bulkOpsMgr != nil {
		if suData, exists := r.bulkOpsMgr.GetResourceResponse("su", suName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached su data for %s from recent operation", suName))
			state = populateSuState(ctx, state, suData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("su") {
		tflog.Info(ctx, fmt.Sprintf("Skipping SU %s verification - trusting recent successful API operation", suName))
		return
	}

	type SusResponse struct {
		Su map[string]map[string]interface{} `json:"su"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "sus", suName,
		func() (SusResponse, error) {
			respAPI, err := r.client.SUsAPI.SusGet(ctx).Execute()
			if err != nil {
				return SusResponse{}, fmt.Errorf("error reading SU: %v", err)
			}
			defer respAPI.Body.Close()

			var res SusResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SusResponse{}, fmt.Errorf("failed to decode SUs response: %v", err)
			}

			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read SU %s", suName))...,
		)
		return
	}

	suData, _, exists := utils.FindResourceByAPIName(
		result.Su,
		suName,
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

	suMap, ok := (interface{}(suData)).(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid SU Data",
			fmt.Sprintf("SU data is not in expected format for %s", suName),
		)
		return
	}

	state = populateSuState(ctx, state, suMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySuResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySuResourceModel

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

	var config veritySuResourceModel
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
	suReq := openapi.SusPutRequestSuValue{}
	hasChanges := false

	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, suTerraformType, name)

	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { suReq.Name = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { suReq.Enable = v }, &hasChanges)
	utils.CompareAndSetNullableNumberField(config.Position, state.Position, configuredAttrs.IsConfigured("position"), func(v *openapi.NullableFloat64) { suReq.Position = *v }, &hasChanges)

	if !utils.HandleOneRefTypeSupported(
		plan.Pod, state.Pod, plan.PodRefType, state.PodRefType,
		func(v *string) { suReq.Pod = v },
		func(v *string) { suReq.PodRefType = v },
		"pod", "pod_ref_type_",
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
			suReq.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "su", name, suReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SU %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sus")

	var minState veritySuResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if suData, exists := bulkMgr.GetResourceResponse("su", name); exists {
			newState := populateSuState(ctx, minState, suData, r.provCtx.mode)
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

	r.Read(ctx, readReq, &readResp)
}

func (r *veritySuResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state veritySuResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "su", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SU %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sus")
	resp.State.RemoveResource(ctx)
}

func (r *veritySuResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateSuState(ctx context.Context, state veritySuResourceModel, data map[string]interface{}, mode string) veritySuResourceModel {
	const resourceType = suResourceType

	state.Name = utils.MapStringFromAPI(data["name"])
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.Pod = utils.MapStringWithMode(data, "pod", resourceType, mode)
	state.PodRefType = utils.MapStringWithMode(data, "pod_ref_type_", resourceType, mode)
	state.Position = utils.MapNumberWithMode(data, "position", resourceType, mode)

	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			state.ObjectProperties = []veritySuObjectPropertiesModel{
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

func (r *veritySuResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySuResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const resourceType = suResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"pod", "pod_ref_type_",
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

	var state veritySuResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySuResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, suTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		NumberFields: []utils.NullableNumberField{
			{AttrName: "position", ConfigVal: config.Position, StateVal: state.Position},
		},
	})
}
