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
	_ resource.Resource                = &verityCommunityListResource{}
	_ resource.ResourceWithConfigure   = &verityCommunityListResource{}
	_ resource.ResourceWithImportState = &verityCommunityListResource{}
)

func NewVerityCommunityListResource() resource.Resource {
	return &verityCommunityListResource{}
}

type verityCommunityListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityCommunityListResourceModel struct {
	Name             types.String                               `tfsdk:"name"`
	Enable           types.Bool                                 `tfsdk:"enable"`
	PermitDeny       types.String                               `tfsdk:"permit_deny"`
	AnyAll           types.String                               `tfsdk:"any_all"`
	StandardExpanded types.String                               `tfsdk:"standard_expanded"`
	Lists            []verityCommunityListListsModel            `tfsdk:"lists"`
	ObjectProperties []verityCommunityListObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityCommunityListListsModel struct {
	Enable                            types.Bool   `tfsdk:"enable"`
	Mode                              types.String `tfsdk:"mode"`
	CommunityStringExpandedExpression types.String `tfsdk:"community_string_expanded_expression"`
	Index                             types.Int64  `tfsdk:"index"`
}

type verityCommunityListObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityCommunityListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_community_list"
}

func (r *verityCommunityListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityCommunityListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Community List",
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
				Description: "List of Community List entries",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable of this Community List",
							Optional:    true,
						},
						"mode": schema.StringAttribute{
							Description: "Mode",
							Optional:    true,
						},
						"community_string_expanded_expression": schema.StringAttribute{
							Description: "Community String in standard mode and Expanded Expression in Expanded mode",
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
				Description: "Object properties for the Community List",
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

func (r *verityCommunityListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityCommunityListResourceModel
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
	communityListProps := &openapi.CommunitylistsPutRequestCommunityListValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		communityListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.PermitDeny.IsNull() {
		communityListProps.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
	}

	if !plan.AnyAll.IsNull() {
		communityListProps.AnyAll = openapi.PtrString(plan.AnyAll.ValueString())
	}

	if !plan.StandardExpanded.IsNull() {
		communityListProps.StandardExpanded = openapi.PtrString(plan.StandardExpanded.ValueString())
	}

	if len(plan.Lists) > 0 {
		lists := make([]openapi.CommunitylistsPutRequestCommunityListValueListsInner, len(plan.Lists))
		for i, listItem := range plan.Lists {
			listEntry := openapi.CommunitylistsPutRequestCommunityListValueListsInner{}

			if !listItem.Index.IsNull() {
				listEntry.Index = openapi.PtrInt32(int32(listItem.Index.ValueInt64()))
			}
			if !listItem.Enable.IsNull() {
				listEntry.Enable = openapi.PtrBool(listItem.Enable.ValueBool())
			}
			if !listItem.Mode.IsNull() {
				listEntry.Mode = openapi.PtrString(listItem.Mode.ValueString())
			}
			if !listItem.CommunityStringExpandedExpression.IsNull() {
				listEntry.CommunityStringExpandedExpression = openapi.PtrString(listItem.CommunityStringExpandedExpression.ValueString())
			}
			lists[i] = listEntry
		}
		communityListProps.Lists = lists
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			objProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objProps.Notes = nil
		}
		communityListProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "community_list", name, *communityListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Community List creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Community List %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Community List %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "community_lists")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityCommunityListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityCommunityListResourceModel
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

	communityListName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("community_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Community List %s verification – trusting recent successful API operation", communityListName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Community Lists for verification of %s", communityListName))

	type CommunityListsResponse struct {
		CommunityList map[string]interface{} `json:"community_list"`
	}

	var result CommunityListsResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		communityListsData, fetchErr := getCachedResponse(ctx, r.provCtx, "communitylists", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Community Lists")
			respAPI, err := r.client.CommunityListsAPI.CommunitylistsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Community Lists: %v", err)
			}
			defer respAPI.Body.Close()

			var res CommunityListsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Community Lists response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Community Lists", len(res.CommunityList)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Community Lists on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = communityListsData.(CommunityListsResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Community List %s", communityListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Community List with ID: %s", communityListName))
	var communityListData map[string]interface{}
	exists := false

	if data, ok := result.CommunityList[communityListName].(map[string]interface{}); ok {
		communityListData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Community List directly by ID: %s", communityListName))
	} else {
		for apiName, cl := range result.CommunityList {
			communityList, ok := cl.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := communityList["name"].(string); ok && name == communityListName {
				communityListData = communityList
				communityListName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Community List with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Community List with ID '%s' not found in API response", communityListName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", communityListData["name"]))

	if enable, ok := communityListData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if permitDeny, ok := communityListData["permit_deny"].(string); ok && permitDeny != "" {
		state.PermitDeny = types.StringValue(permitDeny)
	} else {
		state.PermitDeny = types.StringNull()
	}

	if anyAll, ok := communityListData["any_all"].(string); ok && anyAll != "" {
		state.AnyAll = types.StringValue(anyAll)
	} else {
		state.AnyAll = types.StringNull()
	}

	if standardExpanded, ok := communityListData["standard_expanded"].(string); ok && standardExpanded != "" {
		state.StandardExpanded = types.StringValue(standardExpanded)
	} else {
		state.StandardExpanded = types.StringNull()
	}

	if listsData, ok := communityListData["lists"].([]interface{}); ok {
		var lists []verityCommunityListListsModel
		for _, listItemData := range listsData {
			if listItem, ok := listItemData.(map[string]interface{}); ok {
				var list verityCommunityListListsModel

				if enable, ok := listItem["enable"].(bool); ok {
					list.Enable = types.BoolValue(enable)
				} else {
					list.Enable = types.BoolNull()
				}

				if mode, ok := listItem["mode"].(string); ok {
					list.Mode = types.StringValue(mode)
				} else {
					list.Mode = types.StringNull()
				}

				if communityStringExpandedExpression, ok := listItem["community_string_expanded_expression"].(string); ok {
					list.CommunityStringExpandedExpression = types.StringValue(communityStringExpandedExpression)
				} else {
					list.CommunityStringExpandedExpression = types.StringNull()
				}

				if index, ok := listItem["index"].(float64); ok {
					list.Index = types.Int64Value(int64(index))
				} else if index, ok := listItem["index"].(int); ok {
					list.Index = types.Int64Value(int64(index))
				} else if index, ok := listItem["index"].(int64); ok {
					list.Index = types.Int64Value(index)
				} else {
					list.Index = types.Int64Null()
				}

				lists = append(lists, list)
			}
		}
		state.Lists = lists
	} else {
		state.Lists = nil
	}

	if objProps, ok := communityListData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objProps["notes"].(string); ok {
			state.ObjectProperties = []verityCommunityListObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityCommunityListObjectPropertiesModel{
				{Notes: types.StringValue("")},
			}
		}
	} else {
		state.ObjectProperties = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityCommunityListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityCommunityListResourceModel

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
	communityListProps := openapi.CommunitylistsPutRequestCommunityListValue{}
	hasChanges := false

	if !plan.Enable.Equal(state.Enable) {
		communityListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.PermitDeny.Equal(state.PermitDeny) {
		communityListProps.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
		hasChanges = true
	}

	if !plan.AnyAll.Equal(state.AnyAll) {
		communityListProps.AnyAll = openapi.PtrString(plan.AnyAll.ValueString())
		hasChanges = true
	}

	if !plan.StandardExpanded.Equal(state.StandardExpanded) {
		communityListProps.StandardExpanded = openapi.PtrString(plan.StandardExpanded.ValueString())
		hasChanges = true
	}

	stateListsByIndex := make(map[int64]verityCommunityListListsModel)
	for _, item := range state.Lists {
		if !item.Index.IsNull() {
			idx := item.Index.ValueInt64()
			stateListsByIndex[idx] = item
		}
	}

	var changedLists []openapi.CommunitylistsPutRequestCommunityListValueListsInner
	listsChanged := false

	for _, planItem := range plan.Lists {
		if planItem.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planItem.Index.ValueInt64()
		stateItem, exists := stateListsByIndex[idx]

		if !exists {
			// CREATE: new item, include all fields
			newItem := openapi.CommunitylistsPutRequestCommunityListValueListsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planItem.Enable.IsNull() {
				newItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			}

			if !planItem.Mode.IsNull() {
				newItem.Mode = openapi.PtrString(planItem.Mode.ValueString())
			}

			if !planItem.CommunityStringExpandedExpression.IsNull() {
				newItem.CommunityStringExpandedExpression = openapi.PtrString(planItem.CommunityStringExpandedExpression.ValueString())
			}

			changedLists = append(changedLists, newItem)
			listsChanged = true
			continue
		}

		// UPDATE: existing item, check which fields changed
		updateItem := openapi.CommunitylistsPutRequestCommunityListValueListsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planItem.Enable.Equal(stateItem.Enable) {
			if !planItem.Enable.IsNull() {
				updateItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			}
			fieldChanged = true
		}

		if !planItem.Mode.Equal(stateItem.Mode) {
			if !planItem.Mode.IsNull() {
				updateItem.Mode = openapi.PtrString(planItem.Mode.ValueString())
			}
			fieldChanged = true
		}

		if !planItem.CommunityStringExpandedExpression.Equal(stateItem.CommunityStringExpandedExpression) {
			if !planItem.CommunityStringExpandedExpression.IsNull() {
				updateItem.CommunityStringExpandedExpression = openapi.PtrString(planItem.CommunityStringExpandedExpression.ValueString())
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedLists = append(changedLists, updateItem)
			listsChanged = true
		}
	}

	// DELETE: Check for removed items
	for stateIdx := range stateListsByIndex {
		found := false
		for _, planItem := range plan.Lists {
			if !planItem.Index.IsNull() && planItem.Index.ValueInt64() == stateIdx {
				found = true
				break
			}
		}

		if !found {
			// item removed - include only the index for deletion
			deletedItem := openapi.CommunitylistsPutRequestCommunityListValueListsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedLists = append(changedLists, deletedItem)
			listsChanged = true
		}
	}

	if listsChanged && len(changedLists) > 0 {
		communityListProps.Lists = changedLists
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			op := plan.ObjectProperties[0]
			objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !op.Notes.IsNull() {
				objProps.Notes = openapi.PtrString(op.Notes.ValueString())
			} else {
				objProps.Notes = nil
			}
			communityListProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "community_list", name, communityListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Community List update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Community List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Community List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "community_lists")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityCommunityListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityCommunityListResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "community_list", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Community List deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Community List %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Community List %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "community_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityCommunityListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
