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
	_ resource.Resource                = &verityExtendedCommunityListResource{}
	_ resource.ResourceWithConfigure   = &verityExtendedCommunityListResource{}
	_ resource.ResourceWithImportState = &verityExtendedCommunityListResource{}
)

func NewVerityExtendedCommunityListResource() resource.Resource {
	return &verityExtendedCommunityListResource{}
}

type verityExtendedCommunityListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityExtendedCommunityListResourceModel struct {
	Name             types.String                                       `tfsdk:"name"`
	Enable           types.Bool                                         `tfsdk:"enable"`
	PermitDeny       types.String                                       `tfsdk:"permit_deny"`
	AnyAll           types.String                                       `tfsdk:"any_all"`
	StandardExpanded types.String                                       `tfsdk:"standard_expanded"`
	Lists            []verityExtendedCommunityListListsModel            `tfsdk:"lists"`
	ObjectProperties []verityExtendedCommunityListObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityExtendedCommunityListListsModel struct {
	Enable                        types.Bool   `tfsdk:"enable"`
	Mode                          types.String `tfsdk:"mode"`
	RouteTargetExpandedExpression types.String `tfsdk:"route_target_expanded_expression"`
	Index                         types.Int64  `tfsdk:"index"`
}

type verityExtendedCommunityListObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityExtendedCommunityListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extended_community_list"
}

func (r *verityExtendedCommunityListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityExtendedCommunityListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Extended Community List",
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
			"permit_deny": schema.StringAttribute{
				Description: "Action upon match of Community Strings.",
				Optional:    true,
			},
			"any_all": schema.StringAttribute{
				Description: "BGP does not advertise any or all routes that do not match the Community String",
				Optional:    true,
			},
			"standard_expanded": schema.StringAttribute{
				Description: "Used Community String or Expanded Expression",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"lists": schema.ListNestedBlock{
				Description: "List of Extended Community List entries",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable of this Extended Community List",
							Optional:    true,
						},
						"mode": schema.StringAttribute{
							Description: "Mode",
							Optional:    true,
						},
						"route_target_expanded_expression": schema.StringAttribute{
							Description: "Match against a BGP extended community of type Route Target",
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
				Description: "Object properties for the Extended Community List",
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

func (r *verityExtendedCommunityListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityExtendedCommunityListResourceModel
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
	extCommunityListReq := &openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		extCommunityListReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.PermitDeny.IsNull() {
		extCommunityListReq.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
	}
	if !plan.AnyAll.IsNull() {
		extCommunityListReq.AnyAll = openapi.PtrString(plan.AnyAll.ValueString())
	}
	if !plan.StandardExpanded.IsNull() {
		extCommunityListReq.StandardExpanded = openapi.PtrString(plan.StandardExpanded.ValueString())
	}

	if len(plan.Lists) > 0 {
		lists := make([]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner, len(plan.Lists))
		for i, listItem := range plan.Lists {
			apiListItem := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{}

			if !listItem.Enable.IsNull() {
				apiListItem.Enable = openapi.PtrBool(listItem.Enable.ValueBool())
			}
			if !listItem.Mode.IsNull() {
				apiListItem.Mode = openapi.PtrString(listItem.Mode.ValueString())
			}
			if !listItem.RouteTargetExpandedExpression.IsNull() {
				apiListItem.RouteTargetExpandedExpression = openapi.PtrString(listItem.RouteTargetExpandedExpression.ValueString())
			}
			if !listItem.Index.IsNull() {
				apiListItem.Index = openapi.PtrInt32(int32(listItem.Index.ValueInt64()))
			}

			lists[i] = apiListItem
		}
		extCommunityListReq.Lists = lists
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objectProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			objectProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objectProps.Notes = nil
		}
		extCommunityListReq.ObjectProperties = &objectProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "extended_community_list", name, *extCommunityListReq)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for extended community list creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Extended Community List %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Extended Community List %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "extended_community_lists")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityExtendedCommunityListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityExtendedCommunityListResourceModel
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

	extCommListName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("extended_community_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Extended Community List %s verification – trusting recent successful API operation", extCommListName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Extended Community Lists for verification of %s", extCommListName))

	type ExtendedCommunityListResponse struct {
		ExtendedCommunityList map[string]interface{} `json:"extended_community_list"`
	}

	var result ExtendedCommunityListResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		extCommListData, fetchErr := getCachedResponse(ctx, r.provCtx, "extended_community_lists", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Extended Community Lists")
			respAPI, err := r.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Extended Community Lists: %v", err)
			}
			defer respAPI.Body.Close()

			var res ExtendedCommunityListResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Extended Community Lists response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Extended Community Lists", len(res.ExtendedCommunityList)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Extended Community Lists on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = extCommListData.(ExtendedCommunityListResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Extended Community List %s", extCommListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Extended Community List with ID: %s", extCommListName))
	var extCommListData map[string]interface{}
	exists := false

	if data, ok := result.ExtendedCommunityList[extCommListName].(map[string]interface{}); ok {
		extCommListData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Extended Community List directly by ID: %s", extCommListName))
	} else {
		for apiName, ecl := range result.ExtendedCommunityList {
			extCommList, ok := ecl.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := extCommList["name"].(string); ok && name == extCommListName {
				extCommListData = extCommList
				extCommListName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Extended Community List with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Extended Community List with ID '%s' not found in API response", extCommListName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", extCommListData["name"]))

	if enable, ok := extCommListData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if permitDeny, ok := extCommListData["permit_deny"].(string); ok {
		state.PermitDeny = types.StringValue(permitDeny)
	} else {
		state.PermitDeny = types.StringNull()
	}

	if anyAll, ok := extCommListData["any_all"].(string); ok {
		state.AnyAll = types.StringValue(anyAll)
	} else {
		state.AnyAll = types.StringNull()
	}

	if standardExpanded, ok := extCommListData["standard_expanded"].(string); ok {
		state.StandardExpanded = types.StringValue(standardExpanded)
	} else {
		state.StandardExpanded = types.StringNull()
	}

	if lists, ok := extCommListData["lists"].([]interface{}); ok && len(lists) > 0 {
		var listItems []verityExtendedCommunityListListsModel
		for _, l := range lists {
			listItem, ok := l.(map[string]interface{})
			if !ok {
				continue
			}
			listModel := verityExtendedCommunityListListsModel{}
			if enable, ok := listItem["enable"].(bool); ok {
				listModel.Enable = types.BoolValue(enable)
			} else {
				listModel.Enable = types.BoolNull()
			}
			if mode, ok := listItem["mode"].(string); ok {
				listModel.Mode = types.StringValue(mode)
			} else {
				listModel.Mode = types.StringNull()
			}
			if routeTarget, ok := listItem["route_target_expanded_expression"].(string); ok {
				listModel.RouteTargetExpandedExpression = types.StringValue(routeTarget)
			} else {
				listModel.RouteTargetExpandedExpression = types.StringNull()
			}
			if index, ok := listItem["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					listModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					listModel.Index = types.Int64Value(int64(intVal))
				} else {
					listModel.Index = types.Int64Null()
				}
			} else {
				listModel.Index = types.Int64Null()
			}
			listItems = append(listItems, listModel)
		}
		state.Lists = listItems
	} else {
		state.Lists = nil
	}

	// Only set object_properties if it exists in the API response
	if objectProps, ok := extCommListData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objectProps["notes"].(string); ok {
			state.ObjectProperties = []verityExtendedCommunityListObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityExtendedCommunityListObjectPropertiesModel{
				{Notes: types.StringNull()},
			}
		}
	} else {
		state.ObjectProperties = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityExtendedCommunityListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityExtendedCommunityListResourceModel

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
	extCommListProps := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue{}
	hasChanges := false

	if !plan.Enable.Equal(state.Enable) {
		extCommListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.PermitDeny.Equal(state.PermitDeny) {
		extCommListProps.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
		hasChanges = true
	}

	if !plan.AnyAll.Equal(state.AnyAll) {
		extCommListProps.AnyAll = openapi.PtrString(plan.AnyAll.ValueString())
		hasChanges = true
	}

	if !plan.StandardExpanded.Equal(state.StandardExpanded) {
		extCommListProps.StandardExpanded = openapi.PtrString(plan.StandardExpanded.ValueString())
		hasChanges = true
	}

	oldListsByIndex := make(map[int64]verityExtendedCommunityListListsModel)
	for _, item := range state.Lists {
		if !item.Index.IsNull() {
			idx := item.Index.ValueInt64()
			oldListsByIndex[idx] = item
		}
	}

	var changedLists []openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner
	listsChanged := false

	for _, planItem := range plan.Lists {
		if planItem.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planItem.Index.ValueInt64()
		stateItem, exists := oldListsByIndex[idx]

		if !exists {
			// CREATE: new item, include all fields
			newItem := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planItem.Enable.IsNull() {
				newItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			}

			if !planItem.Mode.IsNull() {
				newItem.Mode = openapi.PtrString(planItem.Mode.ValueString())
			}

			if !planItem.RouteTargetExpandedExpression.IsNull() {
				newItem.RouteTargetExpandedExpression = openapi.PtrString(planItem.RouteTargetExpandedExpression.ValueString())
			}

			changedLists = append(changedLists, newItem)
			listsChanged = true
			continue
		}

		// UPDATE: existing item, check which fields changed
		updateItem := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planItem.Enable.Equal(stateItem.Enable) {
			updateItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			fieldChanged = true
		}

		if !planItem.Mode.Equal(stateItem.Mode) {
			updateItem.Mode = openapi.PtrString(planItem.Mode.ValueString())
			fieldChanged = true
		}

		if !planItem.RouteTargetExpandedExpression.Equal(stateItem.RouteTargetExpandedExpression) {
			updateItem.RouteTargetExpandedExpression = openapi.PtrString(planItem.RouteTargetExpandedExpression.ValueString())
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
			deletedItem := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedLists = append(changedLists, deletedItem)
			listsChanged = true
		}
	}

	if listsChanged && len(changedLists) > 0 {
		extCommListProps.Lists = changedLists
		hasChanges = true
	}

	// Handle ObjectProperties changes
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				objProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				objProps.Notes = nil
			}
			extCommListProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "extended_community_list", name, extCommListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Extended Community List update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Extended Community List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Extended Community List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "extended_community_lists")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityExtendedCommunityListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityExtendedCommunityListResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "extended_community_list", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Extended Community List delete operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Extended Community List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Extended Community List %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "extended_community_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityExtendedCommunityListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
