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
	_ resource.Resource                = &verityPBRoutingACLResource{}
	_ resource.ResourceWithConfigure   = &verityPBRoutingACLResource{}
	_ resource.ResourceWithImportState = &verityPBRoutingACLResource{}
	_ resource.ResourceWithModifyPlan  = &verityPBRoutingACLResource{}
)

const pbRoutingAclResourceType = "policybasedroutingacl"

func NewVerityPBRoutingACLResource() resource.Resource {
	return &verityPBRoutingACLResource{}
}

type verityPBRoutingACLResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityPBRoutingACLResourceModel struct {
	Name        types.String                    `tfsdk:"name"`
	Enable      types.Bool                      `tfsdk:"enable"`
	IpvProtocol types.String                    `tfsdk:"ipv_protocol"`
	NextHopIps  types.String                    `tfsdk:"next_hop_ips"`
	Ipv4Permit  []verityPBRoutingACLFilterModel `tfsdk:"ipv4_permit"`
	Ipv4Deny    []verityPBRoutingACLFilterModel `tfsdk:"ipv4_deny"`
	Ipv6Permit  []verityPBRoutingACLFilterModel `tfsdk:"ipv6_permit"`
	Ipv6Deny    []verityPBRoutingACLFilterModel `tfsdk:"ipv6_deny"`
}

type verityPBRoutingACLFilterModel struct {
	Enable        types.Bool   `tfsdk:"enable"`
	Filter        types.String `tfsdk:"filter"`
	FilterRefType types.String `tfsdk:"filter_ref_type_"`
	Index         types.Int64  `tfsdk:"index"`
}

func (p verityPBRoutingACLFilterModel) GetIndex() types.Int64 {
	return p.Index
}

func (r *verityPBRoutingACLResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pb_routing_acl"
}

func (r *verityPBRoutingACLResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPBRoutingACLResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Policy-Based Routing ACL resource",
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
			"ipv_protocol": schema.StringAttribute{
				Description: "IPv4 or IPv6",
				Optional:    true,
				Computed:    true,
			},
			"next_hop_ips": schema.StringAttribute{
				Description: "Next hop IP addresses",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"ipv4_permit": schema.ListNestedBlock{
				Description: "IPv4 permit filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
							Computed:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
							Computed:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ipv4_deny": schema.ListNestedBlock{
				Description: "IPv4 deny filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
							Computed:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
							Computed:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ipv6_permit": schema.ListNestedBlock{
				Description: "IPv6 permit filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
							Computed:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
							Computed:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ipv6_deny": schema.ListNestedBlock{
				Description: "IPv6 deny filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
							Computed:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
							Computed:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityPBRoutingACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPBRoutingACLResourceModel
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
	pbRoutingACLProps := &openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &pbRoutingACLProps.Enable, TFValue: plan.Enable},
	})

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "IpvProtocol", APIField: &pbRoutingACLProps.IpvProtocol, TFValue: plan.IpvProtocol},
		{FieldName: "NextHopIps", APIField: &pbRoutingACLProps.NextHopIps, TFValue: plan.NextHopIps},
	})

	// Handle ipv4_permit
	if len(plan.Ipv4Permit) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, len(plan.Ipv4Permit))
		for i, filterItem := range plan.Ipv4Permit {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: filterItem.Index},
			})
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: filterItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: filterItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: filterItem.FilterRefType},
			})
			filters[i] = filter
		}
		pbRoutingACLProps.Ipv4Permit = filters
	}

	// Handle ipv4_deny
	if len(plan.Ipv4Deny) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, len(plan.Ipv4Deny))
		for i, filterItem := range plan.Ipv4Deny {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: filterItem.Index},
			})
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: filterItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: filterItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: filterItem.FilterRefType},
			})
			filters[i] = filter
		}
		pbRoutingACLProps.Ipv4Deny = filters
	}

	// Handle ipv6_permit
	if len(plan.Ipv6Permit) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, len(plan.Ipv6Permit))
		for i, filterItem := range plan.Ipv6Permit {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: filterItem.Index},
			})
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: filterItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: filterItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: filterItem.FilterRefType},
			})
			filters[i] = filter
		}
		pbRoutingACLProps.Ipv6Permit = filters
	}

	// Handle ipv6_deny
	if len(plan.Ipv6Deny) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, len(plan.Ipv6Deny))
		for i, filterItem := range plan.Ipv6Deny {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: filterItem.Index},
			})
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: filterItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: filterItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: filterItem.FilterRefType},
			})
			filters[i] = filter
		}
		pbRoutingACLProps.Ipv6Deny = filters
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "pb_routing_acl", name, *pbRoutingACLProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Routing ACL %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_routing_acl")

	var minState verityPBRoutingACLResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pbRoutingACLData, exists := bulkMgr.GetResourceResponse("pb_routing_acl", name); exists {
			state := populatePBRoutingACLState(ctx, minState, pbRoutingACLData, r.provCtx.mode)
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

	r.Read(ctx, readReq, &readResp)
}

func (r *verityPBRoutingACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPBRoutingACLResourceModel
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

	pbRoutingACLName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if pbRoutingACLData, exists := r.bulkOpsMgr.GetResourceResponse("pb_routing_acl", pbRoutingACLName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached pb routing acl data for %s from recent operation", pbRoutingACLName))
			state = populatePBRoutingACLState(ctx, state, pbRoutingACLData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("pb_routing_acl") {
		tflog.Info(ctx, fmt.Sprintf("Skipping PB Routing ACL %s verification – trusting recent successful API operation", pbRoutingACLName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching PB Routing ACL for verification of %s", pbRoutingACLName))

	type PBRoutingACLResponse struct {
		PbRoutingAcl map[string]interface{} `json:"pb_routing_acl"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "pb_routing_acl", pbRoutingACLName,
		func() (PBRoutingACLResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch PB Routing ACL")
			respAPI, err := r.client.PBRoutingACLAPI.PolicybasedroutingaclGet(ctx).Execute()
			if err != nil {
				return PBRoutingACLResponse{}, fmt.Errorf("error reading PB Routing ACL: %v", err)
			}
			defer respAPI.Body.Close()

			var res PBRoutingACLResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PBRoutingACLResponse{}, fmt.Errorf("failed to decode PB Routing ACL response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d PB Routing ACL entries", len(res.PbRoutingAcl)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read PB Routing ACL %s", pbRoutingACLName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for PB Routing ACL with name: %s", pbRoutingACLName))

	pbRoutingACLData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.PbRoutingAcl,
		pbRoutingACLName,
		func(data interface{}) (string, bool) {
			if pbRoutingACL, ok := data.(map[string]interface{}); ok {
				if name, ok := pbRoutingACL["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("PB Routing ACL with name '%s' not found in API response", pbRoutingACLName))
		resp.State.RemoveResource(ctx)
		return
	}

	pbRoutingACLMap, ok := pbRoutingACLData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid PB Routing ACL Data",
			fmt.Sprintf("PB Routing ACL data is not in expected format for %s", pbRoutingACLName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found PB Routing ACL '%s' under API key '%s'", pbRoutingACLName, actualAPIName))

	state = populatePBRoutingACLState(ctx, state, pbRoutingACLMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPBRoutingACLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPBRoutingACLResourceModel

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
	pbRoutingACLProps := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { pbRoutingACLProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.IpvProtocol, state.IpvProtocol, func(v *string) { pbRoutingACLProps.IpvProtocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NextHopIps, state.NextHopIps, func(v *string) { pbRoutingACLProps.NextHopIps = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { pbRoutingACLProps.Enable = v }, &hasChanges)

	// Handle ipv4_permit
	ipv4PermitHandler := utils.IndexedItemHandler[verityPBRoutingACLFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner]{
		CreateNew: func(planItem verityPBRoutingACLFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: planItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: planItem.FilterRefType},
			})

			return filter
		},
		UpdateExisting: func(planItem verityPBRoutingACLFilterModel, stateItem verityPBRoutingACLFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool) {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { filter.Enable = v }, &fieldChanged)

			// Handle filter and filter_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
				func(v *string) { filter.Filter = v },
				func(v *string) { filter.FilterRefType = v },
				"filter", "filter_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return filter, false
			}

			return filter, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: types.Int64Value(index)},
			})
			return filter
		},
	}

	changedIpv4Permit, ipv4PermitChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv4Permit, state.Ipv4Permit, ipv4PermitHandler)
	if ipv4PermitChanged {
		pbRoutingACLProps.Ipv4Permit = changedIpv4Permit
		hasChanges = true
	}

	// Handle ipv4_deny
	ipv4DenyHandler := utils.IndexedItemHandler[verityPBRoutingACLFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner]{
		CreateNew: func(planItem verityPBRoutingACLFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: planItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: planItem.FilterRefType},
			})

			return filter
		},
		UpdateExisting: func(planItem verityPBRoutingACLFilterModel, stateItem verityPBRoutingACLFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool) {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { filter.Enable = v }, &fieldChanged)

			// Handle filter and filter_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
				func(v *string) { filter.Filter = v },
				func(v *string) { filter.FilterRefType = v },
				"filter", "filter_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return filter, false
			}

			return filter, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: types.Int64Value(index)},
			})
			return filter
		},
	}

	changedIpv4Deny, ipv4DenyChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv4Deny, state.Ipv4Deny, ipv4DenyHandler)
	if ipv4DenyChanged {
		pbRoutingACLProps.Ipv4Deny = changedIpv4Deny
		hasChanges = true
	}

	// Handle ipv6_permit
	ipv6PermitHandler := utils.IndexedItemHandler[verityPBRoutingACLFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner]{
		CreateNew: func(planItem verityPBRoutingACLFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: planItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: planItem.FilterRefType},
			})

			return filter
		},
		UpdateExisting: func(planItem verityPBRoutingACLFilterModel, stateItem verityPBRoutingACLFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool) {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { filter.Enable = v }, &fieldChanged)

			// Handle filter and filter_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
				func(v *string) { filter.Filter = v },
				func(v *string) { filter.FilterRefType = v },
				"filter", "filter_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return filter, false
			}

			return filter, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: types.Int64Value(index)},
			})
			return filter
		},
	}

	changedIpv6Permit, ipv6PermitChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv6Permit, state.Ipv6Permit, ipv6PermitHandler)
	if ipv6PermitChanged {
		pbRoutingACLProps.Ipv6Permit = changedIpv6Permit
		hasChanges = true
	}

	// Handle ipv6_deny
	ipv6DenyHandler := utils.IndexedItemHandler[verityPBRoutingACLFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner]{
		CreateNew: func(planItem verityPBRoutingACLFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: planItem.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: planItem.FilterRefType},
			})

			return filter
		},
		UpdateExisting: func(planItem verityPBRoutingACLFilterModel, stateItem verityPBRoutingACLFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool) {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { filter.Enable = v }, &fieldChanged)

			// Handle filter and filter_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
				func(v *string) { filter.Filter = v },
				func(v *string) { filter.FilterRefType = v },
				"filter", "filter_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return filter, false
			}

			return filter, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: types.Int64Value(index)},
			})
			return filter
		},
	}

	changedIpv6Deny, ipv6DenyChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv6Deny, state.Ipv6Deny, ipv6DenyHandler)
	if ipv6DenyChanged {
		pbRoutingACLProps.Ipv6Deny = changedIpv6Deny
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "pb_routing_acl", name, pbRoutingACLProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Routing ACL %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_routing_acl")

	var minState verityPBRoutingACLResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pbRoutingACLData, exists := bulkMgr.GetResourceResponse("pb_routing_acl", name); exists {
			newState := populatePBRoutingACLState(ctx, minState, pbRoutingACLData, r.provCtx.mode)
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

func (r *verityPBRoutingACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPBRoutingACLResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "pb_routing_acl", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Routing ACL %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_routing_acl")
	resp.State.RemoveResource(ctx)
}

func (r *verityPBRoutingACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populatePBRoutingACLState(ctx context.Context, state verityPBRoutingACLResourceModel, data map[string]interface{}, mode string) verityPBRoutingACLResourceModel {
	const resourceType = pbRoutingAclResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// String fields
	state.IpvProtocol = utils.MapStringWithMode(data, "ipv_protocol", resourceType, mode)
	state.NextHopIps = utils.MapStringWithMode(data, "next_hop_ips", resourceType, mode)

	// Helper function to parse filter arrays with mode awareness
	parseFilters := func(apiFilters []interface{}, blockName string) []verityPBRoutingACLFilterModel {
		var filters []verityPBRoutingACLFilterModel
		for _, f := range apiFilters {
			filterMap, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPBRoutingACLFilterModel{
				Enable:        utils.MapBoolWithModeNested(filterMap, "enable", resourceType, blockName+".enable", mode),
				Filter:        utils.MapStringWithModeNested(filterMap, "filter", resourceType, blockName+".filter", mode),
				FilterRefType: utils.MapStringWithModeNested(filterMap, "filter_ref_type_", resourceType, blockName+".filter_ref_type_", mode),
				Index:         utils.MapInt64WithModeNested(filterMap, "index", resourceType, blockName+".index", mode),
			}
			filters = append(filters, filterModel)
		}
		return filters
	}

	// Handle filter arrays with mode awareness
	if utils.FieldAppliesToMode(resourceType, "ipv4_permit", mode) {
		if ipv4Permit, ok := data["ipv4_permit"].([]interface{}); ok && len(ipv4Permit) > 0 {
			state.Ipv4Permit = parseFilters(ipv4Permit, "ipv4_permit")
		} else {
			state.Ipv4Permit = nil
		}
	} else {
		state.Ipv4Permit = nil
	}

	if utils.FieldAppliesToMode(resourceType, "ipv4_deny", mode) {
		if ipv4Deny, ok := data["ipv4_deny"].([]interface{}); ok && len(ipv4Deny) > 0 {
			state.Ipv4Deny = parseFilters(ipv4Deny, "ipv4_deny")
		} else {
			state.Ipv4Deny = nil
		}
	} else {
		state.Ipv4Deny = nil
	}

	if utils.FieldAppliesToMode(resourceType, "ipv6_permit", mode) {
		if ipv6Permit, ok := data["ipv6_permit"].([]interface{}); ok && len(ipv6Permit) > 0 {
			state.Ipv6Permit = parseFilters(ipv6Permit, "ipv6_permit")
		} else {
			state.Ipv6Permit = nil
		}
	} else {
		state.Ipv6Permit = nil
	}

	if utils.FieldAppliesToMode(resourceType, "ipv6_deny", mode) {
		if ipv6Deny, ok := data["ipv6_deny"].([]interface{}); ok && len(ipv6Deny) > 0 {
			state.Ipv6Deny = parseFilters(ipv6Deny, "ipv6_deny")
		} else {
			state.Ipv6Deny = nil
		}
	} else {
		state.Ipv6Deny = nil
	}

	return state
}

func (r *verityPBRoutingACLResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityPBRoutingACLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = pbRoutingAclResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"ipv_protocol", "next_hop_ips",
	)

	nullifier.NullifyBools(
		"enable",
	)
}
