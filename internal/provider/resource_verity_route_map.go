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
	_ resource.Resource                = &verityRouteMapResource{}
	_ resource.ResourceWithConfigure   = &verityRouteMapResource{}
	_ resource.ResourceWithImportState = &verityRouteMapResource{}
	_ resource.ResourceWithModifyPlan  = &verityRouteMapResource{}
)

const routeMapResourceType = "routemaps"

func NewVerityRouteMapResource() resource.Resource {
	return &verityRouteMapResource{}
}

type verityRouteMapResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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

func (m verityRouteMapClausesModel) GetIndex() types.Int64 {
	return m.Index
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
				Computed:    true,
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
							Computed:    true,
						},
						"route_map_clause": schema.StringAttribute{
							Description: "Route Map Clause is a collection match and set rules",
							Optional:    true,
							Computed:    true,
						},
						"route_map_clause_ref_type_": schema.StringAttribute{
							Description: "Object type for route_map_clause field",
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
				Description: "Object properties for the Route Map",
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
	routeMapReq := &openapi.RoutemapsPutRequestRouteMapValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &routeMapReq.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		routeMapReq.ObjectProperties = &objProps
	}

	// Handle route_map_clauses list
	if len(plan.RouteMapClauses) > 0 {
		routeMapClausesList := make([]openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner, len(plan.RouteMapClauses))
		for i, clause := range plan.RouteMapClauses {
			clauseProps := openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &clauseProps.Enable, TFValue: clause.Enable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RouteMapClause", APIField: &clauseProps.RouteMapClause, TFValue: clause.RouteMapClause},
				{FieldName: "RouteMapClauseRefType", APIField: &clauseProps.RouteMapClauseRefType, TFValue: clause.RouteMapClauseRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &clauseProps.Index, TFValue: clause.Index},
			})

			routeMapClausesList[i] = clauseProps
		}
		routeMapReq.RouteMapClauses = routeMapClausesList
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "route_map", name, *routeMapReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_maps")

	var minState verityRouteMapResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if routeMapData, exists := bulkMgr.GetResourceResponse("route_map", name); exists {
			state := populateRouteMapState(ctx, minState, routeMapData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if routeMapData, exists := r.bulkOpsMgr.GetResourceResponse("route_map", name); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached route_map data for %s from recent operation", name))
			state = populateRouteMapState(ctx, state, routeMapData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("route_map") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Route Map %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Route Map for verification of %s", name))

	type RouteMapResponse struct {
		RouteMap map[string]interface{} `json:"route_map"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "route_maps", name,
		func() (RouteMapResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Route Maps")
			respAPI, err := r.client.RouteMapsAPI.RoutemapsGet(ctx).Execute()
			if err != nil {
				return RouteMapResponse{}, fmt.Errorf("error reading Route Map: %v", err)
			}
			defer respAPI.Body.Close()

			var res RouteMapResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return RouteMapResponse{}, fmt.Errorf("failed to decode Route Maps response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Route Maps from API", len(res.RouteMap)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Route Map %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Route Map with name: %s", name))

	routeMapData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.RouteMap,
		name,
		func(data interface{}) (string, bool) {
			if routeMapMap, ok := data.(map[string]interface{}); ok {
				if name, ok := routeMapMap["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Route Map with name '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	routeMapMap, ok := (interface{}(routeMapData)).(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Route Map Data",
			fmt.Sprintf("Route Map data is not in expected format for %s", name),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Route Map '%s' under API key '%s'", name, actualAPIName))

	state = populateRouteMapState(ctx, state, routeMapMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { routeMapProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { routeMapProps.Enable = v }, &hasChanges)

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
			routeMapProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle route_map_clauses
	changedRouteMapClauses, routeMapClausesChanged := utils.ProcessIndexedArrayUpdates(plan.RouteMapClauses, state.RouteMapClauses,
		utils.IndexedItemHandler[verityRouteMapClausesModel, openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner]{
			CreateNew: func(planItem verityRouteMapClausesModel) openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner {
				newClause := openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &newClause.Enable, TFValue: planItem.Enable},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "RouteMapClause", APIField: &newClause.RouteMapClause, TFValue: planItem.RouteMapClause},
					{FieldName: "RouteMapClauseRefType", APIField: &newClause.RouteMapClauseRefType, TFValue: planItem.RouteMapClauseRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newClause.Index, TFValue: planItem.Index},
				})

				return newClause
			},
			UpdateExisting: func(planItem verityRouteMapClausesModel, stateItem verityRouteMapClausesModel) (openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner, bool) {
				updateClause := openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateClause.Enable = v }, &fieldChanged)

				// Handle route_map_clause and route_map_clause_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.RouteMapClause, stateItem.RouteMapClause, planItem.RouteMapClauseRefType, stateItem.RouteMapClauseRefType,
					func(v *string) { updateClause.RouteMapClause = v },
					func(v *string) { updateClause.RouteMapClauseRefType = v },
					"route_map_clause", "route_map_clause_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updateClause, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateClause.Index = v }, &fieldChanged)

				return updateClause, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner {
				return openapi.RoutemapsPutRequestRouteMapValueRouteMapClausesInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if routeMapClausesChanged {
		routeMapProps.RouteMapClauses = changedRouteMapClauses
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "route_map", name, routeMapProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_maps")

	var minState verityRouteMapResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if routeMapData, exists := bulkMgr.GetResourceResponse("route_map", name); exists {
			newState := populateRouteMapState(ctx, minState, routeMapData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "route_map", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_maps")
	resp.State.RemoveResource(ctx)
}

func (r *verityRouteMapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateRouteMapState(ctx context.Context, state verityRouteMapResourceModel, data map[string]interface{}, mode string) verityRouteMapResourceModel {
	const resourceType = routeMapResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityRouteMapObjectPropertiesModel{
				Notes: utils.MapStringWithModeNested(objProps, "notes", resourceType, "object_properties.notes", mode),
			}
			state.ObjectProperties = []verityRouteMapObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle route_map_clauses block
	if utils.FieldAppliesToMode(resourceType, "route_map_clauses", mode) {
		if routeMapClausesData, ok := data["route_map_clauses"].([]interface{}); ok && len(routeMapClausesData) > 0 {
			var routeMapClauses []verityRouteMapClausesModel

			for _, clauseInterface := range routeMapClausesData {
				clause, ok := clauseInterface.(map[string]interface{})
				if !ok {
					continue
				}

				clauseModel := verityRouteMapClausesModel{
					Enable:                utils.MapBoolWithModeNested(clause, "enable", resourceType, "route_map_clauses.enable", mode),
					RouteMapClause:        utils.MapStringWithModeNested(clause, "route_map_clause", resourceType, "route_map_clauses.route_map_clause", mode),
					RouteMapClauseRefType: utils.MapStringWithModeNested(clause, "route_map_clause_ref_type_", resourceType, "route_map_clauses.route_map_clause_ref_type_", mode),
					Index:                 utils.MapInt64WithModeNested(clause, "index", resourceType, "route_map_clauses.index", mode),
				}

				routeMapClauses = append(routeMapClauses, clauseModel)
			}

			state.RouteMapClauses = routeMapClauses
		} else {
			state.RouteMapClauses = nil
		}
	} else {
		state.RouteMapClauses = nil
	}

	return state
}

func (r *verityRouteMapResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityRouteMapResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = routeMapResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
	)
}
