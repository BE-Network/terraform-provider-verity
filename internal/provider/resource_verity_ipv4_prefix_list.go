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
	_ resource.Resource                = &verityIpv4PrefixListResource{}
	_ resource.ResourceWithConfigure   = &verityIpv4PrefixListResource{}
	_ resource.ResourceWithImportState = &verityIpv4PrefixListResource{}
)

func NewVerityIpv4PrefixListResource() resource.Resource {
	return &verityIpv4PrefixListResource{}
}

type verityIpv4PrefixListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
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
						},
						"permit_deny": schema.StringAttribute{
							Description: "Action upon match of Community Strings.",
							Optional:    true,
						},
						"ipv4_prefix": schema.StringAttribute{
							Description: "IPv4 address and subnet to match against",
							Optional:    true,
						},
						"greater_than_equal_value": schema.Int64Attribute{
							Description: "Match IP routes with a subnet mask greater than or equal to the value indicated (maximum: 32)",
							Optional:    true,
						},
						"less_than_equal_value": schema.Int64Attribute{
							Description: "Match IP routes with a subnet mask less than or equal to the value indicated (maximum: 32)",
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
				Description: "Object properties for the IPv4 Prefix List",
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

func (r *verityIpv4PrefixListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityIpv4PrefixListResourceModel
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
	ipv4PrefixListProps := &openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		ipv4PrefixListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if len(plan.Lists) > 0 {
		var lists []openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner
		for _, listItem := range plan.Lists {
			item := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{}

			if !listItem.Enable.IsNull() {
				item.Enable = openapi.PtrBool(listItem.Enable.ValueBool())
			}
			if !listItem.PermitDeny.IsNull() {
				item.PermitDeny = openapi.PtrString(listItem.PermitDeny.ValueString())
			}
			if !listItem.Ipv4Prefix.IsNull() {
				item.Ipv4Prefix = openapi.PtrString(listItem.Ipv4Prefix.ValueString())
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
		ipv4PrefixListProps.Lists = lists
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objectProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			objectProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objectProps.Notes = nil
		}
		ipv4PrefixListProps.ObjectProperties = &objectProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "ipv4_prefix_list", name, ipv4PrefixListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv4 Prefix List create operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create IPv4 Prefix List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv4 Prefix List %s create operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_prefix_lists")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ipv4_prefix_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv4 prefix list %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv4 prefix lists for verification of %s", name))

	type Ipv4PrefixListResponse struct {
		Ipv4PrefixList map[string]map[string]interface{} `json:"ipv4_prefix_list"`
	}

	var result Ipv4PrefixListResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		ipv4PrefixListsData, fetchErr := getCachedResponse(ctx, r.provCtx, "ipv4_prefix_lists", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv4 prefix lists")
			resp, err := r.client.IPv4PrefixListsAPI.Ipv4prefixlistsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading IPv4 prefix lists: %v", err)
			}
			defer resp.Body.Close()

			var result Ipv4PrefixListResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, fmt.Errorf("error decoding IPv4 prefix list response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv4 prefix lists", len(result.Ipv4PrefixList)))
			return result, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch IPv4 prefix lists on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = ipv4PrefixListsData.(Ipv4PrefixListResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv4 Prefix List %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv4 prefix list with ID: %s", name))
	var ipv4PrefixListData map[string]interface{}
	exists := false

	if data, ok := result.Ipv4PrefixList[name]; ok {
		ipv4PrefixListData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found IPv4 prefix list directly by ID: %s", name))
	} else {
		for apiName, p := range result.Ipv4PrefixList {
			if nameVal, ok := p["name"].(string); ok && nameVal == name {
				ipv4PrefixListData = p
				name = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found IPv4 prefix list with name '%s' under API key '%s'", nameVal, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv4 Prefix List with ID '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", ipv4PrefixListData["name"]))

	if enable, ok := ipv4PrefixListData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if listsData, ok := ipv4PrefixListData["lists"].([]interface{}); ok && len(listsData) > 0 {
		var lists []verityIpv4PrefixListListsModel

		for _, item := range listsData {
			listItem, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			listModel := verityIpv4PrefixListListsModel{}

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

			if ipv4Prefix, ok := listItem["ipv4_prefix"].(string); ok {
				listModel.Ipv4Prefix = types.StringValue(ipv4Prefix)
			} else {
				listModel.Ipv4Prefix = types.StringNull()
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

	if objectPropsData, ok := ipv4PrefixListData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objectPropsData["notes"].(string); ok {
			state.ObjectProperties = []verityIpv4PrefixListObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityIpv4PrefixListObjectPropertiesModel{
				{Notes: types.StringValue("")},
			}
		}
	} else {
		state.ObjectProperties = nil
	}

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

	if !plan.Enable.Equal(state.Enable) {
		ipv4PrefixListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	oldListsByIndex := make(map[int64]verityIpv4PrefixListListsModel)
	for _, item := range state.Lists {
		if !item.Index.IsNull() {
			idx := item.Index.ValueInt64()
			oldListsByIndex[idx] = item
		}
	}

	var changedLists []openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner
	listsChanged := false

	for _, planItem := range plan.Lists {
		if planItem.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planItem.Index.ValueInt64()
		stateItem, exists := oldListsByIndex[idx]

		if !exists {
			// CREATE: new item, include all fields
			newItem := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planItem.Enable.IsNull() {
				newItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			}

			if !planItem.PermitDeny.IsNull() {
				newItem.PermitDeny = openapi.PtrString(planItem.PermitDeny.ValueString())
			}

			if !planItem.Ipv4Prefix.IsNull() {
				newItem.Ipv4Prefix = openapi.PtrString(planItem.Ipv4Prefix.ValueString())
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
		updateItem := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{
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

		if !planItem.Ipv4Prefix.Equal(stateItem.Ipv4Prefix) {
			if !planItem.Ipv4Prefix.IsNull() {
				updateItem.Ipv4Prefix = openapi.PtrString(planItem.Ipv4Prefix.ValueString())
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
			deletedItem := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedLists = append(changedLists, deletedItem)
			listsChanged = true
		}
	}

	if listsChanged && len(changedLists) > 0 {
		ipv4PrefixListProps.Lists = changedLists
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
			ipv4PrefixListProps.ObjectProperties = &objectProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "ipv4_prefix_list", name, ipv4PrefixListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv4 Prefix List update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update IPv4 Prefix List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv4 Prefix List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_prefix_lists")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "ipv4_prefix_list", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv4 Prefix List delete operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete IPv4 Prefix List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv4 Prefix List %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_prefix_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv4PrefixListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
