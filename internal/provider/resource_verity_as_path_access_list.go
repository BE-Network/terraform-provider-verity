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

	if !plan.Enable.IsNull() {
		asPathAccessListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.PermitDeny.IsNull() {
		asPathAccessListProps.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
	}

	if len(plan.Lists) > 0 {
		lists := make([]openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner, len(plan.Lists))
		for i, listItem := range plan.Lists {
			lItem := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{}
			if !listItem.Index.IsNull() {
				lItem.Index = openapi.PtrInt32(int32(listItem.Index.ValueInt64()))
			}
			if !listItem.Enable.IsNull() {
				lItem.Enable = openapi.PtrBool(listItem.Enable.ValueBool())
			}
			if !listItem.RegularExpression.IsNull() {
				lItem.RegularExpression = openapi.PtrString(listItem.RegularExpression.ValueString())
			}
			lists[i] = lItem
		}
		asPathAccessListProps.Lists = lists
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			objProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objProps.Notes = nil
		}
		asPathAccessListProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "as_path_access_list", name, *asPathAccessListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for AS Path Access List creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create AS Path Access List %s", name))...,
		)
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

	var result AsPathAccessListsResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		asPathAccessListsData, fetchErr := getCachedResponse(ctx, r.provCtx, "aspathaccesslists", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch AS Path Access Lists")
			respAPI, err := r.client.ASPathAccessListsAPI.AspathaccesslistsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading AS Path Access Lists: %v", err)
			}
			defer respAPI.Body.Close()

			var res AsPathAccessListsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode AS Path Access Lists response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d AS Path Access Lists", len(res.AsPathAccessList)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch AS Path Access Lists on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = asPathAccessListsData.(AsPathAccessListsResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read AS Path Access List %s", asPathAccessListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for AS Path Access List with ID: %s", asPathAccessListName))
	var asPathAccessListData map[string]interface{}
	exists := false

	if data, ok := result.AsPathAccessList[asPathAccessListName].(map[string]interface{}); ok {
		asPathAccessListData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found AS Path Access List directly by ID: %s", asPathAccessListName))
	} else {
		for apiName, apal := range result.AsPathAccessList {
			asPathAccessList, ok := apal.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := asPathAccessList["name"].(string); ok && name == asPathAccessListName {
				asPathAccessListData = asPathAccessList
				asPathAccessListName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found AS Path Access List with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("AS Path Access List with ID '%s' not found in API response", asPathAccessListName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", asPathAccessListData["name"]))

	if enable, ok := asPathAccessListData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if permitDeny, ok := asPathAccessListData["permit_deny"].(string); ok && permitDeny != "" {
		state.PermitDeny = types.StringValue(permitDeny)
	} else {
		state.PermitDeny = types.StringNull()
	}

	if listsData, ok := asPathAccessListData["lists"].([]interface{}); ok {
		var lists []verityAsPathAccessListListsModel
		for _, listItemData := range listsData {
			if listItem, ok := listItemData.(map[string]interface{}); ok {
				var list verityAsPathAccessListListsModel

				if enable, ok := listItem["enable"].(bool); ok {
					list.Enable = types.BoolValue(enable)
				} else {
					list.Enable = types.BoolNull()
				}

				if regexpr, ok := listItem["regular_expression"].(string); ok {
					list.RegularExpression = types.StringValue(regexpr)
				} else {
					list.RegularExpression = types.StringNull()
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

	if objProps, ok := asPathAccessListData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objProps["notes"].(string); ok {
			state.ObjectProperties = []verityAsPathAccessListObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityAsPathAccessListObjectPropertiesModel{
				{Notes: types.StringValue("")},
			}
		}
	} else {
		state.ObjectProperties = nil
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

	if !plan.Enable.Equal(state.Enable) {
		asPathAccessListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.PermitDeny.Equal(state.PermitDeny) {
		asPathAccessListProps.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
		hasChanges = true
	}

	stateListsByIndex := make(map[int64]verityAsPathAccessListListsModel)
	for _, item := range state.Lists {
		if !item.Index.IsNull() {
			idx := item.Index.ValueInt64()
			stateListsByIndex[idx] = item
		}
	}

	var changedLists []openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner
	listsChanged := false

	// Process plan items (CREATE and UPDATE operations)
	for _, planItem := range plan.Lists {
		if planItem.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planItem.Index.ValueInt64()
		stateItem, exists := stateListsByIndex[idx]

		if !exists {
			// CREATE: new item, include all fields
			newItem := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planItem.Enable.IsNull() {
				newItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			} else {
				newItem.Enable = openapi.PtrBool(false)
			}

			if !planItem.RegularExpression.IsNull() && planItem.RegularExpression.ValueString() != "" {
				newItem.RegularExpression = openapi.PtrString(planItem.RegularExpression.ValueString())
			} else {
				newItem.RegularExpression = openapi.PtrString("")
			}

			changedLists = append(changedLists, newItem)
			listsChanged = true
			continue
		}

		// UPDATE: existing item, check which fields changed
		updateItem := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planItem.Enable.Equal(stateItem.Enable) {
			if !planItem.Enable.IsNull() {
				updateItem.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			} else {
				updateItem.Enable = openapi.PtrBool(false)
			}
			fieldChanged = true
		}

		if !planItem.RegularExpression.Equal(stateItem.RegularExpression) {
			if !planItem.RegularExpression.IsNull() && planItem.RegularExpression.ValueString() != "" {
				updateItem.RegularExpression = openapi.PtrString(planItem.RegularExpression.ValueString())
			} else {
				updateItem.RegularExpression = openapi.PtrString("")
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
			deletedItem := openapi.AspathaccesslistsPutRequestAsPathAccessListValueListsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedLists = append(changedLists, deletedItem)
			listsChanged = true
		}
	}

	if listsChanged && len(changedLists) > 0 {
		asPathAccessListProps.Lists = changedLists
		hasChanges = true
	}

	objPropsChanged := false
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		if !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objPropsChanged = true
		}
	} else if len(plan.ObjectProperties) > 0 || len(state.ObjectProperties) > 0 {
		objPropsChanged = true
	}

	if objPropsChanged {
		if len(plan.ObjectProperties) > 0 {
			objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				objProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				objProps.Notes = nil
			}
			asPathAccessListProps.ObjectProperties = &objProps
		}
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "as_path_access_list", name, asPathAccessListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for AS Path Access List update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update AS Path Access List %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "as_path_access_list", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for AS Path Access List deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete AS Path Access List %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("AS Path Access List %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "as_path_access_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityAsPathAccessListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
