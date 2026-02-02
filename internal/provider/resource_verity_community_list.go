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
	_ resource.Resource                = &verityCommunityListResource{}
	_ resource.ResourceWithConfigure   = &verityCommunityListResource{}
	_ resource.ResourceWithImportState = &verityCommunityListResource{}
	_ resource.ResourceWithModifyPlan  = &verityCommunityListResource{}
)

const communityListResourceType = "communitylists"

func NewVerityCommunityListResource() resource.Resource {
	return &verityCommunityListResource{}
}

type verityCommunityListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
				Computed:    true,
			},
			"permit_deny": schema.StringAttribute{
				Description: "Action upon match of Community Strings.",
				Optional:    true,
				Computed:    true,
			},
			"any_all": schema.StringAttribute{
				Description: "BGP does not advertise any or all routes that do not match the Community String",
				Optional:    true,
				Computed:    true,
			},
			"standard_expanded": schema.StringAttribute{
				Description: "Used Community String or Expanded Expression",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
						},
						"mode": schema.StringAttribute{
							Description: "Mode",
							Optional:    true,
							Computed:    true,
						},
						"community_string_expanded_expression": schema.StringAttribute{
							Description: "Community String in standard mode and Expanded Expression in Expanded mode",
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
				Description: "Object properties for the Community List",
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
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		communityListProps.ObjectProperties = &objProps
	}

	// Handle lists
	if len(plan.Lists) > 0 {
		lists := make([]openapi.CommunitylistsPutRequestCommunityListValueListsInner, len(plan.Lists))
		for i, item := range plan.Lists {
			listEntry := openapi.CommunitylistsPutRequestCommunityListValueListsInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &listEntry.Enable, TFValue: item.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Mode", APIField: &listEntry.Mode, TFValue: item.Mode},
				{FieldName: "CommunityStringExpandedExpression", APIField: &listEntry.CommunityStringExpandedExpression, TFValue: item.CommunityStringExpandedExpression},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &listEntry.Index, TFValue: item.Index},
			})
			lists[i] = listEntry
		}
		communityListProps.Lists = lists
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "community_list", name, *communityListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Community List %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "community_lists")

	var minState verityCommunityListResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if communityListData, exists := bulkMgr.GetResourceResponse("community_list", name); exists {
			state := populateCommunityListState(ctx, minState, communityListData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if communityListData, exists := r.bulkOpsMgr.GetResourceResponse("community_list", communityListName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached community list data for %s from recent operation", communityListName))
			state = populateCommunityListState(ctx, state, communityListData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateCommunityListState(ctx, state, communityListMap, r.provCtx.mode)
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
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objProps.Notes},
		}, &objPropsChanged)

		if objPropsChanged {
			communityListProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle lists
	listsHandler := utils.IndexedItemHandler[verityCommunityListListsModel, openapi.CommunitylistsPutRequestCommunityListValueListsInner]{
		CreateNew: func(planItem verityCommunityListListsModel) openapi.CommunitylistsPutRequestCommunityListValueListsInner {
			item := openapi.CommunitylistsPutRequestCommunityListValueListsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &item.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Mode", APIField: &item.Mode, TFValue: planItem.Mode},
				{FieldName: "CommunityStringExpandedExpression", APIField: &item.CommunityStringExpandedExpression, TFValue: planItem.CommunityStringExpandedExpression},
			})

			return item
		},
		UpdateExisting: func(planItem verityCommunityListListsModel, stateItem verityCommunityListListsModel) (openapi.CommunitylistsPutRequestCommunityListValueListsInner, bool) {
			item := openapi.CommunitylistsPutRequestCommunityListValueListsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { item.Enable = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.Mode, stateItem.Mode, func(v *string) { item.Mode = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.CommunityStringExpandedExpression, stateItem.CommunityStringExpandedExpression, func(v *string) { item.CommunityStringExpandedExpression = v }, &fieldChanged)

			return item, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.CommunitylistsPutRequestCommunityListValueListsInner {
			item := openapi.CommunitylistsPutRequestCommunityListValueListsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: types.Int64Value(index)},
			})
			return item
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "community_list", name, communityListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Community List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "community_lists")

	var minState verityCommunityListResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if communityListData, exists := bulkMgr.GetResourceResponse("community_list", name); exists {
			newState := populateCommunityListState(ctx, minState, communityListData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "community_list", name, nil, &resp.Diagnostics)
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

func populateCommunityListState(ctx context.Context, state verityCommunityListResourceModel, data map[string]interface{}, mode string) verityCommunityListResourceModel {
	const resourceType = communityListResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// String fields
	state.PermitDeny = utils.MapStringWithMode(data, "permit_deny", resourceType, mode)
	state.AnyAll = utils.MapStringWithMode(data, "any_all", resourceType, mode)
	state.StandardExpanded = utils.MapStringWithMode(data, "standard_expanded", resourceType, mode)

	// Handle lists array
	if utils.FieldAppliesToMode(resourceType, "lists", mode) {
		if listsData, ok := data["lists"].([]interface{}); ok && len(listsData) > 0 {
			var lists []verityCommunityListListsModel
			for _, l := range listsData {
				listItem, ok := l.(map[string]interface{})
				if !ok {
					continue
				}
				listModel := verityCommunityListListsModel{
					Enable:                            utils.MapBoolWithModeNested(listItem, "enable", resourceType, "lists.enable", mode),
					Mode:                              utils.MapStringWithModeNested(listItem, "mode", resourceType, "lists.mode", mode),
					CommunityStringExpandedExpression: utils.MapStringWithModeNested(listItem, "community_string_expanded_expression", resourceType, "lists.community_string_expanded_expression", mode),
					Index:                             utils.MapInt64WithModeNested(listItem, "index", resourceType, "lists.index", mode),
				}
				lists = append(lists, listModel)
			}
			state.Lists = lists
		} else {
			state.Lists = nil
		}
	} else {
		state.Lists = nil
	}

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityCommunityListObjectPropertiesModel{
				Notes: utils.MapStringWithModeNested(objProps, "notes", resourceType, "object_properties.notes", mode),
			}
			state.ObjectProperties = []verityCommunityListObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityCommunityListResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityCommunityListResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = communityListResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"permit_deny", "any_all", "standard_expanded",
	)

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "lists",
		ItemCount:    len(plan.Lists),
		StringFields: []string{"mode", "community_string_expanded_expression"},
		BoolFields:   []string{"enable"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"notes"},
	})
}
