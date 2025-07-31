package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	_ resource.Resource                = &verityPacketBrokerResource{}
	_ resource.ResourceWithConfigure   = &verityPacketBrokerResource{}
	_ resource.ResourceWithImportState = &verityPacketBrokerResource{}
)

func NewVerityPacketBrokerResource() resource.Resource {
	return &verityPacketBrokerResource{}
}

type verityPacketBrokerResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityPacketBrokerResourceModel struct {
	Name       types.String                    `tfsdk:"name"`
	Enable     types.Bool                      `tfsdk:"enable"`
	Ipv4Permit []verityPacketBrokerFilterModel `tfsdk:"ipv4_permit"`
	Ipv4Deny   []verityPacketBrokerFilterModel `tfsdk:"ipv4_deny"`
	Ipv6Permit []verityPacketBrokerFilterModel `tfsdk:"ipv6_permit"`
	Ipv6Deny   []verityPacketBrokerFilterModel `tfsdk:"ipv6_deny"`
}

type verityPacketBrokerFilterModel struct {
	Enable        types.Bool   `tfsdk:"enable"`
	Filter        types.String `tfsdk:"filter"`
	FilterRefType types.String `tfsdk:"filter_ref_type_"`
	Index         types.Int64  `tfsdk:"index"`
}

func (r *verityPacketBrokerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_packet_broker"
}

func (r *verityPacketBrokerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPacketBrokerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity PB Egress Profile",
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
			"ipv4_permit": schema.ListNestedBlock{
				Description: "IPv4 Permit filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"ipv4_deny": schema.ListNestedBlock{
				Description: "IPv4 Deny filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"ipv6_permit": schema.ListNestedBlock{
				Description: "IPv6 Permit filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"ipv6_deny": schema.ListNestedBlock{
				Description: "IPv6 Deny filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityPacketBrokerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPacketBrokerResourceModel
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
	pbProps := openapi.PacketbrokerPutRequestPbEgressProfileValue{}
	pbProps.Name = openapi.PtrString(name)

	if !plan.Enable.IsNull() {
		pbProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if len(plan.Ipv4Permit) > 0 {
		ipv4Permit := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, len(plan.Ipv4Permit))
		for i, filter := range plan.Ipv4Permit {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv4Permit[i] = filterItem
		}
		pbProps.Ipv4Permit = ipv4Permit
	}

	if len(plan.Ipv4Deny) > 0 {
		ipv4Deny := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, len(plan.Ipv4Deny))
		for i, filter := range plan.Ipv4Deny {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv4Deny[i] = filterItem
		}
		pbProps.Ipv4Deny = ipv4Deny
	}

	if len(plan.Ipv6Permit) > 0 {
		ipv6Permit := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, len(plan.Ipv6Permit))
		for i, filter := range plan.Ipv6Permit {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv6Permit[i] = filterItem
		}
		pbProps.Ipv6Permit = ipv6Permit
	}

	if len(plan.Ipv6Deny) > 0 {
		ipv6Deny := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, len(plan.Ipv6Deny))
		for i, filter := range plan.Ipv6Deny {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv6Deny[i] = filterItem
		}
		pbProps.Ipv6Deny = ipv6Deny
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "packet_broker", name, pbProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for packet broker creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create PB Egress Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Egress Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_egress_profiles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityPacketBrokerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPacketBrokerResourceModel
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

	pbName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("packet_broker") {
		tflog.Info(ctx, fmt.Sprintf("Skipping PB Egress Profile %s verification – trusting recent successful API operation", pbName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching PB Egress Profiles for verification of %s", pbName))

	type PacketBrokerResponse struct {
		PbEgressProfile map[string]interface{} `json:"pb_egress_profile"`
	}

	var result PacketBrokerResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		pbData, fetchErr := getCachedResponse(ctx, r.provCtx, "pb_egress_profiles", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch PB Egress Profiles")
			respAPI, err := r.client.PacketBrokerAPI.PacketbrokerGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading PB Egress Profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res PacketBrokerResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode PB Egress Profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d PB Egress Profiles", len(res.PbEgressProfile)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch PB Egress Profiles on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = pbData.(PacketBrokerResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read PB Egress Profile %s", pbName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for PB Egress Profile with ID: %s", pbName))
	var pbData map[string]interface{}
	exists := false

	if data, ok := result.PbEgressProfile[pbName].(map[string]interface{}); ok {
		pbData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found PB Egress Profile directly by ID: %s", pbName))
	} else {
		for apiName, p := range result.PbEgressProfile {
			profile, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := profile["name"].(string); ok && name == pbName {
				pbData = profile
				pbName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found PB Egress Profile with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("PB Egress Profile with ID '%s' not found in API response", pbName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", pbData["name"]))

	if enable, ok := pbData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if ipv4Permit, ok := pbData["ipv4_permit"].([]interface{}); ok && len(ipv4Permit) > 0 {
		var filters []verityPacketBrokerFilterModel
		for _, f := range ipv4Permit {
			filter, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPacketBrokerFilterModel{}
			if enable, ok := filter["enable"].(bool); ok {
				filterModel.Enable = types.BoolValue(enable)
			} else {
				filterModel.Enable = types.BoolNull()
			}
			if filterVal, ok := filter["filter"].(string); ok {
				filterModel.Filter = types.StringValue(filterVal)
			} else {
				filterModel.Filter = types.StringNull()
			}
			if refType, ok := filter["filter_ref_type_"].(string); ok {
				filterModel.FilterRefType = types.StringValue(refType)
			} else {
				filterModel.FilterRefType = types.StringNull()
			}
			if index, ok := filter["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else {
					filterModel.Index = types.Int64Null()
				}
			} else {
				filterModel.Index = types.Int64Null()
			}
			filters = append(filters, filterModel)
		}
		state.Ipv4Permit = filters
	} else {
		state.Ipv4Permit = nil
	}

	if ipv4Deny, ok := pbData["ipv4_deny"].([]interface{}); ok && len(ipv4Deny) > 0 {
		var filters []verityPacketBrokerFilterModel
		for _, f := range ipv4Deny {
			filter, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPacketBrokerFilterModel{}
			if enable, ok := filter["enable"].(bool); ok {
				filterModel.Enable = types.BoolValue(enable)
			} else {
				filterModel.Enable = types.BoolNull()
			}
			if filterVal, ok := filter["filter"].(string); ok {
				filterModel.Filter = types.StringValue(filterVal)
			} else {
				filterModel.Filter = types.StringNull()
			}
			if refType, ok := filter["filter_ref_type_"].(string); ok {
				filterModel.FilterRefType = types.StringValue(refType)
			} else {
				filterModel.FilterRefType = types.StringNull()
			}
			if index, ok := filter["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else {
					filterModel.Index = types.Int64Null()
				}
			} else {
				filterModel.Index = types.Int64Null()
			}
			filters = append(filters, filterModel)
		}
		state.Ipv4Deny = filters
	} else {
		state.Ipv4Deny = nil
	}

	if ipv6Permit, ok := pbData["ipv6_permit"].([]interface{}); ok && len(ipv6Permit) > 0 {
		var filters []verityPacketBrokerFilterModel
		for _, f := range ipv6Permit {
			filter, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPacketBrokerFilterModel{}
			if enable, ok := filter["enable"].(bool); ok {
				filterModel.Enable = types.BoolValue(enable)
			} else {
				filterModel.Enable = types.BoolNull()
			}
			if filterVal, ok := filter["filter"].(string); ok {
				filterModel.Filter = types.StringValue(filterVal)
			} else {
				filterModel.Filter = types.StringNull()
			}
			if refType, ok := filter["filter_ref_type_"].(string); ok {
				filterModel.FilterRefType = types.StringValue(refType)
			} else {
				filterModel.FilterRefType = types.StringNull()
			}
			if index, ok := filter["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else {
					filterModel.Index = types.Int64Null()
				}
			} else {
				filterModel.Index = types.Int64Null()
			}
			filters = append(filters, filterModel)
		}
		state.Ipv6Permit = filters
	} else {
		state.Ipv6Permit = nil
	}

	if ipv6Deny, ok := pbData["ipv6_deny"].([]interface{}); ok && len(ipv6Deny) > 0 {
		var filters []verityPacketBrokerFilterModel
		for _, f := range ipv6Deny {
			filter, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPacketBrokerFilterModel{}
			if enable, ok := filter["enable"].(bool); ok {
				filterModel.Enable = types.BoolValue(enable)
			} else {
				filterModel.Enable = types.BoolNull()
			}
			if filterVal, ok := filter["filter"].(string); ok {
				filterModel.Filter = types.StringValue(filterVal)
			} else {
				filterModel.Filter = types.StringNull()
			}
			if refType, ok := filter["filter_ref_type_"].(string); ok {
				filterModel.FilterRefType = types.StringValue(refType)
			} else {
				filterModel.FilterRefType = types.StringNull()
			}
			if index, ok := filter["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					filterModel.Index = types.Int64Value(int64(intVal))
				} else {
					filterModel.Index = types.Int64Null()
				}
			} else {
				filterModel.Index = types.Int64Null()
			}
			filters = append(filters, filterModel)
		}
		state.Ipv6Deny = filters
	} else {
		state.Ipv6Deny = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPacketBrokerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPacketBrokerResourceModel

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
	pbProps := openapi.PacketbrokerPutRequestPbEgressProfileValue{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		pbProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		pbProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !r.equalFilterArrays(plan.Ipv4Permit, state.Ipv4Permit) {
		ipv4Permit := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, len(plan.Ipv4Permit))
		for i, filter := range plan.Ipv4Permit {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv4Permit[i] = filterItem
		}
		pbProps.Ipv4Permit = ipv4Permit
		hasChanges = true
	}

	if !r.equalFilterArrays(plan.Ipv4Deny, state.Ipv4Deny) {
		ipv4Deny := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, len(plan.Ipv4Deny))
		for i, filter := range plan.Ipv4Deny {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv4Deny[i] = filterItem
		}
		pbProps.Ipv4Deny = ipv4Deny
		hasChanges = true
	}

	if !r.equalFilterArrays(plan.Ipv6Permit, state.Ipv6Permit) {
		ipv6Permit := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, len(plan.Ipv6Permit))
		for i, filter := range plan.Ipv6Permit {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv6Permit[i] = filterItem
		}
		pbProps.Ipv6Permit = ipv6Permit
		hasChanges = true
	}

	if !r.equalFilterArrays(plan.Ipv6Deny, state.Ipv6Deny) {
		ipv6Deny := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, len(plan.Ipv6Deny))
		for i, filter := range plan.Ipv6Deny {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}
			if !filter.Enable.IsNull() {
				filterItem.Enable = openapi.PtrBool(filter.Enable.ValueBool())
			}
			if !filter.Filter.IsNull() {
				filterItem.Filter = openapi.PtrString(filter.Filter.ValueString())
			}
			if !filter.FilterRefType.IsNull() {
				filterItem.FilterRefType = openapi.PtrString(filter.FilterRefType.ValueString())
			}
			if !filter.Index.IsNull() {
				filterItem.Index = openapi.PtrInt32(int32(filter.Index.ValueInt64()))
			}
			ipv6Deny[i] = filterItem
		}
		pbProps.Ipv6Deny = ipv6Deny
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "packet_broker", name, pbProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for PB Egress Profile update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update PB Egress Profile %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("PB Egress Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_egress_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityPacketBrokerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPacketBrokerResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "packet_broker", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for PB Egress Profile deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete PB Egress Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("PB Egress Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pb_egress_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityPacketBrokerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *verityPacketBrokerResource) equalFilterArrays(a, b []verityPacketBrokerFilterModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Enable.Equal(b[i].Enable) ||
			!a[i].Filter.Equal(b[i].Filter) ||
			!a[i].FilterRefType.Equal(b[i].FilterRefType) ||
			!a[i].Index.Equal(b[i].Index) {
			return false
		}
	}
	return true
}
