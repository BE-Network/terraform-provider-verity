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
	_ resource.Resource                = &verityThresholdResource{}
	_ resource.ResourceWithConfigure   = &verityThresholdResource{}
	_ resource.ResourceWithImportState = &verityThresholdResource{}
)

func NewVerityThresholdResource() resource.Resource {
	return &verityThresholdResource{}
}

type verityThresholdResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityThresholdRulesModel struct {
	Enable           types.Bool   `tfsdk:"enable"`
	Type             types.String `tfsdk:"type"`
	Metric           types.String `tfsdk:"metric"`
	Operation        types.String `tfsdk:"operation"`
	Value            types.String `tfsdk:"value"`
	Threshold        types.String `tfsdk:"threshold"`
	ThresholdRefType types.String `tfsdk:"threshold_ref_type_"`
	Index            types.Int64  `tfsdk:"index"`
}

func (r verityThresholdRulesModel) GetIndex() types.Int64 {
	return r.Index
}

type verityThresholdResourceModel struct {
	Name          types.String                `tfsdk:"name"`
	Enable        types.Bool                  `tfsdk:"enable"`
	Type          types.String                `tfsdk:"type"`
	Operation     types.String                `tfsdk:"operation"`
	Severity      types.String                `tfsdk:"severity"`
	For           types.String                `tfsdk:"for"`
	KeepFiringFor types.String                `tfsdk:"keep_firing_for"`
	Rules         []verityThresholdRulesModel `tfsdk:"rules"`
}

func (r *verityThresholdResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_threshold"
}

func (r *verityThresholdResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityThresholdResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Threshold.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the threshold.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable or disable the threshold.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of elements threshold applies to.",
				Optional:    true,
			},
			"operation": schema.StringAttribute{
				Description: "How to combine rules.",
				Optional:    true,
			},
			"severity": schema.StringAttribute{
				Description: "Severity of the alarm when the threshold is met.",
				Optional:    true,
			},
			"for": schema.StringAttribute{
				Description: "Duration in minutes the threshold must be met before firing the alarm.",
				Optional:    true,
			},
			"keep_firing_for": schema.StringAttribute{
				Description: "Duration in minutes to keep firing the alarm after the threshold is no longer met.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"rules": schema.ListNestedBlock{
				Description: "Rules for the threshold.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable the rule.",
							Optional:    true,
						},
						"type": schema.StringAttribute{
							Description: "Use a metric or a nested threshold.",
							Optional:    true,
						},
						"metric": schema.StringAttribute{
							Description: "Metric threshold is on.",
							Optional:    true,
						},
						"operation": schema.StringAttribute{
							Description: "How to compare the metric to the value.",
							Optional:    true,
						},
						"value": schema.StringAttribute{
							Description: "Value to compare the metric to.",
							Optional:    true,
						},
						"threshold": schema.StringAttribute{
							Description: "Nested threshold reference.",
							Optional:    true,
						},
						"threshold_ref_type_": schema.StringAttribute{
							Description: "Object type for threshold field. Valid values: 'threshold'.",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index of the rule within the rules list.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityThresholdResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityThresholdResourceModel
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
	thresholdProps := &openapi.ThresholdsPutRequestThresholdValue{
		Name: openapi.PtrString(name),
	}

	// Set string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Type", APIField: &thresholdProps.Type, TFValue: plan.Type},
		{FieldName: "Operation", APIField: &thresholdProps.Operation, TFValue: plan.Operation},
		{FieldName: "Severity", APIField: &thresholdProps.Severity, TFValue: plan.Severity},
		{FieldName: "For", APIField: &thresholdProps.For, TFValue: plan.For},
		{FieldName: "KeepFiringFor", APIField: &thresholdProps.KeepFiringFor, TFValue: plan.KeepFiringFor},
	})

	// Set boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &thresholdProps.Enable, TFValue: plan.Enable},
	})

	// Handle rules
	if len(plan.Rules) > 0 {
		rules := make([]openapi.ThresholdsPutRequestThresholdValueRulesInner, len(plan.Rules))
		for i, ruleItem := range plan.Rules {
			rule := openapi.ThresholdsPutRequestThresholdValueRulesInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &rule.Enable, TFValue: ruleItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Type", APIField: &rule.Type, TFValue: ruleItem.Type},
				{FieldName: "Metric", APIField: &rule.Metric, TFValue: ruleItem.Metric},
				{FieldName: "Operation", APIField: &rule.Operation, TFValue: ruleItem.Operation},
				{FieldName: "Value", APIField: &rule.Value, TFValue: ruleItem.Value},
				{FieldName: "Threshold", APIField: &rule.Threshold, TFValue: ruleItem.Threshold},
				{FieldName: "ThresholdRefType", APIField: &rule.ThresholdRefType, TFValue: ruleItem.ThresholdRefType},
			})

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: ruleItem.Index},
			})

			rules[i] = rule
		}
		thresholdProps.Rules = rules
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "threshold", name, *thresholdProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Threshold %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "thresholds")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityThresholdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityThresholdResourceModel
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

	thresholdName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("threshold") {
		tflog.Info(ctx, fmt.Sprintf("Skipping threshold %s verification – trusting recent successful API operation", thresholdName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching thresholds for verification of %s", thresholdName))

	type ThresholdsResponse struct {
		Threshold map[string]interface{} `json:"threshold"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "thresholds", thresholdName,
		func() (ThresholdsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch thresholds")
			respAPI, err := r.client.ThresholdsAPI.ThresholdsGet(ctx).Execute()
			if err != nil {
				return ThresholdsResponse{}, fmt.Errorf("error reading thresholds: %v", err)
			}
			defer respAPI.Body.Close()

			var res ThresholdsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return ThresholdsResponse{}, fmt.Errorf("failed to decode thresholds response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d thresholds", len(res.Threshold)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Threshold %s", thresholdName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for threshold with name: %s", thresholdName))

	thresholdData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Threshold,
		thresholdName,
		func(data interface{}) (string, bool) {
			if threshold, ok := data.(map[string]interface{}); ok {
				if name, ok := threshold["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Threshold with name '%s' not found in API response", thresholdName))
		resp.State.RemoveResource(ctx)
		return
	}

	thresholdMap, ok := thresholdData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Threshold Data",
			fmt.Sprintf("Threshold data is not in expected format for %s", thresholdName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found threshold '%s' under API key '%s'", thresholdName, actualAPIName))

	state.Name = utils.MapStringFromAPI(thresholdMap["name"])

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"type":            &state.Type,
		"operation":       &state.Operation,
		"severity":        &state.Severity,
		"for":             &state.For,
		"keep_firing_for": &state.KeepFiringFor,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(thresholdMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(thresholdMap[apiKey])
	}

	// Handle rules
	if rules, ok := thresholdMap["rules"].([]interface{}); ok && len(rules) > 0 {
		var rulesList []verityThresholdRulesModel

		for _, r := range rules {
			rule, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			rModel := verityThresholdRulesModel{
				Enable:           utils.MapBoolFromAPI(rule["enable"]),
				Type:             utils.MapStringFromAPI(rule["type"]),
				Metric:           utils.MapStringFromAPI(rule["metric"]),
				Operation:        utils.MapStringFromAPI(rule["operation"]),
				Value:            utils.MapStringFromAPI(rule["value"]),
				Threshold:        utils.MapStringFromAPI(rule["threshold"]),
				ThresholdRefType: utils.MapStringFromAPI(rule["threshold_ref_type_"]),
				Index:            utils.MapInt64FromAPI(rule["index"]),
			}

			rulesList = append(rulesList, rModel)
		}

		state.Rules = rulesList
	} else {
		state.Rules = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityThresholdResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityThresholdResourceModel

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
	thresholdProps := openapi.ThresholdsPutRequestThresholdValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { thresholdProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Type, state.Type, func(v *string) { thresholdProps.Type = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Operation, state.Operation, func(v *string) { thresholdProps.Operation = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Severity, state.Severity, func(v *string) { thresholdProps.Severity = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.For, state.For, func(v *string) { thresholdProps.For = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.KeepFiringFor, state.KeepFiringFor, func(v *string) { thresholdProps.KeepFiringFor = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { thresholdProps.Enable = v }, &hasChanges)

	// Handle rules
	rulesHandler := utils.IndexedItemHandler[verityThresholdRulesModel, openapi.ThresholdsPutRequestThresholdValueRulesInner]{
		CreateNew: func(planItem verityThresholdRulesModel) openapi.ThresholdsPutRequestThresholdValueRulesInner {
			rule := openapi.ThresholdsPutRequestThresholdValueRulesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &rule.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Type", APIField: &rule.Type, TFValue: planItem.Type},
				{FieldName: "Metric", APIField: &rule.Metric, TFValue: planItem.Metric},
				{FieldName: "Operation", APIField: &rule.Operation, TFValue: planItem.Operation},
				{FieldName: "Value", APIField: &rule.Value, TFValue: planItem.Value},
				{FieldName: "Threshold", APIField: &rule.Threshold, TFValue: planItem.Threshold},
				{FieldName: "ThresholdRefType", APIField: &rule.ThresholdRefType, TFValue: planItem.ThresholdRefType},
			})

			return rule
		},
		UpdateExisting: func(planItem verityThresholdRulesModel, stateItem verityThresholdRulesModel) (openapi.ThresholdsPutRequestThresholdValueRulesInner, bool) {
			rule := openapi.ThresholdsPutRequestThresholdValueRulesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { rule.Enable = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.Type, stateItem.Type, func(v *string) { rule.Type = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Metric, stateItem.Metric, func(v *string) { rule.Metric = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Operation, stateItem.Operation, func(v *string) { rule.Operation = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Value, stateItem.Value, func(v *string) { rule.Value = v }, &fieldChanged)

			// Handle threshold and threshold_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Threshold, stateItem.Threshold, planItem.ThresholdRefType, stateItem.ThresholdRefType,
				func(v *string) { rule.Threshold = v },
				func(v *string) { rule.ThresholdRefType = v },
				"threshold", "threshold_ref_type_",
				&fieldChanged, &resp.Diagnostics,
			) {
				return rule, false
			}

			return rule, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.ThresholdsPutRequestThresholdValueRulesInner {
			rule := openapi.ThresholdsPutRequestThresholdValueRulesInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: types.Int64Value(index)},
			})
			return rule
		},
	}

	changedRules, rulesChanged := utils.ProcessIndexedArrayUpdates(plan.Rules, state.Rules, rulesHandler)
	if rulesChanged {
		thresholdProps.Rules = changedRules
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "threshold", name, thresholdProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Threshold %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "thresholds")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityThresholdResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityThresholdResourceModel
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "threshold", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Threshold %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "thresholds")
	resp.State.RemoveResource(ctx)
}

func (r *verityThresholdResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
