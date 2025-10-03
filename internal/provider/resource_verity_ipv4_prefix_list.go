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

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ipv4PrefixListProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
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

	// Handle lists
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "ipv4_prefix_list", name, *ipv4PrefixListProps, &resp.Diagnostics)
	if !success {
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

	state.Name = utils.MapStringFromAPI(ipv4PrefixListMap["name"])

	// Handle object properties
	if objectPropsData, ok := ipv4PrefixListMap["object_properties"].(map[string]interface{}); ok {
		notes := utils.MapStringFromAPI(objectPropsData["notes"])
		if notes.IsNull() {
			notes = types.StringValue("")
		}
		state.ObjectProperties = []verityIpv4PrefixListObjectPropertiesModel{
			{Notes: notes},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ipv4PrefixListMap[apiKey])
	}

	// Handle lists
	if listsData, ok := ipv4PrefixListMap["lists"].([]interface{}); ok && len(listsData) > 0 {
		var lists []verityIpv4PrefixListListsModel

		for _, item := range listsData {
			listItem, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			listModel := verityIpv4PrefixListListsModel{
				Enable:                utils.MapBoolFromAPI(listItem["enable"]),
				PermitDeny:            utils.MapStringFromAPI(listItem["permit_deny"]),
				Ipv4Prefix:            utils.MapStringFromAPI(listItem["ipv4_prefix"]),
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

	// Handle lists
	listsHandler := utils.IndexedItemHandler[verityIpv4PrefixListListsModel, openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner]{
		CreateNew: func(planItem verityIpv4PrefixListListsModel) openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner {
			newItem := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
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

			return newItem
		},
		UpdateExisting: func(planItem verityIpv4PrefixListListsModel, stateItem verityIpv4PrefixListListsModel) (openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner, bool) {
			updateItem := openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "ipv4_prefix_list", name, ipv4PrefixListProps, &resp.Diagnostics)
	if !success {
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "ipv4_prefix_list", name, nil, &resp.Diagnostics)
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
