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
	_ resource.Resource                = &verityAsPathAccessListResource{}
	_ resource.ResourceWithConfigure   = &verityAsPathAccessListResource{}
	_ resource.ResourceWithImportState = &verityAsPathAccessListResource{}
)

func NewVerityAsPathAccessListResource() resource.Resource {
	return &verityAsPathAccessListResource{}
}

type verityAsPathAccessListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityAsPathAccessListResourceModel struct {
	Name             types.String                                  `tfsdk:"name"`
	Enable           types.Bool                                    `tfsdk:"enable"`
	PermitDeny       types.String                                  `tfsdk:"permit_deny"`
	Lists            []verityAsPathAccessListListsModel            `tfsdk:"lists"`
	ObjectProperties []verityAsPathAccessListObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityAsPathAccessListListsModel struct {
	Enable            types.Bool   `tfsdk:"enable"`
	RegularExpression types.String `tfsdk:"regular_expression"`
	Index             types.Int64  `tfsdk:"index"`
}

func (l verityAsPathAccessListListsModel) GetIndex() types.Int64 {
	return l.Index
}

type verityAsPathAccessListObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityAsPathAccessListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_as_path_access_list"
}

func (r *verityAsPathAccessListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityAsPathAccessListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity AS Path Access List",
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
		},
		Blocks: map[string]schema.Block{
			"lists": schema.ListNestedBlock{
				Description: "List of AS Path Access List entries",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable this AS Path Access List",
							Optional:    true,
						},
						"regular_expression": schema.StringAttribute{
							Description: "Regular Expression to match BGP Community Strings",
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
				Description: "Object properties for the AS Path Access List",
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

func (r *verityAsPathAccessListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityAsPathAccessListResourceModel
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
	asPathAccessListProps := &openapi.AspathaccesslistsPutRequestAsPathAccessListValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "PermitDeny", APIField: &asPathAccessListProps.PermitDeny, TFValue: plan.PermitDeny},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &asPathAccessListProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		asPathAccessListProps.ObjectProperties = &objProps
	}

	// Handle lists
	if len(plan.Lists) > 0 {
		lists := make([]openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner, len(plan.Lists))
		for i, item := range plan.Lists {
			list := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &list.Enable, TFValue: item.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RegularExpression", APIField: &list.RegularExpression, TFValue: item.RegularExpression},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &list.Index, TFValue: item.Index},
			})
			lists[i] = list
		}
		asPathAccessListProps.Lists = lists
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "as_path_access_list", name, *asPathAccessListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("AS Path Access List %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "as_path_access_lists")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityAsPathAccessListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityAsPathAccessListResourceModel
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

	asPathAccessListName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("as_path_access_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping AS Path Access List %s verification – trusting recent successful API operation", asPathAccessListName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching AS Path Access Lists for verification of %s", asPathAccessListName))

	type AsPathAccessListsResponse struct {
		AsPathAccessList map[string]interface{} `json:"as_path_access_list"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "as_path_access_lists", asPathAccessListName,
		func() (AsPathAccessListsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch AS Path Access Lists")
			respAPI, err := r.client.ASPathAccessListsAPI.AspathaccesslistsGet(ctx).Execute()
			if err != nil {
				return AsPathAccessListsResponse{}, fmt.Errorf("error reading AS Path Access Lists: %v", err)
			}
			defer respAPI.Body.Close()

			var res AsPathAccessListsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return AsPathAccessListsResponse{}, fmt.Errorf("failed to decode AS Path Access Lists response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d AS Path Access Lists", len(res.AsPathAccessList)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read AS Path Access List %s", asPathAccessListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for AS Path Access List with name: %s", asPathAccessListName))

	asPathAccessListData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.AsPathAccessList,
		asPathAccessListName,
		func(data interface{}) (string, bool) {
			if asPathAccessList, ok := data.(map[string]interface{}); ok {
				if name, ok := asPathAccessList["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("AS Path Access List with name '%s' not found in API response", asPathAccessListName))
		resp.State.RemoveResource(ctx)
		return
	}

	asPathAccessListMap, ok := asPathAccessListData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid AS Path Access List Data",
			fmt.Sprintf("AS Path Access List data is not in expected format for %s", asPathAccessListName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found AS Path Access List '%s' under API key '%s'", asPathAccessListName, actualAPIName))

	state.Name = utils.MapStringFromAPI(asPathAccessListMap["name"])

	// Handle object properties
	if objProps, ok := asPathAccessListMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityAsPathAccessListObjectPropertiesModel{
			{Notes: utils.MapStringFromAPI(objProps["notes"])},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"permit_deny": &state.PermitDeny,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(asPathAccessListMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(asPathAccessListMap[apiKey])
	}

	// Handle lists
	if listsData, ok := asPathAccessListMap["lists"].([]interface{}); ok && len(listsData) > 0 {
		var lists []verityAsPathAccessListListsModel

		for _, l := range listsData {
			listItem, ok := l.(map[string]interface{})
			if !ok {
				continue
			}

			listModel := verityAsPathAccessListListsModel{
				Enable:            utils.MapBoolFromAPI(listItem["enable"]),
				RegularExpression: utils.MapStringFromAPI(listItem["regular_expression"]),
				Index:             utils.MapInt64FromAPI(listItem["index"]),
			}

			lists = append(lists, listModel)
		}

		state.Lists = lists
	} else {
		state.Lists = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityAsPathAccessListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityAsPathAccessListResourceModel

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
	asPathAccessListProps := openapi.AspathaccesslistsPutRequestAsPathAccessListValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { asPathAccessListProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PermitDeny, state.PermitDeny, func(val *string) { asPathAccessListProps.PermitDeny = val }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(val *bool) { asPathAccessListProps.Enable = val }, &hasChanges)

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
			asPathAccessListProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle lists
	listsHandler := utils.IndexedItemHandler[verityAsPathAccessListListsModel, openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner]{
		CreateNew: func(planItem verityAsPathAccessListListsModel) openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner {
			item := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &item.Enable, TFValue: planItem.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RegularExpression", APIField: &item.RegularExpression, TFValue: planItem.RegularExpression},
			})

			return item
		},
		UpdateExisting: func(planItem verityAsPathAccessListListsModel, stateItem verityAsPathAccessListListsModel) (openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner, bool) {
			item := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			if !planItem.Enable.Equal(stateItem.Enable) {
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &item.Enable, TFValue: planItem.Enable},
				})
				fieldChanged = true
			}

			if !planItem.RegularExpression.Equal(stateItem.RegularExpression) {
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "RegularExpression", APIField: &item.RegularExpression, TFValue: planItem.RegularExpression},
				})
				fieldChanged = true
			}

			return item, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner {
			item := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: types.Int64Value(index)},
			})
			return item
		},
	}

	changedLists, listsChanged := utils.ProcessIndexedArrayUpdates(plan.Lists, state.Lists, listsHandler)
	if listsChanged {
		asPathAccessListProps.Lists = changedLists
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "as_path_access_list", name, asPathAccessListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("AS Path Access List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "as_path_access_lists")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityAsPathAccessListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityAsPathAccessListResourceModel
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "as_path_access_list", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("AS Path Access List %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "as_path_access_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityAsPathAccessListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
