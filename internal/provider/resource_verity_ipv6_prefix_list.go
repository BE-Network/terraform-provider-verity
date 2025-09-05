package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	bulkOpsMgr           *utils.BulkOperationManager
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

	if !plan.Enable.IsNull() {
		ipv6PrefixListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if len(plan.Lists) > 0 {
		var lists []openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner
		for _, listItem := range plan.Lists {
			item := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{}

			if !listItem.Enable.IsNull() {
				item.Enable = openapi.PtrBool(listItem.Enable.ValueBool())
			}
			if !listItem.PermitDeny.IsNull() {
				item.PermitDeny = openapi.PtrString(listItem.PermitDeny.ValueString())
			}
			if !listItem.Ipv6Prefix.IsNull() {
				item.Ipv6Prefix = openapi.PtrString(listItem.Ipv6Prefix.ValueString())
			}
			if !listItem.GreaterThanEqualValue.IsNull() {
				val := int32(listItem.GreaterThanEqualValue.ValueInt64())
				item.GreaterThanEqualValue = *openapi.NewNullableInt32(&val)
			} else {
				item.GreaterThanEqualValue = *openapi.NewNullableInt32(nil)
			}
			if !listItem.LessThanEqualValue.IsNull() {
				val := int32(listItem.LessThanEqualValue.ValueInt64())
				item.LessThanEqualValue = *openapi.NewNullableInt32(&val)
			} else {
				item.LessThanEqualValue = *openapi.NewNullableInt32(nil)
			}
			if !listItem.Index.IsNull() {
				item.Index = openapi.PtrInt32(int32(listItem.Index.ValueInt64()))
			}
			lists = append(lists, item)
		}
		ipv6PrefixListProps.Lists = lists
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objectProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			objectProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objectProps.Notes = nil
		}
		ipv6PrefixListProps.ObjectProperties = &objectProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "ipv6_prefix_list", name, ipv6PrefixListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv6 Prefix List create operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create IPv6 Prefix List %s", name))...,
		)
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

	var result Ipv6PrefixListResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		ipv6PrefixListsData, fetchErr := getCachedResponse(ctx, r.provCtx, "ipv6_prefix_lists", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv6 prefix lists")
			resp, err := r.client.IPv6PrefixListsAPI.Ipv6prefixlistsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading IPv6 prefix lists: %v", err)
			}
			defer resp.Body.Close()

			var result Ipv6PrefixListResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, fmt.Errorf("error decoding IPv6 prefix list response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv6 prefix lists", len(result.Ipv6PrefixList)))
			return result, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch IPv6 prefix lists on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = ipv6PrefixListsData.(Ipv6PrefixListResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv6 Prefix List %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv6 prefix list with ID: %s", name))
	var ipv6PrefixListData map[string]interface{}
	exists := false

	if data, ok := result.Ipv6PrefixList[name]; ok {
		ipv6PrefixListData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found IPv6 prefix list directly by ID: %s", name))
	} else {
		for apiName, p := range result.Ipv6PrefixList {
			if nameVal, ok := p["name"].(string); ok && nameVal == name {
				ipv6PrefixListData = p
				name = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found IPv6 prefix list with name '%s' under API key '%s'", nameVal, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv6 Prefix List with ID '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", ipv6PrefixListData["name"]))

	if enable, ok := ipv6PrefixListData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if listsData, ok := ipv6PrefixListData["lists"].([]interface{}); ok && len(listsData) > 0 {
		var lists []verityIpv6PrefixListListsModel

		for _, item := range listsData {
			listItem, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			listModel := verityIpv6PrefixListListsModel{}

			if enable, ok := listItem["enable"].(bool); ok {
				listModel.Enable = types.BoolValue(enable)
			} else {
				listModel.Enable = types.BoolNull()
			}

			if permitDeny, ok := listItem["permit_deny"].(string); ok {
				listModel.PermitDeny = types.StringValue(permitDeny)
			} else {
				listModel.PermitDeny = types.StringNull()
			}

			if ipv6Prefix, ok := listItem["ipv6_prefix"].(string); ok {
				listModel.Ipv6Prefix = types.StringValue(ipv6Prefix)
			} else {
				listModel.Ipv6Prefix = types.StringNull()
			}

			if greaterThanEqual, exists := listItem["greater_than_equal_value"]; exists && greaterThanEqual != nil {
				switch v := greaterThanEqual.(type) {
				case int:
					listModel.GreaterThanEqualValue = types.Int64Value(int64(v))
				case int32:
					listModel.GreaterThanEqualValue = types.Int64Value(int64(v))
				case int64:
					listModel.GreaterThanEqualValue = types.Int64Value(v)
				case float64:
					listModel.GreaterThanEqualValue = types.Int64Value(int64(v))
				case string:
					if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
						listModel.GreaterThanEqualValue = types.Int64Value(intVal)
					} else {
						listModel.GreaterThanEqualValue = types.Int64Null()
					}
				default:
					listModel.GreaterThanEqualValue = types.Int64Null()
				}
			} else {
				listModel.GreaterThanEqualValue = types.Int64Null()
			}

			if lessThanEqual, exists := listItem["less_than_equal_value"]; exists && lessThanEqual != nil {
				switch v := lessThanEqual.(type) {
				case int:
					listModel.LessThanEqualValue = types.Int64Value(int64(v))
				case int32:
					listModel.LessThanEqualValue = types.Int64Value(int64(v))
				case int64:
					listModel.LessThanEqualValue = types.Int64Value(v)
				case float64:
					listModel.LessThanEqualValue = types.Int64Value(int64(v))
				case string:
					if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
						listModel.LessThanEqualValue = types.Int64Value(intVal)
					} else {
						listModel.LessThanEqualValue = types.Int64Null()
					}
				default:
					listModel.LessThanEqualValue = types.Int64Null()
				}
			} else {
				listModel.LessThanEqualValue = types.Int64Null()
			}

			if index, exists := listItem["index"]; exists && index != nil {
				switch v := index.(type) {
				case int:
					listModel.Index = types.Int64Value(int64(v))
				case int32:
					listModel.Index = types.Int64Value(int64(v))
				case int64:
					listModel.Index = types.Int64Value(v)
				case float64:
					listModel.Index = types.Int64Value(int64(v))
				case string:
					if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
						listModel.Index = types.Int64Value(intVal)
					} else {
						listModel.Index = types.Int64Null()
					}
				default:
					listModel.Index = types.Int64Null()
				}
			} else {
				listModel.Index = types.Int64Null()
			}

			lists = append(lists, listModel)
		}
		state.Lists = lists
	} else {
		state.Lists = nil
	}

	if objectPropsData, ok := ipv6PrefixListData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objectPropsData["notes"].(string); ok {
			state.ObjectProperties = []verityIpv6PrefixListObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityIpv6PrefixListObjectPropertiesModel{
				{Notes: types.StringValue("")},
			}
		}
	} else {
		state.ObjectProperties = nil
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

	if !plan.Enable.Equal(state.Enable) {
		ipv6PrefixListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	oldListsByIndex := make(map[int64]verityIpv6PrefixListListsModel)
	for _, item := range state.Lists {
		if !item.Index.IsNull() {
			idx := item.Index.ValueInt64()
			oldListsByIndex[idx] = item
		}
	}

	var changedLists []openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner
	listsChanged := false

	for _, planItem := range plan.Lists {
		if planItem.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planItem.Index.ValueInt64()
		stateItem, exists := oldListsByIndex[idx]

		if !exists {
			// CREATE: new item, include all fields
			newItem := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planItem.Enable.IsNull() {
				newItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			}

			if !planItem.PermitDeny.IsNull() {
				newItem.PermitDeny = openapi.PtrString(planItem.PermitDeny.ValueString())
			}

			if !planItem.Ipv6Prefix.IsNull() {
				newItem.Ipv6Prefix = openapi.PtrString(planItem.Ipv6Prefix.ValueString())
			}

			if !planItem.GreaterThanEqualValue.IsNull() {
				val := int32(planItem.GreaterThanEqualValue.ValueInt64())
				newItem.GreaterThanEqualValue = *openapi.NewNullableInt32(&val)
			} else {
				newItem.GreaterThanEqualValue = *openapi.NewNullableInt32(nil)
			}

			if !planItem.LessThanEqualValue.IsNull() {
				val := int32(planItem.LessThanEqualValue.ValueInt64())
				newItem.LessThanEqualValue = *openapi.NewNullableInt32(&val)
			} else {
				newItem.LessThanEqualValue = *openapi.NewNullableInt32(nil)
			}

			changedLists = append(changedLists, newItem)
			listsChanged = true
			continue
		}

		// UPDATE: existing item, check which fields changed
		updateItem := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planItem.Enable.Equal(stateItem.Enable) {
			if !planItem.Enable.IsNull() {
				updateItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			}
			fieldChanged = true
		}

		if !planItem.PermitDeny.Equal(stateItem.PermitDeny) {
			if !planItem.PermitDeny.IsNull() {
				updateItem.PermitDeny = openapi.PtrString(planItem.PermitDeny.ValueString())
			}
			fieldChanged = true
		}

		if !planItem.Ipv6Prefix.Equal(stateItem.Ipv6Prefix) {
			if !planItem.Ipv6Prefix.IsNull() {
				updateItem.Ipv6Prefix = openapi.PtrString(planItem.Ipv6Prefix.ValueString())
			}
			fieldChanged = true
		}

		if !planItem.GreaterThanEqualValue.Equal(stateItem.GreaterThanEqualValue) {
			if !planItem.GreaterThanEqualValue.IsNull() {
				val := int32(planItem.GreaterThanEqualValue.ValueInt64())
				updateItem.GreaterThanEqualValue = *openapi.NewNullableInt32(&val)
			} else {
				updateItem.GreaterThanEqualValue = *openapi.NewNullableInt32(nil)
			}
			fieldChanged = true
		}

		if !planItem.LessThanEqualValue.Equal(stateItem.LessThanEqualValue) {
			if !planItem.LessThanEqualValue.IsNull() {
				val := int32(planItem.LessThanEqualValue.ValueInt64())
				updateItem.LessThanEqualValue = *openapi.NewNullableInt32(&val)
			} else {
				updateItem.LessThanEqualValue = *openapi.NewNullableInt32(nil)
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedLists = append(changedLists, updateItem)
			listsChanged = true
		}
	}

	// DELETE: Check for deleted items
	for stateIdx := range oldListsByIndex {
		found := false
		for _, planItem := range plan.Lists {
			if !planItem.Index.IsNull() && planItem.Index.ValueInt64() == stateIdx {
				found = true
				break
			}
		}

		if !found {
			// item removed - include only the index for deletion
			deletedItem := openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedLists = append(changedLists, deletedItem)
			listsChanged = true
		}
	}

	if listsChanged && len(changedLists) > 0 {
		ipv6PrefixListProps.Lists = changedLists
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objectProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				objectProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				objectProps.Notes = nil
			}
			ipv6PrefixListProps.ObjectProperties = &objectProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "ipv6_prefix_list", name, ipv6PrefixListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv6 Prefix List update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update IPv6 Prefix List %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "ipv6_prefix_list", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv6 Prefix List delete operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete IPv6 Prefix List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv6 Prefix List %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_prefix_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv6PrefixListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
