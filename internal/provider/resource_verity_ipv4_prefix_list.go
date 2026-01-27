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
	_ resource.Resource                = &verityIpv4PrefixListResource{}
	_ resource.ResourceWithConfigure   = &verityIpv4PrefixListResource{}
	_ resource.ResourceWithImportState = &verityIpv4PrefixListResource{}
	_ resource.ResourceWithModifyPlan  = &verityIpv4PrefixListResource{}
)

const ipv4PrefixListResourceType = "ipv4prefixlists"

func NewVerityIpv4PrefixListResource() resource.Resource {
	return &verityIpv4PrefixListResource{}
}

type verityIpv4PrefixListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityIpv4PrefixListResourceModel struct {
	Name             types.String                                `tfsdk:"name"`
	Enable           types.Bool                                  `tfsdk:"enable"`
	Lists            []verityIpv4PrefixListListsModel            `tfsdk:"lists"`
	ObjectProperties []verityIpv4PrefixListObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityIpv4PrefixListListsModel struct {
	Enable                types.Bool   `tfsdk:"enable"`
	PermitDeny            types.String `tfsdk:"permit_deny"`
	Ipv4Prefix            types.String `tfsdk:"ipv4_prefix"`
	GreaterThanEqualValue types.Int64  `tfsdk:"greater_than_equal_value"`
	LessThanEqualValue    types.Int64  `tfsdk:"less_than_equal_value"`
	Index                 types.Int64  `tfsdk:"index"`
}

func (l verityIpv4PrefixListListsModel) GetIndex() types.Int64 {
	return l.Index
}

type verityIpv4PrefixListObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityIpv4PrefixListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipv4_prefix_list"
}

func (r *verityIpv4PrefixListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityIpv4PrefixListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity IPv4 Prefix List",
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
			"lists": schema.ListNestedBlock{
				Description: "List of IPv4 Prefix List entries",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable of this IPv4 Prefix List",
							Optional:    true,
							Computed:    true,
						},
						"permit_deny": schema.StringAttribute{
							Description: "Action upon match of Community Strings.",
							Optional:    true,
							Computed:    true,
						},
						"ipv4_prefix": schema.StringAttribute{
							Description: "IPv4 address and subnet to match against",
							Optional:    true,
							Computed:    true,
						},
						"greater_than_equal_value": schema.Int64Attribute{
							Description: "Match IP routes with a subnet mask greater than or equal to the value indicated (maximum: 32)",
							Optional:    true,
							Computed:    true,
						},
						"less_than_equal_value": schema.Int64Attribute{
							Description: "Match IP routes with a subnet mask less than or equal to the value indicated (maximum: 32)",
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
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the IPv4 Prefix List",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"notes": schema.StringAttribute{
							Description: "User Notes.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityIpv4PrefixListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityIpv4PrefixListResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityIpv4PrefixListResourceModel
	diags = req.Config.Get(ctx, &config)
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
	ipv4PrefixListProps := &openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ipv4PrefixListProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objectProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objectProps.Notes},
		})
		ipv4PrefixListProps.ObjectProperties = &objectProps
	}

	// Parse HCL to detect explicitly configured attributes
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_ipv4_prefix_list", name)

	// Handle lists
	if len(plan.Lists) > 0 {
		listsConfigMap := utils.BuildIndexedConfigMap(config.Lists)

		var lists []openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner
		for _, listItem := range plan.Lists {
			item := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &item.Enable, TFValue: listItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PermitDeny", APIField: &item.PermitDeny, TFValue: listItem.PermitDeny},
				{FieldName: "Ipv4Prefix", APIField: &item.Ipv4Prefix, TFValue: listItem.Ipv4Prefix},
			})

			// Get per-block configured info for nullable Int64 fields
			itemIndex := listItem.Index.ValueInt64()
			configItem := listItem // fallback to plan item
			if cfgItem, ok := listsConfigMap[itemIndex]; ok {
				configItem = cfgItem
			}
			cfg := &utils.IndexedBlockNullableFieldConfig{
				BlockType:       "lists",
				BlockIndex:      itemIndex,
				ConfiguredAttrs: configuredAttrs,
			}
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "GreaterThanEqualValue", APIField: &item.GreaterThanEqualValue, TFValue: configItem.GreaterThanEqualValue, IsConfigured: cfg.IsFieldConfigured("greater_than_equal_value")},
				{FieldName: "LessThanEqualValue", APIField: &item.LessThanEqualValue, TFValue: configItem.LessThanEqualValue, IsConfigured: cfg.IsFieldConfigured("less_than_equal_value")},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: listItem.Index},
			})
			lists = append(lists, item)
		}
		ipv4PrefixListProps.Lists = lists
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "ipv4_prefix_list", name, *ipv4PrefixListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv4 Prefix List %s create operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_prefix_lists")

	var minState verityIpv4PrefixListResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ipv4PrefixListData, exists := bulkMgr.GetResourceResponse("ipv4_prefix_list", name); exists {
			state := populateIpv4PrefixListState(ctx, minState, ipv4PrefixListData, r.provCtx.mode)
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

func (r *verityIpv4PrefixListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityIpv4PrefixListResourceModel
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
		if ipv4PrefixListData, exists := r.bulkOpsMgr.GetResourceResponse("ipv4_prefix_list", name); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached IPv4 Prefix List data for %s from recent operation", name))
			state = populateIpv4PrefixListState(ctx, state, ipv4PrefixListData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ipv4_prefix_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv4 prefix list %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv4 prefix lists for verification of %s", name))

	type Ipv4PrefixListResponse struct {
		Ipv4PrefixList map[string]map[string]interface{} `json:"ipv4_prefix_list"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "ipv4_prefix_lists", name,
		func() (Ipv4PrefixListResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv4 prefix lists")
			resp, err := r.client.IPv4PrefixListsAPI.Ipv4prefixlistsGet(ctx).Execute()
			if err != nil {
				return Ipv4PrefixListResponse{}, fmt.Errorf("error reading IPv4 prefix lists: %v", err)
			}
			defer resp.Body.Close()

			var result Ipv4PrefixListResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return Ipv4PrefixListResponse{}, fmt.Errorf("error decoding IPv4 prefix list response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv4 prefix lists", len(result.Ipv4PrefixList)))
			return result, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv4 Prefix List %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv4 prefix list with name: %s", name))

	ipv4PrefixListData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Ipv4PrefixList,
		name,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv4 Prefix List with name '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	ipv4PrefixListMap, ok := ipv4PrefixListData, true
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid IPv4 Prefix List Data",
			fmt.Sprintf("IPv4 Prefix List data is not in expected format for %s", name),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found IPv4 prefix list '%s' under API key '%s'", name, actualAPIName))

	state = populateIpv4PrefixListState(ctx, state, ipv4PrefixListMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityIpv4PrefixListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityIpv4PrefixListResourceModel

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
	ipv4PrefixListProps := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue{}
	hasChanges := false

	//handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ipv4PrefixListProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ipv4PrefixListProps.Enable = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objectProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objectProps.Notes},
		}, &objPropsChanged)

		if objPropsChanged {
			ipv4PrefixListProps.ObjectProperties = &objectProps
			hasChanges = true
		}
	}

	// Handle lists
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_ipv4_prefix_list", name)
	var config verityIpv4PrefixListResourceModel
	req.Config.Get(ctx, &config)
	listsConfigMap := utils.BuildIndexedConfigMap(config.Lists)

	listsHandler := utils.IndexedItemHandler[verityIpv4PrefixListListsModel, openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner]{
		CreateNew: func(planItem verityIpv4PrefixListListsModel) openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner {
			newItem := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &newItem.Enable, TFValue: planItem.Enable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PermitDeny", APIField: &newItem.PermitDeny, TFValue: planItem.PermitDeny},
				{FieldName: "Ipv4Prefix", APIField: &newItem.Ipv4Prefix, TFValue: planItem.Ipv4Prefix},
			})

			// Get per-block configured info for nullable Int64 fields
			itemIndex := planItem.Index.ValueInt64()
			configItem := planItem // fallback to plan item
			if cfgItem, ok := listsConfigMap[itemIndex]; ok {
				configItem = cfgItem
			}
			cfg := &utils.IndexedBlockNullableFieldConfig{
				BlockType:       "lists",
				BlockIndex:      itemIndex,
				ConfiguredAttrs: configuredAttrs,
			}
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "GreaterThanEqualValue", APIField: &newItem.GreaterThanEqualValue, TFValue: configItem.GreaterThanEqualValue, IsConfigured: cfg.IsFieldConfigured("greater_than_equal_value")},
				{FieldName: "LessThanEqualValue", APIField: &newItem.LessThanEqualValue, TFValue: configItem.LessThanEqualValue, IsConfigured: cfg.IsFieldConfigured("less_than_equal_value")},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &newItem.Index, TFValue: planItem.Index},
			})

			return newItem
		},
		UpdateExisting: func(planItem verityIpv4PrefixListListsModel, stateItem verityIpv4PrefixListListsModel) (openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner, bool) {
			updateItem := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{}
			fieldChanged := false

			// Handle boolean field changes
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateItem.Enable = v }, &fieldChanged)

			// Handle string field changes
			utils.CompareAndSetStringField(planItem.PermitDeny, stateItem.PermitDeny, func(v *string) { updateItem.PermitDeny = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Ipv4Prefix, stateItem.Ipv4Prefix, func(v *string) { updateItem.Ipv4Prefix = v }, &fieldChanged)

			// Handle nullable int64 field changes
			utils.CompareAndSetNullableInt64Field(planItem.GreaterThanEqualValue, stateItem.GreaterThanEqualValue, func(v *openapi.NullableInt32) { updateItem.GreaterThanEqualValue = *v }, &fieldChanged)
			utils.CompareAndSetNullableInt64Field(planItem.LessThanEqualValue, stateItem.LessThanEqualValue, func(v *openapi.NullableInt32) { updateItem.LessThanEqualValue = *v }, &fieldChanged)

			// Handle index field change
			utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateItem.Index = v }, &fieldChanged)

			return updateItem, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner {
			return openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(index)),
			}
		},
	}

	changedLists, listsChanged := utils.ProcessIndexedArrayUpdates(plan.Lists, state.Lists, listsHandler)
	if listsChanged {
		ipv4PrefixListProps.Lists = changedLists
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "ipv4_prefix_list", name, ipv4PrefixListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv4 Prefix List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_prefix_lists")

	var minState verityIpv4PrefixListResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ipv4PrefixListData, exists := bulkMgr.GetResourceResponse("ipv4_prefix_list", name); exists {
			newState := populateIpv4PrefixListState(ctx, minState, ipv4PrefixListData, r.provCtx.mode)
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

func (r *verityIpv4PrefixListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityIpv4PrefixListResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "ipv4_prefix_list", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv4 Prefix List %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_prefix_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv4PrefixListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateIpv4PrefixListState(ctx context.Context, state verityIpv4PrefixListResourceModel, data map[string]interface{}, mode string) verityIpv4PrefixListResourceModel {
	const resourceType = ipv4PrefixListResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityIpv4PrefixListObjectPropertiesModel{
				Notes: utils.MapStringWithModeNested(objProps, "notes", resourceType, "object_properties.notes", mode),
			}
			state.ObjectProperties = []verityIpv4PrefixListObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle lists block
	if utils.FieldAppliesToMode(resourceType, "lists", mode) {
		if listsData, ok := data["lists"].([]interface{}); ok && len(listsData) > 0 {
			var listsList []verityIpv4PrefixListListsModel

			for _, item := range listsData {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				listItem := verityIpv4PrefixListListsModel{
					Enable:                utils.MapBoolWithModeNested(itemMap, "enable", resourceType, "lists.enable", mode),
					PermitDeny:            utils.MapStringWithModeNested(itemMap, "permit_deny", resourceType, "lists.permit_deny", mode),
					Ipv4Prefix:            utils.MapStringWithModeNested(itemMap, "ipv4_prefix", resourceType, "lists.ipv4_prefix", mode),
					GreaterThanEqualValue: utils.MapInt64WithModeNested(itemMap, "greater_than_equal_value", resourceType, "lists.greater_than_equal_value", mode),
					LessThanEqualValue:    utils.MapInt64WithModeNested(itemMap, "less_than_equal_value", resourceType, "lists.less_than_equal_value", mode),
					Index:                 utils.MapInt64WithModeNested(itemMap, "index", resourceType, "lists.index", mode),
				}

				listsList = append(listsList, listItem)
			}

			state.Lists = listsList
		} else {
			state.Lists = nil
		}
	} else {
		state.Lists = nil
	}

	return state
}

func (r *verityIpv4PrefixListResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityIpv4PrefixListResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = ipv4PrefixListResourceType
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

	nullifier.NullifyNestedBlocks(
		"lists", "object_properties",
	)

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityIpv4PrefixListResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityIpv4PrefixListResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_ipv4_prefix_list", name)

	for i, configItem := range config.Lists {
		itemIndex := configItem.Index.ValueInt64()
		var stateItem *verityIpv4PrefixListListsModel
		for j := range state.Lists {
			if state.Lists[j].Index.ValueInt64() == itemIndex {
				stateItem = &state.Lists[j]
				break
			}
		}

		if stateItem != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "lists",
				BlockListPath:   "lists",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: itemIndex, AttrName: "greater_than_equal_value", ConfigVal: configItem.GreaterThanEqualValue, StateVal: stateItem.GreaterThanEqualValue},
					{BlockIndex: itemIndex, AttrName: "less_than_equal_value", ConfigVal: configItem.LessThanEqualValue, StateVal: stateItem.LessThanEqualValue},
				},
			})
		}
	}
}
