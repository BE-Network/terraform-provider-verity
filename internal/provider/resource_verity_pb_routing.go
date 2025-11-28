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
	_ resource.Resource                = &verityPBRoutingResource{}
	_ resource.ResourceWithConfigure   = &verityPBRoutingResource{}
	_ resource.ResourceWithImportState = &verityPBRoutingResource{}
)

func NewVerityPBRoutingResource() resource.Resource {
	return &verityPBRoutingResource{}
}

type verityPBRoutingResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityPBRoutingResourceModel struct {
	Name   types.String                 `tfsdk:"name"`
	Enable types.Bool                   `tfsdk:"enable"`
	Policy []verityPBRoutingPolicyModel `tfsdk:"policy"`
}

type verityPBRoutingPolicyModel struct {
	Enable              types.Bool   `tfsdk:"enable"`
	PbRoutingAcl        types.String `tfsdk:"pb_routing_acl"`
	PbRoutingAclRefType types.String `tfsdk:"pb_routing_acl_ref_type_"`
	Index               types.Int64  `tfsdk:"index"`
}

func (p verityPBRoutingPolicyModel) GetIndex() types.Int64 {
	return p.Index
}

func (r *verityPBRoutingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pb_routing"
}

func (r *verityPBRoutingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPBRoutingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Policy-Based Routing resource",
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
		},
		Blocks: map[string]schema.Block{
			"policy": schema.ListNestedBlock{
				Description: "Policy configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"pb_routing_acl": schema.StringAttribute{
							Description: "Path to the PB Routing ACL",
							Optional:    true,
						},
						"pb_routing_acl_ref_type_": schema.StringAttribute{
							Description: "Object type for pb_routing_acl field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityPBRoutingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPBRoutingResourceModel
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
	pbRoutingProps := &openapi.PolicybasedroutingPutRequestPbRoutingValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &pbRoutingProps.Enable, TFValue: plan.Enable},
	})

	// Handle policy
	if len(plan.Policy) > 0 {
		policies := make([]openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner, len(plan.Policy))
		for i, policyItem := range plan.Policy {
			policy := openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &policy.Index, TFValue: policyItem.Index},
			})
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &policy.Enable, TFValue: policyItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PbRoutingAcl", APIField: &policy.PbRoutingAcl, TFValue: policyItem.PbRoutingAcl},
				{FieldName: "PbRoutingAclRefType", APIField: &policy.PbRoutingAclRefType, TFValue: policyItem.PbRoutingAclRefType},
			})
			policies[i] = policy
		}
		pbRoutingProps.Policy = policies
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "pb_routing", name, *pbRoutingProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Routing %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_routing")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityPBRoutingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPBRoutingResourceModel
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

	pbRoutingName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("pb_routing") {
		tflog.Info(ctx, fmt.Sprintf("Skipping PB Routing %s verification – trusting recent successful API operation", pbRoutingName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching PB Routing for verification of %s", pbRoutingName))

	type PBRoutingResponse struct {
		PbRouting map[string]interface{} `json:"pb_routing"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "pb_routing", pbRoutingName,
		func() (PBRoutingResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch PB Routing")
			respAPI, err := r.client.PBRoutingAPI.PolicybasedroutingGet(ctx).Execute()
			if err != nil {
				return PBRoutingResponse{}, fmt.Errorf("error reading PB Routing: %v", err)
			}
			defer respAPI.Body.Close()

			var res PBRoutingResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PBRoutingResponse{}, fmt.Errorf("failed to decode PB Routing response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d PB Routing entries", len(res.PbRouting)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read PB Routing %s", pbRoutingName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for PB Routing with name: %s", pbRoutingName))

	pbRoutingData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.PbRouting,
		pbRoutingName,
		func(data interface{}) (string, bool) {
			if pbRouting, ok := data.(map[string]interface{}); ok {
				if name, ok := pbRouting["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("PB Routing with name '%s' not found in API response", pbRoutingName))
		resp.State.RemoveResource(ctx)
		return
	}

	pbRoutingMap, ok := pbRoutingData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid PB Routing Data",
			fmt.Sprintf("PB Routing data is not in expected format for %s", pbRoutingName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found PB Routing '%s' under API key '%s'", pbRoutingName, actualAPIName))

	state.Name = utils.MapStringFromAPI(pbRoutingMap["name"])

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(pbRoutingMap[apiKey])
	}

	// Handle policy
	if policies, ok := pbRoutingMap["policy"].([]interface{}); ok && len(policies) > 0 {
		var policyList []verityPBRoutingPolicyModel

		for _, p := range policies {
			policyMap, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			policyModel := verityPBRoutingPolicyModel{
				Enable:              utils.MapBoolFromAPI(policyMap["enable"]),
				PbRoutingAcl:        utils.MapStringFromAPI(policyMap["pb_routing_acl"]),
				PbRoutingAclRefType: utils.MapStringFromAPI(policyMap["pb_routing_acl_ref_type_"]),
				Index:               utils.MapInt64FromAPI(policyMap["index"]),
			}

			policyList = append(policyList, policyModel)
		}

		state.Policy = policyList
	} else {
		state.Policy = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPBRoutingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPBRoutingResourceModel

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
	pbRoutingProps := openapi.PolicybasedroutingPutRequestPbRoutingValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { pbRoutingProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { pbRoutingProps.Enable = v }, &hasChanges)

	// Handle policy
	policyHandler := utils.IndexedItemHandler[verityPBRoutingPolicyModel, openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner]{
		CreateNew: func(planItem verityPBRoutingPolicyModel) openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner {
			policy := openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &policy.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &policy.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PbRoutingAcl", APIField: &policy.PbRoutingAcl, TFValue: planItem.PbRoutingAcl},
				{FieldName: "PbRoutingAclRefType", APIField: &policy.PbRoutingAclRefType, TFValue: planItem.PbRoutingAclRefType},
			})

			return policy
		},
		UpdateExisting: func(planItem verityPBRoutingPolicyModel, stateItem verityPBRoutingPolicyModel) (openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner, bool) {
			policy := openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &policy.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { policy.Enable = v }, &fieldChanged)

			// Handle pb_routing_acl and pb_routing_acl_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.PbRoutingAcl, stateItem.PbRoutingAcl, planItem.PbRoutingAclRefType, stateItem.PbRoutingAclRefType,
				func(v *string) { policy.PbRoutingAcl = v },
				func(v *string) { policy.PbRoutingAclRefType = v },
				"pb_routing_acl", "pb_routing_acl_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return policy, false
			}

			return policy, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner {
			policy := openapi.PolicybasedroutingPutRequestPbRoutingValuePolicyInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &policy.Index, TFValue: types.Int64Value(index)},
			})
			return policy
		},
	}

	changedPolicies, policiesChanged := utils.ProcessIndexedArrayUpdates(plan.Policy, state.Policy, policyHandler)
	if policiesChanged {
		pbRoutingProps.Policy = changedPolicies
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "pb_routing", name, pbRoutingProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Routing %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_routing")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityPBRoutingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPBRoutingResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "pb_routing", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Routing %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_routing")
	resp.State.RemoveResource(ctx)
}

func (r *verityPBRoutingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
