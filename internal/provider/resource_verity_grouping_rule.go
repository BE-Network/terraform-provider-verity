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
	_ resource.Resource                = &verityGroupingRuleResource{}
	_ resource.ResourceWithConfigure   = &verityGroupingRuleResource{}
	_ resource.ResourceWithImportState = &verityGroupingRuleResource{}
)

func NewVerityGroupingRuleResource() resource.Resource {
	return &verityGroupingRuleResource{}
}

type verityGroupingRuleResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityGroupingRuleRulesModel struct {
	Enable               types.Bool   `tfsdk:"enable"`
	RuleInvert           types.Bool   `tfsdk:"rule_invert"`
	RuleType             types.String `tfsdk:"rule_type"`
	RuleValue            types.String `tfsdk:"rule_value"`
	RuleValuePath        types.String `tfsdk:"rule_value_path"`
	RuleValuePathRefType types.String `tfsdk:"rule_value_path_ref_type_"`
	Index                types.Int64  `tfsdk:"index"`
}

func (r verityGroupingRuleRulesModel) GetIndex() types.Int64 {
	return r.Index
}

type verityGroupingRuleResourceModel struct {
	Name      types.String                   `tfsdk:"name"`
	Enable    types.Bool                     `tfsdk:"enable"`
	Type      types.String                   `tfsdk:"type"`
	Operation types.String                   `tfsdk:"operation"`
	Rules     []verityGroupingRuleRulesModel `tfsdk:"rules"`
}

func (r *verityGroupingRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grouping_rule"
}

func (r *verityGroupingRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityGroupingRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Grouping Rule.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the grouping rule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable or disable the grouping rule.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the grouping rule. Valid values: 'and', 'or'.",
				Optional:    true,
			},
			"operation": schema.StringAttribute{
				Description: "The operation of the grouping rule. Valid values: 'permit', 'deny'.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"rules": schema.ListNestedBlock{
				Description: "List of rules within the grouping rule.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable or disable the rule.",
							Optional:    true,
						},
						"rule_invert": schema.BoolAttribute{
							Description: "Invert the rule logic.",
							Optional:    true,
						},
						"rule_type": schema.StringAttribute{
							Description: "The type of the rule. Valid values: 'device_controller', 'device', 'eth_port', 'lag', 'vlan', 'tenant', 'site', 'pod', 'spineps', 'grouping_rule'.",
							Optional:    true,
						},
						"rule_value": schema.StringAttribute{
							Description: "The value for the rule.",
							Optional:    true,
						},
						"rule_value_path": schema.StringAttribute{
							Description: "The path reference for the rule value.",
							Optional:    true,
						},
						"rule_value_path_ref_type_": schema.StringAttribute{
							Description: "The reference type for rule_value_path. Valid values: 'device_controller', 'device', 'eth_port', 'lag', 'vlan', 'tenant', 'site', 'pod', 'spineps', 'grouping_rule'.",
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

func (r *verityGroupingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityGroupingRuleResourceModel
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
	groupingRuleProps := &openapi.GroupingrulesPutRequestGroupingRulesValue{
		Name: openapi.PtrString(name),
	}

	// Set string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Type", APIField: &groupingRuleProps.Type, TFValue: plan.Type},
		{FieldName: "Operation", APIField: &groupingRuleProps.Operation, TFValue: plan.Operation},
	})

	// Set boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &groupingRuleProps.Enable, TFValue: plan.Enable},
	})

	// Handle rules
	if len(plan.Rules) > 0 {
		rules := make([]openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner, len(plan.Rules))
		for i, ruleItem := range plan.Rules {
			rule := openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &rule.Enable, TFValue: ruleItem.Enable},
				{FieldName: "RuleInvert", APIField: &rule.RuleInvert, TFValue: ruleItem.RuleInvert},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RuleType", APIField: &rule.RuleType, TFValue: ruleItem.RuleType},
				{FieldName: "RuleValue", APIField: &rule.RuleValue, TFValue: ruleItem.RuleValue},
				{FieldName: "RuleValuePath", APIField: &rule.RuleValuePath, TFValue: ruleItem.RuleValuePath},
				{FieldName: "RuleValuePathRefType", APIField: &rule.RuleValuePathRefType, TFValue: ruleItem.RuleValuePathRefType},
			})

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: ruleItem.Index},
			})

			rules[i] = rule
		}
		groupingRuleProps.Rules = rules
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "grouping_rule", name, *groupingRuleProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Grouping rule %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "grouping_rules")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityGroupingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityGroupingRuleResourceModel
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

	groupingRuleName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("grouping_rule") {
		tflog.Info(ctx, fmt.Sprintf("Skipping grouping rule %s verification – trusting recent successful API operation", groupingRuleName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching grouping rules for verification of %s", groupingRuleName))

	type GroupingRulesResponse struct {
		GroupingRules map[string]interface{} `json:"grouping_rules"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "grouping_rules", groupingRuleName,
		func() (GroupingRulesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch grouping rules")
			respAPI, err := r.client.GroupingRulesAPI.GroupingrulesGet(ctx).Execute()
			if err != nil {
				return GroupingRulesResponse{}, fmt.Errorf("error reading grouping rules: %v", err)
			}
			defer respAPI.Body.Close()

			var res GroupingRulesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return GroupingRulesResponse{}, fmt.Errorf("failed to decode grouping rules response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d grouping rules", len(res.GroupingRules)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Grouping Rule %s", groupingRuleName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for grouping rule with name: %s", groupingRuleName))

	groupingRuleData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.GroupingRules,
		groupingRuleName,
		func(data interface{}) (string, bool) {
			if groupingRule, ok := data.(map[string]interface{}); ok {
				if name, ok := groupingRule["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Grouping rule with name '%s' not found in API response", groupingRuleName))
		resp.State.RemoveResource(ctx)
		return
	}

	groupingRuleMap, ok := groupingRuleData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Grouping Rule Data",
			fmt.Sprintf("Grouping rule data is not in expected format for %s", groupingRuleName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found grouping rule '%s' under API key '%s'", groupingRuleName, actualAPIName))

	state.Name = utils.MapStringFromAPI(groupingRuleMap["name"])

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"type":      &state.Type,
		"operation": &state.Operation,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(groupingRuleMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(groupingRuleMap[apiKey])
	}

	// Handle rules
	if rules, ok := groupingRuleMap["rules"].([]interface{}); ok && len(rules) > 0 {
		var rulesList []verityGroupingRuleRulesModel

		for _, r := range rules {
			rule, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			rModel := verityGroupingRuleRulesModel{
				Enable:               utils.MapBoolFromAPI(rule["enable"]),
				RuleInvert:           utils.MapBoolFromAPI(rule["rule_invert"]),
				RuleType:             utils.MapStringFromAPI(rule["rule_type"]),
				RuleValue:            utils.MapStringFromAPI(rule["rule_value"]),
				RuleValuePath:        utils.MapStringFromAPI(rule["rule_value_path"]),
				RuleValuePathRefType: utils.MapStringFromAPI(rule["rule_value_path_ref_type_"]),
				Index:                utils.MapInt64FromAPI(rule["index"]),
			}

			rulesList = append(rulesList, rModel)
		}

		state.Rules = rulesList
	} else {
		state.Rules = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityGroupingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityGroupingRuleResourceModel

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
	groupingRuleProps := openapi.GroupingrulesPutRequestGroupingRulesValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { groupingRuleProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Type, state.Type, func(v *string) { groupingRuleProps.Type = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Operation, state.Operation, func(v *string) { groupingRuleProps.Operation = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { groupingRuleProps.Enable = v }, &hasChanges)

	// Handle rules
	rulesHandler := utils.IndexedItemHandler[verityGroupingRuleRulesModel, openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner]{
		CreateNew: func(planItem verityGroupingRuleRulesModel) openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner {
			rule := openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &rule.Enable, TFValue: planItem.Enable},
				{FieldName: "RuleInvert", APIField: &rule.RuleInvert, TFValue: planItem.RuleInvert},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RuleType", APIField: &rule.RuleType, TFValue: planItem.RuleType},
				{FieldName: "RuleValue", APIField: &rule.RuleValue, TFValue: planItem.RuleValue},
				{FieldName: "RuleValuePath", APIField: &rule.RuleValuePath, TFValue: planItem.RuleValuePath},
				{FieldName: "RuleValuePathRefType", APIField: &rule.RuleValuePathRefType, TFValue: planItem.RuleValuePathRefType},
			})

			return rule
		},
		UpdateExisting: func(planItem verityGroupingRuleRulesModel, stateItem verityGroupingRuleRulesModel) (openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner, bool) {
			rule := openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { rule.Enable = v }, &fieldChanged)
			utils.CompareAndSetBoolField(planItem.RuleInvert, stateItem.RuleInvert, func(v *bool) { rule.RuleInvert = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.RuleType, stateItem.RuleType, func(v *string) { rule.RuleType = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.RuleValue, stateItem.RuleValue, func(v *string) { rule.RuleValue = v }, &fieldChanged)

			// Handle rule_value_path and rule_value_path_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.RuleValuePath, stateItem.RuleValuePath, planItem.RuleValuePathRefType, stateItem.RuleValuePathRefType,
				func(v *string) { rule.RuleValuePath = v },
				func(v *string) { rule.RuleValuePathRefType = v },
				"rule_value_path", "rule_value_path_ref_type_",
				&fieldChanged, &resp.Diagnostics,
			) {
				return rule, false
			}

			return rule, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner {
			rule := openapi.GroupingrulesPutRequestGroupingRulesValueRulesInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rule.Index, TFValue: types.Int64Value(index)},
			})
			return rule
		},
	}

	changedRules, rulesChanged := utils.ProcessIndexedArrayUpdates(plan.Rules, state.Rules, rulesHandler)
	if rulesChanged {
		groupingRuleProps.Rules = changedRules
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "grouping_rule", name, groupingRuleProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Grouping rule %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "grouping_rules")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityGroupingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityGroupingRuleResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "grouping_rule", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Grouping rule %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "grouping_rules")
	resp.State.RemoveResource(ctx)
}

func (r *verityGroupingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
