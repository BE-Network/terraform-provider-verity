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

func (l verityExtendedCommunityListListsModel) GetIndex() types.Int64 {
	return l.Index
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
	extCommListProps := &openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "PermitDeny", APIField: &extCommListProps.PermitDeny, TFValue: plan.PermitDeny},
		{FieldName: "AnyAll", APIField: &extCommListProps.AnyAll, TFValue: plan.AnyAll},
		{FieldName: "StandardExpanded", APIField: &extCommListProps.StandardExpanded, TFValue: plan.StandardExpanded},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &extCommListProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		extCommListProps.ObjectProperties = &objProps
	}

	// Handle lists
	if len(plan.Lists) > 0 {
		lists := make([]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner, len(plan.Lists))
		for i, item := range plan.Lists {
			apiListItem := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &apiListItem.Enable, TFValue: item.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Mode", APIField: &apiListItem.Mode, TFValue: item.Mode},
				{FieldName: "RouteTargetExpandedExpression", APIField: &apiListItem.RouteTargetExpandedExpression, TFValue: item.RouteTargetExpandedExpression},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &apiListItem.Index, TFValue: item.Index},
			})

			lists[i] = apiListItem
		}
		extCommListProps.Lists = lists
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "extended_community_list", name, *extCommListProps, &resp.Diagnostics)
	if !success {
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

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "extended_community_lists", extCommListName,
		func() (ExtendedCommunityListResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Extended Community Lists")
			respAPI, err := r.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(ctx).Execute()
			if err != nil {
				return ExtendedCommunityListResponse{}, fmt.Errorf("error reading Extended Community Lists: %v", err)
			}
			defer respAPI.Body.Close()

			var res ExtendedCommunityListResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return ExtendedCommunityListResponse{}, fmt.Errorf("failed to decode Extended Community Lists response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Extended Community Lists", len(res.ExtendedCommunityList)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Extended Community List %s", extCommListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Extended Community List with name: %s", extCommListName))

	extCommListData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.ExtendedCommunityList,
		extCommListName,
		func(data interface{}) (string, bool) {
			if extCommList, ok := data.(map[string]interface{}); ok {
				if name, ok := extCommList["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Extended Community List with name '%s' not found in API response", extCommListName))
		resp.State.RemoveResource(ctx)
		return
	}

	extCommListMap, ok := extCommListData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Extended Community List Data",
			fmt.Sprintf("Extended Community List data is not in expected format for %s", extCommListName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Extended Community List '%s' under API key '%s'", extCommListName, actualAPIName))

	state.Name = utils.MapStringFromAPI(extCommListMap["name"])

	// Handle object properties
	if objectProps, ok := extCommListMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityExtendedCommunityListObjectPropertiesModel{
			{Notes: utils.MapStringFromAPI(objectProps["notes"])},
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
		*stateField = utils.MapStringFromAPI(extCommListMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(extCommListMap[apiKey])
	}

	// Handle lists
	if lists, ok := extCommListMap["lists"].([]interface{}); ok && len(lists) > 0 {
		var listItems []verityExtendedCommunityListListsModel
		for _, l := range lists {
			listItem, ok := l.(map[string]interface{})
			if !ok {
				continue
			}

			listModel := verityExtendedCommunityListListsModel{
				Enable:                        utils.MapBoolFromAPI(listItem["enable"]),
				Mode:                          utils.MapStringFromAPI(listItem["mode"]),
				RouteTargetExpandedExpression: utils.MapStringFromAPI(listItem["route_target_expanded_expression"]),
				Index:                         utils.MapInt64FromAPI(listItem["index"]),
			}

			listItems = append(listItems, listModel)
		}
		state.Lists = listItems
	} else {
		state.Lists = nil
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

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { extCommListProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PermitDeny, state.PermitDeny, func(v *string) { extCommListProps.PermitDeny = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AnyAll, state.AnyAll, func(v *string) { extCommListProps.AnyAll = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.StandardExpanded, state.StandardExpanded, func(v *string) { extCommListProps.StandardExpanded = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { extCommListProps.Enable = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objProps.Notes},
		}, &objPropsChanged)

		if objPropsChanged {
			extCommListProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle lists
	listsHandler := utils.IndexedItemHandler[verityExtendedCommunityListListsModel, openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner]{
		CreateNew: func(planItem verityExtendedCommunityListListsModel) openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner {
			item := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &item.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Mode", APIField: &item.Mode, TFValue: planItem.Mode},
				{FieldName: "RouteTargetExpandedExpression", APIField: &item.RouteTargetExpandedExpression, TFValue: planItem.RouteTargetExpandedExpression},
			})

			return item
		},
		UpdateExisting: func(planItem verityExtendedCommunityListListsModel, stateItem verityExtendedCommunityListListsModel) (openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner, bool) {
			item := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { item.Enable = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.Mode, stateItem.Mode, func(v *string) { item.Mode = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.RouteTargetExpandedExpression, stateItem.RouteTargetExpandedExpression, func(v *string) { item.RouteTargetExpandedExpression = v }, &fieldChanged)

			return item, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner {
			item := openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: types.Int64Value(index)},
			})
			return item
		},
	}

	changedLists, listsChanged := utils.ProcessIndexedArrayUpdates(plan.Lists, state.Lists, listsHandler)
	if listsChanged {
		extCommListProps.Lists = changedLists
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "extended_community_list", name, extCommListProps, &resp.Diagnostics)
	if !success {
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "extended_community_list", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Extended Community List %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "extended_community_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityExtendedCommunityListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
