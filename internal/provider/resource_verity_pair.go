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
	_ resource.Resource                = &verityPairResource{}
	_ resource.ResourceWithConfigure   = &verityPairResource{}
	_ resource.ResourceWithImportState = &verityPairResource{}
	_ resource.ResourceWithModifyPlan  = &verityPairResource{}
)

const pairResourceType = "pairs"

func NewVerityPairResource() resource.Resource {
	return &verityPairResource{}
}

type verityPairResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityPairResourceModel struct {
	Name                types.String `tfsdk:"name"`
	Enable              types.Bool   `tfsdk:"enable"`
	Switchpoint1        types.String `tfsdk:"switchpoint_1"`
	Switchpoint1RefType types.String `tfsdk:"switchpoint_1_ref_type_"`
	Switchpoint2        types.String `tfsdk:"switchpoint_2"`
	Switchpoint2RefType types.String `tfsdk:"switchpoint_2_ref_type_"`
	Lag                 types.String `tfsdk:"lag"`
	LagRefType          types.String `tfsdk:"lag_ref_type_"`
	IsWhiteboxPair      types.Bool   `tfsdk:"is_whitebox_pair"`
}

func (r *verityPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pair"
}

func (r *verityPairResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPairResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Switch Pair.",
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
			"switchpoint_1": schema.StringAttribute{
				Description: "Switchpoint",
				Optional:    true,
				Computed:    true,
			},
			"switchpoint_1_ref_type_": schema.StringAttribute{
				Description: "Object type for switchpoint_1 field",
				Optional:    true,
				Computed:    true,
			},
			"switchpoint_2": schema.StringAttribute{
				Description: "Switchpoint",
				Optional:    true,
				Computed:    true,
			},
			"switchpoint_2_ref_type_": schema.StringAttribute{
				Description: "Object type for switchpoint_2 field",
				Optional:    true,
				Computed:    true,
			},
			"lag": schema.StringAttribute{
				Description: "LAG",
				Optional:    true,
				Computed:    true,
			},
			"lag_ref_type_": schema.StringAttribute{
				Description: "Object type for lag field",
				Optional:    true,
				Computed:    true,
			},
			"is_whitebox_pair": schema.BoolAttribute{
				Description: "Is Whitebox Pair",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *verityPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPairResourceModel
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
	pairProps := &openapi.PairsPutRequestSwitchPairValue{
		Name: openapi.PtrString(name),
	}

	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Switchpoint1", APIField: &pairProps.Switchpoint1, TFValue: plan.Switchpoint1},
		{FieldName: "Switchpoint1RefType", APIField: &pairProps.Switchpoint1RefType, TFValue: plan.Switchpoint1RefType},
		{FieldName: "Switchpoint2", APIField: &pairProps.Switchpoint2, TFValue: plan.Switchpoint2},
		{FieldName: "Switchpoint2RefType", APIField: &pairProps.Switchpoint2RefType, TFValue: plan.Switchpoint2RefType},
		{FieldName: "Lag", APIField: &pairProps.Lag, TFValue: plan.Lag},
		{FieldName: "LagRefType", APIField: &pairProps.LagRefType, TFValue: plan.LagRefType},
	})

	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &pairProps.Enable, TFValue: plan.Enable},
		{FieldName: "IsWhiteboxPair", APIField: &pairProps.IsWhiteboxPair, TFValue: plan.IsWhiteboxPair},
	})

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "pair", name, *pairProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Pair %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pairs")

	var minState verityPairResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pairData, exists := bulkMgr.GetResourceResponse("pair", name); exists {
			state := populatePairState(ctx, minState, pairData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

func (r *verityPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPairResourceModel
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

	pairName := state.Name.ValueString()

	if r.bulkOpsMgr != nil {
		if pairData, exists := r.bulkOpsMgr.GetResourceResponse("pair", pairName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached pair data for %s from recent operation", pairName))
			state = populatePairState(ctx, state, pairData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("pair") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Pair %s verification – trusting recent successful API operation", pairName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Pairs for verification of %s", pairName))

	type PairsResponse struct {
		SwitchPair map[string]interface{} `json:"switch_pair"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "pairs", pairName,
		func() (PairsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Pairs")
			respAPI, err := r.client.SwitchPairsAPI.PairsGet(ctx).Execute()
			if err != nil {
				return PairsResponse{}, fmt.Errorf("error reading Pairs: %v", err)
			}
			defer respAPI.Body.Close()

			var res PairsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PairsResponse{}, fmt.Errorf("failed to decode Pairs response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Pairs", len(res.SwitchPair)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Pair %s", pairName))...,
		)
		return
	}

	pairData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.SwitchPair,
		pairName,
		func(data interface{}) (string, bool) {
			if pair, ok := data.(map[string]interface{}); ok {
				if name, ok := pair["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Pair with name '%s' not found in API response", pairName))
		resp.State.RemoveResource(ctx)
		return
	}

	pairMap, ok := pairData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Pair Data",
			fmt.Sprintf("Pair data is not in expected format for %s", pairName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Pair '%s' under API key '%s'", pairName, actualAPIName))

	state = populatePairState(ctx, state, pairMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPairResourceModel

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
	pairProps := openapi.PairsPutRequestSwitchPairValue{}
	hasChanges := false

	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { pairProps.Name = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { pairProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IsWhiteboxPair, state.IsWhiteboxPair, func(v *bool) { pairProps.IsWhiteboxPair = v }, &hasChanges)

	if !utils.HandleOneRefTypeSupported(
		plan.Switchpoint1, state.Switchpoint1, plan.Switchpoint1RefType, state.Switchpoint1RefType,
		func(v *string) { pairProps.Switchpoint1 = v },
		func(v *string) { pairProps.Switchpoint1RefType = v },
		"switchpoint_1", "switchpoint_1_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.Switchpoint2, state.Switchpoint2, plan.Switchpoint2RefType, state.Switchpoint2RefType,
		func(v *string) { pairProps.Switchpoint2 = v },
		func(v *string) { pairProps.Switchpoint2RefType = v },
		"switchpoint_2", "switchpoint_2_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.Lag, state.Lag, plan.LagRefType, state.LagRefType,
		func(v *string) { pairProps.Lag = v },
		func(v *string) { pairProps.LagRefType = v },
		"lag", "lag_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "pair", name, pairProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Pair %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pairs")

	var minState verityPairResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pairData, exists := bulkMgr.GetResourceResponse("pair", name); exists {
			newState := populatePairState(ctx, minState, pairData, r.provCtx.mode)
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

func (r *verityPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPairResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "pair", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Pair %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pairs")
	resp.State.RemoveResource(ctx)
}

func (r *verityPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populatePairState(ctx context.Context, state verityPairResourceModel, data map[string]interface{}, mode string) verityPairResourceModel {
	const resourceType = pairResourceType

	state.Name = utils.MapStringFromAPI(data["name"])
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.IsWhiteboxPair = utils.MapBoolWithMode(data, "is_whitebox_pair", resourceType, mode)
	state.Switchpoint1 = utils.MapStringWithMode(data, "switchpoint_1", resourceType, mode)
	state.Switchpoint1RefType = utils.MapStringWithMode(data, "switchpoint_1_ref_type_", resourceType, mode)
	state.Switchpoint2 = utils.MapStringWithMode(data, "switchpoint_2", resourceType, mode)
	state.Switchpoint2RefType = utils.MapStringWithMode(data, "switchpoint_2_ref_type_", resourceType, mode)
	state.Lag = utils.MapStringWithMode(data, "lag", resourceType, mode)
	state.LagRefType = utils.MapStringWithMode(data, "lag_ref_type_", resourceType, mode)

	return state
}

func (r *verityPairResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityPairResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const resourceType = pairResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"switchpoint_1", "switchpoint_1_ref_type_",
		"switchpoint_2", "switchpoint_2_ref_type_",
		"lag", "lag_ref_type_",
	)

	nullifier.NullifyBools(
		"enable", "is_whitebox_pair",
	)
}
