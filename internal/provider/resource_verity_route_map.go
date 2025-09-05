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
	_ resource.Resource                = &verityRouteMapResource{}
	_ resource.ResourceWithConfigure   = &verityRouteMapResource{}
	_ resource.ResourceWithImportState = &verityRouteMapResource{}
)

func NewVerityRouteMapResource() resource.Resource {
	return &verityRouteMapResource{}
}

type verityRouteMapResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityRouteMapResourceModel struct {
	Name             types.String                          `tfsdk:"name"`
	Enable           types.Bool                            `tfsdk:"enable"`
	RouteMapClauses  []verityRouteMapClausesModel          `tfsdk:"route_map_clauses"`
	ObjectProperties []verityRouteMapObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityRouteMapClausesModel struct {
	Enable                types.Bool   `tfsdk:"enable"`
	RouteMapClause        types.String `tfsdk:"route_map_clause"`
	RouteMapClauseRefType types.String `tfsdk:"route_map_clause_ref_type_"`
	Index                 types.Int64  `tfsdk:"index"`
}

type verityRouteMapObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityRouteMapResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_map"
}

func (r *verityRouteMapResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityRouteMapResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Route Map",
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
			"route_map_clauses": schema.ListNestedBlock{
				Description: "List of Route Map Clauses",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"route_map_clause": schema.StringAttribute{
							Description: "Route Map Clause is a collection match and set rules",
							Optional:    true,
						},
						"route_map_clause_ref_type_": schema.StringAttribute{
							Description: "Object type for route_map_clause field",
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
				Description: "Object properties for the Route Map",
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

func (r *verityRouteMapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityRouteMapResourceModel
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
	routeMapProps := openapi.RoutemapsPutRequestRouteMapValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		routeMapProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	// Handle route_map_clauses list
	if len(plan.RouteMapClauses) > 0 {
		routeMapClausesList := make([]openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner, len(plan.RouteMapClauses))
		for i, clause := range plan.RouteMapClauses {
			clauseProps := openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{}

			if !clause.Enable.IsNull() {
				clauseProps.Enable = openapi.PtrBool(clause.Enable.ValueBool())
			}

			if !clause.RouteMapClause.IsNull() && clause.RouteMapClause.ValueString() != "" {
				clauseProps.RouteMapClause = openapi.PtrString(clause.RouteMapClause.ValueString())
			}

			if !clause.RouteMapClauseRefType.IsNull() && clause.RouteMapClauseRefType.ValueString() != "" {
				clauseProps.RouteMapClauseRefType = openapi.PtrString(clause.RouteMapClauseRefType.ValueString())
			}

			if !clause.Index.IsNull() {
				clauseProps.Index = openapi.PtrInt32(int32(clause.Index.ValueInt64()))
			}

			routeMapClausesList[i] = clauseProps
		}
		routeMapProps.RouteMapClauses = routeMapClausesList
	}

	// Handle object properties following Gateway pattern
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		routeMapObjProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			routeMapObjProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			routeMapObjProps.Notes = nil
		}
		routeMapProps.ObjectProperties = &routeMapObjProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "route_map", name, routeMapProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Route Map create operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Route Map %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Route Map %s create operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_maps")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityRouteMapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityRouteMapResourceModel
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("route_map") {
		tflog.Info(ctx, fmt.Sprintf("Skipping route map %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching route maps for verification of %s", name))

	type RouteMapResponse struct {
		RouteMap map[string]interface{} `json:"route_map"`
	}

	var result RouteMapResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		routeMapsData, fetchErr := getCachedResponse(ctx, r.provCtx, "route_maps", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch route maps")
			respAPI, err := r.client.RouteMapsAPI.RoutemapsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading route maps: %v", err)
			}
			defer respAPI.Body.Close()

			var res RouteMapResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode route maps response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d route maps", len(res.RouteMap)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch route maps on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = routeMapsData.(RouteMapResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Route Map %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for route map with ID: %s", name))
	var routeMapData map[string]interface{}
	exists := false

	if data, ok := result.RouteMap[name].(map[string]interface{}); ok {
		routeMapData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found route map directly by ID: %s", name))
	} else {
		for apiName, r := range result.RouteMap {
			routeMap, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			if routeMapName, ok := routeMap["name"].(string); ok && routeMapName == name {
				routeMapData = routeMap
				name = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found route map with name '%s' under API key '%s'", routeMapName, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Route Map with ID '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", routeMapData["name"]))

	if enable, ok := routeMapData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	// Handle route_map_clauses
	if routeMapClausesData, ok := routeMapData["route_map_clauses"].([]interface{}); ok && len(routeMapClausesData) > 0 {
		var routeMapClauses []verityRouteMapClausesModel

		for _, clauseInterface := range routeMapClausesData {
			clause, ok := clauseInterface.(map[string]interface{})
			if !ok {
				continue
			}

			clauseModel := verityRouteMapClausesModel{}

			if enable, ok := clause["enable"].(bool); ok {
				clauseModel.Enable = types.BoolValue(enable)
			} else {
				clauseModel.Enable = types.BoolNull()
			}

			if routeMapClause, ok := clause["route_map_clause"].(string); ok {
				clauseModel.RouteMapClause = types.StringValue(routeMapClause)
			} else {
				clauseModel.RouteMapClause = types.StringNull()
			}

			if routeMapClauseRefType, ok := clause["route_map_clause_ref_type_"].(string); ok {
				clauseModel.RouteMapClauseRefType = types.StringValue(routeMapClauseRefType)
			} else {
				clauseModel.RouteMapClauseRefType = types.StringNull()
			}

			if indexValue, exists := clause["index"]; exists && indexValue != nil {
				switch v := indexValue.(type) {
				case int:
					clauseModel.Index = types.Int64Value(int64(v))
				case int32:
					clauseModel.Index = types.Int64Value(int64(v))
				case int64:
					clauseModel.Index = types.Int64Value(v)
				case float64:
					clauseModel.Index = types.Int64Value(int64(v))
				default:
					clauseModel.Index = types.Int64Null()
				}
			} else {
				clauseModel.Index = types.Int64Null()
			}

			routeMapClauses = append(routeMapClauses, clauseModel)
		}

		state.RouteMapClauses = routeMapClauses
	} else {
		state.RouteMapClauses = nil
	}

	// Only set object_properties if it exists in the API response
	if objProps, ok := routeMapData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objProps["notes"].(string); ok {
			state.ObjectProperties = []verityRouteMapObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityRouteMapObjectPropertiesModel{
				{Notes: types.StringNull()},
			}
		}
	} else {
		state.ObjectProperties = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *verityRouteMapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityRouteMapResourceModel

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
	routeMapProps := openapi.RoutemapsPutRequestRouteMapValue{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		routeMapProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		routeMapProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	// Handle route_map_clauses changes using Gateway pattern
	oldRouteMapClausesByIndex := make(map[int64]verityRouteMapClausesModel)
	for _, clause := range state.RouteMapClauses {
		if !clause.Index.IsNull() {
			oldRouteMapClausesByIndex[clause.Index.ValueInt64()] = clause
		}
	}

	var changedRouteMapClauses []openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner
	routeMapClausesChanged := false

	for _, clause := range plan.RouteMapClauses {
		if clause.Index.IsNull() {
			continue
		}

		index := clause.Index.ValueInt64()
		oldClause, exists := oldRouteMapClausesByIndex[index]

		if !exists {
			// new route map clause, include all fields
			clauseProps := openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			if !clause.Enable.IsNull() {
				clauseProps.Enable = openapi.PtrBool(clause.Enable.ValueBool())
			} else {
				clauseProps.Enable = openapi.PtrBool(false)
			}

			if !clause.RouteMapClause.IsNull() {
				clauseProps.RouteMapClause = openapi.PtrString(clause.RouteMapClause.ValueString())
			} else {
				clauseProps.RouteMapClause = openapi.PtrString("")
			}

			if !clause.RouteMapClauseRefType.IsNull() {
				clauseProps.RouteMapClauseRefType = openapi.PtrString(clause.RouteMapClauseRefType.ValueString())
			} else {
				clauseProps.RouteMapClauseRefType = openapi.PtrString("")
			}

			changedRouteMapClauses = append(changedRouteMapClauses, clauseProps)
			routeMapClausesChanged = true
			continue
		}

		// existing route map clause, check which fields changed
		clauseProps := openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !clause.Enable.Equal(oldClause.Enable) {
			clauseProps.Enable = openapi.PtrBool(clause.Enable.ValueBool())
			fieldChanged = true
		}

		if !clause.RouteMapClause.Equal(oldClause.RouteMapClause) {
			if !clause.RouteMapClause.IsNull() {
				clauseProps.RouteMapClause = openapi.PtrString(clause.RouteMapClause.ValueString())
			} else {
				clauseProps.RouteMapClause = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !clause.RouteMapClauseRefType.Equal(oldClause.RouteMapClauseRefType) {
			if !clause.RouteMapClauseRefType.IsNull() {
				clauseProps.RouteMapClauseRefType = openapi.PtrString(clause.RouteMapClauseRefType.ValueString())
			} else {
				clauseProps.RouteMapClauseRefType = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedRouteMapClauses = append(changedRouteMapClauses, clauseProps)
			routeMapClausesChanged = true
		}
	}

	for idx := range oldRouteMapClausesByIndex {
		found := false
		for _, clause := range plan.RouteMapClauses {
			if !clause.Index.IsNull() && clause.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// route map clause removed - include only the index for deletion
			deletedClause := openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedRouteMapClauses = append(changedRouteMapClauses, deletedClause)
			routeMapClausesChanged = true
		}
	}

	if routeMapClausesChanged && len(changedRouteMapClauses) > 0 {
		routeMapProps.RouteMapClauses = changedRouteMapClauses
		hasChanges = true
	}

	// Handle object_properties changes following Gateway pattern
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			routeMapObjProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				routeMapObjProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				routeMapObjProps.Notes = nil
			}
			routeMapProps.ObjectProperties = &routeMapObjProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "route_map", name, routeMapProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Route Map update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Route Map %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Route Map %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_maps")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityRouteMapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityRouteMapResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "route_map", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Route Map delete operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Route Map %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Route Map %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_maps")
	resp.State.RemoveResource(ctx)
}

func (r *verityRouteMapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
