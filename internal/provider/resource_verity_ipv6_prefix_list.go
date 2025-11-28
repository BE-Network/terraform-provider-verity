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
	_ resource.Resource                = &verityIpv6PrefixListResource{}
	_ resource.ResourceWithConfigure   = &verityIpv6PrefixListResource{}
	_ resource.ResourceWithImportState = &verityIpv6PrefixListResource{}
)

func NewVerityIpv6PrefixListResource() resource.Resource {
	return &verityIpv6PrefixListResource{}
}

type verityIpv6PrefixListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityIpv6PrefixListResourceModel struct {
	Name             types.String                                `tfsdk:"name"`
	Enable           types.Bool                                  `tfsdk:"enable"`
	Lists            []verityIpv6PrefixListListsModel            `tfsdk:"lists"`
	ObjectProperties []verityIpv6PrefixListObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityIpv6PrefixListListsModel struct {
	Enable                types.Bool   `tfsdk:"enable"`
	PermitDeny            types.String `tfsdk:"permit_deny"`
	Ipv6Prefix            types.String `tfsdk:"ipv6_prefix"`
	GreaterThanEqualValue types.Int64  `tfsdk:"greater_than_equal_value"`
	LessThanEqualValue    types.Int64  `tfsdk:"less_than_equal_value"`
	Index                 types.Int64  `tfsdk:"index"`
}

func (l verityIpv6PrefixListListsModel) GetIndex() types.Int64 {
	return l.Index
}

type verityIpv6PrefixListObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityIpv6PrefixListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipv6_prefix_list"
}

func (r *verityIpv6PrefixListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityIpv6PrefixListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity IPv6 Prefix List",
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
			"lists": schema.ListNestedBlock{
				Description: "List of IPv6 Prefix List entries",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable of this IPv6 Prefix List",
							Optional:    true,
						},
						"permit_deny": schema.StringAttribute{
							Description: "Action upon match of Community Strings.",
							Optional:    true,
						},
						"ipv6_prefix": schema.StringAttribute{
							Description: "IPv6 address and subnet to match against",
							Optional:    true,
						},
						"greater_than_equal_value": schema.Int64Attribute{
							Description: "Match IP routes with a subnet mask greater than or equal to the value indicated (minimum: 1, maximum: 128)",
							Optional:    true,
						},
						"less_than_equal_value": schema.Int64Attribute{
							Description: "Match IP routes with a subnet mask less than or equal to the value indicated (minimum: 1, maximum: 128)",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the IPv6 Prefix List",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"notes": schema.StringAttribute{
							Description: "User Notes.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityIpv6PrefixListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityIpv6PrefixListResourceModel
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
	ipv6PrefixListProps := &openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ipv6PrefixListProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objectProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objectProps.Notes},
		})
		ipv6PrefixListProps.ObjectProperties = &objectProps
	}

	// Handle lists
	if len(plan.Lists) > 0 {
		var lists []openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner
		for _, listItem := range plan.Lists {
			item := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &item.Enable, TFValue: listItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PermitDeny", APIField: &item.PermitDeny, TFValue: listItem.PermitDeny},
				{FieldName: "Ipv6Prefix", APIField: &item.Ipv6Prefix, TFValue: listItem.Ipv6Prefix},
			})
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "GreaterThanEqualValue", APIField: &item.GreaterThanEqualValue, TFValue: listItem.GreaterThanEqualValue},
				{FieldName: "LessThanEqualValue", APIField: &item.LessThanEqualValue, TFValue: listItem.LessThanEqualValue},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: listItem.Index},
			})

			lists = append(lists, item)
		}
		ipv6PrefixListProps.Lists = lists
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "ipv6_prefix_list", name, *ipv6PrefixListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv6 Prefix List %s create operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_prefix_lists")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityIpv6PrefixListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityIpv6PrefixListResourceModel
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ipv6_prefix_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv6 prefix list %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv6 prefix lists for verification of %s", name))

	type Ipv6PrefixListResponse struct {
		Ipv6PrefixList map[string]map[string]interface{} `json:"ipv6_prefix_list"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "ipv6_prefix_lists", name,
		func() (Ipv6PrefixListResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv6 prefix lists")
			resp, err := r.client.IPv6PrefixListsAPI.Ipv6prefixlistsGet(ctx).Execute()
			if err != nil {
				return Ipv6PrefixListResponse{}, fmt.Errorf("error reading IPv6 prefix lists: %v", err)
			}
			defer resp.Body.Close()

			var result Ipv6PrefixListResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return Ipv6PrefixListResponse{}, fmt.Errorf("error decoding IPv6 prefix list response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv6 prefix lists", len(result.Ipv6PrefixList)))
			return result, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv6 Prefix List %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv6 prefix list with name: %s", name))

	ipv6PrefixListData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Ipv6PrefixList,
		name,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv6 Prefix List with name '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	ipv6PrefixListMap, ok := ipv6PrefixListData, true
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid IPv6 Prefix List Data",
			fmt.Sprintf("IPv6 Prefix List data is not in expected format for %s", name),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found IPv6 prefix list '%s' under API key '%s'", name, actualAPIName))

	state.Name = utils.MapStringFromAPI(ipv6PrefixListMap["name"])

	// Handle object properties
	if objectPropsData, ok := ipv6PrefixListMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityIpv6PrefixListObjectPropertiesModel{
			{Notes: utils.MapStringFromAPI(objectPropsData["notes"])},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ipv6PrefixListMap[apiKey])
	}

	// Handle lists
	if listsData, ok := ipv6PrefixListMap["lists"].([]interface{}); ok && len(listsData) > 0 {
		var lists []verityIpv6PrefixListListsModel

		for _, item := range listsData {
			listItem, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			listModel := verityIpv6PrefixListListsModel{
				Enable:                utils.MapBoolFromAPI(listItem["enable"]),
				PermitDeny:            utils.MapStringFromAPI(listItem["permit_deny"]),
				Ipv6Prefix:            utils.MapStringFromAPI(listItem["ipv6_prefix"]),
				GreaterThanEqualValue: utils.MapInt64FromAPI(listItem["greater_than_equal_value"]),
				LessThanEqualValue:    utils.MapInt64FromAPI(listItem["less_than_equal_value"]),
				Index:                 utils.MapInt64FromAPI(listItem["index"]),
			}

			lists = append(lists, listModel)
		}
		state.Lists = lists
	} else {
		state.Lists = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *verityIpv6PrefixListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityIpv6PrefixListResourceModel

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
	ipv6PrefixListProps := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ipv6PrefixListProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ipv6PrefixListProps.Enable = v }, &hasChanges)

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
			ipv6PrefixListProps.ObjectProperties = &objectProps
			hasChanges = true
		}
	}

	// Handle lists
	listsHandler := utils.IndexedItemHandler[verityIpv6PrefixListListsModel, openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner]{
		CreateNew: func(planItem verityIpv6PrefixListListsModel) openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner {
			newItem := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &newItem.Enable, TFValue: planItem.Enable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PermitDeny", APIField: &newItem.PermitDeny, TFValue: planItem.PermitDeny},
				{FieldName: "Ipv6Prefix", APIField: &newItem.Ipv6Prefix, TFValue: planItem.Ipv6Prefix},
			})

			// Handle nullable int64 fields
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "GreaterThanEqualValue", APIField: &newItem.GreaterThanEqualValue, TFValue: planItem.GreaterThanEqualValue},
				{FieldName: "LessThanEqualValue", APIField: &newItem.LessThanEqualValue, TFValue: planItem.LessThanEqualValue},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &newItem.Index, TFValue: planItem.Index},
			})

			return newItem
		},
		UpdateExisting: func(planItem verityIpv6PrefixListListsModel, stateItem verityIpv6PrefixListListsModel) (openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner, bool) {
			updateItem := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{}
			fieldChanged := false

			// Handle boolean field changes
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateItem.Enable = v }, &fieldChanged)

			// Handle string field changes
			utils.CompareAndSetStringField(planItem.PermitDeny, stateItem.PermitDeny, func(v *string) { updateItem.PermitDeny = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Ipv6Prefix, stateItem.Ipv6Prefix, func(v *string) { updateItem.Ipv6Prefix = v }, &fieldChanged)

			// Handle nullable int64 field changes
			utils.CompareAndSetNullableInt64Field(planItem.GreaterThanEqualValue, stateItem.GreaterThanEqualValue, func(v *openapi.NullableInt32) { updateItem.GreaterThanEqualValue = *v }, &fieldChanged)
			utils.CompareAndSetNullableInt64Field(planItem.LessThanEqualValue, stateItem.LessThanEqualValue, func(v *openapi.NullableInt32) { updateItem.LessThanEqualValue = *v }, &fieldChanged)

			// Handle index field change
			utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateItem.Index = v }, &fieldChanged)

			return updateItem, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner {
			return openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(index)),
			}
		},
	}

	changedLists, listsChanged := utils.ProcessIndexedArrayUpdates(plan.Lists, state.Lists, listsHandler)
	if listsChanged {
		ipv6PrefixListProps.Lists = changedLists
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "ipv6_prefix_list", name, ipv6PrefixListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv6 Prefix List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_prefix_lists")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityIpv6PrefixListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityIpv6PrefixListResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "ipv6_prefix_list", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv6 Prefix List %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_prefix_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv6PrefixListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
