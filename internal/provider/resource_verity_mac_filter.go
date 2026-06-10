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
	_ resource.Resource                = &verityMacFilterResource{}
	_ resource.ResourceWithConfigure   = &verityMacFilterResource{}
	_ resource.ResourceWithImportState = &verityMacFilterResource{}
	_ resource.ResourceWithModifyPlan  = &verityMacFilterResource{}
)

const macFilterResourceType = "macfilters"

func NewVerityMacFilterResource() resource.Resource {
	return &verityMacFilterResource{}
}

type verityMacFilterResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityMacFilterResourceModel struct {
	Name    types.String                  `tfsdk:"name"`
	Enable  types.Bool                    `tfsdk:"enable"`
	Type    types.String                  `tfsdk:"type"`
	Filters []verityMacFilterFiltersModel `tfsdk:"filters"`
}

type verityMacFilterFiltersModel struct {
	FilterNumMac    types.String `tfsdk:"filter_num_mac"`
	FilterNumMask   types.String `tfsdk:"filter_num_mask"`
	FilterNumEnable types.Bool   `tfsdk:"filter_num_enable"`
	Index           types.Int64  `tfsdk:"index"`
}

func (m verityMacFilterFiltersModel) GetIndex() types.Int64 {
	return m.Index
}

func (r *verityMacFilterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mac_filter"
}

func (r *verityMacFilterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityMacFilterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity MAC Filter.",
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
			"type": schema.StringAttribute{
				Description: "Black vs White MAC Filter",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"filters": schema.ListNestedBlock{
				Description: "List of MAC filter entries",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"filter_num_mac": schema.StringAttribute{
							Description: "MAC address descriptor including colons example 01:23:45:67:9a:ab. and * notation accepted example 12:*",
							Optional:    true,
							Computed:    true,
						},
						"filter_num_mask": schema.StringAttribute{
							Description: "Hexidecimal mask including colons example ff:ff:fe:00:00:00. /n and * notation accepted example /16 or 12:*",
							Optional:    true,
							Computed:    true,
						},
						"filter_num_enable": schema.BoolAttribute{
							Description: "Enable of this MAC Filter ",
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

func (r *verityMacFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityMacFilterResourceModel
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
	macFilterProps := &openapi.MacfiltersPutRequestMacFilterValue{
		Name: openapi.PtrString(name),
	}

	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Type", APIField: &macFilterProps.Type, TFValue: plan.Type},
	})

	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &macFilterProps.Enable, TFValue: plan.Enable},
	})

	if len(plan.Filters) > 0 {
		filters := make([]openapi.MacfiltersPutRequestMacFilterValueFiltersInner, len(plan.Filters))
		for i, item := range plan.Filters {
			filter := openapi.MacfiltersPutRequestMacFilterValueFiltersInner{}

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "FilterNumMac", APIField: &filter.FilterNumMac, TFValue: item.FilterNumMac},
				{FieldName: "FilterNumMask", APIField: &filter.FilterNumMask, TFValue: item.FilterNumMask},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "FilterNumEnable", APIField: &filter.FilterNumEnable, TFValue: item.FilterNumEnable},
			})

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: item.Index},
			})

			filters[i] = filter
		}
		macFilterProps.Filters = filters
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "mac_filter", name, *macFilterProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("MAC Filter %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "mac_filters")

	var minState verityMacFilterResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if macFilterData, exists := bulkMgr.GetResourceResponse("mac_filter", name); exists {
			state := populateMacFilterState(ctx, minState, macFilterData, r.provCtx.mode)
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

func (r *verityMacFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityMacFilterResourceModel
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

	macFilterName := state.Name.ValueString()

	if r.bulkOpsMgr != nil {
		if macFilterData, exists := r.bulkOpsMgr.GetResourceResponse("mac_filter", macFilterName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached MAC filter data for %s from recent operation", macFilterName))
			state = populateMacFilterState(ctx, state, macFilterData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("mac_filter") {
		tflog.Info(ctx, fmt.Sprintf("Skipping MAC Filter %s verification – trusting recent successful API operation", macFilterName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching MAC Filters for verification of %s", macFilterName))

	type MacFiltersResponse struct {
		MacFilter map[string]interface{} `json:"mac_filter"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "mac_filters", macFilterName,
		func() (MacFiltersResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch MAC Filters")
			respAPI, err := r.client.MACFiltersAPI.MacfiltersGet(ctx).Execute()
			if err != nil {
				return MacFiltersResponse{}, fmt.Errorf("error reading MAC Filters: %v", err)
			}
			defer respAPI.Body.Close()

			var res MacFiltersResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return MacFiltersResponse{}, fmt.Errorf("failed to decode MAC Filters response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d MAC Filters", len(res.MacFilter)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read MAC Filter %s", macFilterName))...,
		)
		return
	}

	macFilterData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.MacFilter,
		macFilterName,
		func(data interface{}) (string, bool) {
			if macFilter, ok := data.(map[string]interface{}); ok {
				if name, ok := macFilter["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("MAC Filter with name '%s' not found in API response", macFilterName))
		resp.State.RemoveResource(ctx)
		return
	}

	macFilterMap, ok := macFilterData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid MAC Filter Data",
			fmt.Sprintf("MAC Filter data is not in expected format for %s", macFilterName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found MAC Filter '%s' under API key '%s'", macFilterName, actualAPIName))

	state = populateMacFilterState(ctx, state, macFilterMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityMacFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityMacFilterResourceModel

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
	macFilterProps := openapi.MacfiltersPutRequestMacFilterValue{}
	hasChanges := false

	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { macFilterProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Type, state.Type, func(v *string) { macFilterProps.Type = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { macFilterProps.Enable = v }, &hasChanges)

	filtersHandler := utils.IndexedItemHandler[verityMacFilterFiltersModel, openapi.MacfiltersPutRequestMacFilterValueFiltersInner]{
		CreateNew: func(planItem verityMacFilterFiltersModel) openapi.MacfiltersPutRequestMacFilterValueFiltersInner {
			filter := openapi.MacfiltersPutRequestMacFilterValueFiltersInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "FilterNumMac", APIField: &filter.FilterNumMac, TFValue: planItem.FilterNumMac},
				{FieldName: "FilterNumMask", APIField: &filter.FilterNumMask, TFValue: planItem.FilterNumMask},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "FilterNumEnable", APIField: &filter.FilterNumEnable, TFValue: planItem.FilterNumEnable},
			})

			return filter
		},
		UpdateExisting: func(planItem verityMacFilterFiltersModel, stateItem verityMacFilterFiltersModel) (openapi.MacfiltersPutRequestMacFilterValueFiltersInner, bool) {
			filter := openapi.MacfiltersPutRequestMacFilterValueFiltersInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			utils.CompareAndSetStringField(planItem.FilterNumMac, stateItem.FilterNumMac, func(v *string) { filter.FilterNumMac = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.FilterNumMask, stateItem.FilterNumMask, func(v *string) { filter.FilterNumMask = v }, &fieldChanged)
			utils.CompareAndSetBoolField(planItem.FilterNumEnable, stateItem.FilterNumEnable, func(v *bool) { filter.FilterNumEnable = v }, &fieldChanged)

			return filter, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.MacfiltersPutRequestMacFilterValueFiltersInner {
			filter := openapi.MacfiltersPutRequestMacFilterValueFiltersInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &filter.Index, TFValue: types.Int64Value(index)},
			})
			return filter
		},
	}

	changedFilters, filtersChanged := utils.ProcessIndexedArrayUpdates(plan.Filters, state.Filters, filtersHandler)
	if filtersChanged {
		macFilterProps.Filters = changedFilters
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "mac_filter", name, macFilterProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("MAC Filter %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "mac_filters")

	var minState verityMacFilterResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if macFilterData, exists := bulkMgr.GetResourceResponse("mac_filter", name); exists {
			newState := populateMacFilterState(ctx, minState, macFilterData, r.provCtx.mode)
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

func (r *verityMacFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityMacFilterResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "mac_filter", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("MAC Filter %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "mac_filters")
	resp.State.RemoveResource(ctx)
}

func (r *verityMacFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateMacFilterState(ctx context.Context, state verityMacFilterResourceModel, data map[string]interface{}, mode string) verityMacFilterResourceModel {
	const resourceType = macFilterResourceType

	state.Name = utils.MapStringFromAPI(data["name"])
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.Type = utils.MapStringWithMode(data, "type", resourceType, mode)

	if utils.FieldAppliesToMode(resourceType, "filters", mode) {
		if filtersData, ok := data["filters"].([]interface{}); ok && len(filtersData) > 0 {
			var filters []verityMacFilterFiltersModel
			for _, item := range filtersData {
				filterItem, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				filterModel := verityMacFilterFiltersModel{
					FilterNumMac:    utils.MapStringWithModeNested(filterItem, "filter_num_mac", resourceType, "filters.filter_num_mac", mode),
					FilterNumMask:   utils.MapStringWithModeNested(filterItem, "filter_num_mask", resourceType, "filters.filter_num_mask", mode),
					FilterNumEnable: utils.MapBoolWithModeNested(filterItem, "filter_num_enable", resourceType, "filters.filter_num_enable", mode),
					Index:           utils.MapInt64WithModeNested(filterItem, "index", resourceType, "filters.index", mode),
				}
				filters = append(filters, filterModel)
			}
			state.Filters = filters
		} else {
			state.Filters = nil
		}
	} else {
		state.Filters = nil
	}

	return state
}

func (r *verityMacFilterResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityMacFilterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const resourceType = macFilterResourceType
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
		BlockName:    "filters",
		ItemCount:    len(plan.Filters),
		StringFields: []string{"filter_num_mac", "filter_num_mask"},
		BoolFields:   []string{"filter_num_enable"},
		Int64Fields:  []string{"index"},
	})
}
