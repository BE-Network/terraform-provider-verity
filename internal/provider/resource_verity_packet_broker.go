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
)

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
	pbProps := &openapi.PacketbrokerPutRequestPortAclValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &pbProps.Enable, TFValue: plan.Enable},
	})

	if len(plan.Ipv4Permit) > 0 {
		ipv4Permit := make([]openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner, len(plan.Ipv4Permit))
		for i, filter := range plan.Ipv4Permit {
			filterItem := openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{}

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
		ipv4Deny := make([]openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner, len(plan.Ipv4Deny))
		for i, filter := range plan.Ipv4Deny {
			filterItem := openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{}

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
		ipv6Permit := make([]openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner, len(plan.Ipv6Permit))
		for i, filter := range plan.Ipv6Permit {
			filterItem := openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{}

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
		ipv6Deny := make([]openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner, len(plan.Ipv6Deny))
		for i, filter := range plan.Ipv6Deny {
			filterItem := openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{}

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

	state.Name = utils.MapStringFromAPI(pbMap["name"])

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(pbMap[apiKey])
	}

	// Helper function to parse filter arrays
	parseFilters := func(apiFilters []interface{}) []verityPacketBrokerFilterModel {
		var filters []verityPacketBrokerFilterModel
		for _, f := range apiFilters {
			filter, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			filterModel := verityPacketBrokerFilterModel{
				Enable:        utils.MapBoolFromAPI(filter["enable"]),
				Filter:        utils.MapStringFromAPI(filter["filter"]),
				FilterRefType: utils.MapStringFromAPI(filter["filter_ref_type_"]),
				Index:         utils.MapInt64FromAPI(filter["index"]),
			}
			filters = append(filters, filterModel)
		}
		return filters
	}

	// Handle filter arrays
	if ipv4Permit, ok := pbMap["ipv4_permit"].([]interface{}); ok && len(ipv4Permit) > 0 {
		state.Ipv4Permit = parseFilters(ipv4Permit)
	} else {
		state.Ipv4Permit = nil
	}

	if ipv4Deny, ok := pbMap["ipv4_deny"].([]interface{}); ok && len(ipv4Deny) > 0 {
		state.Ipv4Deny = parseFilters(ipv4Deny)
	} else {
		state.Ipv4Deny = nil
	}

	if ipv6Permit, ok := pbMap["ipv6_permit"].([]interface{}); ok && len(ipv6Permit) > 0 {
		state.Ipv6Permit = parseFilters(ipv6Permit)
	} else {
		state.Ipv6Permit = nil
	}

	if ipv6Deny, ok := pbMap["ipv6_deny"].([]interface{}); ok && len(ipv6Deny) > 0 {
		state.Ipv6Deny = parseFilters(ipv6Deny)
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
	pbProps := openapi.PacketbrokerPutRequestPortAclValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { pbProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { pbProps.Enable = v }, &hasChanges)

	// Handle IPv4 Permit using consolidated handler
	changedIpv4Permit, ipv4PermitChanged := utils.ProcessIndexedArrayUpdates(plan.Ipv4Permit, state.Ipv4Permit,
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{}
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner {
				return openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{
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
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{}
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner {
				return openapi.PacketbrokerPutRequestPortAclValueIpv4PermitInner{
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
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{}
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner {
				return openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{
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
		utils.IndexedItemHandler[verityPacketBrokerFilterModel, openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner]{
			CreateNew: func(planItem verityPacketBrokerFilterModel) openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner {
				newFilter := openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{}

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
			UpdateExisting: func(planItem verityPacketBrokerFilterModel, stateItem verityPacketBrokerFilterModel) (openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner, bool) {
				updateFilter := openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{}
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
			CreateDeleted: func(index int64) openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner {
				return openapi.PacketbrokerPutRequestPortAclValueIpv6PermitInner{
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
