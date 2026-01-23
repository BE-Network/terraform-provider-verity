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
	_ resource.Resource                = &verityPacketBrokerResource{}
	_ resource.ResourceWithConfigure   = &verityPacketBrokerResource{}
	_ resource.ResourceWithImportState = &verityPacketBrokerResource{}
	_ resource.ResourceWithModifyPlan  = &verityPacketBrokerResource{}
)

const packetBrokerResourceType = "packetbroker"

func NewVerityPacketBrokerResource() resource.Resource {
	return &verityPacketBrokerResource{}
}

type verityPacketBrokerResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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

func (f verityPacketBrokerFilterModel) GetIndex() types.Int64 {
	return f.Index
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
		Description: "Manages a Verity Packet Broker",
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
				Description: "IPv4 Permit filters",
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
				Description: "IPv4 Deny filters",
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
				Description: "IPv6 Permit filters",
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
				Description: "IPv6 Deny filters",
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
	pbProps := &openapi.PacketbrokerPutRequestPbEgressProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &pbProps.Enable, TFValue: plan.Enable},
	})

	if len(plan.Ipv4Permit) > 0 {
		ipv4Permit := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, len(plan.Ipv4Permit))
		for i, filter := range plan.Ipv4Permit {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filterItem.Enable, TFValue: filter.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filterItem.Filter, TFValue: filter.Filter},
				{FieldName: "FilterRefType", APIField: &filterItem.FilterRefType, TFValue: filter.FilterRefType},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filterItem.Index, TFValue: filter.Index},
			})

			ipv4Permit[i] = filterItem
		}
		pbProps.Ipv4Permit = ipv4Permit
	}

	if len(plan.Ipv4Deny) > 0 {
		ipv4Deny := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, len(plan.Ipv4Deny))
		for i, filter := range plan.Ipv4Deny {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filterItem.Enable, TFValue: filter.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filterItem.Filter, TFValue: filter.Filter},
				{FieldName: "FilterRefType", APIField: &filterItem.FilterRefType, TFValue: filter.FilterRefType},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filterItem.Index, TFValue: filter.Index},
			})

			ipv4Deny[i] = filterItem
		}
		pbProps.Ipv4Deny = ipv4Deny
	}

	if len(plan.Ipv6Permit) > 0 {
		ipv6Permit := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, len(plan.Ipv6Permit))
		for i, filter := range plan.Ipv6Permit {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filterItem.Enable, TFValue: filter.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filterItem.Filter, TFValue: filter.Filter},
				{FieldName: "FilterRefType", APIField: &filterItem.FilterRefType, TFValue: filter.FilterRefType},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filterItem.Index, TFValue: filter.Index},
			})

			ipv6Permit[i] = filterItem
		}
		pbProps.Ipv6Permit = ipv6Permit
	}

	if len(plan.Ipv6Deny) > 0 {
		ipv6Deny := make([]openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, len(plan.Ipv6Deny))
		for i, filter := range plan.Ipv6Deny {
			filterItem := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &filterItem.Enable, TFValue: filter.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Filter", APIField: &filterItem.Filter, TFValue: filter.Filter},
				{FieldName: "FilterRefType", APIField: &filterItem.FilterRefType, TFValue: filter.FilterRefType},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filterItem.Index, TFValue: filter.Index},
			})

			ipv6Deny[i] = filterItem
		}
		pbProps.Ipv6Deny = ipv6Deny
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "packet_broker", name, *pbProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Broker %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_brokers")

	var minState verityPacketBrokerResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pbData, exists := bulkMgr.GetResourceResponse("packet_broker", name); exists {
			state := populatePacketBrokerState(ctx, minState, pbData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if pbData, exists := r.bulkOpsMgr.GetResourceResponse("packet_broker", pbName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached packet broker data for %s from recent operation", pbName))
			state = populatePacketBrokerState(ctx, state, pbData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("packet_broker") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Packet Broker %s verification – trusting recent successful API operation", pbName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Packet Broker for verification of %s", pbName))

	type PacketBrokerResponse struct {
		PacketBroker map[string]interface{} `json:"pb_egress_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "packet_brokers", pbName,
		func() (PacketBrokerResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Packet Broker")
			respAPI, err := r.client.PacketBrokerAPI.PacketbrokerGet(ctx).Execute()
			if err != nil {
				return PacketBrokerResponse{}, fmt.Errorf("error reading Packet Broker: %v", err)
			}
			defer respAPI.Body.Close()

			var res PacketBrokerResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PacketBrokerResponse{}, fmt.Errorf("failed to decode Packet Broker response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Packet Brokers", len(res.PacketBroker)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Packet Broker %s", pbName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Packet Broker with name: %s", pbName))

	pbData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.PacketBroker,
		pbName,
		func(data interface{}) (string, bool) {
			if profile, ok := data.(map[string]interface{}); ok {
				if name, ok := profile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Packet Broker with name '%s' not found in API response", pbName))
		resp.State.RemoveResource(ctx)
		return
	}

	pbMap, ok := pbData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Packet Broker Data",
			fmt.Sprintf("Packet Broker data is not in expected format for %s", pbName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Packet Broker '%s' under API key '%s'", pbName, actualAPIName))

	state = populatePacketBrokerState(ctx, state, pbMap, r.provCtx.mode)
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

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { pbProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { pbProps.Enable = v }, &hasChanges)

	// Handle IPv4 Permit using consolidated handler
	changedIpv4Permit, ipv4PermitChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv4Permit, state.Ipv4Permit,
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateFilter.Enable = v }, &fieldChanged)

				// Handle filter and filter_ref_type_ using multiple ref types supported pattern
				if !utils.HandleMultipleRefTypesSupported(
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner {
				return openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv4PermitChanged {
		pbProps.Ipv4Permit = changedIpv4Permit
		hasChanges = true
	}

	// Handle IPv4 Deny
	changedIpv4Deny, ipv4DenyChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv4Deny, state.Ipv4Deny,
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateFilter.Enable = v }, &fieldChanged)

				// Handle filter and filter_ref_type_ using multiple ref types supported pattern
				if !utils.HandleMultipleRefTypesSupported(
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner {
				return openapi.PacketbrokerPutRequestPbEgressProfileValueIpv4PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv4DenyChanged {
		pbProps.Ipv4Deny = changedIpv4Deny
		hasChanges = true
	}

	// Handle IPv6 Permit
	changedIpv6Permit, ipv6PermitChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv6Permit, state.Ipv6Permit,
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}
				fieldChanged := false

				// Handle boolean field changes

				// Handle filter and filter_ref_type_ using multiple ref types supported pattern
				if !utils.HandleMultipleRefTypesSupported(
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner {
				return openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv6PermitChanged {
		pbProps.Ipv6Permit = changedIpv6Permit
		hasChanges = true
	}

	// Handle IPv6 Deny
	changedIpv6Deny, ipv6DenyChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv6Deny, state.Ipv6Deny,
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateFilter.Enable = v }, &fieldChanged)

				// Handle filter and filter_ref_type_ using multiple ref types supported pattern
				if !utils.HandleMultipleRefTypesSupported(
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner {
				return openapi.PacketbrokerPutRequestPbEgressProfileValueIpv6PermitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ipv6DenyChanged {
		pbProps.Ipv6Deny = changedIpv6Deny
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "packet_broker", name, pbProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Broker %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_brokers")

	var minState verityPacketBrokerResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pbData, exists := bulkMgr.GetResourceResponse("packet_broker", name); exists {
			newState := populatePacketBrokerState(ctx, minState, pbData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "packet_broker", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Broker %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_brokers")
	resp.State.RemoveResource(ctx)
}

func (r *verityPacketBrokerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populatePacketBrokerState(ctx context.Context, state verityPacketBrokerResourceModel, data map[string]interface{}, mode string) verityPacketBrokerResourceModel {
	const resourceType = packetBrokerResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Helper function to parse filter arrays with mode awareness
	parseFilters := func(apiFilters []interface{}, blockName string) []verityPacketBrokerFilterModel {
		var filters []verityPacketBrokerFilterModel
		for _, f := range apiFilters {
			filter, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPacketBrokerFilterModel{
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

func (r *verityPacketBrokerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityPacketBrokerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = packetBrokerResourceType
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
}
