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

func (m verityCommunityListListsModel) GetIndex() types.Int64 {
	return m.Index
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

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "PermitDeny", APIField: &communityListProps.PermitDeny, TFValue: plan.PermitDeny},
		{FieldName: "AnyAll", APIField: &communityListProps.AnyAll, TFValue: plan.AnyAll},
		{FieldName: "StandardExpanded", APIField: &communityListProps.StandardExpanded, TFValue: plan.StandardExpanded},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &communityListProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
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

	// Handle lists
	if len(plan.Lists) > 0 {
		lists := make([]openapi.CommunitylistsPutRequestCommunityListValueListsInner, len(plan.Lists))
		for i, listItem := range plan.Lists {
			listEntry := openapi.CommunitylistsPutRequestCommunityListValueListsInner{}
			if !listItem.Enable.IsNull() {
				listEntry.Enable = openapi.PtrBool(listItem.Enable.ValueBool())
			}
			if !listItem.Mode.IsNull() {
				listEntry.Mode = openapi.PtrString(listItem.Mode.ValueString())
			}
			if !listItem.CommunityStringExpandedExpression.IsNull() {
				listEntry.CommunityStringExpandedExpression = openapi.PtrString(listItem.CommunityStringExpandedExpression.ValueString())
			}
			if !listItem.Index.IsNull() {
				listEntry.Index = openapi.PtrInt32(int32(listItem.Index.ValueInt64()))
			}
			lists[i] = listEntry
		}
		communityListProps.Lists = lists
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "community_list", name, *communityListProps, &resp.Diagnostics)
	if !success {
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

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "community_lists", communityListName,
		func() (CommunityListsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Community Lists")
			respAPI, err := r.client.CommunityListsAPI.CommunitylistsGet(ctx).Execute()
			if err != nil {
				return CommunityListsResponse{}, fmt.Errorf("error reading Community Lists: %v", err)
			}
			defer respAPI.Body.Close()

			var res CommunityListsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return CommunityListsResponse{}, fmt.Errorf("failed to decode Community Lists response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Community Lists", len(res.CommunityList)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Community List %s", communityListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Community List with name: %s", communityListName))

	communityListData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.CommunityList,
		communityListName,
		func(data interface{}) (string, bool) {
			if communityList, ok := data.(map[string]interface{}); ok {
				if name, ok := communityList["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Community List with name '%s' not found in API response", communityListName))
		resp.State.RemoveResource(ctx)
		return
	}

	communityListMap, ok := communityListData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Community List Data",
			fmt.Sprintf("Community List data is not in expected format for %s", communityListName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Community List '%s' under API key '%s'", communityListName, actualAPIName))

	state.Name = utils.MapStringFromAPI(communityListMap["name"])

	// Handle object properties
	if objProps, ok := communityListMap["object_properties"].(map[string]interface{}); ok {
		notes := utils.MapStringFromAPI(objProps["notes"])
		if notes.IsNull() {
			notes = types.StringValue("")
		}
		state.ObjectProperties = []verityCommunityListObjectPropertiesModel{
			{Notes: notes},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"permit_deny":       &state.PermitDeny,
		"any_all":           &state.AnyAll,
		"standard_expanded": &state.StandardExpanded,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(communityListMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(communityListMap[apiKey])
	}

	// Handle lists
	if listsData, ok := communityListMap["lists"].([]interface{}); ok && len(listsData) > 0 {
		var lists []verityCommunityListListsModel

		for _, l := range listsData {
			listItem, ok := l.(map[string]interface{})
			if !ok {
				continue
			}

			listModel := verityCommunityListListsModel{
				Enable:                            utils.MapBoolFromAPI(listItem["enable"]),
				Mode:                              utils.MapStringFromAPI(listItem["mode"]),
				CommunityStringExpandedExpression: utils.MapStringFromAPI(listItem["community_string_expanded_expression"]),
				Index:                             utils.MapInt64FromAPI(listItem["index"]),
			}

			lists = append(lists, listModel)
		}

		state.Lists = lists
	} else {
		state.Lists = nil
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

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { communityListProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PermitDeny, state.PermitDeny, func(v *string) { communityListProps.PermitDeny = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AnyAll, state.AnyAll, func(v *string) { communityListProps.AnyAll = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.StandardExpanded, state.StandardExpanded, func(v *string) { communityListProps.StandardExpanded = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { communityListProps.Enable = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				objProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				objProps.Notes = nil
			}
			communityListProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle lists
	listsHandler := utils.IndexedItemHandler[verityCommunityListListsModel, openapi.CommunitylistsPutRequestCommunityListValueListsInner]{
		CreateNew: func(planItem verityCommunityListListsModel) openapi.CommunitylistsPutRequestCommunityListValueListsInner {
			item := openapi.CommunitylistsPutRequestCommunityListValueListsInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			if !planItem.Enable.IsNull() {
				item.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			} else {
				item.Enable = openapi.PtrBool(false)
			}

			if !planItem.Mode.IsNull() && planItem.Mode.ValueString() != "" {
				item.Mode = openapi.PtrString(planItem.Mode.ValueString())
			} else {
				item.Mode = openapi.PtrString("")
			}

			if !planItem.CommunityStringExpandedExpression.IsNull() && planItem.CommunityStringExpandedExpression.ValueString() != "" {
				item.CommunityStringExpandedExpression = openapi.PtrString(planItem.CommunityStringExpandedExpression.ValueString())
			} else {
				item.CommunityStringExpandedExpression = openapi.PtrString("")
			}

			return item
		},
		UpdateExisting: func(planItem verityCommunityListListsModel, stateItem verityCommunityListListsModel) (openapi.CommunitylistsPutRequestCommunityListValueListsInner, bool) {
			item := openapi.CommunitylistsPutRequestCommunityListValueListsInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			fieldChanged := false

			if !planItem.Enable.Equal(stateItem.Enable) {
				if !planItem.Enable.IsNull() {
					item.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
				} else {
					item.Enable = openapi.PtrBool(false)
				}
				fieldChanged = true
			}

			if !planItem.Mode.Equal(stateItem.Mode) {
				if !planItem.Mode.IsNull() && planItem.Mode.ValueString() != "" {
					item.Mode = openapi.PtrString(planItem.Mode.ValueString())
				} else {
					item.Mode = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.CommunityStringExpandedExpression.Equal(stateItem.CommunityStringExpandedExpression) {
				if !planItem.CommunityStringExpandedExpression.IsNull() && planItem.CommunityStringExpandedExpression.ValueString() != "" {
					item.CommunityStringExpandedExpression = openapi.PtrString(planItem.CommunityStringExpandedExpression.ValueString())
				} else {
					item.CommunityStringExpandedExpression = openapi.PtrString("")
				}
				fieldChanged = true
			}

			return item, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.CommunitylistsPutRequestCommunityListValueListsInner {
			return openapi.CommunitylistsPutRequestCommunityListValueListsInner{
				Index: openapi.PtrInt32(int32(index)),
			}
		},
	}

	changedLists, listsChanged := utils.ProcessIndexedArrayUpdates(plan.Lists, state.Lists, listsHandler)
	if listsChanged {
		communityListProps.Lists = changedLists
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "community_list", name, communityListProps, &resp.Diagnostics)
	if !success {
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "community_list", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Community List %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "community_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityCommunityListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
