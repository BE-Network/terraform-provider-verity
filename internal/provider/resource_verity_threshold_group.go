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
	_ resource.Resource                = &verityThresholdGroupResource{}
	_ resource.ResourceWithConfigure   = &verityThresholdGroupResource{}
	_ resource.ResourceWithImportState = &verityThresholdGroupResource{}
	_ resource.ResourceWithModifyPlan  = &verityThresholdGroupResource{}
)

const thresholdGroupResourceType = "thresholdgroups"

func NewVerityThresholdGroupResource() resource.Resource {
	return &verityThresholdGroupResource{}
}

type verityThresholdGroupResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityThresholdGroupTargetsModel struct {
	Enable               types.Bool   `tfsdk:"enable"`
	Type                 types.String `tfsdk:"type"`
	GroupingRules        types.String `tfsdk:"grouping_rules"`
	GroupingRulesRefType types.String `tfsdk:"grouping_rules_ref_type_"`
	Switchpoint          types.String `tfsdk:"switchpoint"`
	SwitchpointRefType   types.String `tfsdk:"switchpoint_ref_type_"`
	Port                 types.String `tfsdk:"port"`
	Index                types.Int64  `tfsdk:"index"`
}

func (t verityThresholdGroupTargetsModel) GetIndex() types.Int64 {
	return t.Index
}

type verityThresholdGroupThresholdsModel struct {
	Enable           types.Bool   `tfsdk:"enable"`
	SeverityOverride types.String `tfsdk:"severity_override"`
	Threshold        types.String `tfsdk:"threshold"`
	ThresholdRefType types.String `tfsdk:"threshold_ref_type_"`
	Index            types.Int64  `tfsdk:"index"`
}

func (t verityThresholdGroupThresholdsModel) GetIndex() types.Int64 {
	return t.Index
}

type verityThresholdGroupResourceModel struct {
	Name       types.String                          `tfsdk:"name"`
	Enable     types.Bool                            `tfsdk:"enable"`
	Type       types.String                          `tfsdk:"type"`
	Targets    []verityThresholdGroupTargetsModel    `tfsdk:"targets"`
	Thresholds []verityThresholdGroupThresholdsModel `tfsdk:"thresholds"`
}

func (r *verityThresholdGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_threshold_group"
}

func (r *verityThresholdGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityThresholdGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Threshold Group.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the threshold group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable or disable the threshold group.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of elements to apply thresholds to. Valid values: 'interface', 'device'.",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"targets": schema.ListNestedBlock{
				Description: "Targets to apply thresholds to.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable the target.",
							Optional:    true,
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Specific element or Grouping Rules to apply thresholds to.",
							Optional:    true,
							Computed:    true,
						},
						"grouping_rules": schema.StringAttribute{
							Description: "Elements to apply thresholds to.",
							Optional:    true,
							Computed:    true,
						},
						"grouping_rules_ref_type_": schema.StringAttribute{
							Description: "Object type for grouping_rules field. Valid values: 'grouping_rules'.",
							Optional:    true,
							Computed:    true,
						},
						"switchpoint": schema.StringAttribute{
							Description: "Switchpoint to apply thresholds to.",
							Optional:    true,
							Computed:    true,
						},
						"switchpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for switchpoint field. Valid values: 'switchpoint'.",
							Optional:    true,
							Computed:    true,
						},
						"port": schema.StringAttribute{
							Description: "Port to apply thresholds to.",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index of the target within the targets list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"thresholds": schema.ListNestedBlock{
				Description: "Thresholds to apply to this group.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable the threshold.",
							Optional:    true,
							Computed:    true,
						},
						"severity_override": schema.StringAttribute{
							Description: "Override the severity defined in the threshold for this group only. Valid values: '', 'warning', 'notice', 'error', 'critical'.",
							Optional:    true,
							Computed:    true,
						},
						"threshold": schema.StringAttribute{
							Description: "Threshold to apply to this group.",
							Optional:    true,
							Computed:    true,
						},
						"threshold_ref_type_": schema.StringAttribute{
							Description: "Object type for threshold field. Valid values: 'threshold'.",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index of the threshold within the thresholds list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityThresholdGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityThresholdGroupResourceModel
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
	thresholdGroupProps := &openapi.ThresholdgroupsPutRequestThresholdGroupValue{
		Name: openapi.PtrString(name),
	}

	// Set string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Type", APIField: &thresholdGroupProps.Type, TFValue: plan.Type},
	})

	// Set boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &thresholdGroupProps.Enable, TFValue: plan.Enable},
	})

	// Handle targets
	if len(plan.Targets) > 0 {
		targets := make([]openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner, len(plan.Targets))
		for i, targetItem := range plan.Targets {
			target := openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &target.Enable, TFValue: targetItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Type", APIField: &target.Type, TFValue: targetItem.Type},
				{FieldName: "GroupingRules", APIField: &target.GroupingRules, TFValue: targetItem.GroupingRules},
				{FieldName: "GroupingRulesRefType", APIField: &target.GroupingRulesRefType, TFValue: targetItem.GroupingRulesRefType},
				{FieldName: "Switchpoint", APIField: &target.Switchpoint, TFValue: targetItem.Switchpoint},
				{FieldName: "SwitchpointRefType", APIField: &target.SwitchpointRefType, TFValue: targetItem.SwitchpointRefType},
				{FieldName: "Port", APIField: &target.Port, TFValue: targetItem.Port},
			})

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &target.Index, TFValue: targetItem.Index},
			})

			targets[i] = target
		}
		thresholdGroupProps.Targets = targets
	}

	// Handle thresholds
	if len(plan.Thresholds) > 0 {
		thresholds := make([]openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner, len(plan.Thresholds))
		for i, thresholdItem := range plan.Thresholds {
			threshold := openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &threshold.Enable, TFValue: thresholdItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "SeverityOverride", APIField: &threshold.SeverityOverride, TFValue: thresholdItem.SeverityOverride},
				{FieldName: "Threshold", APIField: &threshold.Threshold, TFValue: thresholdItem.Threshold},
				{FieldName: "ThresholdRefType", APIField: &threshold.ThresholdRefType, TFValue: thresholdItem.ThresholdRefType},
			})

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &threshold.Index, TFValue: thresholdItem.Index},
			})

			thresholds[i] = threshold
		}
		thresholdGroupProps.Thresholds = thresholds
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "threshold_group", name, *thresholdGroupProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Threshold group %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "threshold_groups")

	// Set minimal state first (just the identifier)
	var minState verityThresholdGroupResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if thresholdGroupData, exists := bulkMgr.GetResourceResponse("threshold_group", name); exists {
			newState := populateThresholdGroupState(ctx, minState, thresholdGroupData, r.provCtx.mode)
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

func (r *verityThresholdGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityThresholdGroupResourceModel
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

	thresholdGroupName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if thresholdGroupData, exists := r.bulkOpsMgr.GetResourceResponse("threshold_group", thresholdGroupName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached threshold_group data for %s from recent operation", thresholdGroupName))
			state = populateThresholdGroupState(ctx, state, thresholdGroupData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("threshold_group") {
		tflog.Info(ctx, fmt.Sprintf("Skipping threshold group %s verification – trusting recent successful API operation", thresholdGroupName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching threshold groups for verification of %s", thresholdGroupName))

	type ThresholdGroupResponse struct {
		ThresholdGroup map[string]interface{} `json:"threshold_group"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "threshold_groups", thresholdGroupName,
		func() (ThresholdGroupResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch threshold groups")
			respAPI, err := r.client.ThresholdGroupsAPI.ThresholdgroupsGet(ctx).Execute()
			if err != nil {
				return ThresholdGroupResponse{}, fmt.Errorf("error reading threshold groups: %v", err)
			}
			defer respAPI.Body.Close()

			var res ThresholdGroupResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return ThresholdGroupResponse{}, fmt.Errorf("failed to decode threshold groups response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d threshold groups", len(res.ThresholdGroup)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Threshold Group %s", thresholdGroupName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for threshold group with name: %s", thresholdGroupName))

	thresholdGroupData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.ThresholdGroup,
		thresholdGroupName,
		func(data interface{}) (string, bool) {
			if thresholdGroup, ok := data.(map[string]interface{}); ok {
				if name, ok := thresholdGroup["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Threshold group with name '%s' not found in API response", thresholdGroupName))
		resp.State.RemoveResource(ctx)
		return
	}

	thresholdGroupMap, ok := thresholdGroupData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Threshold Group Data",
			fmt.Sprintf("Threshold group data is not in expected format for %s", thresholdGroupName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found threshold group '%s' under API key '%s'", thresholdGroupName, actualAPIName))

	state = populateThresholdGroupState(ctx, state, thresholdGroupMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityThresholdGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityThresholdGroupResourceModel

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
	thresholdGroupProps := openapi.ThresholdgroupsPutRequestThresholdGroupValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { thresholdGroupProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Type, state.Type, func(v *string) { thresholdGroupProps.Type = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { thresholdGroupProps.Enable = v }, &hasChanges)

	// Handle targets
	targetsHandler := utils.IndexedItemHandler[verityThresholdGroupTargetsModel, openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner]{
		CreateNew: func(planItem verityThresholdGroupTargetsModel) openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner {
			target := openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &target.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &target.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Type", APIField: &target.Type, TFValue: planItem.Type},
				{FieldName: "GroupingRules", APIField: &target.GroupingRules, TFValue: planItem.GroupingRules},
				{FieldName: "GroupingRulesRefType", APIField: &target.GroupingRulesRefType, TFValue: planItem.GroupingRulesRefType},
				{FieldName: "Switchpoint", APIField: &target.Switchpoint, TFValue: planItem.Switchpoint},
				{FieldName: "SwitchpointRefType", APIField: &target.SwitchpointRefType, TFValue: planItem.SwitchpointRefType},
				{FieldName: "Port", APIField: &target.Port, TFValue: planItem.Port},
			})

			return target
		},
		UpdateExisting: func(planItem verityThresholdGroupTargetsModel, stateItem verityThresholdGroupTargetsModel) (openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner, bool) {
			target := openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &target.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { target.Enable = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.Type, stateItem.Type, func(v *string) { target.Type = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Port, stateItem.Port, func(v *string) { target.Port = v }, &fieldChanged)

			// Handle grouping_rules and grouping_rules_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.GroupingRules, stateItem.GroupingRules, planItem.GroupingRulesRefType, stateItem.GroupingRulesRefType,
				func(v *string) { target.GroupingRules = v },
				func(v *string) { target.GroupingRulesRefType = v },
				"grouping_rules", "grouping_rules_ref_type_",
				&fieldChanged, &resp.Diagnostics,
			) {
				return target, false
			}

			// Handle switchpoint and switchpoint_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Switchpoint, stateItem.Switchpoint, planItem.SwitchpointRefType, stateItem.SwitchpointRefType,
				func(v *string) { target.Switchpoint = v },
				func(v *string) { target.SwitchpointRefType = v },
				"switchpoint", "switchpoint_ref_type_",
				&fieldChanged, &resp.Diagnostics,
			) {
				return target, false
			}

			return target, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner {
			target := openapi.ThresholdgroupsPutRequestThresholdGroupValueTargetsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &target.Index, TFValue: types.Int64Value(index)},
			})
			return target
		},
	}

	changedTargets, targetsChanged := utils.ProcessIndexedArrayUpdates(plan.Targets, state.Targets, targetsHandler)
	if targetsChanged {
		thresholdGroupProps.Targets = changedTargets
		hasChanges = true
	}

	// Handle thresholds
	thresholdsHandler := utils.IndexedItemHandler[verityThresholdGroupThresholdsModel, openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner]{
		CreateNew: func(planItem verityThresholdGroupThresholdsModel) openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner {
			threshold := openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &threshold.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &threshold.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "SeverityOverride", APIField: &threshold.SeverityOverride, TFValue: planItem.SeverityOverride},
				{FieldName: "Threshold", APIField: &threshold.Threshold, TFValue: planItem.Threshold},
				{FieldName: "ThresholdRefType", APIField: &threshold.ThresholdRefType, TFValue: planItem.ThresholdRefType},
			})

			return threshold
		},
		UpdateExisting: func(planItem verityThresholdGroupThresholdsModel, stateItem verityThresholdGroupThresholdsModel) (openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner, bool) {
			threshold := openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &threshold.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { threshold.Enable = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.SeverityOverride, stateItem.SeverityOverride, func(v *string) { threshold.SeverityOverride = v }, &fieldChanged)

			// Handle threshold and threshold_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Threshold, stateItem.Threshold, planItem.ThresholdRefType, stateItem.ThresholdRefType,
				func(v *string) { threshold.Threshold = v },
				func(v *string) { threshold.ThresholdRefType = v },
				"threshold", "threshold_ref_type_",
				&fieldChanged, &resp.Diagnostics,
			) {
				return threshold, false
			}

			return threshold, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner {
			threshold := openapi.ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &threshold.Index, TFValue: types.Int64Value(index)},
			})
			return threshold
		},
	}

	changedThresholds, thresholdsChanged := utils.ProcessIndexedArrayUpdates(plan.Thresholds, state.Thresholds, thresholdsHandler)
	if thresholdsChanged {
		thresholdGroupProps.Thresholds = changedThresholds
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "threshold_group", name, thresholdGroupProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Threshold group %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "threshold_groups")

	var minState verityThresholdGroupResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if thresholdGroupData, exists := bulkMgr.GetResourceResponse("threshold_group", name); exists {
			newState := populateThresholdGroupState(ctx, minState, thresholdGroupData, r.provCtx.mode)
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

func (r *verityThresholdGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityThresholdGroupResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "threshold_group", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Threshold group %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "threshold_groups")
	resp.State.RemoveResource(ctx)
}

func (r *verityThresholdGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateThresholdGroupState(ctx context.Context, state verityThresholdGroupResourceModel, data map[string]interface{}, mode string) verityThresholdGroupResourceModel {
	const resourceType = thresholdGroupResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// String fields
	state.Type = utils.MapStringWithMode(data, "type", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Handle targets block
	if utils.FieldAppliesToMode(resourceType, "targets", mode) {
		if targets, ok := data["targets"].([]interface{}); ok && len(targets) > 0 {
			var targetsList []verityThresholdGroupTargetsModel

			for _, t := range targets {
				target, ok := t.(map[string]interface{})
				if !ok {
					continue
				}

				tModel := verityThresholdGroupTargetsModel{
					Enable:               utils.MapBoolWithModeNested(target, "enable", resourceType, "targets.enable", mode),
					Type:                 utils.MapStringWithModeNested(target, "type", resourceType, "targets.type", mode),
					GroupingRules:        utils.MapStringWithModeNested(target, "grouping_rules", resourceType, "targets.grouping_rules", mode),
					GroupingRulesRefType: utils.MapStringWithModeNested(target, "grouping_rules_ref_type_", resourceType, "targets.grouping_rules_ref_type_", mode),
					Switchpoint:          utils.MapStringWithModeNested(target, "switchpoint", resourceType, "targets.switchpoint", mode),
					SwitchpointRefType:   utils.MapStringWithModeNested(target, "switchpoint_ref_type_", resourceType, "targets.switchpoint_ref_type_", mode),
					Port:                 utils.MapStringWithModeNested(target, "port", resourceType, "targets.port", mode),
					Index:                utils.MapInt64WithModeNested(target, "index", resourceType, "targets.index", mode),
				}

				targetsList = append(targetsList, tModel)
			}

			state.Targets = targetsList
		} else {
			state.Targets = nil
		}
	} else {
		state.Targets = nil
	}

	// Handle thresholds block
	if utils.FieldAppliesToMode(resourceType, "thresholds", mode) {
		if thresholds, ok := data["thresholds"].([]interface{}); ok && len(thresholds) > 0 {
			var thresholdsList []verityThresholdGroupThresholdsModel

			for _, th := range thresholds {
				threshold, ok := th.(map[string]interface{})
				if !ok {
					continue
				}

				thModel := verityThresholdGroupThresholdsModel{
					Enable:           utils.MapBoolWithModeNested(threshold, "enable", resourceType, "thresholds.enable", mode),
					SeverityOverride: utils.MapStringWithModeNested(threshold, "severity_override", resourceType, "thresholds.severity_override", mode),
					Threshold:        utils.MapStringWithModeNested(threshold, "threshold", resourceType, "thresholds.threshold", mode),
					ThresholdRefType: utils.MapStringWithModeNested(threshold, "threshold_ref_type_", resourceType, "thresholds.threshold_ref_type_", mode),
					Index:            utils.MapInt64WithModeNested(threshold, "index", resourceType, "thresholds.index", mode),
				}

				thresholdsList = append(thresholdsList, thModel)
			}

			state.Thresholds = thresholdsList
		} else {
			state.Thresholds = nil
		}
	} else {
		state.Thresholds = nil
	}

	return state
}

func (r *verityThresholdGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityThresholdGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = thresholdGroupResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"type",
	)

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "targets",
		ItemCount:    len(plan.Targets),
		StringFields: []string{"type", "grouping_rules", "grouping_rules_ref_type_", "switchpoint", "switchpoint_ref_type_", "port"},
		BoolFields:   []string{"enable"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "thresholds",
		ItemCount:    len(plan.Thresholds),
		StringFields: []string{"severity_override", "threshold", "threshold_ref_type_"},
		BoolFields:   []string{"enable"},
		Int64Fields:  []string{"index"},
	})
}
