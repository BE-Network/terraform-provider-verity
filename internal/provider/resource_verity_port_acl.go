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
	_ resource.Resource                = &verityPortAclResource{}
	_ resource.ResourceWithConfigure   = &verityPortAclResource{}
	_ resource.ResourceWithImportState = &verityPortAclResource{}
	_ resource.ResourceWithModifyPlan  = &verityPortAclResource{}
)

const portAclResourceType = "portacls"

func NewVerityPortAclResource() resource.Resource {
	return &verityPortAclResource{}
}

type verityPortAclResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityPortAclResourceModel struct {
	Name       types.String               `tfsdk:"name"`
	Enable     types.Bool                 `tfsdk:"enable"`
	Ipv4Permit []verityPortAclFilterModel `tfsdk:"ipv4_permit"`
	Ipv4Deny   []verityPortAclFilterModel `tfsdk:"ipv4_deny"`
	Ipv6Permit []verityPortAclFilterModel `tfsdk:"ipv6_permit"`
	Ipv6Deny   []verityPortAclFilterModel `tfsdk:"ipv6_deny"`
}

type verityPortAclFilterModel struct {
	Enable        types.Bool   `tfsdk:"enable"`
	Filter        types.String `tfsdk:"filter"`
	FilterRefType types.String `tfsdk:"filter_ref_type_"`
	Index         types.Int64  `tfsdk:"index"`
}

func (m verityPortAclFilterModel) GetIndex() types.Int64 {
	return m.Index
}

func (r *verityPortAclResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_acl"
}

func (r *verityPortAclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPortAclResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Port ACL",
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
		},
		Blocks: map[string]schema.Block{
			"ipv4_permit": schema.ListNestedBlock{
				Description: "List of IPv4 permit filters",
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
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ipv4_deny": schema.ListNestedBlock{
				Description: "List of IPv4 deny filters",
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
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ipv6_permit": schema.ListNestedBlock{
				Description: "List of IPv6 permit filters",
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
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ipv6_deny": schema.ListNestedBlock{
				Description: "List of IPv6 deny filters",
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
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityPortAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPortAclResourceModel
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
	portAclProps := &openapi.PortaclsPutRequestPortAclValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &portAclProps.Enable, TFValue: plan.Enable},
	})

	// Handle IPv4 Permit
	if len(plan.Ipv4Permit) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, len(plan.Ipv4Permit))
		for i, item := range plan.Ipv4Permit {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: item.Enable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: item.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: item.FilterRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: item.Index},
			})

			filters[i] = filter
		}
		portAclProps.Ipv4Permit = filters
	}

	// Handle IPv4 Deny
	if len(plan.Ipv4Deny) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, len(plan.Ipv4Deny))
		for i, item := range plan.Ipv4Deny {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: item.Enable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: item.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: item.FilterRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: item.Index},
			})

			filters[i] = filter
		}
		portAclProps.Ipv4Deny = filters
	}

	// Handle IPv6 Permit
	if len(plan.Ipv6Permit) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, len(plan.Ipv6Permit))
		for i, item := range plan.Ipv6Permit {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: item.Enable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: item.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: item.FilterRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: item.Index},
			})

			filters[i] = filter
		}
		portAclProps.Ipv6Permit = filters
	}

	// Handle IPv6 Deny
	if len(plan.Ipv6Deny) > 0 {
		filters := make([]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, len(plan.Ipv6Deny))
		for i, item := range plan.Ipv6Deny {
			filter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filter.Enable, TFValue: item.Enable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filter.Filter, TFValue: item.Filter},
				{FieldName: "FilterRefType", APIField: &filter.FilterRefType, TFValue: item.FilterRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: item.Index},
			})

			filters[i] = filter
		}
		portAclProps.Ipv6Deny = filters
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "port_acl", name, *portAclProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Port ACL %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "port_acls")

	var minState verityPortAclResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if portAclData, exists := bulkMgr.GetResourceResponse("port_acl", name); exists {
			state := populatePortAclState(ctx, minState, portAclData, r.provCtx.mode)
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

func (r *verityPortAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPortAclResourceModel
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if portAclData, exists := r.bulkOpsMgr.GetResourceResponse("port_acl", name); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached port_acl data for %s from recent operation", name))
			state = populatePortAclState(ctx, state, portAclData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("port_acl") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Port ACL %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Port ACL for verification of %s", name))

	type PortAclsResponse struct {
		PortAcl map[string]map[string]interface{} `json:"port_acl"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "port_acls", name,
		func() (PortAclsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Port ACLs")
			respAPI, err := r.client.PortACLsAPI.PortaclsGet(ctx).Execute()
			if err != nil {
				return PortAclsResponse{}, fmt.Errorf("error reading Port ACL: %v", err)
			}
			defer respAPI.Body.Close()

			var res PortAclsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PortAclsResponse{}, fmt.Errorf("failed to decode Port ACLs response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Port ACLs from API", len(res.PortAcl)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Port ACL %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Port ACL with name: %s", name))

	portAclData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.PortAcl,
		name,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Port ACL with name '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	portAclMap, ok := (interface{}(portAclData)).(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Port ACL Data",
			fmt.Sprintf("Port ACL data is not in expected format for %s", name),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Port ACL '%s' under API key '%s'", name, actualAPIName))

	state = populatePortAclState(ctx, state, portAclMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPortAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPortAclResourceModel

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
	portAclProps := openapi.PortaclsPutRequestPortAclValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { portAclProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { portAclProps.Enable = v }, &hasChanges)

	// Handle IPv4 Permit
	changedIpv4Permits, ipv4PermitsChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv4Permit, state.Ipv4Permit,
		utils.IndexedItemHandler[verityPortAclFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner]{
			CreateNew: func(planItem verityPortAclFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
				newFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &newFilter.Enable, TFValue: planItem.Enable},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Filter", APIField: &newFilter.Filter, TFValue: planItem.Filter},
					{FieldName: "FilterRefType", APIField: &newFilter.FilterRefType, TFValue: planItem.FilterRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newFilter.Index, TFValue: planItem.Index},
				})

				return newFilter
			},
			UpdateExisting: func(planItem verityPortAclFilterModel, stateItem verityPortAclFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool) {
				updateFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateFilter.Enable = v }, &fieldChanged)

				// Handle filter and filter_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
					func(v *string) { updateFilter.Filter = v },
					func(v *string) { updateFilter.FilterRefType = v },
					"filter", "filter_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updateFilter, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateFilter.Index = v }, &fieldChanged)

				return updateFilter, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
				return openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv4PermitsChanged {
		portAclProps.Ipv4Permit = changedIpv4Permits
		hasChanges = true
	}

	// Handle IPv4 Deny
	changedIpv4Denies, ipv4DeniesChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv4Deny, state.Ipv4Deny,
		utils.IndexedItemHandler[verityPortAclFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner]{
			CreateNew: func(planItem verityPortAclFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
				newFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &newFilter.Enable, TFValue: planItem.Enable},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Filter", APIField: &newFilter.Filter, TFValue: planItem.Filter},
					{FieldName: "FilterRefType", APIField: &newFilter.FilterRefType, TFValue: planItem.FilterRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newFilter.Index, TFValue: planItem.Index},
				})

				return newFilter
			},
			UpdateExisting: func(planItem verityPortAclFilterModel, stateItem verityPortAclFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool) {
				updateFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateFilter.Enable = v }, &fieldChanged)

				// Handle filter and filter_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
					func(v *string) { updateFilter.Filter = v },
					func(v *string) { updateFilter.FilterRefType = v },
					"filter", "filter_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updateFilter, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateFilter.Index = v }, &fieldChanged)

				return updateFilter, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner {
				return openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv4DeniesChanged {
		portAclProps.Ipv4Deny = changedIpv4Denies
		hasChanges = true
	}

	// Handle IPv6 Permit
	changedIpv6Permits, ipv6PermitsChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv6Permit, state.Ipv6Permit,
		utils.IndexedItemHandler[verityPortAclFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner]{
			CreateNew: func(planItem verityPortAclFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
				newFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &newFilter.Enable, TFValue: planItem.Enable},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Filter", APIField: &newFilter.Filter, TFValue: planItem.Filter},
					{FieldName: "FilterRefType", APIField: &newFilter.FilterRefType, TFValue: planItem.FilterRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newFilter.Index, TFValue: planItem.Index},
				})

				return newFilter
			},
			UpdateExisting: func(planItem verityPortAclFilterModel, stateItem verityPortAclFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool) {
				updateFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateFilter.Enable = v }, &fieldChanged)

				// Handle filter and filter_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
					func(v *string) { updateFilter.Filter = v },
					func(v *string) { updateFilter.FilterRefType = v },
					"filter", "filter_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updateFilter, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateFilter.Index = v }, &fieldChanged)

				return updateFilter, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
				return openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv6PermitsChanged {
		portAclProps.Ipv6Permit = changedIpv6Permits
		hasChanges = true
	}

	// Handle IPv6 Deny
	changedIpv6Denies, ipv6DeniesChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv6Deny, state.Ipv6Deny,
		utils.IndexedItemHandler[verityPortAclFilterModel, openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner]{
			CreateNew: func(planItem verityPortAclFilterModel) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
				newFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &newFilter.Enable, TFValue: planItem.Enable},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Filter", APIField: &newFilter.Filter, TFValue: planItem.Filter},
					{FieldName: "FilterRefType", APIField: &newFilter.FilterRefType, TFValue: planItem.FilterRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newFilter.Index, TFValue: planItem.Index},
				})

				return newFilter
			},
			UpdateExisting: func(planItem verityPortAclFilterModel, stateItem verityPortAclFilterModel) (openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool) {
				updateFilter := openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateFilter.Enable = v }, &fieldChanged)

				// Handle filter and filter_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Filter, stateItem.Filter, planItem.FilterRefType, stateItem.FilterRefType,
					func(v *string) { updateFilter.Filter = v },
					func(v *string) { updateFilter.FilterRefType = v },
					"filter", "filter_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updateFilter, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateFilter.Index = v }, &fieldChanged)

				return updateFilter, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner {
				return openapi.PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv6DeniesChanged {
		portAclProps.Ipv6Deny = changedIpv6Denies
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "port_acl", name, portAclProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Port ACL %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "port_acls")

	var minState verityPortAclResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if portAclData, exists := bulkMgr.GetResourceResponse("port_acl", name); exists {
			newState := populatePortAclState(ctx, minState, portAclData, r.provCtx.mode)
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

func (r *verityPortAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPortAclResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "port_acl", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Port ACL %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "port_acls")
	resp.State.RemoveResource(ctx)
}

func (r *verityPortAclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populatePortAclState(ctx context.Context, state verityPortAclResourceModel, data map[string]interface{}, mode string) verityPortAclResourceModel {
	const resourceType = portAclResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Helper function to parse filter arrays with mode awareness
	parseFilters := func(apiFilters []interface{}, blockName string) []verityPortAclFilterModel {
		var filters []verityPortAclFilterModel
		for _, f := range apiFilters {
			filter, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPortAclFilterModel{
				Enable:        utils.MapBoolWithModeNested(filter, "enable", resourceType, blockName+".enable", mode),
				Filter:        utils.MapStringWithModeNested(filter, "filter", resourceType, blockName+".filter", mode),
				FilterRefType: utils.MapStringWithModeNested(filter, "filter_ref_type_", resourceType, blockName+".filter_ref_type_", mode),
				Index:         utils.MapInt64WithModeNested(filter, "index", resourceType, blockName+".index", mode),
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

func (r *verityPortAclResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityPortAclResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = portAclResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "ipv4_permit",
		ItemCount:    len(plan.Ipv4Permit),
		StringFields: []string{"filter", "filter_ref_type_"},
		BoolFields:   []string{"enable"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "ipv4_deny",
		ItemCount:    len(plan.Ipv4Deny),
		StringFields: []string{"filter", "filter_ref_type_"},
		BoolFields:   []string{"enable"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "ipv6_permit",
		ItemCount:    len(plan.Ipv6Permit),
		StringFields: []string{"filter", "filter_ref_type_"},
		BoolFields:   []string{"enable"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "ipv6_deny",
		ItemCount:    len(plan.Ipv6Deny),
		StringFields: []string{"filter", "filter_ref_type_"},
		BoolFields:   []string{"enable"},
	})
}
